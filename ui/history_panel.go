package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// HistoryPanel is a tview component that displays past gRPC calls.
type HistoryPanel struct {
	*tview.Table
	history *grpcpkg.History
}

// NewHistoryPanel creates a HistoryPanel backed by the provided History store.
func NewHistoryPanel(h *grpcpkg.History) *HistoryPanel {
	table := tview.NewTable()
	table.SetBorders(false)
	table.SetSelectable(true, false)
	table.SetTitle(" History ")
	table.SetBorder(true)
	table.SetBorderColor(tcell.ColorDarkCyan)

	p := &HistoryPanel{Table: table, history: h}
	p.renderHeaders()
	return p
}

func (p *HistoryPanel) renderHeaders() {
	headers := []string{"Time", "Address", "Service", "Method", "Status"}
	for col, h := range headers {
		cell := tview.NewTableCell(h).
			SetTextColor(tcell.ColorYellow).
			SetSelectable(false).
			SetExpansion(1)
		p.SetCell(0, col, cell)
	}
}

// Refresh re-renders the table from the current history entries.
func (p *HistoryPanel) Refresh() {
	p.Clear()
	p.renderHeaders()
	for row, entry := range p.history.All() {
		status := "OK"
		color := tcell.ColorGreen
		if entry.Error != "" {
			status = "ERR"
			color = tcell.ColorRed
		}
		values := []string{
			entry.Timestamp.Format("15:04:05"),
			entry.Address,
			entry.Service,
			entry.Method,
			status,
		}
		for col, v := range values {
			cell := tview.NewTableCell(fmt.Sprintf(" %s ", v)).
				SetExpansion(1)
			if col == 4 {
				cell.SetTextColor(color)
			}
			p.SetCell(row+1, col, cell)
		}
	}
}
