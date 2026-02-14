package areas

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
	"github.com/thorstenkramm/sithub/internal/spaces"
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

	// Run migrations
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
			note TEXT NOT NULL DEFAULT '',
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

func seedTestBooking(t *testing.T, store *sql.DB, id, itemID, userID, date string) {
	t.Helper()
	now := time.Now().Format(time.RFC3339)
	_, err := store.ExecContext(context.Background(),
		`INSERT INTO bookings (id, item_id, user_id, booking_date, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		id, itemID, userID, date, now, now)
	require.NoError(t, err)
}

func seedTestUser(t *testing.T, store *sql.DB, id, displayName string) {
	t.Helper()
	now := time.Now().Format(time.RFC3339)
	_, err := store.ExecContext(context.Background(),
		`INSERT INTO users (id, email, display_name, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?)`,
		id, id+"@test.local", displayName, now, now)
	require.NoError(t, err)
}

func testConfig() *spaces.Config {
	return &spaces.Config{
		Areas: []spaces.Area{
			{
				ID:   "area-1",
				Name: "Area One",
				ItemGroups: []spaces.ItemGroup{
					{
						ID:   "room-1",
						Name: "Room One",
						Items: []spaces.Item{
							{ID: "desk-1", Name: "Desk 1"},
							{ID: "desk-2", Name: "Desk 2"},
						},
					},
					{
						ID:   "room-2",
						Name: "Room Two",
						Items: []spaces.Item{
							{ID: "desk-3", Name: "Desk 3"},
						},
					},
				},
			},
			{
				ID:   "area-2",
				Name: "Area Two",
				ItemGroups: []spaces.ItemGroup{
					{
						ID:   "room-3",
						Name: "Room Three",
						Items: []spaces.Item{
							{ID: "desk-4", Name: "Desk 4"},
						},
					},
				},
			},
		},
	}
}

func TestPresenceHandlerAreaNotFound(t *testing.T) {
	t.Parallel()

	store := setupTestDB(t)
	cfg := testConfig()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/areas/unknown/presence", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("area_id")
	c.SetParamValues("unknown")

	h := PresenceHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestPresenceHandlerInvalidDate(t *testing.T) {
	t.Parallel()

	store := setupTestDB(t)
	cfg := testConfig()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/areas/area-1/presence?date=invalid", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("area_id")
	c.SetParamValues("area-1")

	h := PresenceHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPresenceHandlerReturnsUsersInArea(t *testing.T) {
	t.Parallel()

	store := setupTestDB(t)
	cfg := testConfig()

	// Seed users
	seedTestUser(t, store, "user-1", "Alice Smith")
	seedTestUser(t, store, "user-2", "Bob Jones")
	seedTestUser(t, store, "user-3", "Carol White")

	// Seed bookings for area-1
	seedTestBooking(t, store, "b1", "desk-1", "user-1", "2025-01-20")
	seedTestBooking(t, store, "b2", "desk-3", "user-2", "2025-01-20")
	// Booking in area-2 should not appear
	seedTestBooking(t, store, "b3", "desk-4", "user-3", "2025-01-20")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/areas/area-1/presence?date=2025-01-20", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("area_id")
	c.SetParamValues("area-1")

	h := PresenceHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get(echo.HeaderContentType), api.JSONAPIContentType)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Data, 2)

	// Sorted by item_id, so desk-1 comes before desk-3
	assert.Equal(t, "presence", resp.Data[0].Type)
	attrs0, ok := resp.Data[0].Attributes.(map[string]any)
	require.True(t, ok, "expected map attributes")
	assert.Equal(t, "Alice Smith", attrs0["user_name"])
	assert.Equal(t, "desk-1", attrs0["item_id"])
	assert.Equal(t, "Desk 1", attrs0["item_name"])
	assert.Equal(t, "room-1", attrs0["item_group_id"])
	assert.Equal(t, "Room One", attrs0["item_group_name"])

	attrs1, ok := resp.Data[1].Attributes.(map[string]any)
	require.True(t, ok, "expected map attributes")
	assert.Equal(t, "Bob Jones", attrs1["user_name"])
	assert.Equal(t, "desk-3", attrs1["item_id"])
	assert.Equal(t, "Room Two", attrs1["item_group_name"])
}

func TestPresenceHandlerExcludesOtherDates(t *testing.T) {
	t.Parallel()

	store := setupTestDB(t)
	cfg := testConfig()

	// Seed user
	seedTestUser(t, store, "user-1", "Alice")

	// Booking for different date
	seedTestBooking(t, store, "b1", "desk-1", "user-1", "2025-01-21")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/areas/area-1/presence?date=2025-01-20", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("area_id")
	c.SetParamValues("area-1")

	h := PresenceHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Empty(t, resp.Data)
}

func TestPresenceHandlerEmptyResult(t *testing.T) {
	t.Parallel()

	store := setupTestDB(t)
	cfg := testConfig()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/areas/area-1/presence?date=2025-01-20", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("area_id")
	c.SetParamValues("area-1")

	h := PresenceHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Empty(t, resp.Data)
}
