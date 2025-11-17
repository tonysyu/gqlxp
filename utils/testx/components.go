package testx

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tonysyu/gqlxp/tui/xplr/components"
)

// RenderMinimalPanel drastically simplifies Panel rendering to create simpler tests
//
// In particular, this will create a panel with the following characteristics:
// 1. Empty lines removed (using NormalizeView)
// 2. Item only show title (no description)
// 3. No selection indicator
// 4. No "status bar" (item count)
// 4. No help
func RenderMinimalPanel(panel *components.Panel) string {
	panel.ListModel.SetDelegate(minimalItemDelegate{})
	panel.ListModel.SetShowStatusBar(false)
	panel.ListModel.ShowHelp()
	content := NormalizeView(panel.View())
	return content
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
