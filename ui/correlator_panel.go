package ui

import (
	"fmt"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// CorrelatorPanel displays tracked correlation IDs in a scrollable table.
type CorrelatorPanel struct {
	primitive   *tview.Flex
	table       *tview.Table
	correlator  *grpcpkg.RequestCorrelator
}

// NewCorrelatorPanel constructs a CorrelatorPanel backed by the given correlator.
func NewCorrelatorPanel(c *grpcpkg.RequestCorrelator) *CorrelatorPanel {
	if c == nil {
		c = grpcpkg.NewRequestCorrelator(0)
	}

	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false)

	headers := []string{"Correlation ID", "Method", "Timestamp"}
	for col, h := range headers {
		table.SetCell(0, col, tview.NewTableCell(h).
			SetTextColor(tview.Styles.SecondaryTextColor).
			SetSelectable(false))
	}

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(table, 0, 1, true)
	flex.SetBorder(true).SetTitle(" Correlation IDs ")

	return &CorrelatorPanel{
		primitive:  flex,
		table:      table,
		correlator: c,
	}
}

// Primitive returns the tview primitive for embedding in layouts.
func (p *CorrelatorPanel) Primitive() tview.Primitive {
	return p.primitive
}

// Refresh re-reads the correlator entries and redraws the table.
func (p *CorrelatorPanel) Refresh(entries []grpcpkg.CorrelationEntry) {
	// clear data rows, keep header
	for p.table.GetRowCount() > 1 {
		p.table.RemoveRow(p.table.GetRowCount() - 1)
	}
	for i, e := range entries {
		row := i + 1
		p.table.SetCell(row, 0, tview.NewTableCell(e.ID))
		p.table.SetCell(row, 1, tview.NewTableCell(e.Method))
		p.table.SetCell(row, 2, tview.NewTableCell(
			fmt.Sprintf("%s", e.Timestamp.Format("15:04:05.000")),
		))
	}
}
