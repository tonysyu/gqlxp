package tui

import (
	"reflect"
	"slices"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tonysyu/igq/gql"
)

const (
	intialPanels   = 2
	maxPanes       = 6
	minPanes       = 1
	helpHeight     = 5
	navbarHeight   = 3
	overlayPadding = 1
	overlayMargin  = 2
)

type GQLType string

const (
	QueryType     GQLType = "Query"
	MutationType  GQLType = "Mutation"
	ObjectType    GQLType = "Object"
	InputType     GQLType = "Input"
	EnumType      GQLType = "Enum"
	ScalarType    GQLType = "Scalar"
	InterfaceType GQLType = "Interface"
	UnionType     GQLType = "Union"
	DirectiveType GQLType = "Directive"
)

// availableGQLTypes defines the ordered list of GQL types for navigation
var availableGQLTypes = []GQLType{QueryType, MutationType, ObjectType, InputType, EnumType, ScalarType, InterfaceType, UnionType, DirectiveType}

var (
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	cursorLineStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("57")).
			Foreground(lipgloss.Color("230"))

	placeholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("238"))

	endOfBufferStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("235"))

	focusedPlaceholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("99"))

	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("238"))

	blurredBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.HiddenBorder())

	navbarStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Margin(0, 0, 1, 0)

	activeTabStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("57")).
			Padding(0, 2).
			Bold(true)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("244")).
				Padding(0, 2)

	overlayStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238")).
			Padding(overlayPadding).
			Margin(overlayMargin)
)

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
	width          int
	height         int
	keymap         keymap
	globalKeyBinds []key.Binding
	help           help.Model
	panels         []Panel
	focus          int
	schema         gql.GraphQLSchema
	fieldType      GQLType
	overlay        overlayModel
}

