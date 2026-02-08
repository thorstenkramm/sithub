package auth

import (
	"context"
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/oauth2"

	"github.com/thorstenkramm/sithub/internal/config"
)

func TestCallbackHandlerSuccess(t *testing.T) {
	db := setupCallbackTestDB(t)
	cfg := &config.Config{EntraID: entraConfig()}
	cfg.EntraID.AdminsGroupID = "admins"
	svc := newAuthService(t, cfg, db)
	httpClient := newAuthTestClient(cfg.EntraID.TokenURL)

	state := "state-123"
	encoded, err := svc.EncodeState(state)
	if err != nil {
		t.Fatalf("encode state: %v", err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/oauth/callback?state="+state+"&code=abc", http.NoBody)
	req.AddCookie(&http.Cookie{Name: stateCookieName, Value: encoded})
	req = req.WithContext(context.WithValue(req.Context(), oauth2.HTTPClient, httpClient))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := CallbackHandler(svc)
	if err := h(c); err != nil {
		t.Fatalf("handler error: %v", err)
	}

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}

	userCookies := rec.Result().Cookies()
	if len(userCookies) == 0 || userCookies[0].Name != userCookieName {
		t.Fatalf("expected user cookie set")
	}

	decoded, err := svc.DecodeUser(userCookies[0].Value)
	if err != nil {
		t.Fatalf("decode user: %v", err)
	}
	if !decoded.IsAdmin {
		t.Fatalf("expected admin user, got %#v", decoded)
	}
}

func setupCallbackTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close() //nolint:errcheck // Cleanup function, error not critical
	})

	_, err = db.Exec(`
		CREATE TABLE users (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL,
			display_name TEXT NOT NULL,
			password_hash TEXT NOT NULL DEFAULT '',
			user_source TEXT NOT NULL CHECK (user_source IN ('internal', 'entraid')),
			entra_id TEXT NOT NULL DEFAULT '',
			is_admin INTEGER NOT NULL DEFAULT 0,
			last_login TEXT NOT NULL DEFAULT '',
			access_token TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE UNIQUE INDEX idx_users_email ON users(email);
		CREATE INDEX idx_users_entra_id ON users(entra_id);
	`)
	if err != nil {
		t.Fatalf("create users table: %v", err)
	}
	return db
}

func entraConfig() config.EntraIDConfig {
	return config.EntraIDConfig{
		AuthorizeURL: "https://example.com/auth",
		TokenURL:     "https://example.com/token",
		RedirectURI:  "https://example.com/callback",
		ClientID:     "client",
		ClientSecret: "secret",
	}
}

func newAuthService(t *testing.T, cfg *config.Config, db *sql.DB) *Service {
	t.Helper()

	svc, err := NewService(cfg, db)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	return svc
}

func newAuthTestClient(tokenURL string) *http.Client {
	return &http.Client{Transport: roundTripper(func(req *http.Request) (*http.Response, error) {
		switch req.URL.String() {
		case tokenURL:
			body := `{"access_token":"token","token_type":"Bearer","expires_in":3600}`
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		case graphMeURLWithSelect:
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(graphMeBody)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		case graphMemberOfURL:
			body := `{"value":[{"@odata.type":"#microsoft.graph.group","id":"admins"}]}`
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		default:
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     http.Header{},
			}, nil
		}
	})}
}
