package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// DeletionMode tells how DeleteOrDeactivate removed a user.
type DeletionMode string

const (
	// DeletionModeDeleted means the user row was removed (no booking referenced the user).
	DeletionModeDeleted DeletionMode = "deleted"
	// DeletionModeDeactivated means the user was soft-deleted to keep booking history intact.
	DeletionModeDeactivated DeletionMode = "deactivated"
)

// Impact previews what deleting a user would do.
type Impact struct {
	// HasBookings is true if any booking references the user, past or future.
	HasBookings bool
	// FutureBookingCount is the number of bookings from today onward that deletion cancels.
	FutureBookingCount int
}

// CanceledBooking is a booking removed while deactivating a user.
type CanceledBooking struct {
	ID          string
	ItemID      string
	UserID      string
	BookingDate string
	IsGuest     bool
	GuestName   string
	GuestEmail  string
}

// DeletionResult reports the outcome of DeleteOrDeactivate.
type DeletionResult struct {
	Mode             DeletionMode
	CanceledBookings []CanceledBooking
}

// Access is the per-request authorization state of a user as stored in the database.
type Access struct {
	// Active is false for deactivated users and for users that no longer exist.
	Active         bool
	IsAdmin        bool
	IsDataExporter bool
}

const (
	hasBookingsQuery = `SELECT EXISTS(SELECT 1 FROM bookings WHERE user_id = ? OR booked_by_user_id = ?)`

	// futureBookingsWhere selects the bookings canceled on deactivation: the user's own bookings
	// (including those a colleague made for them) and guest bookings the user hosts, from today on.
	// Bookings the user made for another colleague are kept. Args: today, id, id.
	futureBookingsWhere = `booking_date >= ? AND (user_id = ? OR (is_guest = 1 AND booked_by_user_id = ?))`
)

// Today returns the current UTC date, matching how the bookings package defines "today".
func Today() string {
	return time.Now().UTC().Format(time.DateOnly)
}

// LookupAccess loads whether a user may still use the API and which roles they hold (admin, data exporter).
// It is the per-request session check: a missing row (hard-deleted user) or a deactivated user
// yields an inactive Access without error, so revocation takes effect on the next request.
func LookupAccess(ctx context.Context, db *sql.DB, id string) (Access, error) {
	var isAdmin, deleted, dataExporter int
	err := db.QueryRowContext(ctx, `SELECT is_admin, deleted, is_data_exporter FROM users WHERE id = ?`, id).
		Scan(&isAdmin, &deleted, &dataExporter)
	if errors.Is(err, sql.ErrNoRows) {
		return Access{}, nil
	}
	if err != nil {
		return Access{}, fmt.Errorf("lookup user access: %w", err)
	}
	return Access{Active: deleted == 0, IsAdmin: isAdmin == 1, IsDataExporter: dataExporter == 1}, nil
}

// IsActive reports whether a user exists and is not deactivated.
// Missing hard-deleted users are inactive without error so stale sessions are revoked immediately.
func IsActive(ctx context.Context, db *sql.DB, id string) (bool, error) {
	access, err := LookupAccess(ctx, db, id)
	if err != nil {
		return false, err
	}
	return access.Active, nil
}

