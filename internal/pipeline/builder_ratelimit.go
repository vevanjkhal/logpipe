package pipeline

import "github.com/yourorg/logpipe/internal/ratelimit"

// WithRateLimit adds a token-bucket rate limiter to the pipeline builder.
// Only ratePerSec log lines per second will be forwarded; excess lines are
// silently dropped. Calling this more than once replaces the previous limiter.
func (b *Builder) WithRateLimit(ratePerSec float64) *Builder {
	b.limiter = ratelimit.New(ratePerSec)
	return b
}
