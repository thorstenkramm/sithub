package middleware

import (
	"encoding/json"
	"net/http"
	"testing"

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
	svc := newTestAuthServiceWithUsersGroup(t)
	rec := runMiddleware(t, RequireAuth(svc), "/api/v1/me", &auth.User{IsPermitted: true})

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

	svc, err := auth.NewService(cfg)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	return svc
}

func newTestAuthServiceWithUsersGroup(t *testing.T) *auth.Service {
	t.Helper()

	cfg := &config.Config{EntraID: config.EntraIDConfig{
		AuthorizeURL: "https://example.com/auth",
		TokenURL:     "https://example.com/token",
		RedirectURI:  "https://example.com/callback",
		ClientID:     "client",
		ClientSecret: "secret",
		UsersGroupID: "users",
	}}

	svc, err := auth.NewService(cfg)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	return svc
}
