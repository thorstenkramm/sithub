package users

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/notifications"
	"github.com/thorstenkramm/sithub/internal/testutil/dbtest"
)

// testUser mirrors auth.User for setting context in tests.
// Implements the GetID interface used by currentUserID().
type testUser struct {
	ID string
}

func (u *testUser) GetID() string { return u.ID }

func setupHandlerDB(t *testing.T) *sql.DB {
	t.Helper()
	return dbtest.New(t)
}

func seedUser(t *testing.T, db *sql.DB, email, displayName, source string) *Record {
	t.Helper()
	hash := ""
	if source == "internal" {
		var err error
		hash, err = HashPassword("TestPassword123!")
		require.NoError(t, err)
	}
	rec, err := CreateLocalUser(t.Context(), db, email, displayName, hash, false)
	if source == "entraid" {
		// Override source for entraid test users
		_, err = db.Exec("UPDATE users SET user_source = 'entraid', entra_id = ? WHERE id = ?", "entra-"+rec.ID, rec.ID)
		require.NoError(t, err)
		rec.UserSource = "entraid"
		rec.EntraID = "entra-" + rec.ID
	}
	require.NoError(t, err)
	return rec
}

func ptr(s string) *string { return &s }

// seedUserRows inserts n users named "User 001".."User n" without bcrypt hashing.
func seedUserRows(t *testing.T, db *sql.DB, n int) {
	t.Helper()
	for i := 1; i <= n; i++ {
		id := fmt.Sprintf("user-%03d", i)
		_, err := db.Exec(`INSERT INTO users (id, email, display_name, user_source, created_at, updated_at)
			VALUES (?, ?, ?, 'internal', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z')`,
			id, id+"@test.com", fmt.Sprintf("User %03d", i))
		require.NoError(t, err)
	}
}

func TestListHandler(t *testing.T) {
	t.Parallel()

	link := func(limit, offset int) string {
		return fmt.Sprintf("/api/v1/users?page%%5Blimit%%5D=%d&page%%5Boffset%%5D=%d", limit, offset)
	}

	cases := []struct {
		name      string
		seed      int
		query     string
		wantFirst string
		wantLen   int
		wantMeta  api.PaginationMeta
		wantLinks api.PaginationLinks
	}{
		{
			name: "default page without params", seed: 45, query: "",
			wantFirst: "User 001", wantLen: 20,
			wantMeta: api.PaginationMeta{Total: 45, Limit: 20, Offset: 0},
			wantLinks: api.PaginationLinks{
				Self: link(20, 0), First: link(20, 0), Next: ptr(link(20, 20)), Last: link(20, 40),
			},
		},
		{
			name: "requested slice", seed: 45, query: "page[limit]=20&page[offset]=40",
			wantFirst: "User 041", wantLen: 5,
			wantMeta: api.PaginationMeta{Total: 45, Limit: 20, Offset: 40},
			wantLinks: api.PaginationLinks{
				Self: link(20, 40), First: link(20, 0), Prev: ptr(link(20, 20)), Last: link(20, 40),
			},
		},
		{
			name: "limit capped at 100", seed: 120, query: "page[limit]=500",
			wantFirst: "User 001", wantLen: 100,
			wantMeta: api.PaginationMeta{Total: 120, Limit: 100, Offset: 0},
			wantLinks: api.PaginationLinks{
				Self: link(100, 0), First: link(100, 0), Next: ptr(link(100, 100)), Last: link(100, 100),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := setupHandlerDB(t)
			seedUserRows(t, db, tc.seed)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/users?"+tc.query, http.NoBody)
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)

			require.NoError(t, ListHandler(db)(c))
			assert.Equal(t, http.StatusOK, rec.Code)

			var resp struct {
				Data []struct {
					Attributes UserAttributes `json:"attributes"`
				} `json:"data"`
				Meta  api.PaginationMeta  `json:"meta"`
				Links api.PaginationLinks `json:"links"`
			}
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			require.Len(t, resp.Data, tc.wantLen)
			assert.Equal(t, tc.wantFirst, resp.Data[0].Attributes.DisplayName)
			assert.Equal(t, tc.wantMeta, resp.Meta)
			assert.Equal(t, tc.wantLinks, resp.Links)
		})
	}
}

