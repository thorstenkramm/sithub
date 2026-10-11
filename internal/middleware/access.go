package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/auth"
)

const loginPath = "/login"

// RedirectForbidden ensures forbidden users are redirected away from SPA routes.
func RedirectForbidden(svc *auth.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			if req == nil || req.Method != http.MethodGet {
				return next(c)
			}

			path := req.URL.Path
			if path == "/access-denied" || isBypassPath(path) {
				return next(c)
			}

			user, ok := c.Get("user").(*auth.User)
			if !ok || user == nil {
				return next(c)
			}

			if target := redirectTarget(svc, c, user); target != "" {
				return c.Redirect(http.StatusFound, target)
			}

			return next(c)
		}
	}
}

// redirectTarget returns where a user must be sent instead of the SPA page, or "" if access is allowed.
// A deactivated or deleted user goes to the login page, as does a failed permission refresh.
func redirectTarget(svc *auth.Service, c echo.Context, user *auth.User) string {
	ctx := c.Request().Context()
	if err := svc.EnsureActive(ctx, user); err != nil {
		return loginPath
	}
	if err := svc.RefreshPermissions(ctx, user); err != nil {
		return loginPath
	}
	if !user.IsPermitted {
		return "/access-denied"
	}
	return ""
}

func isBypassPath(path string) bool {
	return strings.HasPrefix(path, "/api/") ||
		strings.HasPrefix(path, "/oauth/") ||
		strings.HasPrefix(path, "/auth/") ||
		strings.HasPrefix(path, "/assets/") ||
		path == loginPath
}
