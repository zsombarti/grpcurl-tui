package grpc

import (
	"context"
	"fmt"

	"github.com/jhump/protoreflect/desc"
)

// ServiceInfo holds metadata about a discovered gRPC service.
type ServiceInfo struct {
	Name    string
	Methods []MethodInfo
}

// MethodInfo holds metadata about a single gRPC method.
type MethodInfo struct {
	Name            string
	IsClientStream  bool
	IsServerStream  bool
	Descriptor      *desc.MethodDescriptor
}

// ServiceExplorer resolves service and method metadata via reflection.
type ServiceExplorer struct {
	reflector *Reflector
}

// NewServiceExplorer creates a new ServiceExplorer.
func NewServiceExplorer(r *Reflector) *ServiceExplorer {
	return &ServiceExplorer{reflector: r}
}

// ListServices returns a slice of ServiceInfo for all services on the server.
func (s *ServiceExplorer) ListServices(ctx context.Context) ([]ServiceInfo, error) {
	names, err := s.reflector.ListServices(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	services := make([]ServiceInfo, 0, len(names))
	for _, name := range names {
		fd, err := s.reflector.ResolveService(ctx, name)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve service %q: %w", name, err)
		}

		info := ServiceInfo{Name: name}
		for _, md := range fd.GetMethods() {
			info.Methods = append(info.Methods, MethodInfo{
				Name:           md.GetName(),
				IsClientStream: md.IsClientStreaming(),
				IsServerStream: md.IsServerStreaming(),
				Descriptor:     md,
			})
		}
		services = append(services, info)
	}

	return services, nil
}
