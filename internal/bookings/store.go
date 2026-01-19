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
	ID          string
	DeskID      string
	UserID      string
	BookingDate string
	CreatedAt   string
}

// ListUserBookings returns all bookings for a user on or after the given date, ordered by booking_date.
func ListUserBookings(ctx context.Context, store *sql.DB, userID, fromDate string) (result []BookingRecord, err error) {
	query := `SELECT id, desk_id, user_id, booking_date, created_at 
	          FROM bookings 
	          WHERE user_id = ? AND booking_date >= ? 
	          ORDER BY booking_date ASC`

	rows, err := store.QueryContext(ctx, query, userID, fromDate)
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
		if err := rows.Scan(&b.ID, &b.DeskID, &b.UserID, &b.BookingDate, &b.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan user booking: %w", err)
		}
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
	err := store.QueryRowContext(ctx,
		"SELECT id, desk_id, user_id, booking_date, created_at FROM bookings WHERE id = ?",
		bookingID,
	).Scan(&b.ID, &b.DeskID, &b.UserID, &b.BookingDate, &b.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query booking by id: %w", err)
	}
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
