package itemgroups

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/db"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	store, err := db.Open(t.TempDir())
	require.NoError(t, err)
	t.Cleanup(func() {
		if err := store.Close(); err != nil {
			t.Errorf("close store: %v", err)
		}
	})

	_, err = store.Exec(`
		CREATE TABLE IF NOT EXISTS bookings (
			id TEXT PRIMARY KEY,
			item_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			booked_by_user_id TEXT NOT NULL DEFAULT '',
			booking_date TEXT NOT NULL,
			is_guest INTEGER NOT NULL DEFAULT 0,
			guest_name TEXT NOT NULL DEFAULT '',
			guest_email TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			UNIQUE(item_id, booking_date)
		)
	`)
	require.NoError(t, err)

	_, err = store.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL,
			display_name TEXT NOT NULL,
			password_hash TEXT NOT NULL DEFAULT '',
			user_source TEXT NOT NULL DEFAULT 'internal',
			entra_id TEXT NOT NULL DEFAULT '',
			is_admin INTEGER NOT NULL DEFAULT 0,
			last_login TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)
	`)
	require.NoError(t, err)

	return store
}

func seedBooking(t *testing.T, store *sql.DB, id, itemID, userID, date string) {
	t.Helper()
	now := time.Now().Format(time.RFC3339)
	_, err := store.ExecContext(context.Background(),
		`INSERT INTO bookings (id, item_id, user_id, booking_date, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		id, itemID, userID, date, now, now)
	require.NoError(t, err)
}

func seedUser(t *testing.T, store *sql.DB, id, displayName string) {
	t.Helper()
	now := time.Now().Format(time.RFC3339)
	_, err := store.ExecContext(context.Background(),
		`INSERT INTO users (id, email, display_name, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?)`,
		id, id+"@test.local", displayName, now, now)
	require.NoError(t, err)
}

func TestBookingsHandlerItemGroupNotFound(t *testing.T) {
	t.Parallel()

	store := setupTestDB(t)
	cfg := testConfig()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/item-groups/unknown/bookings?date=2025-01-20", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("item_group_id")
	c.SetParamValues("unknown")

	h := BookingsHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestBookingsHandlerInvalidDate(t *testing.T) {
	t.Parallel()

	store := setupTestDB(t)
	cfg := testConfig()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/item-groups/ig-1/bookings?date=bad", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("item_group_id")
	c.SetParamValues("ig-1")

	h := BookingsHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestBookingsHandlerReturnsBookings(t *testing.T) {
	t.Parallel()

	store := setupTestDB(t)
	cfg := testConfig()

	seedUser(t, store, "user-1", "Alice Smith")
	seedUser(t, store, "user-2", "Bob Jones")
	seedBooking(t, store, "b1", "item-1", "user-1", "2025-01-20")
	seedBooking(t, store, "b2", "item-2", "user-2", "2025-01-20")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/item-groups/ig-1/bookings?date=2025-01-20", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("item_group_id")
	c.SetParamValues("ig-1")

	h := BookingsHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get(echo.HeaderContentType), api.JSONAPIContentType)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Data, 2)

	assert.Equal(t, "bookings", resp.Data[0].Type)
	attrs0, ok := resp.Data[0].Attributes.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "item-1", attrs0["item_id"])
	assert.Equal(t, "Desk 1", attrs0["item_name"])
	assert.Equal(t, "Alice Smith", attrs0["user_name"])

	attrs1, ok := resp.Data[1].Attributes.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "item-2", attrs1["item_id"])
	assert.Equal(t, "Bob Jones", attrs1["user_name"])
}

func TestBookingsHandlerEmptyResult(t *testing.T) {
	t.Parallel()

	store := setupTestDB(t)
	cfg := testConfig()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/item-groups/ig-1/bookings?date=2025-01-20", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("item_group_id")
	c.SetParamValues("ig-1")

	h := BookingsHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Empty(t, resp.Data)
}

func TestBookingsHandlerExcludesOtherItemGroups(t *testing.T) {
	t.Parallel()

	store := setupTestDB(t)
	cfg := testConfig()

	seedUser(t, store, "user-1", "Alice")
	// Booking for item-3 which is in ig-2, not ig-1
	seedBooking(t, store, "b1", "item-3", "user-1", "2025-01-20")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/item-groups/ig-1/bookings?date=2025-01-20", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("item_group_id")
	c.SetParamValues("ig-1")

	h := BookingsHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Empty(t, resp.Data)
}
