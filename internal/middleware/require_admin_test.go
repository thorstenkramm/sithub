package middleware

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/auth"
)

func TestRequireAdminAllowsAdmin(t *testing.T) {
	user := &auth.User{
		ID:          "admin1",
		Name:        "Admin User",
		Email:       "admin@test.com",
		IsAdmin:     true,
		IsPermitted: true,
	}

	rec := runMiddleware(t, RequireAdmin(), "/api/v1/admin", user)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRequireAdminBlocksNonAdmin(t *testing.T) {
	user := &auth.User{
		ID:          "user1",
		Name:        "Regular User",
		Email:       "user@test.com",
		IsAdmin:     false,
		IsPermitted: true,
	}

	rec := runMiddleware(t, RequireAdmin(), "/api/v1/admin", user)
	assert.Equal(t, http.StatusForbidden, rec.Code)

	var resp api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Len(t, resp.Errors, 1)
	assert.Equal(t, "forbidden", resp.Errors[0].Code)
}

func TestRequireAdminBlocksMissingUser(t *testing.T) {
	rec := runMiddleware(t, RequireAdmin(), "/api/v1/admin", nil)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var resp api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Len(t, resp.Errors, 1)
	assert.Equal(t, "auth_required", resp.Errors[0].Code)
}
