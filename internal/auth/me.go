package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
)

// GetUserFromContext retrieves the authenticated user from the Echo context.
// Returns nil if no user is present or the type assertion fails.
func GetUserFromContext(c echo.Context) *User {
	user, ok := c.Get("user").(*User)
	if !ok {
		return nil
	}
	return user
}

// MeHandler returns the authenticated user profile.
func MeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := GetUserFromContext(c)
		if user == nil {
			return api.WriteUnauthorized(c)
		}

		resp := api.SingleResponse{
			Data: api.Resource{
				Type: "users",
				ID:   user.ID,
				Attributes: map[string]interface{}{
					"display_name": user.Name,
					"is_admin":     user.IsAdmin,
				},
			},
		}

		c.Response().Header().Set(echo.HeaderContentType, api.JSONAPIContentType)
		return c.JSON(http.StatusOK, resp)
	}
}
