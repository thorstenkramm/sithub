package areas

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/spaces"
)

func TestListHandlerEmpty(t *testing.T) {
	t.Parallel()

	cfg := &spaces.Config{}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/areas", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := ListHandler(cfg)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get(echo.HeaderContentType), api.JSONAPIContentType)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Empty(t, resp.Data)
}

func TestListHandlerReturnsAreas(t *testing.T) {
	t.Parallel()

	cfg := &spaces.Config{
		Areas: []spaces.Area{
			{
				ID:          "a1",
				Name:        "Alpha",
				Description: "Main area",
				FloorPlan:   "floor_plans/alpha.svg",
			},
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/areas", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := ListHandler(cfg)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Data, 1)
	assert.Equal(t, "areas", resp.Data[0].Type)
	assert.Equal(t, "a1", resp.Data[0].ID)

	attrs, ok := resp.Data[0].Attributes.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "Alpha", attrs["name"])
	assert.Equal(t, "Main area", attrs["description"])
	assert.Equal(t, "floor_plans/alpha.svg", attrs["floor_plan"])
}
