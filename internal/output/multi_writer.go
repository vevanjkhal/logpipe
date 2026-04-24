package output

import "fmt"

// MultiWriter fans out log entries to multiple Writer implementations.
type MultiWriter struct {
	writers []Writer
}

// NewMultiWriter creates a MultiWriter that writes to all provided writers.
func NewMultiWriter(writers ...Writer) *MultiWriter {
	return &MultiWriter{writers: writers}
}

// Write sends the entry to every underlying writer, collecting errors.
func (m *MultiWriter) Write(entry LogEntry) error {
	var errs []error
	for _, w := range m.writers {
		if err := w.Write(entry); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("multiwriter: %d write error(s): %v", len(errs), errs)
	}
	return nil
}

// Close closes all underlying writers, collecting errors.
func (m *MultiWriter) Close() error {
	var errs []error
	for _, w := range m.writers {
		if err := w.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("multiwriter: %d close error(s): %v", len(errs), errs)
	}
	return nil
}

// Add appends a writer to the fan-out list at runtime.
func (m *MultiWriter) Add(w Writer) {
	m.writers = append(m.writers, w)
}
