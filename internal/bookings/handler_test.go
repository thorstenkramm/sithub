package bookings

import (
	"bytes"
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
	"github.com/thorstenkramm/sithub/internal/auth"
	"github.com/thorstenkramm/sithub/internal/spaces"
)

func TestCreateHandlerUnauthorized(t *testing.T) {
	t.Parallel()

	cfg := &spaces.Config{}
	store := setupTestStore(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := CreateHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var resp api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Errors, 1)
	assert.Equal(t, "auth_required", resp.Errors[0].Code)
}

func TestCreateHandlerInvalidContentType(t *testing.T) {
	t.Parallel()

	cfg := &spaces.Config{}
	store := setupTestStore(t)

	body := `{"data":{"type":"bookings","attributes":{"desk_id":"desk-1","booking_date":"2026-01-20"}}}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewBufferString(body))
	req.Header.Set(echo.HeaderContentType, "application/json") // Wrong content type
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &auth.User{ID: "user-1", Name: "Test User"})

	h := CreateHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)

	var resp api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Errors, 1)
	assert.Equal(t, "unsupported_media_type", resp.Errors[0].Code)
}

func TestCreateHandlerBadRequestCases(t *testing.T) {
	t.Parallel()

	pastDate := time.Now().UTC().AddDate(0, 0, -1).Format(time.DateOnly)

	tests := []struct {
		name           string
		body           string
		expectedDetail string
	}{
		{
			name:           "invalid JSON",
			body:           `{invalid json`,
			expectedDetail: "Invalid request body",
		},
		{
			name:           "wrong resource type",
			body:           `{"data":{"type":"wrong","attributes":{"desk_id":"desk-1","booking_date":"2026-01-20"}}}`,
			expectedDetail: "type must be 'bookings'",
		},
		{
			name:           "missing desk_id",
			body:           `{"data":{"type":"bookings","attributes":{"booking_date":"2026-01-20"}}}`,
			expectedDetail: "desk_id is required",
		},
		{
			name:           "missing booking_date",
			body:           `{"data":{"type":"bookings","attributes":{"desk_id":"desk-1"}}}`,
			expectedDetail: "booking_date is required",
		},
		{
			name: "invalid date format",
			body: `{"data":{"type":"bookings","attributes":` +
				`{"desk_id":"desk-1","booking_date":"20-01-2026"}}}`,
			expectedDetail: "YYYY-MM-DD",
		},
		{
			name: "past date",
			body: `{"data":{"type":"bookings","attributes":` +
				`{"desk_id":"desk-1","booking_date":"` + pastDate + `"}}}`,
			expectedDetail: "cannot be in the past",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cfg := &spaces.Config{}
			store := setupTestStore(t)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewBufferString(tc.body))
			req.Header.Set(echo.HeaderContentType, api.JSONAPIContentType)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("user", &auth.User{ID: "user-1", Name: "Test User"})

			h := CreateHandler(cfg, store)
			require.NoError(t, h(c))

			assert.Equal(t, http.StatusBadRequest, rec.Code)

			var resp api.ErrorResponse
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			require.Len(t, resp.Errors, 1)
			assert.Contains(t, resp.Errors[0].Detail, tc.expectedDetail)
		})
	}
}

func TestCreateHandlerDeskNotFound(t *testing.T) {
	t.Parallel()

	cfg := &spaces.Config{}
	store := setupTestStore(t)

	futureDate := time.Now().UTC().AddDate(0, 0, 1).Format(time.DateOnly)
	body := `{"data":{"type":"bookings","attributes":{"desk_id":"missing","booking_date":"` + futureDate + `"}}}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewBufferString(body))
	req.Header.Set(echo.HeaderContentType, api.JSONAPIContentType)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &auth.User{ID: "user-1", Name: "Test User"})

	h := CreateHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var resp api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Errors, 1)
	assert.Equal(t, "not_found", resp.Errors[0].Code)
}

