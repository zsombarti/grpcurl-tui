package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// CallRecord captures metadata about a single gRPC call.
type CallRecord struct {
	Method   string
	Duration time.Duration
	Err      error
}

// CallInterceptor is a unary client interceptor that records call metadata.
type CallInterceptor struct {
	records []CallRecord
}

// NewCallInterceptor creates a new CallInterceptor.
func NewCallInterceptor() *CallInterceptor {
	return &CallInterceptor{}
}

// UnaryInterceptor returns a grpc.UnaryClientInterceptor that records each call.
func (ci *CallInterceptor) UnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		ci.records = append(ci.records, CallRecord{
			Method:   method,
			Duration: time.Since(start),
			Err:      err,
		})
		return err
	}
}

// MetadataInterceptor injects outgoing metadata into every call.
type MetadataInterceptor struct {
	pairs map[string]string
}

// NewMetadataInterceptor creates a MetadataInterceptor with the given key-value pairs.
func NewMetadataInterceptor(pairs map[string]string) *MetadataInterceptor {
	return &MetadataInterceptor{pairs: pairs}
}

// UnaryInterceptor returns a grpc.UnaryClientInterceptor that appends metadata.
func (mi *MetadataInterceptor) UnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		md := metadata.New(mi.pairs)
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// Records returns a copy of all recorded call records.
func (ci *CallInterceptor) Records() []CallRecord {
	out := make([]CallRecord, len(ci.records))
	copy(out, ci.records)
	return out
}

// Clear resets the recorded call history.
func (ci *CallInterceptor) Clear() {
	ci.records = nil
}
