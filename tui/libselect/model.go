package libselect

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tonysyu/gqlxp/gql/introspection"
	"github.com/tonysyu/gqlxp/library"
	"github.com/tonysyu/gqlxp/tui/adapters"
	"github.com/tonysyu/gqlxp/tui/config"
)

// Model is the TUI for selecting a schema from the library
type Model struct {
	list   list.Model
	lib    library.Library
	styles config.Styles
	width  int
	height int
	keymap config.LibSelectKeymaps
	errMsg string
}

type schemaListItem struct {
	id          string
	displayName string
	updatedAt   time.Time
	isDefault   bool
}

func (i schemaListItem) Title() string {
	if i.displayName == "" || i.displayName == i.id {
		if i.isDefault {
			return fmt.Sprintf("%s (Default)", i.id)
		}
		return i.id
	}
	if i.isDefault {
		return fmt.Sprintf("%s (Default; id: %s)", i.displayName, i.id)
	}
	return fmt.Sprintf("%s (id: %s)", i.displayName, i.id)
}

func (i schemaListItem) Description() string {
	if i.updatedAt.IsZero() {
		return "last updated: unknown"
	}
	return "last updated: " + i.updatedAt.Format("2006-01-02 15:04")
}

func (i schemaListItem) FilterValue() string { return i.displayName + " " + i.id }

// SchemaSelectedMsg is sent when a schema is selected
type SchemaSelectedMsg struct {
	SchemaID string
	Schema   adapters.SchemaView
	Metadata library.SchemaMetadata
}

// DefaultSchemaSetMsg is sent when a schema is set as the default
type DefaultSchemaSetMsg struct {
	SchemaID string
}

// SchemaUpdatedMsg is sent when a schema is successfully updated from its source URL
type SchemaUpdatedMsg struct {
	SchemaID  string
	UpdatedAt time.Time
}

// schemaUpdateErrMsg carries an error from an update attempt
type schemaUpdateErrMsg struct {
	err error
}

// New creates a new library selection model
func New(lib library.Library) (Model, error) {
	styles := config.DefaultStyles()

	// Load schemas from library
	schemas, err := lib.List()
	if err != nil {
		return Model{}, fmt.Errorf("failed to load schemas: %w", err)
	}

	defaultID, err := lib.GetDefaultSchema()
	if err != nil {
		return Model{}, fmt.Errorf("failed to get default schema: %w", err)
	}

	// Convert to list items
	schemaItems := make([]schemaListItem, len(schemas))
	for i, schema := range schemas {
		schemaItems[i] = schemaListItem{
			id:          schema.ID,
			displayName: schema.DisplayName,
			updatedAt:   schema.UpdatedAt,
			isDefault:   schema.ID == defaultID,
		}
	}

	// Sort so default schema appears first
	sort.SliceStable(schemaItems, func(i, j int) bool {
		return schemaItems[i].isDefault && !schemaItems[j].isDefault
	})

	items := make([]list.Item, len(schemaItems))
	for i, item := range schemaItems {
		items[i] = item
	}

	keymap := config.NewLibSelectKeymaps()

	delegate := list.NewDefaultDelegate()
	listModel := list.New(items, delegate, 0, 0)
	listModel.Title = "Select a Schema"
	listModel.SetShowStatusBar(false)
	listModel.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{keymap.Select, keymap.SetDefault, keymap.UpdateSchema}
	}
	listModel.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{keymap.Select, keymap.SetDefault, keymap.UpdateSchema}
	}

	m := Model{
		list:   listModel,
		lib:    lib,
		styles: styles,
		keymap: keymap,
	}

	return m, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.Select):
			// Load selected schema
			if item, ok := m.list.SelectedItem().(schemaListItem); ok {
				return m, m.loadSchema(item.id)
			}
		case key.Matches(msg, m.keymap.SetDefault):
			if item, ok := m.list.SelectedItem().(schemaListItem); ok {
				return m, m.setDefaultSchema(item.id)
			}
		case key.Matches(msg, m.keymap.UpdateSchema):
			if item, ok := m.list.SelectedItem().(schemaListItem); ok {
				m.errMsg = ""
				m.list.SetSize(m.width, m.height-2)
				return m, m.updateSchema(item.id)
			}
		}
	case DefaultSchemaSetMsg:
		items := m.list.Items()
		for i, item := range items {
			if si, ok := item.(schemaListItem); ok {
				si.isDefault = si.id == msg.SchemaID
				items[i] = si
			}
		}
		cmd := m.list.SetItems(items)
		return m, cmd
	case SchemaUpdatedMsg:
		items := m.list.Items()
		for i, item := range items {
			if si, ok := item.(schemaListItem); ok && si.id == msg.SchemaID {
				si.updatedAt = msg.UpdatedAt
				items[i] = si
			}
		}
		cmd := m.list.SetItems(items)
		return m, cmd
	case schemaUpdateErrMsg:
		m.errMsg = msg.err.Error()
		m.list.SetSize(m.width, m.height-3)
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		listHeight := msg.Height - 2
		if m.errMsg != "" {
			listHeight--
		}
		m.list.SetSize(msg.Width, listHeight)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) setDefaultSchema(schemaID string) tea.Cmd {
	return func() tea.Msg {
		if err := m.lib.SetDefaultSchema(schemaID); err != nil {
			// TODO: Show error in UI
			return nil
		}
		return DefaultSchemaSetMsg{SchemaID: schemaID}
	}
}

func (m Model) updateSchema(schemaID string) tea.Cmd {
	return func() tea.Msg {
		schema, err := m.lib.Get(schemaID)
		if err != nil {
			return schemaUpdateErrMsg{fmt.Errorf("failed to get schema: %w", err)}
		}
		if schema.Metadata.SourceURL == "" {
			return schemaUpdateErrMsg{fmt.Errorf("schema '%s' has no URL to update from", schemaID)}
		}

		resp, err := introspection.FetchSchema(context.Background(), schema.Metadata.SourceURL, introspection.DefaultClientOptions())
		if err != nil {
			return schemaUpdateErrMsg{fmt.Errorf("failed to fetch schema: %w", err)}
		}

		content, err := introspection.ToSDL(resp)
		if err != nil {
			return schemaUpdateErrMsg{fmt.Errorf("failed to convert schema: %w", err)}
		}

		if err := m.lib.UpdateContent(schemaID, content); err != nil {
			return schemaUpdateErrMsg{fmt.Errorf("failed to update schema: %w", err)}
		}

		return SchemaUpdatedMsg{SchemaID: schemaID, UpdatedAt: time.Now()}
	}
}

func (m Model) loadSchema(schemaID string) tea.Cmd {
	return func() tea.Msg {
		schema, err := m.lib.Get(schemaID)
		if err != nil {
			// TODO: Show error in UI
			return tea.Quit()
		}

		parsedSchema, err := adapters.ParseSchema(schema.Content)
		if err != nil {
			// TODO: Show error in UI
			return tea.Quit()
		}

		return SchemaSelectedMsg{
			SchemaID: schemaID,
			Schema:   parsedSchema,
			Metadata: schema.Metadata,
		}
	}
}

func (m Model) View() string {
	if len(m.list.Items()) == 0 {
		emptyMsg := lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("No schemas in library\n\nAdd schemas with: gqlxp library add <id> <file>")
		return emptyMsg
	}
	if m.errMsg != "" {
		errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
		return lipgloss.JoinVertical(lipgloss.Left, m.list.View(), errStyle.Render(m.errMsg))
	}
	return m.list.View()
}
