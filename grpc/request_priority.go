package grpc

import (
	"errors"
	"sort"
	"sync"
)

// PriorityLevel represents the urgency of a request.
type PriorityLevel int

const (
	PriorityLow    PriorityLevel = 0
	PriorityNormal PriorityLevel = 1
	PriorityHigh   PriorityLevel = 2
)

// PriorityEntry holds a queued request payload with its priority.
type PriorityEntry struct {
	Payload  string
	Priority PriorityLevel
	Label    string
}

// RequestPriorityQueue orders pending requests by priority level.
type RequestPriorityQueue struct {
	mu      sync.Mutex
	entries []PriorityEntry
	maxSize int
}

// NewRequestPriorityQueue creates a priority queue with the given capacity.
func NewRequestPriorityQueue(maxSize int) *RequestPriorityQueue {
	if maxSize <= 0 {
		maxSize = 64
	}
	return &RequestPriorityQueue{maxSize: maxSize}
}

// Enqueue adds an entry to the queue, returning an error if full.
func (q *RequestPriorityQueue) Enqueue(entry PriorityEntry) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.entries) >= q.maxSize {
		return errors.New("priority queue is full")
	}
	q.entries = append(q.entries, entry)
	sort.SliceStable(q.entries, func(i, j int) bool {
		return q.entries[i].Priority > q.entries[j].Priority
	})
	return nil
}

// Dequeue removes and returns the highest-priority entry.
func (q *RequestPriorityQueue) Dequeue() (PriorityEntry, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.entries) == 0 {
		return PriorityEntry{}, false
	}
	e := q.entries[0]
	q.entries = q.entries[1:]
	return e, true
}

// Len returns the current number of queued entries.
func (q *RequestPriorityQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.entries)
}

// Clear removes all entries from the queue.
func (q *RequestPriorityQueue) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.entries = nil
}
