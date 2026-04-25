// Package transform provides log line transformation capabilities,
// such as adding fields, redacting sensitive data, or reformatting values.
package transform

import (
	"encoding/json"
	"strings"
)

// Transformer mutates a log line string and returns the result.
type Transformer interface {
	Transform(line string) string
}

// Chain applies a sequence of Transformers in order.
type Chain []Transformer

// Transform applies each transformer in the chain sequentially.
func (c Chain) Transform(line string) string {
	for _, t := range c {
		line = t.Transform(line)
	}
	return line
}

// RedactTransformer replaces occurrences of sensitive keys in JSON log lines.
type RedactTransformer struct {
	Keys    []string
	MaskVal string
}

// NewRedact creates a RedactTransformer that masks the given JSON keys.
func NewRedact(keys []string, maskVal string) *RedactTransformer {
	if maskVal == "" {
		maskVal = "***"
	}
	return &RedactTransformer{Keys: keys, MaskVal: maskVal}
}

// Transform redacts sensitive fields from a JSON log line.
// Non-JSON lines are returned unchanged.
func (r *RedactTransformer) Transform(line string) string {
	var obj map[string]any
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for _, k := range r.Keys {
		if _, ok := obj[k]; ok {
			obj[k] = r.MaskVal
		}
	}
	b, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(b)
}

// AddFieldTransformer injects a static key=value pair into JSON log lines.
type AddFieldTransformer struct {
	Key   string
	Value string
}

// NewAddField creates an AddFieldTransformer.
func NewAddField(key, value string) *AddFieldTransformer {
	return &AddFieldTransformer{Key: key, Value: value}
}

// Transform adds the configured field to a JSON log line.
// Non-JSON lines are returned unchanged.
func (a *AddFieldTransformer) Transform(line string) string {
	var obj map[string]any
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	obj[a.Key] = a.Value
	b, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(b)
}

// PrefixTransformer prepends a static string to every log line.
type PrefixTransformer struct {
	Prefix string
}

// NewPrefix creates a PrefixTransformer.
func NewPrefix(prefix string) *PrefixTransformer {
	return &PrefixTransformer{Prefix: prefix}
}

// Transform prepends the configured prefix to the line.
func (p *PrefixTransformer) Transform(line string) string {
	if p.Prefix == "" {
		return line
	}
	return strings.Join([]string{p.Prefix, line}, " ")
}
