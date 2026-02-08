package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/auth"
)

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

			if err := svc.RefreshPermissions(req.Context(), user); err != nil {
				return c.Redirect(http.StatusFound, "/login")
			}
			if !user.IsPermitted {
				return c.Redirect(http.StatusFound, "/access-denied")
			}

			return next(c)
		}
	}
}

func isBypassPath(path string) bool {
	return strings.HasPrefix(path, "/api/") ||
		strings.HasPrefix(path, "/oauth/") ||
		strings.HasPrefix(path, "/auth/") ||
		strings.HasPrefix(path, "/assets/") ||
		path == "/login"
}
