// Package middleware provides HTTP middleware for SitHub.
package middleware

import (
	"errors"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/auth"
)

// RequireAuth ensures an authenticated and permitted user is present.
func RequireAuth(svc *auth.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if auth.GetAPIKeyFromContext(c) != nil {
				return denyAPIKey(c)
			}
			user := auth.GetUserFromContext(c)
			if user == nil {
				return api.WriteUnauthorized(c)
			}
			if err := svc.EnsureActive(c.Request().Context(), user); err != nil {
				// Expire the cookie only for a definitively inactive user; a transient
				// database error must not destroy an otherwise valid session.
				if errors.Is(err, auth.ErrUserInactive) {
					// #nosec G124 -- NewCookie sets HttpOnly/SameSite and Secure from TLS or force_secure_cookies.
					cookie := svc.NewCookie(c, "sithub_user", "")
					cookie.Expires = time.Unix(0, 0)
					cookie.MaxAge = -1
					c.SetCookie(cookie)
				}
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
