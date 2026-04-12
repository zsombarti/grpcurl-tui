package ui

import (
	"fmt"
	"time"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// DeduplicatorPanel displays and configures the request deduplication window.
type DeduplicatorPanel struct {
	frame       *tview.Frame
	form        *tview.Form
	deduplicator *grpcpkg.RequestDeduplicator
}

// NewDeduplicatorPanel creates a DeduplicatorPanel backed by the given deduplicator.
func NewDeduplicatorPanel(d *grpcpkg.RequestDeduplicator) *DeduplicatorPanel {
	form := tview.NewForm()
	form.SetBorder(false)

	form.AddInputField("Window (ms)", "5000", 10, nil, nil)
	form.AddInputField("Max Size", "256", 10, nil, nil)

	frame := tview.NewFrame(form).
		SetBorders(1, 1, 1, 1, 2, 2).
		AddText("Deduplicator", true, tview.AlignLeft, tcell_colorWhite())

	return &DeduplicatorPanel{frame: frame, form: form, deduplicator: d}
}

// Primitive returns the root tview primitive for layout embedding.
func (p *DeduplicatorPanel) Primitive() tview.Primitive { return p.frame }

// GetPolicy reads the form fields and returns a DeduplicatorPolicy.
// Invalid values fall back to defaults.
func (p *DeduplicatorPanel) GetPolicy() grpcpkg.DeduplicatorPolicy {
	def := grpcpkg.DefaultDeduplicatorPolicy()

	windowItem, ok := p.form.GetFormItemByLabel("Window (ms)").(*tview.InputField)
	if !ok {
		return def
	}
	maxItem, ok2 := p.form.GetFormItemByLabel("Max Size").(*tview.InputField)
	if !ok2 {
		return def
	}

	var windowMs int
	if _, err := fmt.Sscanf(windowItem.GetText(), "%d", &windowMs); err != nil || windowMs <= 0 {
		windowMs = int(def.WindowDuration.Milliseconds())
	}
	var maxSize int
	if _, err := fmt.Sscanf(maxItem.GetText(), "%d", &maxSize); err != nil || maxSize <= 0 {
		maxSize = def.MaxSize
	}

	return grpcpkg.DeduplicatorPolicy{
		WindowDuration: time.Duration(windowMs) * time.Millisecond,
		MaxSize:        maxSize,
	}
}

// StatusText returns a human-readable summary of the current deduplicator state.
func (p *DeduplicatorPanel) StatusText() string {
	pol := p.GetPolicy()
	return fmt.Sprintf("window=%v  maxSize=%d  active=%d",
		pol.WindowDuration, pol.MaxSize, p.deduplicator.Len())
}

// tcell_colorWhite is a local helper to avoid importing tcell directly.
func tcell_colorWhite() int32 { return 0x00ffffff }
