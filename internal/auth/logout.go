package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const schemeHTTPS = "https"

// LogoutHandler clears the authentication cookie and returns 204.
func LogoutHandler(svc *Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie := svc.NewCookie(c, userCookieName, "")
		cookie.MaxAge = -1
		c.SetCookie(cookie)

		return c.NoContent(http.StatusNoContent)
	}
}
