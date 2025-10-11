package tui

import (
	"reflect"
	"slices"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tonysyu/gqlxp/tui/adapters"
	"github.com/tonysyu/gqlxp/tui/components"
	"github.com/tonysyu/gqlxp/tui/config"
)

type gqlType string

const (
	queryType     gqlType = "Query"
	mutationType  gqlType = "Mutation"
	objectType    gqlType = "Object"
	inputType     gqlType = "Input"
	enumType      gqlType = "Enum"
	scalarType    gqlType = "Scalar"
	interfaceType gqlType = "Interface"
	unionType     gqlType = "Union"
	directiveType gqlType = "Directive"
)

// availableGQLTypes defines the ordered list of GQL types for navigation
var availableGQLTypes = []gqlType{queryType, mutationType, objectType, inputType, enumType, scalarType, interfaceType, unionType, directiveType}

var (
	quitKeyBinding = key.NewBinding(
		key.WithKeys("ctrl+c", "ctrl+d"),
		key.WithHelp("ctrl+c", "quit"),
	)
)

type keymap = struct {
	NextPanel, PrevPanel, Quit, ToggleGQLType, ReverseToggleGQLType, ToggleOverlay key.Binding
}

type mainModel struct {
	// Parsed GraphQL schema that's displayed in the TUI.
	schema adapters.SchemaView
	// Panels displaying list-views of GraphQL types.
	// A list of top-level types (see availableGQLTypes) is at the bottom of the stack, and children
	// of those types (e.g. fields, inputs, return types) are displayed in additional panels.
	panelStack []components.Panel
	// Position of the currently focused panel in the panelStack.
	// This may not be the top-most item in the stack.
	stackPosition int
	// Currently displayed GraphQL Type (see availableGQLTypes)
	selectedGQLType gqlType
	// Overlay for displaying ListItem.Details()
	overlay overlayModel

	width          int
	height         int
	styles         config.Styles
	keymap         keymap
	globalKeyBinds []key.Binding
	help           help.Model
}

func newModel(schema adapters.SchemaView) mainModel {
	styles := config.DefaultStyles()
	m := mainModel{
		panelStack:      make([]components.Panel, config.VisiblePanelCount),
		stackPosition:   0,
		help:            help.New(),
		schema:          schema,
		selectedGQLType: queryType,
		styles:          styles,
		overlay:         newOverlayModel(styles),
		keymap: keymap{
			NextPanel: key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "next"),
			),
			PrevPanel: key.NewBinding(
				key.WithKeys("shift+tab"),
				key.WithHelp("shift+tab", "prev"),
			),
			Quit: quitKeyBinding,
			ToggleGQLType: key.NewBinding(
				key.WithKeys("ctrl+t"),
				key.WithHelp("ctrl+t", "toggle type "),
			),
			ReverseToggleGQLType: key.NewBinding(
				key.WithKeys("ctrl+r"),
				key.WithHelp("ctrl+r", "reverse toggle type"),
			),
			ToggleOverlay: key.NewBinding(
				key.WithKeys(" "),
				key.WithHelp("space", "overlay"),
			),
		},
	}

	// Build globalKeyBinds from all keymap fields using reflection
	v := reflect.ValueOf(m.keymap)
	m.globalKeyBinds = make([]key.Binding, v.NumField())
	for i := range v.NumField() {
		m.globalKeyBinds[i] = v.Field(i).Interface().(key.Binding)
	}

	m.resetAndLoadMainPanel()
	return m
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Try overlay first - it intercepts messages when active
	var overlayCmd tea.Cmd
	var intercepted bool
	m.overlay, overlayCmd, intercepted = m.overlay.Update(msg)
	if intercepted {
		return m, overlayCmd
	}

	// Handle global messages
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.ToggleOverlay):
			m.openOverlayForSelectedItem()
		case key.Matches(msg, m.keymap.NextPanel):
			// Move forward in stack if there's at least one more panel ahead
			if m.stackPosition+1 < len(m.panelStack) {
				m.stackPosition++
			}
		case key.Matches(msg, m.keymap.PrevPanel):
			// Move backward in stack if not at the beginning
			if m.stackPosition > 0 {
				m.stackPosition--
			}
		case key.Matches(msg, m.keymap.ToggleGQLType):
			m.incrementGQLTypeIndex(1)
		case key.Matches(msg, m.keymap.ReverseToggleGQLType):
			m.incrementGQLTypeIndex(-1)
		}
	case components.OpenPanelMsg:
		m.handleOpenPanel(msg.Panel)
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}

	m.sizePanels()

	// Update visible panels in the stack
	var newModel tea.Model
	var cmd tea.Cmd

	// Only the left (focused) panel receives input; right panel is display-only
	shouldReceiveMsg := m.shouldFocusedPanelReceiveMessage(msg)
	if shouldReceiveMsg {
		newModel, cmd = m.panelStack[m.stackPosition].Update(msg)
		if panel, ok := newModel.(components.Panel); ok {
			m.panelStack[m.stackPosition] = panel
		}
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *mainModel) openOverlayForSelectedItem() {
	// Always use the left panel (first visible panel in stack)
	if focusedPanel, ok := m.panelStack[m.stackPosition].(*components.ListPanel); ok {
		if selectedItem := focusedPanel.SelectedItem(); selectedItem != nil {
			if listItem, ok := selectedItem.(components.ListItem); ok {
				// Some items don't have details, so these should now open the overlay
				if content := listItem.Details(); content != "" {
					m.overlay.Show(content, m.width, m.height)
				}
			}
		}
	}
}

