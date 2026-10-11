package bookings

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/thorstenkramm/sithub/internal/api"
)

// ExportFilter selects the bookings listed by ListExportBookings.
type ExportFilter struct {
	// From and To are inclusive booking dates in YYYY-MM-DD form.
	From string
	To   string
	// ItemIDs restricts the result to these items. nil means all items; an empty non-nil slice matches nothing.
	ItemIDs []string
}

// ExportRow is one booking with the person it is for.
type ExportRow struct {
	ID          string
	ItemID      string
	BookingDate string
	// PersonUserID is the ID of the user the booking is for, usable to fetch their avatar. It is empty for
	// guest bookings (their user_id is synthetic) and when the user no longer exists.
	PersonUserID string
	// PersonName and PersonEmail belong to the person the booking is for: the guest for guest bookings,
	// otherwise the user_id user (not the booker). Empty if that user no longer exists.
	PersonName  string
	PersonEmail string
	IsGuest     bool
}

// exportPersonName and exportPersonEmail resolve the booked person in SQL so ordering by name works
// across pages. Guest bookings carry a synthetic user_id, hence the LEFT JOIN on users.
const (
	exportPersonUserID = `CASE WHEN b.is_guest = 1 THEN '' ELSE COALESCE(u.id, '') END`
	exportPersonName   = `CASE WHEN b.is_guest = 1 THEN b.guest_name ELSE COALESCE(u.display_name, '') END`
	exportPersonEmail  = `CASE WHEN b.is_guest = 1 THEN b.guest_email ELSE COALESCE(u.email, '') END`
)

// exportWhere builds the WHERE clause shared by the page and count queries so they cannot diverge.
// Soft-deleted users are deliberately not filtered out. The clause contains only hardcoded SQL and
// placeholders; item IDs are bound as arguments.
func exportWhere(filter ExportFilter) (clause string, args []any) {
	clause = "b.booking_date >= ? AND b.booking_date <= ?"
	args = []any{filter.From, filter.To}
	if filter.ItemIDs != nil {
		placeholders, idArgs := api.BuildINClause(filter.ItemIDs)
		clause += " AND b.item_id IN (" + placeholders + ")"
		args = append(args, idArgs...)
	}
	return clause, args
}

// ListExportBookings returns one page of bookings of all users within the filter, sorted by booking date,
// person name (case-insensitive) and booking id, plus the total number of matching bookings.
// A non-positive limit returns every matching booking and ignores offset (CSV export).
func ListExportBookings(
	ctx context.Context, store *sql.DB, filter ExportFilter, limit, offset int,
) (rows []ExportRow, total int, err error) {
	if filter.ItemIDs != nil && len(filter.ItemIDs) == 0 {
		return []ExportRow{}, 0, nil // IN () is invalid SQL and could not match anyway
	}
	where, args := exportWhere(filter)

	countQuery := "SELECT COUNT(*) FROM bookings b WHERE " + where
	if err := store.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count export bookings: %w", err)
	}

	query := strings.Join([]string{
		"SELECT b.id, b.item_id, b.booking_date, b.is_guest,",
		exportPersonUserID + " AS person_user_id,",
		exportPersonName + " AS person_name,",
		exportPersonEmail + " AS person_email",
		"FROM bookings b LEFT JOIN users u ON u.id = b.user_id",
		"WHERE " + where,
		"ORDER BY b.booking_date ASC, person_name COLLATE NOCASE ASC, b.id ASC",
	}, " ")
	pageArgs := append([]any{}, args...)
	if limit > 0 {
		query += " LIMIT ? OFFSET ?"
		pageArgs = append(pageArgs, limit, offset)
	}

	result, err := store.QueryContext(ctx, query, pageArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("query export bookings: %w", err)
	}
	defer func() {
		if closeErr := result.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close export rows: %w", closeErr)
		}
	}()

	rows = []ExportRow{}
	for result.Next() {
		var r ExportRow
		var isGuest int
		err := result.Scan(&r.ID, &r.ItemID, &r.BookingDate, &isGuest, &r.PersonUserID, &r.PersonName, &r.PersonEmail)
		if err != nil {
			return nil, 0, fmt.Errorf("scan export booking: %w", err)
		}
		r.IsGuest = isGuest == 1
		rows = append(rows, r)
	}
	if err := result.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate export bookings: %w", err)
	}
	return rows, total, nil
}
