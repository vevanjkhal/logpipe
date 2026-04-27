package alert

import (
	"testing"
	"time"
)

func TestNew_InvalidThreshold(t *testing.T) {
	_, err := New("test", 0, func(Alert) {})
	if err == nil {
		t.Fatal("expected error for threshold < 1")
	}
}

func TestNew_NilHandler(t *testing.T) {
	_, err := New("test", 1, nil)
	if err == nil {
		t.Fatal("expected error for nil handler")
	}
}

func TestRecord_FiresAtThreshold(t *testing.T) {
	var fired []Alert
	w, err := New("errors", 3, func(a Alert) { fired = append(fired, a) })
	if err != nil {
		t.Fatal(err)
	}

	w.Record()
	w.Record()
	if len(fired) != 0 {
		t.Fatalf("expected 0 alerts before threshold, got %d", len(fired))
	}
	w.Record()
	if len(fired) != 1 {
		t.Fatalf("expected 1 alert at threshold, got %d", len(fired))
	}
	if fired[0].Count != 3 {
		t.Errorf("expected count 3, got %d", fired[0].Count)
	}
	if fired[0].Name != "errors" {
		t.Errorf("expected name 'errors', got %s", fired[0].Name)
	}
}

func TestRecord_ResetsAfterFiring(t *testing.T) {
	var fired int
	w, _ := New("x", 2, func(Alert) { fired++ })
	w.Record()
	w.Record() // fires
	w.Record()
	w.Record() // fires again
	if fired != 2 {
		t.Errorf("expected 2 fires, got %d", fired)
	}
}

func TestRecord_WindowExpiry(t *testing.T) {
	var fired int
	w, _ := New("x", 3, func(Alert) { fired++ }, WithWindow(50*time.Millisecond))

	base := time.Now()
	w.now = func() time.Time { return base }
	w.Record()
	w.Record()

	// Advance past the window — old events should be evicted.
	w.now = func() time.Time { return base.Add(100 * time.Millisecond) }
	w.Record()
	w.Record()

	if fired != 0 {
		t.Errorf("expected 0 fires after window expiry, got %d", fired)
	}
}

func TestCount_ReturnsWindowCount(t *testing.T) {
	w, _ := New("x", 100, func(Alert) {})
	w.Record()
	w.Record()
	if c := w.Count(); c != 2 {
		t.Errorf("expected count 2, got %d", c)
	}
}

func TestWithLevel_SetsLevel(t *testing.T) {
	var got Level
	w, _ := New("x", 1, func(a Alert) { got = a.Level }, WithLevel(LevelError))
	w.Record()
	if got != LevelError {
		t.Errorf("expected LevelError, got %s", got)
	}
}
