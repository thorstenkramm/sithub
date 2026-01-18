package middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/auth"
)

// LoadUser loads the authenticated user from cookies or test auth.
func LoadUser(svc *auth.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if user := svc.TestUser(); user != nil {
				c.Set("user", user)
				return next(c)
			}

			cookie, err := c.Cookie("sithub_user")
			if err == nil {
				user, err := svc.DecodeUser(cookie.Value)
				if err == nil && user != nil {
					c.Set("user", user)
				}
			}
			return next(c)
		}
	}
}
