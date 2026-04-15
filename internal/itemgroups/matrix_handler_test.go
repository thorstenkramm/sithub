package itemgroups

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/areas"
	"github.com/thorstenkramm/sithub/internal/auth"
)

// matrixTestConfig returns a config with reservations and equipment for matrix tests.
func matrixTestConfig() *areas.Config {
	return &areas.Config{
		Areas: []areas.Area{
			{
				ID:   "area-1",
				Name: "Area One",
				ItemGroups: []areas.ItemGroup{
					{
						ID:   "ig-1",
						Name: "Room 101",
						Items: []areas.Item{
							{
								ID:        "item-1",
								Name:      "Desk 1",
								Equipment: []string{"Dock", "Monitor"},
								Warning:   "Near window",
							},
							{ID: "item-2", Name: "Desk 2", Equipment: []string{"Dock"}},
						},
					},
					{
						ID:   "ig-2",
						Name: "Room 102",
						Items: []areas.Item{
							{ID: "item-3", Name: "Desk 3"},
						},
					},
				},
			},
			{
				ID:   "area-reserved",
				Name: "Reserved Area",
				ItemGroups: []areas.ItemGroup{
					{
						ID:   "ig-reserved",
						Name: "VIP Room",
						Items: []areas.Item{
							{
								ID:          "item-r1",
								Name:        "VIP Desk",
								ReservedFor: []string{"vip@test.local"},
							},
						},
					},
				},
			},
		},
	}
}

func newMatrixRequest(t *testing.T, url, areaID string, user *auth.User) (echo.Context, *httptest.ResponseRecorder) {
	t.Helper()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, url, http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("area_id")
	c.SetParamValues(areaID)
	if user != nil {
		c.Set("user", user)
	}
	return c, rec
}

