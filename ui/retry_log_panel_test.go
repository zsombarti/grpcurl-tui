package ui

import (
	"errors"
	"testing"

	"grpcurl-tui/grpc"
)

func newTestRetryLogPanel() *RetryLogPanel {
	return NewRetryLogPanel(grpc.NewRetryLog(10))
}

func TestNewRetryLogPanel_NotNil(t *testing.T) {
	if newTestRetryLogPanel() == nil {
		t.Fatal("expected non-nil RetryLogPanel")
	}
}

func TestRetryLogPanel_Primitive_NotNil(t *testing.T) {
	p := newTestRetryLogPanel()
	if p.Primitive() == nil {
		t.Fatal("expected non-nil primitive")
	}
}

func TestRetryLogPanel_Primitive_IsFrame(t *testing.T) {
	p := newTestRetryLogPanel()
	if _, ok := p.Primitive().(*tview.Frame); !ok {
		t.Fatal("expected *tview.Frame")
	}
}

func TestRetryLogPanel_NilStore_FallsBackToDefault(t *testing.T) {
	p := NewRetryLogPanel(nil)
	if p.store == nil {
		t.Fatal("expected fallback store")
	}
}

func TestRetryLogPanel_Refresh_Empty(t *testing.T) {
	p := newTestRetryLogPanel()
	p.Refresh() // must not panic
	if p.list.GetItemCount() != 0 {
		t.Fatalf("expected 0 items, got %d", p.list.GetItemCount())
	}
}

func TestRetryLogPanel_Refresh_WithEntries(t *testing.T) {
	store := grpc.NewRetryLog(10)
	_ = store.Record("SayHello", 1, errors.New("unavailable"))
	_ = store.Record("SayHello", 2, nil)
	p := NewRetryLogPanel(store)
	p.Refresh()
	if p.list.GetItemCount() != 2 {
		t.Fatalf("expected 2 items, got %d", p.list.GetItemCount())
	}
}

func TestRetryLogPanel_Clear_ResetsView(t *testing.T) {
	store := grpc.NewRetryLog(10)
	_ = store.Record("M", 1, nil)
	p := NewRetryLogPanel(store)
	p.Refresh()
	p.Clear()
	if p.list.GetItemCount() != 0 {
		t.Fatalf("expected 0 after clear, got %d", p.list.GetItemCount())
	}
	if store.Len() != 0 {
		t.Fatalf("expected store empty after clear, got %d", store.Len())
	}
}
