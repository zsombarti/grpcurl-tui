package grpc

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// StreamReceiver is a callback invoked for each message received from a server stream.
type StreamReceiver func(msg proto.Message, err error)

// StreamInvoker handles server-streaming gRPC calls.
type StreamInvoker struct {
	conn *grpc.ClientConn
}

// NewStreamInvoker creates a StreamInvoker backed by the given connection.
func NewStreamInvoker(conn *grpc.ClientConn) *StreamInvoker {
	return &StreamInvoker{conn: conn}
}

// InvokeServerStream opens a server-streaming RPC and calls recv for each
// response message until the stream ends or the context is cancelled.
func (s *StreamInvoker) InvokeServerStream(
	ctx context.Context,
	fullyQualifiedMethod string,
	req proto.Message,
	receiver StreamReceiver,
) error {
	if s.conn == nil {
		return fmt.Errorf("stream invoker: nil connection")
	}
	if receiver == nil {
		return fmt.Errorf("stream invoker: nil receiver callback")
	}

	stream, err := s.conn.NewStream(ctx, &grpc.StreamDesc{
		ServerStreams: true,
	}, fullyQualifiedMethod)
	if err != nil {
		return fmt.Errorf("stream invoker: open stream: %w", err)
	}

	if err := stream.SendMsg(req); err != nil {
		return fmt.Errorf("stream invoker: send request: %w", err)
	}
	if err := stream.CloseSend(); err != nil {
		return fmt.Errorf("stream invoker: close send: %w", err)
	}

	for {
		var raw []byte
		err := stream.RecvMsg(&raw)
		if err == io.EOF {
			break
		}
		if err != nil {
			receiver(nil, fmt.Errorf("stream invoker: recv: %w", err))
			return err
		}
		receiver(req, nil) // placeholder: real impl would unmarshal into dynamic message
	}
	return nil
}
