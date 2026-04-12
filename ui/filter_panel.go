package ui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// FilterPanel provides a UI for building and inspecting RequestFilter rules.
type FilterPanel struct {
	frame  *tview.Frame
	list   *tview.List
	filter *grpcpkg.RequestFilter
}

// NewFilterPanel creates a new FilterPanel backed by the given RequestFilter.
func NewFilterPanel(filter *grpcpkg.RequestFilter) *FilterPanel {
	if filter == nil {
		filter = grpcpkg.NewRequestFilter()
	}
	list := tview.NewList().ShowSecondaryText(false)
	frame := tview.NewFrame(list).
		SetBorders(1, 1, 1, 1, 2, 2).
		AddText("Request Filters", true, tview.AlignCenter, tcell.ColorYellow)
	return &FilterPanel{
		frame:  frame,
		list:   list,
		filter: filter,
	}
}

// Primitive returns the root tview primitive for layout embedding.
func (p *FilterPanel) Primitive() tview.Primitive {
	return p.frame
}

// Refresh re-renders the list from the current filter rules.
func (p *FilterPanel) Refresh() {
	p.list.Clear()
	for i, r := range p.filter.Rules() {
		label := fmt.Sprintf("[%d] %s %s %q", i+1, r.Field, r.Operator, r.Value)
		p.list.AddItem(label, "", 0, nil)
	}
	if p.filter.Len() == 0 {
		p.list.AddItem("(no filters)", "", 0, nil)
	}
}

// AddRule delegates to the underlying filter and refreshes the panel.
func (p *FilterPanel) AddRule(field, operator, value string) error {
	if err := p.filter.AddRule(strings.TrimSpace(field), operator, strings.TrimSpace(value)); err != nil {
		return err
	}
	p.Refresh()
	return nil
}

// Clear removes all rules and refreshes the panel.
func (p *FilterPanel) Clear() {
	p.filter.Clear()
	p.Refresh()
}
