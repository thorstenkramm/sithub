package rooms

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/spaces"
)

// RoomBookingAttributes represents a booking in the room overview.
type RoomBookingAttributes struct {
	DeskID      string `json:"desk_id"`
	DeskName    string `json:"desk_name"`
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	BookingDate string `json:"booking_date"`
	IsGuest     bool   `json:"is_guest,omitempty"`
}

// BookingsHandler returns a JSON:API list of bookings for a room on a given date.
func BookingsHandler(cfg *spaces.Config, store *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		room, err := validateRoomParam(cfg, c.Param("room_id"))
		if err != nil {
			return api.WriteNotFound(c, "Room not found")
		}

		params, err := api.ParseRoomRequest(room.ID, c.QueryParam("date"))
		if err != nil {
			return api.WriteBadRequest(c, "Invalid booking date. Use YYYY-MM-DD.")
		}

		roomBookings, err := findRoomBookings(c.Request().Context(), store, room, params.BookingDate)
		if err != nil {
			return fmt.Errorf("find room bookings: %w", err)
		}

		return api.WriteCollection(c, roomBookings, "write room bookings response")
	}
}

// validateRoomParam checks if a room exists in the config.
func validateRoomParam(cfg *spaces.Config, roomID string) (*spaces.Room, error) {
	room, ok := cfg.FindRoom(roomID)
	if !ok {
		return nil, fmt.Errorf("room %s not found", roomID)
	}
	return room, nil
}

// roomBookingRecord represents a booking row joined with desk info.
type roomBookingRecord struct {
	ID          string
	DeskID      string
	UserID      string
	UserName    string
	BookingDate string
	IsGuest     bool
}

func findRoomBookings(
	ctx context.Context, store *sql.DB, room *spaces.Room, bookingDate string,
) ([]api.Resource, error) {
	// Build list of desk IDs in this room
	deskIDs := make([]string, 0, len(room.Desks))
	deskNames := make(map[string]string, len(room.Desks))
	for _, desk := range room.Desks {
		deskIDs = append(deskIDs, desk.ID)
		deskNames[desk.ID] = desk.Name
	}

	if len(deskIDs) == 0 {
		return []api.Resource{}, nil
	}

	placeholders, args := api.BuildINClause(deskIDs)
	args = append(args, bookingDate)

	//nolint:gosec // G201: placeholders are "?" literals from BuildINClause, not user input
	query := fmt.Sprintf(
		`SELECT id, desk_id, user_id, user_name, booking_date, is_guest 
		 FROM bookings 
		 WHERE desk_id IN (%s) AND booking_date = ?
		 ORDER BY desk_id`,
		placeholders,
	)

	rows, err := store.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query room bookings: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			// Log error but don't override main error
			_ = closeErr
		}
	}()

	var resources []api.Resource
	for rows.Next() {
		var rec roomBookingRecord
		var isGuestInt int
		err := rows.Scan(
			&rec.ID, &rec.DeskID, &rec.UserID, &rec.UserName, &rec.BookingDate, &isGuestInt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan room booking: %w", err)
		}
		rec.IsGuest = isGuestInt == 1

		resources = append(resources, api.Resource{
			Type: "bookings",
			ID:   rec.ID,
			Attributes: RoomBookingAttributes{
				DeskID:      rec.DeskID,
				DeskName:    deskNames[rec.DeskID],
				UserID:      rec.UserID,
				UserName:    rec.UserName,
				BookingDate: rec.BookingDate,
				IsGuest:     rec.IsGuest,
			},
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate room bookings: %w", err)
	}

	return resources, nil
}
