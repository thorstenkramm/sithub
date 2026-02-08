package middleware

import (
	"sync"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
)

// RateLimiter tracks request counts per key within a sliding window.
// RateLimiter is safe for concurrent use.
type RateLimiter struct {
	mu      sync.Mutex
	entries map[string]*rateLimitEntry
	limit   int
	window  time.Duration
}

type rateLimitEntry struct {
	count   int
	resetAt time.Time
}

// NewRateLimiter creates a rate limiter allowing limit requests per window per key.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		entries: make(map[string]*rateLimitEntry),
		limit:   limit,
		window:  window,
	}
	go rl.cleanup()
	return rl
}

// Allow checks whether a request from the given key is allowed.
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	entry, ok := rl.entries[key]
	if !ok || now.After(entry.resetAt) {
		rl.entries[key] = &rateLimitEntry{count: 1, resetAt: now.Add(rl.window)}
		return true
	}
	entry.count++
	return entry.count <= rl.limit
}

// RateLimit returns Echo middleware that limits requests by client IP.
func RateLimit(limiter *RateLimiter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := c.RealIP()
			if !limiter.Allow(key) {
				return api.WriteTooManyRequests(
					c, "Too many login attempts. Please try again later.",
				)
			}
			return next(c)
		}
	}
}

// cleanup periodically removes expired entries to prevent memory leaks.
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, entry := range rl.entries {
			if now.After(entry.resetAt) {
				delete(rl.entries, key)
			}
		}
		rl.mu.Unlock()
	}
}
