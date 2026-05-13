package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
)

// EntraIDConfigured reports whether the service was constructed with a valid
// Entra ID OAuth configuration. Used by ProvidersHandler to tell the login
// page which authentication options to surface.
func (s *Service) EntraIDConfigured() bool {
	return s.oauthConfig != nil
}

const (
	providerEntraID = "entraid"
	providerLocal   = "local"
)

// ProvidersHandler returns GET /api/v1/auth/providers exposing which
// authentication providers are available on this server. The endpoint is
// public (unauthenticated) so the login page can render the correct affordances
// before the user signs in.
func ProvidersHandler(svc *Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		resp := api.SingleResponse{
			Data: api.Resource{
				Type: "auth-providers",
				ID:   "current",
				Attributes: map[string]interface{}{
					providerEntraID: svc.EntraIDConfigured(),
					providerLocal:   true,
				},
			},
		}

		c.Response().Header().Set(echo.HeaderContentType, api.JSONAPIContentType)
		return c.JSON(http.StatusOK, resp)
	}
}
