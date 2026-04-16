package ui

import (
	"fmt"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// QuotaPanel displays per-method quota usage.
type QuotaPanel struct {
	frame *tview.Frame
	table *tview.Table
	quota *grpcpkg.RequestQuota
}

// NewQuotaPanel creates a QuotaPanel backed by the given RequestQuota.
func NewQuotaPanel(q *grpcpkg.RequestQuota) *QuotaPanel {
	if q == nil {
		q = grpcpkg.NewRequestQuota(grpcpkg.DefaultQuotaPolicy())
	}
	table := tview.NewTable().SetBorders(false)
	table.SetFixed(1, 0)

	frame := tview.NewFrame(table).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText("Quota Usage", true, tview.AlignLeft, tcell.ColorYellow)

	p := &QuotaPanel{frame: frame, table: table, quota: q}
	p.renderHeader()
	return p
}

// Primitive returns the tview primitive for layout embedding.
func (p *QuotaPanel) Primitive() tview.Primitive {
	return p.frame
}

// Refresh redraws the quota table from current usage data.
func (p *QuotaPanel) Refresh() {
	p.table.Clear()
	p.renderHeader()
	entries := p.quota.Usage()
	for i, e := range entries {
		row := i + 1
		p.table.SetCell(row, 0, tview.NewTableCell(e.Method).SetExpansion(2))
		p.table.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%d", e.Count)).SetExpansion(1))
		p.table.SetCell(row, 2, tview.NewTableCell(e.WindowEnd.Format("15:04:05")).SetExpansion(1))
	}
}

func (p *QuotaPanel) renderHeader() {
	headers := []string{"Method", "Count", "Window End"}
	for col, h := range headers {
		cell := tview.NewTableCell(h).
			SetSelectable(false).
			SetExpansion(1)
		p.table.SetCell(0, col, cell)
	}
}
