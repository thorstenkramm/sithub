// Package areas provides area handlers.
package areas

import (
	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/spaces"
)

// ListHandler returns a JSON:API list of areas.
func ListHandler(cfg *spaces.Config) echo.HandlerFunc {
	return ListHandlerDynamic(func() *spaces.Config { return cfg })
}

// ListHandlerDynamic returns a JSON:API list of areas using a dynamic config getter.
func ListHandlerDynamic(getConfig spaces.ConfigGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		cfg := getConfig()
		resources := api.MapResources(cfg.Areas, func(area spaces.Area) api.Resource {
			return api.Resource{
				Type:       "areas",
				ID:         area.ID,
				Attributes: spaces.BaseAttributes(area.Name, area.Description, area.FloorPlan),
			}
		})

		return api.WriteCollection(c, resources, "write areas response")
	}
}
