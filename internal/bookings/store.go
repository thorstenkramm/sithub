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
