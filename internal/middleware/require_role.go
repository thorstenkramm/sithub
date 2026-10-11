package middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/auth"
)

// requireRole allows the request only if the authenticated user passes has.
// A missing user yields 401, a user without the role 403. It must run after RequireAuth.
func requireRole(has func(*auth.User) bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if auth.GetAPIKeyFromContext(c) != nil {
				return denyAPIKey(c)
			}
			user := auth.GetUserFromContext(c)
			if user == nil {
				return api.WriteUnauthorized(c)
			}
			if !has(user) {
				return api.WriteForbidden(c)
			}
			return next(c)
		}
	}
}
