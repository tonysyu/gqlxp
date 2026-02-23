package cmdpalette

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tonysyu/gqlxp/tui/config"
	"github.com/tonysyu/gqlxp/tui/utils"
	"github.com/tonysyu/gqlxp/utils/terminal"
)

// commandItem represents a command in the palette
type commandItem struct {
	context string      // "Main", "Panel", "Overlay", "Global"
	title   string      // Help text from key binding
	key     string      // Key combination that triggers this command
	binding key.Binding // The actual key binding
	enabled bool        // Whether command is active in current context
}

func newCommandItem(context, title string, binding key.Binding) commandItem {
	return commandItem{
		context: context,
		title:   title,
		key:     binding.Help().Key,
		binding: binding,
	}
}

// Implement list.Item interface
func (i commandItem) FilterValue() string {
	return fmt.Sprintf("%s %s %s", i.title, i.context, i.Description())
}

// Implement list.DefaultItem interface
func (i commandItem) Title() string {
	return i.title
}

func (i commandItem) Description() string {
	return i.binding.Help().Key
}

func (i commandItem) Context() string {
	return i.context
}

// commandDelegate is a custom delegate for rendering command items
type commandDelegate struct {
	styles config.Styles
}

func (d commandDelegate) Height() int { return 2 }

func (d commandDelegate) Spacing() int { return 1 }

func (d commandDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d commandDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	cmd, ok := item.(commandItem)
	if !ok {
		return
	}

	// Apply dimmed styling if command is not enabled
	baseStyle := lipgloss.NewStyle().Foreground(terminal.ColorDimWhite)
	if !cmd.enabled {
		baseStyle = baseStyle.Foreground(terminal.ColorMidGray)
	}

	// Highlight selected item
	if index == m.Index() {
		baseStyle = baseStyle.Bold(true).Foreground(terminal.ColorDimMagenta)
		if !cmd.enabled {
			baseStyle = baseStyle.Bold(false).Foreground(terminal.ColorDimIndigo)
		}
	}

	title := cmd.Title()
	context := cmd.Context()
	keyCombo := cmd.Description()

	// Format title: "[title]     [context]"
	// Calculate spacing to align context on the right
	const maxWidth = 80
	descWidth := lipgloss.Width(title)
	contextWidth := lipgloss.Width(context)
	spacing := max(maxWidth-descWidth-contextWidth, 2)

	titleLine := fmt.Sprintf("%s%s%s",
		title,
		lipgloss.NewStyle().Width(spacing).Render(""),
		context)

	// Render as two lines: title on top, key combination below
	output := baseStyle.Render(titleLine) + "\n" + baseStyle.Faint(true).Render(keyCombo)
	fmt.Fprint(w, output)
}

// ClosedMsg is sent when the command palette requests to be closed
type ClosedMsg struct{}

// Model manages the command palette
type Model struct {
	list    list.Model
	styles  config.Styles
	keymaps config.CommandPaletteKeymaps
	width   int
	height  int
}

// New creates a new command palette model
func New(styles config.Styles, paletteKeymaps config.CommandPaletteKeymaps, commandKeymaps CommandKeymaps) Model {
	items := buildCommandItems(commandKeymaps)

	delegate := commandDelegate{styles: styles}
	l := list.New(items, delegate, 0, 0)
	l.Title = "Command Palette"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.DisableQuitKeybindings()
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			paletteKeymaps.Execute,
			paletteKeymaps.Close,
			paletteKeymaps.Quit,
		}
	}
	l.AdditionalFullHelpKeys = l.AdditionalShortHelpKeys

	return Model{
		list:    l,
		styles:  styles,
		keymaps: paletteKeymaps,
	}
}

// CommandKeymaps contains all keymaps needed to build the command palette
type CommandKeymaps struct {
	Main    config.MainKeymaps
	Panel   config.PanelKeymaps
	Overlay config.OverlayKeymaps
}

