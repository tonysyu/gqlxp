package libselect

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tonysyu/gqlxp/library"
	"github.com/tonysyu/gqlxp/tui/adapters"
	"github.com/tonysyu/gqlxp/tui/config"
)

var (
	quitKeyBinding = key.NewBinding(
		key.WithKeys("ctrl+c", "ctrl+d"),
		key.WithHelp("ctrl+c", "quit"),
	)
)

// Model is the TUI for selecting a schema from the library
type Model struct {
	list   list.Model
	lib    library.Library
	styles config.Styles
	width  int
	height int
}

type schemaListItem struct {
	id          string
	displayName string
}

func (i schemaListItem) Title() string       { return i.displayName }
func (i schemaListItem) Description() string { return i.id }
func (i schemaListItem) FilterValue() string { return i.displayName + " " + i.id }

// SchemaSelectedMsg is sent when a schema is selected
type SchemaSelectedMsg struct {
	SchemaID string
	Schema   adapters.SchemaView
	Metadata library.SchemaMetadata
}

func newModel(lib library.Library) (Model, error) {
	styles := config.DefaultStyles()

	// Load schemas from library
	schemas, err := lib.List()
	if err != nil {
		return Model{}, fmt.Errorf("failed to load schemas: %w", err)
	}

	// Convert to list items
	items := make([]list.Item, len(schemas))
	for i, schema := range schemas {
		items[i] = schemaListItem{
			id:          schema.ID,
			displayName: schema.DisplayName,
		}
	}

	delegate := list.NewDefaultDelegate()
	listModel := list.New(items, delegate, 0, 0)
	listModel.Title = "Select a Schema"
	listModel.SetShowStatusBar(false)
	listModel.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "select"),
			),
		}
	}

	return Model{
		list:   listModel,
		lib:    lib,
		styles: styles,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, quitKeyBinding):
			return m, tea.Quit
		case msg.String() == "enter":
			// Load selected schema
			if item, ok := m.list.SelectedItem().(schemaListItem); ok {
				return m, m.loadSchema(item.id)
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-2)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
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
	return m.list.View()
}

// Start starts the schema selector TUI
func Start() (tea.Model, error) {
	lib := library.NewLibrary()
	model, err := newModel(lib)
	if err != nil {
		return nil, err
	}
	p := tea.NewProgram(model)
	return p.Run()
}
