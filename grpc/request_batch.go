package grpc

import (
	"context"
	"errors"
	"sync"
	"time"
)

// BatchPolicy defines configuration for batching requests.
type BatchPolicy struct {
	MaxSize     int
	MaxWait     time.Duration
	Concurrency int
}

// DefaultBatchPolicy returns sensible defaults for request batching.
func DefaultBatchPolicy() BatchPolicy {
	return BatchPolicy{
		MaxSize:     20,
		MaxWait:     200 * time.Millisecond,
		Concurrency: 4,
	}
}

// BatchResult holds the outcome of a single batched request.
type BatchResult struct {
	Index   int
	Payload string
	Err     error
}

// RequestBatcher groups multiple payloads and executes them concurrently.
type RequestBatcher struct {
	mu     sync.Mutex
	policy BatchPolicy
	queue  []string
}

// NewRequestBatcher creates a RequestBatcher with the given policy.
// Falls back to defaults if the policy is invalid.
func NewRequestBatcher(p BatchPolicy) *RequestBatcher {
	def := DefaultBatchPolicy()
	if p.MaxSize <= 0 {
		p.MaxSize = def.MaxSize
	}
	if p.MaxWait <= 0 {
		p.MaxWait = def.MaxWait
	}
	if p.Concurrency <= 0 {
		p.Concurrency = def.Concurrency
	}
	return &RequestBatcher{policy: p}
}

// Add enqueues a payload for the next batch run.
func (b *RequestBatcher) Add(payload string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.queue) >= b.policy.MaxSize {
		return errors.New("batch queue full")
	}
	b.queue = append(b.queue, payload)
	return nil
}

// Len returns the current number of queued payloads.
func (b *RequestBatcher) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.queue)
}

// Flush drains the queue and invokes fn for each payload concurrently.
// Results are collected and returned in index order.
func (b *RequestBatcher) Flush(ctx context.Context, fn func(ctx context.Context, payload string) error) []BatchResult {
	b.mu.Lock()
	items := make([]string, len(b.queue))
	copy(items, b.queue)
	b.queue = b.queue[:0]
	b.mu.Unlock()

	results := make([]BatchResult, len(items))
	sem := make(chan struct{}, b.policy.Concurrency)
	var wg sync.WaitGroup

	for i, payload := range items {
		wg.Add(1)
		go func(idx int, p string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			err := fn(ctx, p)
			results[idx] = BatchResult{Index: idx, Payload: p, Err: err}
		}(i, payload)
	}
	wg.Wait()
	return results
}

// Policy returns the active BatchPolicy.
func (b *RequestBatcher) Policy() BatchPolicy {
	return b.policy
}
