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
package pipeline
