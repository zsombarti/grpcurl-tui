package ui

import (
	"github.com/rivo/tview"
)

// RootLayout holds the top-level TUI layout with a header, sidebar, and main panel.
type RootLayout struct {
	*tview.Flex
	app    *tview.Application
	header *tview.TextView
	sidebar *tview.List
	main   *tview.TextView
}

// NewRootLayout constructs and wires up the root TUI layout.
func NewRootLayout(app *tview.Application) *tview.Flex {
	r := &RootLayout{app: app}

	r.header = tview.NewTextView().
		SetText(" grpcurl-tui — Interactive gRPC Explorer").
		SetTextColor(tview.Styles.PrimaryTextColor).
		SetDynamicColors(true)
	r.header.SetBackgroundColor(tview.Styles.ContrastBackgroundColor)

	r.sidebar = tview.NewList().
		AddItem("Connect", "Set gRPC server address", 'c', nil).
		AddItem("Services", "Browse reflected services", 's', nil).
		AddItem("History", "View past requests", 'h', nil).
		AddItem("Quit", "Exit application", 'q', func() {
			app.Stop()
		})
	r.sidebar.SetBorder(true).SetTitle(" Menu ")

	r.main = tview.NewTextView().
		SetText("Welcome to grpcurl-tui!\n\nUse the menu on the left to get started.\n\n" +
			"Steps:\n  1. Connect to a gRPC server\n  2. Browse available services\n  3. Select a method and send requests").
		SetDynamicColors(true).
		SetWordWrap(true)
	r.main.SetBorder(true).SetTitle(" Output ")

	body := tview.NewFlex().
		AddItem(r.sidebar, 24, 0, true).
		AddItem(r.main, 0, 1, false)

	root := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(r.header, 1, 0, false).
		AddItem(body, 0, 1, true)

	return root
}
