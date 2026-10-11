package users

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testToday = "2026-10-08"

type seedBooking struct {
	id, userID, bookedBy, date string
	isGuest                    bool
}

func insertBookingRow(t *testing.T, db *sql.DB, b seedBooking) {
	t.Helper()
	guestName := ""
	if b.isGuest {
		guestName = "Guest " + b.id
	}
	_, err := db.Exec(`INSERT INTO bookings
		(id, item_id, user_id, booked_by_user_id, booking_date, is_guest, guest_name, guest_email,
		 created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, '', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z')`,
		b.id, "item-"+b.id, b.userID, b.bookedBy, b.date, boolToInt(b.isGuest), guestName)
	require.NoError(t, err)
}

func bookingIDs(t *testing.T, db *sql.DB) []string {
	t.Helper()
	rows, err := db.Query(`SELECT id FROM bookings ORDER BY id`)
	require.NoError(t, err)
	defer func() {
		_ = rows.Close() //nolint:errcheck // Test cleanup
	}()
	ids := []string{}
	for rows.Next() {
		var id string
		require.NoError(t, rows.Scan(&id))
		ids = append(ids, id)
	}
	require.NoError(t, rows.Err())
	return ids
}

func TestDeleteOrDeactivate(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		bookings      []seedBooking
		wantMode      DeletionMode
		wantCanceled  []string
		wantRemaining []string
	}{
		{
			name:     "no bookings hard deletes",
			wantMode: DeletionModeDeleted, wantCanceled: []string{}, wantRemaining: []string{},
		},
		{
			name:     "past own booking deactivates without canceling",
			bookings: []seedBooking{{id: "b1", userID: "x", bookedBy: "x", date: "2026-10-01"}},
			wantMode: DeletionModeDeactivated, wantCanceled: []string{}, wantRemaining: []string{"b1"},
		},
		{
			name:     "booking made for a colleague is kept",
			bookings: []seedBooking{{id: "b1", userID: "colleague", bookedBy: "x", date: "2026-10-20"}},
			wantMode: DeletionModeDeactivated, wantCanceled: []string{}, wantRemaining: []string{"b1"},
		},
		{
			name: "future own booking canceled and past kept",
			bookings: []seedBooking{
				{id: "b1", userID: "x", bookedBy: "x", date: "2026-10-01"},
				{id: "b2", userID: "x", bookedBy: "x", date: "2026-10-20"},
			},
			wantMode: DeletionModeDeactivated, wantCanceled: []string{"b2"}, wantRemaining: []string{"b1"},
		},
		{
			name:     "today's booking is canceled",
			bookings: []seedBooking{{id: "b1", userID: "x", bookedBy: "x", date: testToday}},
			wantMode: DeletionModeDeactivated, wantCanceled: []string{"b1"}, wantRemaining: []string{},
		},
		{
			name: "future guest booking hosted by the user is canceled",
			bookings: []seedBooking{
				{id: "b1", userID: "guest-1", bookedBy: "x", date: "2026-10-20", isGuest: true},
			},
			wantMode: DeletionModeDeactivated, wantCanceled: []string{"b1"}, wantRemaining: []string{},
		},
		{
			name: "on-behalf booking made by a colleague for the user is canceled",
			bookings: []seedBooking{
				{id: "b1", userID: "x", bookedBy: "colleague", date: "2026-10-20"},
				{id: "b2", userID: "other", bookedBy: "other", date: "2026-10-20"},
			},
			wantMode: DeletionModeDeactivated, wantCanceled: []string{"b1"}, wantRemaining: []string{"b2"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := setupTestDB(t)
			ctx := context.Background()
			insertUserRow(t, db, "x", "Xavier")
			for _, b := range tc.bookings {
				insertBookingRow(t, db, b)
			}

			impact, err := DeletionImpact(ctx, db, "x", testToday)
			require.NoError(t, err)
			assert.Equal(t, len(tc.bookings) > 0, impact.HasBookings)

			result, err := DeleteOrDeactivate(ctx, db, "x", testToday)
			require.NoError(t, err)
			assert.Equal(t, tc.wantMode, result.Mode)

			canceled := make([]string, 0, len(result.CanceledBookings))
			for _, b := range result.CanceledBookings {
				canceled = append(canceled, b.ID)
			}
			assert.Equal(t, tc.wantCanceled, canceled)
			assert.Equal(t, impact.FutureBookingCount, len(canceled), "preview count matches cancellation")
			assert.Equal(t, tc.wantRemaining, bookingIDs(t, db))

			rec, err := FindByID(ctx, db, "x")
			if tc.wantMode == DeletionModeDeleted {
				assert.ErrorIs(t, err, ErrUserNotFound)
				return
			}
			require.NoError(t, err)
			assert.True(t, rec.Deleted)
		})
	}
}

func TestDeleteOrDeactivate_CanceledBookingDetails(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	insertUserRow(t, db, "x", "Xavier")
	insertBookingRow(t, db, seedBooking{id: "b1", userID: "guest-1", bookedBy: "x", date: "2026-10-20", isGuest: true})

	result, err := DeleteOrDeactivate(context.Background(), db, "x", testToday)
	require.NoError(t, err)
	require.Len(t, result.CanceledBookings, 1)
	assert.Equal(t, CanceledBooking{
		ID: "b1", ItemID: "item-b1", UserID: "guest-1", BookingDate: "2026-10-20",
		IsGuest: true, GuestName: "Guest b1",
	}, result.CanceledBookings[0])
}

