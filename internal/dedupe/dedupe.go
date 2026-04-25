// Package dedupe provides a transformer that suppresses consecutive duplicate log lines.
package dedupe

import (
	"sync"
	"time"
)

// Deduper suppresses repeated identical log lines within a configurable window.
type Deduper struct {
	mu       sync.Mutex
	last     string
	count    int
	window   time.Duration
	expireAt time.Time
	now      func() time.Time
}

// Option configures a Deduper.
type Option func(*Deduper)

// WithWindow sets the time window during which duplicates are suppressed.
func WithWindow(d time.Duration) Option {
	return func(dp *Deduper) { dp.window = d }
}

// New returns a new Deduper. Duplicates within window are suppressed.
// The default window is 5 seconds.
func New(opts ...Option) *Deduper {
	dp := &Deduper{
		window: 5 * time.Second,
		now:    time.Now,
	}
	for _, o := range opts {
		o(dp)
	}
	return dp
}

// IsDuplicate reports whether line is a duplicate of the previous line
// within the active window. It updates internal state accordingly.
func (d *Deduper) IsDuplicate(line string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()

	// Window expired — reset state.
	if now.After(d.expireAt) {
		d.last = ""
		d.count = 0
	}

	if line == d.last && !now.After(d.expireAt) {
		d.count++
		return true
	}

	d.last = line
	d.count = 1
	d.expireAt = now.Add(d.window)
	return false
}

// Suppressed returns how many lines have been suppressed since the last
// unique line was seen.
func (d *Deduper) Suppressed() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.count <= 1 {
		return 0
	}
	return d.count - 1
}

// Reset clears all deduplication state.
func (d *Deduper) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.last = ""
	d.count = 0
	d.expireAt = time.Time{}
}
