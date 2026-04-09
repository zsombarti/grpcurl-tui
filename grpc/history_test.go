package grpc

import (
	"testing"
	"time"
)

func TestNewHistory_Defaults(t *testing.T) {
	h := NewHistory(0)
	if h == nil {
		t.Fatal("expected non-nil History")
	}
	if h.maxSize != 100 {
		t.Errorf("expected default maxSize 100, got %d", h.maxSize)
	}
}

func TestHistory_AddAndLen(t *testing.T) {
	h := NewHistory(10)
	h.Add(HistoryEntry{Service: "svc", Method: "Ping"})
	h.Add(HistoryEntry{Service: "svc", Method: "Pong"})
	if h.Len() != 2 {
		t.Errorf("expected 2 entries, got %d", h.Len())
	}
}

func TestHistory_Eviction(t *testing.T) {
	h := NewHistory(3)
	for i := 0; i < 5; i++ {
		h.Add(HistoryEntry{Method: "m"})
	}
	if h.Len() != 3 {
		t.Errorf("expected 3 entries after eviction, got %d", h.Len())
	}
}

func TestHistory_TimestampAutoSet(t *testing.T) {
	h := NewHistory(5)
	before := time.Now()
	h.Add(HistoryEntry{})
	after := time.Now()
	entries := h.All()
	if len(entries) == 0 {
		t.Fatal("expected at least one entry")
	}
	ts := entries[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", ts, before, after)
	}
}

func TestHistory_Clear(t *testing.T) {
	h := NewHistory(10)
	h.Add(HistoryEntry{Method: "x"})
	h.Clear()
	if h.Len() != 0 {
		t.Errorf("expected 0 entries after clear, got %d", h.Len())
	}
}

func TestHistory_AllSnapshot(t *testing.T) {
	h := NewHistory(10)
	h.Add(HistoryEntry{Method: "a"})
	h.Add(HistoryEntry{Method: "b"})
	snap := h.All()
	if len(snap) != 2 {
		t.Fatalf("expected 2, got %d", len(snap))
	}
	if snap[0].Method != "a" || snap[1].Method != "b" {
		t.Errorf("unexpected order: %v", snap)
	}
}
