package xplr

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tonysyu/gqlxp/library"
	"github.com/tonysyu/gqlxp/tui/adapters"
	"github.com/tonysyu/gqlxp/tui/config"
	"github.com/tonysyu/gqlxp/tui/overlay"
	"github.com/tonysyu/gqlxp/tui/xplr/components"
	"github.com/tonysyu/gqlxp/tui/xplr/navigation"
	"slices"
)

type FavoriteToggledMsg struct {
	Favorites []string
}

// SchemaLoadedMsg is sent when a schema is loaded or updated
type SchemaLoadedMsg struct {
	Schema         adapters.SchemaView
	SchemaID       string
	Favorites      []string
	HasLibraryData bool
}

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
			ToggleFavorite: key.NewBinding(
				key.WithKeys("f"),
				key.WithHelp("f", "favorite"),
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
		m.keymap.ToggleFavorite,
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
	m.schemaID = schemaID
	m.favorites = metadata.Favorites
	m.hasLibraryData = true
	m.resetAndLoadMainPanel()
	return m
}

// SetSchemaID sets the schema ID for library integration
func (m *Model) SetSchemaID(id string) {
	m.schemaID = id
}

// GetSchemaID returns the schema ID
func (m Model) GetSchemaID() string {
	return m.schemaID
}

// SetFavorites sets the favorites list
func (m *Model) SetFavorites(favorites []string) {
	m.favorites = favorites
}

// GetFavorites returns the favorites list
func (m Model) GetFavorites() []string {
	return m.favorites
}

// SetHasLibraryData sets whether this schema has library metadata
func (m *Model) SetHasLibraryData(has bool) {
	m.hasLibraryData = has
}

// HasLibraryData returns whether this schema has library metadata
func (m Model) HasLibraryData() bool {
	return m.hasLibraryData
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
		m.schemaID = msg.SchemaID
		m.favorites = msg.Favorites
		m.hasLibraryData = msg.HasLibraryData
		m.resetAndLoadMainPanel()
		return m, nil
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
	case FavoriteToggledMsg:
		m.favorites = msg.Favorites
		// Refresh panels in place instead of resetting to preserve navigation state
		m.refreshPanelsWithFavorites()
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
		// For top-level panels, we check RefName (field names); otherwise TypeName
		items = wrapItemsWithFavorites(items, m.favorites, true)
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
// Only top-level panels can favorite items
func (m *Model) toggleFavoriteForSelectedItem() tea.Cmd {
	if m.nav.CurrentPanel() == nil {
		return nil
	}

	// Only allow favoriting at the top level
	if !m.nav.IsAtTopLevelPanel() {
		return nil
	}

	panel := m.nav.CurrentPanel()
	if selectedItem := panel.SelectedItem(); selectedItem != nil {
		if listItem, ok := selectedItem.(components.ListItem); ok {
			// For top-level panels, use RefName() to store field names
			favoriteName := listItem.RefName()
			return m.toggleFavorite(favoriteName)
		}
	}
	return nil
}

// refreshPanelsWithFavorites updates all panels in the stack to reflect current favorites
// without resetting navigation state. This preserves panel stack, selections, and scroll positions.
// Only the top-level panel (position 0) can have favorites.
func (m *Model) refreshPanelsWithFavorites() {
	for panelIndex, panel := range m.nav.Stack().All() {
		if panel == nil {
			continue
		}

		items := panel.Items()
		if len(items) == 0 {
			continue
		}

		// Unwrap items to get original items
		unwrappedItems := make([]components.ListItem, len(items))
		for i, item := range items {
			if listItem, ok := item.(components.ListItem); ok {
				unwrappedItems[i] = unwrapFavoritableItem(listItem)
			}
		}

		// Only wrap top-level panel (position 0) with favorites
		var refreshedItems []components.ListItem
		isTopLevel := panelIndex == 0
		if isTopLevel {
			refreshedItems = wrapItemsWithFavorites(unwrappedItems, m.favorites, true)
		} else {
			refreshedItems = unwrappedItems
		}

		// Convert to []list.Item for SetItems
		listItems := make([]list.Item, len(refreshedItems))
		for i, item := range refreshedItems {
			listItems[i] = item
		}

		panel.SetItems(listItems)
	}
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
			return FavoriteToggledMsg{
				Favorites: m.favorites,
			}
		}

		return FavoriteToggledMsg{
			Favorites: schema.Metadata.Favorites,
		}
	}
}
