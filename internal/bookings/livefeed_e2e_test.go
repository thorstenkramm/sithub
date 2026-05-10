package bookings

import (
	"bytes"
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
	"github.com/thorstenkramm/sithub/internal/livefeed"
	"github.com/thorstenkramm/sithub/internal/middleware"
	"github.com/thorstenkramm/sithub/internal/notifications"
)

func TestBookingFlowBroadcastsCreateAndCancelEventsWithActorIdentity(t *testing.T) {
	t.Parallel()

	cfg := testAreasConfig()
	store := setupTestStore(t)
	seedTestUserRecord(t, store, "admin-1", "admin@test.local", "Admin User")
	seedTestUserRecord(t, store, "allowed-user", "allowed@test.local", "Allowed User")

	hub := livefeed.NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go hub.Run(ctx)

	notifier := notifications.MultiNotifier{hub}
	svc := newLiveFeedTestAuthService(t)
	e := echo.New()
	e.Use(middleware.LoadUser(svc))
	e.GET("/api/v1/live", livefeed.Handler(hub), middleware.RequireAuth(svc))
	e.POST("/api/v1/bookings", CreateHandler(cfg, store, notifier), middleware.RequireAuth(svc))
	e.DELETE("/api/v1/bookings/:id", DeleteHandler(store, notifier), middleware.RequireAuth(svc))

	srv := httptest.NewServer(e)
	defer srv.Close()

	adminCookie := liveFeedTestUserCookie(t, svc, &auth.User{
		ID:          "admin-1",
		Name:        "Admin User",
		Email:       "admin@test.local",
		AuthSource:  "internal",
		IsAdmin:     true,
		IsPermitted: true,
	})

	wsURL := strings.Replace(srv.URL, "http://", "ws://", 1) + "/api/v1/live"
	headers := http.Header{}
	headers.Add("Cookie", adminCookie.String())

	dialer := *websocket.DefaultDialer
	dialer.HandshakeTimeout = 2 * time.Second
	conn, resp, err := dialer.Dial(wsURL, headers)
	require.NoError(t, err)
	if resp != nil {
		require.NoError(t, resp.Body.Close())
	}
	t.Cleanup(func() { _ = conn.Close() }) //nolint:errcheck // best-effort cleanup

	require.Eventually(t, func() bool {
		hub.NotifyAsync(&notifications.BookingEvent{
			Event:       notifications.EventBookingCreated,
			BookingID:   "warmup",
			BookingDate: "2026-05-11",
			Timestamp:   "2026-05-10T12:00:00Z",
		})
		_ = conn.SetReadDeadline(time.Now().Add(50 * time.Millisecond)) //nolint:errcheck // deadline-only
		var ev livefeed.Event
		return conn.ReadJSON(&ev) == nil
	}, time.Second, 20*time.Millisecond, "websocket client never registered with hub")

	bookingDate := time.Now().UTC().AddDate(0, 0, 1).Format(time.DateOnly)
	createReqBody := `{"data":{"type":"bookings","attributes":{"item_id":"desk-1","booking_date":"` +
		bookingDate +
		`","for_user_id":"allowed-user"}}}`

	createReq, err := http.NewRequest(
		http.MethodPost,
		srv.URL+"/api/v1/bookings",
		bytes.NewBufferString(createReqBody),
	)
	require.NoError(t, err)
	createReq.Header.Set(echo.HeaderContentType, api.JSONAPIContentType)
	createReq.AddCookie(adminCookie)

	createResp, err := srv.Client().Do(createReq)
	require.NoError(t, err)
	defer createResp.Body.Close() //nolint:errcheck // best-effort cleanup
	require.Equal(t, http.StatusCreated, createResp.StatusCode)

	var created api.SingleResponse
	require.NoError(t, json.NewDecoder(createResp.Body).Decode(&created))
	bookingID := created.Data.ID
	require.NotEmpty(t, bookingID)

	_ = conn.SetReadDeadline(time.Now().Add(time.Second)) //nolint:errcheck // deadline-only
	var createdEvent livefeed.Event
	require.NoError(t, conn.ReadJSON(&createdEvent))
	assert.Equal(t, livefeed.EventBookingCreated, createdEvent.Type)
	assert.Equal(t, bookingID, createdEvent.BookingID)
	assert.Equal(t, "desk-1", createdEvent.ItemID)
	assert.Equal(t, "admin-1", createdEvent.UserID)
	assert.Equal(t, bookingDate, createdEvent.BookingDate)

	deleteReq, err := http.NewRequest(
		http.MethodDelete,
		srv.URL+"/api/v1/bookings/"+bookingID,
		http.NoBody,
	)
	require.NoError(t, err)
	deleteReq.AddCookie(adminCookie)

	deleteResp, err := srv.Client().Do(deleteReq)
	require.NoError(t, err)
	defer deleteResp.Body.Close() //nolint:errcheck // best-effort cleanup
	require.Equal(t, http.StatusNoContent, deleteResp.StatusCode)

	_ = conn.SetReadDeadline(time.Now().Add(time.Second)) //nolint:errcheck // deadline-only
	var canceledEvent livefeed.Event
	require.NoError(t, conn.ReadJSON(&canceledEvent))
	assert.Equal(t, livefeed.EventBookingCanceled, canceledEvent.Type)
	assert.Equal(t, bookingID, canceledEvent.BookingID)
	assert.Equal(t, "desk-1", canceledEvent.ItemID)
	assert.Equal(t, "admin-1", canceledEvent.UserID)
	assert.Equal(t, bookingDate, canceledEvent.BookingDate)
}

func newLiveFeedTestAuthService(t *testing.T) *auth.Service {
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

func liveFeedTestUserCookie(t *testing.T, svc *auth.Service, user *auth.User) *http.Cookie {
	t.Helper()

	encodedUser, err := svc.EncodeUser(user)
	require.NoError(t, err)

	return &http.Cookie{Name: "sithub_user", Value: encodedUser}
}
