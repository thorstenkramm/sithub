package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/config"
)

func TestJSONAPIError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := jsonAPIError(c, http.StatusBadRequest, "Bad", "Details", "code"); err != nil {
		t.Fatalf("jsonAPIError: %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestNewCookie(t *testing.T) {
	cookie := newCookie("name", "value", true)
	if cookie.Name != "name" || cookie.Value != "value" {
		t.Fatalf("unexpected cookie: %#v", cookie)
	}
	if !cookie.HttpOnly || !cookie.Secure {
		t.Fatalf("expected secure http-only cookie")
	}
	if cookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("unexpected samesite: %v", cookie.SameSite)
	}
}

func TestServiceNewCookieForceSecureOverHTTP(t *testing.T) {
	svc := newTestService(t, &config.Config{Main: config.MainConfig{ForceSecureCookies: true}})
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "http://example.com/oauth/login", http.NoBody)
	c := e.NewContext(req, httptest.NewRecorder())

	// State cookie (set during the OAuth login flow) must be Secure when forced.
	state := svc.NewCookie(c, stateCookieName, "value")
	if !state.Secure {
		t.Fatal("expected state cookie to be Secure when force_secure_cookies is enabled")
	}
	// Session cookie likewise.
	session := svc.NewCookie(c, userCookieName, "value")
	if !session.Secure {
		t.Fatal("expected session cookie to be Secure when force_secure_cookies is enabled")
	}
}

func TestServiceNewCookieDefaultStateCookieOverHTTP(t *testing.T) {
	svc := newTestService(t, &config.Config{})
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "http://example.com/oauth/login", http.NoBody)
	c := e.NewContext(req, httptest.NewRecorder())

	state := svc.NewCookie(c, stateCookieName, "value")
	if state.Secure {
		t.Fatal("expected state cookie to follow default HTTP behavior without Secure flag")
	}
}
