package transform

import (
	"encoding/json"
	"time"
)

// TimestampTransformer injects or rewrites a timestamp field in JSON log lines.
// For plain-text lines it prepends the timestamp as a prefix.
type TimestampTransformer struct {
	field  string
	format string
	now    func() time.Time
}

// TimestampOption configures a TimestampTransformer.
type TimestampOption func(*TimestampTransformer)

// WithTimestampField sets the JSON key used for the timestamp (default: "ts").
func WithTimestampField(field string) TimestampOption {
	return func(t *TimestampTransformer) {
		if field != "" {
			t.field = field
		}
	}
}

// WithTimestampFormat sets the time format string (default: time.RFC3339).
func WithTimestampFormat(format string) TimestampOption {
	return func(t *TimestampTransformer) {
		if format != "" {
			t.format = format
		}
	}
}

// NewTimestamp returns a Transformer that stamps each log line with the current time.
func NewTimestamp(opts ...TimestampOption) Transformer {
	t := &TimestampTransformer{
		field:  "ts",
		format: time.RFC3339,
		now:    time.Now,
	}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Transform implements Transformer.
func (t *TimestampTransformer) Transform(line string) (string, bool) {
	ts := t.now().Format(t.format)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err == nil {
		obj[t.field] = ts
		b, err := json.Marshal(obj)
		if err != nil {
			return line, true
		}
		return string(b), true
	}

	// Plain-text fallback: prepend timestamp.
	return ts + " " + line, true
}
