package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/apikeys"
	"github.com/thorstenkramm/sithub/internal/auth"
	"github.com/thorstenkramm/sithub/internal/testutil/dbtest"
)

func TestLoadAPIKey(t *testing.T) {
	t.Parallel()
	db := dbtest.New(t)
	now := time.Date(2026, 10, 8, 9, 0, 0, 0, time.UTC)
	in := apikeys.CreateInput{Name: "Finance", Scope: apikeys.ScopeReadBookings, ExpiresInMonths: 3}
	k, valid, err := apikeys.Create(context.Background(), db, now, in)
	require.NoError(t, err)
	_, expired, err := apikeys.Create(context.Background(), db, now.AddDate(-1, 0, 0), in)
	require.NoError(t, err)
	revoked, rawRevoked, err := apikeys.Create(context.Background(), db, now, in)
	require.NoError(t, err)
	require.NoError(t, apikeys.Delete(context.Background(), db, revoked.ID))
	for _, tc := range []struct {
		name, path, header string
		status             int
		principal          bool
	}{
		{"no header", "/api/v1/me", "", 200, false},
		{"valid", "/api/v1/me", "Bearer " + valid, 200, true},
		{"lowercase", "/api/v1/me", "bearer " + valid, 200, true},
		{"empty", "/api/v1/me", "Bearer", 401, false},
		{"empty spaced", "/api/v1/me", "Bearer ", 401, false},
		{"malformed", "/api/v1/me", "Bearer not-a-uuid", 401, false},
		{"unknown", "/api/v1/me", "Bearer " + uuid.NewString(), 401, false},
		{"expired", "/api/v1/me", "Bearer " + expired, 401, false},
		{"revoked", "/api/v1/me", "Bearer " + rawRevoked, 401, false},
		{"basic", "/api/v1/me", "Basic dXNlcjpwYXNz", 200, false},
		{"static", "/admin/api-keys", "Bearer invalid", 200, false},
		{"double space", "/api/v1/me", "Bearer  " + valid, 401, false},
		{"tab", "/api/v1/me", "Bearer\t" + valid, 401, false},
		{"no separator", "/api/v1/me", "Bearer" + valid, 401, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tc.path, http.NoBody)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}
			c := echo.New().NewContext(req, rec)
			called := false
			h := loadAPIKey(db, func() time.Time { return now })(func(c echo.Context) error {
				called = true
				if tc.principal {
					assert.Equal(t, k.ID, auth.GetAPIKeyFromContext(c).ID)
					assert.Nil(t, auth.GetUserFromContext(c))
				} else {
					assert.Nil(t, auth.GetAPIKeyFromContext(c))
				}
				return c.NoContent(http.StatusOK)
			})
			require.NoError(t, h(c))
			assert.Equal(t, tc.status, rec.Code)
			assert.Equal(t, tc.status == 200, called)
			if tc.status == 401 {
				assert.Equal(t, api.JSONAPIContentType, rec.Header().Get("Content-Type"))
				assert.Equal(t, `Bearer error="invalid_token"`, rec.Header().Get("WWW-Authenticate"))
				assert.Contains(t, rec.Body.String(), `"errors"`)
				detail := "Invalid API key"
				if tc.name == "expired" {
					detail = "API key has expired"
				}
				assert.Contains(t, rec.Body.String(), detail)
			}
		})
	}
}

func TestLoadUserSkipsAPIKey(t *testing.T) {
	t.Parallel()
	c := echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/", http.NoBody), httptest.NewRecorder())
	c.Request().AddCookie(&http.Cookie{Name: "sithub_user", Value: "cookie"})
	auth.SetAPIKeyInContext(c, &auth.APIKeyPrincipal{ID: "key"})
	h := LoadUser(nil)(func(c echo.Context) error {
		assert.Nil(t, auth.GetUserFromContext(c))
		return nil
	})
	require.NoError(t, h(c), "a bearer principal skips cookie decoding entirely")
}

func TestBearerTokenAmbiguity(t *testing.T) {
	t.Parallel()
	for _, headers := range [][]string{
		{"Basic proxy", "Bearer token"}, {"Bearer token", "Basic proxy"},
		{"Bearer one", "Bearer two"}, {"Bearer\ttoken"}, {" Bearer token"},
		{"Bearertoken"}, {"bearer,token"},
	} {
		t.Run(strings.Join(headers, "|"), func(t *testing.T) {
			t.Parallel()
			present, raw := bearerToken(headers)
			assert.True(t, present, "an explicit bearer scheme must not fall back to cookies")
			assert.Empty(t, raw, "ambiguous credentials must be rejected")
		})
	}
}

func TestBearerTokenIgnoresOtherSchemes(t *testing.T) {
	t.Parallel()
	for _, headers := range [][]string{nil, {""}, {"Basic dXNlcjpwYXNz"}, {"Digest x"}, {"Bear token"}} {
		t.Run(strings.Join(headers, "|"), func(t *testing.T) {
			t.Parallel()
			present, _ := bearerToken(headers)
			assert.False(t, present, "non-bearer schemes continue to cookie authentication")
		})
	}
}

func TestLoadAPIKeyDatabaseFailure(t *testing.T) {
	t.Parallel()
	db := dbtest.New(t)
	require.NoError(t, db.Close())
	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+uuid.NewString())
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	require.NoError(t, LoadAPIKey(db)(func(echo.Context) error {
		t.Fatal("storage failures must fail closed")
		return nil
	})(c))
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Empty(t, rec.Header().Get("WWW-Authenticate"), "a storage failure is not an invalid credential")
	assert.NotContains(t, rec.Body.String(), "database is closed")
}
