package rooms

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

func TestListHandlerNotFound(t *testing.T) {
	t.Parallel()

	cfg := &spaces.Config{}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/areas/missing/rooms", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("area_id")
	c.SetParamValues("missing")

	h := ListHandler(cfg)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var resp api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Errors, 1)
	assert.Equal(t, "not_found", resp.Errors[0].Code)
}

func TestListHandlerReturnsRooms(t *testing.T) {
	t.Parallel()

	cfg := &spaces.Config{
		Areas: []spaces.Area{
			{
				ID:   "area-1",
				Name: "Office",
				Rooms: []spaces.Room{
					{
						ID:          "room-1",
						Name:        "Room 1",
						Description: "Main room",
					},
				},
			},
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/areas/area-1/rooms", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("area_id")
	c.SetParamValues("area-1")

	h := ListHandler(cfg)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Data, 1)
	assert.Equal(t, "rooms", resp.Data[0].Type)
	assert.Equal(t, "room-1", resp.Data[0].ID)

	attrs, ok := resp.Data[0].Attributes.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "Room 1", attrs["name"])
	assert.Equal(t, "Main room", attrs["description"])
}