func TestCreateHandlerSuccess(t *testing.T) {
	t.Parallel()

	cfg := testSpacesConfig()
	store := setupTestStore(t)
	seedTestDeskData(t, store, []string{"desk-1"})

	futureDate := time.Now().UTC().AddDate(0, 0, 1).Format(time.DateOnly)
	body := `{"data":{"type":"bookings","attributes":{"desk_id":"desk-1","booking_date":"` + futureDate + `"}}}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewBufferString(body))
	req.Header.Set(echo.HeaderContentType, api.JSONAPIContentType)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &auth.User{ID: "user-1", Name: "Test User"})

	h := CreateHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, api.JSONAPIContentType, rec.Header().Get(echo.HeaderContentType))

	var resp api.SingleResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "bookings", resp.Data.Type)
	assert.NotEmpty(t, resp.Data.ID)

	attrs, ok := resp.Data.Attributes.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "desk-1", attrs["desk_id"])
	assert.Equal(t, "user-1", attrs["user_id"])
	assert.Equal(t, futureDate, attrs["booking_date"])
	assert.NotEmpty(t, attrs["created_at"])
}

func TestCreateHandlerConflictCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		existingUserID string
		requestUserID  string
		expectedDetail string
	}{
		{
			name:           "desk booked by another user",
			existingUserID: "other-user",
			requestUserID:  "user-1",
			expectedDetail: "Desk is already booked for this date",
		},
		{
			name:           "self duplicate booking",
			existingUserID: "user-1",
			requestUserID:  "user-1",
			expectedDetail: "You already have this desk booked",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cfg := testSpacesConfig()
			store := setupTestStore(t)
			seedTestDeskData(t, store, []string{"desk-1"})

			futureDate := time.Now().UTC().AddDate(0, 0, 1).Format(time.DateOnly)
			seedTestBooking(t, store, "existing-booking", "desk-1", tc.existingUserID, futureDate)

			body := `{"data":{"type":"bookings","attributes":{"desk_id":"desk-1","booking_date":"` + futureDate + `"}}}`

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewBufferString(body))
			req.Header.Set(echo.HeaderContentType, api.JSONAPIContentType)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("user", &auth.User{ID: tc.requestUserID, Name: "Test User"})

			h := CreateHandler(cfg, store)
			require.NoError(t, h(c))

			assert.Equal(t, http.StatusConflict, rec.Code)

			var resp api.ErrorResponse
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			require.Len(t, resp.Errors, 1)
			assert.Equal(t, "conflict", resp.Errors[0].Code)
			assert.Contains(t, resp.Errors[0].Detail, tc.expectedDetail)
		})
	}
}

func TestListHandlerUnauthorized(t *testing.T) {
	t.Parallel()

	cfg := &spaces.Config{}
	store := setupTestStore(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/bookings", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := ListHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestListHandlerReturnsUserFutureBookings(t *testing.T) {
	t.Parallel()

	cfg := testSpacesConfig()
	store := setupTestStore(t)
	seedTestDeskData(t, store, []string{"desk-1", "desk-2"})

	today := time.Now().UTC().Format(time.DateOnly)
	tomorrow := time.Now().UTC().AddDate(0, 0, 1).Format(time.DateOnly)
	dayAfter := time.Now().UTC().AddDate(0, 0, 2).Format(time.DateOnly)
	yesterday := time.Now().UTC().AddDate(0, 0, -1).Format(time.DateOnly)
	threeDaysFromNow := time.Now().UTC().AddDate(0, 0, 3).Format(time.DateOnly)

	// User's bookings
	seedTestBooking(t, store, "b1", "desk-1", "user-1", tomorrow)
	seedTestBooking(t, store, "b2", "desk-2", "user-1", dayAfter)
	seedTestBooking(t, store, "b3", "desk-1", "user-1", today)                // Today should be included
	seedTestBooking(t, store, "b4", "desk-2", "user-1", yesterday)            // Past, should be excluded
	seedTestBooking(t, store, "b5", "desk-2", "other-user", threeDaysFromNow) // Other user, should be excluded

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/bookings", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &auth.User{ID: "user-1", Name: "Test User"})

	h := ListHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, api.JSONAPIContentType, rec.Header().Get(echo.HeaderContentType))

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

	// Should return 3 bookings (today, tomorrow, day after), ordered by date
	require.Len(t, resp.Data, 3)

	// First booking should be today
	attrs0, ok := resp.Data[0].Attributes.(map[string]interface{})
	require.True(t, ok, "failed to cast attributes")
	assert.Equal(t, today, attrs0["booking_date"])
	assert.Equal(t, "desk-1", attrs0["desk_id"])
	assert.Equal(t, "Desk 1", attrs0["desk_name"])
	assert.Equal(t, "room-1", attrs0["room_id"])
	assert.Equal(t, "Room 1", attrs0["room_name"])
	assert.Equal(t, "area-1", attrs0["area_id"])
	assert.Equal(t, "Office", attrs0["area_name"])

	// Second booking should be tomorrow
	attrs1, ok := resp.Data[1].Attributes.(map[string]interface{})
	require.True(t, ok, "failed to cast attributes")
	assert.Equal(t, tomorrow, attrs1["booking_date"])

	// Third booking should be day after
	attrs2, ok := resp.Data[2].Attributes.(map[string]interface{})
	require.True(t, ok, "failed to cast attributes")
	assert.Equal(t, dayAfter, attrs2["booking_date"])
}

func TestListHandlerEmptyList(t *testing.T) {
	t.Parallel()

	cfg := testSpacesConfig()
	store := setupTestStore(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/bookings", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &auth.User{ID: "user-1", Name: "Test User"})

	h := ListHandler(cfg, store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Len(t, resp.Data, 0)
}

func TestDeleteHandlerUnauthorized(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/bookings/booking-1", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("booking-1")

	h := DeleteHandler(store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestDeleteHandlerNotFound(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/bookings/nonexistent", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("nonexistent")
	c.Set("user", &auth.User{ID: "user-1", Name: "Test User"})

	h := DeleteHandler(store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteHandlerOtherUsersBooking(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	seedTestDeskData(t, store, []string{"desk-1"})

	tomorrow := time.Now().UTC().AddDate(0, 0, 1).Format(time.DateOnly)
	seedTestBooking(t, store, "booking-1", "desk-1", "other-user", tomorrow)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/bookings/booking-1", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("booking-1")
	c.Set("user", &auth.User{ID: "user-1", Name: "Test User"})

	h := DeleteHandler(store)
	require.NoError(t, h(c))

	// Should return 404 to not reveal booking existence
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteHandlerSuccess(t *testing.T) {
	t.Parallel()

	store := setupTestStore(t)
	seedTestDeskData(t, store, []string{"desk-1"})

	tomorrow := time.Now().UTC().AddDate(0, 0, 1).Format(time.DateOnly)
	seedTestBooking(t, store, "booking-1", "desk-1", "user-1", tomorrow)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/bookings/booking-1", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("booking-1")
	c.Set("user", &auth.User{ID: "user-1", Name: "Test User"})

	h := DeleteHandler(store)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Verify booking is deleted
	ctx := context.Background()
	booking, err := FindBookingByID(ctx, store, "booking-1")
	require.NoError(t, err)
	assert.Nil(t, booking)
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
							{
								ID:        "desk-1",
								Name:      "Desk 1",
								Equipment: []string{"Monitor"},
							},
							{
								ID:        "desk-2",
								Name:      "Desk 2",
								Equipment: []string{"Monitor"},
							},
						},
					},
				},
			},
		},
	}
}