// buildCommandItems constructs command items from all keymaps
func buildCommandItems(keymaps CommandKeymaps) []list.Item {
	var items []list.Item

	// Add main keymaps
	items = append(items, newCommandItem(
		"Main",
		"Go to next panel",
		keymaps.Main.NextPanel,
	))
	items = append(items, newCommandItem(
		"Main",
		"Go to previous panel",
		keymaps.Main.PrevPanel,
	))
	items = append(items, newCommandItem(

		"Main",
		"Go to next GQL type class",
		keymaps.Main.NextGQLType,
	))
	items = append(items, newCommandItem(
		"Main",
		"Go to previous GQL type class",
		keymaps.Main.PrevGQLType,
	))
	items = append(items, newCommandItem(
		"Main",
		"Show detail overlay of focused type",
		keymaps.Main.ToggleOverlay,
	))

	// Add panel keymaps
	items = append(items, newCommandItem(
		"Panel",
		"Go to next panel subtab",
		keymaps.Panel.NextTab,
	))
	items = append(items, newCommandItem(
		"Panel",
		"Go to previous panel subtab",
		keymaps.Panel.PrevTab,
	))

	// Add search keymaps
	items = append(items, newCommandItem(
		"Search",
		"Focus on search input",
		keymaps.Main.SearchFocus,
	))
	items = append(items, newCommandItem(
		"Search",
		"Submit search query",
		keymaps.Main.SearchSubmit,
	))
	items = append(items, newCommandItem(
		"Search",
		"Clear search query",
		keymaps.Main.SearchClear,
	))

	// Add overlay keymaps
	items = append(items, newCommandItem(
		"Overlay",
		"Close detail overlay",
		keymaps.Overlay.Close,
	))

	// Add global keymaps
	items = append(items, newCommandItem(
		"Global",
		"Quit gqlxp app",
		keymaps.Main.Quit,
	))

	return items
}

// Update processes messages and returns (model, cmd).
// xplr.Model is responsible for routing messages here only when the palette is active.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymaps.Close):
			return m, func() tea.Msg { return ClosedMsg{} }
		case key.Matches(msg, m.keymaps.Execute):
			selectedItem := m.list.SelectedItem()
			cmd, ok := selectedItem.(commandItem)
			if ok && cmd.enabled {
				keyMsg := parseKeyString(cmd.key)
				return m, tea.Sequence(
					func() tea.Msg { return ClosedMsg{} },
					func() tea.Msg { return keyMsg },
				)
			}
			return m, func() tea.Msg { return ClosedMsg{} }
		case key.Matches(msg, m.keymaps.Quit):
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateSize()
	}

	// Pass message to list
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// Show configures the command palette with current dimensions and context.
// xplr.Model is responsible for setting its state to xplrCmdPaletteView when calling this.
func (m *Model) Show(width, height int, searchActive bool) {
	m.width = width
	m.height = height
	m.updateSize()
	m.updateCommandAvailability(searchActive)
}

// View renders the command palette
func (m Model) View() string {
	content := m.list.View()
	overlay := m.styles.Overlay.Render(content)
	return utils.CenterOverlay(overlay, m.width, m.height)
}

// updateSize updates the list size based on window dimensions
func (m *Model) updateSize() {
	listWidth := m.width - config.OverlayInsetMargin
	listHeight := m.height - config.OverlayInsetMargin
	if listWidth < 10 {
		listWidth = 10
	}
	if listHeight < 5 {
		listHeight = 5
	}
	m.list.SetSize(listWidth, listHeight)
}

// updateCommandAvailability updates the enabled state of commands based on context
func (m *Model) updateCommandAvailability(searchActive bool) {
	items := m.list.Items()
	for i, item := range items {
		if cmd, ok := item.(commandItem); ok {
			// Determine if command should be enabled
			enabled := true

			if cmd.context == "Search" {
				enabled = searchActive
			}

			if cmd.context == "Overlay" {
				enabled = false
			}

			cmd.enabled = enabled
			items[i] = cmd
		}
	}
	m.list.SetItems(items)
}

// parseKeyString converts a key string to a tea.KeyMsg
func parseKeyString(keyStr string) tea.KeyMsg {
	// Handle common key combinations
	lower := strings.ToLower(keyStr)

	// Handle special keys
	switch lower {
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "esc", "escape":
		return tea.KeyMsg{Type: tea.KeyEsc}
	case "tab":
		return tea.KeyMsg{Type: tea.KeyTab}
	case "space", " ":
		return tea.KeyMsg{Type: tea.KeySpace, Runes: []rune{' '}}
	case "shift+tab":
		return tea.KeyMsg{Type: tea.KeyShiftTab}
	case "shift+left":
		return tea.KeyMsg{Type: tea.KeyShiftLeft}
	case "shift+right":
		return tea.KeyMsg{Type: tea.KeyShiftRight}
	}

	// Handle ctrl+ combinations
	if strings.HasPrefix(lower, "ctrl+") {
		key := strings.TrimPrefix(lower, "ctrl+")
		return tea.KeyMsg{Type: tea.KeyCtrlC + tea.KeyType(key[0]-'a')} // Approximation
	}

	// Handle single character keys and simple combinations
	if len(keyStr) == 1 {
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{rune(keyStr[0])}}
	}

	// Default: treat as runes
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(keyStr)}
}
