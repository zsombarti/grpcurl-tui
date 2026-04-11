package ui

import (
	"fmt"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// SnapshotPanel displays and manages saved request snapshots.
type SnapshotPanel struct {
	store *grpcpkg.SnapshotStore
	list  *tview.List
	flex  *tview.Flex
}

// NewSnapshotPanel constructs a SnapshotPanel backed by the given store.
func NewSnapshotPanel(store *grpcpkg.SnapshotStore) *SnapshotPanel {
	list := tview.NewList().ShowSecondaryText(true)
	list.SetBorder(true).SetTitle(" Snapshots ")

	flex := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(list, 0, 1, true)

	return &SnapshotPanel{store: store, list: list, flex: flex}
}

// Primitive returns the root tview primitive for layout embedding.
func (p *SnapshotPanel) Primitive() tview.Primitive { return p.flex }

// Refresh re-renders the list from the current store contents.
func (p *SnapshotPanel) Refresh() {
	p.list.Clear()
	for _, snap := range p.store.All() {
		name := snap.Name
		secondary := fmt.Sprintf("%s → %s  [%s]", snap.Address, snap.Method, snap.CreatedAt.Format("15:04:05"))
		p.list.AddItem(name, secondary, 0, nil)
	}
}

// SelectedIndex returns the currently highlighted index, or -1 if empty.
func (p *SnapshotPanel) SelectedIndex() int {
	if p.store.Len() == 0 {
		return -1
	}
	return p.list.GetCurrentItem()
}

// Clear wipes the store and refreshes the list.
func (p *SnapshotPanel) Clear() {
	p.store.Clear()
	p.Refresh()
}
