package source

import (
	"bytes"
)

// bytesReader wraps a byte slice in a bytes.Reader, used internally by SyslogSource.
func bytesReader(b []byte) *bytes.Reader {
	return bytes.NewReader(b)
}
