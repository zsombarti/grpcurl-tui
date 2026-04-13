package grpc

import (
	"testing"
)

func TestNewRequestSampler_NotNil(t *testing.T) {
	s, err := NewRequestSampler(DefaultSamplerPolicy())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil sampler")
	}
}

func TestDefaultSamplerPolicy_Values(t *testing.T) {
	p := DefaultSamplerPolicy()
	if p.Rate != 1.0 {
		t.Errorf("expected rate 1.0, got %v", p.Rate)
	}
	if p.MaxSize != 256 {
		t.Errorf("expected max size 256, got %d", p.MaxSize)
	}
}

func TestNewRequestSampler_InvalidPolicy_FallsBackToDefault(t *testing.T) {
	s, err := NewRequestSampler(SamplerPolicy{Rate: -0.5, MaxSize: -1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Policy().Rate != 1.0 {
		t.Errorf("expected fallback rate 1.0, got %v", s.Policy().Rate)
	}
}

func TestRequestSampler_Len_Empty(t *testing.T) {
	s, _ := NewRequestSampler(DefaultSamplerPolicy())
	if s.Len() != 0 {
		t.Errorf("expected 0, got %d", s.Len())
	}
}

func TestRequestSampler_Sample_RateOne_AlwaysSamples(t *testing.T) {
	s, _ := NewRequestSampler(SamplerPolicy{Rate: 1.0, MaxSize: 256})
	for i := 0; i < 10; i++ {
		s.Sample("pkg.Service/Method", map[string]interface{}{"i": i})
	}
	if s.Len() != 10 {
		t.Errorf("expected 10 samples, got %d", s.Len())
	}
}

func TestRequestSampler_Sample_RateZero_NeverSamples(t *testing.T) {
	s, _ := NewRequestSampler(SamplerPolicy{Rate: 0.0, MaxSize: 256})
	for i := 0; i < 20; i++ {
		s.Sample("pkg.Service/Method", nil)
	}
	if s.Len() != 0 {
		t.Errorf("expected 0 samples, got %d", s.Len())
	}
}

func TestRequestSampler_Eviction(t *testing.T) {
	s, _ := NewRequestSampler(SamplerPolicy{Rate: 1.0, MaxSize: 3})
	for i := 0; i < 5; i++ {
		s.Sample("m", map[string]interface{}{"i": i})
	}
	if s.Len() != 3 {
		t.Errorf("expected 3 after eviction, got %d", s.Len())
	}
}

func TestRequestSampler_All_ReturnsCopy(t *testing.T) {
	s, _ := NewRequestSampler(DefaultSamplerPolicy())
	s.Sample("m", map[string]interface{}{"x": 1})
	all := s.All()
	if len(all) != 1 {
		t.Errorf("expected 1, got %d", len(all))
	}
}

func TestRequestSampler_Clear_ResetsLen(t *testing.T) {
	s, _ := NewRequestSampler(DefaultSamplerPolicy())
	s.Sample("m", nil)
	if err := s.Clear(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Len() != 0 {
		t.Errorf("expected 0 after clear, got %d", s.Len())
	}
}

func TestRequestSampler_Clear_EmptyReturnsError(t *testing.T) {
	s, _ := NewRequestSampler(DefaultSamplerPolicy())
	if err := s.Clear(); err == nil {
		t.Fatal("expected error when clearing empty sampler")
	}
}
