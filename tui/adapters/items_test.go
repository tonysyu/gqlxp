package adapters

import (
	"fmt"
	"io"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/tui/xplr/components"
	"github.com/tonysyu/gqlxp/utils/testx"
)

// renderMinimalPanel drastically simplifies Panel rendering to create simpler tests
//
// In particular, this will create a panel with the following characteristics:
// 1. Empty lines removed (using NormalizeView)
// 2. Item only show title (no description)
// 3. No selection indicator
// 4. No "status bar" (item count)
// 4. No help
func renderMinimalPanel(panel *components.Panel) string {
	panel.ListModel.SetDelegate(minimalItemDelegate{})
	panel.ListModel.SetShowStatusBar(false)
	panel.ListModel.ShowHelp()
	content := testx.NormalizeView(panel.View())
	return content
}

func nextPanelTab(panel *components.Panel) *components.Panel {
	panel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'L'}})
	return panel
}

type minimalItemDelegate struct{}

func (d minimalItemDelegate) Height() int                             { return 1 }
func (d minimalItemDelegate) Spacing() int                            { return 0 }
func (d minimalItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d minimalItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(list.DefaultItem)
	if !ok {
		return
	}
	fmt.Fprint(w, i.Title())
}

func TestSimpleItemInterface(t *testing.T) {
	is := is.New(t)

	item := components.NewSimpleItem("Test Title", components.WithDescription("Test Description"))

	is.Equal(item.Title(), "Test Title")
	is.Equal(item.Description(), "Test Description")
	is.Equal(item.FilterValue(), "Test Title")

	// Simple items should not be openable
	panel, ok := item.OpenPanel()
	is.True(!ok)
	is.Equal(panel, nil)
}
