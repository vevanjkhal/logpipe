package pipeline

import (
	"time"

	"github.com/logpipe/logpipe/internal/alert"
	"github.com/logpipe/logpipe/internal/transform"
)

// WithAlert adds a threshold-based alert to the pipeline. handler is called
// when count events are observed within window. If window is zero, the
// default (1 minute) is used.
func (b *Builder) WithAlert(
	name string,
	threshold int,
	window time.Duration,
	level alert.Level,
	handler alert.Handler,
) *Builder {
	opts := []alert.Option{alert.WithLevel(level)}
	if window > 0 {
		opts = append(opts, alert.WithWindow(window))
	}
	w, err := alert.New(name, threshold, handler, opts...)
	if err != nil {
		panic("logpipe/pipeline: WithAlert: " + err.Error())
	}
	b.transformers = append(b.transformers, transform.NewAlert(w))
	return b
}
