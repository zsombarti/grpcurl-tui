package grpc

import (
	"errors"
	"sync"
)

// RouteRule maps a method pattern to a target address.
type RouteRule struct {
	Method  string
	Address string
}

// RequestRouter routes gRPC method calls to different backend addresses
// based on configurable rules.
type RequestRouter struct {
	mu      sync.RWMutex
	rules   []RouteRule
	fallback string
}

// NewRequestRouter creates a new RequestRouter with an optional fallback address.
func NewRequestRouter(fallback string) *RequestRouter {
	return &RequestRouter{
		rules:    make([]RouteRule, 0),
		fallback: fallback,
	}
}

// AddRule registers a routing rule. Method is matched as a prefix or exact string.
func (r *RequestRouter) AddRule(method, address string) error {
	if method == "" {
		return errors.New("request_router: method must not be empty")
	}
	if address == "" {
		return errors.New("request_router: address must not be empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rules = append(r.rules, RouteRule{Method: method, Address: address})
	return nil
}

// Route returns the target address for the given method.
// It returns the fallback address if no rule matches.
func (r *RequestRouter) Route(method string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, rule := range r.rules {
		if rule.Method == method {
			return rule.Address
		}
	}
	return r.fallback
}

// Len returns the number of registered routing rules.
func (r *RequestRouter) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.rules)
}

// Clear removes all routing rules.
func (r *RequestRouter) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rules = r.rules[:0]
}

// SetFallback updates the fallback address.
func (r *RequestRouter) SetFallback(address string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.fallback = address
}

// Fallback returns the current fallback address.
func (r *RequestRouter) Fallback() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.fallback
}
