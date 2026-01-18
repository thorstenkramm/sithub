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

	assertErrorResponse(t, rec, http.StatusForbidden, "forbidden")
}

func TestWriteBadRequest(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/rooms/room-1/desks", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := WriteBadRequest(c, "Invalid booking date"); err != nil {
		t.Fatalf("write bad request: %v", err)
	}

	assertErrorResponse(t, rec, http.StatusBadRequest, "bad_request")
}

func TestWriteNotFound(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/areas/missing/rooms", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := WriteNotFound(c, "Area not found"); err != nil {
		t.Fatalf("write not found: %v", err)
	}

	assertErrorResponse(t, rec, http.StatusNotFound, "not_found")
}

func assertErrorResponse(t *testing.T, rec *httptest.ResponseRecorder, status int, code string) {
	t.Helper()

	if rec.Code != status {
		t.Fatalf("expected %d, got %d", status, rec.Code)
	}

	var resp ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Errors) != 1 || resp.Errors[0].Code != code {
		t.Fatalf("unexpected error response: %#v", resp.Errors)
	}
}
