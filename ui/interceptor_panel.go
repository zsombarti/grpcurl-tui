package ui

import (
	"strings"

	"github.com/rivo/tview"
)

// InterceptorPanel displays and configures active call interceptors in the TUI.
type InterceptorPanel struct {
	root    *tview.Flex
	logView *tview.TextView
	toggle  *tview.Checkbox
}

// NewInterceptorPanel creates a new InterceptorPanel with log view and enable toggle.
func NewInterceptorPanel() *InterceptorPanel {
	p := &InterceptorPanel{}

	p.logView = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWordWrap(true)
	p.logView.SetBorder(true).SetTitle(" Call Log ")

	p.toggle = tview.NewCheckbox().
		SetLabel("Enable call logging ").
		SetChecked(true)

	p.root = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(p.toggle, 1, 0, false).
		AddItem(p.logView, 0, 1, false)

	p.root.SetBorder(true).SetTitle(" Interceptors ")
	return p
}

// Primitive returns the root tview primitive for layout embedding.
func (p *InterceptorPanel) Primitive() tview.Primitive {
	return p.root
}

// IsLoggingEnabled returns whether call logging is currently toggled on.
func (p *InterceptorPanel) IsLoggingEnabled() bool {
	return p.toggle.IsChecked()
}

// AppendLog appends a formatted call log entry to the log view.
func (p *InterceptorPanel) AppendLog(method, duration, errMsg string) {
	var sb strings.Builder
	sb.WriteString("[green]" + method + "[-] ")
	sb.WriteString("[yellow]" + duration + "[-]")
	if errMsg != "" {
		sb.WriteString(" [red]ERR: " + errMsg + "[-]")
	}
	sb.WriteString("\n")
	p.logView.Write([]byte(sb.String())) //nolint:errcheck
}

// Clear resets the log view.
func (p *InterceptorPanel) Clear() {
	p.logView.Clear()
}

// Len returns the number of lines currently in the log view.
func (p *InterceptorPanel) Len() int {
	text := p.logView.GetText(false)
	if text == "" {
		return 0
	}
	return strings.Count(text, "\n")
}
