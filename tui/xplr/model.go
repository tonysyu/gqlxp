package xplr

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tonysyu/gqlxp/library"
	"github.com/tonysyu/gqlxp/search"
	"github.com/tonysyu/gqlxp/tui/adapters"
	"github.com/tonysyu/gqlxp/tui/config"
	"github.com/tonysyu/gqlxp/tui/xplr/cmdpalette"
	"github.com/tonysyu/gqlxp/tui/xplr/components"
	"github.com/tonysyu/gqlxp/tui/xplr/navigation"
	"github.com/tonysyu/gqlxp/tui/xplr/overlay"
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

type xplrState uint

const (
	xplrNormalView     xplrState = iota
	xplrOverlayView              // overlay is displayed
	xplrCmdPaletteView           // command palette is displayed
)

// Model is the main schema explorer model
type Model struct {
	state xplrState
	// Parsed GraphQL schema that's displayed in the TUI.
	schema adapters.SchemaView
	// Navigation manager coordinates panel stack, breadcrumbs, and type selection
	nav *navigation.NavigationManager
	// Overlay for displaying ListItem.Details()
	overlay overlay.Model
	// Command palette for discovering and executing commands
	commandPalette cmdpalette.Model

	// Library integration (optional)
	SchemaID       string // Schema ID if loaded from library
	HasLibraryData bool   // Whether this schema has library metadata

	// Search state
	searchInput   components.SearchInput
	searchFocused bool
	searchResults []components.ListItem
	searchBaseDir string // Base directory for search indexes

	width          int
	height         int
	Styles         config.Styles
	keymap         config.MainKeymaps
	globalKeyBinds []key.Binding
	help           help.Model
}

// NewEmpty creates a new schema explorer model without a schema
// The schema can be loaded later via SchemaLoadedMsg
func NewEmpty() Model {
	styles := config.DefaultStyles()
	mainKeymap := config.NewMainKeymaps()
	panelKeymap := config.NewPanelKeymaps()
	overlayKeymap := config.NewOverlayKeymaps()
	paletteKeymap := config.NewCommandPaletteKeymaps()

	m := Model{
		help:        help.New(),
		Styles:      styles,
		overlay:     overlay.New(styles),
		nav:         navigation.NewNavigationManager(config.VisiblePanelCount),
		searchInput: components.NewSearchInput(),
		keymap:      mainKeymap,
	}

	// Create command palette with all keymaps
	m.commandPalette = cmdpalette.New(styles, paletteKeymap, cmdpalette.CommandKeymaps{
		Main:    mainKeymap,
		Panel:   panelKeymap,
		Overlay: overlayKeymap,
	})

	// Build globalKeyBinds from all keymap fields
	m.globalKeyBinds = []key.Binding{
		m.keymap.NextPanel,
		m.keymap.PrevPanel,
		m.keymap.Quit,
		m.keymap.NextGQLType,
		m.keymap.PrevGQLType,
		m.keymap.ToggleOverlay,
		m.keymap.CommandPalette,
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

// IsOverlayVisible returns whether the overlay is currently displayed
func (m Model) IsOverlayVisible() bool {
	return m.state == xplrOverlayView
}

// SetOverlayStyle sets the overlay style (primarily for testing)
func (m *Model) SetOverlayStyle(style lipgloss.Style) {
	m.overlay.Styles.Overlay = style
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	// Handle close messages from sub-views before routing
	switch msg.(type) {
	case overlay.ClosedMsg, cmdpalette.ClosedMsg:
		m.state = xplrNormalView
		return m, nil
	}

	// Route to the active sub-view
	switch m.state {
	case xplrCmdPaletteView:
		var cmd tea.Cmd
		m.commandPalette, cmd = m.commandPalette.Update(msg)
		return m, cmd
	case xplrOverlayView:
		var cmd tea.Cmd
		m.overlay, cmd = m.overlay.Update(msg)
		return m, cmd
	}

	var cmds []tea.Cmd

	// Handle global messages (only reached in xplrNormalView)
	switch msg := msg.(type) {
	case SchemaLoadedMsg:
		// Update schema and related properties
		m.schema = msg.Schema
		m.SchemaID = msg.SchemaID
		m.HasLibraryData = msg.HasLibraryData
		m.resetAndLoadMainPanel()
		return m, nil
	case searchResultsMsg:
		// Update search results
		if msg.err != nil {
			// Handle search error - show empty results
			m.searchResults = nil
		} else {
			m.searchResults = m.convertSearchResultsToListItems(msg.results)
		}
		// Reload panel if we're on the search tab
		if m.nav.CurrentType() == navigation.SearchType {
			m.loadMainPanel()
		}
		return m, nil
	case tea.KeyMsg:
		// Handle global keys that should work even when search is focused
		switch {
		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.CommandPalette):
			searchActive := m.nav.CurrentType() == navigation.SearchType
			m.commandPalette.Show(m.width, m.height, searchActive)
			m.state = xplrCmdPaletteView
			return m, nil
		case key.Matches(msg, m.keymap.NextGQLType):
			cmds = m.cycleGQLType(true, cmds)
		case key.Matches(msg, m.keymap.PrevGQLType):
			cmds = m.cycleGQLType(false, cmds)
		default:
			// Delegate to appropriate handler based on search focus state
			if m.nav.CurrentType() == navigation.SearchType && m.searchFocused {
				return m.handleSearchFocused(msg)
			}
			return m.handleNormal(msg, cmds)
		}
	case components.OpenPanelMsg:
		m.handleOpenPanel(msg.Panel)
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}

	m.sizePanels()

	// Update visible panels in the stack
	var cmd tea.Cmd

	// Only the left (focused) panel receives input; right panel is display-only
	shouldReceiveMsg := m.shouldFocusedPanelReceiveMessage(msg)
	if shouldReceiveMsg && m.nav.CurrentPanel() != nil {
		currentPanel := m.nav.CurrentPanel()
		_, cmd = currentPanel.Update(msg)
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
				m.overlay.Show(content, m.width, m.height)
				m.state = xplrOverlayView
			}
		}
	}
}