func TestMatrixHandlerAreaNotFound(t *testing.T) {
	t.Parallel()
	store := setupTestDB(t)
	cfg := matrixTestConfig()

	c, rec := newMatrixRequest(t,
		"/api/v1/areas/unknown/item-groups/matrix?week=2026-W04", "unknown", nil)

	h := MatrixHandler(cfg, store)
	require.NoError(t, h(c))
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestMatrixHandlerInvalidWeek(t *testing.T) {
	t.Parallel()
	store := setupTestDB(t)
	cfg := matrixTestConfig()

	c, rec := newMatrixRequest(t,
		"/api/v1/areas/area-1/item-groups/matrix?week=bad", "area-1", nil)

	h := MatrixHandler(cfg, store)
	require.NoError(t, h(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestMatrixHandlerFiveDayOutput(t *testing.T) {
	t.Parallel()
	store := setupTestDB(t)
	cfg := matrixTestConfig()
	user := &auth.User{ID: "user-1"}
	seedUser(t, store, "user-1", "Ada Lovelace")

	c, rec := newMatrixRequest(t,
		"/api/v1/areas/area-1/item-groups/matrix?week=2026-W04", "area-1", user)

	h := MatrixHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get(echo.HeaderContentType), api.JSONAPIContentType)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Data, 2) // ig-1 and ig-2

	attrs0 := resourceAttrs(t, resp.Data[0])
	days := attrSlice(t, attrs0, "days")
	require.Len(t, days, 5) // Mon-Fri only
	assertDay(t, days[0], "2026-01-19", "MO")
	assertDay(t, days[4], "2026-01-23", "FR")
}

func TestMatrixHandlerSevenDayOutput(t *testing.T) {
	t.Parallel()
	store := setupTestDB(t)
	cfg := matrixTestConfig()

	c, rec := newMatrixRequest(t,
		"/api/v1/areas/area-1/item-groups/matrix?week=2026-W04&days=7", "area-1", nil)

	h := MatrixHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

	attrs0 := resourceAttrs(t, resp.Data[0])
	days := attrSlice(t, attrs0, "days")
	require.Len(t, days, 7)
	assertDay(t, days[5], "2026-01-24", "SA")
	assertDay(t, days[6], "2026-01-25", "SU")
}

func TestMatrixHandlerConfiguredOrder(t *testing.T) {
	t.Parallel()
	store := setupTestDB(t)
	cfg := matrixTestConfig()

	c, rec := newMatrixRequest(t,
		"/api/v1/areas/area-1/item-groups/matrix?week=2026-W04", "area-1", nil)

	h := MatrixHandler(cfg, store)
	require.NoError(t, h(c))

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

	// Verify item group order matches config
	assert.Equal(t, "ig-1", resp.Data[0].ID)
	assert.Equal(t, "ig-2", resp.Data[1].ID)
	assert.Equal(t, matrixResourceType, resp.Data[0].Type)

	// Verify item order within ig-1
	attrs0 := resourceAttrs(t, resp.Data[0])
	items := attrSlice(t, attrs0, "items")
	require.Len(t, items, 2)
	assert.Equal(t, "item-1", itemAt(t, items, 0)["item_id"])
	assert.Equal(t, "Desk 1", itemAt(t, items, 0)["item_name"])
	assert.Equal(t, "item-2", itemAt(t, items, 1)["item_id"])
}

func TestMatrixHandlerOccupiedCellsWithDisplayName(t *testing.T) {
	t.Parallel()
	store := setupTestDB(t)
	cfg := matrixTestConfig()

	seedUser(t, store, "user-1", "Ada Lovelace")
	seedBooking(t, store, "b1", "item-1", "user-1", "2026-01-19")

	user := &auth.User{ID: "user-1"}
	c, rec := newMatrixRequest(t,
		"/api/v1/areas/area-1/item-groups/matrix?week=2026-W04", "area-1", user)

	h := MatrixHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

	attrs0 := resourceAttrs(t, resp.Data[0])
	items := attrSlice(t, attrs0, "items")
	item0 := itemAt(t, items, 0)
	cells := cellsOf(t, item0)
	monCell := cellAt(t, cells, 0)

	assert.Equal(t, "occupied", monCell["availability"])
	assert.Equal(t, "Ada Lovelace", monCell["booker_name"])
	assert.Equal(t, "user-1", monCell["booker_user_id"])
	assert.Equal(t, true, monCell["booked_by_me"])
	// Owner sees booking_id
	assert.Equal(t, "b1", monCell["booking_id"])

	// Tuesday should be free
	tueCell := cellAt(t, cells, 1)
	assert.Equal(t, "free", tueCell["availability"])
	assert.Equal(t, false, tueCell["booked_by_me"])
}

func TestMatrixHandlerGuestBookingShowsGuestName(t *testing.T) {
	t.Parallel()
	store := setupTestDB(t)
	cfg := matrixTestConfig()

	// Seed a guest booking
	now := time.Now().Format(time.RFC3339)
	_, err := store.ExecContext(context.Background(),
		`INSERT INTO bookings (id, item_id, user_id, booking_date, is_guest, guest_name, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		"bg1", "item-1", "", "2026-01-19", 1, "Guest Person", now, now)
	require.NoError(t, err)

	c, rec := newMatrixRequest(t,
		"/api/v1/areas/area-1/item-groups/matrix?week=2026-W04", "area-1", nil)

	h := MatrixHandler(cfg, store)
	require.NoError(t, h(c))

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

	attrs0 := resourceAttrs(t, resp.Data[0])
	items := attrSlice(t, attrs0, "items")
	item0 := itemAt(t, items, 0)
	cells := cellsOf(t, item0)
	monCell := cellAt(t, cells, 0)

	assert.Equal(t, "occupied", monCell["availability"])
	assert.Equal(t, "Guest Person", monCell["booker_name"])
	// Guest bookings should not have booker_user_id
	assert.Empty(t, monCell["booker_user_id"])
}

func TestMatrixHandlerReservedFreeCell(t *testing.T) {
	t.Parallel()
	store := setupTestDB(t)
	cfg := matrixTestConfig()

	// Create user who is NOT in the reserved list
	seedUser(t, store, "outsider", "Outsider User")
	user := &auth.User{ID: "outsider"}

	c, rec := newMatrixRequest(t,
		"/api/v1/areas/area-reserved/item-groups/matrix?week=2026-W04", "area-reserved", user)

	h := MatrixHandler(cfg, store)
	require.NoError(t, h(c))

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Data, 1)

	attrs := resourceAttrs(t, resp.Data[0])
	items := attrSlice(t, attrs, "items")

	// Item should be marked reserved since user is not in reserved_for list
	assert.Equal(t, true, itemAt(t, items, 0)["reserved"])
}

func TestMatrixHandlerBookingIDVisibility(t *testing.T) {
	t.Parallel()
	store := setupTestDB(t)
	cfg := matrixTestConfig()

	seedUser(t, store, "owner", "Owner User")
	seedUser(t, store, "viewer", "Viewer User")
	seedBooking(t, store, "b1", "item-1", "owner", "2026-01-19")

	tests := []struct {
		name     string
		user     *auth.User
		expectID bool
	}{
		{"owner sees booking_id", &auth.User{ID: "owner"}, true},
		{"admin sees booking_id", &auth.User{ID: "viewer", IsAdmin: true}, true},
		{"other user does not see booking_id", &auth.User{ID: "viewer"}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c, rec := newMatrixRequest(t,
				"/api/v1/areas/area-1/item-groups/matrix?week=2026-W04", "area-1", tc.user)

			h := MatrixHandler(cfg, store)
			require.NoError(t, h(c))

			var resp api.CollectionResponse
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

			attrs0 := resourceAttrs(t, resp.Data[0])
			items := attrSlice(t, attrs0, "items")
			item0 := itemAt(t, items, 0)
			cells := cellsOf(t, item0)
			monCell := cellAt(t, cells, 0)

			if tc.expectID {
				assert.Equal(t, "b1", monCell["booking_id"])
			} else {
				assert.Empty(t, monCell["booking_id"])
			}
		})
	}
}

func TestMatrixHandlerEquipmentAndWarning(t *testing.T) {
	t.Parallel()
	store := setupTestDB(t)
	cfg := matrixTestConfig()

	c, rec := newMatrixRequest(t,
		"/api/v1/areas/area-1/item-groups/matrix?week=2026-W04", "area-1", nil)

	h := MatrixHandler(cfg, store)
	require.NoError(t, h(c))

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

	attrs0 := resourceAttrs(t, resp.Data[0])
	items := attrSlice(t, attrs0, "items")
	item0 := itemAt(t, items, 0)

	// item-1 has equipment and warning
	equip := equipOf(t, item0)
	assert.Len(t, equip, 2)
	assert.Equal(t, "Dock", equip[0])
	assert.Equal(t, "Monitor", equip[1])
	assert.Equal(t, "Near window", item0["warning"])

	// item-3 (ig-2) has no equipment
	attrs1 := resourceAttrs(t, resp.Data[1])
	items1 := attrSlice(t, attrs1, "items")
	equip3 := equipOf(t, itemAt(t, items1, 0))
	assert.Empty(t, equip3)
}

// --- test helpers ---

func resourceAttrs(t *testing.T, r api.Resource) map[string]any {
	t.Helper()
	attrs, ok := r.Attributes.(map[string]any)
	require.True(t, ok)
	return attrs
}

func attrSlice(t *testing.T, attrs map[string]any, key string) []any {
	t.Helper()
	s, ok := attrs[key].([]any)
	require.True(t, ok, "expected []any for key %q", key)
	return s
}

func assertDay(t *testing.T, day any, date, weekday string) {
	t.Helper()
	d, ok := day.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, date, d["date"])
	assert.Equal(t, weekday, d["weekday"])
}

func itemAt(t *testing.T, items []any, idx int) map[string]any {
	t.Helper()
	item, ok := items[idx].(map[string]any)
	require.True(t, ok, "expected map at items[%d]", idx)
	return item
}

func cellsOf(t *testing.T, item map[string]any) []any {
	t.Helper()
	cells, ok := item["cells"].([]any)
	require.True(t, ok, "expected cells array")
	return cells
}

func cellAt(t *testing.T, cells []any, idx int) map[string]any {
	t.Helper()
	cell, ok := cells[idx].(map[string]any)
	require.True(t, ok, "expected map at cells[%d]", idx)
	return cell
}

func equipOf(t *testing.T, item map[string]any) []any {
	t.Helper()
	equip, ok := item["equipment"].([]any)
	require.True(t, ok, "expected equipment array")
	return equip
}
