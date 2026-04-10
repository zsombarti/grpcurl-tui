package ui_test

import (
	"testing"

	"github.com/rivo/tview"

	"grpcurl-tui/grpc"
	"grpcurl-tui/ui"
)

func TestNewBookmarksPanel_NotNil(t *testing.T) {
	store := grpc.NewBookmarkStore(10)
	panel := ui.NewBookmarksPanel(store)
	if panel == nil {
		t.Fatal("expected non-nil BookmarksPanel")
	}
}

func TestBookmarksPanel_Primitive_NotNil(t *testing.T) {
	store := grpc.NewBookmarkStore(10)
	panel := ui.NewBookmarksPanel(store)
	prim := panel.Primitive()
	if prim == nil {
		t.Fatal("expected non-nil primitive")
	}
}

func TestBookmarksPanel_Primitive_IsFlex(t *testing.T) {
	store := grpc.NewBookmarkStore(10)
	panel := ui.NewBookmarksPanel(store)
	prim := panel.Primitive()
	if _, ok := prim.(*tview.Flex); !ok {
		t.Fatalf("expected *tview.Flex, got %T", prim)
	}
}

func TestBookmarksPanel_Refresh_Empty(t *testing.T) {
	store := grpc.NewBookmarkStore(10)
	panel := ui.NewBookmarksPanel(store)
	// Refresh on empty store should not panic
	panel.Refresh()
}

func TestBookmarksPanel_Refresh_WithEntries(t *testing.T) {
	store := grpc.NewBookmarkStore(10)
	_ = store.Add(grpc.Bookmark{
		Name:    "local",
		Address: "localhost:50051",
		Note:    "local dev server",
	})
	_ = store.Add(grpc.Bookmark{
		Name:    "staging",
		Address: "staging.example.com:443",
		Note:    "staging env",
	})
	panel := ui.NewBookmarksPanel(store)
	panel.Refresh()
	if panel.Len() != 2 {
		t.Fatalf("expected 2 entries after refresh, got %d", panel.Len())
	}
}

func TestBookmarksPanel_Refresh_Idempotent(t *testing.T) {
	store := grpc.NewBookmarkStore(10)
	_ = store.Add(grpc.Bookmark{
		Name:    "prod",
		Address: "prod.example.com:443",
	})
	panel := ui.NewBookmarksPanel(store)
	panel.Refresh()
	panel.Refresh()
	if panel.Len() != 1 {
		t.Fatalf("expected 1 entry after double refresh, got %d", panel.Len())
	}
}

func TestBookmarksPanel_SelectedAddress_Empty(t *testing.T) {
	store := grpc.NewBookmarkStore(10)
	panel := ui.NewBookmarksPanel(store)
	addr := panel.SelectedAddress()
	if addr != "" {
		t.Fatalf("expected empty address on empty panel, got %q", addr)
	}
}
