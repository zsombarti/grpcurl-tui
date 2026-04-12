package ui

import (
	"testing"

	grpcpkg "grpcurl-tui/grpc"
)

func newTestRouterPanel() (*RouterPanel, *grpcpkg.RequestRouter) {
	router := grpcpkg.NewRequestRouter("localhost:50051")
	panel := NewRouterPanel(router)
	return panel, router
}

func TestNewRouterPanel_NotNil(t *testing.T) {
	p, _ := newTestRouterPanel()
	if p == nil {
		t.Fatal("expected non-nil RouterPanel")
	}
}

func TestRouterPanel_Primitive_NotNil(t *testing.T) {
	p, _ := newTestRouterPanel()
	if p.Primitive() == nil {
		t.Fatal("expected non-nil Primitive")
	}
}

func TestRouterPanel_Refresh_Empty(t *testing.T) {
	p, _ := newTestRouterPanel()
	// Should not panic with zero rules.
	p.Refresh()
	if p.table.GetRowCount() < 1 {
		t.Fatal("expected at least header row")
	}
}

func TestRouterPanel_Refresh_WithRules(t *testing.T) {
	p, router := newTestRouterPanel()
	_ = router.AddRule("/svc.Foo/Bar", "host-a:1111")
	_ = router.AddRule("/svc.Foo/Baz", "host-b:2222")
	p.Refresh()
	// header row + 2 rule rows (from routeSnapshot using Len)
	if p.table.GetRowCount() < 3 {
		t.Fatalf("expected at least 3 rows, got %d", p.table.GetRowCount())
	}
}

func TestRouterPanel_Refresh_Idempotent(t *testing.T) {
	p, router := newTestRouterPanel()
	_ = router.AddRule("/a/B", "x:1")
	p.Refresh()
	count1 := p.table.GetRowCount()
	p.Refresh()
	count2 := p.table.GetRowCount()
	if count1 != count2 {
		t.Fatalf("Refresh not idempotent: %d vs %d", count1, count2)
	}
}

func TestRouterPanel_FallbackInput_DefaultValue(t *testing.T) {
	p, _ := newTestRouterPanel()
	got := p.fallbackInput.GetText()
	if got != "localhost:50051" {
		t.Fatalf("expected fallback 'localhost:50051', got %q", got)
	}
}
