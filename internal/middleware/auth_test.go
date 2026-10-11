package middleware

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/auth"
	"github.com/thorstenkramm/sithub/internal/config"
	"github.com/thorstenkramm/sithub/internal/testutil/dbtest"
)

func TestRequireAuthBlocksMissingUser(t *testing.T) {
	svc := newTestAuthService(t)
	rec := runMiddleware(t, RequireAuth(svc), "/api/v1/me", nil)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}

	var resp api.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Errors) != 1 || resp.Errors[0].Code != "auth_required" {
		t.Fatalf("unexpected error response: %#v", resp.Errors)
	}
}

func TestRequireAuthAllowsUser(t *testing.T) {
	svc := newTestAuthService(t)
	rec := runMiddleware(t, RequireAuth(svc), "/api/v1/me", &auth.User{IsPermitted: true})

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestRequireAuthBlocksForbiddenUser(t *testing.T) {
	svc := newTestAuthService(t)
	rec := runMiddleware(t, RequireAuth(svc), "/api/v1/me", &auth.User{IsPermitted: false})

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}

	var resp api.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Errors) != 1 || resp.Errors[0].Code != "forbidden" {
		t.Fatalf("unexpected error response: %#v", resp.Errors)
	}
}

func TestRequireAuthRequiresAccessTokenForGroupChecks(t *testing.T) {
	db := setupMiddlewareTestDB(t)
	seedMiddlewareUser(t, db, "u1", "ada@example.com")
	svc := newTestAuthServiceWithUsersGroup(t, db)
	rec := runMiddleware(t, RequireAuth(svc), "/api/v1/me", &auth.User{ID: "u1", IsPermitted: true})

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestRequireAuthBlocksDeletedAndMissingUsers(t *testing.T) {
	db := setupMiddlewareTestDB(t)
	seedMiddlewareUser(t, db, "active", "active@example.com")
	seedMiddlewareUser(t, db, "deleted", "deleted@example.com")
	_, err := db.Exec(`UPDATE users SET deleted = 1 WHERE id = 'deleted'`)
	if err != nil {
		t.Fatalf("mark deleted: %v", err)
	}
	svc := newTestAuthServiceWithStore(t, db)

	for _, tc := range []struct {
		name string
		user *auth.User
		want int
	}{
		{
			name: "active local user", want: http.StatusOK,
			user: &auth.User{ID: "active", AuthSource: "internal", IsPermitted: true},
		},
		{
			name: "deleted local user", want: http.StatusUnauthorized,
			user: &auth.User{ID: "deleted", AuthSource: "internal", IsPermitted: true},
		},
		{
			name: "missing hard deleted user", want: http.StatusUnauthorized,
			user: &auth.User{ID: "missing", AuthSource: "internal", IsPermitted: true},
		},
		{
			name: "deleted entra user", want: http.StatusUnauthorized,
			user: &auth.User{ID: "deleted", AuthSource: "entraid", IsPermitted: true},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec := runMiddleware(t, RequireAuth(svc), "/api/v1/me", tc.user)
			if rec.Code != tc.want {
				t.Fatalf("expected %d, got %d", tc.want, rec.Code)
			}
		})
	}
}

func TestRequireAuthExpiresCookieForInactiveUser(t *testing.T) {
	db := setupMiddlewareTestDB(t)
	seedMiddlewareUser(t, db, "deleted", "deleted@example.com")
	_, err := db.Exec(`UPDATE users SET deleted = 1 WHERE id = 'deleted'`)
	if err != nil {
		t.Fatalf("mark deleted: %v", err)
	}
	svc := newTestAuthServiceWithStore(t, db)

	rec := runMiddleware(t, RequireAuth(svc), "/api/v1/me", &auth.User{
		ID: "deleted", AuthSource: "internal", IsPermitted: true,
	})

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
	cookies := rec.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != "sithub_user" || cookies[0].MaxAge != -1 {
		t.Fatalf("expected expired sithub_user cookie, got %#v", cookies)
	}
}

