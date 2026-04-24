package metrics_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/logpipe/logpipe/internal/metrics"
)

func TestCounters_IncAndSnapshot(t *testing.T) {
	var c metrics.Counters

	for i := 0; i < 5; i++ {
		c.IncRead()
	}
	for i := 0; i < 2; i++ {
		c.IncFiltered()
	}
	for i := 0; i < 3; i++ {
		c.IncWritten()
	}
	c.IncError()

	snap := c.Snapshot()

	if snap.LinesRead != 5 {
		t.Errorf("expected LinesRead=5, got %d", snap.LinesRead)
	}
	if snap.LinesFiltered != 2 {
		t.Errorf("expected LinesFiltered=2, got %d", snap.LinesFiltered)
	}
	if snap.LinesWritten != 3 {
		t.Errorf("expected LinesWritten=3, got %d", snap.LinesWritten)
	}
	if snap.Errors != 1 {
		t.Errorf("expected Errors=1, got %d", snap.Errors)
	}
}

func TestSnapshot_WriteTo(t *testing.T) {
	snap := metrics.Snapshot{
		LinesRead:     10,
		LinesFiltered: 4,
		LinesWritten:  6,
		Errors:        0,
	}

	var buf bytes.Buffer
	_, err := snap.WriteTo(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{
		"lines_read=10",
		"lines_filtered=4",
		"lines_written=6",
		"errors=0",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q; got: %s", want, out)
		}
	}
}

func TestCounters_ZeroValue(t *testing.T) {
	var c metrics.Counters
	snap := c.Snapshot()

	if snap.LinesRead != 0 || snap.LinesFiltered != 0 ||
		snap.LinesWritten != 0 || snap.Errors != 0 {
		t.Error("expected all zero counters on fresh Counters")
	}
}
