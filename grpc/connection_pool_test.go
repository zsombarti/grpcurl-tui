package grpc

import (
	"testing"
)

func TestNewConnectionPool_Defaults(t *testing.T) {
	p := NewConnectionPool(0)
	if p == nil {
		t.Fatal("expected non-nil pool")
	}
	if p.maxSize != 10 {
		t.Errorf("expected default maxSize 10, got %d", p.maxSize)
	}
}

func TestNewConnectionPool_CustomSize(t *testing.T) {
	p := NewConnectionPool(5)
	if p.maxSize != 5 {
		t.Errorf("expected maxSize 5, got %d", p.maxSize)
	}
}

func TestConnectionPool_Len_Empty(t *testing.T) {
	p := NewConnectionPool(10)
	if p.Len() != 0 {
		t.Errorf("expected 0 connections, got %d", p.Len())
	}
}

func TestConnectionPool_Get_InvalidAddress(t *testing.T) {
	p := NewConnectionPool(10)
	_, err := p.Get("://invalid-address")
	if err == nil {
		t.Fatal("expected error for invalid address, got nil")
	}
	if p.Len() != 0 {
		t.Errorf("expected pool to remain empty after failed Get")
	}
}

func TestConnectionPool_Full(t *testing.T) {
	p := NewConnectionPool(1)
	// Manually insert a fake client to fill the pool.
	p.mu.Lock()
	p.clients["fake:1234"] = &Client{}
	p.mu.Unlock()

	_, err := p.Get("localhost:9999")
	if err == nil {
		t.Fatal("expected pool-full error, got nil")
	}
}

func TestConnectionPool_Remove_NilSafe(t *testing.T) {
	p := NewConnectionPool(10)
	// Should not panic when removing a non-existent address.
	p.Remove("nonexistent:1234")
}

func TestConnectionPool_CloseAll(t *testing.T) {
	p := NewConnectionPool(10)
	// Insert fake clients.
	p.mu.Lock()
	p.clients["a:1"] = &Client{}
	p.clients["b:2"] = &Client{}
	p.mu.Unlock()

	p.CloseAll()
	if p.Len() != 0 {
		t.Errorf("expected 0 connections after CloseAll, got %d", p.Len())
	}
}
