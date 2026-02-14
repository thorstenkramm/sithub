package itemgroups

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/api"
)

func TestParseISOWeek(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		wantDate  string
		wantErr   bool
	}{
		{"valid week 1 2026", "2026-W01", "2025-12-29", false},
		{"valid week 12 2026", "2026-W12", "2026-03-16", false},
		{"valid week with leading zero", "2026-W03", "2026-01-12", false},
		{"valid week 53 2026", "2026-W53", "2026-12-28", false},
		{"empty defaults to current week", "", "", false},
		{"invalid format", "2026-12", "", true},
		{"invalid week 53 for non-53 year", "2025-W53", "", true},
		{"invalid week number", "2026-W54", "", true},
		{"invalid week zero", "2026-W00", "", true},
		{"invalid year", "abc-W12", "", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, err := parseISOWeek(tc.input)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tc.input == "" {
				// Just verify it returns a Monday
				assert.Equal(t, time.Monday, result.Weekday())
			} else {
				assert.Equal(t, tc.wantDate, result.Format(time.DateOnly))
				assert.Equal(t, time.Monday, result.Weekday())
			}
		})
	}
}

func TestWeekdayDates(t *testing.T) {
	t.Parallel()

	monday := time.Date(2026, 3, 16, 0, 0, 0, 0, time.UTC)
	days := weekdayDates(monday)

	require.Len(t, days, 5)
	assert.Equal(t, "2026-03-16", days[0].Format(time.DateOnly))
	assert.Equal(t, "2026-03-17", days[1].Format(time.DateOnly))
	assert.Equal(t, "2026-03-18", days[2].Format(time.DateOnly))
	assert.Equal(t, "2026-03-19", days[3].Format(time.DateOnly))
	assert.Equal(t, "2026-03-20", days[4].Format(time.DateOnly))
}

func TestWeekdayAbbreviation(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "MO", weekdayAbbreviation(time.Monday))
	assert.Equal(t, "TU", weekdayAbbreviation(time.Tuesday))
	assert.Equal(t, "WE", weekdayAbbreviation(time.Wednesday))
	assert.Equal(t, "TH", weekdayAbbreviation(time.Thursday))
	assert.Equal(t, "FR", weekdayAbbreviation(time.Friday))
}

func TestAvailabilityHandlerAreaNotFound(t *testing.T) {
	t.Parallel()

	store := setupTestDB(t)
	cfg := testConfig()

	e := echo.New()
	url := "/api/v1/areas/unknown/item-groups/availability?week=2026-W04"
	req := httptest.NewRequest(http.MethodGet, url, http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("area_id")
	c.SetParamValues("unknown")

	h := AvailabilityHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAvailabilityHandlerInvalidWeek(t *testing.T) {
	t.Parallel()

	store := setupTestDB(t)
	cfg := testConfig()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/areas/area-1/item-groups/availability?week=bad", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("area_id")
	c.SetParamValues("area-1")

	h := AvailabilityHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAvailabilityHandlerReturnsPerDayAvailability(t *testing.T) {
	t.Parallel()

	store := setupTestDB(t)
	cfg := testConfig()

	// Book item-1 on Monday 2026-01-19 and item-2 on Monday 2026-01-19
	// This means ig-1 is fully booked on Monday (2 of 2 items booked)
	seedBooking(t, store, "b1", "item-1", "user-1", "2026-01-19")
	seedBooking(t, store, "b2", "item-2", "user-1", "2026-01-19")
	// Book item-1 on Tuesday 2026-01-20 (1 of 2 booked)
	seedBooking(t, store, "b3", "item-1", "user-1", "2026-01-20")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet,
		"/api/v1/areas/area-1/item-groups/availability?week=2026-W04", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("area_id")
	c.SetParamValues("area-1")

	h := AvailabilityHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get(echo.HeaderContentType), api.JSONAPIContentType)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Data, 2) // 2 item groups in area-1

	// ig-1: 2 items total
	assert.Equal(t, "item-group-availability", resp.Data[0].Type)
	assert.Equal(t, "ig-1", resp.Data[0].ID)

	attrs0, ok := resp.Data[0].Attributes.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "Room 101", attrs0["item_group_name"])

	days0, ok := attrs0["days"].([]any)
	require.True(t, ok)
	require.Len(t, days0, 5)

	// Monday: 2 booked, 0 available
	mon, ok := days0[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "2026-01-19", mon["date"])
	assert.Equal(t, "MO", mon["weekday"])
	assert.Equal(t, float64(2), mon["total"])
	assert.Equal(t, float64(0), mon["available"])

	// Tuesday: 1 booked, 1 available
	tue, ok := days0[1].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "2026-01-20", tue["date"])
	assert.Equal(t, float64(1), tue["available"])

	// Wednesday: 0 booked, 2 available
	wed, ok := days0[2].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "2026-01-21", wed["date"])
	assert.Equal(t, float64(2), wed["available"])

	// ig-2: 1 item total, no bookings
	attrs1, ok := resp.Data[1].Attributes.(map[string]any)
	require.True(t, ok)
	days1, ok := attrs1["days"].([]any)
	require.True(t, ok)
	for _, d := range days1 {
		day, ok := d.(map[string]any)
		require.True(t, ok)
		assert.Equal(t, float64(1), day["total"])
		assert.Equal(t, float64(1), day["available"])
	}
}

func TestAvailabilityHandlerNoBookings(t *testing.T) {
	t.Parallel()

	store := setupTestDB(t)
	cfg := testConfig()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet,
		"/api/v1/areas/area-1/item-groups/availability?week=2026-W04", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("area_id")
	c.SetParamValues("area-1")

	h := AvailabilityHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Data, 2)

	// All items should be available
	attrs0, ok := resp.Data[0].Attributes.(map[string]any)
	require.True(t, ok)
	days0, ok := attrs0["days"].([]any)
	require.True(t, ok)
	for _, d := range days0 {
		day, ok := d.(map[string]any)
		require.True(t, ok)
		assert.Equal(t, float64(2), day["total"])
		assert.Equal(t, float64(2), day["available"])
	}
}

func TestAvailabilityHandlerDefaultsToCurrentWeek(t *testing.T) {
	t.Parallel()

	store := setupTestDB(t)
	cfg := testConfig()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet,
		"/api/v1/areas/area-1/item-groups/availability", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("area_id")
	c.SetParamValues("area-1")

	h := AvailabilityHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Data, 2)

	// Verify the first day is a Monday
	attrs0, ok := resp.Data[0].Attributes.(map[string]any)
	require.True(t, ok)
	days0, ok := attrs0["days"].([]any)
	require.True(t, ok)
	firstDay, ok := days0[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "MO", firstDay["weekday"])
}
