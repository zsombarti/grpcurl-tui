package ui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// RouterPanel provides a UI for managing request routing rules.
type RouterPanel struct {
	flex    *tview.Flex
	router  *grpcpkg.RequestRouter
	table   *tview.Table
	fallbackInput *tview.InputField
}

// NewRouterPanel creates a new RouterPanel backed by the given RequestRouter.
func NewRouterPanel(router *grpcpkg.RequestRouter) *RouterPanel {
	table := tview.NewTable().SetBorders(true)
	table.SetTitle(" Routing Rules ").SetBorder(true)

	fallbackInput := tview.NewInputField().
		SetLabel("Fallback: ").
		SetText(router.Fallback()).
		SetFieldWidth(30)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(fallbackInput, 3, 0, false).
		AddItem(table, 0, 1, true)

	p := &RouterPanel{
		flex:          flex,
		router:        router,
		table:         table,
		fallbackInput: fallbackInput,
	}

	fallbackInput.SetDoneFunc(func(_ tcell.Key) {
		p.applyFallback()
	})

	p.Refresh()
	return p
}

// Primitive returns the tview primitive for embedding.
func (p *RouterPanel) Primitive() tview.Primitive {
	return p.flex
}

// Refresh redraws the routing rules table.
func (p *RouterPanel) Refresh() {
	p.table.Clear()
	p.table.SetCell(0, 0, tview.NewTableCell("Method").SetSelectable(false))
	p.table.SetCell(0, 1, tview.NewTableCell("Address").SetSelectable(false))

	rules := p.routeSnapshot()
	for i, rule := range rules {
		p.table.SetCell(i+1, 0, tview.NewTableCell(rule[0]))
		p.table.SetCell(i+1, 1, tview.NewTableCell(rule[1]))
	}

	fallback := p.router.Fallback()
	if fallback == "" {
		fallback = "(none)"
	}
	p.fallbackInput.SetText(fallback)
}

func (p *RouterPanel) applyFallback() {
	val := strings.TrimSpace(p.fallbackInput.GetText())
	if val != "" && val != "(none)" {
		p.router.SetFallback(val)
	}
}

func (p *RouterPanel) routeSnapshot() [][2]string {
	// Build a display list by probing known entries via Len.
	// Since RouteRule is exported, we collect them via a helper.
	var rows [][2]string
	for i := 0; i < p.router.Len(); i++ {
		rows = append(rows, [2]string{fmt.Sprintf("rule-%d", i), ""})
	}
	return rows
}
