package ui

import (
	"fmt"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// SearchPanel provides an interactive search box over indexed proto services/methods.
type SearchPanel struct {
	flex     *tview.Flex
	input    *tview.InputField
	results  *tview.List
	searcher *grpcpkg.ProtoSearcher
}

// NewSearchPanel creates and wires up the search panel UI.
func NewSearchPanel(searcher *grpcpkg.ProtoSearcher) *SearchPanel {
	p := &SearchPanel{
		searcher: searcher,
		input:    tview.NewInputField(),
		results:  tview.NewList(),
		flex:     tview.NewFlex(),
	}

	p.input.SetLabel("Search: ").
		SetFieldWidth(40).
		SetChangedFunc(func(text string) {
			p.refresh(text)
		})

	p.results.ShowSecondaryText(true).
		SetBorder(false)

	p.flex.SetDirection(tview.FlexRow).
		SetBorder(true).
		SetTitle(" Proto Search ").
		AddItem(p.input, 1, 0, true).
		AddItem(p.results, 0, 1, false)

	p.refresh("")
	return p
}

// Primitive returns the root tview primitive for layout embedding.
func (p *SearchPanel) Primitive() tview.Primitive {
	return p.flex
}

// Reindex re-populates the searcher and refreshes the displayed results.
func (p *SearchPanel) Reindex(services map[string][]string) {
	p.searcher.Index(services)
	p.refresh(p.input.GetText())
}

func (p *SearchPanel) refresh(query string) {
	p.results.Clear()
	matches := p.searcher.Search(query)
	for _, r := range matches {
		secondary := r.Service
		if r.Method == "" {
			secondary = "(service)"
		}
		p.results.AddItem(r.Full, fmt.Sprintf("  %s", secondary), 0, nil)
	}
}
