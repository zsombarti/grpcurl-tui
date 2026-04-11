package ui

import (
	"testing"

	grpcpkg "grpcurl-tui/grpc"
)

func newTestSnapshotPanel() *SnapshotPanel {
	return NewSnapshotPanel(grpcpkg.NewSnapshotStore(10))
}

func TestNewSnapshotPanel_NotNil(t *testing.T) {
	p := newTestSnapshotPanel()
	if p == nil {
		t.Fatal("expected non-nil SnapshotPanel")
	}
}

func TestSnapshotPanel_Primitive_NotNil(t *testing.T) {
	p := newTestSnapshotPanel()
	if p.Primitive() == nil {
		t.Fatal("expected non-nil primitive")
	}
}

func TestSnapshotPanel_SelectedIndex_EmptyStore(t *testing.T) {
	p := newTestSnapshotPanel()
	if p.SelectedIndex() != -1 {
		t.Fatalf("expected -1 for empty store, got %d", p.SelectedIndex())
	}
}

func TestSnapshotPanel_Refresh_Empty(t *testing.T) {
	p := newTestSnapshotPanel()
	p.Refresh()
	if p.list.GetItemCount() != 0 {
		t.Fatalf("expected 0 items, got %d", p.list.GetItemCount())
	}
}

func TestSnapshotPanel_Refresh_WithEntries(t *testing.T) {
	p := newTestSnapshotPanel()
	p.store.Save(grpcpkg.Snapshot{Name: "snap1", Address: "localhost:50051", Method: "Ping"})
	p.store.Save(grpcpkg.Snapshot{Name: "snap2", Address: "localhost:50051", Method: "Pong"})
	p.Refresh()
	if p.list.GetItemCount() != 2 {
		t.Fatalf("expected 2 items, got %d", p.list.GetItemCount())
	}
}

func TestSnapshotPanel_Clear_ResetsStore(t *testing.T) {
	p := newTestSnapshotPanel()
	p.store.Save(grpcpkg.Snapshot{Name: "x"})
	p.Refresh()
	p.Clear()
	if p.store.Len() != 0 {
		t.Fatalf("expected store len 0, got %d", p.store.Len())
	}
	if p.list.GetItemCount() != 0 {
		t.Fatalf("expected list count 0, got %d", p.list.GetItemCount())
	}
}