// shouldFocusedPanelReceiveMessage determines if the focused panel should receive a message
func (m *mainModel) shouldFocusedPanelReceiveMessage(msg tea.Msg) bool {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global navigation keys handled by main model should not go to panels
		for _, binding := range m.globalKeyBinds {
			if key.Matches(msg, binding) {
				return false
			}
		}
		return true
	case components.OpenPanelMsg:
		// OpenPanelMsg is handled by main model, not individual panels
		return false
	default:
		// Unknown message types go to focued panel (safe default)
		return true
	}
}

func (m *mainModel) sizePanels() {
	panelWidth := m.width / config.VisiblePanelCount
	panelHeight := m.height - config.HelpHeight - config.NavbarHeight
	// Size only the visible panels (config.DisplayedPanels = 2)
	m.panelStack[m.stackPosition].SetSize(panelWidth, panelHeight)
	// The right panel might not exist, so check before resizing
	if len(m.panelStack) > m.stackPosition+1 {
		m.panelStack[m.stackPosition+1].SetSize(panelWidth, panelHeight)
	}
}

// handleOpenPanel handles when an item is opened
// The new panel is added to the stack after the currently focused panel
func (m *mainModel) handleOpenPanel(newPanel components.Panel) {
	// Truncate stack to keep only up to and including the current left panel
	m.panelStack = m.panelStack[:m.stackPosition+1]
	// Append the new panel - it will show on the right
	m.panelStack = append(m.panelStack, newPanel)

	m.sizePanels()
}

// resetAndLoadMainPanel defines initial panels and loads currently selected GQL type.
// This method is called on initilization and when switching types, so that detail panels get
// cleared out to avoid inconsistencies across panels.
func (m *mainModel) resetAndLoadMainPanel() {
	// Reset stack to initial state with empty panels
	m.panelStack = make([]components.Panel, config.VisiblePanelCount)
	for i := range config.VisiblePanelCount {
		m.panelStack[i] = components.NewStringPanel("")
	}
	m.stackPosition = 0

	// Load initial fields based on currently selected GQL type
	m.loadMainPanel()
}

// loadMainPanel loads the the currently selected GQL type in the main (left-most) panel
func (m *mainModel) loadMainPanel() {
	var items []components.ListItem
	var title string

	switch m.selectedGQLType {
	case queryType:
		items = m.schema.GetQueryItems()
		title = "Query Fields"
	case mutationType:
		items = m.schema.GetMutationItems()
		title = "Mutation Fields"
	case objectType:
		items = m.schema.GetObjectItems()
		title = "Object Types"
	case inputType:
		items = m.schema.GetInputItems()
		title = "Input Types"
	case enumType:
		items = m.schema.GetEnumItems()
		title = "Enum Types"
	case scalarType:
		items = m.schema.GetScalarItems()
		title = "Scalar Types"
	case interfaceType:
		items = m.schema.GetInterfaceItems()
		title = "Interface Types"
	case unionType:
		items = m.schema.GetUnionItems()
		title = "Union Types"
	case directiveType:
		items = m.schema.GetDirectiveItems()
		title = "Directive Types"
	}

	m.panelStack[0] = components.NewListPanel(items, title)
	// Reset to the beginning of the stack
	m.stackPosition = 0

	// Auto-open detail panel for the first item if available
	if len(items) > 0 {
		if firstItem, ok := items[0].(components.ListItem); ok {
			if newPanel, ok := firstItem.Open(); ok {
				m.handleOpenPanel(newPanel)
			}
		}
	}
}

// incrementGQLTypeIndex cycles through available GQL types with wraparound
func (m *mainModel) incrementGQLTypeIndex(offset int) {
	// Find current GQL type index
	currentIndex := slices.IndexFunc(availableGQLTypes, func(fieldType gqlType) bool {
		return m.selectedGQLType == fieldType
	})

	newIndex := (currentIndex + offset)
	// Force new index to wraparound, if is out-of-bounds on either the beginning or end:
	if newIndex < 0 {
		newIndex = len(availableGQLTypes) - 1
	} else if newIndex >= len(availableGQLTypes) {
		newIndex = 0
	}
	m.selectedGQLType = availableGQLTypes[newIndex]

	m.resetAndLoadMainPanel()
	m.sizePanels()
}

// renderGQLTypeNavbar creates the navbar showing GQL types
func (m *mainModel) renderGQLTypeNavbar() string {
	var tabs []string

	for _, fieldType := range availableGQLTypes {
		var style lipgloss.Style
		if m.selectedGQLType == fieldType {
			style = m.styles.ActiveTab
		} else {
			style = m.styles.InactiveTab
		}
		tabs = append(tabs, style.Render(string(fieldType)))
	}

	navbar := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	return m.styles.Navbar.Render(navbar)
}

func (m mainModel) View() string {
	help := m.help.ShortHelpView([]key.Binding{
		m.keymap.NextPanel,
		m.keymap.PrevPanel,
		m.keymap.ToggleGQLType,
		m.keymap.ToggleOverlay,
		m.keymap.Quit,
	})

	// Show overlay if active, and return immediately
	if m.overlay.IsActive() {
		return m.overlay.View()
	}

	views := []string{m.styles.FocusedBorder.Render(m.panelStack[m.stackPosition].View())}
	if len(m.panelStack) > m.stackPosition+1 {
		views = append(views, m.styles.BlurredBorder.Render(m.panelStack[m.stackPosition+1].View()))
	}

	navbar := m.renderGQLTypeNavbar()
	panels := lipgloss.JoinHorizontal(lipgloss.Top, views...)
	mainView := navbar + "\n" + panels + "\n\n" + help

	return mainView
}
