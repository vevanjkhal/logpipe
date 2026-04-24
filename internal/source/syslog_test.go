package source

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestSyslogSource_Name(t *testing.T) {
	s, err := NewSyslogSource(":0")
	if err != nil {
		t.Fatalf("NewSyslogSource: %v", err)
	}
	defer s.conn.Close()

	if s.Name() == "" {
		t.Error("expected non-empty name")
	}
}

func TestSyslogSource_Lines(t *testing.T) {
	s, err := NewSyslogSource(":0")
	if err != nil {
		t.Fatalf("NewSyslogSource: %v", err)
	}

	port := s.conn.LocalAddr().(*net.UDPAddr).Port

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	lines := s.Lines(ctx)

	msg := "hello syslog"
	go func() {
		conn, err := net.Dial("udp", fmt.Sprintf("127.0.0.1:%d", port))
		if err != nil {
			return
		}
		defer conn.Close()
		_, _ = fmt.Fprint(conn, msg)
	}()

	select {
	case got, ok := <-lines:
		if !ok {
			t.Fatal("channel closed before receiving line")
		}
		if got != msg {
			t.Errorf("got %q, want %q", got, msg)
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for syslog line")
	}
}

func TestSyslogSource_ContextCancel(t *testing.T) {
	s, err := NewSyslogSource(":0")
	if err != nil {
		t.Fatalf("NewSyslogSource: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	lines := s.Lines(ctx)
	cancel()

	select {
	case _, ok := <-lines:
		if ok {
			t.Error("expected channel to be closed after cancel")
		}
	case <-time.After(time.Second):
		t.Error("timed out waiting for channel to close after cancel")
	}
}
