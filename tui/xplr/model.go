package xplr

import (
	"reflect"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tonysyu/gqlxp/library"
	"github.com/tonysyu/gqlxp/tui/adapters"
	"github.com/tonysyu/gqlxp/tui/config"
	"github.com/tonysyu/gqlxp/tui/overlay"
	"github.com/tonysyu/gqlxp/tui/xplr/components"
	"github.com/tonysyu/gqlxp/tui/xplr/navigation"
	"slices"
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

type SetGQLTypeMsg struct {
	GQLType gqlType
}

type FavoriteToggledMsg struct {
	Favorites []string
}

var (
	quitKeyBinding = key.NewBinding(
		key.WithKeys("ctrl+c", "ctrl+d"),
		key.WithHelp("ctrl+c", "quit"),
	)
)

type keymap = struct {
	NextPanel, PrevPanel, Quit, ToggleGQLType, ReverseToggleGQLType, ToggleOverlay, ToggleFavorite key.Binding
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
	schemaID       string   // Schema ID if loaded from library
	favorites      []string // List of favorited type names
	hasLibraryData bool     // Whether this schema has library metadata

	width          int
	height         int
	Styles         config.Styles
	keymap         keymap
	globalKeyBinds []key.Binding
	help           help.Model
}

// New creates a new schema explorer model
func New(schema adapters.SchemaView) Model {
	styles := config.DefaultStyles()
	m := Model{
		help:    help.New(),
		schema:  schema,
		Styles:  styles,
		Overlay: overlay.New(styles),
		nav:     navigation.NewNavigationManager(config.VisiblePanelCount),
		keymap: keymap{
			NextPanel: key.NewBinding(
				key.WithKeys("tab", "]"),
				key.WithHelp("tab", "next"),
			),
			PrevPanel: key.NewBinding(
				key.WithKeys("shift+tab", "["),
				key.WithHelp("shift+tab", "prev"),
			),
			Quit: quitKeyBinding,
			ToggleGQLType: key.NewBinding(
				key.WithKeys("ctrl+t", "}"),
				key.WithHelp("ctrl+t", "toggle type"),
			),
			ReverseToggleGQLType: key.NewBinding(
				key.WithKeys("ctrl+r", "{"),
				key.WithHelp("ctrl+r", "reverse toggle type"),
			),
			ToggleOverlay: key.NewBinding(
				key.WithKeys(" "),
				key.WithHelp("space", "overlay"),
			),
			ToggleFavorite: key.NewBinding(
				key.WithKeys("f"),
				key.WithHelp("f", "favorite"),
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

// SetSchemaID sets the schema ID for library integration
func (m *Model) SetSchemaID(id string) {
	m.schemaID = id
}

// SetFavorites sets the favorites list
func (m *Model) SetFavorites(favorites []string) {
	m.favorites = favorites
}

// SetHasLibraryData sets whether this schema has library metadata
func (m *Model) SetHasLibraryData(has bool) {
	m.hasLibraryData = has
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.ToggleOverlay):
			m.openOverlayForSelectedItem()
		case key.Matches(msg, m.keymap.ToggleFavorite):
			if m.hasLibraryData {
				return m, m.toggleFavoriteForSelectedItem()
			}
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
	case SetGQLTypeMsg:
		m.nav.SwitchType(gqlTypeToNavType(msg.GQLType))
		m.resetAndLoadMainPanel()
	case FavoriteToggledMsg:
		m.favorites = msg.Favorites
		m.resetAndLoadMainPanel()
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

	// Wrap items with favorites indicator if library data is available
	if m.hasLibraryData {
		items = wrapItemsWithFavorites(items, m.favorites)
	}

	m.nav.SetCurrentPanel(components.NewPanel(items, title))
	m.updatePanelFocusStates()

	// Auto-open detail panel for the first item if available
	if len(items) > 0 {
		if firstItem, ok := items[0].(components.ListItem); ok {
			if newPanel, ok := firstItem.OpenPanel(); ok {
				m.handleOpenPanel(newPanel)
			}
		}
	}
}

// toggleFavoriteForSelectedItem toggles favorite status for selected item
func (m *Model) toggleFavoriteForSelectedItem() tea.Cmd {
	if m.nav.CurrentPanel() == nil {
		return nil
	}
	panel := m.nav.CurrentPanel()
	if selectedItem := panel.SelectedItem(); selectedItem != nil {
		if listItem, ok := selectedItem.(components.ListItem); ok {
			typeName := listItem.TypeName()
			return m.toggleFavorite(typeName)
		}
	}
	return nil
}

// toggleFavorite toggles favorite status and saves to library
func (m *Model) toggleFavorite(typeName string) tea.Cmd {
	return func() tea.Msg {
		lib := library.NewLibrary()

		if slices.Contains(m.favorites, typeName) {
			_ = lib.RemoveFavorite(m.schemaID, typeName)
		} else {
			_ = lib.AddFavorite(m.schemaID, typeName)
		}

		// Reload favorites from library to get fresh state
		schema, err := lib.Get(m.schemaID)
		if err != nil {
			// On error, return current favorites unchanged
			return FavoriteToggledMsg{Favorites: m.favorites}
		}

		return FavoriteToggledMsg{Favorites: schema.Metadata.Favorites}
	}
}

// gqlTypeToNavType converts old gqlType to navigation.GQLType
func gqlTypeToNavType(t gqlType) navigation.GQLType {
	return navigation.GQLType(t)
}

// navTypeToGQLType converts navigation.GQLType to old gqlType
func navTypeToGQLType(t navigation.GQLType) gqlType {
	return gqlType(t)
}

// renderGQLTypeNavbar creates the navbar showing GQL types
func (m *Model) renderGQLTypeNavbar() string {
	var tabs []string

	for _, fieldType := range m.nav.AllTypes() {
		var style lipgloss.Style
		if m.nav.CurrentType() == fieldType {
			style = m.Styles.ActiveTab
		} else {
			style = m.Styles.InactiveTab
		}
		tabs = append(tabs, style.Render(string(fieldType)))
	}

	navbar := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	return m.Styles.Navbar.Render(navbar)
}

// renderBreadcrumbs renders the breadcrumb trail
func (m *Model) renderBreadcrumbs() string {
	crumbs := m.nav.Breadcrumbs()
	if len(crumbs) == 0 {
		return ""
	}

	// Build breadcrumb parts with separators
	var parts []string
	separator := " > "

	for i, crumb := range crumbs {
		if i > 0 {
			parts = append(parts, separator)
		}
		// Apply special color to the last breadcrumb
		if i == len(crumbs)-1 {
			parts = append(parts, m.Styles.CurrentBreadcrumb.Render(crumb))
		} else {
			parts = append(parts, crumb)
		}
	}

	breadcrumbText := lipgloss.JoinHorizontal(lipgloss.Left, parts...)
	return m.Styles.Breadcrumbs.Render(breadcrumbText)
}

func (m Model) View() string {
	// Build help key bindings
	helpBindings := []key.Binding{
		m.keymap.NextPanel,
		m.keymap.PrevPanel,
		m.keymap.ToggleGQLType,
		m.keymap.ToggleOverlay,
	}
	// Add library-specific keybindings if available
	if m.hasLibraryData {
		helpBindings = append(helpBindings, m.keymap.ToggleFavorite)
	}
	helpBindings = append(helpBindings, m.keymap.Quit)

	help := m.help.ShortHelpView(helpBindings)

	// Show overlay if active, and return immediately
	if m.Overlay.IsActive() {
		return m.Overlay.View()
	}

	var views []string
	if m.nav.CurrentPanel() != nil {
		views = append(views, m.nav.CurrentPanel().View())
	}
	if m.nav.NextPanel() != nil {
		views = append(views, m.nav.NextPanel().View())
	}

	navbar := m.renderGQLTypeNavbar()
	breadcrumbs := m.renderBreadcrumbs()
	panels := lipgloss.JoinHorizontal(lipgloss.Top, views...)

	mainView := lipgloss.JoinVertical(0, navbar, breadcrumbs, panels, help)
	return mainView
}
