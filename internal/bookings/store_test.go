package bookings

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/db"
)

func TestFindBookedDeskIDs(t *testing.T) {
	store := setupTestStore(t)
	seedTestDeskData(t, store, []string{"desk-1"})
	seedTestBooking(t, store, "booking-1", "desk-1", "user-1", "2026-01-20")

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
