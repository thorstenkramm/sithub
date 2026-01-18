package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/auth"
	"github.com/thorstenkramm/sithub/internal/config"
)

func TestLoadUserFromTestAuth(t *testing.T) {
	cfg := &config.Config{EntraID: config.EntraIDConfig{
		AuthorizeURL: "https://example.com/auth",
		TokenURL:     "https://example.com/token",
		RedirectURI:  "https://example.com/callback",
		ClientID:     "client",
		ClientSecret: "secret",
	}, TestAuth: config.TestAuthConfig{
		Enabled:  true,
		UserID:   "u-test",
		UserName: "Test User",
	}}

	svc, err := auth.NewService(cfg)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	h := LoadUser(svc)(func(c echo.Context) error {
		if c.Get("user") == nil {
			t.Fatal("expected user in context")
		}
		return c.NoContent(http.StatusOK)
	})

	if err := h(c); err != nil {
		t.Fatalf("handler error: %v", err)
	}
}

func TestLoadUserFromCookie(t *testing.T) {
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

	userCookie, err := svc.EncodeUser(auth.User{ID: "u1", Name: "Ada"})
	if err != nil {
		t.Fatalf("encode user: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	req.AddCookie(&http.Cookie{Name: "sithub_user", Value: userCookie})
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	h := LoadUser(svc)(func(c echo.Context) error {
		user, ok := c.Get("user").(*auth.User)
		if !ok || user == nil {
			t.Fatal("expected user in context")
		}
		if user.ID != "u1" {
			t.Fatalf("unexpected user id: %s", user.ID)
		}
		return c.NoContent(http.StatusOK)
	})

	if err := h(c); err != nil {
		t.Fatalf("handler error: %v", err)
	}
}
