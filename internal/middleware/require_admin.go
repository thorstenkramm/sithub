package middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/auth"
)

// RequireAdmin ensures the authenticated user has admin privileges.
func RequireAdmin() echo.MiddlewareFunc {
	return requireRole(func(user *auth.User) bool { return user.IsAdmin })
}
