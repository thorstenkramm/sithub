package system

import (
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
		return api.WriteSingle(c, http.StatusOK, api.Resource{
			Type: resourceTypeVersion,
			ID:   resourceTypeVersion,
			Attributes: VersionAttributes{
				Version: version,
			},
		}, "write version response")
	}
}
