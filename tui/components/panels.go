package components

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/tonysyu/gqlxp/tui/config"
	"github.com/tonysyu/gqlxp/utils/text"
)

// Panel represents a generic panel that can be displayed in the TUI
type Panel interface {
	tea.Model
	SetSize(width, height int)
}

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
	title             string
	description       string
	resultType        ListItem // Virtual item displayed at top
	focusOnResultType bool     // Track whether focus is on result type or list
	styles            config.Styles
	width             int
	height            int
}

// OpenPanelFromItem tries to create an OpenPanelMsg tea.Cmd from an item.
// Return nil if item is not ListItem or Open doesn't return a Panel.
func OpenPanelFromItem(item list.Item) tea.Cmd {
	if listItem, ok := item.(ListItem); ok {
		if newPanel, ok := listItem.Open(); ok {
			return func() tea.Msg { return OpenPanelMsg{Panel: newPanel} }
		}
	}
	return nil
}

func NewListPanel[T list.Item](choices []T, title string) *ListPanel {
	items := make([]list.Item, len(choices))
	for i, choice := range choices {
		items[i] = choice
	}
	m := list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.DisableQuitKeybindings()
	m.SetShowTitle(false)
	return &ListPanel{
		model:             m,
		lastSelectedIndex: -1, // Initialize to -1 to trigger opening on first selection
		title:             title,
		styles:            config.DefaultStyles(),
	}
}

func (lp *ListPanel) Init() tea.Cmd {
	return nil
}

func (lp *ListPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle navigation when result type is present
	if lp.resultType != nil {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch {
			case key.Matches(keyMsg, lp.model.KeyMap.CursorDown):
				if lp.focusOnResultType {
					// Move from result type to first list item
					lp.focusOnResultType = false
					if len(lp.model.Items()) > 0 {
						lp.model.Select(0)
						lp.lastSelectedIndex = 0
						return lp, lp.OpenSelectedItem()
					}
					return lp, nil
				}
				// Otherwise, let list handle it below

			case key.Matches(keyMsg, lp.model.KeyMap.CursorUp):
				if !lp.focusOnResultType && lp.model.Index() == 0 {
					// Move from first list item back to result type
					lp.model.Select(-1)
					lp.lastSelectedIndex = -1
					lp.focusOnResultType = true
					return lp, lp.OpenSelectedItem()
				}
				// Otherwise, let list handle it below
			}
		}
	}

	// Only update list if focus is on list (or no result type)
	if !lp.focusOnResultType {
		var cmd tea.Cmd
		lp.model, cmd = lp.model.Update(msg)
		// Check if selection has changed and auto-open detail panel
		currentIndex := lp.model.Index()
		if currentIndex != lp.lastSelectedIndex && currentIndex >= 0 {
			lp.lastSelectedIndex = currentIndex
			if openCmd := lp.OpenSelectedItem(); openCmd != nil {
				return lp, tea.Batch(cmd, openCmd)
			}
		}
		return lp, cmd
	}

	return lp, nil
}

func (lp *ListPanel) SetSize(width, height int) {
	lp.width = width
	lp.height = height
}

func (lp *ListPanel) SetTitle(title string) {
	lp.title = title
}

func (lp *ListPanel) Title() string {
	return lp.title
}

func (lp *ListPanel) SetDescription(description string) {
	lp.description = description
}

func (lp *ListPanel) Description() string {
	return lp.description
}

func (lp *ListPanel) SetResultType(item ListItem) {
	lp.resultType = item
	lp.focusOnResultType = true // Start with focus on result type
}

// SelectedItem returns the currently selected item in the list
func (lp *ListPanel) SelectedItem() list.Item {
	if lp.focusOnResultType {
		return lp.resultType
	}
	return lp.model.SelectedItem()
}

func (lp *ListPanel) OpenSelectedItem() tea.Cmd {
	return OpenPanelFromItem(lp.SelectedItem())
}

// Items returns the items in the list
func (lp *ListPanel) Items() []list.Item {
	return lp.model.Items()
}

// View renders the list panel
func (lp *ListPanel) View() string {
	availableHeight := lp.height
	parts := []string{}

	title := lp.styles.PanelTitle.Render(lp.Title())
	parts = append(parts, title)
	availableHeight -= lipgloss.Height(title)

	if lp.Description() != "" {
		desc := wordwrap.String(lp.Description(), lp.width)
		parts = append(parts, desc)
		availableHeight -= lipgloss.Height(desc)
	}

	// Render result type section if present
	if lp.resultType != nil {
		sectionLabel := lp.styles.SectionLabel.Render("Result Type")
		parts = append(parts, "", sectionLabel, "") // Appending empty string adds new lines
		availableHeight -= lipgloss.Height(sectionLabel)

		// Render result type with focus indicator
		resultTypeText := lp.resultType.Title()
		if lp.focusOnResultType {
			resultTypeText = lp.styles.FocusedItem.Render(resultTypeText)
		} else {
			resultTypeText = lp.styles.UnfocusedItem.Render(" " + resultTypeText)
		}
		parts = append(parts, resultTypeText)
		availableHeight -= lipgloss.Height(resultTypeText)
	}

	// Input Arguments section label if list has items
	if len(lp.model.Items()) > 0 {
		sectionLabel := lp.styles.SectionLabel.Render("Input Arguments")
		parts = append(parts, "", sectionLabel) // Appending empty string adds new lines
		availableHeight -= lipgloss.Height(sectionLabel)
		lp.model.SetWidth(lp.width)
		lp.model.SetHeight(availableHeight)
		parts = append(parts, lp.model.View())
	}

	return text.JoinLines(parts...)
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
