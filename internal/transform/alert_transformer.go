package transform

import (
	"github.com/logpipe/logpipe/internal/alert"
)

// alertTransformer watches lines and fires an alert when the count
// within the configured window exceeds the threshold. Lines are always
// passed through unchanged.
type alertTransformer struct {
	watcher *alert.Watcher
}

// NewAlert returns a Transformer that records every line with the given
// Watcher and passes the line through unmodified.
func NewAlert(w *alert.Watcher) Transformer {
	return &alertTransformer{watcher: w}
}

func (t *alertTransformer) Transform(line string) (string, bool) {
	t.watcher.Record()
	return line, true
}
