package ui

import (
	"fmt"
	"strings"
	"sync"

	"github.com/rivo/tview"
)

// StreamPanel displays live messages arriving from a server-streaming RPC.
type StreamPanel struct {
	*tview.TextView
	mu    sync.Mutex
	lines []string
}

// NewStreamPanel creates a StreamPanel ready to display streaming responses.
func NewStreamPanel() *StreamPanel {
	tv := tview.NewTextView()
	tv.SetDynamicColors(true)
	tv.SetScrollable(true)
	tv.SetBorder(true)
	tv.SetTitle(" Stream Output ")
	tv.SetTitleAlign(tview.AlignLeft)

	return &StreamPanel{
		TextView: tv,
		lines:    make([]string, 0),
	}
}

// AppendMessage adds a new JSON message line to the panel.
func (s *StreamPanel) AppendMessage(index int, msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	line := fmt.Sprintf("[green][%d][white] %s", index, msg)
	s.lines = append(s.lines, line)
	s.TextView.SetText(strings.Join(s.lines, "\n"))
	s.TextView.ScrollToEnd()
}

// AppendError appends a red-coloured error line to the panel.
func (s *StreamPanel) AppendError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	line := fmt.Sprintf("[red]ERROR: %s[white]", err.Error())
	s.lines = append(s.lines, line)
	s.TextView.SetText(strings.Join(s.lines, "\n"))
	s.TextView.ScrollToEnd()
}

// Clear resets the panel content.
func (s *StreamPanel) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lines = s.lines[:0]
	s.TextView.SetText("")
}

// Len returns the number of messages currently displayed.
func (s *StreamPanel) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.lines)
}
