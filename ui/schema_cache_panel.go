package ui

import (
	"fmt"
	"time"

	"github.com/rivo/tview"

	grpcclient "grpcurl-tui/grpc"
)

// SchemaCachePanel displays cache statistics and allows the user to flush the schema cache.
type SchemaCachePanel struct {
	flex  *tview.Flex
	info  *tview.TextView
	cache *grpcclient.SchemaCache
}

// NewSchemaCachePanel creates a panel bound to the provided SchemaCache.
func NewSchemaCachePanel(cache *grpcclient.SchemaCache) *SchemaCachePanel {
	if cache == nil {
		cache = grpcclient.NewSchemaCache(5 * time.Minute)
	}

	info := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(false)
	info.SetBorder(true).SetTitle(" Schema Cache ")

	flushBtn := tview.NewButton("Flush Cache").
		SetLabelColor(tview.Styles.PrimaryTextColor)

	p := &SchemaCachePanel{
		flex:  tview.NewFlex().SetDirection(tview.FlexRow),
		info:  info,
		cache: cache,
	}

	flushBtn.SetSelectedFunc(func() {
		p.cache.Flush()
		p.Refresh()
	})

	p.flex.AddItem(info, 0, 1, false)
	p.flex.AddItem(flushBtn, 1, 0, true)
	p.flex.SetBorder(true).SetTitle(" Schema Cache Panel ")

	p.Refresh()
	return p
}

// Primitive returns the root tview primitive for embedding in a layout.
func (p *SchemaCachePanel) Primitive() tview.Primitive {
	return p.flex
}

// Refresh updates the info view with the current cache statistics.
func (p *SchemaCachePanel) Refresh() {
	p.info.Clear()
	fmt.Fprintf(p.info, "[yellow]Cached Entries:[white] %d\n", p.cache.Len())
	fmt.Fprintf(p.info, "[yellow]Status:[white]         active\n")
}
