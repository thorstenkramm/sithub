package system

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
)

func TestPing(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/ping", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := Ping(c); err != nil {
		t.Fatalf("ping: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	if ct := rec.Header().Get(echo.HeaderContentType); ct != api.JSONAPIContentType {
		t.Fatalf("unexpected content type: %s", ct)
	}

	var resp api.SingleResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Data.Type != "ping" {
		t.Fatalf("unexpected type: %s", resp.Data.Type)
	}
}
