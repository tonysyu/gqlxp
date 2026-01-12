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

// Tab represents a labeled tab with associated content
type Tab struct {
	Label   string
	Content []ListItem
}

type keymap = struct {
	nextTab, prevTab key.Binding
}

// Panel wraps a list.Model to provide panel functionality in the TUI
type Panel struct {
	ListModel         list.Model
	title             string
	description       string
	lastSelectedIndex int               // Track the last selected index to detect changes
	wasFiltering      bool              // Track whether we were in filtering mode to detect exits
	tabs              []Tab             // Tabs with labels and content
	activeTab         int               // Currently selected tab index
	isFocused         bool              // Track whether this panel is in focus
	focusedDelegate   list.ItemDelegate // Item delegate for rendering when panel is focused
	blurredDelegate   list.ItemDelegate // Item delegate for rendering when panel is blurred
	wrapperStyle      lipgloss.Style    // Current style of wrapper (focused or blurred)
	keymap            keymap
	styles            config.Styles
	width             int
	height            int
}

// OpenPanelFromItem tries to create an OpenPanelMsg tea.Cmd from an item.
// Return nil if item is not ListItem or Open doesn't return a Panel.
func OpenPanelFromItem(item list.Item) tea.Cmd {
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
	m.SetShowStatusBar(false)
	styles := config.DefaultStyles()
	panel := &Panel{
		ListModel:         m,
		lastSelectedIndex: -1, // Initialize to -1 to trigger opening on first selection
		title:             title,
		blurredDelegate:   blurredItemDelegate,
		focusedDelegate:   list.NewDefaultDelegate(),
		wrapperStyle:      styles.BlurredPanel,
		// Only include single keymaps for short help (overridden for full help)
		keymap: keymap{
			nextTab: key.NewBinding(
				key.WithKeys("L", "shift+right"),
				key.WithHelp("L", "next tab"),
			),
			prevTab: key.NewBinding(
				key.WithKeys("H", "shift+left"),
				key.WithHelp("H", "prev tab"),
			),
		},
		styles: styles,
	}
	panel.configureTabHelp()
	return panel
}

func newBlurredItemDelegate() list.ItemDelegate {
	delegate := list.NewDefaultDelegate()
	// Match "Selected" styles to "Normal" styles, since blurred items shouldn't render as selected
	delegate.Styles.SelectedTitle = delegate.Styles.NormalTitle
	delegate.Styles.SelectedDesc = delegate.Styles.NormalDesc
	return delegate
}

// configureTabHelp sets up the help display for tab navigation
func (p *Panel) configureTabHelp() {
	p.ListModel.AdditionalShortHelpKeys = func() []key.Binding {
		// Only show tab navigation help when there are multiple tabs
		if len(p.tabs) <= 1 {
			return nil
		}
		return []key.Binding{p.keymap.nextTab, p.keymap.prevTab}
	}
	p.ListModel.AdditionalFullHelpKeys = func() []key.Binding {
		// Only show tab navigation help when there are multiple tabs
		if len(p.tabs) <= 1 {
			return nil
		}
		// Include additional keymaps when displaying full help
		return []key.Binding{
			key.NewBinding(
				key.WithKeys(p.keymap.nextTab.Keys()...),
				key.WithHelp("L/⇧+→", p.keymap.nextTab.Help().Desc),
			),
			key.NewBinding(
				key.WithKeys(p.keymap.prevTab.Keys()...),
				key.WithHelp("H/⇧+←", p.keymap.prevTab.Help().Desc),
			),
		}
	}
}

func (p *Panel) Init() tea.Cmd {
	return nil
}

func (p *Panel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle tab navigation with Shift-H (previous) and Shift-L (next)
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if len(p.tabs) > 1 {
			switch {
			case key.Matches(keyMsg, p.keymap.prevTab):
				if p.activeTab > 0 {
					p.activeTab--
					p.switchToActiveTab()
					return p, p.OpenSelectedItem()
				}
				return p, nil
			case key.Matches(keyMsg, p.keymap.nextTab):
				if p.activeTab < len(p.tabs)-1 {
					p.activeTab++
					p.switchToActiveTab()
					return p, p.OpenSelectedItem()
				}
				return p, nil
			}
		}
	}

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

// switchToActiveTab updates the list content to show the active tab's items
func (p *Panel) switchToActiveTab() {
	if p.activeTab >= 0 && p.activeTab < len(p.tabs) {
		// Convert []ListItem to []list.Item
		content := p.tabs[p.activeTab].Content
		items := make([]list.Item, len(content))
		for i, item := range content {
			items[i] = item
		}
		p.ListModel.SetItems(items)
		p.ListModel.Select(0)
		p.lastSelectedIndex = 0
	}
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

// SetTabs configures the panel with multiple tabs
func (p *Panel) SetTabs(tabs []Tab) {
	p.tabs = tabs
	p.activeTab = 0
	if len(tabs) > 0 {
		p.switchToActiveTab()
	}
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

// SelectItemByName selects an item by its RefName
// Returns true if found and selected, false otherwise
func (p *Panel) SelectItemByName(name string) bool {
	for i, item := range p.ListModel.Items() {
		if listItem, ok := item.(ListItem); ok {
			if listItem.RefName() == name {
				p.ListModel.Select(i)
				p.lastSelectedIndex = i
				return true
			}
		}
	}
	return false
}

// View renders the panel
func (p *Panel) View() string {
	availableHeight := p.height
	parts := []string{}

	// Render title
	truncatedTitle := text.Truncate(p.Title(), p.width-2*config.PanelTitleHPadding)
	title := p.styles.PanelTitle.Render(truncatedTitle)
	parts = append(parts, title)
	availableHeight -= lipgloss.Height(title)

	// Render description if present
	if p.Description() != "" {
		desc := text.WrapAndTruncate(p.Description(), p.width, maxDescriptionHeight)
		parts = append(parts, desc)
		availableHeight -= lipgloss.Height(desc)
	}

	// Render tabs if configured (even if only one tab for consistency)
	if len(p.tabs) > 0 {
		tabBar := p.renderTabBar()
		parts = append(parts, "", tabBar)
		availableHeight -= 1 + lipgloss.Height(tabBar) // 1 for empty line
	}

	// Render list content
	if len(p.ListModel.Items()) > 0 {
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

// renderTabBar creates a tab bar display with active/inactive styling
func (p *Panel) renderTabBar() string {
	if len(p.tabs) == 0 {
		return ""
	}

	// Calculate max width per tab to fit all tabs in available width
	// Account for padding in tab styles by using a small buffer
	const styleOverhead = 2 // Approximate horizontal padding per tab
	availableWidth := max(p.width - (len(p.tabs) * styleOverhead), len(p.tabs))
	maxTabWidth := availableWidth / len(p.tabs)

	var tabParts []string
	for i, tab := range p.tabs {
		// Truncate label to fit within calculated max width
		label := text.Truncate(tab.Label, maxTabWidth)
		if i == p.activeTab {
			tabParts = append(tabParts, p.styles.ActiveSubTab.Render(label))
		} else {
			tabParts = append(tabParts, p.styles.InactiveTab.Render(label))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, tabParts...)
}
