package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
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
