// Package ratelimit provides a token-bucket rate limiter for log lines.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter controls how many log lines per second are allowed through the
// pipeline. It uses a simple token-bucket algorithm.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens per nanosecond
	lastTick time.Time
	dropped  int64
}

// New creates a Limiter that allows up to ratePerSec log lines per second.
// A burst of up to ratePerSec lines is permitted before throttling begins.
func New(ratePerSec float64) *Limiter {
	if ratePerSec <= 0 {
		ratePerSec = 1
	}
	return &Limiter{
		tokens:   ratePerSec,
		max:      ratePerSec,
		rate:     ratePerSec / float64(time.Second),
		lastTick: time.Now(),
	}
}

// Allow returns true if a log line should be forwarded, false if it should
// be dropped due to rate limiting.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(l.lastTick)
	l.lastTick = now

	l.tokens += float64(elapsed) * l.rate
	if l.tokens > l.max {
		l.tokens = l.max
	}

	if l.tokens >= 1 {
		l.tokens--
		return true
	}

	l.dropped++
	return false
}

// Dropped returns the total number of lines dropped since the Limiter was
// created.
func (l *Limiter) Dropped() int64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.dropped
}

// Reset resets the dropped counter and refills the token bucket to maximum
// capacity.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.tokens = l.max
	l.dropped = 0
	l.lastTick = time.Now()
}
