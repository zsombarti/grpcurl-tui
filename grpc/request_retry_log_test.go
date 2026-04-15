package grpc

import (
	"errors"
	"testing"
)

func TestNewRetryLog_NotNil(t *testing.T) {
	rl := NewRetryLog(0)
	if rl == nil {
		t.Fatal("expected non-nil RetryLog")
	}
}

func TestNewRetryLog_DefaultMaxSize(t *testing.T) {
	rl := NewRetryLog(0)
	if rl.maxSize != defaultRetryLogMaxSize {
		t.Fatalf("expected %d, got %d", defaultRetryLogMaxSize, rl.maxSize)
	}
}

func TestRetryLog_Len_Empty(t *testing.T) {
	rl := NewRetryLog(10)
	if rl.Len() != 0 {
		t.Fatalf("expected 0, got %d", rl.Len())
	}
}

func TestRetryLog_Record_And_Len(t *testing.T) {
	rl := NewRetryLog(10)
	if err := rl.Record("SomeMethod", 1, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rl.Len() != 1 {
		t.Fatalf("expected 1, got %d", rl.Len())
	}
}

func TestRetryLog_Record_EmptyMethod_ReturnsError(t *testing.T) {
	rl := NewRetryLog(10)
	if err := rl.Record("", 1, nil); err == nil {
		t.Fatal("expected error for empty method")
	}
}

func TestRetryLog_Eviction(t *testing.T) {
	rl := NewRetryLog(3)
	for i := 0; i < 5; i++ {
		_ = rl.Record("M", i, nil)
	}
	if rl.Len() != 3 {
		t.Fatalf("expected 3 after eviction, got %d", rl.Len())
	}
}

func TestRetryLog_Entries_ContainsError(t *testing.T) {
	rl := NewRetryLog(10)
	sentinel := errors.New("timeout")
	_ = rl.Record("Foo", 2, sentinel)
	entries := rl.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Err != sentinel {
		t.Fatalf("expected sentinel error, got %v", entries[0].Err)
	}
}

func TestRetryLog_Clear(t *testing.T) {
	rl := NewRetryLog(10)
	_ = rl.Record("M", 1, nil)
	rl.Clear()
	if rl.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", rl.Len())
	}
}

func TestRetryLog_TimestampAutoSet(t *testing.T) {
	rl := NewRetryLog(10)
	_ = rl.Record("M", 1, nil)
	if rl.Entries()[0].Timestamp.IsZero() {
		t.Fatal("expected non-zero timestamp")
	}
}
