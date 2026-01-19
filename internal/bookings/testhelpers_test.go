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

// seedTestDeskData creates test area, room, and desks. Uses fixed area-1 and room-1.
func seedTestDeskData(t *testing.T, store *sql.DB, deskIDs []string) {
	t.Helper()

	const areaID = "area-1"
	const roomID = "room-1"
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := store.Exec(
		"INSERT INTO areas (id, name, created_at, updated_at) VALUES (?, ?, ?, ?)",
		areaID,
		"Area",
		now,
		now,
	)
	require.NoError(t, err)

	_, err = store.Exec(
		"INSERT INTO rooms (id, area_id, name, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		roomID,
		areaID,
		"Room",
		now,
		now,
	)
	require.NoError(t, err)

	for _, deskID := range deskIDs {
		_, err = store.Exec(
			"INSERT INTO desks (id, room_id, name, equipment, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
			deskID,
			roomID,
			deskID,
			"",
			now,
			now,
		)
		require.NoError(t, err)
	}
}

func seedTestBooking(t *testing.T, store *sql.DB, bookingID, deskID, userID, bookingDate string) {
	seedTestBookingFull(t, store, bookingID, deskID, userID, "Test User", userID, "Test User", bookingDate)
}

func seedTestBookingFull(
	t *testing.T, store *sql.DB, bookingID, deskID,
	userID, userName, bookedByUserID, bookedByUserName, bookingDate string,
) {
	t.Helper()

	now := time.Now().UTC().Format(time.RFC3339)
	_, err := store.Exec(`
		INSERT INTO bookings 
		(id, desk_id, user_id, user_name, booked_by_user_id, booked_by_user_name, booking_date, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		bookingID, deskID, userID, userName, bookedByUserID, bookedByUserName, bookingDate, now, now,
	)
	require.NoError(t, err)
}
