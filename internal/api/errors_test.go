//revive:disable-next-line var-naming
package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestWriteForbidden(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := WriteForbidden(c); err != nil {
		t.Fatalf("write forbidden: %v", err)
	}

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}

	var resp ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Errors) != 1 || resp.Errors[0].Code != "forbidden" {
		t.Fatalf("unexpected error response: %#v", resp.Errors)
	}
}
