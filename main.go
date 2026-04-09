package main

import (
	"log"
	"os"

	"github.com/rivo/tview"
	"grpcurl-tui/ui"
)

func main() {
	app := tview.NewApplication()

	root := ui.NewRootLayout(app)

	if err := app.SetRoot(root, true).EnableMouse(true).Run(); err != nil {
		log.Printf("error running application: %v", err)
		os.Exit(1)
	}
}
