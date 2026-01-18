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

	db := setupTestDB(t)
	repo := NewRepository(db)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/areas", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := ListHandler(repo)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get(echo.HeaderContentType), api.JSONAPIContentType)

	var resp api.CollectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Empty(t, resp.Data)
}

func TestListHandlerReturnsAreas(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	_, err := db.Exec(
		`INSERT INTO areas (id, name, sort_order, created_at, updated_at)
		VALUES ('a1', 'Alpha', 0, '2026-01-18T00:00:00Z', '2026-01-18T00:00:00Z')`,
	)
	require.NoError(t, err)

	repo := NewRepository(db)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/areas", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := ListHandler(repo)
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
}
