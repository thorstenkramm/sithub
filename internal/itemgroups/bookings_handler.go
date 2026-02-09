package itemgroups

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/spaces"
	"github.com/thorstenkramm/sithub/internal/users"
)

// ItemGroupBookingAttributes represents a booking in the item group overview.
type ItemGroupBookingAttributes struct {
	ItemID      string `json:"item_id"`
	ItemName    string `json:"item_name"`
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	BookingDate string `json:"booking_date"`
	IsGuest     bool   `json:"is_guest,omitempty"`
}

// BookingsHandler returns a JSON:API list of bookings for an item group on a given date.
func BookingsHandler(cfg *spaces.Config, store *sql.DB) echo.HandlerFunc {
	return BookingsHandlerDynamic(func() *spaces.Config { return cfg }, store)
}

// BookingsHandlerDynamic returns a JSON:API list of bookings for an item group on a given date.
func BookingsHandlerDynamic(getConfig spaces.ConfigGetter, store *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		cfg := getConfig()
		ig, err := validateItemGroupParam(cfg, c.Param("item_group_id"))
		if err != nil {
			return api.WriteNotFound(c, "Item group not found")
		}

		params, err := api.ParseItemGroupRequest(ig.ID, c.QueryParam("date"))
		if err != nil {
			return api.WriteBadRequest(c, "Invalid booking date. Use YYYY-MM-DD.")
		}

		igBookings, err := findItemGroupBookings(c.Request().Context(), store, ig, params.BookingDate)
		if err != nil {
			return fmt.Errorf("find item group bookings: %w", err)
		}

		return api.WriteCollection(c, igBookings, "write item group bookings response")
	}
}

// validateItemGroupParam checks if an item group exists in the config.
func validateItemGroupParam(cfg *spaces.Config, itemGroupID string) (*spaces.ItemGroup, error) {
	ig, ok := cfg.FindItemGroup(itemGroupID)
	if !ok {
		return nil, fmt.Errorf("item group %s not found", itemGroupID)
	}
	return ig, nil
}

// itemGroupBookingRecord represents a booking row joined with item info.
type itemGroupBookingRecord struct {
	ID          string
	ItemID      string
	UserID      string
	BookingDate string
	IsGuest     bool
}

func findItemGroupBookings(
	ctx context.Context, store *sql.DB, ig *spaces.ItemGroup, bookingDate string,
) ([]api.Resource, error) {
	// Build list of item IDs in this item group
	itemIDs := make([]string, 0, len(ig.Items))
	itemNames := make(map[string]string, len(ig.Items))
	for _, item := range ig.Items {
		itemIDs = append(itemIDs, item.ID)
		itemNames[item.ID] = item.Name
	}

	if len(itemIDs) == 0 {
		return []api.Resource{}, nil
	}

	placeholders, args := api.BuildINClause(itemIDs)
	args = append(args, bookingDate)

	//nolint:gosec // G201: placeholders are "?" literals from BuildINClause, not user input
	query := fmt.Sprintf(
		`SELECT id, item_id, user_id, booking_date, is_guest
		 FROM bookings
		 WHERE item_id IN (%s) AND booking_date = ?
		 ORDER BY item_id`,
		placeholders,
	)

	rows, err := store.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query item group bookings: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			// Log error but don't override main error
			_ = closeErr
		}
	}()

	var records []itemGroupBookingRecord
	var userIDs []string
	for rows.Next() {
		var rec itemGroupBookingRecord
		var isGuestInt int
		err := rows.Scan(
			&rec.ID, &rec.ItemID, &rec.UserID, &rec.BookingDate, &isGuestInt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan item group booking: %w", err)
		}
		rec.IsGuest = isGuestInt == 1
		records = append(records, rec)
		userIDs = append(userIDs, rec.UserID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate item group bookings: %w", err)
	}

	// Look up display names from users table
	displayNames, err := users.FindDisplayNames(ctx, store, userIDs)
	if err != nil {
		return nil, fmt.Errorf("find display names: %w", err)
	}

	resources := make([]api.Resource, 0, len(records))
	for _, rec := range records {
		resources = append(resources, api.Resource{
			Type: "bookings",
			ID:   rec.ID,
			Attributes: ItemGroupBookingAttributes{
				ItemID:      rec.ItemID,
				ItemName:    itemNames[rec.ItemID],
				UserID:      rec.UserID,
				UserName:    displayNames[rec.UserID],
				BookingDate: rec.BookingDate,
				IsGuest:     rec.IsGuest,
			},
		})
	}

	return resources, nil
}
