// Package bookings provides booking lookup helpers.
package bookings

import (
	"context"
	"database/sql"
	"fmt"
)

// FindBookedItemIDs returns the item IDs with bookings on the given date.
func FindBookedItemIDs(ctx context.Context, store *sql.DB, bookingDate string) (booked map[string]struct{}, err error) {
	rows, err := store.QueryContext(ctx, "SELECT item_id FROM bookings WHERE booking_date = ?", bookingDate)
	if err != nil {
		return nil, fmt.Errorf("query bookings: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close bookings rows: %w", closeErr)
		}
	}()

	booked = make(map[string]struct{})
	for rows.Next() {
		var itemID string
		if err := rows.Scan(&itemID); err != nil {
			return nil, fmt.Errorf("scan booking: %w", err)
		}
		booked[itemID] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate bookings: %w", err)
	}

	return booked, nil
}

// BookingRecord represents a booking row from the database.
type BookingRecord struct {
	ID             string
	ItemID         string
	UserID         string
	BookingDate    string
	BookedByUserID string
	IsGuest        bool
	GuestName      string
	GuestEmail     string
	CreatedAt      string
	UpdatedAt      string
}

// ListUserBookings returns all bookings for a user on or after the given date, ordered by booking_date.
// Includes bookings where user_id matches OR booked_by_user_id matches.
func ListUserBookings(ctx context.Context, store *sql.DB, userID, fromDate string) (result []BookingRecord, err error) {
	return ListUserBookingsRange(ctx, store, userID, fromDate, "")
}

// ListUserBookingsRange returns bookings for a user within a date range, ordered by booking_date.
// If toDate is empty, returns all bookings from fromDate onwards.
// Includes bookings where user_id matches OR booked_by_user_id matches.
func ListUserBookingsRange(
	ctx context.Context, store *sql.DB, userID, fromDate, toDate string,
) (result []BookingRecord, err error) {
	var query string
	var args []interface{}

	if toDate != "" {
		query = `SELECT id, item_id, user_id, booking_date, booked_by_user_id,
		                is_guest, guest_name, guest_email, created_at, updated_at
		         FROM bookings
		         WHERE (user_id = ? OR booked_by_user_id = ?) AND booking_date >= ? AND booking_date <= ?
		         ORDER BY booking_date DESC`
		args = []interface{}{userID, userID, fromDate, toDate}
	} else {
		query = `SELECT id, item_id, user_id, booking_date, booked_by_user_id,
		                is_guest, guest_name, guest_email, created_at, updated_at
		         FROM bookings
		         WHERE (user_id = ? OR booked_by_user_id = ?) AND booking_date >= ?
		         ORDER BY booking_date ASC`
		args = []interface{}{userID, userID, fromDate}
	}

	rows, err := store.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query user bookings: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close user bookings rows: %w", closeErr)
		}
	}()

	for rows.Next() {
		var b BookingRecord
		var isGuestInt int
		err := rows.Scan(
			&b.ID, &b.ItemID, &b.UserID, &b.BookingDate,
			&b.BookedByUserID, &isGuestInt, &b.GuestName, &b.GuestEmail, &b.CreatedAt, &b.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan user booking: %w", err)
		}
		b.IsGuest = isGuestInt == 1
		result = append(result, b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate user bookings: %w", err)
	}

	return result, nil
}

// FindBookingByID returns a booking by its ID, or nil if not found.
func FindBookingByID(ctx context.Context, store *sql.DB, bookingID string) (*BookingRecord, error) {
	var b BookingRecord
	var isGuestInt int
	err := store.QueryRowContext(ctx,
		`SELECT id, item_id, user_id, booking_date, booked_by_user_id,
		        is_guest, guest_name, guest_email, created_at, updated_at
		 FROM bookings WHERE id = ?`,
		bookingID,
	).Scan(&b.ID, &b.ItemID, &b.UserID, &b.BookingDate, &b.BookedByUserID,
		&isGuestInt, &b.GuestName, &b.GuestEmail, &b.CreatedAt, &b.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query booking by id: %w", err)
	}
	b.IsGuest = isGuestInt == 1
	return &b, nil
}

// DeleteBooking removes a booking by its ID.
func DeleteBooking(ctx context.Context, store *sql.DB, bookingID string) error {
	_, err := store.ExecContext(ctx, "DELETE FROM bookings WHERE id = ?", bookingID)
	if err != nil {
		return fmt.Errorf("delete booking: %w", err)
	}
	return nil
}

// ItemBookingInfo contains booking details for an item.
type ItemBookingInfo struct {
	BookingID  string
	UserID     string
	BookerName string
}

// FindItemBookings returns booking info for items on a given date, keyed by item ID.
// Note: BookerName will be empty; the caller must look up display names separately if needed.
func FindItemBookings(
	ctx context.Context, store *sql.DB, bookingDate string,
) (result map[string]ItemBookingInfo, err error) {
	query := `SELECT id, item_id, user_id FROM bookings WHERE booking_date = ?`

	rows, err := store.QueryContext(ctx, query, bookingDate)
	if err != nil {
		return nil, fmt.Errorf("query item bookings: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close item bookings rows: %w", closeErr)
		}
	}()

	result = make(map[string]ItemBookingInfo)
	for rows.Next() {
		var info ItemBookingInfo
		var itemID string
		if err := rows.Scan(&info.BookingID, &itemID, &info.UserID); err != nil {
			return nil, fmt.Errorf("scan item booking: %w", err)
		}
		result[itemID] = info
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate item bookings: %w", err)
	}

	return result, nil
}
