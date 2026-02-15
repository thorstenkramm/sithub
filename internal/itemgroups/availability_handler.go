package itemgroups

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/spaces"
)

// DayAvailability holds availability data for a single day within an item group.
type DayAvailability struct {
	Date      string `json:"date"`
	Weekday   string `json:"weekday"`
	Total     int    `json:"total"`
	Available int    `json:"available"`
}

// ItemGroupAvailabilityAttributes holds per-item-group weekly availability.
type ItemGroupAvailabilityAttributes struct {
	ItemGroupID   string            `json:"item_group_id"`
	ItemGroupName string            `json:"item_group_name"`
	Days          []DayAvailability `json:"days"`
}

// AvailabilityHandler returns weekly availability for item groups in an area.
func AvailabilityHandler(cfg *spaces.Config, store *sql.DB) echo.HandlerFunc {
	return AvailabilityHandlerDynamic(func() *spaces.Config { return cfg }, store)
}

// AvailabilityHandlerDynamic returns weekly availability using dynamic config.
func AvailabilityHandlerDynamic(getConfig spaces.ConfigGetter, store *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		cfg := getConfig()
		areaID := c.Param("area_id")
		area, ok := cfg.FindArea(areaID)
		if !ok {
			return api.WriteNotFound(c, "Area not found")
		}

		weekParam := c.QueryParam("week")
		monday, err := parseISOWeek(weekParam)
		if err != nil {
			return api.WriteBadRequest(c, "Invalid week parameter. Use ISO 8601 format: YYYY-Www (e.g., 2026-W12).")
		}

		dayCount := 5
		if daysParam := c.QueryParam("days"); daysParam == "7" {
			dayCount = 7
		}
		weekdays := weekdayDates(monday, dayCount)
		ctx := c.Request().Context()

		resources, err := buildAvailabilityResources(ctx, store, area, weekdays)
		if err != nil {
			return fmt.Errorf("build availability: %w", err)
		}

		return api.WriteCollection(c, resources, "write availability response")
	}
}

// parseISOWeek parses an ISO 8601 week string (e.g., "2026-W12") and returns
// the Monday of that week. If the input is empty, returns the Monday of the
// current week.
func parseISOWeek(s string) (time.Time, error) {
	if strings.TrimSpace(s) == "" {
		now := time.Now()
		return mondayOfWeek(now), nil
	}

	// Expected format: YYYY-Www (e.g., 2026-W12 or 2026-W03)
	parts := strings.SplitN(s, "-W", 2)
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid ISO week format: %s", s)
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid year in ISO week: %s", s)
	}

	week, err := strconv.Atoi(parts[1])
	if err != nil || week < 1 || week > 53 {
		return time.Time{}, fmt.Errorf("invalid week number in ISO week: %s", s)
	}
	maxWeek := isoWeeksInYear(year)
	if week > maxWeek {
		return time.Time{}, fmt.Errorf("invalid week number in ISO week: %s", s)
	}

	// ISO 8601: Week 1 contains January 4th.
	// Find January 4th, then back up to Monday of that week, then add (week-1) weeks.
	jan4 := time.Date(year, time.January, 4, 0, 0, 0, 0, time.UTC)
	jan4Monday := mondayOfWeek(jan4)
	monday := jan4Monday.AddDate(0, 0, (week-1)*7)

	return monday, nil
}

// mondayOfWeek returns the Monday of the ISO week containing t.
func mondayOfWeek(t time.Time) time.Time {
	year, month, day := t.Date()
	base := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	weekday := base.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	offset := int(weekday) - int(time.Monday)
	return base.AddDate(0, 0, -offset)
}

// weekdayDates returns dates for a week starting at monday.
// count is typically 5 (Mon-Fri) or 7 (Mon-Sun).
func weekdayDates(monday time.Time, count int) []time.Time {
	days := make([]time.Time, count)
	for i := range count {
		days[i] = monday.AddDate(0, 0, i)
	}
	return days
}

// weekdayAbbreviation returns a two-letter weekday abbreviation.
func weekdayAbbreviation(d time.Weekday) string {
	switch d {
	case time.Monday:
		return "MO"
	case time.Tuesday:
		return "TU"
	case time.Wednesday:
		return "WE"
	case time.Thursday:
		return "TH"
	case time.Friday:
		return "FR"
	case time.Saturday:
		return "SA"
	case time.Sunday:
		return "SU"
	default:
		return d.String()[:2]
	}
}

