package grpc

import (
	"fmt"
	"sync"
	"time"
)

// RequestLogEntry holds a single logged request/response pair.
type RequestLogEntry struct {
	Method    string
	Payload   string
	Response  string
	Error     string
	Timestamp time.Time
	Duration  time.Duration
}

// RequestLogger records outgoing requests and their responses for display.
type RequestLogger struct {
	mu      sync.RWMutex
	entries []RequestLogEntry
	maxSize int
}

const defaultRequestLogMaxSize = 200

// NewRequestLogger creates a RequestLogger with an optional max size.
func NewRequestLogger(maxSize int) *RequestLogger {
	if maxSize <= 0 {
		maxSize = defaultRequestLogMaxSize
	}
	return &RequestLogger{maxSize: maxSize}
}

// Log appends a new entry; oldest entry is evicted when capacity is reached.
func (l *RequestLogger) Log(method, payload, response string, err error, d time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	entry := RequestLogEntry{
		Method:    method,
		Payload:   payload,
		Response:  response,
		Error:     errStr,
		Timestamp: time.Now(),
		Duration:  d,
	}
	if len(l.entries) >= l.maxSize {
		l.entries = l.entries[1:]
	}
	l.entries = append(l.entries, entry)
}

// Len returns the number of log entries.
func (l *RequestLogger) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.entries)
}

// Entries returns a copy of all log entries.
func (l *RequestLogger) Entries() []RequestLogEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]RequestLogEntry, len(l.entries))
	copy(out, l.entries)
	return out
}

// Clear removes all log entries.
func (l *RequestLogger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = nil
}

// Summary returns a one-line summary of the entry at index i.
func (l *RequestLogger) Summary(i int) (string, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if i < 0 || i >= len(l.entries) {
		return "", fmt.Errorf("request_logger: index %d out of range", i)
	}
	e := l.entries[i]
	status := "OK"
	if e.Error != "" {
		status = "ERR"
	}
	return fmt.Sprintf("[%s] %s %s (%s)", e.Timestamp.Format("15:04:05"), status, e.Method, e.Duration.Round(time.Millisecond)), nil
}
