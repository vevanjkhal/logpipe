package transform

import (
	"time"

	"github.com/logpipe/logpipe/internal/dedupe"
)

// dedupeTransformer wraps a Deduper and implements Transformer.
// Duplicate lines within the window are replaced with an empty string
// so the pipeline can discard them.
type dedupeTransformer struct {
	dd *dedupe.Deduper
}

// NewDedupe returns a Transformer that suppresses consecutive duplicate lines
// within the given time window. Pass 0 to use the default (5 s).
func NewDedupe(window time.Duration) Transformer {
	opts := []dedupe.Option{}
	if window > 0 {
		opts = append(opts, dedupe.WithWindow(window))
	}
	return &dedupeTransformer{dd: dedupe.New(opts...)}
}

// Transform returns an empty string when the line is a duplicate within the
// active window, causing the pipeline to drop it. Otherwise it returns the
// line unchanged.
func (t *dedupeTransformer) Transform(line string) string {
	if t.dd.IsDuplicate(line) {
		return ""
	}
	return line
}
