package transform

import (
	"encoding/json"
	"fmt"
	"strings"
)

// JSONPathTransformer extracts a field from a JSON log line and promotes it
// (or a derived string) as the new log line output.
type JSONPathTransformer struct {
	field    string
	fallback bool // if true, pass the original line when field is missing
}

// NewJSONPath returns a Transformer that extracts the value at the top-level
// JSON field and emits it as the new line. If fallback is true, lines that are
// not JSON or that lack the field are passed through unchanged; otherwise they
// are suppressed (empty string returned).
func NewJSONPath(field string, fallback bool) (*JSONPathTransformer, error) {
	field = strings.TrimSpace(field)
	if field == "" {
		return nil, fmt.Errorf("jsonpath: field name must not be empty")
	}
	return &JSONPathTransformer{field: field, fallback: fallback}, nil
}

// Transform implements the Transformer interface.
func (j *JSONPathTransformer) Transform(line string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		if j.fallback {
			return line
		}
		return ""
	}

	val, ok := obj[j.field]
	if !ok {
		if j.fallback {
			return line
		}
		return ""
	}

	switch v := val.(type) {
	case string:
		return v
	case nil:
		if j.fallback {
			return line
		}
		return ""
	default:
		b, err := json.Marshal(v)
		if err != nil {
			if j.fallback {
				return line
			}
			return ""
		}
		return string(b)
	}
}
