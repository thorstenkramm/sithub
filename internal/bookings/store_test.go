package bookings

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestFindDeskBookingsReturnsBookingInfo(t *testing.T) {
	store := setupTestStore(t)
	seedTestDeskData(t, store, []string{"desk-1", "desk-2"})
	seedTestBooking(t, store, "booking-1", "desk-1", "user-1", "2026-01-20")

	result, err := FindDeskBookings(t.Context(), store, "2026-01-20")
	require.NoError(t, err)
	require.Contains(t, result, "desk-1")
	assert.Equal(t, "booking-1", result["desk-1"].BookingID)
	assert.Equal(t, "user-1", result["desk-1"].UserID)
	// desk-2 should not be in results (not booked)
	require.NotContains(t, result, "desk-2")
}

func TestFindDeskBookingsReturnsEmptyMapForNoBookings(t *testing.T) {
	store := setupTestStore(t)
	seedTestDeskData(t, store, []string{"desk-1"})

	result, err := FindDeskBookings(t.Context(), store, "2026-01-20")
	require.NoError(t, err)
	require.Empty(t, result)
}

func TestFindDeskBookingsReturnsErrorOnClosedDB(t *testing.T) {
	store, err := db.Open(t.TempDir())
	require.NoError(t, err)
	require.NoError(t, store.Close())

	_, err = FindDeskBookings(t.Context(), store, "2026-01-20")
	require.Error(t, err)
}