// DeletionImpact previews DeleteOrDeactivate for the given user and date ("YYYY-MM-DD").
// Returns ErrUserNotFound for an unknown user.
func DeletionImpact(ctx context.Context, db *sql.DB, id, today string) (Impact, error) {
	if _, err := FindByID(ctx, db, id); err != nil {
		return Impact{}, err
	}
	var impact Impact
	if err := db.QueryRowContext(ctx, hasBookingsQuery, id, id).Scan(&impact.HasBookings); err != nil {
		return Impact{}, fmt.Errorf("check user bookings: %w", err)
	}
	err := db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM bookings WHERE `+futureBookingsWhere, today, id, id).Scan(&impact.FutureBookingCount)
	if err != nil {
		return Impact{}, fmt.Errorf("count future bookings: %w", err)
	}
	return impact, nil
}

// DeleteOrDeactivate removes a user. A user no booking references is deleted from the database.
// Otherwise the user is deactivated (soft-deleted) so booking history keeps its name, and their
// bookings from today onward are canceled (see futureBookingsWhere). Everything runs in one
// transaction so a booking created concurrently cannot end up referencing a removed user.
// Returns ErrUserNotFound for an unknown user. Deactivating an already deactivated user is a no-op.
func DeleteOrDeactivate(ctx context.Context, db *sql.DB, id, today string) (result DeletionResult, err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return DeletionResult{}, fmt.Errorf("begin delete user: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback() //nolint:errcheck // Rollback after failure, original error is returned
		}
	}()

	result, err = deleteOrDeactivateTx(ctx, tx, id, today)
	if err != nil {
		return DeletionResult{}, err
	}
	if err = tx.Commit(); err != nil {
		return DeletionResult{}, fmt.Errorf("commit delete user: %w", err)
	}
	return result, nil
}

func deleteOrDeactivateTx(ctx context.Context, tx *sql.Tx, id, today string) (DeletionResult, error) {
	var exists bool
	if err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`, id).Scan(&exists); err != nil {
		return DeletionResult{}, fmt.Errorf("check user: %w", err)
	}
	if !exists {
		return DeletionResult{}, ErrUserNotFound
	}

	var hasBookings bool
	if err := tx.QueryRowContext(ctx, hasBookingsQuery, id, id).Scan(&hasBookings); err != nil {
		return DeletionResult{}, fmt.Errorf("check user bookings: %w", err)
	}
	if !hasBookings {
		if _, err := tx.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id); err != nil {
			return DeletionResult{}, fmt.Errorf("delete user: %w", err)
		}
		return DeletionResult{Mode: DeletionModeDeleted, CanceledBookings: []CanceledBooking{}}, nil
	}

	now := time.Now().UTC().Format(time.RFC3339)
	if _, err := tx.ExecContext(ctx, `UPDATE users SET deleted = 1, updated_at = ? WHERE id = ?`, now, id); err != nil {
		return DeletionResult{}, fmt.Errorf("deactivate user: %w", err)
	}

	canceled, err := selectFutureBookings(ctx, tx, id, today)
	if err != nil {
		return DeletionResult{}, err
	}
	for _, b := range canceled {
		if _, err := tx.ExecContext(ctx, `DELETE FROM bookings WHERE id = ?`, b.ID); err != nil {
			return DeletionResult{}, fmt.Errorf("cancel booking %s: %w", b.ID, err)
		}
	}
	return DeletionResult{Mode: DeletionModeDeactivated, CanceledBookings: canceled}, nil
}

func selectFutureBookings(ctx context.Context, tx *sql.Tx, id, today string) ([]CanceledBooking, error) {
	rows, err := tx.QueryContext(ctx,
		`SELECT id, item_id, user_id, booking_date, is_guest, guest_name, guest_email
		FROM bookings WHERE `+futureBookingsWhere+` ORDER BY booking_date, id`, today, id, id)
	if err != nil {
		return nil, fmt.Errorf("select future bookings: %w", err)
	}
	defer func() {
		_ = rows.Close() //nolint:errcheck // Defer close, error not critical
	}()

	result := []CanceledBooking{}
	for rows.Next() {
		var b CanceledBooking
		var isGuest int
		err := rows.Scan(&b.ID, &b.ItemID, &b.UserID, &b.BookingDate, &isGuest, &b.GuestName, &b.GuestEmail)
		if err != nil {
			return nil, fmt.Errorf("scan future booking: %w", err)
		}
		b.IsGuest = isGuest == 1
		result = append(result, b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate future bookings: %w", err)
	}
	return result, nil
}
