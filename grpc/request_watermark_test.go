package grpc

import (
	"testing"
	"time"
)

func TestNewRequestWatermark_NotNil(t *testing.T) {
	w := NewRequestWatermark(0)
	if w == nil {
		t.Fatal("expected non-nil")
	}
}

func TestNewRequestWatermark_DefaultMaxSize(t *testing.T) {
	w := NewRequestWatermark(0)
	if w.max != defaultWatermarkMaxSize {
		t.Fatalf("expected %d, got %d", defaultWatermarkMaxSize, w.max)
	}
}

func TestRequestWatermark_Len_Empty(t *testing.T) {
	w := NewRequestWatermark(0)
	if w.Len() != 0 {
		t.Fatal("expected 0")
	}
}

func TestRequestWatermark_Record_And_Len(t *testing.T) {
	w := NewRequestWatermark(0)
	if err := w.Record("svc/Method", 10*time.Millisecond); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.Len() != 1 {
		t.Fatalf("expected 1, got %d", w.Len())
	}
}

func TestRequestWatermark_Record_EmptyMethod_ReturnsError(t *testing.T) {
	w := NewRequestWatermark(0)
	if err := w.Record("", 5*time.Millisecond); err == nil {
		t.Fatal("expected error for empty method")
	}
}

func TestRequestWatermark_MinMax(t *testing.T) {
	w := NewRequestWatermark(0)
	_ = w.Record("svc/A", 20*time.Millisecond)
	_ = w.Record("svc/A", 5*time.Millisecond)
	_ = w.Record("svc/A", 50*time.Millisecond)
	e, ok := w.Get("svc/A")
	if !ok {
		t.Fatal("expected entry")
	}
	if e.Min != 5*time.Millisecond {
		t.Fatalf("expected min 5ms, got %v", e.Min)
	}
	if e.Max != 50*time.Millisecond {
		t.Fatalf("expected max 50ms, got %v", e.Max)
	}
	if e.Count != 3 {
		t.Fatalf("expected count 3, got %d", e.Count)
	}
}

func TestRequestWatermark_Clear(t *testing.T) {
	w := NewRequestWatermark(0)
	_ = w.Record("svc/A", 10*time.Millisecond)
	w.Clear()
	if w.Len() != 0 {
		t.Fatal("expected 0 after clear")
	}
}

func TestRequestWatermark_MaxBuckets(t *testing.T) {
	w := NewRequestWatermark(2)
	_ = w.Record("svc/A", 1*time.Millisecond)
	_ = w.Record("svc/B", 1*time.Millisecond)
	if err := w.Record("svc/C", 1*time.Millisecond); err == nil {
		t.Fatal("expected error when max buckets reached")
	}
}
