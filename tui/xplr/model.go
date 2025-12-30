package xplr

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tonysyu/gqlxp/library"
	"github.com/tonysyu/gqlxp/tui/adapters"
	"github.com/tonysyu/gqlxp/tui/config"
	"github.com/tonysyu/gqlxp/tui/overlay"
	"github.com/tonysyu/gqlxp/tui/xplr/components"
	"github.com/tonysyu/gqlxp/tui/xplr/navigation"
)

// SchemaLoadedMsg is sent when a schema is loaded or updated
type SchemaLoadedMsg struct {
	Schema         adapters.SchemaView
	SchemaID       string
	HasLibraryData bool
}

// SelectionTarget specifies a type and optional field to pre-select in the TUI
type SelectionTarget struct {
	TypeName  string
	FieldName string
}

type keymap = struct {
	NextPanel, PrevPanel, Quit, ToggleGQLType, ReverseToggleGQLType, ToggleOverlay key.Binding
}

// Model is the main schema explorer model
type Model struct {
	// Parsed GraphQL schema that's displayed in the TUI.
	schema adapters.SchemaView
	// Navigation manager coordinates panel stack, breadcrumbs, and type selection
	nav *navigation.NavigationManager
	// Overlay for displaying ListItem.Details()
	Overlay overlay.Model

	// Library integration (optional)
	SchemaID       string // Schema ID if loaded from library
	HasLibraryData bool   // Whether this schema has library metadata

	width          int
	height         int
	Styles         config.Styles
	keymap         keymap
	globalKeyBinds []key.Binding
	help           help.Model
}

// NewEmpty creates a new schema explorer model without a schema
// The schema can be loaded later via SchemaLoadedMsg
func NewEmpty() Model {
	styles := config.DefaultStyles()
	m := Model{
		help:    help.New(),
		Styles:  styles,
		Overlay: overlay.New(styles),
		nav:     navigation.NewNavigationManager(config.VisiblePanelCount),
		keymap: keymap{
			NextPanel: key.NewBinding(
				key.WithKeys("]", "tab"),
				key.WithHelp("]/tab", "next"),
			),
			PrevPanel: key.NewBinding(
				key.WithKeys("shift+tab", "["),
				key.WithHelp("[/⇧+tab", "prev"),
			),
			Quit: key.NewBinding(
				key.WithKeys("ctrl+c", "ctrl+d"),
				key.WithHelp("⌃+c", "quit"),
			),
			ToggleGQLType: key.NewBinding(
				key.WithKeys("ctrl+t", "}"),
				key.WithHelp("}/⌃+T", "next type"),
			),
			ReverseToggleGQLType: key.NewBinding(
				key.WithKeys("ctrl+r", "{"),
				key.WithHelp("{/⌃+r", "prev type"),
			),
			ToggleOverlay: key.NewBinding(
				key.WithKeys(" "),
				key.WithHelp("space", "overlay"),
			),
		},
	}

	// Build globalKeyBinds from all keymap fields
	m.globalKeyBinds = []key.Binding{
		m.keymap.NextPanel,
		m.keymap.PrevPanel,
		m.keymap.Quit,
		m.keymap.ToggleGQLType,
		m.keymap.ReverseToggleGQLType,
		m.keymap.ToggleOverlay,
	}

	// Don't load panels until schema is provided
	return m
}

// New creates a new schema explorer model
func New(schema adapters.SchemaView) Model {
	m := NewEmpty()
	m.schema = schema
	m.resetAndLoadMainPanel()
	return m
}

// NewFromSchemaLibrary creates a new schema explorer model with library metadata
func NewFromSchemaLibrary(schema adapters.SchemaView, schemaID string, metadata library.SchemaMetadata) Model {
	m := NewEmpty()
	m.schema = schema
	m.SchemaID = schemaID
	m.HasLibraryData = true
	m.resetAndLoadMainPanel()
	return m
}

// Width returns the current width
func (m Model) Width() int {
	return m.width
}

// Height returns the current height
func (m Model) Height() int {
	return m.height
}

// CurrentType returns the currently selected GraphQL type
func (m Model) CurrentType() string {
	return string(m.nav.CurrentType())
}

