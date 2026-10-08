// Package system provides system health endpoints.
package system

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
)

const resourceTypePing = "ping"

// Ping returns a JSON:API health check response.
func Ping(c echo.Context) error {
	//nolint:wrapcheck // Terminal response
	return api.WriteSingle(c, http.StatusOK, api.Resource{
		Type: resourceTypePing,
		ID:   resourceTypePing,
		Attributes: map[string]string{
			"status": "ok",
		},
	}, "write ping response")
}
