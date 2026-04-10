package ui

import (
	"strconv"
	"time"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// TimeoutPanel provides UI controls for configuring gRPC timeout settings.
type TimeoutPanel struct {
	form           *tview.Form
	dialField      *tview.InputField
	requestField   *tview.InputField
}

// NewTimeoutPanel creates and returns a new TimeoutPanel with default values.
func NewTimeoutPanel() *TimeoutPanel {
	defaults := grpcpkg.DefaultTimeoutPolicy()

	dialField := tview.NewInputField().
		SetLabel("Dial Timeout (s): ").
		SetText(strconv.Itoa(int(defaults.DialTimeout.Seconds()))).
		SetFieldWidth(10)

	requestField := tview.NewInputField().
		SetLabel("Request Timeout (s): ").
		SetText(strconv.Itoa(int(defaults.RequestTimeout.Seconds()))).
		SetFieldWidth(10)

	form := tview.NewForm().
		AddFormItem(dialField).
		AddFormItem(requestField)
	form.SetBorder(true).SetTitle(" Timeouts ")

	return &TimeoutPanel{
		form:         form,
		dialField:    dialField,
		requestField: requestField,
	}
}

// GetPolicy reads the current form values and returns a TimeoutPolicy.
// Falls back to defaults for invalid or zero inputs.
func (p *TimeoutPanel) GetPolicy() grpcpkg.TimeoutPolicy {
	defaults := grpcpkg.DefaultTimeoutPolicy()

	dialSec, err := strconv.Atoi(p.dialField.GetText())
	if err != nil || dialSec <= 0 {
		dialSec = int(defaults.DialTimeout.Seconds())
	}

	reqSec, err := strconv.Atoi(p.requestField.GetText())
	if err != nil || reqSec <= 0 {
		reqSec = int(defaults.RequestTimeout.Seconds())
	}

	return grpcpkg.TimeoutPolicy{
		DialTimeout:    time.Duration(dialSec) * time.Second,
		RequestTimeout: time.Duration(reqSec) * time.Second,
	}
}

// Primitive returns the underlying tview primitive for embedding in layouts.
func (p *TimeoutPanel) Primitive() tview.Primitive {
	return p.form
}
