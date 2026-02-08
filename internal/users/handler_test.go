package users

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/api"
)

// testUser mirrors auth.User for setting context in tests.
// Implements the GetID interface used by currentUserID().
type testUser struct {
	ID string
}

func (u *testUser) GetID() string { return u.ID }

func setupHandlerDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close() //nolint:errcheck // Cleanup function, error not critical
	})

	_, err = db.Exec(`
		CREATE TABLE users (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL,
			display_name TEXT NOT NULL,
			password_hash TEXT NOT NULL DEFAULT '',
			user_source TEXT NOT NULL CHECK (user_source IN ('internal', 'entraid')),
			entra_id TEXT NOT NULL DEFAULT '',
			is_admin INTEGER NOT NULL DEFAULT 0,
			last_login TEXT NOT NULL DEFAULT '',
			access_token TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE UNIQUE INDEX idx_users_email ON users(email);
		CREATE INDEX idx_users_entra_id ON users(entra_id);
	`)
	require.NoError(t, err)
	return db
}

func seedUser(t *testing.T, db *sql.DB, email, displayName, source string, isAdmin bool) *Record {
	t.Helper()
	hash := ""
	if source == "internal" {
		var err error
		hash, err = HashPassword("TestPassword123!")
		require.NoError(t, err)
	}
	rec, err := CreateLocalUser(t.Context(), db, email, displayName, hash, isAdmin)
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

func TestListHandler(t *testing.T) {
	db := setupHandlerDB(t)
	seedUser(t, db, "alice@test.com", "Alice", "internal", false)
	seedUser(t, db, "bob@test.com", "Bob", "internal", true)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := ListHandler(db)(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Len(t, resp.Data, 2)
}

func TestGetHandler(t *testing.T) {
	db := setupHandlerDB(t)
	user := seedUser(t, db, "alice@test.com", "Alice", "internal", false)

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
	seedUser(t, db, "dup@test.com", "Existing", "internal", false)

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
	user := seedUser(t, db, "alice@test.com", "Alice", "internal", false)

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
	user := seedUser(t, db, "alice@test.com", "Alice", "internal", false)

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
	user := seedUser(t, db, "entra@test.com", "Entra User", "entraid", false)

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

func TestDeleteHandler(t *testing.T) {
	db := setupHandlerDB(t)
	user := seedUser(t, db, "alice@test.com", "Alice", "internal", false)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/"+user.ID, http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(user.ID)
	// Set a different user as the caller
	c.Set("user", &testUser{ID: "admin-user"})

	err := DeleteHandler(db)(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Verify deletion
	_, err = FindByID(t.Context(), db, user.ID)
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestDeleteHandlerPreventsSelfDeletion(t *testing.T) {
	db := setupHandlerDB(t)
	user := seedUser(t, db, "admin@test.com", "Admin", "internal", true)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/"+user.ID, http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(user.ID)
	// Set the same user as the caller
	c.Set("user", &testUser{ID: user.ID})

	err := DeleteHandler(db)(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteHandlerEntraIDUserRejected(t *testing.T) {
	db := setupHandlerDB(t)
	user := seedUser(t, db, "entra@test.com", "Entra User", "entraid", false)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/"+user.ID, http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(user.ID)
	c.Set("user", &testUser{ID: "admin-user"})

	err := DeleteHandler(db)(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteHandlerNotFound(t *testing.T) {
	db := setupHandlerDB(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/nonexistent", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("nonexistent")

	err := DeleteHandler(db)(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
