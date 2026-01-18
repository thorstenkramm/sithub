package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
)

func TestMeHandlerUnauthorized(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := MeHandler()
	if err := h(c); err != nil {
		t.Fatalf("handler error: %v", err)
	}

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

func TestMeHandlerAuthorized(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &User{ID: "u1", Name: "Ada", IsAdmin: true})

	h := MeHandler()
	if err := h(c); err != nil {
		t.Fatalf("handler error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp api.SingleResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	attrs, ok := resp.Data.Attributes.(map[string]interface{})
	if !ok {
		t.Fatalf("unexpected attributes type: %T", resp.Data.Attributes)
	}
	if attrs["display_name"] != "Ada" {
		t.Fatalf("unexpected display_name: %v", attrs["display_name"])
	}
	if attrs["is_admin"] != true {
		t.Fatalf("unexpected is_admin: %v", attrs["is_admin"])
	}
}
