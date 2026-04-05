package system

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/api"
)

func TestSettingsHandlerReturnsConfig(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/settings", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := SettingsHandler(7)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, api.JSONAPIContentType, rec.Header().Get(echo.HeaderContentType))

	var resp api.SingleResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "settings", resp.Data.Type)

	attrs, ok := resp.Data.Attributes.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, float64(7), attrs["weeks_in_advanced"])
}
