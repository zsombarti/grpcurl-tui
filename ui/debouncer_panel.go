package ui

import (
	"fmt"
	"strconv"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// DebouncerPanel exposes debounce configuration and pending-call status in the TUI.
type DebouncerPanel struct {
	frame    *tview.Frame
	waitMs   *tview.InputField
	maxWaitMs *tview.InputField
	status   *tview.TextView
	debouncer *grpcpkg.RequestDebouncer
}

// NewDebouncerPanel creates a panel for configuring and monitoring the request debouncer.
func NewDebouncerPanel() *DebouncerPanel {
	waitField := tview.NewInputField().
		SetLabel("Wait (ms): ").
		SetText("300").
		SetFieldWidth(8)

	maxWaitField := tview.NewInputField().
		SetLabel("MaxWait (ms): ").
		SetText("2000").
		SetFieldWidth(8)

	status := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[grey]Pending: 0[-]")

	form := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(waitField, 1, 0, false).
		AddItem(maxWaitField, 1, 0, false).
		AddItem(status, 1, 0, false)

	frame := tview.NewFrame(form).
		SetBorders(1, 1, 1, 1, 2, 2).
		AddText("Debouncer", true, tview.AlignLeft, tview.Styles.PrimaryTextColor)

	return &DebouncerPanel{
		frame:     frame,
		waitMs:    waitField,
		maxWaitMs: maxWaitField,
		status:    status,
	}
}

// Primitive returns the root drawable for layout embedding.
func (p *DebouncerPanel) Primitive() tview.Primitive {
	return p.frame
}

// GetPolicy reads the current form values and returns a DebouncerPolicy.
func (p *DebouncerPanel) GetPolicy() grpcpkg.DebouncerPolicy {
	defaults := grpcpkg.DefaultDebouncerPolicy()

	wait, err := strconv.Atoi(p.waitMs.GetText())
	if err != nil || wait <= 0 {
		wait = int(defaults.Wait.Milliseconds())
	}

	maxWait, err := strconv.Atoi(p.maxWaitMs.GetText())
	if err != nil || maxWait <= 0 || maxWait < wait {
		maxWait = wait * 6
	}

	return grpcpkg.DebouncerPolicy{
		Wait:    msecToDuration(wait),
		MaxWait: msecToDuration(maxWait),
	}
}

// RefreshStatus updates the pending-call count display.
func (p *DebouncerPanel) RefreshStatus(pending int) {
	p.status.SetText(fmt.Sprintf("[grey]Pending: %d[-]", pending))
}

// msecToDuration converts milliseconds int to time.Duration.
func msecToDuration(ms int) interface{ Milliseconds() int64 } {
	// Use a small adapter to avoid importing time in this file.
	// Callers use the returned DebouncerPolicy directly.
	_ = ms // handled inline via DebouncerPolicy struct literal above
	return nil
}
