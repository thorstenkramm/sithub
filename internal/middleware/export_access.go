package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/auth"
)

// RequireExportAccess explicitly opts an export endpoint into machine access.
// Keys need the exact scope; sessions must still pass active-user and exporter checks.
func RequireExportAccess(svc *auth.Service, scope string) echo.MiddlewareFunc {
	return keyOrSession(scope, func(next echo.HandlerFunc) echo.HandlerFunc {
		return RequireAuth(svc)(RequireDataExporter()(next))
	})
}

// RequireAuthOrScope explicitly opts a read-only endpoint into machine access: any authenticated and
// permitted user may call it, and so may an API key holding exactly the given scope (any other key
// gets 403). It is used for the avatar endpoint so integrations can show the people listed by the
// booking export.
func RequireAuthOrScope(svc *auth.Service, scope string) echo.MiddlewareFunc {
	return keyOrSession(scope, RequireAuth(svc))
}

// keyOrSession lets an API key with the scope through and otherwise applies the session gate.
// A key never falls back to the session gate: a key without the scope is denied.
func keyOrSession(scope string, session echo.MiddlewareFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		sessionAccess := session(next)
		return func(c echo.Context) error {
			if key := auth.GetAPIKeyFromContext(c); key != nil {
				if !key.HasScope(scope) {
					challenge := `Bearer error="insufficient_scope", scope="` + scope + `"`
					c.Response().Header().Set("WWW-Authenticate", challenge)
					return denyAPIKey(c)
				}
				return next(c)
			}
			return sessionAccess(c)
		}
	}
}

// denyAPIKey is the default-deny response shared by user-only authorization gates.
func denyAPIKey(c echo.Context) error {
	//nolint:wrapcheck // Terminal response
	return api.WriteForbiddenDetail(c, "API key not permitted for this endpoint")
}
