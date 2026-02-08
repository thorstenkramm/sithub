package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/config"
	"github.com/thorstenkramm/sithub/internal/users"
)

func TestLocalLoginHandlerSuccess(t *testing.T) {
	db := setupTestDB(t)
	hash, err := users.HashPassword("TestPassword123!")
	require.NoError(t, err)
	_, err = users.CreateLocalUser(t.Context(), db, "alice@test.com", "Alice", hash, false)
	require.NoError(t, err)

	cfg := &config.Config{}
	svc, err := NewService(cfg, db)
	require.NoError(t, err)

	body := `{"email":"alice@test.com","password":"TestPassword123!"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = LocalLoginHandler(svc)(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp api.SingleResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "users", resp.Data.Type)
	attrs, ok := resp.Data.Attributes.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "Alice", attrs["display_name"])
	assert.Equal(t, "alice@test.com", attrs["email"])
	assert.Equal(t, false, attrs["is_admin"])
	assert.Equal(t, "internal", attrs["auth_source"])
	assert.Equal(t, "user", attrs["role"])

	// Check cookie was set
	cookies := rec.Result().Cookies()
	assert.NotEmpty(t, cookies)
	var found bool
	for _, cookie := range cookies {
		if cookie.Name == userCookieName {
			found = true
			assert.NotEmpty(t, cookie.Value)
		}
	}
	assert.True(t, found, "Expected user cookie to be set")
}

func TestLocalLoginHandlerAdminRole(t *testing.T) {
	db := setupTestDB(t)
	hash, err := users.HashPassword("AdminPass123!!")
	require.NoError(t, err)
	_, err = users.CreateLocalUser(t.Context(), db, "admin@test.com", "Admin", hash, true)
	require.NoError(t, err)

	cfg := &config.Config{}
	svc, err := NewService(cfg, db)
	require.NoError(t, err)

	body := `{"email":"admin@test.com","password":"AdminPass123!!"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = LocalLoginHandler(svc)(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp api.SingleResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	attrs, ok := resp.Data.Attributes.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, true, attrs["is_admin"])
	assert.Equal(t, "admin", attrs["role"])
}

func TestLocalLoginHandlerInvalidPassword(t *testing.T) {
	db := setupTestDB(t)
	hash, err := users.HashPassword("CorrectPassword123!")
	require.NoError(t, err)
	_, err = users.CreateLocalUser(t.Context(), db, "alice@test.com", "Alice", hash, false)
	require.NoError(t, err)

	cfg := &config.Config{}
	svc, err := NewService(cfg, db)
	require.NoError(t, err)

	body := `{"email":"alice@test.com","password":"WrongPassword123!"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = LocalLoginHandler(svc)(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var resp api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Len(t, resp.Errors, 1)
	assert.Equal(t, "invalid_credentials", resp.Errors[0].Code)
}

func TestLocalLoginHandlerUserNotFound(t *testing.T) {
	db := setupTestDB(t)
	cfg := &config.Config{}
	svc, err := NewService(cfg, db)
	require.NoError(t, err)

	body := `{"email":"nonexistent@test.com","password":"TestPassword123!"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = LocalLoginHandler(svc)(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var resp api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Len(t, resp.Errors, 1)
	assert.Equal(t, "invalid_credentials", resp.Errors[0].Code)
}

func TestLocalLoginHandlerEntraIDUserRejected(t *testing.T) {
	db := setupTestDB(t)
	// Create an Entra ID user
	_, err := db.Exec(`
		INSERT INTO users (id, email, display_name, password_hash,
		  user_source, entra_id, is_admin, created_at, updated_at)
		VALUES ('u1', 'entra@test.com', 'Entra User', '',
		  'entraid', 'entra-u1', 0, datetime('now'), datetime('now'))
	`)
	require.NoError(t, err)

	cfg := &config.Config{}
	svc, err := NewService(cfg, db)
	require.NoError(t, err)

	body := `{"email":"entra@test.com","password":"SomePassword123!"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = LocalLoginHandler(svc)(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var resp api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Len(t, resp.Errors, 1)
	assert.Equal(t, "wrong_auth_source", resp.Errors[0].Code)
	assert.Contains(t, resp.Errors[0].Detail, "Entra ID")
}

func TestLocalLoginHandlerInvalidRequestBody(t *testing.T) {
	db := setupTestDB(t)
	svc, err := NewService(&config.Config{}, db)
	require.NoError(t, err)

	for _, body := range []string{`{invalid json}`, `not json at all`} {
		e := echo.New()
		req := httptest.NewRequest(
			http.MethodPost, "/api/v1/auth/login",
			strings.NewReader(body),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		err = LocalLoginHandler(svc)(e.NewContext(req, rec))
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestLocalLoginHandlerMissingFields(t *testing.T) {
	db := setupTestDB(t)
	cfg := &config.Config{}
	svc, err := NewService(cfg, db)
	require.NoError(t, err)

	tests := []struct {
		name string
		body string
	}{
		{"missing email", `{"password":"TestPassword123!"}`},
		{"missing password", `{"email":"alice@test.com"}`},
		{"empty email", `{"email":"","password":"TestPassword123!"}`},
		{"empty password", `{"email":"alice@test.com","password":""}`},
		{"whitespace email", `{"email":"   ","password":"TestPassword123!"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err = LocalLoginHandler(svc)(c)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	}
}
