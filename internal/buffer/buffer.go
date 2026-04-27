// Package buffer provides a ring-buffer backed line store for logpipe.
// It retains the last N log lines in memory, useful for replaying recent
// history to newly connected outputs or for crash-time diagnostics.
package buffer

import "sync"

// RingBuffer stores the last N lines in a fixed-size circular buffer.
type RingBuffer struct {
	mu       sync.RWMutex
	lines    []string
	cap      int
	head     int // next write position
	count    int // number of valid entries
}

// New creates a RingBuffer that retains at most capacity lines.
// If capacity is less than 1 it is set to 1.
func New(capacity int) *RingBuffer {
	if capacity < 1 {
		capacity = 1
	}
	return &RingBuffer{
		lines: make([]string, capacity),
		cap:   capacity,
	}
}

// Write appends a line to the ring buffer, evicting the oldest entry when full.
func (r *RingBuffer) Write(line string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lines[r.head] = line
	r.head = (r.head + 1) % r.cap
	if r.count < r.cap {
		r.count++
	}
}

// Snapshot returns a copy of all buffered lines in insertion order (oldest first).
func (r *RingBuffer) Snapshot() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.count == 0 {
		return nil
	}
	out := make([]string, r.count)
	start := (r.head - r.count + r.cap) % r.cap
	for i := 0; i < r.count; i++ {
		out[i] = r.lines[(start+i)%r.cap]
	}
	return out
}

// Len returns the number of lines currently stored.
func (r *RingBuffer) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.count
}

// Cap returns the maximum number of lines the buffer can hold.
func (r *RingBuffer) Cap() int { return r.cap }

// Reset clears the buffer without releasing the underlying memory.
func (r *RingBuffer) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.head = 0
	r.count = 0
}
