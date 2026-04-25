// Package sampling provides probabilistic log line sampling.
// It allows passing only a fraction of log lines through the pipeline,
// which is useful for high-volume sources where full fidelity is not required.
package sampling

import (
	"errors"
	"math/rand"
	"sync"
)

// Sampler decides whether a given log line should be passed through.
type Sampler struct {
	mu      sync.Mutex
	rate    float64 // fraction in (0.0, 1.0]
	rng     *rand.Rand
	dropped uint64
	passed  uint64
}

// Option configures a Sampler.
type Option func(*Sampler)

// WithSeed sets the random seed for reproducible sampling.
func WithSeed(seed int64) Option {
	return func(s *Sampler) {
		s.rng = rand.New(rand.NewSource(seed)) //nolint:gosec
	}
}

// New creates a Sampler that passes approximately rate*100% of lines.
// rate must be in the range (0.0, 1.0]; a rate of 1.0 passes everything.
func New(rate float64, opts ...Option) (*Sampler, error) {
	if rate <= 0 || rate > 1 {
		return nil, errors.New("sampling: rate must be in the range (0, 1]")
	}
	s := &Sampler{
		rate: rate,
		rng:  rand.New(rand.NewSource(rand.Int63())), //nolint:gosec
	}
	for _, o := range opts {
		o(s)
	}
	return s, nil
}

// Allow returns true if the line should pass through based on the configured rate.
func (s *Sampler) Allow(_ string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.rng.Float64() < s.rate {
		s.passed++
		return true
	}
	s.dropped++
	return false
}

// Snapshot returns the current passed and dropped counts.
func (s *Sampler) Snapshot() (passed, dropped uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.passed, s.dropped
}

// Reset zeroes the counters.
func (s *Sampler) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.passed = 0
	s.dropped = 0
}
