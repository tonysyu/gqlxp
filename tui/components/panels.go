package components

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// OpenPanelMsg is sent when an item should be opened
type OpenPanelMsg struct {
	Panel Panel
}

var _ Panel = (*ListPanel)(nil)
var _ Panel = (*stringPanel)(nil)

// ListPanel wraps a list.Model to implement the Panel interface
type ListPanel struct {
	model             list.Model
	lastSelectedIndex int // Track the last selected index to detect changes
}

func NewListPanel[T list.Item](choices []T, title string) *ListPanel {
	items := make([]list.Item, len(choices))
	for i, choice := range choices {
		items[i] = choice
	}
	m := list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.DisableQuitKeybindings()
	m.Title = title
	return &ListPanel{
		model:             m,
		lastSelectedIndex: -1, // Initialize to -1 to trigger opening on first selection
	}
}

func (lp *ListPanel) Init() tea.Cmd {
	return nil
}

func (lp *ListPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Update the list model first
	lp.model, cmd = lp.model.Update(msg)

	// Check if selection has changed and auto-open detail panel
	currentIndex := lp.model.Index()
	if currentIndex != lp.lastSelectedIndex && currentIndex >= 0 {
		lp.lastSelectedIndex = currentIndex
		if selectedItem := lp.model.SelectedItem(); selectedItem != nil {
			if listItem, ok := selectedItem.(ListItem); ok {
				if newPanel, ok := listItem.Open(); ok {
					return lp, tea.Batch(cmd, func() tea.Msg {
						return OpenPanelMsg{Panel: newPanel}
					})
				}
			}
		}
	}

	return lp, cmd
}

func (lp *ListPanel) SetSize(width, height int) {
	lp.model.SetWidth(width)
	lp.model.SetHeight(height)
}

func (lp *ListPanel) SetTitle(title string) {
	lp.model.Title = title
}

// SelectedItem returns the currently selected item in the list
func (lp *ListPanel) SelectedItem() list.Item {
	return lp.model.SelectedItem()
}

// Title returns the title of the list panel
func (lp *ListPanel) Title() string {
	return lp.model.Title
}

// Items returns the items in the list
func (lp *ListPanel) Items() []list.Item {
	return lp.model.Items()
}

// View renders the list panel
func (lp *ListPanel) View() string {
	return lp.model.View()
}

// stringPanel displays a simple string content
type stringPanel struct {
	content string
	width   int
	height  int
}

func NewStringPanel(content string) *stringPanel {
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
