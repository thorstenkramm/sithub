package itemgroups

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

func testConfig() *spaces.Config {
	return &spaces.Config{
		Areas: []spaces.Area{
			{
				ID:   "area-1",
				Name: "Area One",
				ItemGroups: []spaces.ItemGroup{
					{
						ID:   "ig-1",
						Name: "Room 101",
						Items: []spaces.Item{
							{ID: "item-1", Name: "Desk 1"},
							{ID: "item-2", Name: "Desk 2"},
						},
					},
					{
						ID:          "ig-2",
						Name:        "Room 102",
						Description: "Corner office",
						Items: []spaces.Item{
							{ID: "item-3", Name: "Desk 3"},
						},
					},
				},
			},
			{
				ID:   "area-2",
				Name: "Area Two",
				ItemGroups: []spaces.ItemGroup{
					{
						ID:   "ig-3",
						Name: "Parking Level 1",
						Items: []spaces.Item{
							{ID: "item-4", Name: "Lot A"},
						},
					},
				},
			},
		},
	}
}

func TestListHandlerReturnsItemGroupsForArea(t *testing.T) {
	t.Parallel()

	cfg := testConfig()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/areas/area-1/item-groups", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("area_id")
	c.SetParamValues("area-1")

	h := ListHandler(cfg)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get(echo.HeaderContentType), api.JSONAPIContentType)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Data, 2)

	assert.Equal(t, "item-groups", resp.Data[0].Type)
	assert.Equal(t, "ig-1", resp.Data[0].ID)
	attrs0, ok := resp.Data[0].Attributes.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "Room 101", attrs0["name"])

	assert.Equal(t, "ig-2", resp.Data[1].ID)
	attrs1, ok := resp.Data[1].Attributes.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "Room 102", attrs1["name"])
	assert.Equal(t, "Corner office", attrs1["description"])
}

func TestListHandlerAreaNotFound(t *testing.T) {
	t.Parallel()

	cfg := testConfig()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/areas/unknown/item-groups", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("area_id")
	c.SetParamValues("unknown")

	h := ListHandler(cfg)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestListHandlerEmptyItemGroups(t *testing.T) {
	t.Parallel()

	cfg := &spaces.Config{
		Areas: []spaces.Area{
			{ID: "area-empty", Name: "Empty Area", ItemGroups: nil},
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/areas/area-empty/item-groups", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("area_id")
	c.SetParamValues("area-empty")

	h := ListHandler(cfg)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Empty(t, resp.Data)
}

func TestListHandlerDynamic(t *testing.T) {
	t.Parallel()

	cfg := testConfig()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/areas/area-2/item-groups", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("area_id")
	c.SetParamValues("area-2")

	h := ListHandlerDynamic(func() *spaces.Config { return cfg })
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Data, 1)
	assert.Equal(t, "ig-3", resp.Data[0].ID)
}
