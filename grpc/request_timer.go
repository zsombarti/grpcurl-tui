package grpc

import (
	"sync"
	"time"
)

// RequestTimer tracks per-request latency measurements.
type RequestTimer struct {
	mu      sync.Mutex
	samples []time.Duration
	maxSize int
}

// NewRequestTimer creates a RequestTimer with the given sample capacity.
// If maxSize <= 0, it defaults to 100.
func NewRequestTimer(maxSize int) *RequestTimer {
	if maxSize <= 0 {
		maxSize = 100
	}
	return &RequestTimer{
		samples: make([]time.Duration, 0, maxSize),
		maxSize: maxSize,
	}
}

// Record stores a latency sample, evicting the oldest when full.
func (t *RequestTimer) Record(d time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.samples) >= t.maxSize {
		t.samples = t.samples[1:]
	}
	t.samples = append(t.samples, d)
}

// Len returns the number of recorded samples.
func (t *RequestTimer) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.samples)
}

// Average returns the mean latency across all samples.
// Returns 0 if no samples exist.
func (t *RequestTimer) Average() time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.samples) == 0 {
		return 0
	}
	var total time.Duration
	for _, s := range t.samples {
		total += s
	}
	return total / time.Duration(len(t.samples))
}

// Max returns the maximum recorded latency.
// Returns 0 if no samples exist.
func (t *RequestTimer) Max() time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.samples) == 0 {
		return 0
	}
	max := t.samples[0]
	for _, s := range t.samples[1:] {
		if s > max {
			max = s
		}
	}
	return max
}

// Clear removes all recorded samples.
func (t *RequestTimer) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.samples = t.samples[:0]
}
