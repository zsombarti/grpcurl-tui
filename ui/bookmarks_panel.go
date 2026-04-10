package ui

import (
	"fmt"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// BookmarksPanel displays and manages saved request bookmarks.
type BookmarksPanel struct {
	list  *tview.List
	store *grpcpkg.BookmarkStore
}

// NewBookmarksPanel creates a BookmarksPanel backed by the given store.
func NewBookmarksPanel(store *grpcpkg.BookmarkStore) *BookmarksPanel {
	if store == nil {
		store = grpcpkg.NewBookmarkStore(50)
	}
	list := tview.NewList()
	list.SetBorder(true).SetTitle(" Bookmarks ")
	list.ShowSecondaryText(true)
	return &BookmarksPanel{list: list, store: store}
}

// Primitive returns the underlying tview widget.
func (p *BookmarksPanel) Primitive() tview.Primitive {
	return p.list
}

// Refresh reloads the list from the bookmark store.
func (p *BookmarksPanel) Refresh() {
	p.list.Clear()
	for _, b := range p.store.All() {
		secondary := fmt.Sprintf("%s  %s", b.Address, b.Method)
		p.list.AddItem(b.Name, secondary, 0, nil)
	}
}

// AddBookmark saves a bookmark and refreshes the view.
func (p *BookmarksPanel) AddBookmark(b grpcpkg.Bookmark) error {
	if err := p.store.Add(b); err != nil {
		return err
	}
	p.Refresh()
	return nil
}

// DeleteSelected removes the currently highlighted bookmark.
func (p *BookmarksPanel) DeleteSelected() error {
	idx := p.list.GetCurrentItem()
	all := p.store.All()
	if idx < 0 || idx >= len(all) {
		return fmt.Errorf("no bookmark selected")
	}
	if err := p.store.Delete(all[idx].Name); err != nil {
		return err
	}
	p.Refresh()
	return nil
}

// SelectedBookmark returns the currently highlighted bookmark, or nil.
func (p *BookmarksPanel) SelectedBookmark() *grpcpkg.Bookmark {
	idx := p.list.GetCurrentItem()
	all := p.store.All()
	if idx < 0 || idx >= len(all) {
		return nil
	}
	b := all[idx]
	return &b
}
