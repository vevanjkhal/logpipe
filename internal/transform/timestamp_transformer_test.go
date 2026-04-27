package transform_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/transform"
)

func fixedTime() time.Time {
	return time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
}

func newTestTimestamp(opts ...transform.TimestampOption) transform.Transformer {
	// Inject a fixed clock via the unexported option exposed for tests.
	all := append([]transform.TimestampOption{transform.WithTimestampClock(fixedTime)}, opts...)
	return transform.NewTimestamp(all...)
}

func TestTimestampTransformer_PlainText(t *testing.T) {
	tr := newTestTimestamp()
	out, keep := tr.Transform("hello world")
	if !keep {
		t.Fatal("expected line to be kept")
	}
	if !strings.HasPrefix(out, "2024-06-01T12:00:00Z ") {
		t.Fatalf("unexpected output: %q", out)
	}
	if !strings.HasSuffix(out, "hello world") {
		t.Fatalf("unexpected suffix in output: %q", out)
	}
}

func TestTimestampTransformer_JSON(t *testing.T) {
	tr := newTestTimestamp()
	out, keep := tr.Transform(`{"level":"info","msg":"ok"}`)
	if !keep {
		t.Fatal("expected line to be kept")
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["ts"] != "2024-06-01T12:00:00Z" {
		t.Fatalf("unexpected ts value: %v", obj["ts"])
	}
}

func TestTimestampTransformer_CustomField(t *testing.T) {
	tr := newTestTimestamp(transform.WithTimestampField("timestamp"))
	out, _ := tr.Transform(`{"msg":"hi"}`)
	var obj map[string]interface{}
	_ = json.Unmarshal([]byte(out), &obj)
	if _, ok := obj["timestamp"]; !ok {
		t.Fatalf("expected 'timestamp' field in %q", out)
	}
}

func TestTimestampTransformer_CustomFormat(t *testing.T) {
	tr := newTestTimestamp(transform.WithTimestampFormat("2006-01-02"))
	out, _ := tr.Transform("plain line")
	if !strings.HasPrefix(out, "2024-06-01 ") {
		t.Fatalf("unexpected format in output: %q", out)
	}
}
