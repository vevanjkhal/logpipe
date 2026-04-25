package routing_test

import (
	"bytes"
	"testing"

	"github.com/yourorg/logpipe/internal/routing"
)

func TestDispatcher_RoutesToRegisteredSink(t *testing.T) {
	r := routing.New("default",
		routing.Rule{Contains: "ERROR", Destination: "errors"},
	)
	var errBuf, defBuf bytes.Buffer
	d := routing.NewDispatcher(r, &defBuf)
	d.Register("errors", &errBuf)

	if err := d.Dispatch("ERROR something bad"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if errBuf.String() != "ERROR something bad\n" {
		t.Fatalf("errors sink got %q", errBuf.String())
	}
	if defBuf.Len() != 0 {
		t.Fatalf("default sink should be empty, got %q", defBuf.String())
	}
}

func TestDispatcher_FallsBackWhenSinkNotRegistered(t *testing.T) {
	r := routing.New("default",
		routing.Rule{Contains: "ERROR", Destination: "errors"},
	)
	var defBuf bytes.Buffer
	d := routing.NewDispatcher(r, &defBuf)
	// "errors" sink is intentionally NOT registered

	if err := d.Dispatch("ERROR no sink"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if defBuf.String() != "ERROR no sink\n" {
		t.Fatalf("fallback got %q", defBuf.String())
	}
}

func TestDispatcher_DefaultDestinationToFallback(t *testing.T) {
	r := routing.New("default")
	var defBuf bytes.Buffer
	d := routing.NewDispatcher(r, &defBuf)

	if err := d.Dispatch("INFO hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if defBuf.String() != "INFO hello\n" {
		t.Fatalf("fallback got %q", defBuf.String())
	}
}

func TestDispatcher_MultipleRules(t *testing.T) {
	r := routing.New("default",
		routing.Rule{Contains: "ERROR", Destination: "errors"},
		routing.Rule{Contains: "WARN", Destination: "warnings"},
	)
	var errBuf, warnBuf, defBuf bytes.Buffer
	d := routing.NewDispatcher(r, &defBuf)
	d.Register("errors", &errBuf)
	d.Register("warnings", &warnBuf)

	_ = d.Dispatch("ERROR bad")
	_ = d.Dispatch("WARN careful")
	_ = d.Dispatch("INFO ok")

	if errBuf.String() != "ERROR bad\n" {
		t.Errorf("errors: %q", errBuf.String())
	}
	if warnBuf.String() != "WARN careful\n" {
		t.Errorf("warnings: %q", warnBuf.String())
	}
	if defBuf.String() != "INFO ok\n" {
		t.Errorf("default: %q", defBuf.String())
	}
}
