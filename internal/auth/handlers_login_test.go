package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/config"
)

func TestLoginHandlerRedirects(t *testing.T) {
	cfg := &config.Config{EntraID: testEntraConfig()}
	svc := newTestService(t, cfg)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/oauth/login", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := LoginHandler(svc)
	if err := h(c); err != nil {
		t.Fatalf("handler error: %v", err)
	}

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}
	location := rec.Header().Get("Location")
	if !strings.HasPrefix(location, cfg.EntraID.AuthorizeURL) {
		t.Fatalf("unexpected redirect: %s", location)
	}

	cookies := rec.Result().Cookies()
	if len(cookies) == 0 || cookies[0].Name != stateCookieName {
		t.Fatalf("expected state cookie set")
	}
}

func TestCallbackHandlerBadRequest(t *testing.T) {
	cfg := &config.Config{EntraID: testEntraConfig()}
	svc := newTestService(t, cfg)

	tests := []struct {
		name string
		path string
	}{
		{name: "missing params", path: "/oauth/callback"},
		{name: "missing cookie", path: "/oauth/callback?state=s1&code=c1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, tt.path, http.NoBody)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h := CallbackHandler(svc)
			if err := h(c); err != nil {
				t.Fatalf("handler error: %v", err)
			}

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d", rec.Code)
			}
		})
	}
}

func TestCallbackHandlerInvalidState(t *testing.T) {
	cfg := &config.Config{EntraID: testEntraConfig()}
	svc := newTestService(t, cfg)

	encoded, err := svc.EncodeState("expected")
	if err != nil {
		t.Fatalf("encode state: %v", err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/oauth/callback?state=other&code=c1", http.NoBody)
	req.AddCookie(&http.Cookie{Name: stateCookieName, Value: encoded})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := CallbackHandler(svc)
	if err := h(c); err != nil {
		t.Fatalf("handler error: %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func testEntraConfig() config.EntraIDConfig {
	return config.EntraIDConfig{
		AuthorizeURL: "https://example.com/auth",
		TokenURL:     "https://example.com/token",
		RedirectURI:  "https://example.com/callback",
		ClientID:     "client",
		ClientSecret: "secret",
	}
}

func newTestService(t *testing.T, cfg *config.Config) *Service {
	t.Helper()

	svc, err := NewService(cfg)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	return svc
}
