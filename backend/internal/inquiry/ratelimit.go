package inquiry

import (
	"sync"
	"time"
)

// RateLimiter is a simple in-memory fixed-window limiter keyed by client
// identifier (typically the client IP). It bounds abuse per API instance; a
// multi-instance deployment would need a shared store.
type RateLimiter struct {
	mu       sync.Mutex
	limit    int
	window   time.Duration
	now      func() time.Time
	counters map[string]counter
}

type counter struct {
	count       int
	windowStart time.Time
}

// NewRateLimiter creates a limiter allowing at most limit requests per window.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:    limit,
		window:   window,
		now:      time.Now,
		counters: make(map[string]counter),
	}
}

// Allow reports whether a request from key is permitted, and how long until the
// current window resets when it is not.
func (l *RateLimiter) Allow(key string) (bool, time.Duration) {
	if l == nil || l.limit <= 0 {
		return true, 0
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	entry, ok := l.counters[key]
	if !ok || now.Sub(entry.windowStart) >= l.window {
		l.counters[key] = counter{count: 1, windowStart: now}
		return true, 0
	}
	if entry.count >= l.limit {
		return false, l.window - now.Sub(entry.windowStart)
	}
	entry.count++
	l.counters[key] = entry
	return true, 0
}