func TestRequireAuthExpiredCookieSecurity(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name        string
		url         string
		forceSecure bool
		wantSecure  bool
	}{
		{name: "local HTTP", url: "http://localhost/api/v1/me"},
		{name: "HTTPS", url: "https://example.com/api/v1/me", wantSecure: true},
		{name: "TLS proxy", url: "http://localhost/api/v1/me", forceSecure: true, wantSecure: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			svc, err := auth.NewService(&config.Config{
				Main: config.MainConfig{ForceSecureCookies: tc.forceSecure},
			}, setupMiddlewareTestDB(t))
			require.NoError(t, err)
			rec := runMiddleware(t, RequireAuth(svc), tc.url, &auth.User{
				ID: "missing", AuthSource: "internal", IsPermitted: true,
			})
			require.Equal(t, http.StatusUnauthorized, rec.Code)
			cookies := rec.Result().Cookies()
			require.Len(t, cookies, 1)
			cookie := cookies[0]
			assert.Equal(t, tc.wantSecure, cookie.Secure)
			assert.True(t, cookie.HttpOnly)
			assert.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
			assert.Equal(t, "sithub_user", cookie.Name)
			assert.Equal(t, "/", cookie.Path)
			assert.Empty(t, cookie.Value)
			assert.Equal(t, -1, cookie.MaxAge)
			assert.True(t, cookie.Expires.Before(time.Now()))
		})
	}
}

func TestRequireAuthRefreshesAdminDemotion(t *testing.T) {
	db := setupMiddlewareTestDB(t)
	seedMiddlewareUser(t, db, "admin", "admin@example.com")
	svc := newTestAuthServiceWithStore(t, db)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &auth.User{ID: "admin", AuthSource: "internal", IsPermitted: true, IsAdmin: true})

	h := RequireAuth(svc)(RequireAdmin()(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}))
	if err := h(c); err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected demoted admin to get 403, got %d", rec.Code)
	}
}

func newTestAuthService(t *testing.T) *auth.Service {
	t.Helper()

	cfg := &config.Config{EntraID: config.EntraIDConfig{
		AuthorizeURL: "https://example.com/auth",
		TokenURL:     "https://example.com/token",
		RedirectURI:  "https://example.com/callback",
		ClientID:     "client",
		ClientSecret: "secret",
	}}

	svc, err := auth.NewService(cfg, nil)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	return svc
}

func newTestAuthServiceWithUsersGroup(t *testing.T, db *sql.DB) *auth.Service {
	t.Helper()

	cfg := &config.Config{EntraID: config.EntraIDConfig{
		AuthorizeURL: "https://example.com/auth",
		TokenURL:     "https://example.com/token",
		RedirectURI:  "https://example.com/callback",
		ClientID:     "client",
		ClientSecret: "secret",
		UsersGroupID: "users",
	}}

	svc, err := auth.NewService(cfg, db)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	return svc
}

func newTestAuthServiceWithStore(t *testing.T, db *sql.DB) *auth.Service {
	t.Helper()

	cfg := &config.Config{EntraID: config.EntraIDConfig{
		AuthorizeURL: "https://example.com/auth",
		TokenURL:     "https://example.com/token",
		RedirectURI:  "https://example.com/callback",
		ClientID:     "client",
		ClientSecret: "secret",
	}}

	svc, err := auth.NewService(cfg, db)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	return svc
}

func setupMiddlewareTestDB(t *testing.T) *sql.DB {
	t.Helper()
	return dbtest.New(t)
}

func seedMiddlewareUser(t *testing.T, db *sql.DB, id, email string) {
	t.Helper()
	_, err := db.Exec(`
		INSERT INTO users (id, email, display_name, password_hash,
			user_source, entra_id, is_admin, last_login, access_token,
			created_at, updated_at)
		VALUES (?, ?, 'Test', '', 'entraid', '', 0, '', '',
			datetime('now'), datetime('now'))`,
		id, email,
	)
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}
}
