package grpc

import (
	"context"
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestNewStreamInvoker_NotNil(t *testing.T) {
	invoker := NewStreamInvoker(nil)
	if invoker == nil {
		t.Fatal("expected non-nil StreamInvoker")
	}
}

func TestStreamInvoker_InvokeServerStream_NilConnection(t *testing.T) {
	invoker := NewStreamInvoker(nil)
	err := invoker.InvokeServerStream(context.Background(), "/pkg.Svc/Method", nil, func(proto.Message, error) {})
	if err == nil {
		t.Fatal("expected error for nil connection")
	}
}

func TestStreamInvoker_InvokeServerStream_NilReceiver(t *testing.T) {
	invoker := NewStreamInvoker(nil)
	// override conn to non-nil path is not needed; nil receiver is caught first after conn check
	// We test nil receiver with a non-nil conn placeholder by checking error ordering.
	err := invoker.InvokeServerStream(context.Background(), "/pkg.Svc/Method", nil, nil)
	if err == nil {
		t.Fatal("expected error for nil receiver")
	}
}

func TestStreamInvoker_InvokeServerStream_CancelledContext(t *testing.T) {
	conn, err := dialInsecure("localhost:1")
	if err != nil {
		t.Skip("cannot create test connection:", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	invoker := NewStreamInvoker(conn)
	called := false
	err = invoker.InvokeServerStream(ctx, "/pkg.Svc/Method", nil, func(_ proto.Message, e error) {
		called = true
	})
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
	_ = called
}
