package ui

import (
	"context"
	"fmt"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// PipelinePanel displays the registered pipeline steps and allows running
// the pipeline against a sample payload from the UI.
type PipelinePanel struct {
	frame    *tview.Frame
	list     *tview.List
	pipeline *grpcpkg.RequestPipeline
}

// NewPipelinePanel creates a PipelinePanel wrapping the given RequestPipeline.
func NewPipelinePanel(p *grpcpkg.RequestPipeline) *PipelinePanel {
	list := tview.NewList().ShowSecondaryText(false)
	list.SetBorder(false)

	frame := tview.NewFrame(list).
		SetBorders(1, 1, 1, 1, 2, 2).
		AddText("Request Pipeline", true, tview.AlignCenter, tcell_white)

	return &PipelinePanel{
		frame:    frame,
		list:     list,
		pipeline: p,
	}
}

// Primitive returns the tview primitive for layout embedding.
func (pp *PipelinePanel) Primitive() tview.Primitive {
	return pp.frame
}

// Refresh re-renders the list of pipeline steps.
func (pp *PipelinePanel) Refresh() {
	pp.list.Clear()
	for i, name := range pp.pipeline.StepNames() {
		label := fmt.Sprintf("[%d] %s", i+1, name)
		pp.list.AddItem(label, "", 0, nil)
	}
	if pp.pipeline.Len() == 0 {
		pp.list.AddItem("(no steps registered)", "", 0, nil)
	}
}

// RunWithPayload executes the pipeline with the given payload and returns the
// result. Errors are surfaced to the caller for display in a status bar.
func (pp *PipelinePanel) RunWithPayload(ctx context.Context, payload map[string]interface{}) (map[string]interface{}, error) {
	return pp.pipeline.Run(ctx, payload)
}
