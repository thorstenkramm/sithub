package middleware

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/auth"
)

func TestRequireDataExporter(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		user     *auth.User
		wantCode int
		wantErr  string
	}{
		{
			name: "data exporter is allowed", wantCode: http.StatusOK,
			user: &auth.User{ID: "u1", IsDataExporter: true, IsPermitted: true},
		},
		{
			name: "admin without the flag is forbidden", wantCode: http.StatusForbidden, wantErr: "forbidden",
			user: &auth.User{ID: "u2", IsAdmin: true, IsPermitted: true},
		},
		{
			name: "regular user is forbidden", wantCode: http.StatusForbidden, wantErr: "forbidden",
			user: &auth.User{ID: "u3", IsPermitted: true},
		},
		{name: "missing user is unauthorized", wantCode: http.StatusUnauthorized, wantErr: "auth_required"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := runMiddleware(t, RequireDataExporter(), "/api/v1/bookings/export", tc.user)
			assert.Equal(t, tc.wantCode, rec.Code)
			if tc.wantErr == "" {
				return
			}
			var resp api.ErrorResponse
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			require.Len(t, resp.Errors, 1)
			assert.Equal(t, tc.wantErr, resp.Errors[0].Code)
		})
	}
}

// TestRequireAuthRefreshesDataExporterFlag proves the flag comes from the database on every request:
// the cookie user carries a stale value in both directions.
func TestRequireAuthRefreshesDataExporterFlag(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		dbFlag   int
		cookie   bool
		wantCode int
	}{
		{name: "grant takes effect without re-login", dbFlag: 1, cookie: false, wantCode: http.StatusOK},
		{name: "revoke takes effect without re-login", dbFlag: 0, cookie: true, wantCode: http.StatusForbidden},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := setupMiddlewareTestDB(t)
			seedMiddlewareUser(t, db, "u1", "u1@example.com")
			setDataExporter(t, db, "u1", tc.dbFlag)
			svc := newTestAuthServiceWithStore(t, db)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/bookings/export", http.NoBody)
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)
			c.Set("user", &auth.User{ID: "u1", AuthSource: "internal", IsPermitted: true, IsDataExporter: tc.cookie})

			h := RequireAuth(svc)(RequireDataExporter()(func(c echo.Context) error {
				return c.NoContent(http.StatusOK)
			}))
			require.NoError(t, h(c))
			assert.Equal(t, tc.wantCode, rec.Code)
		})
	}
}

func setDataExporter(t *testing.T, db *sql.DB, id string, flag int) {
	t.Helper()
	_, err := db.Exec(`UPDATE users SET is_data_exporter = ? WHERE id = ?`, flag, id)
	require.NoError(t, err)
}
