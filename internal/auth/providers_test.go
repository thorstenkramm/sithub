package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/config"
)

func TestProvidersHandlerEntraIDConfigured(t *testing.T) {
	cfg := &config.Config{EntraID: testEntraConfig()}
	svc := newTestService(t, cfg)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/providers", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, ProvidersHandler(svc)(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/vnd.api+json", rec.Header().Get(echo.HeaderContentType))

	var body struct {
		Data struct {
			Type       string `json:"type"`
			ID         string `json:"id"`
			Attributes struct {
				EntraID bool `json:"entraid"`
				Local   bool `json:"local"`
			} `json:"attributes"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "auth-providers", body.Data.Type)
	assert.Equal(t, "current", body.Data.ID)
	assert.True(t, body.Data.Attributes.EntraID)
	assert.True(t, body.Data.Attributes.Local)
}

func TestProvidersHandlerEntraIDNotConfigured(t *testing.T) {
	cfg := &config.Config{} // No EntraID config
	svc := newTestService(t, cfg)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/providers", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, ProvidersHandler(svc)(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	var body struct {
		Data struct {
			Attributes struct {
				EntraID bool `json:"entraid"`
				Local   bool `json:"local"`
			} `json:"attributes"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.False(t, body.Data.Attributes.EntraID)
	assert.True(t, body.Data.Attributes.Local)
}
