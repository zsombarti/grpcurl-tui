package grpc

import (
	"context"
	"testing"
	"time"
)

func TestDefaultThrottlePolicy_Values(t *testing.T) {
	p := DefaultThrottlePolicy()
	if p.MaxConcurrent != 10 {
		t.Fatalf("expected MaxConcurrent=10, got %d", p.MaxConcurrent)
	}
	if p.QueueTimeout != 5*time.Second {
		t.Fatalf("expected QueueTimeout=5s, got %v", p.QueueTimeout)
	}
}

func TestNewRequestThrottle_NotNil(t *testing.T) {
	th := NewRequestThrottle(DefaultThrottlePolicy())
	if th == nil {
		t.Fatal("expected non-nil RequestThrottle")
	}
}

func TestNewRequestThrottle_InvalidPolicy_FallsBackToDefault(t *testing.T) {
	th := NewRequestThrottle(ThrottlePolicy{MaxConcurrent: -1, QueueTimeout: -1})
	def := DefaultThrottlePolicy()
	if th.Policy().MaxConcurrent != def.MaxConcurrent {
		t.Fatalf("expected MaxConcurrent=%d, got %d", def.MaxConcurrent, th.Policy().MaxConcurrent)
	}
}

func TestRequestThrottle_Acquire_And_Active(t *testing.T) {
	th := NewRequestThrottle(ThrottlePolicy{MaxConcurrent: 3, QueueTimeout: time.Second})
	ctx := context.Background()

	if err := th.Acquire(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th.Active() != 1 {
		t.Fatalf("expected active=1, got %d", th.Active())
	}
	th.Release()
	if th.Active() != 0 {
		t.Fatalf("expected active=0 after release, got %d", th.Active())
	}
}

func TestRequestThrottle_Acquire_CancelledContext(t *testing.T) {
	th := NewRequestThrottle(ThrottlePolicy{MaxConcurrent: 1, QueueTimeout: 2 * time.Second})
	ctx := context.Background()
	_ = th.Acquire(ctx) // fill the slot

	cancelled, cancel := context.WithCancel(context.Background())
	cancel()

	if err := th.Acquire(cancelled); err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestRequestThrottle_Acquire_QueueTimeout(t *testing.T) {
	th := NewRequestThrottle(ThrottlePolicy{MaxConcurrent: 1, QueueTimeout: 50 * time.Millisecond})
	ctx := context.Background()
	_ = th.Acquire(ctx) // fill the only slot

	if err := th.Acquire(ctx); err == nil {
		t.Fatal("expected queue timeout error")
	}
}

func TestRequestThrottle_Release_Idempotent(t *testing.T) {
	th := NewRequestThrottle(DefaultThrottlePolicy())
	// Release without Acquire should not panic
	th.Release()
	if th.Active() != 0 {
		t.Fatalf("expected active=0, got %d", th.Active())
	}
}
