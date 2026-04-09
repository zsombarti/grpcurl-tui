package grpc

import (
	"context"
	"fmt"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

// Reflector wraps gRPC server reflection to discover services and methods.
type Reflector struct {
	conn *grpc.ClientConn
}

// NewReflector creates a new Reflector for the given gRPC connection.
func NewReflector(conn *grpc.ClientConn) *Reflector {
	return &Reflector{conn: conn}
}

// ListServices returns the fully-qualified names of all services exposed by the server.
func (r *Reflector) ListServices(ctx context.Context) ([]string, error) {
	stub := reflectpb.NewServerReflectionClient(r.conn)
	client := grpcreflect.NewClient(ctx, stub)
	defer client.Reset()

	services, err := client.ListServices()
	if err != nil {
		return nil, fmt.Errorf("reflection list services: %w", err)
	}
	return services, nil
}

// ResolveService returns the ServiceDescriptor for the named service.
func (r *Reflector) ResolveService(ctx context.Context, serviceName string) (*desc.ServiceDescriptor, error) {
	stub := reflectpb.NewServerReflectionClient(r.conn)
	client := grpcreflect.NewClient(ctx, stub)
	defer client.Reset()

	fd, err := client.FileContainingSymbol(serviceName)
	if err != nil {
		return nil, fmt.Errorf("reflection resolve service %q: %w", serviceName, err)
	}

	sd := fd.FindSymbol(serviceName)
	if sd == nil {
		return nil, fmt.Errorf("symbol %q not found in file descriptor", serviceName)
	}

	svcDesc, ok := sd.(*desc.ServiceDescriptor)
	if !ok {
		return nil, fmt.Errorf("symbol %q is not a service", serviceName)
	}

	return svcDesc, nil
}
