package ui

import (
	"fmt"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// RequestLoggerPanel displays logged requests in a scrollable list.
type RequestLoggerPanel struct {
	list   *tview.List
	logger *grpcpkg.RequestLogger
}

// NewRequestLoggerPanel creates a panel backed by the given RequestLogger.
func NewRequestLoggerPanel(logger *grpcpkg.RequestLogger) *RequestLoggerPanel {
	list := tview.NewList().
		ShowSecondaryText(true).
		SetHighlightFullLine(true)
	list.SetBorder(true).SetTitle(" Request Log ")
	return &RequestLoggerPanel{list: list, logger: logger}
}

// Primitive returns the underlying tview primitive for layout inclusion.
func (p *RequestLoggerPanel) Primitive() tview.Primitive {
	return p.list
}

// Refresh reloads all entries from the logger into the list.
func (p *RequestLoggerPanel) Refresh() {
	p.list.Clear()
	entries := p.logger.Entries()
	if len(entries) == 0 {
		p.list.AddItem("(no requests logged)", "", 0, nil)
		return
	}
	for i, e := range entries {
		summary, err := p.logger.Summary(i)
		if err != nil {
			continue
		}
		secondary := e.Response
		if e.Error != "" {
			secondary = fmt.Sprintf("error: %s", e.Error)
		}
		if len(secondary) > 80 {
			secondary = secondary[:80] + "..."
		}
		p.list.AddItem(summary, secondary, 0, nil)
	}
}

// SelectedEntry returns the RequestLogEntry for the currently selected row,
// or nil if the list is empty.
func (p *RequestLoggerPanel) SelectedEntry() *grpcpkg.RequestLogEntry {
	idx := p.list.GetCurrentItem()
	entries := p.logger.Entries()
	if idx < 0 || idx >= len(entries) {
		return nil
	}
	e := entries[idx]
	return &e
}

// Clear wipes the underlying logger and refreshes the panel.
func (p *RequestLoggerPanel) Clear() {
	p.logger.Clear()
	p.Refresh()
}
