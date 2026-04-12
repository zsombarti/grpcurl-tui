package ui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// ProfilerPanel displays per-request performance metrics.
type ProfilerPanel struct {
	frame    *tview.Frame
	table    *tview.Table
	profiler *grpcpkg.RequestProfiler
}

// NewProfilerPanel creates a ProfilerPanel backed by the given RequestProfiler.
func NewProfilerPanel(profiler *grpcpkg.RequestProfiler) *ProfilerPanel {
	table := tview.NewTable().SetBorders(false).SetFixed(1, 0)
	table.SetCell(0, 0, tview.NewTableCell("Method").SetSelectable(false))
	table.SetCell(0, 1, tview.NewTableCell("Duration").SetSelectable(false))
	table.SetCell(0, 2, tview.NewTableCell("ReqSize").SetSelectable(false))
	table.SetCell(0, 3, tview.NewTableCell("RespSize").SetSelectable(false))
	table.SetCell(0, 4, tview.NewTableCell("Error").SetSelectable(false))

	frame := tview.NewFrame(table).
		SetBorders(1, 1, 1, 1, 2, 2).
		AddText("Request Profiler", true, tview.AlignCenter, tcell.ColorYellow)

	return &ProfilerPanel{
		frame:    frame,
		table:    table,
		profiler: profiler,
	}
}

// Primitive returns the tview primitive for embedding in layouts.
func (p *ProfilerPanel) Primitive() tview.Primitive {
	return p.frame
}

// Refresh reloads all profiler entries into the table.
func (p *ProfilerPanel) Refresh() {
	// Clear data rows (keep header at row 0)
	for p.table.GetRowCount() > 1 {
		p.table.RemoveRow(p.table.GetRowCount() - 1)
	}

	entries := p.profiler.All()
	for i, e := range entries {
		row := i + 1
		errStr := ""
		errColor := tcell.ColorWhite
		if e.Error {
			errStr = "yes"
			errColor = tcell.ColorRed
		}
		p.table.SetCell(row, 0, tview.NewTableCell(e.Method))
		p.table.SetCell(row, 1, tview.NewTableCell(roundDuration(e.Duration)))
		p.table.SetCell(row, 2, tview.NewTableCell(fmt.Sprintf("%d B", e.RequestSize)))
		p.table.SetCell(row, 3, tview.NewTableCell(fmt.Sprintf("%d B", e.ResponseSize)))
		p.table.SetCell(row, 4, tview.NewTableCell(errStr).SetTextColor(errColor))
	}
}

// Clear removes all profiler entries from the table and resets the underlying profiler.
func (p *ProfilerPanel) Clear() {
	for p.table.GetRowCount() > 1 {
		p.table.RemoveRow(p.table.GetRowCount() - 1)
	}
	p.profiler.Reset()
}

func roundDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
	return fmt.Sprintf("%dms", d.Milliseconds())
}
