package ui

import (
	"fmt"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// CheckpointPanel displays and manages named request checkpoints.
type CheckpointPanel struct {
	flex        *tview.Flex
	list        *tview.List
	nameInput   *tview.InputField
	checkpointer *grpcpkg.RequestCheckpointer
}

// NewCheckpointPanel creates a new CheckpointPanel backed by the given checkpointer.
// If checkpointer is nil a default one is created.
func NewCheckpointPanel(checkpointer *grpcpkg.RequestCheckpointer) *CheckpointPanel {
	if checkpointer == nil {
		checkpointer = grpcpkg.NewRequestCheckpointer(0)
	}

	nameInput := tview.NewInputField().
		SetLabel("Name: ").
		SetFieldWidth(24)

	list := tview.NewList().ShowSecondaryText(true)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText("[yellow]Checkpoints[-]"), 1, 0, false).
		AddItem(nameInput, 1, 0, true).
		AddItem(list, 0, 1, false)

	p := &CheckpointPanel{
		flex:         flex,
		list:         list,
		nameInput:    nameInput,
		checkpointer: checkpointer,
	}
	return p
}

// Primitive returns the root tview primitive for layout embedding.
func (p *CheckpointPanel) Primitive() tview.Primitive {
	return p.flex
}

// Refresh redraws the list from the current checkpointer state.
func (p *CheckpointPanel) Refresh() {
	p.list.Clear()
	for _, cp := range p.checkpointer.All() {
		name := cp.Name
		secondary := fmt.Sprintf("%s  %s", cp.Method, cp.CreatedAt.Format("15:04:05"))
		p.list.AddItem(name, secondary, 0, nil)
	}
}

// CurrentName returns the value currently typed in the name input field.
func (p *CheckpointPanel) CurrentName() string {
	return p.nameInput.GetText()
}

// SelectedName returns the name of the currently highlighted checkpoint, or
// an empty string if the list is empty.
func (p *CheckpointPanel) SelectedName() string {
	if p.checkpointer.Len() == 0 {
		return ""
	}
	idx := p.list.GetCurrentItem()
	all := p.checkpointer.All()
	if idx < 0 || idx >= len(all) {
		return ""
	}
	return all[idx].Name
}
