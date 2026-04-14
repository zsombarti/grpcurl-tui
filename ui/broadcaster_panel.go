package ui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"

	"grpcurl-tui/grpc"
)

// BroadcasterPanel displays registered broadcast targets and the latest
// broadcast results inside a tview Flex layout.
type BroadcasterPanel struct {
	flex        *tview.Flex
	targetView  *tview.TextView
	resultView  *tview.TextView
	broadcaster *grpc.RequestBroadcaster
}

// NewBroadcasterPanel creates a BroadcasterPanel backed by the given
// RequestBroadcaster. If broadcaster is nil a default one is created.
func NewBroadcasterPanel(broadcaster *grpc.RequestBroadcaster) *BroadcasterPanel {
	if broadcaster == nil {
		broadcaster = grpc.NewRequestBroadcaster(grpc.DefaultBroadcasterPolicy())
	}

	targetView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	targetView.SetBorder(true).SetTitle(" Targets ")

	resultView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	resultView.SetBorder(true).SetTitle(" Last Broadcast Results ")

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(targetView, 0, 1, false).
		AddItem(resultView, 0, 2, false)

	return &BroadcasterPanel{
		flex:        flex,
		targetView:  targetView,
		resultView:  resultView,
		broadcaster: broadcaster,
	}
}

// Primitive returns the root tview primitive for embedding in a layout.
func (p *BroadcasterPanel) Primitive() tview.Primitive { return p.flex }

// RefreshTargets re-renders the target list from the broadcaster.
func (p *BroadcasterPanel) RefreshTargets(targets []string) {
	p.targetView.Clear()
	if len(targets) == 0 {
		fmt.Fprint(p.targetView, "[grey]No targets registered[-]")
		return
	}
	for i, t := range targets {
		fmt.Fprintf(p.targetView, "[green]%d.[-] %s\n", i+1, t)
	}
}

// ShowResults renders the broadcast results in the result view.
func (p *BroadcasterPanel) ShowResults(results []grpc.BroadcastResult) {
	p.resultView.Clear()
	if len(results) == 0 {
		fmt.Fprint(p.resultView, "[grey]No results yet[-]")
		return
	}
	var sb strings.Builder
	for _, r := range results {
		if r.Err != nil {
			sb.WriteString(fmt.Sprintf("[red]%-20s ERR: %v[-]\n", r.Target, r.Err))
		} else {
			sb.WriteString(fmt.Sprintf("[green]%-20s OK  keys=%d[-]\n", r.Target, len(r.Payload)))
		}
	}
	fmt.Fprint(p.resultView, sb.String())
}
