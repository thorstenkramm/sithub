package middleware

import (
	"net/http"
	"testing"

	"github.com/thorstenkramm/sithub/internal/auth"
)

func TestRedirectForbiddenSkipsBypassPaths(t *testing.T) {
	svc := newTestAuthService(t)
	rec := runMiddleware(t, RedirectForbidden(svc), "/api/v1/me", &auth.User{IsPermitted: false})

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestRedirectForbiddenSendsToAccessDenied(t *testing.T) {
	svc := newTestAuthService(t)
	rec := runMiddleware(t, RedirectForbidden(svc), "/", &auth.User{IsPermitted: false})

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/access-denied" {
		t.Fatalf("expected /access-denied, got %s", loc)
	}
}
