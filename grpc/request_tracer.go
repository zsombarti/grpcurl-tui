package grpc

import (
	"sync"
	"time"
)

// TraceSpan represents a single traced request span.
type TraceSpan struct {
	TraceID   string
	Method    string
	Address   string
	StartedAt time.Time
	EndedAt   time.Time
	Duration  time.Duration
	Error     string
	Metadata  map[string]string
}

// RequestTracer records trace spans for outgoing gRPC requests.
type RequestTracer struct {
	mu      sync.RWMutex
	spans   []TraceSpan
	maxSize int
}

// NewRequestTracer creates a new RequestTracer with the given max size.
// If maxSize <= 0 it defaults to 100.
func NewRequestTracer(maxSize int) *RequestTracer {
	if maxSize <= 0 {
		maxSize = 100
	}
	return &RequestTracer{maxSize: maxSize}
}

// Start begins a new trace span and returns its index.
func (t *RequestTracer) Start(traceID, method, address string, meta map[string]string) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	span := TraceSpan{
		TraceID:   traceID,
		Method:    method,
		Address:   address,
		StartedAt: time.Now(),
		Metadata:  meta,
	}
	if len(t.spans) >= t.maxSize {
		t.spans = t.spans[1:]
	}
	t.spans = append(t.spans, span)
	return len(t.spans) - 1
}

// Finish closes the span at the given index, recording duration and optional error.
func (t *RequestTracer) Finish(index int, errMsg string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if index < 0 || index >= len(t.spans) {
		return
	}
	t.spans[index].EndedAt = time.Now()
	t.spans[index].Duration = t.spans[index].EndedAt.Sub(t.spans[index].StartedAt)
	t.spans[index].Error = errMsg
}

// Spans returns a copy of all recorded spans.
func (t *RequestTracer) Spans() []TraceSpan {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]TraceSpan, len(t.spans))
	copy(out, t.spans)
	return out
}

// Len returns the number of recorded spans.
func (t *RequestTracer) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.spans)
}

// Clear removes all recorded spans.
func (t *RequestTracer) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.spans = nil
}
