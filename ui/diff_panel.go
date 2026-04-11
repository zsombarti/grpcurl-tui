package ui

import (
	"fmt"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// DiffPanel displays a side-by-side diff of two request payloads.
type DiffPanel struct {
	frame  *tview.Frame
	text   *tview.TextView
	differ *grpcpkg.RequestDiffer
}

// NewDiffPanel creates a new DiffPanel.
func NewDiffPanel() *DiffPanel {
	tv := tview.NewTextView()
	tv.SetDynamicColors(true)
	tv.SetScrollable(true)
	tv.SetWrap(true)

	frame := tview.NewFrame(tv).
		SetBorders(0, 0, 0, 0, 0, 0)
	frame.SetBorder(true).SetTitle(" Request Diff ")

	return &DiffPanel{
		frame:  frame,
		text:   tv,
		differ: grpcpkg.NewRequestDiffer(),
	}
}

// Primitive returns the tview primitive for layout embedding.
func (p *DiffPanel) Primitive() tview.Primitive {
	return p.frame
}

// Compare diffs left and right JSON payloads and renders the result.
func (p *DiffPanel) Compare(left, right string) error {
	res, err := p.differ.Diff(left, right)
	if err != nil {
		p.text.SetText(fmt.Sprintf("[red]Error: %v[-]", err))
		return err
	}
	p.render(res)
	return nil
}

// Clear resets the panel content.
func (p *DiffPanel) Clear() {
	p.text.Clear()
}

func (p *DiffPanel) render(res *grpcpkg.DiffResult) {
	p.text.Clear()
	if len(res.Added) == 0 && len(res.Removed) == 0 && len(res.Changed) == 0 {
		fmt.Fprint(p.text, "[green]No differences[-]")
		return
	}
	for k, v := range res.Added {
		fmt.Fprintf(p.text, "[green]+ %s: %v[-]\n", k, v)
	}
	for k, v := range res.Removed {
		fmt.Fprintf(p.text, "[red]- %s: %v[-]\n", k, v)
	}
	for k, v := range res.Changed {
		fmt.Fprintf(p.text, "[yellow]~ %s: %v -> %v[-]\n", k, v[0], v[1])
	}
}
