//nolint:dupl // Test files often have similar test structure for different entities
package admin

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/auth"
	"github.com/thorstenkramm/sithub/internal/db"
	"github.com/thorstenkramm/sithub/internal/spaces"
)

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

func setAdminUser(c echo.Context) {
	c.Set("user", &auth.User{
		ID:      "admin-1",
		Name:    "Admin",
		IsAdmin: true,
	})
}

func setNonAdminUser(c echo.Context) {
	c.Set("user", &auth.User{
		ID:      "user-1",
		Name:    "User",
		IsAdmin: false,
	})
}

//nolint:unparam // Param needed for test clarity
func seedTestArea(t *testing.T, store *sql.DB, areaID, name string) {
	t.Helper()
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := store.Exec(
		`INSERT INTO areas (id, name, description, floor_plan, created_at, updated_at) 
		 VALUES (?, ?, ?, ?, ?, ?)`,
		areaID, name, "", "", now, now,
	)
	require.NoError(t, err)
}

//nolint:unparam // Param needed for test clarity
func seedTestRoom(t *testing.T, store *sql.DB, roomID, areaID, name string) {
	t.Helper()
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := store.Exec(
		`INSERT INTO rooms (id, area_id, name, description, floor_plan, created_at, updated_at) 
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		roomID, areaID, name, "", "", now, now,
	)
	require.NoError(t, err)
}

func seedTestDesk(t *testing.T, store *sql.DB, deskID, roomID, name string) {
	t.Helper()
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := store.Exec(
		`INSERT INTO desks (id, room_id, name, equipment, warning, created_at, updated_at) 
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		deskID, roomID, name, "[]", "", now, now,
	)
	require.NoError(t, err)
}

// --- Area Tests ---

func TestListAreasHandler_ForbiddenForNonAdmin(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := spaces.NewStore(store)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/areas", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	setNonAdminUser(c)

	h := ListAreasHandler(spacesStore)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestListAreasHandler_Success(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := spaces.NewStore(store)
	seedTestArea(t, store, "area-1", "Office")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/areas", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	setAdminUser(c)

	h := ListAreasHandler(spacesStore)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Data, 1)
	assert.Equal(t, "areas", resp.Data[0].Type)
	assert.Equal(t, "area-1", resp.Data[0].ID)
}

