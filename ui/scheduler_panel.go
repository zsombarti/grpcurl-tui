package ui

import (
	"fmt"
	"time"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// SchedulerPanel provides a TUI panel for managing scheduled gRPC requests.
type SchedulerPanel struct {
	frame     *tview.Frame
	list      *tview.List
	scheduler *grpcpkg.RequestScheduler
}

// NewSchedulerPanel creates a new SchedulerPanel backed by the given scheduler.
func NewSchedulerPanel(scheduler *grpcpkg.RequestScheduler) *SchedulerPanel {
	list := tview.NewList()
	list.SetBorder(false)
	list.ShowSecondaryText(true)

	frame := tview.NewFrame(list)
	frame.SetBorder(true)
	frame.SetTitle(" Scheduler ")

	return &SchedulerPanel{
		frame:     frame,
		list:      list,
		scheduler: scheduler,
	}
}

// Primitive returns the root tview primitive for layout embedding.
func (p *SchedulerPanel) Primitive() tview.Primitive {
	return p.frame
}

// Refresh re-renders the list of scheduled jobs from the scheduler.
func (p *SchedulerPanel) Refresh() {
	p.list.Clear()
	jobs := p.scheduler.List()
	if len(jobs) == 0 {
		p.list.AddItem("(no scheduled jobs)", "", 0, nil)
		return
	}
	for _, job := range jobs {
		title := fmt.Sprintf("[%s] %s → %s", job.ID, job.Address, job.Method)
		sub := fmt.Sprintf("every %s  |  added: %s", job.Interval.String(), job.CreatedAt.Format(time.RFC3339))
		p.list.AddItem(title, sub, 0, nil)
	}
}

// SelectedID returns the job ID of the currently selected list item, or empty string.
func (p *SchedulerPanel) SelectedID() string {
	jobs := p.scheduler.List()
	idx := p.list.GetCurrentItem()
	if idx < 0 || idx >= len(jobs) {
		return ""
	}
	return jobs[idx].ID
}

// RemoveSelected stops and removes the currently selected scheduled job.
func (p *SchedulerPanel) RemoveSelected() bool {
	id := p.SelectedID()
	if id == "" {
		return false
	}
	ok := p.scheduler.Remove(id)
	if ok {
		p.Refresh()
	}
	return ok
}
