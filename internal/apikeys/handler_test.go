package apikeys

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
	"github.com/thorstenkramm/sithub/internal/auth"
	"github.com/thorstenkramm/sithub/internal/testutil/dbtest"
)

func callHandler(t *testing.T, handler echo.HandlerFunc, body, id string, user *auth.User) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body)), rec)
	c.SetParamNames("id")
	c.SetParamValues(id)
	if user != nil {
		c.Set("user", user)
	}
	require.NoError(t, handler(c))
	return rec
}

func TestHandlersLifecycle(t *testing.T) {
	t.Parallel()
	db := dbtest.New(t)
	user := &auth.User{ID: "admin"}
	body := `{"data":{"type":"api-keys","attributes":{"name":"Export","scope":"read-bookings","expires_in_months":12}}}`
	rec := callHandler(t, CreateHandler(db), body, "", user)
	require.Equal(t, http.StatusCreated, rec.Code)
	assert.NotContains(t, rec.Body.String(), "key_hash")
	var response api.SingleResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
	attrs, ok := response.Data.Attributes.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "admin", attrs["created_by"])
	key, ok := attrs["key"].(string)
	require.True(t, ok)
	assert.Len(t, key, 36)
	assert.Equal(t, "no-store", rec.Header().Get("Cache-Control"))
	rec = callHandler(t, ListHandler(db), "", "", user)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotContains(t, rec.Body.String(), `"key"`)
	assert.NotContains(t, rec.Body.String(), "key_hash")
	assert.NotContains(t, rec.Body.String(), key)
	assert.Contains(t, rec.Body.String(), key[:10])
	rec = callHandler(t, DeleteHandler(db), "", response.Data.ID, user)
	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Empty(t, rec.Body.String())
	rec = callHandler(t, DeleteHandler(db), "", response.Data.ID, user)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCreateHandlerValidation(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ body, detail string }{
		{"{", "Invalid request body"},
		{`{"data":{"type":"wrong"}}`, "data.type must be api-keys"},
		{`{"data":{"type":"api-keys","attributes":{}}}`, "name is required"},
		{`{"data":{"type":"api-keys","attributes":{"name":"` + strings.Repeat("a", 101) + `"}}}`,
			"name must be at most 100 characters"},
		{`{"data":{"type":"api-keys","attributes":{"name":"Name","scope":"write"}}}`, "scope must be read-bookings"},
		{`{"data":{"type":"api-keys","attributes":{"name":"Name","scope":"read-bookings","expires_in_months":6}}}`,
			"expires_in_months must be one of 3, 12, 24"},
	} {
		t.Run(tc.detail, func(t *testing.T) {
			t.Parallel()
			rec := callHandler(t, CreateHandler(nil), tc.body, "", &auth.User{ID: "admin"})
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.detail)
		})
	}
	rec := callHandler(t, CreateHandler(nil), `{}`, "", nil)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestHandlersDatabaseFailure(t *testing.T) {
	t.Parallel()
	db := dbtest.New(t)
	require.NoError(t, db.Close())
	for _, handler := range []echo.HandlerFunc{ListHandler(db), DeleteHandler(db), CreateHandler(db)} {
		rec := callHandler(t, handler,
			`{"data":{"type":"api-keys","attributes":{"name":"Export","scope":"read-bookings","expires_in_months":3}}}`,
			"missing", &auth.User{ID: "admin"})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}
}
