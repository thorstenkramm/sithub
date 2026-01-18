// Package middleware provides HTTP middleware for SitHub.
package middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/auth"
)

// RequireAuth ensures an authenticated and permitted user is present.
func RequireAuth(svc *auth.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get("user").(*auth.User)
			if !ok || user == nil {
				return api.WriteUnauthorized(c)
			}
			if err := svc.RefreshPermissions(c.Request().Context(), user); err != nil {
				return api.WriteUnauthorized(c)
			}
			if !user.IsPermitted {
				return api.WriteForbidden(c)
			}
			return next(c)
		}
	}
}
