package grpc

import (
	"errors"
	"sync"
	"time"
)

const defaultWatermarkMaxSize = 256

// WatermarkEntry records a high/low watermark observation for a method.
type WatermarkEntry struct {
	Method    string
	Min       time.Duration
	Max       time.Duration
	Last      time.Duration
	Count     int64
	Timestamp time.Time
}

// RequestWatermark tracks min/max latency watermarks per gRPC method.
type RequestWatermark struct {
	mu      sync.Mutex
	entries map[string]*WatermarkEntry
	max     int
}

// NewRequestWatermark creates a new RequestWatermark with an optional max bucket count.
func NewRequestWatermark(maxSize int) *RequestWatermark {
	if maxSize <= 0 {
		maxSize = defaultWatermarkMaxSize
	}
	return &RequestWatermark{
		entries: make(map[string]*WatermarkEntry, maxSize),
		max:     maxSize,
	}
}

// Record updates the watermark for the given method with the observed duration.
func (w *RequestWatermark) Record(method string, d time.Duration) error {
	if method == "" {
		return errors.New("watermark: method must not be empty")
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if e, ok := w.entries[method]; ok {
		if d < e.Min {
			e.Min = d
		}
		if d > e.Max {
			e.Max = d
		}
		e.Last = d
		e.Count++
		e.Timestamp = time.Now()
		return nil
	}
	if len(w.entries) >= w.max {
		return errors.New("watermark: max bucket count reached")
	}
	w.entries[method] = &WatermarkEntry{
		Method:    method,
		Min:       d,
		Max:       d,
		Last:      d,
		Count:     1,
		Timestamp: time.Now(),
	}
	return nil
}

// Get returns the watermark entry for the given method.
func (w *RequestWatermark) Get(method string) (*WatermarkEntry, bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	e, ok := w.entries[method]
	return e, ok
}

// Len returns the number of tracked methods.
func (w *RequestWatermark) Len() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.entries)
}

// All returns a snapshot of all watermark entries.
func (w *RequestWatermark) All() []WatermarkEntry {
	w.mu.Lock()
	defer w.mu.Unlock()
	out := make([]WatermarkEntry, 0, len(w.entries))
	for _, e := range w.entries {
		out = append(out, *e)
	}
	return out
}

// Clear removes all recorded watermarks.
func (w *RequestWatermark) Clear() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.entries = make(map[string]*WatermarkEntry, w.max)
}
