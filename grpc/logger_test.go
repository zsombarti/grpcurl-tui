package grpc

import (
	"strings"
	"testing"
)

func TestNewLogger_NotNil(t *testing.T) {
	l := NewLogger(100, LogLevelDebug)
	if l == nil {
		t.Fatal("expected non-nil Logger")
	}
}

func TestNewLogger_DefaultMaxSize(t *testing.T) {
	l := NewLogger(0, LogLevelDebug)
	if l.maxSize != 200 {
		t.Fatalf("expected default maxSize 200, got %d", l.maxSize)
	}
}

func TestLogger_LogLevel_String(t *testing.T) {
	cases := map[LogLevel]string{
		LogLevelDebug: "DEBUG",
		LogLevelInfo:  "INFO",
		LogLevelWarn:  "WARN",
		LogLevelError: "ERROR",
	}
	for level, want := range cases {
		if got := level.String(); got != want {
			t.Errorf("level %d: got %q, want %q", level, got, want)
		}
	}
}

func TestLogger_Append_And_Len(t *testing.T) {
	l := NewLogger(10, LogLevelDebug)
	l.Info("hello %s", "world")
	l.Warn("something off")
	if l.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", l.Len())
	}
}

func TestLogger_Eviction(t *testing.T) {
	l := NewLogger(3, LogLevelDebug)
	l.Debug("a")
	l.Debug("b")
	l.Debug("c")
	l.Debug("d")
	if l.Len() != 3 {
		t.Fatalf("expected 3 entries after eviction, got %d", l.Len())
	}
	entries := l.Entries()
	if entries[0].Message != "b" {
		t.Errorf("expected oldest entry to be evicted, got %q", entries[0].Message)
	}
}

func TestLogger_MinLevel_Filters(t *testing.T) {
	l := NewLogger(50, LogLevelWarn)
	l.Debug("debug msg")
	l.Info("info msg")
	l.Warn("warn msg")
	l.Error("error msg")
	if l.Len() != 2 {
		t.Fatalf("expected 2 entries (warn+error), got %d", l.Len())
	}
}

func TestLogger_Entries_Snapshot(t *testing.T) {
	l := NewLogger(10, LogLevelDebug)
	l.Error("boom")
	entries := l.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if !strings.Contains(entries[0].Message, "boom") {
		t.Errorf("unexpected message: %q", entries[0].Message)
	}
	if entries[0].Level != LogLevelError {
		t.Errorf("expected ERROR level, got %s", entries[0].Level)
	}
}

func TestLogger_Clear(t *testing.T) {
	l := NewLogger(10, LogLevelDebug)
	l.Info("a")
	l.Info("b")
	l.Clear()
	if l.Len() != 0 {
		t.Fatalf("expected 0 entries after Clear, got %d", l.Len())
	}
}
