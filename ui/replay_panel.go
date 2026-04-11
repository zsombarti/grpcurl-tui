package ui

import (
	"context"
	"fmt"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// ReplayPanel provides a UI panel for replaying historical gRPC requests.
type ReplayPanel struct {
	frame    *tview.Frame
	list     *tview.List
	replayer *grpcpkg.RequestReplayer
	results  []grpcpkg.ReplayResult
}

// NewReplayPanel creates a new ReplayPanel.
func NewReplayPanel(replayer *grpcpkg.RequestReplayer) *ReplayPanel {
	list := tview.NewList()
	list.ShowSecondaryText(true)
	frame := tview.NewFrame(list).
		SetBorders(1, 1, 1, 1, 2, 2).
		AddText("Replay Results", true, tview.AlignLeft, 0)
	return &ReplayPanel{
		frame:    frame,
		list:     list,
		replayer: replayer,
	}
}

// Primitive returns the root tview primitive for layout embedding.
func (p *ReplayPanel) Primitive() tview.Primitive {
	return p.frame
}

// RunAll triggers replay of all history entries and refreshes the panel.
func (p *ReplayPanel) RunAll(ctx context.Context) error {
	results, err := p.replayer.ReplayAll(ctx)
	if err != nil {
		return err
	}
	p.results = results
	p.refresh()
	return nil
}

// RunAt replays the entry at the given index and appends the result.
func (p *ReplayPanel) RunAt(ctx context.Context, index int) error {
	res, err := p.replayer.ReplayAt(ctx, index)
	if err != nil {
		return err
	}
	p.results = append(p.results, *res)
	p.refresh()
	return nil
}

// Clear removes all replay results from the panel.
func (p *ReplayPanel) Clear() {
	p.results = nil
	p.list.Clear()
}

func (p *ReplayPanel) refresh() {
	p.list.Clear()
	for _, r := range p.results {
		main := fmt.Sprintf("[%d] %s → %s", r.Index, r.Address, r.Method)
		sub := fmt.Sprintf("duration: %s", r.Duration)
		if r.Err != nil {
			sub = fmt.Sprintf("error: %v", r.Err)
		}
		p.list.AddItem(main, sub, 0, nil)
	}
}
