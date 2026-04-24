// Package pipeline wires together a log source, an optional filter, and an
// output writer into a single processing unit.
//
// Basic usage:
//
//	p, err := pipeline.NewBuilder().
//		Source(src).
//		Level("warn").
//		Keyword("timeout").
//		Stdout().
//		Build()
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer p.Close()
//	if err := p.Run(ctx); err != nil {
//		log.Fatal(err)
//	}
//
// The pipeline reads lines from the source concurrently, passes each line
// through the configured filter, and forwards matching lines to the writer.
// Cancelling the supplied context causes Run to return promptly.
//
// Filtering
//
// Filters are applied in the order they are registered. A log line must
// satisfy all configured filters to be forwarded to the writer. Currently
// supported filters are:
//
//   - Level: discards lines whose log level is below the specified minimum
//     (e.g. "warn" drops debug and info lines).
//   - Keyword: discards lines that do not contain the given substring.
//
// Error handling
//
// Build returns an error if no source or no writer has been configured.
// Run returns the first non-nil error encountered while reading from the
// source or writing to the writer; context cancellation is not treated as
// an error.
package pipeline
