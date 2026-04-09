package ui

import (
	"fmt"
	"strconv"
	"time"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// RetryPanel provides a UI form for configuring a RetryPolicy.
type RetryPanel struct {
	*tview.Form
	maxAttemptsField *tview.InputField
	initialDelayField *tview.InputField
	maxDelayField    *tview.InputField
	multiplierField  *tview.InputField
}

// NewRetryPanel creates and returns a RetryPanel pre-populated with defaults.
func NewRetryPanel() *RetryPanel {
	defaults := grpcpkg.DefaultRetryPolicy()
	p := &RetryPanel{
		Form:              tview.NewForm(),
		maxAttemptsField:  tview.NewInputField(),
		initialDelayField: tview.NewInputField(),
		maxDelayField:    tview.NewInputField(),
		multiplierField:  tview.NewInputField(),
	}

	p.maxAttemptsField.SetLabel("Max Attempts").SetText(strconv.Itoa(defaults.MaxAttempts))
	p.initialDelayField.SetLabel("Initial Delay (ms)").SetText(fmt.Sprintf("%d", defaults.InitialDelay.Milliseconds()))
	p.maxDelayField.SetLabel("Max Delay (ms)").SetText(fmt.Sprintf("%d", defaults.MaxDelay.Milliseconds()))
	p.multiplierField.SetLabel("Multiplier").SetText(fmt.Sprintf("%.1f", defaults.Multiplier))

	p.Form.SetBorder(true).SetTitle(" Retry Policy ").SetTitleAlign(tview.AlignLeft)
	p.Form.AddFormItem(p.maxAttemptsField)
	p.Form.AddFormItem(p.initialDelayField)
	p.Form.AddFormItem(p.maxDelayField)
	p.Form.AddFormItem(p.multiplierField)

	return p
}

// GetPolicy reads the current form values and returns a RetryPolicy.
// Invalid or missing values fall back to defaults.
func (p *RetryPanel) GetPolicy() grpcpkg.RetryPolicy {
	defaults := grpcpkg.DefaultRetryPolicy()

	maxAttempts, err := strconv.Atoi(p.maxAttemptsField.GetText())
	if err != nil || maxAttempts < 1 {
		maxAttempts = defaults.MaxAttempts
	}

	initialMs, err := strconv.ParseInt(p.initialDelayField.GetText(), 10, 64)
	if err != nil || initialMs < 0 {
		initialMs = defaults.InitialDelay.Milliseconds()
	}

	maxMs, err := strconv.ParseInt(p.maxDelayField.GetText(), 10, 64)
	if err != nil || maxMs < 0 {
		maxMs = defaults.MaxDelay.Milliseconds()
	}

	multiplier, err := strconv.ParseFloat(p.multiplierField.GetText(), 64)
	if err != nil || multiplier < 1.0 {
		multiplier = defaults.Multiplier
	}

	return grpcpkg.RetryPolicy{
		MaxAttempts:  maxAttempts,
		InitialDelay: time.Duration(initialMs) * time.Millisecond,
		MaxDelay:     time.Duration(maxMs) * time.Millisecond,
		Multiplier:   multiplier,
	}
}
