package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogoutHandlerClearsCookie(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := LogoutHandler()(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	cookies := rec.Result().Cookies()
	assert.NotEmpty(t, cookies)

	var found bool
	for _, cookie := range cookies {
		if cookie.Name != userCookieName {
			continue
		}
		found = true
		assert.Empty(t, cookie.Value)
		assert.Equal(t, -1, cookie.MaxAge)
		assert.Equal(t, "/", cookie.Path)
		assert.True(t, cookie.HttpOnly)
		assert.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
	}
	assert.True(t, found, "Expected user cookie to be cleared")
}

func TestLogoutHandlerWithHTTPS(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "https://example.com/api/v1/auth/logout", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := LogoutHandler()(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	cookies := rec.Result().Cookies()
	var found bool
	for _, cookie := range cookies {
		if cookie.Name == userCookieName {
			found = true
			assert.True(t, cookie.Secure, "Expected Secure flag for HTTPS")
		}
	}
	assert.True(t, found)
}
