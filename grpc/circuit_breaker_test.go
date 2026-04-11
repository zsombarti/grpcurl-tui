package grpc

import (
	"testing"
	"time"
)

func TestDefaultCircuitBreakerPolicy_Values(t *testing.T) {
	p := DefaultCircuitBreakerPolicy()
	if p.MaxFailures != 5 {
		t.Fatalf("expected MaxFailures=5, got %d", p.MaxFailures)
	}
	if p.OpenDuration != 10*time.Second {
		t.Fatalf("expected OpenDuration=10s, got %v", p.OpenDuration)
	}
}

func TestNewCircuitBreaker_NotNil(t *testing.T) {
	cb := NewCircuitBreaker(DefaultCircuitBreakerPolicy())
	if cb == nil {
		t.Fatal("expected non-nil CircuitBreaker")
	}
}

func TestNewCircuitBreaker_InvalidPolicy_FallsBackToDefault(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerPolicy{MaxFailures: -1, OpenDuration: -1})
	def := DefaultCircuitBreakerPolicy()
	if cb.policy.MaxFailures != def.MaxFailures {
		t.Fatalf("expected fallback MaxFailures=%d, got %d", def.MaxFailures, cb.policy.MaxFailures)
	}
	if cb.policy.OpenDuration != def.OpenDuration {
		t.Fatalf("expected fallback OpenDuration=%v, got %v", def.OpenDuration, cb.policy.OpenDuration)
	}
}

func TestCircuitBreaker_InitialState_Closed(t *testing.T) {
	cb := NewCircuitBreaker(DefaultCircuitBreakerPolicy())
	if cb.State() != StateClosed {
		t.Fatalf("expected initial state Closed, got %v", cb.State())
	}
}

func TestCircuitBreaker_TripsOpen_AfterMaxFailures(t *testing.T) {
	p := CircuitBreakerPolicy{MaxFailures: 3, OpenDuration: 10 * time.Second}
	cb := NewCircuitBreaker(p)
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}
	if cb.State() != StateOpen {
		t.Fatalf("expected state Open after max failures, got %v", cb.State())
	}
	if err := cb.Allow(); err == nil {
		t.Fatal("expected error when circuit is open")
	}
}

func TestCircuitBreaker_RecordSuccess_ResetsClosed(t *testing.T) {
	p := CircuitBreakerPolicy{MaxFailures: 2, OpenDuration: 10 * time.Second}
	cb := NewCircuitBreaker(p)
	cb.RecordFailure()
	cb.RecordFailure()
	if cb.State() != StateOpen {
		t.Fatal("expected Open state")
	}
	cb.RecordSuccess()
	if cb.State() != StateClosed {
		t.Fatalf("expected Closed after success, got %v", cb.State())
	}
	if cb.Failures() != 0 {
		t.Fatalf("expected 0 failures after success, got %d", cb.Failures())
	}
}

func TestCircuitBreaker_HalfOpen_AfterOpenDuration(t *testing.T) {
	p := CircuitBreakerPolicy{MaxFailures: 1, OpenDuration: 10 * time.Millisecond}
	cb := NewCircuitBreaker(p)
	cb.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	if err := cb.Allow(); err != nil {
		t.Fatalf("expected Allow to succeed in half-open, got %v", err)
	}
	if cb.State() != StateHalfOpen {
		t.Fatalf("expected HalfOpen state, got %v", cb.State())
	}
}
