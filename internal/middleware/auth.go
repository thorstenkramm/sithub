// Package middleware provides HTTP middleware for SitHub.
package middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
)

// RequireAuth ensures an authenticated user is present.
func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Get("user") == nil {
			return api.WriteUnauthorized(c)
		}
		return next(c)
	}
}
