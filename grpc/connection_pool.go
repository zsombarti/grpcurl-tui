package grpc

import (
	"fmt"
	"sync"
)

// ConnectionPool manages a pool of gRPC clients keyed by address.
type ConnectionPool struct {
	mu      sync.RWMutex
	clients map[string]*Client
	maxSize int
}

// NewConnectionPool creates a new ConnectionPool with the given maximum size.
func NewConnectionPool(maxSize int) *ConnectionPool {
	if maxSize <= 0 {
		maxSize = 10
	}
	return &ConnectionPool{
		clients: make(map[string]*Client),
		maxSize: maxSize,
	}
}

// Get returns an existing client for the given address, or creates a new one.
func (p *ConnectionPool) Get(address string) (*Client, error) {
	p.mu.RLock()
	if c, ok := p.clients[address]; ok {
		p.mu.RUnlock()
		return c, nil
	}
	p.mu.RUnlock()

	p.mu.Lock()
	defer p.mu.Unlock()

	// Double-check after acquiring write lock.
	if c, ok := p.clients[address]; ok {
		return c, nil
	}

	if len(p.clients) >= p.maxSize {
		return nil, fmt.Errorf("connection pool full (max %d)", p.maxSize)
	}

	c, err := NewClient(address)
	if err != nil {
		return nil, fmt.Errorf("failed to create client for %q: %w", address, err)
	}

	p.clients[address] = c
	return c, nil
}

// Remove closes and removes the client for the given address.
func (p *ConnectionPool) Remove(address string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if c, ok := p.clients[address]; ok {
		c.Close()
		delete(p.clients, address)
	}
}

// Len returns the number of active connections in the pool.
func (p *ConnectionPool) Len() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.clients)
}

// CloseAll closes all clients in the pool.
func (p *ConnectionPool) CloseAll() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for addr, c := range p.clients {
		c.Close()
		delete(p.clients, addr)
	}
}
