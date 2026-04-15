package grpc

import (
	"testing"
	"time"
)

func TestNewCircuitLog_NotNil(t *testing.T) {
	log := NewCircuitLog(10)
	if log == nil {
		t.Fatal("expected non-nil CircuitLog")
	}
}

func TestNewCircuitLog_DefaultMaxSize(t *testing.T) {
	log := NewCircuitLog(0)
	if log.MaxSize() <= 0 {
		t.Fatal("expected positive default max size")
	}
}

func TestCircuitLog_Len_Empty(t *testing.T) {
	log := NewCircuitLog(10)
	if log.Len() != 0 {
		t.Fatalf("expected 0, got %d", log.Len())
	}
}

func TestCircuitLog_Record_And_Len(t *testing.T) {
	log := NewCircuitLog(10)
	log.Record("myMethod", "open", "threshold exceeded")
	if log.Len() != 1 {
		t.Fatalf("expected 1, got %d", log.Len())
	}
}

func TestCircuitLog_Record_EmptyMethod_ReturnsError(t *testing.T) {
	log := NewCircuitLog(10)
	err := log.Record("", "open", "reason")
	if err == nil {
		t.Fatal("expected error for empty method")
	}
}

func TestCircuitLog_Eviction(t *testing.T) {
	log := NewCircuitLog(2)
	log.Record("a", "open", "r1")
	log.Record("b", "closed", "r2")
	log.Record("c", "half-open", "r3")
	if log.Len() != 2 {
		t.Fatalf("expected 2 after eviction, got %d", log.Len())
	}
}

func TestCircuitLog_TimestampAutoSet(t *testing.T) {
	log := NewCircuitLog(10)
	before := time.Now()
	log.Record("m", "open", "test")
	entries := log.Entries()
	if len(entries) == 0 {
		t.Fatal("expected at least one entry")
	}
	if entries[0].Timestamp.Before(before) {
		t.Fatal("timestamp should be set automatically")
	}
}

func TestCircuitLog_Entries_Order(t *testing.T) {
	log := NewCircuitLog(10)
	log.Record("first", "open", "r1")
	log.Record("second", "closed", "r2")
	entries := log.Entries()
	if entries[0].Method != "first" {
		t.Fatalf("expected first entry method 'first', got '%s'", entries[0].Method)
	}
}
