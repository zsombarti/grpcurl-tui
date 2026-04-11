package grpc

import (
	"testing"
	"time"
)

func TestNewRequestTimer_NotNil(t *testing.T) {
	timer := NewRequestTimer(10)
	if timer == nil {
		t.Fatal("expected non-nil RequestTimer")
	}
}

func TestNewRequestTimer_DefaultMaxSize(t *testing.T) {
	timer := NewRequestTimer(0)
	if timer.maxSize != 100 {
		t.Fatalf("expected maxSize 100, got %d", timer.maxSize)
	}
}

func TestRequestTimer_Len_Empty(t *testing.T) {
	timer := NewRequestTimer(10)
	if timer.Len() != 0 {
		t.Fatalf("expected 0, got %d", timer.Len())
	}
}

func TestRequestTimer_Record_And_Len(t *testing.T) {
	timer := NewRequestTimer(10)
	timer.Record(10 * time.Millisecond)
	timer.Record(20 * time.Millisecond)
	if timer.Len() != 2 {
		t.Fatalf("expected 2, got %d", timer.Len())
	}
}

func TestRequestTimer_Eviction(t *testing.T) {
	timer := NewRequestTimer(2)
	timer.Record(1 * time.Millisecond)
	timer.Record(2 * time.Millisecond)
	timer.Record(3 * time.Millisecond)
	if timer.Len() != 2 {
		t.Fatalf("expected 2 after eviction, got %d", timer.Len())
	}
}

func TestRequestTimer_Average(t *testing.T) {
	timer := NewRequestTimer(10)
	timer.Record(10 * time.Millisecond)
	timer.Record(30 * time.Millisecond)
	avg := timer.Average()
	if avg != 20*time.Millisecond {
		t.Fatalf("expected 20ms average, got %v", avg)
	}
}

func TestRequestTimer_Average_Empty(t *testing.T) {
	timer := NewRequestTimer(10)
	if timer.Average() != 0 {
		t.Fatal("expected 0 average for empty timer")
	}
}

func TestRequestTimer_Max(t *testing.T) {
	timer := NewRequestTimer(10)
	timer.Record(5 * time.Millisecond)
	timer.Record(50 * time.Millisecond)
	timer.Record(15 * time.Millisecond)
	if timer.Max() != 50*time.Millisecond {
		t.Fatalf("expected 50ms max, got %v", timer.Max())
	}
}

func TestRequestTimer_Max_Empty(t *testing.T) {
	timer := NewRequestTimer(10)
	if timer.Max() != 0 {
		t.Fatal("expected 0 max for empty timer")
	}
}

func TestRequestTimer_Clear(t *testing.T) {
	timer := NewRequestTimer(10)
	timer.Record(1 * time.Millisecond)
	timer.Clear()
	if timer.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", timer.Len())
	}
}
