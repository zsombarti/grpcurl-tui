package grpc

import (
	"context"
	"errors"
	"sync"
	"time"
)

// RateLimitPolicy holds configuration for rate limiting.
type RateLimitPolicy struct {
	RequestsPerSecond float64
	Burst             int
}

// DefaultRateLimitPolicy returns a sensible default rate limit policy.
func DefaultRateLimitPolicy() RateLimitPolicy {
	return RateLimitPolicy{
		RequestsPerSecond: 10.0,
		Burst:             5,
	}
}

// RateLimiter controls the rate of outgoing gRPC requests.
type RateLimiter struct {
	mu       sync.Mutex
	policy   RateLimitPolicy
	tokens   float64
	lastTick time.Time
}

// NewRateLimiter creates a RateLimiter with the given policy.
// If rps <= 0 or burst <= 0, the default policy is used.
func NewRateLimiter(policy RateLimitPolicy) *RateLimiter {
	if policy.RequestsPerSecond <= 0 || policy.Burst <= 0 {
		policy = DefaultRateLimitPolicy()
	}
	return &RateLimiter{
		policy:   policy,
		tokens:   float64(policy.Burst),
		lastTick: time.Now(),
	}
}

// Wait blocks until a token is available or the context is cancelled.
func (r *RateLimiter) Wait(ctx context.Context) error {
	for {
		if err := ctx.Err(); err != nil {
			return errors.New("rate limiter: context cancelled")
		}
		r.mu.Lock()
		now := time.Now()
		elapsed := now.Sub(r.lastTick).Seconds()
		r.tokens += elapsed * r.policy.RequestsPerSecond
		if r.tokens > float64(r.policy.Burst) {
			r.tokens = float64(r.policy.Burst)
		}
		r.lastTick = now
		if r.tokens >= 1.0 {
			r.tokens -= 1.0
			r.mu.Unlock()
			return nil
		}
		r.mu.Unlock()
		select {
		case <-ctx.Done():
			return errors.New("rate limiter: context cancelled")
		case <-time.After(time.Duration(float64(time.Second) / r.policy.RequestsPerSecond)):
		}
	}
}

// Policy returns the current rate limit policy.
func (r *RateLimiter) Policy() RateLimitPolicy {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.policy
}
