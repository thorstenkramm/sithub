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
)

func TestListHandlerEmpty(t *testing.T) {
	t.Parallel()

	cfg := &Config{}

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

	cfg := &Config{
		Areas: []Area{
			{
				ID:          "a1",
				Name:        "Alpha",
				Description: "Main area",
				FloorPlan:   "alpha.svg",
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
	assert.Equal(t, "alpha.svg", attrs["floor_plan"])
}

func TestListHandlerReturnsIcon(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Areas: []Area{
			{ID: "a1", Name: "Alpha", Icon: "mdi-garage"},
			{ID: "a2", Name: "Beta"},
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/areas", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := ListHandler(cfg)
	require.NoError(t, h(c))

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Data, 2)

	attrs0, _ := resp.Data[0].Attributes.(map[string]interface{})
	assert.Equal(t, "mdi-garage", attrs0["icon"])

	attrs1, _ := resp.Data[1].Attributes.(map[string]interface{})
	assert.Nil(t, attrs1["icon"], "icon should be absent when not set")
}
