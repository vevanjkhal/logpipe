package pipeline

import (
	"fmt"

	"github.com/logpipe/logpipe/internal/transform"
)

// WithThrottle adds a rate-limiting transformer to the pipeline that allows at
// most rate lines per second. Lines exceeding the limit are silently dropped.
// A rate of 0 is treated as "no limit" and the option becomes a no-op.
// A negative rate is treated as an error.
func (b *Builder) WithThrottle(rate int) *Builder {
	if b.err != nil {
		return b
	}
	if rate < 0 {
		b.err = fmt.Errorf("pipeline: WithThrottle: rate must be non-negative, got %d", rate)
		return b
	}
	if rate == 0 {
		return b
	}
	tr, err := transform.NewThrottle(rate)
	if err != nil {
		b.err = fmt.Errorf("pipeline: WithThrottle: %w", err)
		return b
	}
	b.transformers = append(b.transformers, tr)
	return b
}
