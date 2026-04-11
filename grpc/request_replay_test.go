package grpc

import (
	"context"
	"testing"
	"time"
)

func newTestReplayer(t *testing.T) (*RequestReplayer, *History) {
	t.Helper()
	invoker := NewMethodInvoker()
	h := NewHistory(5)
	return NewRequestReplayer(invoker, h), h
}

func TestNewRequestReplayer_NotNil(t *testing.T) {
	r, _ := newTestReplayer(t)
	if r == nil {
		t.Fatal("expected non-nil RequestReplayer")
	}
}

func TestRequestReplayer_ReplayAll_EmptyHistory(t *testing.T) {
	r, _ := newTestReplayer(t)
	_, err := r.ReplayAll(context.Background())
	if err == nil {
		t.Fatal("expected error for empty history")
	}
}

func TestRequestReplayer_ReplayAll_WithEntries(t *testing.T) {
	r, h := newTestReplayer(t)
	h.Add(HistoryEntry{Address: "localhost:50051", Method: "/svc/Method", Payload: `{}`, Timestamp: time.Now()})
	h.Add(HistoryEntry{Address: "localhost:50051", Method: "/svc/Other", Payload: `{"id":1}`, Timestamp: time.Now()})
	results, err := r.ReplayAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestRequestReplayer_ReplayAt_OutOfRange(t *testing.T) {
	r, _ := newTestReplayer(t)
	_, err := r.ReplayAt(context.Background(), 0)
	if err == nil {
		t.Fatal("expected error for out-of-range index")
	}
}

func TestRequestReplayer_ReplayAt_ValidIndex(t *testing.T) {
	r, h := newTestReplayer(t)
	h.Add(HistoryEntry{Address: "localhost:9090", Method: "/pkg.Svc/Do", Payload: `{}`, Timestamp: time.Now()})
	res, err := r.ReplayAt(context.Background(), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Method != "/pkg.Svc/Do" {
		t.Errorf("expected method '/pkg.Svc/Do', got %q", res.Method)
	}
}

func TestRequestReplayer_ReplayAt_CancelledContext(t *testing.T) {
	r, h := newTestReplayer(t)
	h.Add(HistoryEntry{Address: "localhost:9090", Method: "/pkg.Svc/Do", Payload: `{}`, Timestamp: time.Now()})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := r.ReplayAt(ctx, 0)
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}
