package ui

import (
	"fmt"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// AuditPanel displays the request audit log in a scrollable table.
type AuditPanel struct {
	table   *tview.Table
	auditor *grpcpkg.RequestAuditor
}

// NewAuditPanel creates an AuditPanel backed by the given RequestAuditor.
func NewAuditPanel(auditor *grpcpkg.RequestAuditor) *AuditPanel {
	table := tview.NewTable().SetBorders(false).SetSelectable(true, false)
	table.SetTitle(" Audit Log ").SetBorder(true)

	// Header row
	headers := []string{"Time", "Method", "Address", "Status", "Duration", "Note"}
	for col, h := range headers {
		table.SetCell(0, col, tview.NewTableCell(h).
			SetTextColor(tview.Styles.SecondaryTextColor).
			SetSelectable(false))
	}

	return &AuditPanel{table: table, auditor: auditor}
}

// Primitive returns the underlying tview primitive for layout embedding.
func (p *AuditPanel) Primitive() tview.Primitive {
	return p.table
}

// Refresh re-renders the table from the current audit log.
func (p *AuditPanel) Refresh() {
	// Clear data rows (keep header at row 0)
	for p.table.GetRowCount() > 1 {
		p.table.RemoveRow(p.table.GetRowCount() - 1)
	}

	entries := p.auditor.Entries()
	for i, e := range entries {
		row := i + 1
		p.table.SetCell(row, 0, tview.NewTableCell(e.Timestamp.Format("15:04:05")))
		p.table.SetCell(row, 1, tview.NewTableCell(e.Method))
		p.table.SetCell(row, 2, tview.NewTableCell(e.Address))
		p.table.SetCell(row, 3, tview.NewTableCell(e.Status))
		p.table.SetCell(row, 4, tview.NewTableCell(fmt.Sprintf("%v", e.Duration.Round(1e6))))
		p.table.SetCell(row, 5, tview.NewTableCell(e.Note))
	}

	// Scroll to the latest entry
	if len(entries) > 0 {
		p.table.ScrollToEnd()
	}
}

// Clear wipes the audit log and refreshes the table.
func (p *AuditPanel) Clear() {
	p.auditor.Clear()
	p.Refresh()
}
