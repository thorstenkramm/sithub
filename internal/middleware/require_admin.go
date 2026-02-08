package middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/auth"
)

// RequireAdmin ensures the authenticated user has admin privileges.
func RequireAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := auth.GetUserFromContext(c)
			if user == nil {
				return api.WriteUnauthorized(c)
			}
			if !user.IsAdmin {
				return api.WriteForbidden(c)
			}
			return next(c)
		}
	}
}
