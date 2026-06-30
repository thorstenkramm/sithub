package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/auth"
	"github.com/thorstenkramm/sithub/internal/config"
)

// TestInvalidSessionCookieRejected verifies that a tampered or expired session
// cookie fails to decode in LoadUser, leaving no user in context, so RequireAuth
// rejects the request with 401. A genuinely expired securecookie value fails the
// same Decode path, producing identical behavior.
func TestInvalidSessionCookieRejected(t *testing.T) {
	svc, err := auth.NewService(&config.Config{}, nil)
	require.NoError(t, err)

	e := echo.New()
	protected := LoadUser(svc)(RequireAuth(svc)(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", http.NoBody)
	req.AddCookie(&http.Cookie{Name: "sithub_user", Value: "tampered-or-expired-cookie-value"})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, protected(c))
	assert.Equal(t, http.StatusUnauthorized, rec.Code,
		"a tampered/expired session cookie must be rejected as unauthenticated")
}
