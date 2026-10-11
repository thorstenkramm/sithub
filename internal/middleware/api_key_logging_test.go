package middleware

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thorstenkramm/sithub/internal/apikeys"
	"github.com/thorstenkramm/sithub/internal/testutil/dbtest"
)

func TestAPIKeyLoggingRedaction(t *testing.T) {
	// This test changes the global logger; it must not run in parallel.
	var output bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug})))
	t.Cleanup(func() { slog.SetDefault(previous) })
	db := dbtest.New(t)
	now := time.Now()
	_, valid, err := apikeys.Create(context.Background(), db, now, apikeys.CreateInput{
		Name: "Finance", Scope: apikeys.ScopeReadBookings, ExpiresInMonths: 3,
	})
	require.NoError(t, err)
	_, expired, err := apikeys.Create(context.Background(), db, now.AddDate(-1, 0, 0), apikeys.CreateInput{
		Name: "Old", Scope: apikeys.ScopeReadBookings, ExpiresInMonths: 3,
	})
	require.NoError(t, err)
	keys := []string{valid, uuid.NewString(), expired, "malformed-secret"}
	for _, raw := range keys {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users?token="+raw, http.NoBody)
		req.Header.Set("Authorization", "Bearer "+raw)
		c := echo.New().NewContext(req, httptest.NewRecorder())
		c.SetPath("/api/v1/users")
		require.NoError(t, LoadAPIKey(db)(func(c echo.Context) error { return c.NoContent(200) })(c))
	}
	for _, raw := range keys {
		assert.NotContains(t, output.String(), raw)
		assert.NotContains(t, output.String(), raw[:10])
		assert.NotContains(t, output.String(), apikeys.HashKey(raw))
	}
	for _, reason := range []string{"malformed", "unknown", "expired"} {
		assert.Contains(t, output.String(), "reason="+reason)
	}
	assert.Contains(t, output.String(), "path=/api/v1/users")
	assert.Contains(t, output.String(), "remote_ip=")
}
