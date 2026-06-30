package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/auth"
	"github.com/thorstenkramm/sithub/internal/config"
)

func TestRateLimiterAllowsUnderLimit(t *testing.T) {
	rl := NewRateLimiter(3, time.Minute)
	for i := range 3 {
		assert.True(t, rl.Allow("ip1"), "request %d should be allowed", i+1)
	}
}

func TestRateLimiterBlocksOverLimit(t *testing.T) {
	rl := NewRateLimiter(2, time.Minute)
	require.True(t, rl.Allow("ip1"))
	require.True(t, rl.Allow("ip1"))
	assert.False(t, rl.Allow("ip1"), "third request should be blocked")
}

func TestRateLimiterIsolatesKeys(t *testing.T) {
	rl := NewRateLimiter(1, time.Minute)
	require.True(t, rl.Allow("ip1"))
	assert.True(t, rl.Allow("ip2"), "different key should be independent")
}

func TestRateLimiterResetsAfterWindow(t *testing.T) {
	rl := NewRateLimiter(1, 50*time.Millisecond)
	require.True(t, rl.Allow("ip1"))
	assert.False(t, rl.Allow("ip1"))

	time.Sleep(60 * time.Millisecond)
	assert.True(t, rl.Allow("ip1"), "should be allowed after window reset")
}

// TestRateLimitMiddlewareReturns429 exercises the production configuration
// (60 requests/minute, matching the /api/v1/auth/login limiter): the 61st
// request from the same client IP must receive HTTP 429.
func TestRateLimitMiddlewareReturns429(t *testing.T) {
	limiter := NewRateLimiter(60, time.Minute)
	e := echo.New()
	svc, err := auth.NewService(&config.Config{}, nil)
	require.NoError(t, err)
	e.POST("/api/v1/auth/login", auth.LocalLoginHandler(svc), RateLimit(limiter))

	const clientIP = "203.0.113.7:12345"
	doRequest := func() int {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(`{`))
		req.RemoteAddr = clientIP
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		return rec.Code
	}

	for i := 1; i <= 60; i++ {
		require.Equal(t, http.StatusBadRequest, doRequest(), "request %d within limit should reach login handler", i)
	}
	assert.Equal(t, http.StatusTooManyRequests, doRequest(), "61st request must be rate limited")
}
