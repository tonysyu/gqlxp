package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Panel represents a generic panel that can be displayed in the TUI
type Panel interface {
	tea.Model
	SetSize(width, height int)
}

// listPanel wraps a list.Model to implement the Panel interface
type listPanel struct {
	list.Model
}

func newListPanel[T list.Item](choices []T) *listPanel {
	items := make([]list.Item, len(choices))
	for i, choice := range choices {
		items[i] = choice
	}
	return &listPanel{
		Model: list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
}

func (lp *listPanel) Init() tea.Cmd {
	return nil
}

func (lp *listPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	lp.Model, cmd = lp.Model.Update(msg)
	return lp, cmd
}

func (lp *listPanel) SetSize(width, height int) {
	lp.Model.SetWidth(width)
	lp.Model.SetHeight(height)
}

// stringPanel displays a simple string content
type stringPanel struct {
	content string
	width   int
	height  int
}

func newStringPanel(content string) *stringPanel {
	return &stringPanel{content: content}
}

func (sp *stringPanel) Init() tea.Cmd {
	return nil
}

func (sp *stringPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return sp, nil
}

func (sp *stringPanel) View() string {
	style := lipgloss.NewStyle().
		Width(sp.width).
		Height(sp.height).
		Padding(1)
	return style.Render(sp.content)
}

func (sp *stringPanel) SetSize(width, height int) {
	sp.width = width
	sp.height = height
}

