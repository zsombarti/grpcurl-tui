package grpc

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewRetryer_NotNil(t *testing.T) {
	r := NewRetryer(DefaultRetryPolicy())
	if r == nil {
		t.Fatal("expected non-nil Retryer")
	}
}

func TestDefaultRetryPolicy_Values(t *testing.T) {
	p := DefaultRetryPolicy()
	if p.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", p.MaxAttempts)
	}
	if p.Multiplier != 2.0 {
		t.Errorf("expected Multiplier=2.0, got %f", p.Multiplier)
	}
}

func TestRetryer_Do_SuccessOnFirstAttempt(t *testing.T) {
	r := NewRetryer(DefaultRetryPolicy())
	calls := 0
	err := r.Do(context.Background(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestRetryer_Do_RetriesOnFailure(t *testing.T) {
	policy := RetryPolicy{MaxAttempts: 3, InitialDelay: time.Millisecond, MaxDelay: 10 * time.Millisecond, Multiplier: 1.5}
	r := NewRetryer(policy)
	calls := 0
	sentinel := errors.New("transient")
	err := r.Do(context.Background(), func() error {
		calls++
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestRetryer_Do_CancelledContext(t *testing.T) {
	policy := RetryPolicy{MaxAttempts: 5, InitialDelay: time.Second, MaxDelay: 5 * time.Second, Multiplier: 2.0}
	r := NewRetryer(policy)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := r.Do(ctx, func() error {
		return errors.New("should not run")
	})
	if err == nil {
		t.Fatal("expected error from cancelled context")
	}
}

func TestRetryer_Do_SuccessOnRetry(t *testing.T) {
	policy := RetryPolicy{MaxAttempts: 3, InitialDelay: time.Millisecond, MaxDelay: 10 * time.Millisecond, Multiplier: 1.0}
	r := NewRetryer(policy)
	calls := 0
	err := r.Do(context.Background(), func() error {
		calls++
		if calls < 3 {
			return errors.New("not yet")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error after retry, got %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}
