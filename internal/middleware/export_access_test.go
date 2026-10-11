package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thorstenkramm/sithub/internal/apikeys"
	"github.com/thorstenkramm/sithub/internal/auth"
	"github.com/thorstenkramm/sithub/internal/config"
)

func TestUserGatesDenyAPIKeys(t *testing.T) {
	t.Parallel()
	for name, mw := range map[string]echo.MiddlewareFunc{
		"auth": RequireAuth(nil), "admin": RequireAdmin(), "exporter": RequireDataExporter(),
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/api/v1/users", http.NoBody), rec)
			auth.SetAPIKeyInContext(c, &auth.APIKeyPrincipal{ID: "key", Scope: apikeys.ScopeReadBookings})
			require.NoError(t, mw(func(echo.Context) error { t.Fatal("key bypassed user gate"); return nil })(c))
			assert.Equal(t, http.StatusForbidden, rec.Code)
			assert.Contains(t, rec.Body.String(), "API key not permitted for this endpoint")
		})
	}
}

func TestRequireAuthOrScope(t *testing.T) {
	t.Parallel()
	svc, err := auth.NewService(&config.Config{}, nil)
	require.NoError(t, err)
	for _, tc := range []struct {
		name   string
		key    *auth.APIKeyPrincipal
		user   *auth.User
		status int
	}{
		{"key scope", &auth.APIKeyPrincipal{Scope: apikeys.ScopeReadBookings}, nil, 200},
		{"wrong scope", &auth.APIKeyPrincipal{Scope: "other"}, nil, 403},
		{"regular user", nil, &auth.User{AuthSource: "internal", IsPermitted: true}, 200},
		{"forbidden user", nil, &auth.User{AuthSource: "internal", IsPermitted: false}, 403},
		{"nobody", nil, nil, 401},
		{"key wins", &auth.APIKeyPrincipal{Scope: "other"}, &auth.User{AuthSource: "internal", IsPermitted: true}, 403},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/api/v1/avatars/u1", http.NoBody), rec)
			auth.SetAPIKeyInContext(c, tc.key)
			if tc.user != nil {
				c.Set("user", tc.user)
			}
			h := RequireAuthOrScope(svc, apikeys.ScopeReadBookings)(func(c echo.Context) error {
				return c.NoContent(http.StatusOK)
			})
			require.NoError(t, h(c))
			assert.Equal(t, tc.status, rec.Code)
			if tc.key != nil && tc.status == http.StatusForbidden {
				assert.Contains(t, rec.Header().Get("WWW-Authenticate"), `error="insufficient_scope"`)
			}
		})
	}
}

func TestRequireExportAccess(t *testing.T) {
	t.Parallel()
	svc, err := auth.NewService(&config.Config{}, nil)
	require.NoError(t, err)
	for _, tc := range []struct {
		name   string
		key    *auth.APIKeyPrincipal
		user   *auth.User
		status int
	}{
		{"key scope", &auth.APIKeyPrincipal{Scope: apikeys.ScopeReadBookings}, nil, 200},
		{"wrong scope", &auth.APIKeyPrincipal{Scope: "other"}, nil, 403},
		{"exporter", nil, &auth.User{AuthSource: "internal", IsPermitted: true, IsDataExporter: true}, 200},
		{"admin", nil, &auth.User{AuthSource: "internal", IsPermitted: true, IsAdmin: true}, 403},
		{"nobody", nil, nil, 401},
		{"key wins", &auth.APIKeyPrincipal{Scope: "other"}, &auth.User{IsDataExporter: true}, 403},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/api/v1/bookings/export", http.NoBody), rec)
			auth.SetAPIKeyInContext(c, tc.key)
			if tc.user != nil {
				c.Set("user", tc.user)
			}
			h := RequireExportAccess(svc, apikeys.ScopeReadBookings)(func(c echo.Context) error {
				return c.NoContent(http.StatusOK)
			})
			require.NoError(t, h(c))
			assert.Equal(t, tc.status, rec.Code)
			if tc.key != nil && tc.status == http.StatusForbidden {
				assert.Contains(t, rec.Header().Get("WWW-Authenticate"), `error="insufficient_scope"`,
					"a key outside its scope gets an RFC 6750 insufficient_scope challenge")
			}
		})
	}
}
