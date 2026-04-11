package ui

import (
	"fmt"

	"github.com/rivo/tview"
)

// ConnectionPanel displays the active connection pool status and
// provides an input field to add new gRPC server addresses.
type ConnectionPanel struct {
	*tview.Flex
	input      *tview.InputField
	statusView *tview.TextView
	onConnect  func(address string)
}

// NewConnectionPanel creates a new ConnectionPanel.
// onConnect is called with the entered address when the user submits.
func NewConnectionPanel(onConnect func(address string)) *ConnectionPanel {
	if onConnect == nil {
		onConnect = func(string) {}
	}

	input := tview.NewInputField().
		SetLabel("Address: ").
		SetPlaceholder("localhost:50051").
		SetFieldWidth(40)

	statusView := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[grey]No active connections[-]")

	panel := &ConnectionPanel{
		Flex:       tview.NewFlex().SetDirection(tview.FlexRow),
		input:      input,
		statusView: statusView,
		onConnect:  onConnect,
	}

	input.SetDoneFunc(func(key rune) {
		addr := input.GetText()
		if addr != "" {
			onConnect(addr)
			input.SetText("")
		}
	})

	panel.Flex.
		AddItem(tview.NewTextView().SetText("gRPC Connection").SetTextAlign(tview.AlignCenter), 1, 0, false).
		AddItem(input, 1, 0, true).
		AddItem(statusView, 0, 1, false)

	panel.Flex.SetBorder(true).SetTitle(" Connections ")

	return panel
}

// SetStatus updates the status text shown in the panel.
func (p *ConnectionPanel) SetStatus(activeCount int, addresses []string) {
	if activeCount == 0 {
		p.statusView.SetText("[grey]No active connections[-]")
		return
	}
	text := fmt.Sprintf("[green]%d active connection(s)[-]\n", activeCount)
	for _, addr := range addresses {
		text += fmt.Sprintf("  • %s\n", addr)
	}
	p.statusView.SetText(text)
}

// Focus sets keyboard focus to the address input field.
func (p *ConnectionPanel) Focus(delegate func(tview.Primitive)) {
	delegate(p.input)
}

// GetAddress returns the current text in the address input field.
func (p *ConnectionPanel) GetAddress() string {
	return p.input.GetText()
}
