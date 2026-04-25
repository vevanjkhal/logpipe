package transform

import (
	"testing"
)

func TestMaskTransformer_PlainText(t *testing.T) {
	m, err := NewMask(`\d{4}-\d{4}-\d{4}-\d{4}`, "[CARD]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := "payment with card 1234-5678-9012-3456 processed"
	got := m.Transform(input)
	want := "payment with card [CARD] processed"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMaskTransformer_JSON(t *testing.T) {
	m, err := NewMask(`\d{3}-\d{2}-\d{4}`, "[SSN]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := `{"user":"alice","ssn":"123-45-6789","score":99}`
	got := m.Transform(input)
	if got == input {
		t.Error("expected SSN to be masked in JSON output")
	}
	if contains(got, "123-45-6789") {
		t.Errorf("SSN not masked in output: %s", got)
	}
}

func TestMaskTransformer_DefaultMask(t *testing.T) {
	m, err := NewMask(`secret`, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := m.Transform("my secret value")
	want := "my *** value"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMaskTransformer_InvalidPattern(t *testing.T) {
	_, err := NewMask(`[invalid`, "")
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestTruncateTransformer_ShortLine(t *testing.T) {
	tr := NewTruncate(50, "")
	input := "short line"
	if got := tr.Transform(input); got != input {
		t.Errorf("got %q, want %q", got, input)
	}
}

func TestTruncateTransformer_LongLine(t *testing.T) {
	tr := NewTruncate(10, "...")
	input := "this is a very long log line that exceeds the limit"
	got := tr.Transform(input)
	if len(got) > 13 { // 10 + len("...")
		t.Errorf("output too long: %q", got)
	}
	if got[len(got)-3:] != "..." {
		t.Errorf("expected suffix '...', got %q", got)
	}
}

func TestTruncateTransformer_Defaults(t *testing.T) {
	tr := NewTruncate(0, "")
	if tr.maxLen != 256 {
		t.Errorf("expected default maxLen 256, got %d", tr.maxLen)
	}
	if tr.suffix != "..." {
		t.Errorf("expected default suffix '...', got %q", tr.suffix)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
