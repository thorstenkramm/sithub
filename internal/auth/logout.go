package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const schemeHTTPS = "https"

// LogoutHandler clears the authentication cookie and returns 204.
func LogoutHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie := &http.Cookie{
			Name:     userCookieName,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   c.Scheme() == schemeHTTPS,
			SameSite: http.SameSiteLaxMode,
		}
		c.SetCookie(cookie)

		return c.NoContent(http.StatusNoContent)
	}
}