// cycleGQLType handles cycling between GQL types (forward or backward)
func (m *Model) cycleGQLType(forward bool, cmds []tea.Cmd) []tea.Cmd {
	// Blur search input when switching away from Search tab
	if m.searchFocused {
		m.searchFocused = false
		m.searchInput.Blur()
		m.updateKeybindings()
	}

	// Cycle type
	if forward {
		m.nav.CycleTypeForward()
	} else {
		m.nav.CycleTypeBackward()
	}
	m.resetAndLoadMainPanel()

	// Focus search input when switching to Search tab
	if m.nav.CurrentType() == navigation.SearchType {
		m.searchFocused = true
		m.updateKeybindings()
		cmds = append(cmds, m.searchInput.Focus())
	}

	return cmds
}

// handleSearchFocused handles messages when the search input is focused
func (m *Model) handleSearchFocused(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.SearchSubmit):
			// Execute search and transfer focus to results
			query := m.searchInput.Value()
			if query != "" {
				m.searchFocused = false
				m.searchInput.Blur()
				m.updateKeybindings()
				cmds = append(cmds, m.executeSearch(query))
			}
			return *m, tea.Batch(cmds...)
		case key.Matches(msg, m.keymap.SearchClear):
			// Clear input and keep focus
			m.searchInput.SetValue("")
			return *m, nil
		default:
			// Pass message to search input
			var cmd tea.Cmd
			m.searchInput, cmd = m.searchInput.Update(msg)
			return *m, cmd
		}
	}

	return *m, nil
}

// handleNormal handles messages in normal mode (when search is not focused)
func (m *Model) handleNormal(msg tea.Msg, cmds []tea.Cmd) (Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m.updateFocusedPanel(msg, cmds)
	}

	// Handle Search tab specific keys - early return if handled
	if m.nav.CurrentType() == navigation.SearchType && key.Matches(keyMsg, m.keymap.SearchFocus) && !m.searchFocused {
		m.searchFocused = true
		m.updateKeybindings()
		cmds = append(cmds, m.searchInput.Focus())
		return *m, tea.Batch(cmds...)
	}

	// Handle remaining global keys
	switch {
	case key.Matches(keyMsg, m.keymap.ToggleOverlay):
		m.openOverlayForSelectedItem()
	case key.Matches(keyMsg, m.keymap.NextPanel):
		cmds = m.handleNextPanel(cmds)
	case key.Matches(keyMsg, m.keymap.PrevPanel):
		m.handlePrevPanel()
	}

	return m.updateFocusedPanel(msg, cmds)
}

// handleNextPanel moves forward in the navigation stack
func (m *Model) handleNextPanel(cmds []tea.Cmd) []tea.Cmd {
	if !m.nav.NavigateForward() {
		return cmds
	}

	m.updatePanelFocusStates()

	// Open up child panel for ResultType if it exists
	focusedPanel := m.nav.CurrentPanel()
	if focusedPanel == nil {
		return cmds
	}

	if openCmd := focusedPanel.OpenSelectedItem(); openCmd != nil {
		cmds = append(cmds, openCmd)
	}

	return cmds
}

// handlePrevPanel moves backward in the navigation stack
func (m *Model) handlePrevPanel() {
	if m.nav.NavigateBackward() {
		m.updatePanelFocusStates()
	}
}

// updateFocusedPanel updates the currently focused panel with the message
func (m *Model) updateFocusedPanel(msg tea.Msg, cmds []tea.Cmd) (Model, tea.Cmd) {
	shouldReceiveMsg := m.shouldFocusedPanelReceiveMessage(msg)
	if !shouldReceiveMsg || m.nav.CurrentPanel() == nil {
		return *m, tea.Batch(cmds...)
	}

	currentPanel := m.nav.CurrentPanel()
	_, cmd := currentPanel.Update(msg)
	cmds = append(cmds, cmd)

	return *m, tea.Batch(cmds...)
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

	// Reserve space for search input if on Search tab
	const searchInputHeight = 3 // Height for search input field
	if m.nav.CurrentType() == navigation.SearchType {
		panelHeight -= searchInputHeight
	}

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
			panelHeight-m.Styles.BlurredPanel.GetVerticalFrameSize(),
		)
	}
}

