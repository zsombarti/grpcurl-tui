package grpc

import (
	"testing"
)

func TestNewRequestPriorityQueue_NotNil(t *testing.T) {
	q := NewRequestPriorityQueue(10)
	if q == nil {
		t.Fatal("expected non-nil queue")
	}
}

func TestNewRequestPriorityQueue_DefaultMaxSize(t *testing.T) {
	q := NewRequestPriorityQueue(0)
	if q.maxSize != 64 {
		t.Fatalf("expected default maxSize 64, got %d", q.maxSize)
	}
}

func TestRequestPriorityQueue_Len_Empty(t *testing.T) {
	q := NewRequestPriorityQueue(10)
	if q.Len() != 0 {
		t.Fatalf("expected 0, got %d", q.Len())
	}
}

func TestRequestPriorityQueue_Enqueue_And_Len(t *testing.T) {
	q := NewRequestPriorityQueue(10)
	_ = q.Enqueue(PriorityEntry{Payload: "a", Priority: PriorityNormal})
	if q.Len() != 1 {
		t.Fatalf("expected 1, got %d", q.Len())
	}
}

func TestRequestPriorityQueue_Dequeue_OrderByPriority(t *testing.T) {
	q := NewRequestPriorityQueue(10)
	_ = q.Enqueue(PriorityEntry{Label: "low", Priority: PriorityLow})
	_ = q.Enqueue(PriorityEntry{Label: "high", Priority: PriorityHigh})
	_ = q.Enqueue(PriorityEntry{Label: "normal", Priority: PriorityNormal})

	e, ok := q.Dequeue()
	if !ok || e.Label != "high" {
		t.Fatalf("expected high priority first, got %q", e.Label)
	}
	e, _ = q.Dequeue()
	if e.Label != "normal" {
		t.Fatalf("expected normal second, got %q", e.Label)
	}
}

func TestRequestPriorityQueue_Dequeue_Empty(t *testing.T) {
	q := NewRequestPriorityQueue(10)
	_, ok := q.Dequeue()
	if ok {
		t.Fatal("expected false on empty dequeue")
	}
}

func TestRequestPriorityQueue_Full_ReturnsError(t *testing.T) {
	q := NewRequestPriorityQueue(1)
	_ = q.Enqueue(PriorityEntry{Label: "first"})
	err := q.Enqueue(PriorityEntry{Label: "second"})
	if err == nil {
		t.Fatal("expected error when queue is full")
	}
}

func TestRequestPriorityQueue_Clear_ResetsLen(t *testing.T) {
	q := NewRequestPriorityQueue(10)
	_ = q.Enqueue(PriorityEntry{Label: "x"})
	q.Clear()
	if q.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", q.Len())
	}
}
