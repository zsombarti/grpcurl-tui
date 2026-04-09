package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

// ServiceInfo holds metadata about a discovered gRPC service.
type ServiceInfo struct {
	Name    string
	Methods []string
}

// Reflector uses gRPC server reflection to discover services and methods.
type Reflector struct {
	conn *grpc.ClientConn
}

// NewReflector creates a new Reflector for the given connection.
func NewReflector(conn *grpc.ClientConn) *Reflector {
	return &Reflector{conn: conn}
}

// ListServices returns all services exposed by the gRPC server via reflection.
func (r *Reflector) ListServices(ctx context.Context) ([]ServiceInfo, error) {
	stub := grpc_reflection_v1alpha.NewServerReflectionClient(r.conn)

	stream, err := stub.ServerReflectionInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("opening reflection stream: %w", err)
	}
	defer stream.CloseSend()

	err = stream.Send(&grpc_reflection_v1alpha.ServerReflectionRequest{
		MessageRequest: &grpc_reflection_v1alpha.ServerReflectionRequest_ListServices{
			ListServices: "",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("sending list services request: %w", err)
	}

	resp, err := stream.Recv()
	if err != nil {
		return nil, fmt.Errorf("receiving list services response: %w", err)
	}

	listResp, ok := resp.MessageResponse.(*grpc_reflection_v1alpha.ServerReflectionResponse_ListServicesResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type: %T", resp.MessageResponse)
	}

	var services []ServiceInfo
	for _, svc := range listResp.ListServicesResponse.Service {
		services = append(services, ServiceInfo{
			Name:    svc.Name,
			Methods: []string{},
		})
	}

	return services, nil
}
