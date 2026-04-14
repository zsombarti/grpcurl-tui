package grpc

import (
	"testing"
)

func TestNewRequestFence_NotNil(t *testing.T) {
	f := NewRequestFence(DefaultFencePolicy())
	if f == nil {
		t.Fatal("expected non-nil RequestFence")
	}
}

func TestDefaultFencePolicy_Values(t *testing.T) {
	p := DefaultFencePolicy()
	if p.MaxConcurrent != 10 {
		t.Errorf("expected MaxConcurrent=10, got %d", p.MaxConcurrent)
	}
}

func TestNewRequestFence_InvalidPolicy_FallsBackToDefault(t *testing.T) {
	f := NewRequestFence(FencePolicy{MaxConcurrent: -5})
	if f.Policy().MaxConcurrent != 10 {
		t.Errorf("expected fallback MaxConcurrent=10, got %d", f.Policy().MaxConcurrent)
	}
}

func TestRequestFence_Active_InitiallyZero(t *testing.T) {
	f := NewRequestFence(DefaultFencePolicy())
	if f.Active() != 0 {
		t.Errorf("expected Active=0, got %d", f.Active())
	}
}

func TestRequestFence_Acquire_And_Active(t *testing.T) {
	f := NewRequestFence(DefaultFencePolicy())
	if err := f.Acquire(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Active() != 1 {
		t.Errorf("expected Active=1, got %d", f.Active())
	}
}

func TestRequestFence_Release_DecrementsActive(t *testing.T) {
	f := NewRequestFence(DefaultFencePolicy())
	_ = f.Acquire()
	f.Release()
	if f.Active() != 0 {
		t.Errorf("expected Active=0 after release, got %d", f.Active())
	}
}

func TestRequestFence_Release_NeverNegative(t *testing.T) {
	f := NewRequestFence(DefaultFencePolicy())
	f.Release()
	if f.Active() != 0 {
		t.Errorf("expected Active=0, got %d", f.Active())
	}
}

func TestRequestFence_Acquire_BlocksAtMax(t *testing.T) {
	f := NewRequestFence(FencePolicy{MaxConcurrent: 2})
	_ = f.Acquire()
	_ = f.Acquire()
	if err := f.Acquire(); err == nil {
		t.Error("expected error when fence is full, got nil")
	}
	if f.Active() != 2 {
		t.Errorf("expected Active=2, got %d", f.Active())
	}
}
