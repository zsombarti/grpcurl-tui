package grpc

import (
	"context"
	"errors"
	"testing"
)

func TestNewRequestPipeline_NotNil(t *testing.T) {
	p := NewRequestPipeline()
	if p == nil {
		t.Fatal("expected non-nil pipeline")
	}
}

func TestRequestPipeline_Len_Empty(t *testing.T) {
	p := NewRequestPipeline()
	if p.Len() != 0 {
		t.Fatalf("expected 0, got %d", p.Len())
	}
}

func TestRequestPipeline_AddStep_IncrementsLen(t *testing.T) {
	p := NewRequestPipeline()
	p.AddStep(PipelineStep{Name: "s1", Handler: func(_ context.Context, m map[string]interface{}) (map[string]interface{}, error) { return m, nil }})
	if p.Len() != 1 {
		t.Fatalf("expected 1, got %d", p.Len())
	}
}

func TestRequestPipeline_Run_PassesPayloadThrough(t *testing.T) {
	p := NewRequestPipeline()
	p.AddStep(PipelineStep{
		Name: "add-key",
		Handler: func(_ context.Context, m map[string]interface{}) (map[string]interface{}, error) {
			m["injected"] = true
			return m, nil
		},
	})
	out, err := p.Run(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["injected"] != true {
		t.Fatal("expected injected key to be true")
	}
}

func TestRequestPipeline_Run_StepError_Propagates(t *testing.T) {
	p := NewRequestPipeline()
	p.AddStep(PipelineStep{
		Name: "fail",
		Handler: func(_ context.Context, m map[string]interface{}) (map[string]interface{}, error) {
			return nil, errors.New("step failure")
		},
	})
	_, err := p.Run(context.Background(), map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error from failing step")
	}
}

func TestRequestPipeline_Run_CancelledContext(t *testing.T) {
	p := NewRequestPipeline()
	p.AddStep(PipelineStep{
		Name: "never-runs",
		Handler: func(_ context.Context, m map[string]interface{}) (map[string]interface{}, error) {
			return m, nil
		},
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := p.Run(ctx, map[string]interface{}{})
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}

func TestRequestPipeline_Clear_ResetsLen(t *testing.T) {
	p := NewRequestPipeline()
	p.AddStep(PipelineStep{Name: "s1", Handler: func(_ context.Context, m map[string]interface{}) (map[string]interface{}, error) { return m, nil }})
	p.Clear()
	if p.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", p.Len())
	}
}

func TestRequestPipeline_StepNames(t *testing.T) {
	p := NewRequestPipeline()
	p.AddStep(PipelineStep{Name: "alpha", Handler: func(_ context.Context, m map[string]interface{}) (map[string]interface{}, error) { return m, nil }})
	p.AddStep(PipelineStep{Name: "beta", Handler: func(_ context.Context, m map[string]interface{}) (map[string]interface{}, error) { return m, nil }})
	names := p.StepNames()
	if len(names) != 2 || names[0] != "alpha" || names[1] != "beta" {
		t.Fatalf("unexpected step names: %v", names)
	}
}
