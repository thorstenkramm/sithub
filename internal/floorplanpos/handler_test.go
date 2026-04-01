package floorplanpos

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/api"
)

func TestListHandlerRequiresFloorPlan(t *testing.T) {
	t.Parallel()

	store := setupTestDB(t)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/floor-plan-positions", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, ListHandler(store)(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateHandlerReturnsResource(t *testing.T) {
	t.Parallel()

	store := setupTestDB(t)
	e := echo.New()
	body := `{"data":{"type":"floor-plan-positions","attributes":{` +
		`"floor_plan":"office.svg","item_id":"desk-1","label":"D1",` +
		`"x":10,"y":20,"width":5,"height":3,"border_width":4}}}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/floor-plan-positions", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, CreateHandler(store)(c))
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Header().Get(echo.HeaderContentType), api.JSONAPIContentType)

	var resp api.SingleResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "floor-plan-positions", resp.Data.Type)

	attrs, ok := resp.Data.Attributes.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "office.svg", attrs["floor_plan"])
	assert.Equal(t, "desk-1", attrs["item_id"])
	assert.Equal(t, float64(4), attrs["border_width"])
}

func TestUpdateHandlerNotFound(t *testing.T) {
	t.Parallel()

	store := setupTestDB(t)
	e := echo.New()
	body := `{"data":{"type":"floor-plan-positions","attributes":{"x":25}}}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/floor-plan-positions/missing", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("missing")

	require.NoError(t, UpdateHandler(store)(c))
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteHandlerDeletesExistingPosition(t *testing.T) {
	t.Parallel()

	store := setupTestDB(t)
	ctx := t.Context()
	pos, err := Create(ctx, store, &CreateInput{
		FloorPlan: "office.svg",
		ItemID:    "desk-1",
		X:         10,
		Y:         20,
		Width:     5,
		Height:    3,
	})
	require.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/floor-plan-positions/"+pos.ID, http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(pos.ID)

	require.NoError(t, DeleteHandler(store)(c))
	assert.Equal(t, http.StatusNoContent, rec.Code)

	positions, err := FindByFloorPlan(ctx, store, "office.svg")
	require.NoError(t, err)
	assert.Empty(t, positions)
}
