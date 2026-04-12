package grpc

import (
	"testing"
	"time"
)

func TestNewRequestTracer_NotNil(t *testing.T) {
	tr := NewRequestTracer(10)
	if tr == nil {
		t.Fatal("expected non-nil RequestTracer")
	}
}

func TestNewRequestTracer_DefaultMaxSize(t *testing.T) {
	tr := NewRequestTracer(0)
	if tr.maxSize != 100 {
		t.Fatalf("expected maxSize 100, got %d", tr.maxSize)
	}
}

func TestRequestTracer_Len_Empty(t *testing.T) {
	tr := NewRequestTracer(10)
	if tr.Len() != 0 {
		t.Fatalf("expected 0, got %d", tr.Len())
	}
}

func TestRequestTracer_Start_And_Len(t *testing.T) {
	tr := NewRequestTracer(10)
	tr.Start("trace-1", "/pkg.Svc/Method", "localhost:50051", nil)
	if tr.Len() != 1 {
		t.Fatalf("expected 1, got %d", tr.Len())
	}
}

func TestRequestTracer_Finish_SetsDuration(t *testing.T) {
	tr := NewRequestTracer(10)
	idx := tr.Start("trace-1", "/pkg.Svc/Method", "localhost:50051", nil)
	time.Sleep(2 * time.Millisecond)
	tr.Finish(idx, "")
	spans := tr.Spans()
	if spans[idx].Duration <= 0 {
		t.Fatal("expected positive duration")
	}
}

func TestRequestTracer_Finish_RecordsError(t *testing.T) {
	tr := NewRequestTracer(10)
	idx := tr.Start("trace-2", "/pkg.Svc/Method", "localhost:50051", nil)
	tr.Finish(idx, "connection refused")
	spans := tr.Spans()
	if spans[idx].Error != "connection refused" {
		t.Fatalf("expected error 'connection refused', got %q", spans[idx].Error)
	}
}

func TestRequestTracer_Finish_OutOfRange_NoOp(t *testing.T) {
	tr := NewRequestTracer(10)
	tr.Finish(99, "should not panic")
}

func TestRequestTracer_Eviction(t *testing.T) {
	tr := NewRequestTracer(3)
	for i := 0; i < 5; i++ {
		tr.Start("id", "/Svc/M", "addr", nil)
	}
	if tr.Len() != 3 {
		t.Fatalf("expected 3 after eviction, got %d", tr.Len())
	}
}

func TestRequestTracer_Clear(t *testing.T) {
	tr := NewRequestTracer(10)
	tr.Start("t1", "/Svc/M", "addr", nil)
	tr.Clear()
	if tr.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", tr.Len())
	}
}

func TestRequestTracer_Spans_ReturnsCopy(t *testing.T) {
	tr := NewRequestTracer(10)
	tr.Start("t1", "/Svc/M", "addr", map[string]string{"key": "val"})
	spans := tr.Spans()
	spans[0].Method = "mutated"
	original := tr.Spans()
	if original[0].Method == "mutated" {
		t.Fatal("Spans should return a copy, not a reference")
	}
}
