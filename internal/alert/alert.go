// Package alert provides threshold-based alerting for log pipelines.
// When a counter exceeds a configured threshold within a window, an alert
// is fired to a registered handler function.
package alert

import (
	"fmt"
	"sync"
	"time"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// Alert carries information about a triggered threshold.
type Alert struct {
	Name      string
	Level     Level
	Count     int
	Threshold int
	FiredAt   time.Time
	Message   string
}

// Handler is called when an alert fires.
type Handler func(Alert)

// Option configures a Watcher.
type Option func(*Watcher)

// WithWindow sets the rolling window duration (default: 1 minute).
func WithWindow(d time.Duration) Option {
	return func(w *Watcher) { w.window = d }
}

// WithLevel sets the alert level (default: LevelWarn).
func WithLevel(l Level) Option {
	return func(w *Watcher) { w.level = l }
}

// Watcher counts events and fires an alert when the threshold is exceeded.
type Watcher struct {
	mu        sync.Mutex
	name      string
	threshold int
	window    time.Duration
	level     Level
	handler   Handler
	events    []time.Time
	now       func() time.Time
}

// New creates a Watcher that fires handler when count exceeds threshold.
func New(name string, threshold int, handler Handler, opts ...Option) (*Watcher, error) {
	if threshold < 1 {
		return nil, fmt.Errorf("alert: threshold must be >= 1, got %d", threshold)
	}
	if handler == nil {
		return nil, fmt.Errorf("alert: handler must not be nil")
	}
	w := &Watcher{
		name:      name,
		threshold: threshold,
		window:    time.Minute,
		level:     LevelWarn,
		handler:   handler,
		now:       time.Now,
	}
	for _, o := range opts {
		o(w)
	}
	return w, nil
}

// Record registers one event and fires the handler if the threshold is met.
func (w *Watcher) Record() {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := w.now()
	cutoff := now.Add(-w.window)

	// Evict events outside the window.
	valid := w.events[:0]
	for _, t := range w.events {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	w.events = append(valid, now)

	count := len(w.events)
	if count >= w.threshold {
		w.handler(Alert{
			Name:      w.name,
			Level:     w.level,
			Count:     count,
			Threshold: w.threshold,
			FiredAt:   now,
			Message:   fmt.Sprintf("%s: %d events in %s (threshold %d)", w.name, count, w.window, w.threshold),
		})
		w.events = w.events[:0] // reset after firing
	}
}

// Count returns the number of events currently within the window.
func (w *Watcher) Count() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	cutoff := w.now().Add(-w.window)
	n := 0
	for _, t := range w.events {
		if t.After(cutoff) {
			n++
		}
	}
	return n
}
