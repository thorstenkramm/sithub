package middleware

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/auth"
	"github.com/thorstenkramm/sithub/internal/config"
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

func setupMiddlewareTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close() //nolint:errcheck // Cleanup
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
	`)
	if err != nil {
		t.Fatalf("create users table: %v", err)
	}
	return db
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
