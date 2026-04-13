package ui

import (
	"fmt"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// SamplerPanel displays sampled gRPC requests in the TUI.
type SamplerPanel struct {
	frame   *tview.Frame
	table   *tview.Table
	sampler *grpcpkg.RequestSampler
}

// NewSamplerPanel creates a new SamplerPanel backed by the given sampler.
func NewSamplerPanel(sampler *grpcpkg.RequestSampler) *SamplerPanel {
	if sampler == nil {
		s, _ := grpcpkg.NewRequestSampler(grpcpkg.DefaultSamplerPolicy())
		sampler = s
	}
	table := tview.NewTable().SetBorders(false)
	table.SetFixed(1, 0)

	// Header row
	table.SetCell(0, 0, tview.NewTableCell("[::b]#").SetExpansion(0))
	table.SetCell(0, 1, tview.NewTableCell("[::b]Method").SetExpansion(2))
	table.SetCell(0, 2, tview.NewTableCell("[::b]Timestamp").SetExpansion(1))

	frame := tview.NewFrame(table).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText("Sampled Requests", true, tview.AlignLeft, tview.Styles.PrimaryTextColor)

	return &SamplerPanel{
		frame:   frame,
		table:   table,
		sampler: sampler,
	}
}

// Primitive returns the tview primitive for embedding in layouts.
func (p *SamplerPanel) Primitive() tview.Primitive {
	return p.frame
}

// Refresh redraws the table with current sampled requests.
func (p *SamplerPanel) Refresh() {
	// Clear data rows (keep header at row 0)
	for p.table.GetRowCount() > 1 {
		p.table.RemoveRow(1)
	}
	samples := p.sampler.All()
	for i, sr := range samples {
		row := i + 1
		p.table.SetCell(row, 0, tview.NewTableCell(fmt.Sprintf("%d", row)))
		p.table.SetCell(row, 1, tview.NewTableCell(sr.Method))
		p.table.SetCell(row, 2, tview.NewTableCell(sr.Timestamp.Format("15:04:05.000")))
	}
}

// Clear removes all samples from the underlying sampler and refreshes.
func (p *SamplerPanel) Clear() {
	_ = p.sampler.Clear()
	p.Refresh()
}

// SampleCount returns the number of currently stored samples.
func (p *SamplerPanel) SampleCount() int {
	return p.sampler.Len()
}