func TestDeleteOrDeactivate_AlreadyDeactivatedIsIdempotent(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()
	insertUserRow(t, db, "x", "Xavier")
	insertBookingRow(t, db, seedBooking{id: "b1", userID: "x", bookedBy: "x", date: "2026-10-01"})

	_, err := DeleteOrDeactivate(ctx, db, "x", testToday)
	require.NoError(t, err)
	result, err := DeleteOrDeactivate(ctx, db, "x", testToday)
	require.NoError(t, err)
	assert.Equal(t, DeletionModeDeactivated, result.Mode)
	assert.Empty(t, result.CanceledBookings)
}

func TestDeleteOrDeactivate_UnknownUser(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)

	_, err := DeleteOrDeactivate(context.Background(), db, "missing", testToday)
	assert.ErrorIs(t, err, ErrUserNotFound)

	_, err = DeletionImpact(context.Background(), db, "missing", testToday)
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestListColleagues_ExcludesDeleted(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	insertUserRow(t, db, "a", "Alice")
	insertUserRow(t, db, "b", "Bob")
	_, err := db.Exec(`UPDATE users SET deleted = 1 WHERE id = 'b'`)
	require.NoError(t, err)

	colleagues, err := ListColleagues(context.Background(), db)
	require.NoError(t, err)
	require.Len(t, colleagues, 1)
	assert.Equal(t, "a", colleagues[0].ID)
}

func TestUpsertEntraIDUser_ReactivatesDeletedUser(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()

	rec, err := UpsertEntraIDUser(ctx, db, "entra-1", "e@example.com", "Entra", false)
	require.NoError(t, err)
	_, err = db.Exec(`UPDATE users SET deleted = 1 WHERE id = ?`, rec.ID)
	require.NoError(t, err)

	rec, err = UpsertEntraIDUser(ctx, db, "entra-1", "e@example.com", "Entra", false)
	require.NoError(t, err)
	assert.False(t, rec.Deleted)

	colleagues, err := ListColleagues(ctx, db)
	require.NoError(t, err)
	assert.Len(t, colleagues, 1)
}

func TestLookupAccess(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()
	insertUserRow(t, db, "active", "Active")
	insertUserRow(t, db, "admin", "Admin")
	insertUserRow(t, db, "gone", "Gone")
	_, err := db.Exec(`UPDATE users SET is_admin = 1 WHERE id = 'admin'`)
	require.NoError(t, err)
	_, err = db.Exec(`UPDATE users SET deleted = 1 WHERE id = 'gone'`)
	require.NoError(t, err)

	cases := []struct {
		id   string
		want Access
	}{
		{id: "active", want: Access{Active: true}},
		{id: "admin", want: Access{Active: true, IsAdmin: true}},
		{id: "gone", want: Access{Active: false}},
		{id: "missing", want: Access{Active: false}},
	}
	for _, tc := range cases {
		t.Run(tc.id, func(t *testing.T) {
			t.Parallel()
			got, err := LookupAccess(ctx, db, tc.id)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestUpdateUser_IsDataExporter(t *testing.T) {
	t.Parallel()

	yes, no := true, false
	cases := []struct {
		name    string
		initial bool
		field   *bool
		want    bool
	}{
		{name: "grant", initial: false, field: &yes, want: true},
		{name: "revoke", initial: true, field: &no, want: false},
		{name: "nil leaves granted flag", initial: true, field: nil, want: true},
		{name: "nil leaves revoked flag", initial: false, field: nil, want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := setupTestDB(t)
			insertUserRow(t, db, "u1", "User One")
			_, err := db.Exec(`UPDATE users SET is_data_exporter = ? WHERE id = 'u1'`, boolToInt(tc.initial))
			require.NoError(t, err)

			rec, err := UpdateUser(context.Background(), db, "u1", UpdateFields{IsDataExporter: tc.field})
			require.NoError(t, err)
			assert.Equal(t, tc.want, rec.IsDataExporter)
		})
	}
}

func TestUpsertEntraIDUser_KeepsDataExporterFlag(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	ctx := context.Background()

	rec, err := UpsertEntraIDUser(ctx, db, "entra-1", "e@example.com", "Entra", false)
	require.NoError(t, err)
	_, err = db.Exec(`UPDATE users SET is_data_exporter = 1 WHERE id = ?`, rec.ID)
	require.NoError(t, err)

	rec, err = UpsertEntraIDUser(ctx, db, "entra-1", "e@example.com", "Entra Renamed", true)
	require.NoError(t, err)
	assert.True(t, rec.IsDataExporter, "an Entra ID login must never overwrite the role")
}

func TestLookupAccess_IncludesDataExporter(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	insertUserRow(t, db, "u1", "User One")
	_, err := db.Exec(`UPDATE users SET is_data_exporter = 1 WHERE id = 'u1'`)
	require.NoError(t, err)

	got, err := LookupAccess(context.Background(), db, "u1")
	require.NoError(t, err)
	assert.Equal(t, Access{Active: true, IsDataExporter: true}, got)
}
