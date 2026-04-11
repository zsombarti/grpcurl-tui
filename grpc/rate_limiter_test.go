package grpc

import (
	"context"
	"testing"
	"time"
)

func TestNewRateLimiter_NotNil(t *testing.T) {
	rl := NewRateLimiter(DefaultRateLimitPolicy())
	if rl == nil {
		t.Fatal("expected non-nil RateLimiter")
	}
}

func TestDefaultRateLimitPolicy_Values(t *testing.T) {
	p := DefaultRateLimitPolicy()
	if p.RequestsPerSecond != 10.0 {
		t.Errorf("expected 10.0 rps, got %f", p.RequestsPerSecond)
	}
	if p.Burst != 5 {
		t.Errorf("expected burst 5, got %d", p.Burst)
	}
}

func TestNewRateLimiter_InvalidPolicy_FallsBackToDefault(t *testing.T) {
	rl := NewRateLimiter(RateLimitPolicy{RequestsPerSecond: -1, Burst: 0})
	policy := rl.Policy()
	if policy.RequestsPerSecond != 10.0 {
		t.Errorf("expected fallback rps 10.0, got %f", policy.RequestsPerSecond)
	}
	if policy.Burst != 5 {
		t.Errorf("expected fallback burst 5, got %d", policy.Burst)
	}
}

func TestRateLimiter_Wait_AllowsImmediately(t *testing.T) {
	rl := NewRateLimiter(RateLimitPolicy{RequestsPerSecond: 100, Burst: 10})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := rl.Wait(ctx); err != nil {
		t.Fatalf("expected no error on first wait, got: %v", err)
	}
}

func TestRateLimiter_Wait_CancelledContext(t *testing.T) {
	rl := NewRateLimiter(RateLimitPolicy{RequestsPerSecond: 0.001, Burst: 1})
	// Drain the single burst token first.
	ctx := context.Background()
	_ = rl.Wait(ctx)

	cancelled, cancel := context.WithCancel(context.Background())
	cancel()
	err := rl.Wait(cancelled)
	if err == nil {
		t.Fatal("expected error on cancelled context")
	}
}

func TestRateLimiter_Policy_RoundTrip(t *testing.T) {
	policy := RateLimitPolicy{RequestsPerSecond: 5.0, Burst: 3}
	rl := NewRateLimiter(policy)
	got := rl.Policy()
	if got.RequestsPerSecond != 5.0 || got.Burst != 3 {
		t.Errorf("policy round-trip failed: got %+v", got)
	}
}
