// Package bookings provides booking lookup helpers.
package bookings

import (
	"context"
	"database/sql"
	"fmt"
)

// FindBookedDeskIDs returns the desk IDs with bookings on the given date.
func FindBookedDeskIDs(ctx context.Context, store *sql.DB, bookingDate string) (booked map[string]struct{}, err error) {
	rows, err := store.QueryContext(ctx, "SELECT desk_id FROM bookings WHERE booking_date = ?", bookingDate)
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
		var deskID string
		if err := rows.Scan(&deskID); err != nil {
			return nil, fmt.Errorf("scan booking: %w", err)
		}
		booked[deskID] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate bookings: %w", err)
	}

	return booked, nil
}

// BookingRecord represents a booking row from the database.
type BookingRecord struct {
	ID               string
	DeskID           string
	UserID           string
	BookingDate      string
	BookedByUserID   string
	BookedByUserName string
	IsGuest          bool
	GuestEmail       string
	CreatedAt        string
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
		query = `SELECT id, desk_id, user_id, booking_date, booked_by_user_id, booked_by_user_name, 
		                is_guest, guest_email, created_at 
		         FROM bookings 
		         WHERE (user_id = ? OR booked_by_user_id = ?) AND booking_date >= ? AND booking_date <= ?
		         ORDER BY booking_date DESC`
		args = []interface{}{userID, userID, fromDate, toDate}
	} else {
		query = `SELECT id, desk_id, user_id, booking_date, booked_by_user_id, booked_by_user_name, 
		                is_guest, guest_email, created_at 
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
			&b.ID, &b.DeskID, &b.UserID, &b.BookingDate,
			&b.BookedByUserID, &b.BookedByUserName, &isGuestInt, &b.GuestEmail, &b.CreatedAt,
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
		`SELECT id, desk_id, user_id, booking_date, booked_by_user_id, booked_by_user_name, 
		        is_guest, guest_email, created_at 
		 FROM bookings WHERE id = ?`,
		bookingID,
	).Scan(&b.ID, &b.DeskID, &b.UserID, &b.BookingDate, &b.BookedByUserID, &b.BookedByUserName,
		&isGuestInt, &b.GuestEmail, &b.CreatedAt)

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

// DeskBookingInfo contains booking details for a desk.
type DeskBookingInfo struct {
	BookingID  string
	UserID     string
	BookerName string
}

// FindDeskBookings returns booking info for desks on a given date, keyed by desk ID.
// Note: BookerName will be empty; the caller must look up display names separately if needed.
func FindDeskBookings(
	ctx context.Context, store *sql.DB, bookingDate string,
) (result map[string]DeskBookingInfo, err error) {
	query := `SELECT id, desk_id, user_id FROM bookings WHERE booking_date = ?`

	rows, err := store.QueryContext(ctx, query, bookingDate)
	if err != nil {
		return nil, fmt.Errorf("query desk bookings: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close desk bookings rows: %w", closeErr)
		}
	}()

	result = make(map[string]DeskBookingInfo)
	for rows.Next() {
		var info DeskBookingInfo
		var deskID string
		if err := rows.Scan(&info.BookingID, &deskID, &info.UserID); err != nil {
			return nil, fmt.Errorf("scan desk booking: %w", err)
		}
		result[deskID] = info
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate desk bookings: %w", err)
	}

	return result, nil
}
