package grpc

import (
	"context"
	"fmt"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// MethodInvoker handles invoking gRPC methods dynamically.
type MethodInvoker struct {
	conn *grpc.ClientConn
}

// NewMethodInvoker creates a new MethodInvoker for the given connection.
func NewMethodInvoker(conn *grpc.ClientConn) *MethodInvoker {
	return &MethodInvoker{conn: conn}
}

// InvokeUnary invokes a unary gRPC method described by the given MethodDescriptor
// with the provided request message bytes, returning the response bytes.
func (m *MethodInvoker) InvokeUnary(ctx context.Context, md *desc.MethodDescriptor, reqBytes []byte) ([]byte, error) {
	if md == nil {
		return nil, fmt.Errorf("method descriptor must not be nil")
	}

	stub := grpcdynamic.NewStub(m.conn)

	reqMsg := md.GetInputType().NewMessage()
	if err := proto.Unmarshal(reqBytes, reqMsg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %w", err)
	}

	respMsg, err := stub.InvokeRpc(ctx, md, reqMsg)
	if err != nil {
		return nil, fmt.Errorf("rpc invocation failed: %w", err)
	}

	respBytes, err := proto.Marshal(respMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	return respBytes, nil
}
