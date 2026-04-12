package ui

import (
	"strings"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// TaggingPanel provides a UI for viewing and adding tags to a named request.
type TaggingPanel struct {
	store  *grpcpkg.TagStore
	frame  *tview.Frame
	flex   *tview.Flex
	keyIn  *tview.InputField
	tagIn  *tview.InputField
	list   *tview.TextView
}

// NewTaggingPanel creates a TaggingPanel backed by the given TagStore.
func NewTaggingPanel(store *grpcpkg.TagStore) *TaggingPanel {
	if store == nil {
		store = grpcpkg.NewTagStore(20)
	}

	keyIn := tview.NewInputField().
		SetLabel("Request Key: ").
		SetFieldWidth(30)

	tagIn := tview.NewInputField().
		SetLabel("Tags (comma-separated): ").
		SetFieldWidth(40)

	list := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(keyIn, 1, 0, true).
		AddItem(tagIn, 1, 0, false).
		AddItem(list, 0, 1, false)

	frame := tview.NewFrame(flex).
		SetBorders(1, 1, 1, 1, 2, 2).
		AddText("Request Tags", true, tview.AlignCenter, tview.Styles.PrimaryTextColor)

	p := &TaggingPanel{
		store: store,
		frame: frame,
		flex:  flex,
		keyIn: keyIn,
		tagIn: tagIn,
		list:  list,
	}

	tagIn.SetDoneFunc(func(_ tcellKey) { p.commit() })
	return p
}

// Primitive returns the root tview primitive for layout embedding.
func (p *TaggingPanel) Primitive() tview.Primitive { return p.frame }

// commit reads the current input fields and adds the tags to the store.
func (p *TaggingPanel) commit() {
	key := strings.TrimSpace(p.keyIn.GetText())
	raw := p.tagIn.GetText()
	if key == "" || strings.TrimSpace(raw) == "" {
		return
	}
	parts := strings.Split(raw, ",")
	_ = p.store.Add(key, parts...)
	p.tagIn.SetText("")
	p.Refresh()
}

// Refresh redraws the tag list from the store.
func (p *TaggingPanel) Refresh() {
	p.list.Clear()
	for _, k := range p.store.Keys() {
		tags := p.store.Get(k)
		p.list.Write([]byte("[yellow]" + k + "[white]: " + strings.Join(tags, ", ") + "\n"))
	}
}

// tcellKey is an alias kept local to satisfy the SetDoneFunc signature.
type tcellKey = rune
