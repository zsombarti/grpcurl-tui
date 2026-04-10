package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// dialInsecure is a test helper that creates a non-blocking insecure client
// connection to addr. It does NOT wait for the connection to be established.
func dialInsecure(addr string) (*grpc.ClientConn, error) {
	return grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}
