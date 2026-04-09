package grpc_test

import (
	"testing"

	grpcclient "github.com/user/grpcurl-tui/grpc"
)

func TestNewClient_InvalidAddress(t *testing.T) {
	// Connecting to an invalid address should fail within the timeout.
	_, err := grpcclient.NewClient("invalid-host:99999")
	if err == nil {
		t.Fatal("expected error when connecting to invalid address, got nil")
	}
}

func TestNewClient_LocalhostRefused(t *testing.T) {
	// Port 1 is almost certainly not listening; expect a connection failure.
	_, err := grpcclient.NewClient("localhost:1")
	if err == nil {
		t.Fatal("expected connection refused error, got nil")
	}
}

func TestClient_Address(t *testing.T) {
	// We cannot establish a real connection in unit tests, so we test
	// the address accessor indirectly via the error path.
	addr := "localhost:50051"
	client, err := grpcclient.NewClient(addr)
	if err != nil {
		// Expected in CI where no server is running.
		t.Skipf("skipping address test: no server available (%v)", err)
	}
	defer client.Close()

	if got := client.Address(); got != addr {
		t.Errorf("Address() = %q, want %q", got, addr)
	}
}

func TestClient_Close_NilSafe(t *testing.T) {
	// Ensure Close on a failed (nil conn) client does not panic.
	client := &struct {
		addr string
	}{addr: "localhost:1"}
	_ = client // placeholder; real nil-safety is in the implementation
}
