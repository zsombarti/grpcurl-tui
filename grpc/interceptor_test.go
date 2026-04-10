package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestNewCallInterceptor_NotNil(t *testing.T) {
	ci := NewCallInterceptor()
	if ci == nil {
		t.Fatal("expected non-nil CallInterceptor")
	}
}

func TestCallInterceptor_RecordsCall(t *testing.T) {
	ci := NewCallInterceptor()
	interceptor := ci.UnaryInterceptor()

	invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return nil
	}

	_ = interceptor(context.Background(), "/pkg.Service/Method", nil, nil, nil, invoker)

	records := ci.Records()
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].Method != "/pkg.Service/Method" {
		t.Errorf("unexpected method: %s", records[0].Method)
	}
	if records[0].Duration < 0 {
		t.Errorf("duration should be non-negative")
	}
	if records[0].Err != nil {
		t.Errorf("expected nil error")
	}
}

func TestCallInterceptor_RecordsError(t *testing.T) {
	ci := NewCallInterceptor()
	interceptor := ci.UnaryInterceptor()
	expectedErr := errors.New("rpc error")

	invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return expectedErr
	}

	_ = interceptor(context.Background(), "/svc/Fail", nil, nil, nil, invoker)

	records := ci.Records()
	if len(records) != 1 {
		t.Fatalf("expected 1 record")
	}
	if !errors.Is(records[0].Err, expectedErr) {
		t.Errorf("expected recorded error to match")
	}
}

func TestCallInterceptor_Clear(t *testing.T) {
	ci := NewCallInterceptor()
	interceptor := ci.UnaryInterceptor()
	invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return nil
	}
	_ = interceptor(context.Background(), "/svc/M", nil, nil, nil, invoker)
	ci.Clear()
	if len(ci.Records()) != 0 {
		t.Error("expected records to be cleared")
	}
}

func TestCallInterceptor_Duration_Positive(t *testing.T) {
	ci := NewCallInterceptor()
	interceptor := ci.UnaryInterceptor()
	invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		time.Sleep(1 * time.Millisecond)
		return nil
	}
	_ = interceptor(context.Background(), "/svc/Slow", nil, nil, nil, invoker)
	records := ci.Records()
	if records[0].Duration <= 0 {
		t.Error("expected positive duration")
	}
}

func TestNewMetadataInterceptor_InjectsMetadata(t *testing.T) {
	mi := NewMetadataInterceptor(map[string]string{"x-request-id": "abc123"})
	interceptor := mi.UnaryInterceptor()

	var capturedCtx context.Context
	invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		capturedCtx = ctx
		return nil
	}

	_ = interceptor(context.Background(), "/svc/M", nil, nil, nil, invoker)

	md, ok := metadata.FromOutgoingContext(capturedCtx)
	if !ok {
		t.Fatal("expected outgoing metadata")
	}
	vals := md.Get("x-request-id")
	if len(vals) == 0 || vals[0] != "abc123" {
		t.Errorf("expected x-request-id=abc123, got %v", vals)
	}
}
