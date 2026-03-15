// Package areas provides area configuration, handlers, and domain types.
package areas

import (
	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
)

// ListHandler returns a JSON:API list of areas.
func ListHandler(cfg *Config) echo.HandlerFunc {
	return ListHandlerDynamic(func() *Config { return cfg })
}

// ListHandlerDynamic returns a JSON:API list of areas using a dynamic config getter.
func ListHandlerDynamic(getConfig ConfigGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		cfg := getConfig()
		resources := api.MapResources(cfg.Areas, func(area Area) api.Resource {
			return api.Resource{
				Type:       "areas",
				ID:         area.ID,
				Attributes: BaseAttributes(area.Name, area.Description, area.FloorPlan),
			}
		})

		return api.WriteCollection(c, resources, "write areas response")
	}
}
