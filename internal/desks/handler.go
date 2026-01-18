// Package desks provides desk handlers.
package desks

import (
	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/spaces"
)

// ListHandler returns a JSON:API list of desks for a room.
func ListHandler(cfg *spaces.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		roomID := c.Param("room_id")
		room, ok := cfg.FindRoom(roomID)
		if !ok {
			return api.WriteNotFound(c, "Room not found")
		}

		resources := api.MapResources(room.Desks, func(desk spaces.Desk) api.Resource {
			return api.Resource{
				Type:       "desks",
				ID:         desk.ID,
				Attributes: spaces.DeskAttributes(desk.Name, desk.Equipment, desk.Warning),
			}
		})

		return api.WriteCollection(c, resources, "write desks response")
	}
}
