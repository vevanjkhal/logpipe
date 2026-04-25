package transform

import (
	"fmt"

	"github.com/yourorg/logpipe/internal/sampling"
)

// samplingTransformer wraps a Sampler and implements Transformer.
// Lines that are not selected by the sampler are suppressed (empty string returned).
type samplingTransformer struct {
	sampler *sampling.Sampler
}

// NewSampling returns a Transformer that probabilistically drops log lines.
// rate must be in (0.0, 1.0]; a rate of 1.0 keeps all lines.
func NewSampling(rate float64, opts ...sampling.Option) (Transformer, error) {
	s, err := sampling.New(rate, opts...)
	if err != nil {
		return nil, fmt.Errorf("transform: %w", err)
	}
	return &samplingTransformer{sampler: s}, nil
}

// Transform returns the line unchanged if the sampler allows it,
// or an empty string to signal suppression to the pipeline.
func (t *samplingTransformer) Transform(line string) string {
	if t.sampler.Allow(line) {
		return line
	}
	return ""
}
