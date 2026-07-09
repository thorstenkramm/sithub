package system

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
)

const resourceTypeVersion = "version"

// VersionAttributes contains the running application version.
type VersionAttributes struct {
	Version string `json:"version"`
}

// Version returns a handler that reports the running application version.
// The version value is captured at wiring time (injected via build ldflags).
func Version(version string) echo.HandlerFunc {
	return func(c echo.Context) error {
		resp := api.SingleResponse{
			Data: api.Resource{
				Type: resourceTypeVersion,
				ID:   resourceTypeVersion,
				Attributes: VersionAttributes{
					Version: version,
				},
			},
		}

		c.Response().Header().Set(echo.HeaderContentType, api.JSONAPIContentType)
		if err := c.JSON(http.StatusOK, resp); err != nil {
			return fmt.Errorf("write version response: %w", err)
		}
		return nil
	}
}
