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

func TestFindBookedDeskIDs(t *testing.T) {
	store := setupStore(t)
	seedDeskData(t, store, "area-1", "room-1", "desk-1")
	seedBooking(t, store, "booking-1", "desk-1", "user-1", "2026-01-20")

	booked, err := FindBookedDeskIDs(t.Context(), store, "2026-01-20")
	require.NoError(t, err)
	require.Contains(t, booked, "desk-1")
}

func TestFindBookedDeskIDsReturnsErrorOnClosedDB(t *testing.T) {
	store, err := db.Open(t.TempDir())
	require.NoError(t, err)
	require.NoError(t, store.Close())

	_, err = FindBookedDeskIDs(t.Context(), store, "2026-01-20")
	require.Error(t, err)
}

func setupStore(t *testing.T) *sql.DB {
	t.Helper()

	dir := t.TempDir()
	store, err := db.Open(dir)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, store.Close())
	})

	migrationsPath := resolveMigrationsPath(t)
	require.NoError(t, db.RunMigrations(store, migrationsPath))

	return store
}

func resolveMigrationsPath(t *testing.T) string {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok)

	root := filepath.Clean(filepath.Join(filepath.Dir(filename), "..", ".."))
	return filepath.Join(root, "migrations")
}

func seedDeskData(t *testing.T, store *sql.DB, areaID, roomID, deskID string) {
	t.Helper()

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

func seedBooking(t *testing.T, store *sql.DB, bookingID, deskID, userID, bookingDate string) {
	t.Helper()

	now := time.Now().UTC().Format(time.RFC3339)
	_, err := store.Exec(
		"INSERT INTO bookings (id, desk_id, user_id, booking_date, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		bookingID,
		deskID,
		userID,
		bookingDate,
		now,
		now,
	)
	require.NoError(t, err)
}
