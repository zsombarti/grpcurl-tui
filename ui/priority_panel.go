package ui

import (
	"fmt"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// PriorityPanel displays and manages the request priority queue.
type PriorityPanel struct {
	frame *tview.Frame
	list  *tview.List
	queue *grpcpkg.RequestPriorityQueue
}

// NewPriorityPanel creates a panel backed by the given priority queue.
func NewPriorityPanel(queue *grpcpkg.RequestPriorityQueue) *PriorityPanel {
	list := tview.NewList().ShowSecondaryText(false)
	frame := tview.NewFrame(list).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText("Priority Queue", true, tview.AlignLeft, 0)
	return &PriorityPanel{frame: frame, list: list, queue: queue}
}

// Primitive returns the tview primitive for layout embedding.
func (p *PriorityPanel) Primitive() tview.Primitive {
	return p.frame
}

// Refresh redraws the list from the current queue snapshot.
func (p *PriorityPanel) Refresh(entries []grpcpkg.PriorityEntry) {
	p.list.Clear()
	for i, e := range entries {
		label := priorityLabel(e.Priority)
		text := fmt.Sprintf("[%d] [%s] %s", i+1, label, e.Label)
		p.list.AddItem(text, "", 0, nil)
	}
}

// Len returns the number of items currently displayed.
func (p *PriorityPanel) Len() int {
	return p.list.GetItemCount()
}

func priorityLabel(level grpcpkg.PriorityLevel) string {
	switch level {
	case grpcpkg.PriorityHigh:
		return "HIGH"
	case grpcpkg.PriorityNormal:
		return "NORMAL"
	default:
		return "LOW"
	}
}
