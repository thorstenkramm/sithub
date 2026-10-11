package startup

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thorstenkramm/sithub/internal/apikeys"
	"github.com/thorstenkramm/sithub/internal/auth"
	"github.com/thorstenkramm/sithub/internal/livefeed"
	"github.com/thorstenkramm/sithub/internal/notifications"
)

func keyTestRouter(t *testing.T) (*echo.Echo, *auth.Service, string) {
	t.Helper()
	db := setupStartupTestStore(t)
	svc := newTestAuthService(t)
	e := echo.New()
	registerAuthMiddleware(e, svc, db)
	registerRoutes(e, svc, testAreasConfig(), t.TempDir(), t.TempDir(), db,
		notifications.NewNotifier(""), livefeed.NewHub(), nil, "test")
	_, raw, err := apikeys.Create(context.Background(), db, time.Now(), apikeys.CreateInput{
		Name: "Finance", Scope: apikeys.ScopeReadBookings, ExpiresInMonths: 3,
	})
	require.NoError(t, err)
	return e, svc, raw
}

func TestBearerRouteSweep(t *testing.T) {
	t.Parallel()
	e, svc, raw := keyTestRouter(t)
	admin := testUserCookie(t, svc, &auth.User{ID: "admin", AuthSource: "internal", IsAdmin: true, IsPermitted: true})
	allow := map[string]bool{
		"GET /api/v1/bookings/export": true, "GET /api/v1/bookings/export.csv": true,
		"GET /api/v1/avatars/:user_id": true,
		"GET /api/v1/ping":             true, "POST /api/v1/auth/login": true,
		"POST /api/v1/auth/logout": true, "GET /api/v1/auth/providers": true,
	}
	registered := map[string]bool{}
	checked := 0
	for _, route := range e.Routes() {
		registered[route.Method+" "+route.Path] = true
		if !strings.HasPrefix(route.Path, "/api/") || allow[route.Method+" "+route.Path] {
			continue
		}
		checked++
		t.Run(route.Method+" "+route.Path, func(t *testing.T) {
			req := httptest.NewRequest(route.Method, route.Path, http.NoBody)
			req.Header.Set("Authorization", "Bearer "+raw)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusForbidden, rec.Code, "every non-public, non-export API route denies keys")
			assert.Contains(t, rec.Body.String(), "API key not permitted for this endpoint",
				"the 403 must come from the API-key default-deny gate")

			req = httptest.NewRequest(route.Method, route.Path, http.NoBody)
			req.AddCookie(admin)
			req.Header.Set("Authorization", "Bearer "+uuid.NewString())
			rec = httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusUnauthorized, rec.Code, "an unknown bearer key never falls back to the cookie")
		})
	}
	for route := range allow {
		assert.True(t, registered[route], "allowlisted route %q must exist", route)
	}
	assert.Greater(t, checked, 20, "the sweep must exercise the real registered API routes")
}

func TestBearerCookiePrecedence(t *testing.T) {
	t.Parallel()
	e, svc, raw := keyTestRouter(t)
	cookie := testUserCookie(t, svc, &auth.User{ID: "admin", AuthSource: "internal", IsAdmin: true, IsPermitted: true})
	for _, tc := range []struct {
		name, header, msg string
		status            int
	}{
		{"invalid bearer", "Bearer bad-key", "an invalid bearer key must not fall back to the admin cookie", 401},
		{"valid bearer", "Bearer " + raw, "a valid key wins over the admin cookie and is denied", 403},
		{"basic proxy", "Basic proxy", "non-bearer schemes keep cookie authentication", 200},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/users", http.NoBody)
			req.AddCookie(cookie)
			req.Header.Set("Authorization", tc.header)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			assert.Equal(t, tc.status, rec.Code, tc.msg)
		})
	}
}

const exportRange = "filter[from]=2026-10-01&filter[to]=2026-10-08"

func TestBearerExports(t *testing.T) {
	t.Parallel()
	e, _, raw := keyTestRouter(t)
	for _, tc := range []struct {
		path    string
		status  int
		content string
	}{
		{"/api/v1/bookings/export?" + exportRange + "&page[limit]=1&page[offset]=0", 200, `"meta"`},
		{"/api/v1/bookings/export.csv?" + exportRange + "&filter[area]=area-1", 200, "Name,Email"},
		{"/api/v1/bookings/export?filter[from]=2026-10-08&filter[to]=2026-10-01", 400, `"errors"`},
		{"/api/v1/bookings/export.csv", 400, `"errors"`},
	} {
		t.Run(tc.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, http.NoBody)
			req.Header.Set("Authorization", "Bearer "+raw)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			assert.Equal(t, tc.status, rec.Code, "a read-bookings key reaches the export handlers and their validation")
			assert.Contains(t, rec.Body.String(), tc.content)
		})
	}
}

func TestBearerAvatar(t *testing.T) {
	t.Parallel()
	db := setupStartupTestStore(t)
	svc := newTestAuthService(t)
	avatarsDir := t.TempDir()
	e := echo.New()
	registerAuthMiddleware(e, svc, db)
	registerRoutes(e, svc, testAreasConfig(), t.TempDir(), avatarsDir, db,
		notifications.NewNotifier(""), livefeed.NewHub(), nil, "test")
	_, raw, err := apikeys.Create(context.Background(), db, time.Now(), apikeys.CreateInput{
		Name: "Signage", Scope: apikeys.ScopeReadBookings, ExpiresInMonths: 3,
	})
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(avatarsDir, "user-1.png"), paddedPNG(t, 1024), 0o600))

	for _, tc := range []struct {
		name, path, header string
		status             int
	}{
		{"key reads an avatar", "/api/v1/avatars/user-1", "Bearer " + raw, http.StatusOK},
		{"key gets 404 for a missing avatar", "/api/v1/avatars/nobody", "Bearer " + raw, http.StatusNotFound},
		{"unknown key", "/api/v1/avatars/user-1", "Bearer " + uuid.NewString(), http.StatusUnauthorized},
		{"anonymous", "/api/v1/avatars/user-1", "", http.StatusUnauthorized},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, http.NoBody)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			assert.Equal(t, tc.status, rec.Code)
		})
	}

	// The same key still cannot change avatars.
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/me/avatar", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+raw)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}
