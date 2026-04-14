package grpc

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestNewRequestDebouncer_NotNil(t *testing.T) {
	d, err := NewRequestDebouncer(DefaultDebouncerPolicy(), func(string, map[string]interface{}) {})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil debouncer")
	}
}

func TestDefaultDebouncerPolicy_Values(t *testing.T) {
	p := DefaultDebouncerPolicy()
	if p.Wait <= 0 {
		t.Error("expected positive Wait")
	}
	if p.MaxWait <= p.Wait {
		t.Error("expected MaxWait > Wait")
	}
}

func TestNewRequestDebouncer_NilFlush_ReturnsError(t *testing.T) {
	_, err := NewRequestDebouncer(DefaultDebouncerPolicy(), nil)
	if err == nil {
		t.Fatal("expected error for nil flush")
	}
}

func TestRequestDebouncer_Len_Empty(t *testing.T) {
	d, _ := NewRequestDebouncer(DefaultDebouncerPolicy(), func(string, map[string]interface{}) {})
	if d.Len() != 0 {
		t.Errorf("expected 0, got %d", d.Len())
	}
}

func TestRequestDebouncer_Debounce_EmptyMethod_ReturnsError(t *testing.T) {
	d, _ := NewRequestDebouncer(DefaultDebouncerPolicy(), func(string, map[string]interface{}) {})
	if err := d.Debounce("", nil); err == nil {
		t.Fatal("expected error for empty method")
	}
}

func TestRequestDebouncer_Debounce_FlushesAfterWait(t *testing.T) {
	var called int32
	policy := DebouncerPolicy{Wait: 50 * time.Millisecond, MaxWait: 500 * time.Millisecond}
	d, _ := NewRequestDebouncer(policy, func(method string, _ map[string]interface{}) {
		atomic.AddInt32(&called, 1)
	})

	_ = d.Debounce("svc.Method", map[string]interface{}{"k": "v"})
	if d.Len() != 1 {
		t.Errorf("expected 1 pending, got %d", d.Len())
	}

	time.Sleep(120 * time.Millisecond)
	if atomic.LoadInt32(&called) != 1 {
		t.Errorf("expected flush called once, got %d", atomic.LoadInt32(&called))
	}
	if d.Len() != 0 {
		t.Errorf("expected 0 pending after flush, got %d", d.Len())
	}
}

func TestRequestDebouncer_Cancel_RemovesPending(t *testing.T) {
	var called int32
	policy := DebouncerPolicy{Wait: 100 * time.Millisecond, MaxWait: 1 * time.Second}
	d, _ := NewRequestDebouncer(policy, func(string, map[string]interface{}) {
		atomic.AddInt32(&called, 1)
	})

	_ = d.Debounce("svc.Method", nil)
	d.Cancel("svc.Method")
	if d.Len() != 0 {
		t.Errorf("expected 0 pending after cancel, got %d", d.Len())
	}
	time.Sleep(150 * time.Millisecond)
	if atomic.LoadInt32(&called) != 0 {
		t.Errorf("expected flush not called after cancel, got %d", atomic.LoadInt32(&called))
	}
}

func TestRequestDebouncer_Debounce_CollapsesCalls(t *testing.T) {
	var flushed []map[string]interface{}
	policy := DebouncerPolicy{Wait: 60 * time.Millisecond, MaxWait: 600 * time.Millisecond}
	d, _ := NewRequestDebouncer(policy, func(_ string, p map[string]interface{}) {
		flushed = append(flushed, p)
	})

	_ = d.Debounce("svc.M", map[string]interface{}{"n": 1})
	time.Sleep(20 * time.Millisecond)
	_ = d.Debounce("svc.M", map[string]interface{}{"n": 2})
	time.Sleep(20 * time.Millisecond)
	_ = d.Debounce("svc.M", map[string]interface{}{"n": 3})

	time.Sleep(150 * time.Millisecond)
	if len(flushed) != 1 {
		t.Errorf("expected 1 flush, got %d", len(flushed))
	}
	if flushed[0]["n"] != 3 {
		t.Errorf("expected last payload, got %v", flushed[0])
	}
}
