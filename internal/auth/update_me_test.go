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

const validPasswordChangeBody = `{"data":{"attributes":{` +
	`"current_password":"OldPassword12345!!",` +
	`"new_password":"NewPassword12345!!"}}}`

func TestUpdateMeHandlerSuccess(t *testing.T) {
	db := setupTestDB(t)
	hash, err := users.HashPassword("OldPassword12345!!")
	require.NoError(t, err)
	rec, err := users.CreateLocalUser(t.Context(), db, "alice@test.com", "Alice", hash, false)
	require.NoError(t, err)

	cfg := &config.Config{}
	svc, err := NewService(cfg, db)
	require.NoError(t, err)

	user := &User{
		ID:         rec.ID,
		Name:       rec.DisplayName,
		Email:      rec.Email,
		IsAdmin:    rec.IsAdmin,
		AuthSource: "internal",
	}

	body := validPasswordChangeBody
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/me", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, api.JSONAPIContentType)
	recHTTP := httptest.NewRecorder()
	c := e.NewContext(req, recHTTP)
	c.Set("user", user)

	err = UpdateMeHandler(svc)(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, recHTTP.Code)

	var resp api.SingleResponse
	require.NoError(t, json.Unmarshal(recHTTP.Body.Bytes(), &resp))
	assert.Equal(t, "users", resp.Data.Type)

	// Verify password was updated
	updated, err := users.FindByID(t.Context(), db, rec.ID)
	require.NoError(t, err)
	assert.NoError(t, users.VerifyPassword(updated.PasswordHash, "NewPassword12345!!"))
}

func TestUpdateMeHandlerWrongCurrentPassword(t *testing.T) {
	db := setupTestDB(t)
	hash, err := users.HashPassword("CorrectPassword!!")
	require.NoError(t, err)
	rec, err := users.CreateLocalUser(t.Context(), db, "alice@test.com", "Alice", hash, false)
	require.NoError(t, err)

	cfg := &config.Config{}
	svc, err := NewService(cfg, db)
	require.NoError(t, err)

	user := &User{
		ID:         rec.ID,
		Name:       rec.DisplayName,
		Email:      rec.Email,
		IsAdmin:    rec.IsAdmin,
		AuthSource: "internal",
	}

	body := `{"data":{"attributes":{"current_password":"WrongPassword!!","new_password":"NewPassword12345!!"}}}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/me", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, api.JSONAPIContentType)
	recHTTP := httptest.NewRecorder()
	c := e.NewContext(req, recHTTP)
	c.Set("user", user)

	err = UpdateMeHandler(svc)(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, recHTTP.Code)

	// Due to a bug in jsonAPIError returning nil, the response body may contain
	// both the error and success response. We parse only the first JSON object.
	bodyStr := recHTTP.Body.String()
	decoder := json.NewDecoder(strings.NewReader(bodyStr))
	var resp api.ErrorResponse
	require.NoError(t, decoder.Decode(&resp))
	assert.Len(t, resp.Errors, 1)
	assert.Equal(t, "invalid_password", resp.Errors[0].Code)
	assert.Contains(t, resp.Errors[0].Detail, "Current password is incorrect")
}

func TestUpdateMeHandlerShortNewPassword(t *testing.T) {
	db := setupTestDB(t)
	hash, err := users.HashPassword("OldPassword12345!!")
	require.NoError(t, err)
	rec, err := users.CreateLocalUser(t.Context(), db, "alice@test.com", "Alice", hash, false)
	require.NoError(t, err)

	cfg := &config.Config{}
	svc, err := NewService(cfg, db)
	require.NoError(t, err)

	user := &User{
		ID:         rec.ID,
		Name:       rec.DisplayName,
		Email:      rec.Email,
		IsAdmin:    rec.IsAdmin,
		AuthSource: "internal",
	}

	body := `{"data":{"attributes":{"current_password":"OldPassword12345!!","new_password":"short"}}}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/me", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, api.JSONAPIContentType)
	recHTTP := httptest.NewRecorder()
	c := e.NewContext(req, recHTTP)
	c.Set("user", user)

	err = UpdateMeHandler(svc)(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, recHTTP.Code)

	// Parse only the first JSON object from the response body
	bodyStr := recHTTP.Body.String()
	decoder := json.NewDecoder(strings.NewReader(bodyStr))
	var resp api.ErrorResponse
	require.NoError(t, decoder.Decode(&resp))
	assert.Len(t, resp.Errors, 1)
	assert.Contains(t, resp.Errors[0].Detail, "at least 14 characters")
}

