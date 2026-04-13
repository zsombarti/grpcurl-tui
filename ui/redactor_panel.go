package ui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// RedactorPanel provides a UI for configuring the request redactor.
type RedactorPanel struct {
	frame    *tview.Frame
	form     *tview.Form
actor *grpcpkg.RequestRedactor
	fieldBuf string
}

// NewRedactorPanel creates a new RedactorPanel backed by the given redactor.
func NewRedactorPanel(redactor *grpcpkg.RequestRedactor) *RedactorPanel {
	if redactor == nil {
		redactor = grpcpkg.NewRequestRedactor(grpcpkg.DefaultRedactorPolicy())
	}

	p := &RedactorPanel{redactor: redactor}

	p.form = tview.NewForm()
	p.form.SetBorder(false	p.form.AddInputField("Field to redact", "", 30, nil, func(text string) {
		p.fieldBuf = strings.TrimSpace(text)
	})

	p.form.AddButton("Add Field", func() {
		if p.fieldBuf != "" {
			_ = p.redactor.AddField(p.fieldBuf)
			p.fieldBuf = ""
			if item, ok := p.form.GetFormItemByLabel("Field to redact").(*tview.InputField); ok {
				item.SetText("")
			}
			p.refreshStatus()
		}
	})

	p.frame = tview.NewFrame(p.form).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText("Redactor", true, tview.AlignLeft, tcell.ColorYellow)

	p.refreshStatus()
	return p
}

// Primitive returns the root tview primitive for embedding in layouts.
func (p *RedactorPanel) Primitive() tview.Primitive {
	return p.frame
}

// Redactor returns the underlying RequestRedactor.
func (p *RedactorPanel) Redactor() *grpcpkg.RequestRedactor {
	return p.redactor
}

func (p *RedactorPanel) refreshStatus() {
	p.frame.Clear()
	p.frame.AddText("Redactor", true, tview.AlignLeft, tcell.ColorYellow)
	p.frame.AddText(
		fmt.Sprintf("Fields registered: %d", p.redactor.Len()),
		false, tview.AlignLeft, tcell.ColorWh
	)
}
