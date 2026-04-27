package transform

import (
	"fmt"

	"github.com/logpipe/logpipe/internal/ratelimit"
)

// ThrottleTransformer wraps a rate limiter and drops lines that exceed the
// configured rate, passing all other lines through unchanged.
type ThrottleTransformer struct {
	limiter *ratelimit.RateLimiter
}

// NewThrottle returns a ThrottleTransformer that allows at most rate lines per
// second. If rate is less than 1 it is clamped to 1.
func NewThrottle(rate int) (*ThrottleTransformer, error) {
	if rate < 0 {
		return nil, fmt.Errorf("transform: throttle rate must be >= 0, got %d", rate)
	}
	rl, err := ratelimit.New(rate)
	if err != nil {
		return nil, fmt.Errorf("transform: %w", err)
	}
	return &ThrottleTransformer{limiter: rl}, nil
}

// Transform returns the line unchanged if the rate limiter allows it, or an
// empty string to signal that the line should be dropped.
func (t *ThrottleTransformer) Transform(line string) (string, error) {
	if !t.limiter.Allow() {
		return "", nil
	}
	return line, nil
}

// Dropped returns the number of lines dropped since the last Reset.
func (t *ThrottleTransformer) Dropped() int64 {
	return t.limiter.Dropped()
}

// Reset refills the token bucket and clears the dropped counter.
func (t *ThrottleTransformer) Reset() {
	t.limiter.Reset()
}
