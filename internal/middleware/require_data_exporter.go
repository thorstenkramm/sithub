package middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/auth"
)

// RequireDataExporter ensures the authenticated user holds the data exporter role.
// Administrators are NOT implicitly exporters: the flag must be granted explicitly.
// It must run after RequireAuth, which re-reads the flag from the database on every request,
// so a grant or revocation applies without a new login.
func RequireDataExporter() echo.MiddlewareFunc {
	return requireRole(func(user *auth.User) bool { return user.IsDataExporter })
}
