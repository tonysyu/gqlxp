package components

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tonysyu/gqlxp/tui/config"
	"github.com/tonysyu/gqlxp/utils/text"
)

const (
	maxDescriptionHeight = 5 // Maximum height for description (in lines)
)

// OpenPanelMsg is sent when an item should be opened
type OpenPanelMsg struct {
	Panel *Panel
}

// Panel wraps a list.Model to provide panel functionality in the TUI
type Panel struct {
	ListModel         list.Model
	title             string
	description       string
	lastSelectedIndex int               // Track the last selected index to detect changes
	wasFiltering      bool              // Track whether we were in filtering mode to detect exits
	resultType        ListItem          // Virtual item displayed at top
	focusOnResultType bool              // Track whether focus is on result type or list
	isFocused         bool              // Track whether this panel is in focus
	focusedDelegate   list.ItemDelegate // Item delegate for rendering when panel is focused
	blurredDelegate   list.ItemDelegate // Item delegate for rendering when panel is blurred
	wrapperStyle      lipgloss.Style    // Current style of wrapper (focused or blurred)
	styles            config.Styles
	width             int
	height            int
}

// OpenPanelFromItem tries to create an OpenPanelMsg tea.Cmd from an item.
// Return nil if item is not ListItem or Open doesn't return a Panel.
func OpenPanelFromItem(item list.Item) tea.Cmd {
	// FIXME: This should also clear old panels if current item can't be opened
	if listItem, ok := item.(ListItem); ok {
		if newPanel, ok := listItem.OpenPanel(); ok {
			return func() tea.Msg { return OpenPanelMsg{Panel: newPanel} }
		}
	}
	return nil
}

func NewEmptyPanel(content string) *Panel {
	return NewPanel([]ListItem{}, content)
}

func NewPanel[T list.Item](choices []T, title string) *Panel {
	items := make([]list.Item, len(choices))
	for i, choice := range choices {
		items[i] = choice
	}
	blurredItemDelegate := newBlurredItemDelegate()
	m := list.New(items, blurredItemDelegate, 0, 0)
	m.DisableQuitKeybindings()
	m.SetShowTitle(false)
	m.SetShowHelp(false)
	styles := config.DefaultStyles()
	return &Panel{
		ListModel:         m,
		lastSelectedIndex: -1, // Initialize to -1 to trigger opening on first selection
		title:             title,
		blurredDelegate:   blurredItemDelegate,
		focusedDelegate:   list.NewDefaultDelegate(),
		wrapperStyle:      styles.BlurredPanel,
		styles:            styles,
	}
}

func newBlurredItemDelegate() list.ItemDelegate {
	delegate := list.NewDefaultDelegate()
	// Match "Selected" styles to "Normal" styles, since blurred items shouldn't render as selected
	delegate.Styles.SelectedTitle = delegate.Styles.NormalTitle
	delegate.Styles.SelectedDesc = delegate.Styles.NormalDesc
	return delegate
}

func (p *Panel) Init() tea.Cmd {
	return nil
}

func (p *Panel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle navigation when result type is present
	if p.resultType != nil {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch {
			case key.Matches(keyMsg, p.ListModel.KeyMap.CursorDown):
				if p.focusOnResultType {
					// Move from result type to first list item
					p.focusOnResultType = false
					if len(p.ListModel.Items()) > 0 {
						p.ListModel.Select(0)
						p.lastSelectedIndex = 0
						return p, p.OpenSelectedItem()
					}
					return p, nil
				}
				// Otherwise, let list handle it below

			case key.Matches(keyMsg, p.ListModel.KeyMap.CursorUp):
				if !p.focusOnResultType && p.ListModel.Index() == 0 {
					// Move from first list item back to result type
					p.ListModel.Select(-1)
					p.lastSelectedIndex = -1
					p.focusOnResultType = true
					return p, p.OpenSelectedItem()
				}
				// Otherwise, let list handle it below
			}
		}
	}

	// Only update list if focus is on list (or no result type)
	if !p.focusOnResultType {
		var cmd tea.Cmd
		p.ListModel, cmd = p.ListModel.Update(msg)

		// Check if we just exited filtering mode
		// Filtering mode is when the user is actively typing the filter
		isFiltering := p.ListModel.FilterState() == list.Filtering
		exitedFiltering := p.wasFiltering && !isFiltering
		p.wasFiltering = isFiltering

		// Check if selection has changed and auto-open detail panel
		currentIndex := p.ListModel.Index()
		// Refresh if index changed OR if we just exited filtering (since filtering could change
		// the selected item)
		if (currentIndex != p.lastSelectedIndex || exitedFiltering) && currentIndex >= 0 {
			p.lastSelectedIndex = currentIndex
			if openCmd := p.OpenSelectedItem(); openCmd != nil {
				return p, tea.Batch(cmd, openCmd)
			}
		}
		return p, cmd
	}

	return p, nil
}

