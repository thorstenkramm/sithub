package auth

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestPNG(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			img.Set(x, y, color.RGBA{R: 255, A: 255})
		}
	}
	var buf bytes.Buffer
	require.NoError(t, png.Encode(&buf, img))
	return buf.Bytes()
}

func TestServeAvatarFound(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	// Create a test avatar file
	avatarPath := filepath.Join(dir, "user-1.png")
	require.NoError(t, os.WriteFile(avatarPath, createTestPNG(t), 0o600))

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/avatars/user-1", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("user_id")
	c.SetParamValues("user-1")

	h := ServeAvatarHandler(dir)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "max-age=300", rec.Header().Get("Cache-Control"))
}

func TestServeAvatarNotFound(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/avatars/missing", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("user_id")
	c.SetParamValues("missing")

	h := ServeAvatarHandler(dir)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUploadAvatarSuccess(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	imgData := createTestPNG(t)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("avatar", "photo.png")
	require.NoError(t, err)
	_, err = part.Write(imgData)
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/me/avatar", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &User{ID: "user-1", Name: "Test"})

	h := UploadAvatarHandler(dir)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify file exists
	_, err = os.Stat(filepath.Join(dir, "user-1.png"))
	assert.NoError(t, err)
}

func TestUploadAvatarUnauthorized(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/me/avatar", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := UploadAvatarHandler(dir)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestDeleteAvatarSuccess(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	// Create avatar first
	avatarPath := filepath.Join(dir, "user-1.png")
	require.NoError(t, os.WriteFile(avatarPath, createTestPNG(t), 0o600))

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/me/avatar", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &User{ID: "user-1", Name: "Test"})

	h := DeleteAvatarHandler(dir)
	require.NoError(t, h(c))

	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Verify file removed
	_, err := os.Stat(avatarPath)
	assert.True(t, os.IsNotExist(err))
}

func TestSyncAvatarSuccess(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	imgData := createTestPNG(t)

	mockClient := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(imgData)),
			}, nil
		},
	}

	SyncAvatar(t.Context(), mockClient, "user-1", dir)

	_, err := os.Stat(filepath.Join(dir, "user-1.png"))
	assert.NoError(t, err)
}

func TestSyncAvatarNotFound(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	mockClient := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(bytes.NewReader(nil)),
			}, nil
		},
	}

	SyncAvatar(t.Context(), mockClient, "user-1", dir)

	_, err := os.Stat(filepath.Join(dir, "user-1.png"))
	assert.True(t, os.IsNotExist(err))
}

type mockHTTPClient struct {
	doFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.doFunc(req)
}
