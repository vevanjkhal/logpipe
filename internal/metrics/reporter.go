package metrics

import (
	"context"
	"io"
	"time"
)

// Reporter periodically writes counter snapshots to a writer.
type Reporter struct {
	counters  *Counters
	w         io.Writer
	interval  time.Duration
}

// NewReporter creates a Reporter that writes snapshots of c to w every interval.
func NewReporter(c *Counters, w io.Writer, interval time.Duration) *Reporter {
	if interval <= 0 {
		interval = 10 * time.Second
	}
	return &Reporter{
		counters: c,
		w:        w,
		interval: interval,
	}
}

// Run starts the reporting loop and blocks until ctx is cancelled.
func (r *Reporter) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Emit one final snapshot before exiting.
			r.counters.Snapshot().WriteTo(r.w) //nolint:errcheck
			return
		case <-ticker.C:
			r.counters.Snapshot().WriteTo(r.w) //nolint:errcheck
		}
	}
}