func (p *Panel) SetSize(width, height int) {
	p.width = width
	p.height = height
}

func (p *Panel) SetTitle(title string) {
	p.title = title
}

func (p *Panel) Title() string {
	return p.title
}

func (p *Panel) SetDescription(description string) {
	p.description = description
}

func (p *Panel) Description() string {
	return p.description
}

func (p *Panel) SetObjectType(item ListItem) {
	p.resultType = item
	// If there are no items in the list, focus on result type; otherwise focus on first list item
	p.focusOnResultType = len(p.ListModel.Items()) == 0
}

// Update items to display with focused style (opposite of SetBlurred)
func (p *Panel) SetFocused() {
	p.wrapperStyle = p.styles.FocusedPanel
	p.ListModel.SetDelegate(p.focusedDelegate)
	p.ListModel.SetShowHelp(true)
	p.isFocused = true
}

// Update items to display with blurred style (opposite of SetFocused)
func (p *Panel) SetBlurred() {
	p.wrapperStyle = p.styles.BlurredPanel
	p.ListModel.SetDelegate(p.blurredDelegate)
	p.ListModel.SetShowHelp(false)
	p.isFocused = false
}

// SelectedItem returns the currently selected item in the list
func (p *Panel) SelectedItem() list.Item {
	if p.focusOnResultType {
		return p.resultType
	}
	return p.ListModel.SelectedItem()
}

func (p *Panel) OpenSelectedItem() tea.Cmd {
	return OpenPanelFromItem(p.SelectedItem())
}

// Items returns the items in the list
func (p *Panel) Items() []list.Item {
	return p.ListModel.Items()
}

// SetItems replaces the items in the list while preserving the current selection index
func (p *Panel) SetItems(items []list.Item) {
	currentIndex := p.ListModel.Index()
	p.ListModel.SetItems(items)
	// Restore selection if it's still valid
	if currentIndex >= 0 && currentIndex < len(items) {
		p.ListModel.Select(currentIndex)
	}
}

// View renders the panel
func (p *Panel) View() string {
	const (
		sectionLabelHeight = 1 // Fixed height for section labels
		emptyLineHeight    = 1 // Fixed height for empty lines between sections
	)

	availableHeight := p.height
	parts := []string{}

	truncatedTitle := text.Truncate(p.Title(), p.width-2*config.PanelTitleHPadding)
	title := p.styles.PanelTitle.Render(truncatedTitle)
	parts = append(parts, title)
	availableHeight -= lipgloss.Height(title)

	if p.Description() != "" {
		desc := text.WrapAndTruncate(p.Description(), p.width, maxDescriptionHeight)

		parts = append(parts, desc)
		availableHeight -= lipgloss.Height(desc)
	}

	// Render result type section if present (fixed heights)
	if p.resultType != nil {
		sectionLabel := p.styles.SectionLabel.Render("Result Type")
		parts = append(parts, "", sectionLabel, "") // Appending empty string adds new lines
		availableHeight -= emptyLineHeight + sectionLabelHeight + emptyLineHeight

		// Render result type with focus indicator
		truncatedResultType := text.Truncate(p.resultType.Title(), p.width-config.ItemLeftPadding)
		if p.focusOnResultType && p.isFocused {
			truncatedResultType = p.styles.FocusedItem.Render(truncatedResultType)
		} else {
			truncatedResultType = p.styles.UnfocusedItem.Render(truncatedResultType)
		}
		parts = append(parts, truncatedResultType)
		availableHeight -= lipgloss.Height(truncatedResultType)
	}

	// Input Arguments section label if list has items (fixed heights)
	if len(p.ListModel.Items()) > 0 {
		if p.resultType != nil {
			// Only render sectionLabel if it needs to be differentiated from "Result Type"
			sectionLabel := p.styles.SectionLabel.Render("Input Arguments")
			parts = append(parts, "", sectionLabel) // Appending empty string adds new lines
			availableHeight -= emptyLineHeight + sectionLabelHeight
		}
		p.ListModel.SetWidth(p.width)
		p.ListModel.SetHeight(availableHeight)
		parts = append(parts, p.ListModel.View())
	}

	content := text.JoinLines(parts...)

	// Apply fixed width & height style to avoid jittery panel borders
	style := lipgloss.NewStyle().Width(p.width).Height(p.height)
	innerPanel := style.Render(content)
	return p.wrapperStyle.Render(innerPanel)
}
