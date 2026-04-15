package ui

import (
	"fmt"

	"github.com/rivo/tview"

	"grpcurl-tui/grpc"
)

// RetryLogPanel displays a live view of retry log entries.
type RetryLogPanel struct {
	frame *tview.Frame
	list  *tview.List
	store *grpc.RetryLog
}

// NewRetryLogPanel creates a RetryLogPanel backed by the given RetryLog.
// A non-nil store is required; pass NewRetryLog(0) for a default one.
func NewRetryLogPanel(store *grpc.RetryLog) *RetryLogPanel {
	if store == nil {
		store = grpc.NewRetryLog(0)
	}
	list := tview.NewList().ShowSecondaryText(true)
	frame := tview.NewFrame(list).
		SetBorders(1, 1, 1, 1, 2, 2).
		AddText("Retry Log", true, tview.AlignLeft, 0)
	return &RetryLogPanel{frame: frame, list: list, store: store}
}

// Primitive returns the root tview primitive for layout embedding.
func (p *RetryLogPanel) Primitive() tview.Primitive {
	return p.frame
}

// Refresh repaints the list from the current store contents.
func (p *RetryLogPanel) Refresh() {
	p.list.Clear()
	for _, e := range p.store.Entries() {
		errText := ""
		if e.Err != nil {
			errText = e.Err.Error()
		}
		primary := fmt.Sprintf("[%s] attempt #%d", e.Method, e.Attempt)
		secondary := fmt.Sprintf("%s  err: %s", e.Timestamp.Format("15:04:05"), errText)
		p.list.AddItem(primary, secondary, 0, nil)
	}
}

// Clear wipes the backing store and repaints.
func (p *RetryLogPanel) Clear() {
	p.store.Clear()
	p.Refresh()
}
