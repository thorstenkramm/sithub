package desks

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

func TestListHandlerNotFound(t *testing.T) {
	t.Parallel()

	cfg := &spaces.Config{}
	store := setupStore(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/rooms/missing/desks", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("room_id")
	c.SetParamValues("missing")

	h := ListHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var resp api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Errors, 1)
	assert.Equal(t, "not_found", resp.Errors[0].Code)
}

func TestListHandlerRejectsInvalidDate(t *testing.T) {
	t.Parallel()

	cfg := &spaces.Config{
		Areas: []spaces.Area{
			{
				ID:   "area-1",
				Name: "Office",
				Rooms: []spaces.Room{
					{
						ID:   "room-1",
						Name: "Room 1",
					},
				},
			},
		},
	}
	store := setupStore(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/rooms/room-1/desks?date=not-a-date", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("room_id")
	c.SetParamValues("room-1")

	h := ListHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Errors, 1)
	assert.Equal(t, "bad_request", resp.Errors[0].Code)
}

func TestListHandlerReturnsDesks(t *testing.T) {
	t.Parallel()

	cfg := &spaces.Config{
		Areas: []spaces.Area{
			{
				ID:   "area-1",
				Name: "Office",
				Rooms: []spaces.Room{
					{
						ID:   "room-1",
						Name: "Room 1",
						Desks: []spaces.Desk{
							{
								ID:        "desk-1",
								Name:      "Desk 1",
								Equipment: []string{"Monitor", "Keyboard"},
								Warning:   "USB-C only",
							},
						},
					},
				},
			},
		},
	}
	store := setupStore(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/rooms/room-1/desks", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("room_id")
	c.SetParamValues("room-1")

	h := ListHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Data, 1)
	assert.Equal(t, "desks", resp.Data[0].Type)
	assert.Equal(t, "desk-1", resp.Data[0].ID)

	attrs, ok := resp.Data[0].Attributes.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "Desk 1", attrs["name"])
	assert.Equal(t, "USB-C only", attrs["warning"])
	assert.Equal(t, "available", attrs["availability"])
}

func TestListHandlerShowsBookedDesk(t *testing.T) {
	t.Parallel()

	cfg := &spaces.Config{
		Areas: []spaces.Area{
			{
				ID:   "area-1",
				Name: "Office",
				Rooms: []spaces.Room{
					{
						ID:   "room-1",
						Name: "Room 1",
						Desks: []spaces.Desk{
							{
								ID:        "desk-1",
								Name:      "Desk 1",
								Equipment: []string{"Monitor"},
							},
							{
								ID:        "desk-2",
								Name:      "Desk 2",
								Equipment: []string{"Keyboard"},
							},
						},
					},
				},
			},
		},
	}
	store := setupStore(t)
	// No need to seed desk data - desk_id is just a string reference now
	seedBooking(t, store, "booking-1", "desk-1", "user-1", "2026-01-20")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/rooms/room-1/desks?date=2026-01-20", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("room_id")
	c.SetParamValues("room-1")

	h := ListHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Data, 2)

	attrsDesk1, ok := resp.Data[0].Attributes.(map[string]interface{})
	require.True(t, ok)
	attrsDesk2, ok := resp.Data[1].Attributes.(map[string]interface{})
	require.True(t, ok)

	availabilityDesk1, ok := attrsDesk1["availability"].(string)
	require.True(t, ok)
	availabilityDesk2, ok := attrsDesk2["availability"].(string)
	require.True(t, ok)

	availability := map[string]string{
		resp.Data[0].ID: availabilityDesk1,
		resp.Data[1].ID: availabilityDesk2,
	}
	assert.Equal(t, "occupied", availability["desk-1"])
	assert.Equal(t, "available", availability["desk-2"])
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

func seedBooking(t *testing.T, store *sql.DB, bookingID, deskID, userID, bookingDate string) {
	t.Helper()

	now := time.Now().UTC().Format(time.RFC3339)
	_, err := store.Exec(`
		INSERT INTO bookings 
		(id, desk_id, user_id, user_name, booked_by_user_id, booked_by_user_name, 
		 booking_date, is_guest, guest_email, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		bookingID, deskID, userID, "Test User", userID, "Test User",
		bookingDate, 0, "", now, now,
	)
	require.NoError(t, err)
}
