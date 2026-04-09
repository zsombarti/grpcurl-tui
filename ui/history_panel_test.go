package ui

import (
	"testing"

	grpcpkg "grpcurl-tui/grpc"
)

func TestNewHistoryPanel_NotNil(t *testing.T) {
	h := grpcpkg.NewHistory(10)
	p := NewHistoryPanel(h)
	if p == nil {
		t.Fatal("expected non-nil HistoryPanel")
	}
	if p.Table == nil {
		t.Fatal("expected embedded Table to be initialised")
	}
}

func TestHistoryPanel_Refresh_Empty(t *testing.T) {
	h := grpcpkg.NewHistory(10)
	p := NewHistoryPanel(h)
	// Should not panic on empty history
	p.Refresh()
	// Only the header row should be present (row 0)
	if p.GetRowCount() != 1 {
		t.Errorf("expected 1 row (header), got %d", p.GetRowCount())
	}
}

func TestHistoryPanel_Refresh_WithEntries(t *testing.T) {
	h := grpcpkg.NewHistory(10)
	h.Add(grpcpkg.HistoryEntry{Address: "localhost:50051", Service: "Greeter", Method: "SayHello"})
	h.Add(grpcpkg.HistoryEntry{Address: "localhost:50051", Service: "Greeter", Method: "SayBye", Error: "rpc error"})
	p := NewHistoryPanel(h)
	p.Refresh()
	// header + 2 data rows
	if p.GetRowCount() != 3 {
		t.Errorf("expected 3 rows, got %d", p.GetRowCount())
	}
}

func TestHistoryPanel_Refresh_Idempotent(t *testing.T) {
	h := grpcpkg.NewHistory(10)
	h.Add(grpcpkg.HistoryEntry{Method: "Ping"})
	p := NewHistoryPanel(h)
	p.Refresh()
	p.Refresh()
	if p.GetRowCount() != 2 {
		t.Errorf("expected 2 rows after double refresh, got %d", p.GetRowCount())
	}
}
