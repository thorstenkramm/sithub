package livefeed

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/auth"
	"github.com/thorstenkramm/sithub/internal/config"
	"github.com/thorstenkramm/sithub/internal/middleware"
	"github.com/thorstenkramm/sithub/internal/notifications"
)

func TestHandlerBroadcastsEventToAuthorizedClient(t *testing.T) {
	t.Parallel()

	hub := NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go hub.Run(ctx)

	svc := newTestAuthService(t)
	e := echo.New()
	e.Use(middleware.LoadUser(svc))
	e.GET("/api/v1/live", Handler(hub), middleware.RequireAuth(svc))

	srv := httptest.NewServer(e)
	defer srv.Close()

	wsURL := strings.Replace(srv.URL, "http://", "ws://", 1) + "/api/v1/live"
	headers := http.Header{}
	headers.Add("Cookie", testUserCookie(t, svc, &auth.User{
		ID:          "observer-1",
		Name:        "Observer",
		AuthSource:  "internal",
		IsPermitted: true,
	}).String())

	dialer := *websocket.DefaultDialer
	dialer.HandshakeTimeout = 2 * time.Second
	conn, resp, err := dialer.Dial(wsURL, headers)
	require.NoError(t, err)
	if resp != nil {
		require.NoError(t, resp.Body.Close())
	}
	t.Cleanup(func() { _ = conn.Close() }) //nolint:errcheck // best-effort cleanup

	// The hub registers the client asynchronously; events sent before the
	// registration lands are silently dropped. Keep sending warmup events from
	// a background ticker and block on a SINGLE read with a generous deadline.
	// gorilla/websocket treats a timed-out read as FATAL for the connection,
	// so the previous poll (repeated ReadJSON with a 50 ms deadline inside
	// Eventually) poisoned the connection on its first miss and could never
	// recover — the root cause of the -race CI flake.
	stopWarmup := make(chan struct{})
	warmupDone := make(chan struct{})
	go func() {
		defer close(warmupDone)
		ticker := time.NewTicker(20 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-stopWarmup:
				return
			case <-ticker.C:
				hub.NotifyAsync(&notifications.BookingEvent{
					Event:       notifications.EventBookingCreated,
					BookingID:   "warmup",
					BookingDate: "2026-05-11",
					Timestamp:   "2026-05-10T12:00:00Z",
				})
			}
		}
	}()
	_ = conn.SetReadDeadline(time.Now().Add(10 * time.Second)) //nolint:errcheck // deadline-only
	var warm Event
	require.NoError(t, conn.ReadJSON(&warm), "websocket client never registered with hub")
	close(stopWarmup)
	<-warmupDone

	// Everything enqueued from here on comes after any leftover warmups, so
	// skipping warmups below is guaranteed to terminate at "b1".
	hub.NotifyAsync(&notifications.BookingEvent{
		Event:          notifications.EventBookingCreated,
		BookingID:      "b1",
		ItemID:         "desk-1",
		UserID:         "owner-1",
		BookedByUserID: "admin-1",
		BookingDate:    "2026-05-11",
		Timestamp:      "2026-05-10T12:00:00Z",
	})

	_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second)) //nolint:errcheck // deadline-only
	var got Event
	for {
		require.NoError(t, conn.ReadJSON(&got))
		if got.BookingID != "warmup" {
			break
		}
	}
	assert.Equal(t, EventBookingCreated, got.Type)
	assert.Equal(t, "b1", got.BookingID)
	assert.Equal(t, "desk-1", got.ItemID)
	assert.Equal(t, "admin-1", got.UserID)
	assert.Equal(t, "2026-05-11", got.BookingDate)
}

func TestHandlerUpgradeFailureReturnsJSONAPIBadRequest(t *testing.T) {
	t.Parallel()

	hub := NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go hub.Run(ctx)

	svc := newTestAuthService(t)
	e := echo.New()
	e.Use(middleware.LoadUser(svc))
	e.GET("/api/v1/live", Handler(hub), middleware.RequireAuth(svc))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/live", http.NoBody)
	req.AddCookie(testUserCookie(t, svc, &auth.User{
		ID:          "user-1",
		Name:        "Test User",
		AuthSource:  "internal",
		IsPermitted: true,
	}))
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, api.JSONAPIContentType, rec.Header().Get(echo.HeaderContentType))

	var resp api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Errors, 1)
	assert.Equal(t, "bad_request", resp.Errors[0].Code)
	assert.Equal(t, "WebSocket upgrade failed", resp.Errors[0].Detail)
}

func newTestAuthService(t *testing.T) *auth.Service {
	t.Helper()

	svc, err := auth.NewService(&config.Config{EntraID: config.EntraIDConfig{
		AuthorizeURL: "https://example.com/auth",
		TokenURL:     "https://example.com/token",
		RedirectURI:  "https://example.com/callback",
		ClientID:     "client",
		ClientSecret: "secret",
	}}, nil)
	require.NoError(t, err)
	return svc
}

func testUserCookie(t *testing.T, svc *auth.Service, user *auth.User) *http.Cookie {
	t.Helper()

	encodedUser, err := svc.EncodeUser(user)
	require.NoError(t, err)

	return &http.Cookie{Name: "sithub_user", Value: encodedUser}
}
