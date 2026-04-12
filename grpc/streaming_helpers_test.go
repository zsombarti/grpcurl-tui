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

// mustDialInsecure is like dialInsecure but panics if the connection cannot be
// created. Intended for use in tests where a failure to dial is unrecoverable.
func mustDialInsecure(addr string) *grpc.ClientConn {
	conn, err := dialInsecure(addr)
	if err != nil {
		panic("mustDialInsecure: " + err.Error())
	}
	return conn
}
