package ui

import (
	"context"
	"testing"

	grpcpkg "grpcurl-tui/grpc"
)

func newTestPipelinePanel() (*PipelinePanel, *grpcpkg.RequestPipeline) {
	p := grpcpkg.NewRequestPipeline()
	return NewPipelinePanel(p), p
}

func TestNewPipelinePanel_NotNil(t *testing.T) {
	pp, _ := newTestPipelinePanel()
	if pp == nil {
		t.Fatal("expected non-nil PipelinePanel")
	}
}

func TestPipelinePanel_Primitive_NotNil(t *testing.T) {
	pp, _ := newTestPipelinePanel()
	if pp.Primitive() == nil {
		t.Fatal("expected non-nil primitive")
	}
}

func TestPipelinePanel_Refresh_Empty(t *testing.T) {
	pp, _ := newTestPipelinePanel()
	// Should not panic on empty pipeline
	pp.Refresh()
	if pp.list.GetItemCount() != 1 {
		t.Fatalf("expected 1 placeholder item, got %d", pp.list.GetItemCount())
	}
}

func TestPipelinePanel_Refresh_WithSteps(t *testing.T) {
	pp, p := newTestPipelinePanel()
	p.AddStep(grpcpkg.PipelineStep{
		Name: "enrich",
		Handler: func(_ context.Context, m map[string]interface{}) (map[string]interface{}, error) {
			return m, nil
		},
	})
	p.AddStep(grpcpkg.PipelineStep{
		Name: "validate",
		Handler: func(_ context.Context, m map[string]interface{}) (map[string]interface{}, error) {
			return m, nil
		},
	})
	pp.Refresh()
	if pp.list.GetItemCount() != 2 {
		t.Fatalf("expected 2 items, got %d", pp.list.GetItemCount())
	}
}

func TestPipelinePanel_RunWithPayload_Success(t *testing.T) {
	pp, p := newTestPipelinePanel()
	p.AddStep(grpcpkg.PipelineStep{
		Name: "tag",
		Handler: func(_ context.Context, m map[string]interface{}) (map[string]interface{}, error) {
			m["tagged"] = true
			return m, nil
		},
	})
	out, err := pp.RunWithPayload(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["tagged"] != true {
		t.Fatal("expected tagged key in output")
	}
}

func TestPipelinePanel_RunWithPayload_CancelledContext(t *testing.T) {
	pp, p := newTestPipelinePanel()
	p.AddStep(grpcpkg.PipelineStep{
		Name: "noop",
		Handler: func(_ context.Context, m map[string]interface{}) (map[string]interface{}, error) {
			return m, nil
		},
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := pp.RunWithPayload(ctx, map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error from cancelled context")
	}
}
