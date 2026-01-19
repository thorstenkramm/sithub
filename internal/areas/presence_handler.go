package areas

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/spaces"
)

// PresenceAttributes represents a user present in the area.
type PresenceAttributes struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	DeskID   string `json:"desk_id"`
	DeskName string `json:"desk_name"`
	RoomID   string `json:"room_id"`
	RoomName string `json:"room_name"`
}

// PresenceHandler returns a JSON:API list of users present in an area on a given date.
func PresenceHandler(cfg *spaces.Config, store *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		areaID := c.Param("area_id")
		area, ok := cfg.FindArea(areaID)
		if !ok {
			return api.WriteNotFound(c, "Area not found")
		}

		params, err := api.ParseRoomRequest(areaID, c.QueryParam("date"))
		if err != nil {
			return api.WriteBadRequest(c, "Invalid date. Use YYYY-MM-DD.")
		}

		presence, err := findAreaPresence(c.Request().Context(), store, area, params.BookingDate)
		if err != nil {
			return fmt.Errorf("find area presence: %w", err)
		}

		return api.WriteCollection(c, presence, "write area presence response")
	}
}

// deskDetails holds desk information for presence display.
type deskDetails struct {
	DeskName string
	RoomID   string
	RoomName string
}

// buildDeskIndex creates a map of desk IDs to their details for an area.
func buildDeskIndex(area *spaces.Area) (deskIDs []string, deskInfo map[string]deskDetails) {
	var totalDesks int
	for _, room := range area.Rooms {
		totalDesks += len(room.Desks)
	}

	deskIDs = make([]string, 0, totalDesks)
	deskInfo = make(map[string]deskDetails, totalDesks)

	for _, room := range area.Rooms {
		for _, desk := range room.Desks {
			deskIDs = append(deskIDs, desk.ID)
			deskInfo[desk.ID] = deskDetails{
				DeskName: desk.Name,
				RoomID:   room.ID,
				RoomName: room.Name,
			}
		}
	}

	return deskIDs, deskInfo
}

func findAreaPresence(
	ctx context.Context, store *sql.DB, area *spaces.Area, bookingDate string,
) ([]api.Resource, error) {
	deskIDs, deskInfo := buildDeskIndex(area)
	if len(deskIDs) == 0 {
		return []api.Resource{}, nil
	}

	placeholders, args := api.BuildINClause(deskIDs)
	args = append(args, bookingDate)

	//nolint:gosec // G201: placeholders are "?" literals from BuildINClause, not user input
	query := fmt.Sprintf(
		`SELECT id, desk_id, user_id, user_name 
		 FROM bookings 
		 WHERE desk_id IN (%s) AND booking_date = ?
		 ORDER BY user_name, desk_id`,
		placeholders,
	)

	rows, err := store.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query area presence: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			_ = err
		}
	}()

	return scanPresenceRows(rows, deskInfo)
}

func scanPresenceRows(rows *sql.Rows, deskInfo map[string]deskDetails) ([]api.Resource, error) {
	var resources []api.Resource
	for rows.Next() {
		var bookingID, deskID, userID, userName string
		if err := rows.Scan(&bookingID, &deskID, &userID, &userName); err != nil {
			return nil, fmt.Errorf("scan area presence: %w", err)
		}

		info := deskInfo[deskID]
		resources = append(resources, api.Resource{
			Type: "presence",
			ID:   bookingID,
			Attributes: PresenceAttributes{
				UserID:   userID,
				UserName: userName,
				DeskID:   deskID,
				DeskName: info.DeskName,
				RoomID:   info.RoomID,
				RoomName: info.RoomName,
			},
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate area presence: %w", err)
	}

	return resources, nil
}
