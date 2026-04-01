package floorplanpos

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
)

func toResource(p Position) api.Resource { //nolint:gocritic // value needed for MapResources
	attrs := map[string]interface{}{
		"floor_plan": p.FloorPlan,
		"item_id":    p.ItemID,
		"x":          p.X,
		"y":          p.Y,
		"width":      p.Width,
		"height":       p.Height,
		"border_width": p.BorderWidth,
		"created_at":   p.CreatedAt,
		"updated_at":   p.UpdatedAt,
	}
	if p.Label != "" {
		attrs["label"] = p.Label
	}
	return api.Resource{
		Type:       "floor-plan-positions",
		ID:         p.ID,
		Attributes: attrs,
	}
}

// ListHandler returns positions for a floor plan.
// GET /api/v1/floor-plan-positions?floor_plan=<filename>
func ListHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		floorPlan := c.QueryParam("floor_plan")
		if floorPlan == "" {
			return api.WriteBadRequest(c, "Missing floor_plan query parameter")
		}

		positions, err := FindByFloorPlan(c.Request().Context(), db, floorPlan)
		if err != nil {
			return api.WriteInternalError(c, "list positions", err)
		}

		resources := api.MapResources(positions, toResource)
		return api.WriteCollection(c, resources, "write positions response")
	}
}

type createRequest struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			FloorPlan string  `json:"floor_plan"`
			ItemID    string  `json:"item_id"`
			Label     string  `json:"label"`
			X         float64 `json:"x"`
			Y         float64 `json:"y"`
			Width       float64 `json:"width"`
			Height      float64 `json:"height"`
			BorderWidth int     `json:"border_width"`
		} `json:"attributes"`
	} `json:"data"`
}

// CreateHandler creates a new position.
// POST /api/v1/floor-plan-positions
func CreateHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req createRequest
		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			return api.WriteBadRequest(c, "Invalid request body")
		}

		a := req.Data.Attributes
		if a.FloorPlan == "" || a.ItemID == "" {
			return api.WriteBadRequest(c, "floor_plan and item_id are required")
		}

		pos, err := Create(c.Request().Context(), db, &CreateInput{
			FloorPlan:   a.FloorPlan,
			ItemID:      a.ItemID,
			Label:       a.Label,
			X:           a.X,
			Y:           a.Y,
			Width:       a.Width,
			Height:      a.Height,
			BorderWidth: a.BorderWidth,
		})
		if err != nil {
			return api.WriteInternalError(c, "create position", err)
		}

		resource := toResource(*pos)
		return api.WriteSingle(c, http.StatusCreated, resource, "write position response")
	}
}

type updateRequest struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			Label       *string  `json:"label"`
			X           *float64 `json:"x"`
			Y           *float64 `json:"y"`
			Width       *float64 `json:"width"`
			Height      *float64 `json:"height"`
			BorderWidth *int     `json:"border_width"`
		} `json:"attributes"`
	} `json:"data"`
}

// UpdateHandler updates a position.
// PUT /api/v1/floor-plan-positions/:id
func UpdateHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		var req updateRequest
		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			return api.WriteBadRequest(c, "Invalid request body")
		}

		a := req.Data.Attributes
		pos, err := Update(c.Request().Context(), db, id, UpdateInput{
			Label:       a.Label,
			X:           a.X,
			Y:           a.Y,
			Width:       a.Width,
			Height:      a.Height,
			BorderWidth: a.BorderWidth,
		})
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				return api.WriteNotFound(c, "Position not found")
			}
			return api.WriteInternalError(c, "update position", err)
		}

		resource := toResource(*pos)
		return api.WriteSingle(c, http.StatusOK, resource, "write position response")
	}
}

// DeleteHandler removes a position.
// DELETE /api/v1/floor-plan-positions/:id
func DeleteHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		if err := Delete(c.Request().Context(), db, id); err != nil {
			if errors.Is(err, ErrNotFound) {
				return api.WriteNotFound(c, "Position not found")
			}
			return api.WriteInternalError(c, "delete position", err)
		}

		return c.NoContent(http.StatusNoContent)
	}
}
