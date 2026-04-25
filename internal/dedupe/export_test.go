// export_test.go exposes internal fields for white-box testing.
package dedupe

import "time"

// testDeduper wraps Deduper and allows the test clock to be advanced.
type testDeduper struct {
	*Deduper
}

func newTestDeduper(base time.Time, window time.Duration) *testDeduper {
	d := New(WithWindow(window))
	d.now = func() time.Time { return base }
	return &testDeduper{d}
}

func (td *testDeduper) SetNow(t time.Time) {
	td.Deduper.now = func() time.Time { return t }
}
