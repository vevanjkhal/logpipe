package transform_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/alert"
	"github.com/logpipe/logpipe/internal/transform"
)

func TestAlertTransformer_PassesLineThrough(t *testing.T) {
	w, err := alert.New("test", 100, func(alert.Alert) {})
	if err != nil {
		t.Fatal(err)
	}
	tr := transform.NewAlert(w)
	out, ok := tr.Transform("hello world")
	if !ok {
		t.Fatal("expected line to pass through")
	}
	if out != "hello world" {
		t.Errorf("expected 'hello world', got %q", out)
	}
}

func TestAlertTransformer_RecordsEvents(t *testing.T) {
	w, _ := alert.New("test", 100, func(alert.Alert) {})
	tr := transform.NewAlert(w)

	for i := 0; i < 5; i++ {
		tr.Transform("line")
	}
	if c := w.Count(); c != 5 {
		t.Errorf("expected count 5, got %d", c)
	}
}

func TestAlertTransformer_FiresHandler(t *testing.T) {
	var fired int
	w, _ := alert.New("test", 3, func(alert.Alert) { fired++ })
	tr := transform.NewAlert(w)

	tr.Transform("a")
	tr.Transform("b")
	if fired != 0 {
		t.Fatal("should not have fired yet")
	}
	tr.Transform("c")
	if fired != 1 {
		t.Errorf("expected 1 fire, got %d", fired)
	}
}
