package itemgroups

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/areas"
	"github.com/thorstenkramm/sithub/internal/auth"
	"github.com/thorstenkramm/sithub/internal/bookings"
	"github.com/thorstenkramm/sithub/internal/users"
)

// MatrixDayMeta describes a single visible day in the matrix header.
type MatrixDayMeta struct {
	Date    string `json:"date"`
	Weekday string `json:"weekday"`
}

// MatrixCell holds booking state for one item on one day.
type MatrixCell struct {
	Date         string `json:"date"`
	Availability string `json:"availability"`
	BookerName   string `json:"booker_name,omitempty"`
	BookerUserID string `json:"booker_user_id,omitempty"`
	BookedByMe   bool   `json:"booked_by_me"`
	BookingID    string `json:"booking_id,omitempty"`
}

// MatrixItem holds metadata and cells for a single item row.
type MatrixItem struct {
	ItemID    string       `json:"item_id"`
	ItemName  string       `json:"item_name"`
	Equipment []string     `json:"equipment"`
	Warning   string       `json:"warning,omitempty"`
	Reserved  bool         `json:"reserved,omitempty"`
	Cells     []MatrixCell `json:"cells"`
}

// MatrixAttributes holds the attributes for an item-group-weekly-matrix resource.
type MatrixAttributes struct {
	ItemGroupID   string          `json:"item_group_id"`
	ItemGroupName string          `json:"item_group_name"`
	Days          []MatrixDayMeta `json:"days"`
	Items         []MatrixItem    `json:"items"`
}

const matrixResourceType = "item-group-weekly-matrix"

// MatrixHandler returns a weekly desk matrix for item groups in an area.
func MatrixHandler(cfg *areas.Config, store *sql.DB) echo.HandlerFunc {
	return MatrixHandlerDynamic(func() *areas.Config { return cfg }, store)
}

// MatrixHandlerDynamic returns a weekly desk matrix using dynamic config.
func MatrixHandlerDynamic(getConfig areas.ConfigGetter, store *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		cfg := getConfig()
		areaID := c.Param("area_id")
		area, ok := cfg.FindArea(areaID)
		if !ok {
			return api.WriteNotFound(c, "Area not found")
		}

		monday, err := parseISOWeek(c.QueryParam("week"))
		if err != nil {
			return api.WriteBadRequest(c, "Invalid week parameter. Use ISO 8601 format: YYYY-Www (e.g., 2026-W12).")
		}

		dayCount := 5
		if c.QueryParam("days") == "7" {
			dayCount = 7
		}
		weekdays := weekdayDates(monday, dayCount)

		ctx := c.Request().Context()
		user := auth.GetUserFromContext(c)
		isAdmin := user != nil && user.IsAdmin

		currentUserID, userEmail := resolveMatrixUser(ctx, store, user)

		resources, err := buildMatrixResources(ctx, store, area, weekdays, isAdmin, currentUserID, userEmail)
		if err != nil {
			return fmt.Errorf("build matrix: %w", err)
		}

		return api.WriteCollection(c, resources, "write matrix response")
	}
}

// resolveMatrixUser returns the current user's ID and email from the database.
func resolveMatrixUser(ctx context.Context, store *sql.DB, user *auth.User) (userID, email string) {
	if user == nil {
		return "", ""
	}
	rec, err := users.FindByID(ctx, store, user.ID)
	if err != nil || rec == nil {
		return user.ID, ""
	}
	return user.ID, rec.Email
}