// updateKeybindings enables or disables key bindings based on search focus state
func (m *Model) updateKeybindings() {
	if m.searchFocused {
		// When search is focused, disable panel navigation keys
		m.keymap.NextPanel.SetEnabled(false)
		m.keymap.PrevPanel.SetEnabled(false)
		m.keymap.ToggleOverlay.SetEnabled(false)
		// Keep global keys enabled (Quit, ToggleGQLType, etc.)
		m.keymap.Quit.SetEnabled(true)
		m.keymap.NextGQLType.SetEnabled(true)
		m.keymap.PrevGQLType.SetEnabled(true)
	} else {
		// Normal mode: enable all keys
		m.keymap.NextPanel.SetEnabled(true)
		m.keymap.PrevPanel.SetEnabled(true)
		m.keymap.ToggleOverlay.SetEnabled(true)
		m.keymap.Quit.SetEnabled(true)
		m.keymap.NextGQLType.SetEnabled(true)
		m.keymap.PrevGQLType.SetEnabled(true)
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

// SetSearchBaseDir sets the base directory for search indexes
func (m *Model) SetSearchBaseDir(baseDir string) {
	m.searchBaseDir = baseDir
}

// executeSearch performs a search query and updates the search results
func (m *Model) executeSearch(query string) tea.Cmd {
	if query == "" || m.SchemaID == "" || m.searchBaseDir == "" {
		return nil
	}

	return func() tea.Msg {
		// Create searcher
		searcher := search.NewSearcher(m.searchBaseDir)
		defer searcher.Close()

		// Perform search
		results, err := searcher.Search(m.SchemaID, query, 50)
		if err != nil {
			// If index doesn't exist, need to create it
			// For now, return empty results - indexing should be handled separately
			return searchResultsMsg{results: nil, err: err}
		}

		return searchResultsMsg{results: results, err: nil}
	}
}

// searchResultsMsg is sent when search results are ready
type searchResultsMsg struct {
	results []search.SearchResult
	err     error
}

// convertSearchResultsToListItems converts search results to list items
func (m *Model) convertSearchResultsToListItems(results []search.SearchResult) []components.ListItem {
	return adapters.AdaptSearchResults(results, &m.schema)
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
	case navigation.SearchType:
		// Use cached search results or show empty state
		if m.searchResults != nil {
			items = m.searchResults
		} else {
			items = []components.ListItem{}
		}
		title = "Search Results"
	}

	m.nav.SetCurrentPanel(components.NewPanel(items, title))
	m.updatePanelFocusStates()

	// Auto-open detail panel for the first item if available (but not for Search tab)
	if len(items) > 0 && m.nav.CurrentType() != navigation.SearchType {
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
		m.selectQueryOrMutationField(currentPanel, target.FieldName)
		return
	}

	// For other types (Object, Input, Enum, etc.), select the type in the current panel
	if !currentPanel.SelectItemByName(target.TypeName) {
		return
	}

	// If a field name is specified, navigate forward and select the field
	if target.FieldName != "" {
		m.selectTypeField(currentPanel, target.FieldName)
	}
}

// selectQueryOrMutationField selects a field in Query or Mutation type panels
func (m *Model) selectQueryOrMutationField(panel *components.Panel, fieldName string) {
	if fieldName == "" {
		return
	}

	// Select the field directly in the current panel
	if !panel.SelectItemByName(fieldName) {
		return
	}

	// Navigate forward to show the field's details
	openCmd := panel.OpenSelectedItem()
	if openCmd == nil {
		return
	}

	msg, ok := openCmd().(components.OpenPanelMsg)
	if !ok {
		return
	}

	m.handleOpenPanel(msg.Panel)
	if m.nav.NavigateForward() {
		m.updatePanelFocusStates()
	}
}

// selectTypeField opens a type's panel and selects a specific field within it
func (m *Model) selectTypeField(panel *components.Panel, fieldName string) {
	// Open child panel for the selected item
	openCmd := panel.OpenSelectedItem()
	if openCmd == nil {
		return
	}

	// Execute the command to populate the next panel
	msg, ok := openCmd().(components.OpenPanelMsg)
	if !ok {
		return
	}

	// Add the new panel to the stack first
	m.handleOpenPanel(msg.Panel)

	// Then navigate forward (this adds breadcrumb from current panel's selected item)
	if !m.nav.NavigateForward() {
		return
	}

	m.updatePanelFocusStates()

	// Select the field in the newly opened panel (now at current position)
	currentPanel := m.nav.CurrentPanel()
	if currentPanel != nil {
		currentPanel.SelectItemByName(fieldName)
	}
}
