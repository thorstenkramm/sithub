// Package areas provides area handlers.
package areas

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/spaces"
)

// ListHandler returns a JSON:API list of areas.
func ListHandler(cfg *spaces.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		resources := make([]api.Resource, 0, len(cfg.Areas))
		for _, area := range cfg.Areas {
			attrs := map[string]interface{}{
				"name": area.Name,
			}
			if area.Description != "" {
				attrs["description"] = area.Description
			}
			if area.FloorPlan != "" {
				attrs["floor_plan"] = area.FloorPlan
			}
			resources = append(resources, api.Resource{
				Type:       "areas",
				ID:         area.ID,
				Attributes: attrs,
			})
		}

		resp := api.CollectionResponse{Data: resources}
		c.Response().Header().Set(echo.HeaderContentType, api.JSONAPIContentType)
		if err := c.JSON(http.StatusOK, resp); err != nil {
			return fmt.Errorf("write areas response: %w", err)
		}
		return nil
	}
}
