package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAPIKeyPrincipalContext(t *testing.T) {
	t.Parallel()
	c := echo.New().NewContext(httptest.NewRequest("GET", "/", http.NoBody), httptest.NewRecorder())
	assert.Nil(t, GetAPIKeyFromContext(c))
	p := &APIKeyPrincipal{ID: "key-id", Name: "Finance", Scope: "read-bookings"}
	SetAPIKeyInContext(c, p)
	assert.Same(t, p, GetAPIKeyFromContext(c))
	assert.Nil(t, GetUserFromContext(c), "a key must never become a user")
	assert.True(t, p.HasScope("read-bookings"))
	assert.False(t, p.HasScope("READ-BOOKINGS"))
	assert.False(t, p.HasScope(""))
}
