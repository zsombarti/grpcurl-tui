package ui

import (
	"testing"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

func newTestSamplerPanel(t *testing.T) *SamplerPanel {
	t.Helper()
	s, err := grpcpkg.NewRequestSampler(grpcpkg.DefaultSamplerPolicy())
	if err != nil {
		t.Fatalf("failed to create sampler: %v", err)
	}
	return NewSamplerPanel(s)
}

func TestNewSamplerPanel_NotNil(t *testing.T) {
	p := newTestSamplerPanel(t)
	if p == nil {
		t.Fatal("expected non-nil SamplerPanel")
	}
}

func TestSamplerPanel_Primitive_NotNil(t *testing.T) {
	p := newTestSamplerPanel(t)
	if p.Primitive() == nil {
		t.Fatal("expected non-nil primitive")
	}
}

func TestSamplerPanel_Primitive_IsFrame(t *testing.T) {
	p := newTestSamplerPanel(t)
	_, ok := p.Primitive().(*tview.Frame)
	if !ok {
		t.Fatal("expected primitive to be *tview.Frame")
	}
}

func TestSamplerPanel_NilSampler_FallsBackToDefault(t *testing.T) {
	p := NewSamplerPanel(nil)
	if p == nil {
		t.Fatal("expected non-nil panel even with nil sampler")
	}
}

func TestSamplerPanel_SampleCount_Initial_Zero(t *testing.T) {
	p := newTestSamplerPanel(t)
	if p.SampleCount() != 0 {
		t.Errorf("expected 0 samples initially, got %d", p.SampleCount())
	}
}

func TestSamplerPanel_Refresh_Empty(t *testing.T) {
	p := newTestSamplerPanel(t)
	p.Refresh() // should not panic
}

func TestSamplerPanel_Refresh_WithSamples(t *testing.T) {
	s, _ := grpcpkg.NewRequestSampler(grpcpkg.SamplerPolicy{Rate: 1.0, MaxSize: 256})
	s.Sample("pkg.Svc/MethodA", map[string]interface{}{"key": "val"})
	s.Sample("pkg.Svc/MethodB", nil)
	p := NewSamplerPanel(s)
	p.Refresh()
	if p.SampleCount() != 2 {
		t.Errorf("expected 2 samples, got %d", p.SampleCount())
	}
}

func TestSamplerPanel_Clear_ResetsCount(t *testing.T) {
	s, _ := grpcpkg.NewRequestSampler(grpcpkg.SamplerPolicy{Rate: 1.0, MaxSize: 256})
	s.Sample("m", nil)
	p := NewSamplerPanel(s)
	p.Clear()
	if p.SampleCount() != 0 {
		t.Errorf("expected 0 after clear, got %d", p.SampleCount())
	}
}
