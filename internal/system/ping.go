// Package system provides system health endpoints.
package system

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
)

// Ping returns a JSON:API health check response.
func Ping(c echo.Context) error {
	resp := api.SingleResponse{
		Data: api.Resource{
			Type: "ping",
			ID:   "ping",
			Attributes: map[string]string{
				"status": "ok",
			},
		},
	}

	c.Response().Header().Set(echo.HeaderContentType, api.JSONAPIContentType)
	if err := c.JSON(http.StatusOK, resp); err != nil {
		return fmt.Errorf("write ping response: %w", err)
	}
	return nil
}
