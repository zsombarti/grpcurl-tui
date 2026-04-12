package grpc

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestNewRequestScheduler_NotNil(t *testing.T) {
	s := NewRequestScheduler(DefaultSchedulerPolicy())
	if s == nil {
		t.Fatal("expected non-nil scheduler")
	}
}

func TestDefaultSchedulerPolicy_Values(t *testing.T) {
	p := DefaultSchedulerPolicy()
	if p.MaxJobs <= 0 {
		t.Errorf("expected positive MaxJobs, got %d", p.MaxJobs)
	}
	if p.DefaultInterval <= 0 {
		t.Errorf("expected positive DefaultInterval, got %v", p.DefaultInterval)
	}
}

func TestNewRequestScheduler_InvalidPolicy_FallsBackToDefault(t *testing.T) {
	s := NewRequestScheduler(SchedulerPolicy{MaxJobs: -1})
	if s.policy.MaxJobs != DefaultSchedulerPolicy().MaxJobs {
		t.Errorf("expected fallback to default MaxJobs")
	}
}

func TestRequestScheduler_Len_Empty(t *testing.T) {
	s := NewRequestScheduler(DefaultSchedulerPolicy())
	if s.Len() != 0 {
		t.Errorf("expected 0, got %d", s.Len())
	}
}

func TestRequestScheduler_Add_And_Len(t *testing.T) {
	s := NewRequestScheduler(DefaultSchedulerPolicy())
	err := s.Add(ScheduledRequest{ID: "job1", Interval: 100 * time.Millisecond}, func(_ ScheduledRequest) {})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Len() != 1 {
		t.Errorf("expected 1, got %d", s.Len())
	}
	s.StopAll()
}

func TestRequestScheduler_Add_EmptyID_ReturnsError(t *testing.T) {
	s := NewRequestScheduler(DefaultSchedulerPolicy())
	err := s.Add(ScheduledRequest{ID: ""}, func(_ ScheduledRequest) {})
	if err == nil {
		t.Fatal("expected error for empty ID")
	}
}

func TestRequestScheduler_Add_DuplicateID_ReturnsError(t *testing.T) {
	s := NewRequestScheduler(DefaultSchedulerPolicy())
	_ = s.Add(ScheduledRequest{ID: "dup", Interval: time.Second}, func(_ ScheduledRequest) {})
	err := s.Add(ScheduledRequest{ID: "dup", Interval: time.Second}, func(_ ScheduledRequest) {})
	if err == nil {
		t.Fatal("expected error for duplicate ID")
	}
	s.StopAll()
}

func TestRequestScheduler_MaxJobs_Enforced(t *testing.T) {
	s := NewRequestScheduler(SchedulerPolicy{MaxJobs: 2, DefaultInterval: time.Second})
	_ = s.Add(ScheduledRequest{ID: "a", Interval: time.Second}, func(_ ScheduledRequest) {})
	_ = s.Add(ScheduledRequest{ID: "b", Interval: time.Second}, func(_ ScheduledRequest) {})
	err := s.Add(ScheduledRequest{ID: "c", Interval: time.Second}, func(_ ScheduledRequest) {})
	if err == nil {
		t.Fatal("expected error when max jobs reached")
	}
	s.StopAll()
}

func TestRequestScheduler_Remove(t *testing.T) {
	s := NewRequestScheduler(DefaultSchedulerPolicy())
	_ = s.Add(ScheduledRequest{ID: "r1", Interval: time.Second}, func(_ ScheduledRequest) {})
	if !s.Remove("r1") {
		t.Fatal("expected Remove to return true")
	}
	if s.Len() != 0 {
		t.Errorf("expected 0 after remove, got %d", s.Len())
	}
}

func TestRequestScheduler_Fires(t *testing.T) {
	s := NewRequestScheduler(DefaultSchedulerPolicy())
	var count int32
	_ = s.Add(ScheduledRequest{ID: "fire", Interval: 30 * time.Millisecond}, func(_ ScheduledRequest) {
		atomic.AddInt32(&count, 1)
	})
	time.Sleep(120 * time.Millisecond)
	s.StopAll()
	if atomic.LoadInt32(&count) < 2 {
		t.Errorf("expected at least 2 fires, got %d", count)
	}
}
