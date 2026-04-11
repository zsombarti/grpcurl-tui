package ui

import (
	"fmt"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// ProtoFilePanel displays loaded proto file descriptors and allows clearing them.
type ProtoFilePanel struct {
	flex   *tview.Flex
	list   *tview.List
	status *tview.TextView
	loader *grpcpkg.ProtoFileLoader
}

// NewProtoFilePanel creates and returns a new ProtoFilePanel.
func NewProtoFilePanel(loader *grpcpkg.ProtoFileLoader) *ProtoFilePanel {
	list := tview.NewList().ShowSecondaryText(false)
	list.SetBorder(true).SetTitle(" Proto Files ")

	status := tview.NewTextView()
	status.SetBorder(false)
	status.SetText("No files loaded.")

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(list, 0, 1, true).
		AddItem(status, 1, 0, false)

	p := &ProtoFilePanel{
		flex:   flex,
		list:   list,
		status: status,
		loader: loader,
	}
	p.Refresh()
	return p
}

// Primitive returns the root tview primitive for layout embedding.
func (p *ProtoFilePanel) Primitive() tview.Primitive {
	return p.flex
}

// Refresh redraws the list from the current loader state.
func (p *ProtoFilePanel) Refresh() {
	p.list.Clear()
	names := p.loader.Names()
	for _, name := range names {
		p.list.AddItem(name, "", 0, nil)
	}
	if len(names) == 0 {
		p.status.SetText("No files loaded.")
	} else {
		p.status.SetText(fmt.Sprintf("%d file(s) loaded.", len(names)))
	}
}

// Clear removes all loaded descriptors and refreshes the panel.
func (p *ProtoFilePanel) Clear() {
	p.loader.Clear()
	p.Refresh()
}
