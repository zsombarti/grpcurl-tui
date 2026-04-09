package grpc

import (
	"fmt"
	"sync"
)

// Manager maintains a registry of named gRPC clients.
type Manager struct {
	mu      sync.RWMutex
	clients map[string]*Client
}

// NewManager creates an empty client manager.
func NewManager() *Manager {
	return &Manager{
		clients: make(map[string]*Client),
	}
}

// Connect dials the given address, stores the client under that address key,
// and returns it. If a healthy connection already exists it is returned as-is.
func (m *Manager) Connect(address string) (*Client, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if existing, ok := m.clients[address]; ok && existing.IsReady() {
		return existing, nil
	}

	client, err := NewClient(address)
	if err != nil {
		return nil, fmt.Errorf("manager: %w", err)
	}

	m.clients[address] = client
	return client, nil
}

// Get returns the client registered under address, or nil if not found.
func (m *Manager) Get(address string) *Client {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.clients[address]
}

// Addresses returns all registered target addresses.
func (m *Manager) Addresses() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	addrs := make([]string, 0, len(m.clients))
	for addr := range m.clients {
		addrs = append(addrs, addr)
	}
	return addrs
}

// Disconnect closes and removes the client for the given address.
func (m *Manager) Disconnect(address string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	client, ok := m.clients[address]
	if !ok {
		return fmt.Errorf("manager: no client for %s", address)
	}

	delete(m.clients, address)
	return client.Close()
}

// CloseAll closes every managed connection.
func (m *Manager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for addr, client := range m.clients {
		client.Close() //nolint:errcheck
		delete(m.clients, addr)
	}
}
