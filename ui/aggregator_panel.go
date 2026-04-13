package ui

import (
	"fmt"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// AggregatorPanel displays aggregated request results in the TUI.
type AggregatorPanel struct {
	frame      *tview.Frame
	list       *tview.List
	aggregator *grpcpkg.RequestAggregator
}

// NewAggregatorPanel creates a new AggregatorPanel backed by the given aggregator.
func NewAggregatorPanel(aggregator *grpcpkg.RequestAggregator) *AggregatorPanel {
	if aggregator == nil {
		aggregator = grpcpkg.NewRequestAggregator(0)
	}
	list := tview.NewList().ShowSecondaryText(true)
	frame := tview.NewFrame(list).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText("Aggregated Results", true, tview.AlignLeft, tcell_white)
	return &AggregatorPanel{
		frame:      frame,
		list:       list,
		aggregator: aggregator,
	}
}

// Primitive returns the tview primitive for embedding in layouts.
func (p *AggregatorPanel) Primitive() tview.Primitive {
	return p.frame
}

// Refresh re-renders the list from the current aggregator state.
func (p *AggregatorPanel) Refresh() {
	p.list.Clear()
	results := p.aggregator.All()
	if len(results) == 0 {
		p.list.AddItem("No aggregated results", "", 0, nil)
		return
	}
	for _, r := range results {
		main := fmt.Sprintf("[%s] %d payload(s)", r.Method, len(r.Payloads))
		sub := fmt.Sprintf("merged keys: %d  at: %s", len(r.Merged), r.CreatedAt.Format("15:04:05"))
		p.list.AddItem(main, sub, 0, nil)
	}
}

// Clear removes all entries from the underlying aggregator and refreshes the panel.
func (p *AggregatorPanel) Clear() {
	p.aggregator.Clear()
	p.Refresh()
}
