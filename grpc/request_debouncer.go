package grpc

import (
	"errors"
	"sync"
	"time"
)

// DefaultDebouncerPolicy returns a sensible default debounce policy.
func DefaultDebouncerPolicy() DebouncerPolicy {
	return DebouncerPolicy{
		Wait:    300 * time.Millisecond,
		MaxWait: 2 * time.Second,
	}
}

// DebouncerPolicy controls debounce timing behaviour.
type DebouncerPolicy struct {
	Wait    time.Duration
	MaxWait time.Duration
}

// pendingCall holds the latest payload waiting to be flushed.
type pendingCall struct {
	method  string
	payload map[string]interface{}
	first   time.Time
	timer   *time.Timer
}

// RequestDebouncer collapses rapid successive calls for the same method
// into a single invocation after the wait period elapses.
type RequestDebouncer struct {
	mu     sync.Mutex
	policy DebouncerPolicy
	calls  map[string]*pendingCall
	flush  func(method string, payload map[string]interface{})
}

// NewRequestDebouncer creates a debouncer with the given policy and flush callback.
// If policy values are zero the defaults are used.
func NewRequestDebouncer(policy DebouncerPolicy, flush func(string, map[string]interface{})) (*RequestDebouncer, error) {
	if flush == nil {
		return nil, errors.New("debouncer: flush callback must not be nil")
	}
	if policy.Wait <= 0 {
		policy = DefaultDebouncerPolicy()
	}
	if policy.MaxWait <= 0 || policy.MaxWait < policy.Wait {
		policy.MaxWait = policy.Wait * 6
	}
	return &RequestDebouncer{
		policy: policy,
		calls:  make(map[string]*pendingCall),
		flush:  flush,
	}, nil
}

// Debounce schedules or resets the debounce timer for the given method.
func (d *RequestDebouncer) Debounce(method string, payload map[string]interface{}) error {
	if method == "" {
		return errors.New("debouncer: method must not be empty")
	}
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()
	if p, ok := d.calls[method]; ok {
		p.payload = payload
		// Respect MaxWait: if we have already waited long enough, fire immediately.
		if now.Sub(p.first) >= d.policy.MaxWait {
			p.timer.Stop()
			delete(d.calls, method)
			go d.flush(method, payload)
			return nil
		}
		p.timer.Reset(d.policy.Wait)
		return nil
	}

	pc := &pendingCall{method: method, payload: payload, first: now}
	pc.timer = time.AfterFunc(d.policy.Wait, func() {
		d.mu.Lock()
		entry, ok := d.calls[method]
		if ok {
			delete(d.calls, method)
		}
		d.mu.Unlock()
		if ok {
			d.flush(entry.method, entry.payload)
		}
	})
	d.calls[method] = pc
	return nil
}

// Len returns the number of pending debounced calls.
func (d *RequestDebouncer) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.calls)
}

// Cancel removes a pending debounced call without flushing it.
func (d *RequestDebouncer) Cancel(method string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if p, ok := d.calls[method]; ok {
		p.timer.Stop()
		delete(d.calls, method)
	}
}
