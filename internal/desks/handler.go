// Package desks provides desk handlers.
package desks

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/auth"
	"github.com/thorstenkramm/sithub/internal/bookings"
	"github.com/thorstenkramm/sithub/internal/spaces"
)

// ListHandler returns a JSON:API list of desks for a room.
// For admin users, includes booking details (booking_id, booker_name) for occupied desks.
func ListHandler(cfg *spaces.Config, store *sql.DB) echo.HandlerFunc {
	return ListHandlerDynamic(func() *spaces.Config { return cfg }, store)
}

// ListHandlerDynamic returns a JSON:API list of desks for a room using dynamic config.
func ListHandlerDynamic(getConfig spaces.ConfigGetter, store *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		cfg := getConfig()
		roomID := c.Param("room_id")
		room, ok := cfg.FindRoom(roomID)
		if !ok {
			return api.WriteNotFound(c, "Room not found")
		}

		bookingDate, err := api.ParseBookingDate(c.QueryParam("date"))
		if err != nil {
			return api.WriteBadRequest(c, "Invalid booking date. Use YYYY-MM-DD.")
		}

		ctx := c.Request().Context()
		user := auth.GetUserFromContext(c)
		isAdmin := user != nil && user.IsAdmin

		deskBookings, err := loadDeskBookings(ctx, store, bookingDate, isAdmin)
		if err != nil {
			return err
		}

		resources := buildDeskResources(room.Desks, deskBookings, isAdmin)
		return api.WriteCollection(c, resources, "write desks response")
	}
}

func loadDeskBookings(
	ctx context.Context, store *sql.DB, bookingDate string, isAdmin bool,
) (map[string]bookings.DeskBookingInfo, error) {
	if isAdmin {
		info, err := bookings.FindDeskBookings(ctx, store, bookingDate)
		if err != nil {
			return nil, fmt.Errorf("list desk bookings: %w", err)
		}
		return info, nil
	}

	bookedDesks, err := bookings.FindBookedDeskIDs(ctx, store, bookingDate)
	if err != nil {
		return nil, fmt.Errorf("list booked desks: %w", err)
	}
	// Convert to DeskBookingInfo map for uniform handling
	result := make(map[string]bookings.DeskBookingInfo, len(bookedDesks))
	for deskID := range bookedDesks {
		result[deskID] = bookings.DeskBookingInfo{}
	}
	return result, nil
}

func buildDeskResources(
	desks []spaces.Desk, deskBookings map[string]bookings.DeskBookingInfo, isAdmin bool,
) []api.Resource {
	return api.MapResources(desks, func(desk spaces.Desk) api.Resource {
		attrs := spaces.DeskAttributes(desk.Name, desk.Equipment, desk.Warning, "")
		if info, booked := deskBookings[desk.ID]; booked {
			attrs["availability"] = "occupied"
			if isAdmin {
				attrs["booking_id"] = info.BookingID
				attrs["booker_name"] = info.BookerName
			}
		} else {
			attrs["availability"] = "available"
		}
		return api.Resource{
			Type:       "desks",
			ID:         desk.ID,
			Attributes: attrs,
		}
	})
}
