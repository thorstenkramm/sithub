package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
)

// MeHandler returns the authenticated user profile.
func MeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*User)
		if !ok || user == nil {
			return api.WriteUnauthorized(c)
		}

		resp := api.SingleResponse{
			Data: api.Resource{
				Type: "users",
				ID:   user.ID,
				Attributes: map[string]string{
					"display_name": user.Name,
				},
			},
		}

		c.Response().Header().Set(echo.HeaderContentType, api.JSONAPIContentType)
		return c.JSON(http.StatusOK, resp)
	}
}
