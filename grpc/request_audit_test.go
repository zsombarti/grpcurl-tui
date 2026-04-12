package grpc

import (
	"testing"
	"time"
)

func TestNewRequestAuditor_NotNil(t *testing.T) {
	a := NewRequestAuditor(0)
	if a == nil {
		t.Fatal("expected non-nil auditor")
	}
}

func TestNewRequestAuditor_DefaultMaxSize(t *testing.T) {
	a := NewRequestAuditor(0)
	if a.maxSize != defaultAuditMaxSize {
		t.Fatalf("expected maxSize %d, got %d", defaultAuditMaxSize, a.maxSize)
	}
}

func TestRequestAuditor_Len_Empty(t *testing.T) {
	a := NewRequestAuditor(10)
	if a.Len() != 0 {
		t.Fatalf("expected 0, got %d", a.Len())
	}
}

func TestRequestAuditor_Record_And_Len(t *testing.T) {
	a := NewRequestAuditor(10)
	a.Record(AuditEntry{Method: "/pkg.Svc/Method", Address: "localhost:50051", Status: "ok"})
	if a.Len() != 1 {
		t.Fatalf("expected 1, got %d", a.Len())
	}
}

func TestRequestAuditor_TimestampAutoSet(t *testing.T) {
	a := NewRequestAuditor(10)
	before := time.Now()
	a.Record(AuditEntry{Method: "/pkg.Svc/M"})
	after := time.Now()
	entries := a.Entries()
	if entries[0].Timestamp.Before(before) || entries[0].Timestamp.After(after) {
		t.Fatal("timestamp not within expected range")
	}
}

func TestRequestAuditor_Eviction(t *testing.T) {
	a := NewRequestAuditor(3)
	for i := 0; i < 5; i++ {
		a.Record(AuditEntry{Method: "/M", Note: string(rune('A' + i))})
	}
	if a.Len() != 3 {
		t.Fatalf("expected 3 after eviction, got %d", a.Len())
	}
	if a.Entries()[0].Note != "C" {
		t.Fatalf("expected oldest surviving entry Note=C, got %s", a.Entries()[0].Note)
	}
}

func TestRequestAuditor_Clear(t *testing.T) {
	a := NewRequestAuditor(10)
	a.Record(AuditEntry{Method: "/M"})
	a.Clear()
	if a.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", a.Len())
	}
}

func TestRequestAuditor_Entries_IsCopy(t *testing.T) {
	a := NewRequestAuditor(10)
	a.Record(AuditEntry{Method: "/M"})
	e := a.Entries()
	e[0].Method = "mutated"
	if a.Entries()[0].Method == "mutated" {
		t.Fatal("Entries should return a copy, not a reference")
	}
}
