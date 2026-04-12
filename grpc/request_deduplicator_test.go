package grpc

import (
	"testing"
	"time"
)

func TestNewRequestDeduplicator_NotNil(t *testing.T) {
	d := NewRequestDeduplicator(DefaultDeduplicatorPolicy())
	if d == nil {
		t.Fatal("expected non-nil RequestDeduplicator")
	}
}

func TestDefaultDeduplicatorPolicy_Values(t *testing.T) {
	p := DefaultDeduplicatorPolicy()
	if p.WindowDuration != 5*time.Second {
		t.Errorf("expected 5s window, got %v", p.WindowDuration)
	}
	if p.MaxSize != 256 {
		t.Errorf("expected MaxSize 256, got %d", p.MaxSize)
	}
}

func TestNewRequestDeduplicator_InvalidPolicy_FallsBackToDefault(t *testing.T) {
	d := NewRequestDeduplicator(DeduplicatorPolicy{WindowDuration: -1, MaxSize: -5})
	if d.policy.WindowDuration != DefaultDeduplicatorPolicy().WindowDuration {
		t.Errorf("expected default window duration")
	}
	if d.policy.MaxSize != DefaultDeduplicatorPolicy().MaxSize {
		t.Errorf("expected default max size")
	}
}

func TestRequestDeduplicator_Len_Empty(t *testing.T) {
	d := NewRequestDeduplicator(DefaultDeduplicatorPolicy())
	if d.Len() != 0 {
		t.Errorf("expected 0, got %d", d.Len())
	}
}

func TestRequestDeduplicator_IsDuplicate_FalseOnFirstSeen(t *testing.T) {
	d := NewRequestDeduplicator(DefaultDeduplicatorPolicy())
	if d.IsDuplicate("/pkg.Service/Method", `{"id":1}`) {
		t.Error("first occurrence should not be a duplicate")
	}
}

func TestRequestDeduplicator_IsDuplicate_TrueAfterRecord(t *testing.T) {
	d := NewRequestDeduplicator(DefaultDeduplicatorPolicy())
	d.Record("/pkg.Service/Method", `{"id":1}`)
	if !d.IsDuplicate("/pkg.Service/Method", `{"id":1}`) {
		t.Error("expected duplicate after Record")
	}
}

func TestRequestDeduplicator_IsDuplicate_DifferentPayload(t *testing.T) {
	d := NewRequestDeduplicator(DefaultDeduplicatorPolicy())
	d.Record("/pkg.Service/Method", `{"id":1}`)
	if d.IsDuplicate("/pkg.Service/Method", `{"id":2}`) {
		t.Error("different payload should not be a duplicate")
	}
}

func TestRequestDeduplicator_Eviction_AfterWindow(t *testing.T) {
	p := DeduplicatorPolicy{WindowDuration: 50 * time.Millisecond, MaxSize: 256}
	d := NewRequestDeduplicator(p)
	d.Record("/svc/Method", `{}`)
	time.Sleep(80 * time.Millisecond)
	if d.IsDuplicate("/svc/Method", `{}`) {
		t.Error("entry should have expired after window")
	}
	if d.Len() != 0 {
		t.Errorf("expected Len 0 after eviction, got %d", d.Len())
	}
}

func TestRequestDeduplicator_MaxSize_Eviction(t *testing.T) {
	p := DeduplicatorPolicy{WindowDuration: 10 * time.Second, MaxSize: 3}
	d := NewRequestDeduplicator(p)
	d.Record("m", "a")
	d.Record("m", "b")
	d.Record("m", "c")
	d.Record("m", "d") // should evict "a"
	if d.Len() != 3 {
		t.Errorf("expected Len 3, got %d", d.Len())
	}
	if d.IsDuplicate("m", "a") {
		t.Error("oldest entry should have been evicted")
	}
}