func buildMatrixResources(
	ctx context.Context, store *sql.DB, area *areas.Area, weekdays []time.Time,
	isAdmin bool, currentUserID, userEmail string,
) ([]api.Resource, error) {
	// Collect all item IDs for one batch query.
	allItemIDs := collectAreaItemIDs(area)

	// Build date strings for the visible week.
	dateStrings := make([]string, len(weekdays))
	for i, d := range weekdays {
		dateStrings[i] = d.Format(time.DateOnly)
	}

	// One query: fetch all bookings for these items and dates.
	matrixBookings, err := bookings.FindMatrixBookings(ctx, store, allItemIDs, dateStrings)
	if err != nil {
		return nil, fmt.Errorf("find matrix bookings: %w", err)
	}

	// Resolve display names for all booker user IDs.
	resolveMatrixBookerNames(ctx, store, matrixBookings)

	// Build day metadata (shared across all groups).
	daysMeta := make([]MatrixDayMeta, len(weekdays))
	for i, d := range weekdays {
		daysMeta[i] = MatrixDayMeta{
			Date:    dateStrings[i],
			Weekday: weekdayAbbreviation(d.Weekday()),
		}
	}

	// Build one resource per item group in configured order.
	resources := make([]api.Resource, 0, len(area.ItemGroups))
	for i := range area.ItemGroups {
		ig := &area.ItemGroups[i]
		items := buildMatrixItems(ig, area, matrixBookings, dateStrings, isAdmin, currentUserID, userEmail)

		resources = append(resources, api.Resource{
			Type: matrixResourceType,
			ID:   ig.ID,
			Attributes: MatrixAttributes{
				ItemGroupID:   ig.ID,
				ItemGroupName: ig.Name,
				Days:          daysMeta,
				Items:         items,
			},
		})
	}

	return resources, nil
}

func collectAreaItemIDs(area *areas.Area) []string {
	var total int
	for i := range area.ItemGroups {
		total += len(area.ItemGroups[i].Items)
	}
	ids := make([]string, 0, total)
	for i := range area.ItemGroups {
		for _, item := range area.ItemGroups[i].Items {
			ids = append(ids, item.ID)
		}
	}
	return ids
}

// resolveMatrixBookerNames populates BookerName for all matrix bookings.
func resolveMatrixBookerNames(
	ctx context.Context, store *sql.DB, mb map[string]bookings.MatrixBookingInfo,
) {
	userIDs := make([]string, 0, len(mb))
	for _, info := range mb {
		if !info.IsGuest && info.UserID != "" {
			userIDs = append(userIDs, info.UserID)
		}
	}

	var names map[string]string
	if len(userIDs) > 0 {
		var err error
		names, err = users.FindDisplayNames(ctx, store, userIDs)
		if err != nil {
			names = make(map[string]string)
		}
	}

	for key, info := range mb {
		if info.IsGuest {
			info.BookerName = info.GuestName
		} else if name, ok := names[info.UserID]; ok {
			info.BookerName = name
		}
		mb[key] = info
	}
}

func buildMatrixItems(
	ig *areas.ItemGroup, parentArea *areas.Area,
	mb map[string]bookings.MatrixBookingInfo, dateStrings []string,
	isAdmin bool, currentUserID, userEmail string,
) []MatrixItem {
	items := make([]MatrixItem, 0, len(ig.Items))
	for j := range ig.Items {
		item := &ig.Items[j]

		// Check reservation at item level.
		reserved := false
		if userEmail != "" {
			loc := &areas.ItemLocation{Area: parentArea, ItemGroup: ig, Item: item}
			reserved = areas.IsReserved(loc, userEmail)
		}

		cells := buildMatrixCells(item.ID, mb, dateStrings, isAdmin, currentUserID)

		equip := item.Equipment
		if equip == nil {
			equip = []string{}
		}

		items = append(items, MatrixItem{
			ItemID:    item.ID,
			ItemName:  item.Name,
			Equipment: equip,
			Warning:   item.Warning,
			Reserved:  reserved,
			Cells:     cells,
		})
	}
	return items
}

func buildMatrixCells(
	itemID string, mb map[string]bookings.MatrixBookingInfo,
	dateStrings []string, isAdmin bool, currentUserID string,
) []MatrixCell {
	cells := make([]MatrixCell, len(dateStrings))
	for i, dateStr := range dateStrings {
		key := itemID + "|" + dateStr
		info, occupied := mb[key]

		cell := MatrixCell{
			Date:         dateStr,
			Availability: "free",
		}

		if occupied {
			cell.Availability = "occupied"
			cell.BookerName = info.BookerName
			if !info.IsGuest {
				cell.BookerUserID = info.UserID
			}
			cell.BookedByMe = currentUserID != "" && info.UserID == currentUserID

			// Only expose booking_id to the booking owner or admins.
			if isAdmin || (currentUserID != "" && info.UserID == currentUserID) {
				cell.BookingID = info.BookingID
			}
		}

		cells[i] = cell
	}
	return cells
}
