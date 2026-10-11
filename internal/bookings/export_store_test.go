package bookings

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func seedExportUser(t *testing.T, store *sql.DB, id, name, email string, deleted bool) {
	t.Helper()
	flag := 0
	if deleted {
		flag = 1
	}
	_, err := store.Exec(`INSERT INTO users (id, email, display_name, user_source, deleted, created_at, updated_at)
		VALUES (?, ?, ?, 'internal', ?, '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z')`, id, email, name, flag)
	require.NoError(t, err)
}

func exportIDs(rows []ExportRow) []string {
	ids := make([]string, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.ID)
	}
	return ids
}

func TestListExportBookings_FiltersOrderAndTotal(t *testing.T) {
	t.Parallel()
	store := setupTestStore(t)
	seedExportUser(t, store, "u-bob", "bob", "bob@example.com", false)
	seedExportUser(t, store, "u-alice", "Alice", "alice@example.com", false)
	seedExportUser(t, store, "u-zed", "Zed", "zed@example.com", false)

	seedTestBooking(t, store, "b-before", "desk-1", "u-alice", "2026-10-03")
	seedTestBooking(t, store, "b-from", "desk-1", "u-zed", "2026-10-04")
	seedTestBooking(t, store, "b-mid-bob", "desk-2", "u-bob", "2026-10-06")
	seedTestBooking(t, store, "b-mid-alice", "desk-3", "u-alice", "2026-10-06")
	seedTestBooking(t, store, "b-to", "desk-1", "u-alice", "2026-10-08")
	seedTestBooking(t, store, "b-after", "desk-1", "u-alice", "2026-10-09")

	rows, total, err := ListExportBookings(context.Background(), store,
		ExportFilter{From: "2026-10-04", To: "2026-10-08"}, 20, 0)
	require.NoError(t, err)

	assert.Equal(t, 4, total)
	// Date ascending, then name case-insensitively ("Alice" before "bob").
	assert.Equal(t, []string{"b-from", "b-mid-alice", "b-mid-bob", "b-to"}, exportIDs(rows))
}

func TestListExportBookings_AreaItemFilter(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		itemIDs   []string
		wantIDs   []string
		wantTotal int
	}{
		{name: "nil means all areas", itemIDs: nil, wantIDs: []string{"b1", "b2", "b3"}, wantTotal: 3},
		{name: "area items only", itemIDs: []string{"desk-1", "desk-2"}, wantIDs: []string{"b1", "b2"}, wantTotal: 2},
		{name: "empty area returns nothing", itemIDs: []string{}, wantIDs: []string{}, wantTotal: 0},
		{name: "unrelated items", itemIDs: []string{"other"}, wantIDs: []string{}, wantTotal: 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			store := setupTestStore(t)
			seedExportUser(t, store, "u1", "Ann", "ann@example.com", false)
			seedTestBooking(t, store, "b1", "desk-1", "u1", "2026-10-05")
			seedTestBooking(t, store, "b2", "desk-2", "u1", "2026-10-06")
			seedTestBooking(t, store, "b3", "desk-9", "u1", "2026-10-07")

			rows, total, err := ListExportBookings(context.Background(), store,
				ExportFilter{From: "2026-10-01", To: "2026-10-31", ItemIDs: tc.itemIDs}, 20, 0)
			require.NoError(t, err)
			assert.Equal(t, tc.wantTotal, total)
			assert.Equal(t, tc.wantIDs, exportIDs(rows))
		})
	}
}

func TestListExportBookings_PersonResolution(t *testing.T) {
	t.Parallel()
	store := setupTestStore(t)
	seedExportUser(t, store, "u-booker", "Booker", "booker@example.com", false)
	seedExportUser(t, store, "u-colleague", "Colleague", "colleague@example.com", false)
	seedExportUser(t, store, "u-gone", "Gone Person", "gone@example.com", true)

	seedTestBooking(t, store, "own", "d1", "u-booker", "2026-10-05")
	seedTestBookingFull(t, store, "on-behalf", "d2", "u-colleague", "u-booker", "2026-10-06")
	seedTestBookingWithGuest(t, store, "guest", "d3", "guest-1a2b3c4d", "u-booker", "2026-10-07",
		true, "Guest Person", "guest@example.com")
	seedTestBookingWithGuest(t, store, "guest-no-mail", "d4", "guest-aabbccdd", "u-booker", "2026-10-08",
		true, "Mailless Guest", "")
	seedTestBooking(t, store, "soft-deleted", "d5", "u-gone", "2026-10-09")
	seedTestBooking(t, store, "orphan", "d6", "u-hard-deleted", "2026-10-10")

	rows, total, err := ListExportBookings(context.Background(), store,
		ExportFilter{From: "2026-10-01", To: "2026-10-31"}, 20, 0)
	require.NoError(t, err)
	require.Equal(t, 6, total)

	got := map[string][3]interface{}{}
	ids := map[string]string{}
	for _, r := range rows {
		got[r.ID] = [3]interface{}{r.PersonName, r.PersonEmail, r.IsGuest}
		ids[r.ID] = r.PersonUserID
	}
	assert.Equal(t, [3]interface{}{"Booker", "booker@example.com", false}, got["own"])
	assert.Equal(t, [3]interface{}{"Colleague", "colleague@example.com", false}, got["on-behalf"],
		"on-behalf bookings belong to the colleague, not the booker")
	assert.Equal(t, [3]interface{}{"Guest Person", "guest@example.com", true}, got["guest"])
	assert.Equal(t, [3]interface{}{"Mailless Guest", "", true}, got["guest-no-mail"])
	assert.Equal(t, [3]interface{}{"Gone Person", "gone@example.com", false}, got["soft-deleted"],
		"soft-deleted users stay in the export")
	assert.Equal(t, [3]interface{}{"", "", false}, got["orphan"], "a missing user row yields empty name and email")

	assert.Equal(t, map[string]string{
		"own": "u-booker", "on-behalf": "u-colleague", "guest": "", "guest-no-mail": "",
		"soft-deleted": "u-gone", "orphan": "",
	}, ids, "the user id is the booked person's, empty for guests and removed users")
}

func TestListExportBookings_Paging(t *testing.T) {
	t.Parallel()
	store := setupTestStore(t)
	seedExportUser(t, store, "u1", "Same Name", "same@example.com", false)
	for _, id := range []string{"b5", "b3", "b1", "b4", "b2"} {
		seedTestBooking(t, store, id, "desk-"+id, "u1", "2026-10-05")
	}

	cases := []struct {
		name          string
		limit, offset int
		want          []string
	}{
		{name: "first page", limit: 2, offset: 0, want: []string{"b1", "b2"}},
		{name: "middle page", limit: 2, offset: 2, want: []string{"b3", "b4"}},
		{name: "last partial page", limit: 2, offset: 4, want: []string{"b5"}},
		{name: "past the end", limit: 2, offset: 10, want: []string{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rows, total, err := ListExportBookings(context.Background(), store,
				ExportFilter{From: "2026-10-01", To: "2026-10-31"}, tc.limit, tc.offset)
			require.NoError(t, err)
			assert.Equal(t, 5, total, "total ignores paging")
			assert.Equal(t, tc.want, exportIDs(rows), "equal date and name are ordered by booking id")
		})
	}
}
