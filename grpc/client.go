package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

// Client wraps a gRPC connection with connection state management.
type Client struct {
	conn    *grpc.ClientConn
	address string
}

// NewClient creates a new gRPC client connected to the given address.
func NewClient(address string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", address, err)
	}

	return &Client{
		conn:    conn,
		address: address,
	}, nil
}

// Address returns the target address of the client.
func (c *Client) Address() string {
	return c.address
}

// State returns the current connectivity state of the connection.
func (c *Client) State() connectivity.State {
	return c.conn.GetState()
}

// IsReady reports whether the connection is in the READY state.
func (c *Client) IsReady() bool {
	return c.conn.GetState() == connectivity.Ready
}

// Conn returns the underlying gRPC client connection.
func (c *Client) Conn() *grpc.ClientConn {
	return c.conn
}

// Close tears down the client connection.
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
