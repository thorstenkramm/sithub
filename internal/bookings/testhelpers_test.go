package bookings

import (
	"database/sql"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/db"
)

func setupTestStore(t *testing.T) *sql.DB {
	t.Helper()

	dir := t.TempDir()
	store, err := db.Open(dir)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, store.Close())
	})

	migrationsPath := resolveTestMigrationsPath(t)
	require.NoError(t, db.RunMigrations(store, migrationsPath))

	return store
}

func resolveTestMigrationsPath(t *testing.T) string {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok)

	root := filepath.Clean(filepath.Join(filepath.Dir(filename), "..", ".."))
	return filepath.Join(root, "migrations")
}

func seedTestBooking(t *testing.T, store *sql.DB, bookingID, deskID, userID, bookingDate string) {
	seedTestBookingFull(t, store, bookingID, deskID, userID, "Test User", userID, "Test User", bookingDate)
}

func seedTestBookingFull(
	t *testing.T, store *sql.DB, bookingID, deskID,
	userID, userName, bookedByUserID, bookedByUserName, bookingDate string,
) {
	seedTestBookingWithGuest(t, store, bookingID, deskID, userID, userName,
		bookedByUserID, bookedByUserName, bookingDate, false, "")
}

func seedTestBookingWithGuest(
	t *testing.T, store *sql.DB, bookingID, deskID,
	userID, userName, bookedByUserID, bookedByUserName, bookingDate string,
	isGuest bool, guestEmail string,
) {
	t.Helper()

	now := time.Now().UTC().Format(time.RFC3339)
	isGuestInt := 0
	if isGuest {
		isGuestInt = 1
	}
	_, err := store.Exec(`
		INSERT INTO bookings 
		(id, desk_id, user_id, user_name, booked_by_user_id, booked_by_user_name, 
		 booking_date, is_guest, guest_email, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		bookingID, deskID, userID, userName, bookedByUserID, bookedByUserName,
		bookingDate, isGuestInt, guestEmail, now, now,
	)
	require.NoError(t, err)
}
