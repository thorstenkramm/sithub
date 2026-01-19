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
			desk_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			user_name TEXT NOT NULL DEFAULT '',
			booking_date TEXT NOT NULL,
			created_at TEXT NOT NULL,
			UNIQUE(desk_id, booking_date)
		)
	`)
	require.NoError(t, err)

	return store
}

func seedTestBooking(t *testing.T, store *sql.DB, id, deskID, userID, userName, date string) {
	t.Helper()
	_, err := store.ExecContext(context.Background(),
		`INSERT INTO bookings (id, desk_id, user_id, user_name, booking_date, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		id, deskID, userID, userName, date, time.Now().Format(time.RFC3339))
	require.NoError(t, err)
}

func testConfig() *spaces.Config {
	return &spaces.Config{
		Areas: []spaces.Area{
			{
				ID:   "area-1",
				Name: "Area One",
				Rooms: []spaces.Room{
					{
						ID:   "room-1",
						Name: "Room One",
						Desks: []spaces.Desk{
							{ID: "desk-1", Name: "Desk 1"},
							{ID: "desk-2", Name: "Desk 2"},
						},
					},
					{
						ID:   "room-2",
						Name: "Room Two",
						Desks: []spaces.Desk{
							{ID: "desk-3", Name: "Desk 3"},
						},
					},
				},
			},
			{
				ID:   "area-2",
				Name: "Area Two",
				Rooms: []spaces.Room{
					{
						ID:   "room-3",
						Name: "Room Three",
						Desks: []spaces.Desk{
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

	// Seed bookings for area-1
	seedTestBooking(t, store, "b1", "desk-1", "user-1", "Alice Smith", "2025-01-20")
	seedTestBooking(t, store, "b2", "desk-3", "user-2", "Bob Jones", "2025-01-20")
	// Booking in area-2 should not appear
	seedTestBooking(t, store, "b3", "desk-4", "user-3", "Carol White", "2025-01-20")

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

	// Sorted by user_name, so Alice comes before Bob
	assert.Equal(t, "presence", resp.Data[0].Type)
	attrs0, ok := resp.Data[0].Attributes.(map[string]any)
	require.True(t, ok, "expected map attributes")
	assert.Equal(t, "Alice Smith", attrs0["user_name"])
	assert.Equal(t, "desk-1", attrs0["desk_id"])
	assert.Equal(t, "Desk 1", attrs0["desk_name"])
	assert.Equal(t, "room-1", attrs0["room_id"])
	assert.Equal(t, "Room One", attrs0["room_name"])

	attrs1, ok := resp.Data[1].Attributes.(map[string]any)
	require.True(t, ok, "expected map attributes")
	assert.Equal(t, "Bob Jones", attrs1["user_name"])
	assert.Equal(t, "desk-3", attrs1["desk_id"])
	assert.Equal(t, "Room Two", attrs1["room_name"])
}

func TestPresenceHandlerExcludesOtherDates(t *testing.T) {
	t.Parallel()

	store := setupTestDB(t)
	cfg := testConfig()

	// Booking for different date
	seedTestBooking(t, store, "b1", "desk-1", "user-1", "Alice", "2025-01-21")

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
