package ui

import (
	"fmt"

	"github.com/rivo/tview"

	"grpcurl-tui/grpc"
)

// CircuitLogPanel displays circuit breaker state transition logs.
type CircuitLogPanel struct {
	frame *tview.Frame
	list  *tview.List
	store *grpc.CircuitLog
}

// NewCircuitLogPanel creates a new CircuitLogPanel backed by the given CircuitLog.
func NewCircuitLogPanel(store *grpc.CircuitLog) *CircuitLogPanel {
	if store == nil {
		store = grpc.NewCircuitLog(50)
	}

	list := tview.NewList().ShowSecondaryText(true)
	frame := tview.NewFrame(list).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText("Circuit Breaker Log", true, tview.AlignLeft, 0)

	p := &CircuitLogPanel{
		frame: frame,
		list:  list,
		store: store,
	}
	p.Refresh()
	return p
}

// Primitive returns the tview primitive for embedding in a layout.
func (p *CircuitLogPanel) Primitive() tview.Primitive {
	return p.frame
}

// Refresh reloads all entries from the backing store into the list.
func (p *CircuitLogPanel) Refresh() {
	p.list.Clear()
	for _, e := range p.store.Entries() {
		main := fmt.Sprintf("[%s] %s", e.State, e.Method)
		secondary := fmt.Sprintf("%s — %s", e.Timestamp.Format("15:04:05"), e.Reason)
		p.list.AddItem(main, secondary, 0, nil)
	}
	if p.list.GetItemCount() == 0 {
		p.list.AddItem("No circuit log entries.", "", 0, nil)
	}
}

// Clear removes all entries from the store and refreshes the view.
func (p *CircuitLogPanel) Clear() {
	p.store.Clear()
	p.Refresh()
}
