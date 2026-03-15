package areas

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFloorPlanHandlerServesImage(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	content := []byte("PNG image data")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "plan.png"), content, 0o600))

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/floor-plans/plan.png", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("filename")
	c.SetParamValues("plan.png")

	h := FloorPlanHandler(dir)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "image/png", rec.Header().Get(echo.HeaderContentType))
	assert.Equal(t, content, rec.Body.Bytes())
}

func TestFloorPlanHandlerServesSVG(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	content := []byte("<svg></svg>")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "plan.svg"), content, 0o600))

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/floor-plans/plan.svg", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("filename")
	c.SetParamValues("plan.svg")

	h := FloorPlanHandler(dir)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "image/svg+xml", rec.Header().Get(echo.HeaderContentType))
}

func TestFloorPlanHandlerNotFound(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/floor-plans/missing.png", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("filename")
	c.SetParamValues("missing.png")

	h := FloorPlanHandler(dir)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestFloorPlanHandlerNotConfigured(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/floor-plans/plan.png", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("filename")
	c.SetParamValues("plan.png")

	h := FloorPlanHandler("")
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestFloorPlanHandlerUnsupportedFormat(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "plan.gif"), []byte("GIF"), 0o600))

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/floor-plans/plan.gif", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("filename")
	c.SetParamValues("plan.gif")

	h := FloorPlanHandler(dir)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestFloorPlanHandlerPathTraversal(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/floor-plans/..%2Fetc%2Fpasswd", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("filename")
	c.SetParamValues("../etc/passwd")

	h := FloorPlanHandler(dir)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