func TestListHandler_InvalidPagination(t *testing.T) {
	t.Parallel()

	cases := []struct {
		query string
		param string
	}{
		{query: "page[limit]=abc", param: "page[limit]"},
		{query: "page[limit]=0", param: "page[limit]"},
		{query: "page[offset]=-5", param: "page[offset]"},
	}

	for _, tc := range cases {
		t.Run(tc.query, func(t *testing.T) {
			t.Parallel()
			db := setupHandlerDB(t)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/users?"+tc.query, http.NoBody)
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)

			require.NoError(t, ListHandler(db)(c))
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			var resp api.ErrorResponse
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			require.Len(t, resp.Errors, 1)
			assert.Equal(t, "400", resp.Errors[0].Status)
			assert.Contains(t, resp.Errors[0].Detail, tc.param)
		})
	}
}

func TestGetHandler(t *testing.T) {
	db := setupHandlerDB(t)
	user := seedUser(t, db, "alice@test.com", "Alice", "internal")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+user.ID, http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(user.ID)

	err := GetHandler(db)(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp api.SingleResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, user.ID, resp.Data.ID)
}

func TestGetHandlerNotFound(t *testing.T) {
	db := setupHandlerDB(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/nonexistent", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("nonexistent")

	err := GetHandler(db)(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCreateHandler(t *testing.T) {
	db := setupHandlerDB(t)

	body := `{"data":{"type":"users","attributes":{"email":"new@test.com",` +
		`"display_name":"New User","password":"SuperSecure12345!","is_admin":false}}}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, api.JSONAPIContentType)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := CreateHandler(db)(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp api.SingleResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "users", resp.Data.Type)
	assert.NotEmpty(t, resp.Data.ID)
}

func TestCreateHandlerDuplicateEmail(t *testing.T) {
	db := setupHandlerDB(t)
	seedUser(t, db, "dup@test.com", "Existing", "internal")

	body := `{"data":{"type":"users","attributes":{"email":"dup@test.com",` +
		`"display_name":"Duplicate","password":"SuperSecure12345!","is_admin":false}}}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, api.JSONAPIContentType)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := CreateHandler(db)(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestCreateHandlerShortPassword(t *testing.T) {
	db := setupHandlerDB(t)

	body := `{"data":{"type":"users","attributes":{"email":"new@test.com","display_name":"New","password":"short"}}}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, api.JSONAPIContentType)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := CreateHandler(db)(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateHandler(t *testing.T) {
	db := setupHandlerDB(t)
	user := seedUser(t, db, "alice@test.com", "Alice", "internal")

	body := `{"data":{"attributes":{"display_name":"Alice Updated","is_admin":true}}}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/users/"+user.ID, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, api.JSONAPIContentType)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(user.ID)

	err := UpdateHandler(db)(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify the update
	updated, err := FindByID(t.Context(), db, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "Alice Updated", updated.DisplayName)
	assert.True(t, updated.IsAdmin)
}

func TestUpdateHandlerPasswordReset(t *testing.T) {
	db := setupHandlerDB(t)
	user := seedUser(t, db, "alice@test.com", "Alice", "internal")

	body := `{"data":{"attributes":{"password":"NewSecurePassword!!"}}}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/users/"+user.ID, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, api.JSONAPIContentType)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(user.ID)

	err := UpdateHandler(db)(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify new password works
	updated, err := FindByID(t.Context(), db, user.ID)
	require.NoError(t, err)
	assert.NoError(t, VerifyPassword(updated.PasswordHash, "NewSecurePassword!!"))
}

func TestUpdateHandlerPasswordResetEntraIDUser(t *testing.T) {
	db := setupHandlerDB(t)
	user := seedUser(t, db, "entra@test.com", "Entra User", "entraid")

	body := `{"data":{"attributes":{"password":"NewSecurePassword!!"}}}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/users/"+user.ID, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, api.JSONAPIContentType)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(user.ID)

	err := UpdateHandler(db)(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateHandlerNotFound(t *testing.T) {
	db := setupHandlerDB(t)

	body := `{"data":{"attributes":{"display_name":"Nobody"}}}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/users/unknown", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, api.JSONAPIContentType)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("unknown")

	err := UpdateHandler(db)(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateHandlerDataExporter(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		body    string
		initial int
		want    bool
	}{
		{name: "grant", body: `{"data":{"attributes":{"is_data_exporter":true}}}`, initial: 0, want: true},
		{name: "revoke", body: `{"data":{"attributes":{"is_data_exporter":false}}}`, initial: 1, want: false},
		{
			name: "omitted keeps flag", initial: 1, want: true,
			body: `{"data":{"attributes":{"display_name":"Renamed"}}}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := setupHandlerDB(t)
			insertUserRow(t, db, "u1", "User One")
			_, err := db.Exec(`UPDATE users SET is_data_exporter = ? WHERE id = 'u1'`, tc.initial)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch, "/api/v1/users/u1", strings.NewReader(tc.body))
			req.Header.Set(echo.HeaderContentType, api.JSONAPIContentType)
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues("u1")

			require.NoError(t, UpdateHandler(db)(c))
			assert.Equal(t, http.StatusOK, rec.Code)

			var resp struct {
				Data struct {
					Attributes UserAttributes `json:"attributes"`
				} `json:"data"`
			}
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			assert.Equal(t, tc.want, resp.Data.Attributes.IsDataExporter)
			assert.Contains(t, rec.Body.String(), `"is_data_exporter"`)
		})
	}
}

func TestListAndGetHandlersExposeDataExporter(t *testing.T) {
	t.Parallel()
	db := setupHandlerDB(t)
	insertUserRow(t, db, "u1", "User One")
	_, err := db.Exec(`UPDATE users SET is_data_exporter = 1 WHERE id = 'u1'`)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", http.NoBody)
	rec := httptest.NewRecorder()
	require.NoError(t, ListHandler(db)(echo.New().NewContext(req, rec)))
	assert.Contains(t, rec.Body.String(), `"is_data_exporter":true`)

	req = httptest.NewRequest(http.MethodGet, "/api/v1/users/u1", http.NoBody)
	rec = httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("u1")
	require.NoError(t, GetHandler(db)(c))
	assert.Contains(t, rec.Body.String(), `"is_data_exporter":true`)
}

// recordingNotifier captures booking events synchronously.
type recordingNotifier struct {
	mu     sync.Mutex
	events []notifications.BookingEvent
}

func (n *recordingNotifier) NotifyAsync(event *notifications.BookingEvent) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.events = append(n.events, *event)
}

func newDeleteContext(t *testing.T, method, path, userID, callerID string) (echo.Context, *httptest.ResponseRecorder) {
	t.Helper()
	req := httptest.NewRequest(method, path, http.NoBody)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(userID)
	if callerID != "" {
		c.Set("user", &testUser{ID: callerID})
	}
	return c, rec
}

func TestDeleteHandler(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		source       string
		bookings     []seedBooking
		wantStatus   int
		wantDeleted  bool
		wantCanceled []string
	}{
		{name: "local user without bookings is hard deleted", source: "internal", wantStatus: http.StatusNoContent},
		{name: "entra user without bookings is hard deleted", source: "entraid", wantStatus: http.StatusNoContent},
		{
			name: "user with bookings is deactivated and future bookings canceled", source: "entraid",
			bookings: []seedBooking{
				{id: "past", userID: "x", bookedBy: "x", date: "2000-01-01"},
				{id: "future", userID: "x", bookedBy: "x", date: "2099-01-01"},
			},
			wantStatus: http.StatusOK, wantDeleted: true, wantCanceled: []string{"future"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := setupHandlerDB(t)
			insertUserRow(t, db, "x", "Xavier")
			_, err := db.Exec(`UPDATE users SET user_source = ? WHERE id = 'x'`, tc.source)
			require.NoError(t, err)
			for _, b := range tc.bookings {
				insertBookingRow(t, db, b)
			}
			notifier := &recordingNotifier{}

			c, rec := newDeleteContext(t, http.MethodDelete, "/api/v1/users/x", "x", "admin-user")
			require.NoError(t, DeleteHandler(db, notifier)(c))
			assert.Equal(t, tc.wantStatus, rec.Code)

			canceled := make([]string, 0, len(notifier.events))
			for _, ev := range notifier.events {
				assert.Equal(t, notifications.EventBookingCanceled, ev.Event)
				assert.Equal(t, "admin-user", ev.CanceledByUserID)
				canceled = append(canceled, ev.BookingID)
			}
			if tc.wantCanceled == nil {
				tc.wantCanceled = []string{}
			}
			assert.Equal(t, tc.wantCanceled, canceled)

			if !tc.wantDeleted {
				assert.Empty(t, rec.Body.String())
				_, err := FindByID(t.Context(), db, "x")
				assert.ErrorIs(t, err, ErrUserNotFound)
				return
			}
			var resp struct {
				Data struct {
					Attributes UserAttributes `json:"attributes"`
				} `json:"data"`
			}
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			assert.True(t, resp.Data.Attributes.Deleted)
		})
	}
}

func TestDeleteHandlerPreventsSelfDeletion(t *testing.T) {
	t.Parallel()
	db := setupHandlerDB(t)
	insertUserRow(t, db, "admin", "Admin")
	insertBookingRow(t, db, seedBooking{id: "b1", userID: "admin", bookedBy: "admin", date: "2099-01-01"})
	notifier := &recordingNotifier{}

	c, rec := newDeleteContext(t, http.MethodDelete, "/api/v1/users/admin", "admin", "admin")
	require.NoError(t, DeleteHandler(db, notifier)(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Cannot delete your own account")

	found, err := FindByID(t.Context(), db, "admin")
	require.NoError(t, err)
	assert.False(t, found.Deleted)
	assert.Equal(t, []string{"b1"}, bookingIDs(t, db))
	assert.Empty(t, notifier.events)
}

func TestDeleteHandlerNotFound(t *testing.T) {
	t.Parallel()
	db := setupHandlerDB(t)

	c, rec := newDeleteContext(t, http.MethodDelete, "/api/v1/users/nonexistent", "nonexistent", "admin-user")
	require.NoError(t, DeleteHandler(db, &recordingNotifier{})(c))
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeletionImpactHandler(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		bookings   []seedBooking
		userID     string
		wantStatus int
		wantAttrs  DeletionImpactAttributes
	}{
		{
			name: "no bookings", userID: "x", wantStatus: http.StatusOK,
			wantAttrs: DeletionImpactAttributes{Mode: "delete"},
		},
		{
			name: "with bookings", userID: "x", wantStatus: http.StatusOK,
			bookings: []seedBooking{
				{id: "past", userID: "x", bookedBy: "x", date: "2000-01-01"},
				{id: "f1", userID: "x", bookedBy: "x", date: "2099-01-01"},
				{id: "f2", userID: "guest-1", bookedBy: "x", date: "2099-01-02", isGuest: true},
			},
			wantAttrs: DeletionImpactAttributes{HasBookings: true, FutureBookingCount: 2, Mode: "deactivate"},
		},
		{name: "unknown user", userID: "missing", wantStatus: http.StatusNotFound},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := setupHandlerDB(t)
			insertUserRow(t, db, "x", "Xavier")
			for _, b := range tc.bookings {
				insertBookingRow(t, db, b)
			}

			path := "/api/v1/users/" + tc.userID + "/deletion-impact"
			c, rec := newDeleteContext(t, http.MethodGet, path, tc.userID, "admin-user")
			require.NoError(t, DeletionImpactHandler(db)(c))
			assert.Equal(t, tc.wantStatus, rec.Code)
			if tc.wantStatus != http.StatusOK {
				return
			}

			var resp struct {
				Data struct {
					ID         string                   `json:"id"`
					Type       string                   `json:"type"`
					Attributes DeletionImpactAttributes `json:"attributes"`
				} `json:"data"`
			}
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			assert.Equal(t, "x", resp.Data.ID)
			assert.Equal(t, "user-deletion-impacts", resp.Data.Type)
			assert.Equal(t, tc.wantAttrs, resp.Data.Attributes)
		})
	}
}
