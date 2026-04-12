package grpc

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewRequestBatcher_NotNil(t *testing.T) {
	b := NewRequestBatcher(DefaultBatchPolicy())
	if b == nil {
		t.Fatal("expected non-nil RequestBatcher")
	}
}

func TestDefaultBatchPolicy_Values(t *testing.T) {
	p := DefaultBatchPolicy()
	if p.MaxSize != 20 {
		t.Errorf("expected MaxSize 20, got %d", p.MaxSize)
	}
	if p.MaxWait != 200*time.Millisecond {
		t.Errorf("expected MaxWait 200ms, got %v", p.MaxWait)
	}
	if p.Concurrency != 4 {
		t.Errorf("expected Concurrency 4, got %d", p.Concurrency)
	}
}

func TestNewRequestBatcher_InvalidPolicy_FallsBackToDefault(t *testing.T) {
	b := NewRequestBatcher(BatchPolicy{MaxSize: -1, MaxWait: -1, Concurrency: -1})
	def := DefaultBatchPolicy()
	if b.policy.MaxSize != def.MaxSize {
		t.Errorf("expected MaxSize %d, got %d", def.MaxSize, b.policy.MaxSize)
	}
}

func TestRequestBatcher_Len_Empty(t *testing.T) {
	b := NewRequestBatcher(DefaultBatchPolicy())
	if b.Len() != 0 {
		t.Errorf("expected 0, got %d", b.Len())
	}
}

func TestRequestBatcher_Add_And_Len(t *testing.T) {
	b := NewRequestBatcher(DefaultBatchPolicy())
	if err := b.Add(`{"key":"value"}`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Len() != 1 {
		t.Errorf("expected 1, got %d", b.Len())
	}
}

func TestRequestBatcher_Add_Full_ReturnsError(t *testing.T) {
	b := NewRequestBatcher(BatchPolicy{MaxSize: 2, MaxWait: time.Millisecond, Concurrency: 1})
	_ = b.Add("a")
	_ = b.Add("b")
	if err := b.Add("c"); err == nil {
		t.Error("expected error when queue is full")
	}
}

func TestRequestBatcher_Flush_DrainQueue(t *testing.T) {
	b := NewRequestBatcher(DefaultBatchPolicy())
	_ = b.Add("payload1")
	_ = b.Add("payload2")

	results := b.Flush(context.Background(), func(_ context.Context, p string) error {
		return nil
	})

	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
	if b.Len() != 0 {
		t.Errorf("expected queue drained, got %d", b.Len())
	}
}

func TestRequestBatcher_Flush_PropagatesError(t *testing.T) {
	b := NewRequestBatcher(DefaultBatchPolicy())
	_ = b.Add("bad")

	results := b.Flush(context.Background(), func(_ context.Context, p string) error {
		return errors.New("invoke failed")
	})

	if results[0].Err == nil {
		t.Error("expected error in result")
	}
}

func TestRequestBatcher_Policy_RoundTrip(t *testing.T) {
	p := BatchPolicy{MaxSize: 5, MaxWait: 50 * time.Millisecond, Concurrency: 2}
	b := NewRequestBatcher(p)
	got := b.Policy()
	if got.MaxSize != 5 || got.Concurrency != 2 {
		t.Errorf("policy round-trip mismatch: %+v", got)
	}
}
