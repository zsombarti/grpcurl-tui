package grpc

import (
	"testing"
	"time"
)

func TestNewRequestProfiler_NotNil(t *testing.T) {
	p := NewRequestProfiler(10)
	if p == nil {
		t.Fatal("expected non-nil profiler")
	}
}

func TestNewRequestProfiler_DefaultMaxSize(t *testing.T) {
	p := NewRequestProfiler(0)
	if p.maxSize != 100 {
		t.Fatalf("expected maxSize 100, got %d", p.maxSize)
	}
}

func TestRequestProfiler_Len_Empty(t *testing.T) {
	p := NewRequestProfiler(10)
	if p.Len() != 0 {
		t.Fatalf("expected 0, got %d", p.Len())
	}
}

func TestRequestProfiler_Record_And_Len(t *testing.T) {
	p := NewRequestProfiler(10)
	p.Record(ProfileEntry{Method: "/pkg.Svc/Hello", Duration: 5 * time.Millisecond})
	if p.Len() != 1 {
		t.Fatalf("expected 1, got %d", p.Len())
	}
}

func TestRequestProfiler_TimestampAutoSet(t *testing.T) {
	p := NewRequestProfiler(10)
	before := time.Now()
	p.Record(ProfileEntry{Method: "/pkg.Svc/Hello"})
	after := time.Now()
	entries := p.All()
	if entries[0].Timestamp.Before(before) || entries[0].Timestamp.After(after) {
		t.Fatal("timestamp not auto-set correctly")
	}
}

func TestRequestProfiler_Eviction(t *testing.T) {
	p := NewRequestProfiler(3)
	for i := 0; i < 5; i++ {
		p.Record(ProfileEntry{Method: "m", Duration: time.Duration(i) * time.Millisecond})
	}
	if p.Len() != 3 {
		t.Fatalf("expected 3 after eviction, got %d", p.Len())
	}
}

func TestRequestProfiler_Summary_Empty(t *testing.T) {
	p := NewRequestProfiler(10)
	avg, rate := p.Summary()
	if avg != 0 || rate != 0 {
		t.Fatal("expected zero summary for empty profiler")
	}
}

func TestRequestProfiler_Summary_WithEntries(t *testing.T) {
	p := NewRequestProfiler(10)
	p.Record(ProfileEntry{Duration: 10 * time.Millisecond, Error: false})
	p.Record(ProfileEntry{Duration: 20 * time.Millisecond, Error: true})
	avg, rate := p.Summary()
	if avg != 15*time.Millisecond {
		t.Fatalf("expected 15ms avg, got %v", avg)
	}
	if rate != 0.5 {
		t.Fatalf("expected 0.5 error rate, got %v", rate)
	}
}

func TestRequestProfiler_Clear(t *testing.T) {
	p := NewRequestProfiler(10)
	p.Record(ProfileEntry{Method: "m"})
	p.Clear()
	if p.Len() != 0 {
		t.Fatal("expected 0 after clear")
	}
}