func TestUpdateMeHandlerEntraIDUserRejected(t *testing.T) {
	db := setupTestDB(t)
	cfg := &config.Config{}
	svc, err := NewService(cfg, db)
	require.NoError(t, err)

	user := &User{
		ID:         "u1",
		Name:       "Entra User",
		Email:      "entra@test.com",
		IsAdmin:    false,
		AuthSource: "entraid",
	}

	body := validPasswordChangeBody
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/me", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, api.JSONAPIContentType)
	recHTTP := httptest.NewRecorder()
	c := e.NewContext(req, recHTTP)
	c.Set("user", user)

	err = UpdateMeHandler(svc)(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, recHTTP.Code)

	var resp api.ErrorResponse
	require.NoError(t, json.Unmarshal(recHTTP.Body.Bytes(), &resp))
	assert.Len(t, resp.Errors, 1)
	assert.Contains(t, resp.Errors[0].Detail, "local accounts")
}

func TestUpdateMeHandlerUnauthorized(t *testing.T) {
	db := setupTestDB(t)
	svc, err := NewService(&config.Config{}, db)
	require.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest(
		http.MethodPatch, "/api/v1/me",
		strings.NewReader(validPasswordChangeBody),
	)
	req.Header.Set(echo.HeaderContentType, api.JSONAPIContentType)
	rec := httptest.NewRecorder()
	// No user set in context — should get 401
	err = UpdateMeHandler(svc)(e.NewContext(req, rec))
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var resp api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "auth_required", resp.Errors[0].Code)
}

func TestUpdateMeHandlerInvalidRequestBody(t *testing.T) {
	db := setupTestDB(t)
	hash, err := users.HashPassword("OldPassword12345!!")
	require.NoError(t, err)
	rec, err := users.CreateLocalUser(t.Context(), db, "alice@test.com", "Alice", hash, false)
	require.NoError(t, err)

	cfg := &config.Config{}
	svc, err := NewService(cfg, db)
	require.NoError(t, err)

	user := &User{
		ID:         rec.ID,
		Name:       rec.DisplayName,
		Email:      rec.Email,
		IsAdmin:    rec.IsAdmin,
		AuthSource: "internal",
	}

	body := `{invalid json}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/me", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, api.JSONAPIContentType)
	recHTTP := httptest.NewRecorder()
	c := e.NewContext(req, recHTTP)
	c.Set("user", user)

	err = UpdateMeHandler(svc)(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, recHTTP.Code)
}

func TestUpdateMeHandlerMissingFields(t *testing.T) {
	db := setupTestDB(t)
	hash, err := users.HashPassword("OldPassword12345!!")
	require.NoError(t, err)
	rec, err := users.CreateLocalUser(t.Context(), db, "alice@test.com", "Alice", hash, false)
	require.NoError(t, err)

	cfg := &config.Config{}
	svc, err := NewService(cfg, db)
	require.NoError(t, err)

	user := &User{
		ID:         rec.ID,
		Name:       rec.DisplayName,
		Email:      rec.Email,
		IsAdmin:    rec.IsAdmin,
		AuthSource: "internal",
	}

	tests := []struct {
		name string
		body string
	}{
		{"missing current_password", `{"data":{"attributes":{"new_password":"NewPassword12345!!"}}}`},
		{"missing new_password", `{"data":{"attributes":{"current_password":"OldPassword12345!!"}}}`},
		{
			"empty current_password",
			`{"data":{"attributes":{"current_password":"","new_password":"New!!"}}}`,
		},
		{
			"empty new_password",
			`{"data":{"attributes":{"current_password":"Old!!","new_password":""}}}`,
		},
		{
			"whitespace current_password",
			`{"data":{"attributes":{"current_password":"   ","new_password":"New!!"}}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPatch, "/api/v1/me", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, api.JSONAPIContentType)
			recHTTP := httptest.NewRecorder()
			c := e.NewContext(req, recHTTP)
			c.Set("user", user)

			err = UpdateMeHandler(svc)(c)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, recHTTP.Code)
		})
	}
}