// isoWeeksInYear returns the number of ISO weeks in the given year (52 or 53).
func isoWeeksInYear(year int) int {
	// ISO 8601: week with Jan 4 is always week 1; Dec 28 always lies in the last ISO week.
	dec28 := time.Date(year, time.December, 28, 0, 0, 0, 0, time.UTC)
	_, week := dec28.ISOWeek()
	return week
}

func buildAvailabilityResources(
	ctx context.Context, store *sql.DB, area *spaces.Area, weekdays []time.Time,
) ([]api.Resource, error) {
	// Collect all item IDs for this area to query bookings.
	var totalItems int
	for _, ig := range area.ItemGroups {
		totalItems += len(ig.Items)
	}
	allItemIDs := make([]string, 0, totalItems)
	itemGroupItems := make(map[string][]string, len(area.ItemGroups))
	for _, ig := range area.ItemGroups {
		igItemIDs := make([]string, 0, len(ig.Items))
		for _, item := range ig.Items {
			allItemIDs = append(allItemIDs, item.ID)
			igItemIDs = append(igItemIDs, item.ID)
		}
		itemGroupItems[ig.ID] = igItemIDs
	}

	// Query booking counts per item per day for the week.
	bookingCounts, err := countBookingsPerItemPerDay(ctx, store, allItemIDs, weekdays)
	if err != nil {
		return nil, err
	}

	resources := make([]api.Resource, 0, len(area.ItemGroups))
	for _, ig := range area.ItemGroups {
		igItemIDs := itemGroupItems[ig.ID]
		totalItems := len(igItemIDs)

		days := make([]DayAvailability, len(weekdays))
		for i, day := range weekdays {
			dateStr := day.Format(time.DateOnly)
			bookedCount := 0
			for _, itemID := range igItemIDs {
				if bookingCounts[itemID+"|"+dateStr] > 0 {
					bookedCount++
				}
			}
			available := totalItems - bookedCount
			if available < 0 {
				available = 0
			}
			days[i] = DayAvailability{
				Date:      dateStr,
				Weekday:   weekdayAbbreviation(day.Weekday()),
				Total:     totalItems,
				Available: available,
			}
		}

		resources = append(resources, api.Resource{
			Type: "item-group-availability",
			ID:   ig.ID,
			Attributes: ItemGroupAvailabilityAttributes{
				ItemGroupID:   ig.ID,
				ItemGroupName: ig.Name,
				Days:          days,
			},
		})
	}

	return resources, nil
}

// countBookingsPerItemPerDay returns a map of "itemID|date" -> booking count.
func countBookingsPerItemPerDay(
	ctx context.Context, store *sql.DB, itemIDs []string, weekdays []time.Time,
) (map[string]int, error) {
	if len(itemIDs) == 0 || len(weekdays) == 0 {
		return make(map[string]int), nil
	}

	itemPlaceholders, itemArgs := api.BuildINClause(itemIDs)

	dateStrings := make([]string, len(weekdays))
	for i, d := range weekdays {
		dateStrings[i] = d.Format(time.DateOnly)
	}
	datePlaceholders, dateArgs := api.BuildINClause(dateStrings)

	args := make([]any, 0, len(itemArgs)+len(dateArgs))
	args = append(args, itemArgs...)
	args = append(args, dateArgs...)

	//nolint:gosec // G201: placeholders are "?" literals from BuildINClause
	query := fmt.Sprintf(
		`SELECT item_id, booking_date, COUNT(*) as cnt
		 FROM bookings
		 WHERE item_id IN (%s) AND booking_date IN (%s)
		 GROUP BY item_id, booking_date`,
		itemPlaceholders, datePlaceholders,
	)

	rows, err := store.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query booking counts: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			_ = err
		}
	}()

	result := make(map[string]int)
	for rows.Next() {
		var itemID, bookingDate string
		var cnt int
		if err := rows.Scan(&itemID, &bookingDate, &cnt); err != nil {
			return nil, fmt.Errorf("scan booking count: %w", err)
		}
		result[itemID+"|"+bookingDate] = cnt
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate booking counts: %w", err)
	}

	return result, nil
}
