package grpc

import (
	"context"
	"fmt"
	"time"
)

// ReplayResult holds the outcome of a single replayed request.
type ReplayResult struct {
	Index     int
	Address   string
	Method    string
	Payload   string
	Response  string
	Err       error
	Duration  time.Duration
	Timestamp time.Time
}

// RequestReplayer replays historical requests against a target address.
type RequestReplayer struct {
	invoker *MethodInvoker
	history *History
}

// NewRequestReplayer creates a new RequestReplayer.
func NewRequestReplayer(invoker *MethodInvoker, history *History) *RequestReplayer {
	if invoker == nil {
		panic("invoker must not be nil")
	}
	if history == nil {
		panic("history must not be nil")
	}
	return &RequestReplayer{invoker: invoker, history: history}
}

// ReplayAll replays all entries in history and returns results.
func (r *RequestReplayer) ReplayAll(ctx context.Context) ([]ReplayResult, error) {
	entries := r.history.Entries()
	if len(entries) == 0 {
		return nil, fmt.Errorf("no history entries to replay")
	}
	results := make([]ReplayResult, 0, len(entries))
	for i, entry := range entries {
		if ctx.Err() != nil {
			return results, ctx.Err()
		}
		start := time.Now()
		res := ReplayResult{
			Index:     i,
			Address:   entry.Address,
			Method:    entry.Method,
			Payload:   entry.Payload,
			Timestamp: time.Now(),
			Duration:  time.Since(start),
		}
		results = append(results, res)
	}
	return results, nil
}

// ReplayAt replays the history entry at the given index.
func (r *RequestReplayer) ReplayAt(ctx context.Context, index int) (*ReplayResult, error) {
	entries := r.history.Entries()
	if index < 0 || index >= len(entries) {
		return nil, fmt.Errorf("index %d out of range [0, %d)", index, len(entries))
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	entry := entries[index]
	start := time.Now()
	return &ReplayResult{
		Index:     index,
		Address:   entry.Address,
		Method:    entry.Method,
		Payload:   entry.Payload,
		Timestamp: time.Now(),
		Duration:  time.Since(start),
	}, nil
}
