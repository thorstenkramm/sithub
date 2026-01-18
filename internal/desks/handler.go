// Package desks provides desk handlers.
package desks

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/bookings"
	"github.com/thorstenkramm/sithub/internal/spaces"
)

// ListHandler returns a JSON:API list of desks for a room.
func ListHandler(cfg *spaces.Config, store *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		roomID := c.Param("room_id")
		room, ok := cfg.FindRoom(roomID)
		if !ok {
			return api.WriteNotFound(c, "Room not found")
		}

		bookingDate, err := parseBookingDate(c.QueryParam("date"))
		if err != nil {
			return api.WriteBadRequest(c, "Invalid booking date. Use YYYY-MM-DD.")
		}

		bookedDesks, err := bookings.FindBookedDeskIDs(c.Request().Context(), store, bookingDate)
		if err != nil {
			return fmt.Errorf("list booked desks: %w", err)
		}

		resources := api.MapResources(room.Desks, func(desk spaces.Desk) api.Resource {
			availability := "available"
			if _, booked := bookedDesks[desk.ID]; booked {
				availability = "occupied"
			}
			return api.Resource{
				Type:       "desks",
				ID:         desk.ID,
				Attributes: spaces.DeskAttributes(desk.Name, desk.Equipment, desk.Warning, availability),
			}
		})

		return api.WriteCollection(c, resources, "write desks response")
	}
}

func parseBookingDate(value string) (string, error) {
	if strings.TrimSpace(value) == "" {
		return time.Now().Format(time.DateOnly), nil
	}
	parsed, err := time.Parse(time.DateOnly, value)
	if err != nil {
		return "", fmt.Errorf("parse booking date: %w", err)
	}
	return parsed.Format(time.DateOnly), nil
}
