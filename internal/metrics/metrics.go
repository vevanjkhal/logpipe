// Package metrics provides lightweight in-process counters for tracking
// pipeline activity such as lines read, filtered, and written.
package metrics

import (
	"fmt"
	"io"
	"sync/atomic"
)

// Counters holds atomic counters for pipeline statistics.
type Counters struct {
	LinesRead     atomic.Int64
	LinesFiltered atomic.Int64
	LinesWritten  atomic.Int64
	Errors        atomic.Int64
}

// IncRead increments the lines-read counter.
func (c *Counters) IncRead() { c.LinesRead.Add(1) }

// IncFiltered increments the lines-filtered counter.
func (c *Counters) IncFiltered() { c.LinesFiltered.Add(1) }

// IncWritten increments the lines-written counter.
func (c *Counters) IncWritten() { c.LinesWritten.Add(1) }

// IncError increments the error counter.
func (c *Counters) IncError() { c.Errors.Add(1) }

// Snapshot returns a point-in-time copy of the counters.
func (c *Counters) Snapshot() Snapshot {
	return Snapshot{
		LinesRead:     c.LinesRead.Load(),
		LinesFiltered: c.LinesFiltered.Load(),
		LinesWritten:  c.LinesWritten.Load(),
		Errors:        c.Errors.Load(),
	}
}

// Snapshot is an immutable point-in-time view of Counters.
type Snapshot struct {
	LinesRead     int64
	LinesFiltered int64
	LinesWritten  int64
	Errors        int64
}

// WriteTo writes a human-readable summary to w.
func (s Snapshot) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprintf(w,
		"lines_read=%d lines_filtered=%d lines_written=%d errors=%d\n",
		s.LinesRead, s.LinesFiltered, s.LinesWritten, s.Errors,
	)
	return int64(n), err
}
