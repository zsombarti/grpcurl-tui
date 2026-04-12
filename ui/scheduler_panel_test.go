package ui

import (
	"testing"
	"time"

	grpcpkg "grpcurl-tui/grpc"
)

func newTestSchedulerPanel() *SchedulerPanel {
	s := grpcpkg.NewRequestScheduler(grpcpkg.DefaultSchedulerPolicy())
	return NewSchedulerPanel(s)
}

func TestNewSchedulerPanel_NotNil(t *testing.T) {
	p := newTestSchedulerPanel()
	if p == nil {
		t.Fatal("expected non-nil panel")
	}
}

func TestSchedulerPanel_Primitive_NotNil(t *testing.T) {
	p := newTestSchedulerPanel()
	if p.Primitive() == nil {
		t.Fatal("expected non-nil primitive")
	}
}

func TestSchedulerPanel_Refresh_Empty(t *testing.T) {
	p := newTestSchedulerPanel()
	// Should not panic on empty scheduler
	p.Refresh()
	if p.list.GetItemCount() != 1 {
		t.Errorf("expected 1 placeholder item, got %d", p.list.GetItemCount())
	}
}

func TestSchedulerPanel_Refresh_WithEntries(t *testing.T) {
	s := grpcpkg.NewRequestScheduler(grpcpkg.DefaultSchedulerPolicy())
	p := NewSchedulerPanel(s)
	_ = s.Add(grpcpkg.ScheduledRequest{
		ID:       "j1",
		Address:  "localhost:50051",
		Method:   "SayHello",
		Interval: 2 * time.Second,
	}, func(_ grpcpkg.ScheduledRequest) {})
	p.Refresh()
	if p.list.GetItemCount() != 1 {
		t.Errorf("expected 1 job item, got %d", p.list.GetItemCount())
	}
	s.StopAll()
}

func TestSchedulerPanel_SelectedID_Empty(t *testing.T) {
	p := newTestSchedulerPanel()
	p.Refresh()
	if p.SelectedID() != "" {
		t.Errorf("expected empty string for empty scheduler")
	}
}

func TestSchedulerPanel_RemoveSelected_NoOp(t *testing.T) {
	p := newTestSchedulerPanel()
	p.Refresh()
	if p.RemoveSelected() {
		t.Error("expected false when nothing to remove")
	}
}
