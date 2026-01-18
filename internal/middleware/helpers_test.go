package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/auth"
)

func runMiddleware(
	t *testing.T,
	mw echo.MiddlewareFunc,
	path string,
	user *auth.User,
) *httptest.ResponseRecorder {
	t.Helper()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, path, http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if user != nil {
		c.Set("user", user)
	}

	h := mw(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	if err := h(c); err != nil {
		t.Fatalf("handler error: %v", err)
	}

	return rec
}
