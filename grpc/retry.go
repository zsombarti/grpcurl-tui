package grpc

import (
	"context"
	"time"
)

// RetryPolicy defines the configuration for retry behaviour.
type RetryPolicy struct {
	MaxAttempts int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

// DefaultRetryPolicy returns a sensible default retry policy.
func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxAttempts:  3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     2 * time.Second,
		Multiplier:   2.0,
	}
}

// Retryer executes an operation with retry logic based on a RetryPolicy.
type Retryer struct {
	policy RetryPolicy
}

// NewRetryer creates a new Retryer with the given policy.
func NewRetryer(policy RetryPolicy) *Retryer {
	return &Retryer{policy: policy}
}

// Do runs fn up to MaxAttempts times, backing off between attempts.
// It returns the last error if all attempts fail.
func (r *Retryer) Do(ctx context.Context, fn func() error) error {
	delay := r.policy.InitialDelay
	var lastErr error

	for attempt := 0; attempt < r.policy.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if attempt < r.policy.MaxAttempts-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
			delay = time.Duration(float64(delay) * r.policy.Multiplier)
			if delay > r.policy.MaxDelay {
				delay = r.policy.MaxDelay
			}
		}
	}
	return lastErr
}
