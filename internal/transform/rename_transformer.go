package transform

import (
	"encoding/json"
)

// renameTransformer renames a JSON field from one key to another.
// If the line is not valid JSON or the source field is missing, the line is
// passed through unchanged.
type renameTransformer struct {
	from string
	to   string
}

// NewRename returns a Transformer that renames the JSON field `from` to `to`.
// If the input line is not JSON, or the field does not exist, it is passed
// through as-is.
func NewRename(from, to string) (Transformer, error) {
	if from == "" {
		return nil, errorf("rename: 'from' field name must not be empty")
	}
	if to == "" {
		return nil, errorf("rename: 'to' field name must not be empty")
	}
	return &renameTransformer{from: from, to: to}, nil
}

func (r *renameTransformer) Transform(line string) (string, bool) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		// Not JSON — pass through unchanged.
		return line, true
	}

	val, ok := obj[r.from]
	if !ok {
		// Field not present — pass through unchanged.
		return line, true
	}

	delete(obj, r.from)
	obj[r.to] = val

	out, err := json.Marshal(obj)
	if err != nil {
		return line, true
	}
	return string(out), true
}
