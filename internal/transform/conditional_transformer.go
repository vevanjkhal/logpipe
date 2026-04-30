package transform

import (
	"encoding/json"
	"strings"
)

// ConditionalTransformer applies an inner Transformer only when a condition
// is satisfied. Lines that do not match the condition are passed through
// unchanged.
type ConditionalTransformer struct {
	field   string
	value   string
	inner   Transformer
	invert  bool
}

// ConditionalOption configures a ConditionalTransformer.
type ConditionalOption func(*ConditionalTransformer)

// WithInvert reverses the condition so that the inner transformer is applied
// when the field does NOT match the expected value.
func WithInvert() ConditionalOption {
	return func(c *ConditionalTransformer) {
		c.invert = true
	}
}

// NewConditional returns a Transformer that applies inner only when the JSON
// field `field` equals `value` (or, when inverted, when it does not).
// For plain-text lines the condition is evaluated against a simple substring
// match on `value`.
func NewConditional(field, value string, inner Transformer, opts ...ConditionalOption) Transformer {
	c := &ConditionalTransformer{
		field: field,
		value: value,
		inner: inner,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// Transform implements Transformer.
func (c *ConditionalTransformer) Transform(line string) (string, error) {
	matched := c.matches(line)
	if c.invert {
		matched = !matched
	}
	if matched {
		return c.inner.Transform(line)
	}
	return line, nil
}

func (c *ConditionalTransformer) matches(line string) bool {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err == nil {
		v, ok := obj[c.field]
		if !ok {
			return false
		}
		switch val := v.(type) {
		case string:
			return val == c.value
		default:
			b, _ := json.Marshal(v)
			return string(b) == c.value
		}
	}
	// plain-text fallback: substring match
	return strings.Contains(line, c.value)
}