func NewModel(schema gql.GraphQLSchema) mainModel {
	m := mainModel{
		panels:    make([]Panel, intialPanels),
		help:      help.New(),
		schema:    schema,
		fieldType: QueryType,
		overlay:   newOverlayModel(),
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
	for i := 0; i < v.NumField(); i++ {
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
			content := "No item selected"
			if focusedPanel, ok := m.panels[m.focus].(*listPanel); ok {
				if selectedItem := focusedPanel.Model.SelectedItem(); selectedItem != nil {
					if listItem, ok := selectedItem.(ListItem); ok {
						content = "# " + listItem.Title() + "\n\n" + listItem.Description()
					}
				}
			}
			m.overlay.Show(content, m.width, m.height)
		case key.Matches(msg, m.keymap.NextPanel):
			m.focus++
			if m.focus > len(m.panels)-1 {
				m.focus = 0
			}
		case key.Matches(msg, m.keymap.PrevPanel):
			m.focus--
			if m.focus < 0 {
				m.focus = len(m.panels) - 1
			}
		case key.Matches(msg, m.keymap.ToggleGQLType):
			m.incrementGQLTypeIndex(1)
		case key.Matches(msg, m.keymap.ReverseToggleGQLType):
			m.incrementGQLTypeIndex(-1)
		}
	case openPanelMsg:
		m.handleOpenPanel(msg.panel)
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}

	m.sizePanels()

	// Update panels based on message type and focus
	for i := range m.panels {
		var newModel tea.Model
		var cmd tea.Cmd

		shouldReceiveMsg := m.shouldPanelReceiveMessage(i, msg)
		if shouldReceiveMsg {
			newModel, cmd = m.panels[i].Update(msg)
			if panel, ok := newModel.(Panel); ok {
				m.panels[i] = panel
			}
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// shouldPanelReceiveMessage determines if a panel should receive a message
// based on the panel index, current focus, and message type
func (m *mainModel) shouldPanelReceiveMessage(panelIndex int, msg tea.Msg) bool {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global navigation keys handled by main model should not go to panels
		for _, binding := range m.globalKeyBinds {
			if key.Matches(msg, binding) {
				return false
			}
		}
		// All other key messages should only go to the focused panel
		return panelIndex == m.focus
	case openPanelMsg:
		// openPanelMsg is handled by main model, not individual panels
		return false
	default:
		// Unknown message types go to all panels (safe default)
		return true
	}
}

func (m *mainModel) sizePanels() {
	for i := range m.panels {
		m.panels[i].SetSize(m.width/len(m.panels), m.height-helpHeight-navbarHeight)
	}
}

// addPanel adds a new panel to the model
func (m *mainModel) addPanel(panel Panel) {
	if len(m.panels) < maxPanes {
		m.panels = append(m.panels, panel)
		m.sizePanels()
	}
}

// handleOpenPanel handles when an item is opened
func (m *mainModel) handleOpenPanel(newPanel Panel) {
	nextPanelIndex := m.focus + 1

	// If there's a next panel, replace it
	if nextPanelIndex < len(m.panels) {
		m.panels[nextPanelIndex] = newPanel
	} else if len(m.panels) < maxPanes {
		// Add a new panel if we haven't reached the max
		m.addPanel(newPanel)
	}

	m.sizePanels()
}

// resetAndLoadMainPanel defines initial panels and loads currently selected GQL type.
// This method is called on initilization and when switching types, so that detail panels get
// cleared out to avoid inconsistencies across panels.
func (m *mainModel) resetAndLoadMainPanel() {
	// Initialize panels with empty list models
	for i := range intialPanels {
		m.panels[i] = newStringPanel("")
	}

	// Load initial fields based on currently selected GQL type
	m.loadMainPanel()
}

// loadMainPanel loads the the currently selected GQL type in the main (left-most) panel
func (m *mainModel) loadMainPanel() {
	var items []ListItem
	var title string

	switch m.fieldType {
	case QueryType:
		items = adaptFieldDefinitions(gql.CollectAndSortMapValues(m.schema.Query))
		title = "Query Fields"
	case MutationType:
		items = adaptFieldDefinitions(gql.CollectAndSortMapValues(m.schema.Mutation))
		title = "Mutation Fields"
	case ObjectType:
		items = adaptObjectDefinitions(gql.CollectAndSortMapValues(m.schema.Object))
		title = "Object Types"
	case InputType:
		items = adaptInputDefinitions(gql.CollectAndSortMapValues(m.schema.Input))
		title = "Input Types"
	case EnumType:
		items = adaptEnumDefinitions(gql.CollectAndSortMapValues(m.schema.Enum))
		title = "Enum Types"
	case ScalarType:
		items = adaptScalarDefinitions(gql.CollectAndSortMapValues(m.schema.Scalar))
		title = "Scalar Types"
	case InterfaceType:
		items = adaptInterfaceDefinitions(gql.CollectAndSortMapValues(m.schema.Interface))
		title = "Interface Types"
	case UnionType:
		items = adaptUnionDefinitions(gql.CollectAndSortMapValues(m.schema.Union))
		title = "Union Types"
	case DirectiveType:
		items = adaptDirectiveDefinitions(gql.CollectAndSortMapValues(m.schema.Directive))
		title = "Directive Types"
	}

	m.panels[0] = newListPanel(items, title)
	// Move focus to the main panel when switching fields.
	m.focus = 0

	// Auto-open detail panel for the first item if available
	if len(items) > 0 {
		if firstItem, ok := items[0].(ListItem); ok {
			if newPanel, ok := firstItem.Open(); ok {
				m.handleOpenPanel(newPanel)
			}
		}
	}
}

// incrementGQLTypeIndex cycles through available GQL types with wraparound
func (m *mainModel) incrementGQLTypeIndex(offset int) {
	// Find current GQL type index
	currentIndex := slices.IndexFunc(availableGQLTypes, func(fieldType GQLType) bool {
		return m.fieldType == fieldType
	})

	newIndex := (currentIndex + offset)
	// Force new index to wraparound, if is out-of-bounds on either the beginning or end:
	if newIndex < 0 {
		newIndex = len(availableGQLTypes) - 1
	} else if newIndex >= len(availableGQLTypes) {
		newIndex = 0
	}
	m.fieldType = availableGQLTypes[newIndex]

	m.resetAndLoadMainPanel()
	m.sizePanels()
}

// renderGQLTypeNavbar creates the navbar showing GQL types
func (m *mainModel) renderGQLTypeNavbar() string {
	var tabs []string

	for _, fieldType := range availableGQLTypes {
		var style lipgloss.Style
		if m.fieldType == fieldType {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		tabs = append(tabs, style.Render(string(fieldType)))
	}

	navbar := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	return navbarStyle.Render(navbar)
}

func (m mainModel) View() string {
	help := m.help.ShortHelpView([]key.Binding{
		m.keymap.NextPanel,
		m.keymap.PrevPanel,
		m.keymap.ToggleGQLType,
		m.keymap.ToggleOverlay,
		m.keymap.Quit,
	})

	var views []string
	for i := range m.panels {
		panelView := m.panels[i].View()
		if i == m.focus && !m.overlay.IsActive() {
			panelView = focusedBorderStyle.Render(panelView)
		} else {
			panelView = blurredBorderStyle.Render(panelView)
		}
		views = append(views, panelView)
	}

	navbar := m.renderGQLTypeNavbar()
	panels := lipgloss.JoinHorizontal(lipgloss.Top, views...)
	mainView := navbar + "\n" + panels + "\n\n" + help

	// Show overlay if active
	if m.overlay.IsActive() {
		overlayContent := m.overlay.View()
		overlay := overlayStyle.Render(overlayContent)

		// Center the overlay on screen
		overlayHeight := lipgloss.Height(overlay)
		overlayWidth := lipgloss.Width(overlay)

		verticalMargin := (m.height - overlayHeight) / 2
		horizontalMargin := (m.width - overlayWidth) / 2

		if verticalMargin < 0 {
			verticalMargin = 0
		}
		if horizontalMargin < 0 {
			horizontalMargin = 0
		}

		// Position the overlay over the main view
		positionedOverlay := lipgloss.NewStyle().
			MarginTop(verticalMargin).
			MarginLeft(horizontalMargin).
			Render(overlay)

		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, positionedOverlay)
	}

	return mainView
}
