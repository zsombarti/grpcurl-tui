package ui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// SanitizerPanel displays sanitized request payload previews.
type SanitizerPanel struct {
	frame     *tview.Frame
	textView  *tview.TextView
	sanitizer *grpcpkg.RequestSanitizer
}

// NewSanitizerPanel creates a SanitizerPanel with the given sanitizer.
func NewSanitizerPanel(sanitizer *grpcpkg.RequestSanitizer) *SanitizerPanel {
	if sanitizer == nil {
		sanitizer = grpcpkg.NewRequestSanitizer(grpcpkg.DefaultSanitizerPolicy())
	}
	tv := tview.NewTextView()
	tv.SetDynamicColors(true)
	tv.SetScrollable(true)
	tv.SetWrap(true)

	frame := tview.NewFrame(tv).
		SetBorders(1, 1, 1, 1, 2, 2).
		AddText("Sanitized Payload", true, tview.AlignCenter, tview.Styles.PrimaryTextColor)

	return &SanitizerPanel{
		frame:     frame,
		textView:  tv,
		sanitizer: sanitizer,
	}
}

// Primitive returns the root tview primitive for layout embedding.
func (p *SanitizerPanel) Primitive() tview.Primitive { return p.frame }

// Preview sanitizes the given payload and renders it in the panel.
func (p *SanitizerPanel) Preview(payload map[string]interface{}) {
	p.textView.Clear()
	if len(payload) == 0 {
		fmt.Fprint(p.textView, "[grey](empty payload)[-]")
		return
	}
	sanitized := p.sanitizer.Sanitize(payload)
	var sb strings.Builder
	for k, v := range sanitized {
		switch val := v.(type) {
		case string:
			if val == p.sanitizer.Policy().RedactedValue {
				sb.WriteString(fmt.Sprintf("[red]%s[-]: %s\n", k, val))
			} else {
				sb.WriteString(fmt.Sprintf("[green]%s[-]: %s\n", k, val))
			}
		default:
			sb.WriteString(fmt.Sprintf("[green]%s[-]: %v\n", k, v))
		}
	}
	fmt.Fprint(p.textView, sb.String())
}

// Clear resets the panel content.
func (p *SanitizerPanel) Clear() {
	p.textView.Clear()
}
