package rooms

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/db"
	"github.com/thorstenkramm/sithub/internal/spaces"
)

func TestBookingsHandlerRoomNotFound(t *testing.T) {
	t.Parallel()

	cfg := &spaces.Config{}
	store := setupTestStore(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/rooms/nonexistent/bookings", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("room_id")
	c.SetParamValues("nonexistent")

	h := BookingsHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestBookingsHandlerInvalidDate(t *testing.T) {
	t.Parallel()

	cfg := testSpacesConfig()
	store := setupTestStore(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/rooms/room-1/bookings?date=invalid", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("room_id")
	c.SetParamValues("room-1")

	h := BookingsHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestBookingsHandlerReturnsBookingsForRoom(t *testing.T) {
	t.Parallel()

	cfg := testSpacesConfig()
	store := setupTestStore(t)
	// No need to seed space data - it comes from config now

	tomorrow := time.Now().UTC().AddDate(0, 0, 1).Format(time.DateOnly)
	seedTestBooking(t, store, "booking-1", "desk-1", "user-1", "Alice Smith", tomorrow)
	seedTestBooking(t, store, "booking-2", "desk-2", "user-2", "Bob Jones", tomorrow)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/rooms/room-1/bookings?date="+tomorrow, http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("room_id")
	c.SetParamValues("room-1")

	h := BookingsHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, api.JSONAPIContentType, rec.Header().Get(echo.HeaderContentType))

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Data, 2)

	// First booking
	attrs0, ok := resp.Data[0].Attributes.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "desk-1", attrs0["desk_id"])
	assert.Equal(t, "Desk 1", attrs0["desk_name"])
	assert.Equal(t, "user-1", attrs0["user_id"])
	assert.Equal(t, "Alice Smith", attrs0["user_name"])

	// Second booking
	attrs1, ok := resp.Data[1].Attributes.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "desk-2", attrs1["desk_id"])
	assert.Equal(t, "Bob Jones", attrs1["user_name"])
}

func TestBookingsHandlerExcludesOtherDates(t *testing.T) {
	t.Parallel()

	cfg := testSpacesConfig()
	store := setupTestStore(t)
	// No need to seed space data - it comes from config now

	tomorrow := time.Now().UTC().AddDate(0, 0, 1).Format(time.DateOnly)
	dayAfter := time.Now().UTC().AddDate(0, 0, 2).Format(time.DateOnly)
	seedTestBooking(t, store, "booking-1", "desk-1", "user-1", "Alice", tomorrow)
	seedTestBooking(t, store, "booking-2", "desk-2", "user-2", "Bob", dayAfter)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/rooms/room-1/bookings?date="+tomorrow, http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("room_id")
	c.SetParamValues("room-1")

	h := BookingsHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Data, 1)
	assert.Equal(t, "booking-1", resp.Data[0].ID)
}

func TestBookingsHandlerEmptyResult(t *testing.T) {
	t.Parallel()

	cfg := testSpacesConfig()
	store := setupTestStore(t)
	// No need to seed space data - it comes from config now

	tomorrow := time.Now().UTC().AddDate(0, 0, 1).Format(time.DateOnly)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/rooms/room-1/bookings?date="+tomorrow, http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("room_id")
	c.SetParamValues("room-1")

	h := BookingsHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Len(t, resp.Data, 0)
}

func testSpacesConfig() *spaces.Config {
	return &spaces.Config{
		Areas: []spaces.Area{
			{
				ID:   "area-1",
				Name: "Office",
				Rooms: []spaces.Room{
					{
						ID:   "room-1",
						Name: "Room 1",
						Desks: []spaces.Desk{
							{ID: "desk-1", Name: "Desk 1", Equipment: []string{"Monitor"}},
							{ID: "desk-2", Name: "Desk 2", Equipment: []string{"Monitor"}},
						},
					},
				},
			},
		},
	}
}

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

func seedTestBooking(t *testing.T, store *sql.DB, bookingID, deskID, userID, userName, bookingDate string) {
	t.Helper()

	now := time.Now().UTC().Format(time.RFC3339)
	_, err := store.Exec(`
		INSERT INTO bookings 
		(id, desk_id, user_id, user_name, booked_by_user_id, booked_by_user_name, 
		 booking_date, is_guest, guest_email, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		bookingID, deskID, userID, userName, userID, userName,
		bookingDate, 0, "", now, now,
	)
	require.NoError(t, err)
}
