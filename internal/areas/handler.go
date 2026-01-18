// Package areas provides area handlers.
package areas

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
)

// ListHandler returns a JSON:API list of areas.
func ListHandler(repo *Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		areas, err := repo.List(c.Request().Context())
		if err != nil {
			return fmt.Errorf("list areas: %w", err)
		}

		resources := make([]api.Resource, 0, len(areas))
		for _, area := range areas {
			resources = append(resources, api.Resource{
				Type: "areas",
				ID:   area.ID,
				Attributes: map[string]interface{}{
					"name":       area.Name,
					"sort_order": area.SortOrder,
					"created_at": area.CreatedAt,
					"updated_at": area.UpdatedAt,
				},
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
