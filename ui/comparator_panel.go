package ui

import (
	"fmt"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// ComparatorPanel provides a UI panel for comparing two request payloads side by side.
type ComparatorPanel struct {
	flex       *tview.Flex
	leftInput  *tview.TextArea
	rightInput *tview.TextArea
	resultView *tview.TextView
	comparator *grpcpkg.RequestComparator
}

// NewComparatorPanel creates and returns a new ComparatorPanel.
func NewComparatorPanel() *ComparatorPanel {
	p := &ComparatorPanel{
		comparator: grpcpkg.NewRequestComparator(),
		leftInput:  tview.NewTextArea(),
		rightInput: tview.NewTextArea(),
		resultView: tview.NewTextView(),
	}

	p.leftInput.SetTitle(" Request A ").SetBorder(true)
	p.rightInput.SetTitle(" Request B ").SetBorder(true)
	p.resultView.SetTitle(" Comparison Result ").SetBorder(true)
	p.resultView.SetDynamicColors(true)

	inputRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(p.leftInput, 0, 1, true).
		AddItem(p.rightInput, 0, 1, false)

	p.flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(inputRow, 0, 3, true).
		AddItem(p.resultView, 5, 0, false)

	return p
}

// Primitive returns the root tview primitive for embedding in a layout.
func (p *ComparatorPanel) Primitive() tview.Primitive {
	return p.flex
}

// Compare runs the comparison between the two input payloads and updates the result view.
func (p *ComparatorPanel) Compare() {
	a := p.leftInput.GetText()
	b := p.rightInput.GetText()

	if a == "" || b == "" {
		p.resultView.SetText("[yellow]Enter JSON in both panels to compare.")
		return
	}

	result, err := p.comparator.Compare(a, b)
	if err != nil {
		p.resultView.SetText(fmt.Sprintf("[red]Error: %s", err.Error()))
		return
	}

	if result.Match {
		p.resultView.SetText("[green]✔ Requests are identical (100% match).")
		return
	}

	text := fmt.Sprintf("[orange]%.0f%% similar — %d difference(s):\n", result.Similarity*100, len(result.Differences))
	for _, d := range result.Differences {
		text += fmt.Sprintf("  [white]• %s\n", d)
	}
	p.resultView.SetText(text)
}

// Clear resets both input areas and the result view.
func (p *ComparatorPanel) Clear() {
	p.leftInput.SetText("", false)
	p.rightInput.SetText("", false)
	p.resultView.Clear()
}
