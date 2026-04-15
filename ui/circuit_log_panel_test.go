package ui

import (
	"testing"

	"github.com/rivo/tview"

	"grpcurl-tui/grpc"
)

func newTestCircuitLogPanel() *CircuitLogPanel {
	return NewCircuitLogPanel(grpc.NewCircuitLog(10))
}

func TestNewCircuitLogPanel_NotNil(t *testing.T) {
	p := newTestCircuitLogPanel()
	if p == nil {
		t.Fatal("expected non-nil CircuitLogPanel")
	}
}

func TestCircuitLogPanel_Primitive_NotNil(t *testing.T) {
	p := newTestCircuitLogPanel()
	if p.Primitive() == nil {
		t.Fatal("expected non-nil primitive")
	}
}

func TestCircuitLogPanel_Primitive_IsFrame(t *testing.T) {
	p := newTestCircuitLogPanel()
	if _, ok := p.Primitive().(*tview.Frame); !ok {
		t.Fatal("expected primitive to be *tview.Frame")
	}
}

func TestCircuitLogPanel_NilStore_FallsBackToDefault(t *testing.T) {
	p := NewCircuitLogPanel(nil)
	if p == nil {
		t.Fatal("expected non-nil panel even with nil store")
	}
	if p.store == nil {
		t.Fatal("expected store to be initialised")
	}
}

func TestCircuitLogPanel_Refresh_Empty(t *testing.T) {
	p := newTestCircuitLogPanel()
	p.Refresh()
	if p.list.GetItemCount() != 1 {
		t.Fatalf("expected 1 placeholder item, got %d", p.list.GetItemCount())
	}
}

func TestCircuitLogPanel_Refresh_WithEntries(t *testing.T) {
	store := grpc.NewCircuitLog(10)
	store.Record("helloworld.Greeter/SayHello", "open", "error threshold")
	store.Record("helloworld.Greeter/SayHello", "closed", "recovery")
	p := NewCircuitLogPanel(store)
	p.Refresh()
	if p.list.GetItemCount() != 2 {
		t.Fatalf("expected 2 items, got %d", p.list.GetItemCount())
	}
}

func TestCircuitLogPanel_Clear_ResetsView(t *testing.T) {
	store := grpc.NewCircuitLog(10)
	store.Record("m", "open", "r")
	p := NewCircuitLogPanel(store)
	p.Clear()
	if p.list.GetItemCount() != 1 {
		t.Fatalf("expected 1 placeholder after clear, got %d", p.list.GetItemCount())
	}
}
