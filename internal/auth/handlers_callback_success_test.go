package auth

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"

	"github.com/thorstenkramm/sithub/internal/config"
)

func TestCallbackHandlerSuccess(t *testing.T) {
	cfg := &config.Config{EntraID: entraConfig()}
	cfg.EntraID.AdminsGroupID = "admins"
	svc := newAuthService(t, cfg)
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

func TestCallbackHandlerTestAuth(t *testing.T) {
	cfg := &config.Config{EntraID: config.EntraIDConfig{
		AuthorizeURL: "https://example.com/auth",
		TokenURL:     "https://example.com/token",
		RedirectURI:  "https://example.com/callback",
		ClientID:     "client",
		ClientSecret: "secret",
	}, TestAuth: config.TestAuthConfig{
		Enabled:  true,
		UserID:   "u-123",
		UserName: "Ada Lovelace",
	}}
	svc := newAuthService(t, cfg)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/oauth/callback", http.NoBody)
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

func newAuthService(t *testing.T, cfg *config.Config) *Service {
	t.Helper()

	svc, err := NewService(cfg)
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
		case "https://graph.microsoft.com/v1.0/me":
			body := `{"id":"u1","displayName":"Ada"}`
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		case "https://graph.microsoft.com/v1.0/me/memberOf?$select=id":
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
