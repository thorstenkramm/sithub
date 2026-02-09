// Package items provides item handlers.
package items

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/auth"
	"github.com/thorstenkramm/sithub/internal/bookings"
	"github.com/thorstenkramm/sithub/internal/spaces"
	"github.com/thorstenkramm/sithub/internal/users"
)

// ListHandler returns a JSON:API list of items for an item group.
// For admin users, includes booking details (booking_id, booker_name) for occupied items.
func ListHandler(cfg *spaces.Config, store *sql.DB) echo.HandlerFunc {
	return ListHandlerDynamic(func() *spaces.Config { return cfg }, store)
}

// ListHandlerDynamic returns a JSON:API list of items for an item group using dynamic config.
func ListHandlerDynamic(getConfig spaces.ConfigGetter, store *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		cfg := getConfig()
		itemGroupID := c.Param("item_group_id")
		ig, ok := cfg.FindItemGroup(itemGroupID)
		if !ok {
			return api.WriteNotFound(c, "Item group not found")
		}

		bookingDate, err := api.ParseBookingDate(c.QueryParam("date"))
		if err != nil {
			return api.WriteBadRequest(c, "Invalid booking date. Use YYYY-MM-DD.")
		}

		ctx := c.Request().Context()
		user := auth.GetUserFromContext(c)
		isAdmin := user != nil && user.IsAdmin

		itemBookings, err := loadItemBookings(ctx, store, bookingDate, isAdmin)
		if err != nil {
			return err
		}

		// Resolve booker display names for admin users
		if isAdmin {
			resolveBookerNames(ctx, store, itemBookings)
		}

		resources := buildItemResources(ig.Items, itemBookings, isAdmin)
		return api.WriteCollection(c, resources, "write items response")
	}
}

func loadItemBookings(
	ctx context.Context, store *sql.DB, bookingDate string, isAdmin bool,
) (map[string]bookings.ItemBookingInfo, error) {
	if isAdmin {
		info, err := bookings.FindItemBookings(ctx, store, bookingDate)
		if err != nil {
			return nil, fmt.Errorf("list item bookings: %w", err)
		}
		return info, nil
	}

	bookedItems, err := bookings.FindBookedItemIDs(ctx, store, bookingDate)
	if err != nil {
		return nil, fmt.Errorf("list booked items: %w", err)
	}
	// Convert to ItemBookingInfo map for uniform handling
	result := make(map[string]bookings.ItemBookingInfo, len(bookedItems))
	for itemID := range bookedItems {
		result[itemID] = bookings.ItemBookingInfo{}
	}
	return result, nil
}

// resolveBookerNames looks up display names for all user IDs in the bookings map
// and populates the BookerName field. Lookup errors are silently ignored to avoid
// breaking the items response for a non-critical display field.
func resolveBookerNames(
	ctx context.Context, store *sql.DB, itemBookings map[string]bookings.ItemBookingInfo,
) {
	userIDs := make([]string, 0, len(itemBookings))
	for _, info := range itemBookings {
		if info.UserID != "" {
			userIDs = append(userIDs, info.UserID)
		}
	}
	if len(userIDs) == 0 {
		return
	}
	names, err := users.FindDisplayNames(ctx, store, userIDs)
	if err != nil {
		return
	}
	for itemID, info := range itemBookings {
		if name, ok := names[info.UserID]; ok {
			info.BookerName = name
			itemBookings[itemID] = info
		}
	}
}

func buildItemResources(
	items []spaces.Item, itemBookings map[string]bookings.ItemBookingInfo, isAdmin bool,
) []api.Resource {
	return api.MapResources(items, func(item spaces.Item) api.Resource {
		attrs := spaces.ItemAttributes(item.Name, item.Equipment, item.Warning, "")
		if info, booked := itemBookings[item.ID]; booked {
			attrs["availability"] = "occupied"
			if isAdmin {
				attrs["booking_id"] = info.BookingID
				attrs["booker_name"] = info.BookerName
			}
		} else {
			attrs["availability"] = "available"
		}
		return api.Resource{
			Type:       "items",
			ID:         item.ID,
			Attributes: attrs,
		}
	})
}
