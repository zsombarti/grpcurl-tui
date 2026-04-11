package grpc

import (
	"strings"
	"sync"
)

// ProtoSearchResult holds a matched service or method name.
type ProtoSearchResult struct {
	Service string
	Method  string
	Full    string
}

// ProtoSearcher filters services and methods by a query string.
type ProtoSearcher struct {
	mu      sync.RWMutex
	entries []ProtoSearchResult
}

// NewProtoSearcher creates a new ProtoSearcher with no entries.
func NewProtoSearcher() *ProtoSearcher {
	return &ProtoSearcher{
		entries: make([]ProtoSearchResult, 0),
	}
}

// Index populates the searcher with service/method pairs.
func (p *ProtoSearcher) Index(services map[string][]string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.entries = p.entries[:0]
	for svc, methods := range services {
		if len(methods) == 0 {
			p.entries = append(p.entries, ProtoSearchResult{
				Service: svc,
				Full:    svc,
			})
			continue
		}
		for _, m := range methods {
			p.entries = append(p.entries, ProtoSearchResult{
				Service: svc,
				Method:  m,
				Full:    svc + "/" + m,
			})
		}
	}
}

// Search returns all entries whose Full name contains query (case-insensitive).
func (p *ProtoSearcher) Search(query string) []ProtoSearchResult {
	p.mu.RLock()
	defer p.mu.RUnlock()
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		result := make([]ProtoSearchResult, len(p.entries))
		copy(result, p.entries)
		return result
	}
	var out []ProtoSearchResult
	for _, e := range p.entries {
		if strings.Contains(strings.ToLower(e.Full), q) {
			out = append(out, e)
		}
	}
	return out
}

// Len returns the total number of indexed entries.
func (p *ProtoSearcher) Len() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.entries)
}
