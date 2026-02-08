package middleware

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
