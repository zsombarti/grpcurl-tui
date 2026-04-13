package ui

import (
	"testing"

	grpcpkg "grpcurl-tui/grpc"
)

func newTestAggregatorPanel() *AggregatorPanel {
	return NewAggregatorPanel(grpcpkg.NewRequestAggregator(10))
}

func TestNewAggregatorPanel_NotNil(t *testing.T) {
	p := newTestAggregatorPanel()
	if p == nil {
		t.Fatal("expected non-nil panel")
	}
}

func TestAggregatorPanel_Primitive_NotNil(t *testing.T) {
	p := newTestAggregatorPanel()
	if p.Primitive() == nil {
		t.Fatal("expected non-nil primitive")
	}
}

func TestAggregatorPanel_Primitive_IsFrame(t *testing.T) {
	p := newTestAggregatorPanel()
	if _, ok := p.Primitive().(*tview.Frame); !ok {
		t.Fatal("expected *tview.Frame")
	}
}

func TestAggregatorPanel_NilAggregator_FallsBackToDefault(t *testing.T) {
	p := NewAggregatorPanel(nil)
	if p.aggregator == nil {
		t.Fatal("expected fallback aggregator")
	}
}

func TestAggregatorPanel_Refresh_Empty(t *testing.T) {
	p := newTestAggregatorPanel()
	p.Refresh() // should not panic
	if p.list.GetItemCount() != 1 {
		t.Fatalf("expected 1 placeholder item, got %d", p.list.GetItemCount())
	}
}

func TestAggregatorPanel_Refresh_WithEntries(t *testing.T) {
	a := grpcpkg.NewRequestAggregator(10)
	a.Aggregate("pkg.Service/Method", []map[string]interface{}{{"field": "value"}})
	a.Aggregate("pkg.Service/Other", []map[string]interface{}{{"a": 1}, {"b": 2}})
	p := NewAggregatorPanel(a)
	p.Refresh()
	if p.list.GetItemCount() != 2 {
		t.Fatalf("expected 2 items, got %d", p.list.GetItemCount())
	}
}

func TestAggregatorPanel_Clear_ResetsPanel(t *testing.T) {
	a := grpcpkg.NewRequestAggregator(10)
	a.Aggregate("M", []map[string]interface{}{{"x": 1}})
	p := NewAggregatorPanel(a)
	p.Refresh()
	p.Clear()
	if p.list.GetItemCount() != 1 {
		t.Fatalf("expected 1 placeholder after clear, got %d", p.list.GetItemCount())
	}
	if a.Len() != 0 {
		t.Fatalf("expected aggregator empty after clear, got %d", a.Len())
	}
}
