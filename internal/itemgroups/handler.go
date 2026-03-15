// Package itemgroups provides item group handlers.
package itemgroups

import (
	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/areas"
)

// ListHandler returns a JSON:API list of item groups for an area.
func ListHandler(cfg *areas.Config) echo.HandlerFunc {
	return ListHandlerDynamic(func() *areas.Config { return cfg })
}

// ListHandlerDynamic returns a JSON:API list of item groups for an area using dynamic config.
func ListHandlerDynamic(getConfig areas.ConfigGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		cfg := getConfig()
		areaID := c.Param("area_id")
		area, ok := cfg.FindArea(areaID)
		if !ok {
			return api.WriteNotFound(c, "Area not found")
		}

		resources := api.MapResources(area.ItemGroups, func(ig areas.ItemGroup) api.Resource {
			return api.Resource{
				Type:       "item-groups",
				ID:         ig.ID,
				Attributes: areas.BaseAttributes(ig.Name, ig.Description, ig.FloorPlan),
			}
		})

		return api.WriteCollection(c, resources, "write item groups response")
	}
}
