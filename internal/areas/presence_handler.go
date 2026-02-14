package areas

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/spaces"
	"github.com/thorstenkramm/sithub/internal/users"
)

// PresenceAttributes represents a user present in the area.
type PresenceAttributes struct {
	UserID        string `json:"user_id"`
	UserName      string `json:"user_name"`
	ItemID        string `json:"item_id"`
	ItemName      string `json:"item_name"`
	ItemGroupID   string `json:"item_group_id"`
	ItemGroupName string `json:"item_group_name"`
	Note          string `json:"note"`
}

// PresenceHandler returns a JSON:API list of users present in an area on a given date.
func PresenceHandler(cfg *spaces.Config, store *sql.DB) echo.HandlerFunc {
	return PresenceHandlerDynamic(func() *spaces.Config { return cfg }, store)
}

// PresenceHandlerDynamic returns a JSON:API list of users present in an area on a given date.
func PresenceHandlerDynamic(getConfig spaces.ConfigGetter, store *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		cfg := getConfig()
		areaID := c.Param("area_id")
		area, ok := cfg.FindArea(areaID)
		if !ok {
			return api.WriteNotFound(c, "Area not found")
		}

		params, err := api.ParseItemGroupRequest(areaID, c.QueryParam("date"))
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

// itemDetails holds item information for presence display.
type itemDetails struct {
	ItemName      string
	ItemGroupID   string
	ItemGroupName string
}

// buildItemIndex creates a map of item IDs to their details for an area.
func buildItemIndex(area *spaces.Area) (itemIDs []string, itemInfo map[string]itemDetails) {
	var totalItems int
	for _, ig := range area.ItemGroups {
		totalItems += len(ig.Items)
	}

	itemIDs = make([]string, 0, totalItems)
	itemInfo = make(map[string]itemDetails, totalItems)

	for _, ig := range area.ItemGroups {
		for _, item := range ig.Items {
			itemIDs = append(itemIDs, item.ID)
			itemInfo[item.ID] = itemDetails{
				ItemName:      item.Name,
				ItemGroupID:   ig.ID,
				ItemGroupName: ig.Name,
			}
		}
	}

	return itemIDs, itemInfo
}

func findAreaPresence(
	ctx context.Context, store *sql.DB, area *spaces.Area, bookingDate string,
) ([]api.Resource, error) {
	itemIDs, itemInfo := buildItemIndex(area)
	if len(itemIDs) == 0 {
		return []api.Resource{}, nil
	}

	placeholders, args := api.BuildINClause(itemIDs)
	args = append(args, bookingDate)

	//nolint:gosec // G201: placeholders are "?" literals from BuildINClause, not user input
	query := fmt.Sprintf(
		`SELECT id, item_id, user_id, note
		 FROM bookings
		 WHERE item_id IN (%s) AND booking_date = ?
		 ORDER BY item_id`,
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

	return scanPresenceRows(ctx, store, rows, itemInfo)
}

func scanPresenceRows(
	ctx context.Context, store *sql.DB, rows *sql.Rows, itemInfo map[string]itemDetails,
) ([]api.Resource, error) {
	type booking struct {
		bookingID string
		itemID    string
		userID    string
		note      string
	}

	var bookingList []booking
	userIDSet := make(map[string]struct{})

	for rows.Next() {
		var bookingID, itemID, userID, note string
		if err := rows.Scan(&bookingID, &itemID, &userID, &note); err != nil {
			return nil, fmt.Errorf("scan area presence: %w", err)
		}
		bookingList = append(bookingList, booking{bookingID: bookingID, itemID: itemID, userID: userID, note: note})
		userIDSet[userID] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate area presence: %w", err)
	}

	// Collect unique user IDs and look up display names.
	userIDs := make([]string, 0, len(userIDSet))
	for uid := range userIDSet {
		userIDs = append(userIDs, uid)
	}

	displayNames := make(map[string]string)
	if len(userIDs) > 0 {
		var err error
		displayNames, err = users.FindDisplayNames(ctx, store, userIDs)
		if err != nil {
			return nil, fmt.Errorf("find display names: %w", err)
		}
	}

	resources := make([]api.Resource, 0, len(bookingList))
	for _, b := range bookingList {
		info := itemInfo[b.itemID]
		resources = append(resources, api.Resource{
			Type: "presence",
			ID:   b.bookingID,
			Attributes: PresenceAttributes{
				UserID:        b.userID,
				UserName:      displayNames[b.userID],
				ItemID:        b.itemID,
				ItemName:      info.ItemName,
				ItemGroupID:   info.ItemGroupID,
				ItemGroupName: info.ItemGroupName,
				Note:          b.note,
			},
		})
	}

	return resources, nil
}
