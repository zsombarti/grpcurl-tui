package ui

import (
	"fmt"
	"sort"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// WatermarkPanel displays per-method latency watermarks in a table.
type WatermarkPanel struct {
	table     *tview.Table
	watermark *grpcpkg.RequestWatermark
}

// NewWatermarkPanel creates a WatermarkPanel backed by the given RequestWatermark.
func NewWatermarkPanel(w *grpcpkg.RequestWatermark) *WatermarkPanel {
	if w == nil {
		w = grpcpkg.NewRequestWatermark(0)
	}
	table := tview.NewTable().SetBorders(false).SetFixed(1, 0)
	table.SetTitle(" Latency Watermarks ").SetBorder(true)
	p := &WatermarkPanel{table: table, watermark: w}
	p.setHeaders()
	return p
}

func (p *WatermarkPanel) setHeaders() {
	headers := []string{"Method", "Min", "Max", "Last", "Count"}
	for i, h := range headers {
		p.table.SetCell(0, i, tview.NewTableCell(h).SetSelectable(false).SetAttributes(1))
	}
}

// Primitive returns the underlying tview primitive.
func (p *WatermarkPanel) Primitive() tview.Primitive {
	return p.table
}

// Refresh redraws the table from the current watermark data.
func (p *WatermarkPanel) Refresh() {
	p.table.Clear()
	p.setHeaders()
	entries := p.watermark.All()
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Method < entries[j].Method
	})
	for row, e := range entries {
		p.table.SetCell(row+1, 0, tview.NewTableCell(e.Method))
		p.table.SetCell(row+1, 1, tview.NewTableCell(fmt.Sprintf("%v", e.Min.Round(1))))
		p.table.SetCell(row+1, 2, tview.NewTableCell(fmt.Sprintf("%v", e.Max.Round(1))))
		p.table.SetCell(row+1, 3, tview.NewTableCell(fmt.Sprintf("%v", e.Last.Round(1))))
		p.table.SetCell(row+1, 4, tview.NewTableCell(fmt.Sprintf("%d", e.Count)))
	}
}