func TestCreateAreaHandler_Success(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := spaces.NewStore(store)
	cfg := &spaces.Config{}
	configHolder := NewConfigHolder(cfg)

	body := `{"data":{"type":"areas","id":"new-area","attributes":{"name":"New Area"}}}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/areas", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	setAdminUser(c)

	h := CreateAreaHandler(spacesStore, configHolder)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp api.SingleResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "areas", resp.Data.Type)
	assert.Equal(t, "new-area", resp.Data.ID)
}

func TestCreateAreaHandler_InvalidType(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := spaces.NewStore(store)
	cfg := &spaces.Config{}
	configHolder := NewConfigHolder(cfg)

	body := `{"data":{"type":"invalid","attributes":{"name":"Test"}}}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/areas", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	setAdminUser(c)

	h := CreateAreaHandler(spacesStore, configHolder)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateAreaHandler_Success(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := spaces.NewStore(store)
	seedTestArea(t, store, "area-1", "Office")
	cfg := &spaces.Config{}
	configHolder := NewConfigHolder(cfg)

	body := `{"data":{"type":"areas","attributes":{"name":"Updated Office"}}}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/areas/area-1", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("area_id")
	c.SetParamValues("area-1")
	setAdminUser(c)

	h := UpdateAreaHandler(spacesStore, configHolder)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp api.SingleResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	attrs, ok := resp.Data.Attributes.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "Updated Office", attrs["name"])
}

func TestUpdateAreaHandler_NotFound(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := spaces.NewStore(store)
	cfg := &spaces.Config{}
	configHolder := NewConfigHolder(cfg)

	body := `{"data":{"type":"areas","attributes":{"name":"Test"}}}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/areas/missing", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("area_id")
	c.SetParamValues("missing")
	setAdminUser(c)

	h := UpdateAreaHandler(spacesStore, configHolder)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteAreaHandler_Success(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := spaces.NewStore(store)
	seedTestArea(t, store, "area-1", "Office")
	cfg := &spaces.Config{}
	configHolder := NewConfigHolder(cfg)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/areas/area-1", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("area_id")
	c.SetParamValues("area-1")
	setAdminUser(c)

	h := DeleteAreaHandler(spacesStore, configHolder)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestDeleteAreaHandler_NotFound(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := spaces.NewStore(store)
	cfg := &spaces.Config{}
	configHolder := NewConfigHolder(cfg)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/areas/missing", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("area_id")
	c.SetParamValues("missing")
	setAdminUser(c)

	h := DeleteAreaHandler(spacesStore, configHolder)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// --- Room Tests ---

func TestListRoomsHandler_Success(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := spaces.NewStore(store)
	seedTestArea(t, store, "area-1", "Office")
	seedTestRoom(t, store, "room-1", "area-1", "Room 1")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/areas/area-1/rooms", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("area_id")
	c.SetParamValues("area-1")
	setAdminUser(c)

	h := ListRoomsHandler(spacesStore)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Data, 1)
	assert.Equal(t, "rooms", resp.Data[0].Type)
}

func TestCreateRoomHandler_Success(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := spaces.NewStore(store)
	seedTestArea(t, store, "area-1", "Office")
	cfg := &spaces.Config{}
	configHolder := NewConfigHolder(cfg)

	body := `{"data":{"type":"rooms","id":"new-room","attributes":{"name":"New Room"}}}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/areas/area-1/rooms", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("area_id")
	c.SetParamValues("area-1")
	setAdminUser(c)

	h := CreateRoomHandler(spacesStore, configHolder)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestCreateRoomHandler_AreaNotFound(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := spaces.NewStore(store)
	cfg := &spaces.Config{}
	configHolder := NewConfigHolder(cfg)

	body := `{"data":{"type":"rooms","attributes":{"name":"New Room"}}}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/areas/missing/rooms", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("area_id")
	c.SetParamValues("missing")
	setAdminUser(c)

	h := CreateRoomHandler(spacesStore, configHolder)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateRoomHandler_Success(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := spaces.NewStore(store)
	seedTestArea(t, store, "area-1", "Office")
	seedTestRoom(t, store, "room-1", "area-1", "Room 1")
	cfg := &spaces.Config{}
	configHolder := NewConfigHolder(cfg)

	body := `{"data":{"type":"rooms","attributes":{"name":"Updated Room"}}}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/rooms/room-1", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("room_id")
	c.SetParamValues("room-1")
	setAdminUser(c)

	h := UpdateRoomHandler(spacesStore, configHolder)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUpdateRoomHandler_NotFound(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := spaces.NewStore(store)
	cfg := &spaces.Config{}
	configHolder := NewConfigHolder(cfg)

	body := `{"data":{"type":"rooms","attributes":{"name":"Test"}}}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/rooms/missing", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("room_id")
	c.SetParamValues("missing")
	setAdminUser(c)

	h := UpdateRoomHandler(spacesStore, configHolder)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteRoomHandler_Success(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := spaces.NewStore(store)
	seedTestArea(t, store, "area-1", "Office")
	seedTestRoom(t, store, "room-1", "area-1", "Room 1")
	cfg := &spaces.Config{}
	configHolder := NewConfigHolder(cfg)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/rooms/room-1", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("room_id")
	c.SetParamValues("room-1")
	setAdminUser(c)

	h := DeleteRoomHandler(spacesStore, configHolder)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

// --- Desk Tests ---

func TestListDesksHandler_Success(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := spaces.NewStore(store)
	seedTestArea(t, store, "area-1", "Office")
	seedTestRoom(t, store, "room-1", "area-1", "Room 1")
	seedTestDesk(t, store, "desk-1", "room-1", "Desk 1")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/rooms/room-1/desks", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("room_id")
	c.SetParamValues("room-1")
	setAdminUser(c)

	h := ListDesksHandler(spacesStore)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Data, 1)
	assert.Equal(t, "desks", resp.Data[0].Type)
}

func TestCreateDeskHandler_Success(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := spaces.NewStore(store)
	seedTestArea(t, store, "area-1", "Office")
	seedTestRoom(t, store, "room-1", "area-1", "Room 1")
	cfg := &spaces.Config{}
	configHolder := NewConfigHolder(cfg)

	body := `{"data":{"type":"desks","id":"new-desk","attributes":{"name":"New Desk","equipment":["monitor"]}}}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/rooms/room-1/desks", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("room_id")
	c.SetParamValues("room-1")
	setAdminUser(c)

	h := CreateDeskHandler(spacesStore, configHolder)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestCreateDeskHandler_RoomNotFound(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := spaces.NewStore(store)
	cfg := &spaces.Config{}
	configHolder := NewConfigHolder(cfg)

	body := `{"data":{"type":"desks","attributes":{"name":"New Desk"}}}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/rooms/missing/desks", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("room_id")
	c.SetParamValues("missing")
	setAdminUser(c)

	h := CreateDeskHandler(spacesStore, configHolder)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateDeskHandler_Success(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := spaces.NewStore(store)
	seedTestArea(t, store, "area-1", "Office")
	seedTestRoom(t, store, "room-1", "area-1", "Room 1")
	seedTestDesk(t, store, "desk-1", "room-1", "Desk 1")
	cfg := &spaces.Config{}
	configHolder := NewConfigHolder(cfg)

	body := `{"data":{"type":"desks","attributes":{"name":"Updated Desk","warning":"Near exit"}}}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/desks/desk-1", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("desk_id")
	c.SetParamValues("desk-1")
	setAdminUser(c)

	h := UpdateDeskHandler(spacesStore, configHolder)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUpdateDeskHandler_NotFound(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := spaces.NewStore(store)
	cfg := &spaces.Config{}
	configHolder := NewConfigHolder(cfg)

	body := `{"data":{"type":"desks","attributes":{"name":"Test"}}}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/desks/missing", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("desk_id")
	c.SetParamValues("missing")
	setAdminUser(c)

	h := UpdateDeskHandler(spacesStore, configHolder)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteDeskHandler_Success(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := spaces.NewStore(store)
	seedTestArea(t, store, "area-1", "Office")
	seedTestRoom(t, store, "room-1", "area-1", "Room 1")
	seedTestDesk(t, store, "desk-1", "room-1", "Desk 1")
	cfg := &spaces.Config{}
	configHolder := NewConfigHolder(cfg)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/desks/desk-1", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("desk_id")
	c.SetParamValues("desk-1")
	setAdminUser(c)

	h := DeleteDeskHandler(spacesStore, configHolder)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestDeleteDeskHandler_NotFound(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := spaces.NewStore(store)
	cfg := &spaces.Config{}
	configHolder := NewConfigHolder(cfg)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/desks/missing", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("desk_id")
	c.SetParamValues("missing")
	setAdminUser(c)

	h := DeleteDeskHandler(spacesStore, configHolder)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// --- ConfigHolder Tests ---

func TestConfigHolder_GetAndReload(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	spacesStore := spaces.NewStore(store)

	// Seed initial data
	seedTestArea(t, store, "area-1", "Office")

	initialCfg := &spaces.Config{}
	holder := NewConfigHolder(initialCfg)

	// Initial config should be empty
	cfg := holder.Get()
	assert.Empty(t, cfg.Areas)

	// Reload should load data from database
	ctx := context.Background()
	err := holder.Reload(ctx, spacesStore)
	require.NoError(t, err)

	cfg = holder.Get()
	require.Len(t, cfg.Areas, 1)
	assert.Equal(t, "area-1", cfg.Areas[0].ID)
}
