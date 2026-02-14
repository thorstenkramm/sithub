package bookings

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/db"
)

func TestFindBookedItemIDs(t *testing.T) {
	store := setupTestStore(t)
	// No need to seed item data - item_id is just a string reference now
	seedTestBooking(t, store, "booking-1", "desk-1", "user-1", "2026-01-20")

	booked, err := FindBookedItemIDs(t.Context(), store, "2026-01-20")
	require.NoError(t, err)
	require.Contains(t, booked, "desk-1")
}

func TestFindBookedItemIDsReturnsErrorOnClosedDB(t *testing.T) {
	store, err := db.Open(t.TempDir())
	require.NoError(t, err)
	require.NoError(t, store.Close())

	_, err = FindBookedItemIDs(t.Context(), store, "2026-01-20")
	require.Error(t, err)
}

func TestFindItemBookingsReturnsBookingInfo(t *testing.T) {
	store := setupTestStore(t)
	seedTestBooking(t, store, "booking-1", "desk-1", "user-1", "2026-01-20")

	result, err := FindItemBookings(t.Context(), store, "2026-01-20")
	require.NoError(t, err)
	require.Contains(t, result, "desk-1")
	assert.Equal(t, "booking-1", result["desk-1"].BookingID)
	assert.Equal(t, "user-1", result["desk-1"].UserID)
	assert.False(t, result["desk-1"].IsGuest)
	assert.Empty(t, result["desk-1"].GuestName)
	require.NotContains(t, result, "desk-2")
}

func TestFindItemBookingsReturnsGuestInfo(t *testing.T) {
	store := setupTestStore(t)
	seedTestBookingWithGuest(t, store, "booking-1", "desk-1", "booker-1",
		"booker-1", "2026-01-20", true, "John Visitor", "john@example.com")

	result, err := FindItemBookings(t.Context(), store, "2026-01-20")
	require.NoError(t, err)
	require.Contains(t, result, "desk-1")
	assert.Equal(t, "booking-1", result["desk-1"].BookingID)
	assert.True(t, result["desk-1"].IsGuest)
	assert.Equal(t, "John Visitor", result["desk-1"].GuestName)
}

func TestFindItemBookingsReturnsEmptyMapForNoBookings(t *testing.T) {
	store := setupTestStore(t)
	// No need to seed item data - item_id is just a string reference now

	result, err := FindItemBookings(t.Context(), store, "2026-01-20")
	require.NoError(t, err)
	require.Empty(t, result)
}

func TestFindItemBookingsReturnsErrorOnClosedDB(t *testing.T) {
	store, err := db.Open(t.TempDir())
	require.NoError(t, err)
	require.NoError(t, store.Close())

	_, err = FindItemBookings(t.Context(), store, "2026-01-20")
	require.Error(t, err)
}
