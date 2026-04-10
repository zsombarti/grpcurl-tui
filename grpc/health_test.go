package grpc

import (
	"context"
	"testing"
	"time"
)

func TestNewHealthChecker_NotNil(t *testing.T) {
	hc := NewHealthChecker(3 * time.Second)
	if hc == nil {
		t.Fatal("expected non-nil HealthChecker")
	}
}

func TestNewHealthChecker_DefaultTimeout(t *testing.T) {
	hc := NewHealthChecker(0)
	if hc.timeout != 5*time.Second {
		t.Fatalf("expected default timeout 5s, got %v", hc.timeout)
	}
}

func TestNewHealthChecker_NegativeTimeout(t *testing.T) {
	hc := NewHealthChecker(-1 * time.Second)
	if hc.timeout != 5*time.Second {
		t.Fatalf("expected default timeout 5s, got %v", hc.timeout)
	}
}

func TestHealthChecker_Check_InvalidAddress(t *testing.T) {
	hc := NewHealthChecker(500 * time.Millisecond)
	status := hc.Check(context.Background(), "invalid-address-$$:9999")
	if status.Err == nil {
		t.Fatal("expected error for invalid address")
	}
	if status.Address != "invalid-address-$$:9999" {
		t.Errorf("expected address preserved, got %q", status.Address)
	}
	if status.Status != "UNAVAILABLE" {
		t.Errorf("expected UNAVAILABLE status, got %q", status.Status)
	}
}

func TestHealthChecker_Check_CancelledContext(t *testing.T) {
	hc := NewHealthChecker(5 * time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	status := hc.Check(ctx, "localhost:50051")
	if status.Err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestHealthChecker_Check_Latency_NonNegative(t *testing.T) {
	hc := NewHealthChecker(300 * time.Millisecond)
	status := hc.Check(context.Background(), "localhost:1")
	if status.Latency < 0 {
		t.Errorf("expected non-negative latency, got %v", status.Latency)
	}
}
