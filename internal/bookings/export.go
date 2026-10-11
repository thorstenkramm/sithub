package bookings

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/areas"
)

// BookingExportAttributes is the export projection of a booking. Name and email belong to the person the
// booking is for (the guest for guest bookings, the colleague for on-behalf bookings).
type BookingExportAttributes struct {
	BookingDate string `json:"booking_date"`
	// UserID identifies the person for avatar lookup; empty for guests and removed users.
	UserID        string `json:"user_id"`
	UserName      string `json:"user_name"`
	UserEmail     string `json:"user_email"`
	IsGuest       bool   `json:"is_guest"`
	ItemID        string `json:"item_id"`
	ItemName      string `json:"item_name"`
	ItemGroupID   string `json:"item_group_id"`
	ItemGroupName string `json:"item_group_name"`
	AreaID        string `json:"area_id"`
	AreaName      string `json:"area_name"`
	// BookedItem is "<item group name> / <item name>" for display and CSV export.
	BookedItem string `json:"booked_item"`
}

// exportParams are the validated query parameters of the export endpoints.
type exportParams struct {
	filter ExportFilter
	page   api.Page
}

// FormatBookedItem renders the booked item as "<item group name> / <item name>".
// Without a known location (the item is no longer configured) the raw item ID is returned.
func FormatBookedItem(loc *areas.ItemLocation, itemID string) string {
	if loc == nil {
		return itemID
	}
	return loc.ItemGroup.Name + " / " + loc.Item.Name
}

// parseExportParams validates from, to, area and the pagination parameters of an export request.
// All failures are validation errors whose detail is shown to the API client unchanged.
func parseExportParams(c echo.Context, cfg *areas.Config) (exportParams, error) {
	return parseExportRequest(c, cfg, true)
}

// Export filters use the JSON:API filter[...] family. Plain names like "from" are removed from URLs by
// privacy browser extensions (tracking-parameter lists), which silently broke the export.
const (
	filterFromParam = "filter[from]"
	filterToParam   = "filter[to]"
	filterAreaParam = "filter[area]"
)

// parseExportRequest shares filter validation; CSV deliberately ignores all page parameters.
func parseExportRequest(c echo.Context, cfg *areas.Config, paginated bool) (exportParams, error) {
	from, err := requiredDateParam(c, filterFromParam)
	if err != nil {
		return exportParams{}, err
	}
	to, err := requiredDateParam(c, filterToParam)
	if err != nil {
		return exportParams{}, err
	}
	if from > to {
		return exportParams{}, errBadRequest(fmt.Sprintf("%q must not be after %q", filterFromParam, filterToParam))
	}

	filter := ExportFilter{From: from, To: to}
	if areaID := c.QueryParam(filterAreaParam); areaID != "" {
		area, ok := cfg.FindArea(areaID)
		if !ok {
			return exportParams{}, errBadRequest("unknown area: " + areaID)
		}
		filter.ItemIDs = collectAreaItemIDs(area)
		if filter.ItemIDs == nil {
			filter.ItemIDs = []string{} // an area without items matches no booking
		}
	}

	if !paginated {
		return exportParams{filter: filter}, nil
	}
	page, err := api.ParsePage(c)
	if err != nil {
		return exportParams{}, errBadRequest(err.Error())
	}
	return exportParams{filter: filter, page: page}, nil
}

// requiredDateParam reads a mandatory YYYY-MM-DD query parameter.
func requiredDateParam(c echo.Context, name string) (string, error) {
	value := c.QueryParam(name)
	if value == "" {
		return "", errBadRequest(fmt.Sprintf("%q is required", name))
	}
	if _, err := time.Parse(time.DateOnly, value); err != nil {
		return "", errBadRequest(fmt.Sprintf("invalid %q date, use YYYY-MM-DD", name))
	}
	return value, nil
}

// ExportHandlerDynamic returns a handler listing the bookings of all users for a date range and area.
//
// The handler is deliberately identity-agnostic: it reads only query parameters, configuration and the
// database. Authorization is applied by the route middleware, so the same handler can later be served to
// API keys by changing only that middleware.
func ExportHandlerDynamic(getConfig areas.ConfigGetter, store *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		cfg := getConfig()
		params, err := parseExportParams(c, cfg)
		if err != nil {
			var validation validationError
			if errors.As(err, &validation) {
				return api.WriteBadRequest(c, validation.detail)
			}
			return fmt.Errorf("parse export params: %w", err)
		}

		rows, total, err := ListExportBookings(
			c.Request().Context(), store, params.filter, params.page.Limit, params.page.Offset)
		if err != nil {
			return fmt.Errorf("list export bookings: %w", err)
		}

		warned := map[string]bool{}
		resources := api.MapResources(rows, func(row ExportRow) api.Resource {
			return api.Resource{
				Type:       resourceTypeBooking,
				ID:         row.ID,
				Attributes: exportAttributes(cfg, &row, warned),
			}
		})
		return api.WritePaginatedCollection(c, resources, params.page, total, "write booking export response")
	}
}

// exportAttributes maps a row to its export attributes. A booking whose item is no longer configured is
// kept (meta.total counts it) and shows the raw item ID; warned makes the log warning once per item.
func exportAttributes(cfg *areas.Config, row *ExportRow, warned map[string]bool) BookingExportAttributes {
	attrs := BookingExportAttributes{
		BookingDate: row.BookingDate,
		UserID:      row.PersonUserID,
		UserName:    row.PersonName,
		UserEmail:   row.PersonEmail,
		IsGuest:     row.IsGuest,
		ItemID:      row.ItemID,
	}
	loc, ok := cfg.FindItemLocation(row.ItemID)
	if !ok {
		if !warned[row.ItemID] {
			warned[row.ItemID] = true
			slog.Warn("booking references an item missing from the areas config", "item_id", row.ItemID)
		}
		attrs.ItemName = row.ItemID
		attrs.BookedItem = FormatBookedItem(nil, row.ItemID)
		return attrs
	}
	attrs.ItemName = loc.Item.Name
	attrs.ItemGroupID = loc.ItemGroup.ID
	attrs.ItemGroupName = loc.ItemGroup.Name
	attrs.AreaID = loc.Area.ID
	attrs.AreaName = loc.Area.Name
	attrs.BookedItem = FormatBookedItem(loc, row.ItemID)
	return attrs
}
