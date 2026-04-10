package grpc

import (
	"context"
	"testing"
	"time"
)

func TestNewReflector_NotNil(t *testing.T) {
	client, err := NewClient("localhost:50051")
	if err != nil {
		t.Fatalf("NewClient returned unexpected error: %v", err)
	}
	defer client.Close()

	reflector := NewReflector(client.conn)
	if reflector == nil {
		t.Fatal("expected non-nil Reflector")
	}
}

func TestReflector_ListServices_Timeout(t *testing.T) {
	client, err := NewClient("localhost:50051")
	if err != nil {
		t.Fatalf("NewClient returned unexpected error: %v", err)
	}
	defer client.Close()

	reflector := NewReflector(client.conn)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// No server is running, so this should fail gracefully.
	_, err = reflector.ListServices(ctx)
	if err == nil {
		t.Fatal("expected error when no server is running, got nil")
	}
}

func TestReflector_ListServices_CancelledContext(t *testing.T) {
	client, err := NewClient("localhost:50051")
	if err != nil {
		t.Fatalf("NewClient returned unexpected error: %v", err)
	}
	defer client.Close()

	reflector := NewReflector(client.conn)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	_, err = reflector.ListServices(ctx)
	if err == nil {
		t.Fatal("expected error for cancelled context, got nil")
	}
}

// newTestReflector is a helper that creates a Client and Reflector for use in
// tests, returning a cleanup function that should be deferred by the caller.
func newTestReflector(t *testing.T) (*Reflector, func()) {
	t.Helper()
	client, err := NewClient("localhost:50051")
	if err != nil {
		t.Fatalf("NewClient returned unexpected error: %v", err)
	}
	reflector := NewReflector(client.conn)
	return reflector, func() { client.Close() }
}
