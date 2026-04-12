package ui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

// EnricherPanel displays registered enrichment steps and allows the user
// to inspect which steps will be applied to outgoing requests.
type EnricherPanel struct {
	flex  *tview.Flex
	list  *tview.List
	status *tview.TextView
	steps []string
}

// NewEnricherPanel creates a new EnricherPanel.
func NewEnricherPanel() *EnricherPanel {
	p := &EnricherPanel{
		flex:   tview.NewFlex(),
		list:   tview.NewList(),
		status: tview.NewTextView(),
	}

	p.list.SetBorder(true).SetTitle(" Enrichment Steps ")
	p.list.ShowSecondaryText(false)

	p.status.SetBorder(false)
	p.status.SetTextColor(tview.Styles.SecondaryTextColor)
	p.status.SetText("No enrichment steps registered.")

	p.flex.SetDirection(tview.FlexRow).
		AddItem(p.list, 0, 1, true).
		AddItem(p.status, 1, 0, false)

	return p
}

// Primitive returns the root tview primitive for layout embedding.
func (p *EnricherPanel) Primitive() tview.Primitive {
	return p.flex
}

// Refresh updates the list with the provided step names.
func (p *EnricherPanel) Refresh(stepNames []string) {
	p.list.Clear()
	p.steps = make([]string, len(stepNames))
	copy(p.steps, stepNames)

	for i, name := range stepNames {
		label := fmt.Sprintf("[%d] %s", i+1, name)
		p.list.AddItem(label, "", 0, nil)
	}

	if len(stepNames) == 0 {
		p.status.SetText("No enrichment steps registered.")
	} else {
		p.status.SetText(fmt.Sprintf("%d step(s) active: %s",
			len(stepNames), strings.Join(stepNames, ", ")))
	}
}

// StepCount returns the number of steps currently displayed.
func (p *EnricherPanel) StepCount() int {
	return len(p.steps)
}
