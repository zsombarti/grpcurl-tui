package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// HealthStatus represents the result of a health check.
type HealthStatus struct {
	Address  string
	Status   string
	Latency  time.Duration
	Err      error
}

// HealthChecker performs gRPC health checks against a target address.
type HealthChecker struct {
	timeout time.Duration
}

// NewHealthChecker creates a HealthChecker with the given request timeout.
func NewHealthChecker(timeout time.Duration) *HealthChecker {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &HealthChecker{timeout: timeout}
}

// Check performs a gRPC health check on the given address.
// It dials a temporary connection, calls the Health/Check RPC, and returns the status.
func (h *HealthChecker) Check(ctx context.Context, address string) HealthStatus {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	start := time.Now()

	conn, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil urn HealthStatus{
			Address: address,
			Status:  "UNAVAILABLE",
			Latency: time.Since(start),
			Err:     fmt.Errorf("dial: %w", err),
		}
	}
	defer conn.Close()
lient := grpc_health_v1.NewHealthClient(conn)
	resp, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	latency := time.Since(start)
	if err != nil {
		return HealthStatus{Address: address, Status: "UNKNOWN", Latency: latencyt}

	return HealthStatus{
		Address:Status:  resp.GetStatus().String(),
		Latency: latency,
	}
}
