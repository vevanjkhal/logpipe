package transform

import "strings"

// DropTransformer discards any line that contains one of the configured
// substrings (case-sensitive). It is the inverse of a keep/allow-list filter.
type DropTransformer struct {
	substrings []string
}

// NewDrop returns a Transformer that drops lines containing any of the given
// substrings. An empty slice means nothing is dropped.
func NewDrop(substrings []string) *DropTransformer {
	copied := make([]string, len(substrings))
	copy(copied, substrings)
	return &DropTransformer{substrings: copied}
}

// Transform returns an empty string when the line matches any of the drop
// substrings, signalling to the pipeline that the line should be discarded.
func (d *DropTransformer) Transform(line string) string {
	for _, s := range d.substrings {
		if s != "" && strings.Contains(line, s) {
			return ""
		}
	}
	return line
}