// SwitchToType switches to the specified GraphQL type
// This is primarily used for testing
func (m *Model) SwitchToType(typeName string) {
	m.nav.SwitchType(navigation.GQLType(typeName))
	m.resetAndLoadMainPanel()
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	// Try overlay first - it intercepts messages when active
	var overlayCmd tea.Cmd
	var intercepted bool
	m.Overlay, overlayCmd, intercepted = m.Overlay.Update(msg)
	if intercepted {
		return m, overlayCmd
	}

	var cmds []tea.Cmd

	// Handle global messages
	switch msg := msg.(type) {
	case SchemaLoadedMsg:
		// Update schema and related properties
		m.schema = msg.Schema
		m.SchemaID = msg.SchemaID
		m.HasLibraryData = msg.HasLibraryData
		m.resetAndLoadMainPanel()
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.ToggleOverlay):
			m.openOverlayForSelectedItem()
		case key.Matches(msg, m.keymap.NextPanel):
			// Move forward in stack if there's at least one more panel ahead
			if m.nav.NavigateForward() {
				m.updatePanelFocusStates()
				// Open up child panel for ResultType if it exists
				focusedPanel := m.nav.CurrentPanel()
				if focusedPanel != nil {
					if openCmd := focusedPanel.OpenSelectedItem(); openCmd != nil {
						cmds = append(cmds, openCmd)
					}
				}
			}
		case key.Matches(msg, m.keymap.PrevPanel):
			// Move backward in stack if not at the beginning
			if m.nav.NavigateBackward() {
				m.updatePanelFocusStates()
			}
		case key.Matches(msg, m.keymap.ToggleGQLType):
			m.nav.CycleTypeForward()
			m.resetAndLoadMainPanel()
		case key.Matches(msg, m.keymap.ReverseToggleGQLType):
			m.nav.CycleTypeBackward()
			m.resetAndLoadMainPanel()
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
	if shouldReceiveMsg && m.nav.CurrentPanel() != nil {
		currentPanel := m.nav.CurrentPanel()
		newModel, cmd = currentPanel.Update(msg)
		if panel, ok := newModel.(*components.Panel); ok {
			m.nav.SetCurrentPanel(panel)
		}
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) openOverlayForSelectedItem() {
	// Always use the left panel (first visible panel in stack)
	if m.nav.CurrentPanel() == nil {
		return
	}
	focusedPanel := m.nav.CurrentPanel()
	if selectedItem := focusedPanel.SelectedItem(); selectedItem != nil {
		if listItem, ok := selectedItem.(components.ListItem); ok {
			// Some items don't have details, so these should now open the overlay
			if content := listItem.Details(); content != "" {
				m.Overlay.Show(content, m.width, m.height)
			}
		}
	}
}

// shouldFocusedPanelReceiveMessage determines if the focused panel should receive a message
func (m *Model) shouldFocusedPanelReceiveMessage(msg tea.Msg) bool {
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

func (m *Model) sizePanels() {
	panelWidth := m.width / config.VisiblePanelCount
	panelHeight := m.height - config.HelpHeight - config.NavbarHeight - config.BreadcrumbsHeight
	// Size only the visible panels (config.VisiblePanelCount = 2)
	if m.nav.CurrentPanel() != nil {
		m.nav.CurrentPanel().SetSize(
			panelWidth-m.Styles.FocusedPanel.GetHorizontalFrameSize(),
			panelHeight-m.Styles.FocusedPanel.GetVerticalFrameSize(),
		)
	}
	// The right panel might not exist, so check before resizing
	if m.nav.NextPanel() != nil {
		m.nav.NextPanel().SetSize(
			panelWidth-m.Styles.BlurredPanel.GetHorizontalFrameSize(),
			panelHeight-m.Styles.BlurredPanel.GetHorizontalFrameSize(),
		)
	}
}

// updatePanelFocusStates updates focus state for all visible panels based on stackPosition
func (m *Model) updatePanelFocusStates() {
	// Blur all panels first
	for _, panel := range m.nav.Stack().All() {
		if panel != nil {
			panel.SetBlurred()
		}
	}

	// Set focused state only for the currently focused panel
	if m.nav.CurrentPanel() != nil {
		m.nav.CurrentPanel().SetFocused()
	}
}

// handleOpenPanel handles when an item is opened
// The new panel is added to the stack after the currently focused panel
func (m *Model) handleOpenPanel(newPanel *components.Panel) {
	m.nav.OpenPanel(newPanel)
	m.sizePanels()
}

// resetAndLoadMainPanel defines initial panels and loads currently selected GQL type.
// This method is called on initilization and when switching types, so that detail panels get
// cleared out to avoid inconsistencies across panels.
func (m *Model) resetAndLoadMainPanel() {
	m.nav.Reset()
	m.loadMainPanel()
}

// loadMainPanel loads the the currently selected GQL type in the main (left-most) panel
func (m *Model) loadMainPanel() {
	var items []components.ListItem
	var title string

	switch m.nav.CurrentType() {
	case navigation.QueryType:
		items = m.schema.GetQueryItems()
		title = "Query Fields"
	case navigation.MutationType:
		items = m.schema.GetMutationItems()
		title = "Mutation Fields"
	case navigation.ObjectType:
		items = m.schema.GetObjectItems()
		title = "Object Types"
	case navigation.InputType:
		items = m.schema.GetInputItems()
		title = "Input Types"
	case navigation.EnumType:
		items = m.schema.GetEnumItems()
		title = "Enum Types"
	case navigation.ScalarType:
		items = m.schema.GetScalarItems()
		title = "Scalar Types"
	case navigation.InterfaceType:
		items = m.schema.GetInterfaceItems()
		title = "Interface Types"
	case navigation.UnionType:
		items = m.schema.GetUnionItems()
		title = "Union Types"
	case navigation.DirectiveType:
		items = m.schema.GetDirectiveItems()
		title = "Directive Types"
	}

	m.nav.SetCurrentPanel(components.NewPanel(items, title))
	m.updatePanelFocusStates()

	// Auto-open detail panel for the first item if available
	if len(items) > 0 {
		if newPanel, ok := items[0].OpenPanel(); ok {
			m.handleOpenPanel(newPanel)
		}
	}
}

// ApplySelection applies a selection target to the model
// This navigates to the specified type category and selects the item
func (m *Model) ApplySelection(target SelectionTarget) {
	if target.TypeName == "" {
		return
	}

	// Find which GQL type category contains the target type
	gqlType, found := m.schema.FindTypeCategory(target.TypeName)
	if !found {
		return
	}

	// Switch to that type category
	m.nav.SwitchType(gqlType)
	m.resetAndLoadMainPanel()

	currentPanel := m.nav.CurrentPanel()
	if currentPanel == nil {
		return
	}

	// For Query and Mutation types, the fields are shown directly in the first panel
	// So if target.TypeName is "Query" or "Mutation", we skip selecting it and go straight to the field
	if gqlType == navigation.QueryType || gqlType == navigation.MutationType {
		if target.FieldName != "" {
			// Select the field directly in the current panel
			if !currentPanel.SelectItemByName(target.FieldName) {
				return
			}
			// Navigate forward to show the field's details
			if openCmd := currentPanel.OpenSelectedItem(); openCmd != nil {
				if msg, ok := openCmd().(components.OpenPanelMsg); ok {
					m.handleOpenPanel(msg.Panel)
					if m.nav.NavigateForward() {
						m.updatePanelFocusStates()
					}
				}
			}
		}
		return
	}

	// For other types (Object, Input, Enum, etc.), select the type in the current panel
	if !currentPanel.SelectItemByName(target.TypeName) {
		return
	}

	// If a field name is specified, navigate forward and select the field
	if target.FieldName != "" {
		// Open child panel for the selected item
		if openCmd := currentPanel.OpenSelectedItem(); openCmd != nil {
			// Execute the command to populate the next panel
			if msg, ok := openCmd().(components.OpenPanelMsg); ok {
				// Add the new panel to the stack first
				m.handleOpenPanel(msg.Panel)
				// Then navigate forward (this adds breadcrumb from current panel's selected item)
				if m.nav.NavigateForward() {
					m.updatePanelFocusStates()
					// Select the field in the newly opened panel (now at current position)
					if currentPanel := m.nav.CurrentPanel(); currentPanel != nil {
						currentPanel.SelectItemByName(target.FieldName)
					}
				}
			}
		}
	}
}
