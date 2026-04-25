package routing_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/routing"
)

func TestRouter_NoRules_ReturnsDefault(t *testing.T) {
	r := routing.New("default")
	if got := r.Route("anything"); got != "default" {
		t.Fatalf("expected default, got %q", got)
	}
}

func TestRouter_PlainTextMatch(t *testing.T) {
	r := routing.New("default",
		routing.Rule{Contains: "ERROR", Destination: "errors"},
	)
	if got := r.Route("2024/01/01 ERROR something broke"); got != "errors" {
		t.Fatalf("expected errors, got %q", got)
	}
}

func TestRouter_PlainTextNoMatch_ReturnsDefault(t *testing.T) {
	r := routing.New("default",
		routing.Rule{Contains: "ERROR", Destination: "errors"},
	)
	if got := r.Route("INFO all good"); got != "default" {
		t.Fatalf("expected default, got %q", got)
	}
}

func TestRouter_JSONFieldMatch(t *testing.T) {
	r := routing.New("default",
		routing.Rule{Field: "level", Contains: "error", Destination: "errors"},
	)
	line := `{"level":"error","msg":"oops"}`
	if got := r.Route(line); got != "errors" {
		t.Fatalf("expected errors, got %q", got)
	}
}

func TestRouter_JSONFieldMissing_ReturnsDefault(t *testing.T) {
	r := routing.New("default",
		routing.Rule{Field: "level", Contains: "error", Destination: "errors"},
	)
	line := `{"msg":"no level field"}`
	if got := r.Route(line); got != "default" {
		t.Fatalf("expected default, got %q", got)
	}
}

func TestRouter_FirstMatchWins(t *testing.T) {
	r := routing.New("default",
		routing.Rule{Contains: "ERROR", Destination: "first"},
		routing.Rule{Contains: "ERROR", Destination: "second"},
	)
	if got := r.Route("ERROR"); got != "first" {
		t.Fatalf("expected first, got %q", got)
	}
}

func TestRouter_InvalidJSON_FallsBackToDefault(t *testing.T) {
	r := routing.New("default",
		routing.Rule{Field: "level", Contains: "error", Destination: "errors"},
	)
	if got := r.Route("not json at all"); got != "default" {
		t.Fatalf("expected default, got %q", got)
	}
}
