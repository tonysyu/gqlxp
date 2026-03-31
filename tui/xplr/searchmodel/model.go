package searchmodel

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	gosearch "github.com/tonysyu/gqlxp/search"
	"github.com/tonysyu/gqlxp/tui/adapters"
	"github.com/tonysyu/gqlxp/tui/config"
	"github.com/tonysyu/gqlxp/tui/xplr/components"
)

// ResultsReadyMsg is emitted when converted search results are ready.
// The parent stores them via StoreResults and reloads the panel.
type ResultsReadyMsg struct {
	Items []components.ListItem
}

type Model struct {
	input    components.SearchInput
	results  []components.ListItem
	schemaID string
	baseDir  string
	schema   *adapters.SchemaView
	keymap   config.MainKeymaps
}

func New(keymap config.MainKeymaps) Model {
	return Model{
		input:  components.NewSearchInput(),
		keymap: keymap,
	}
}

func (m Model) SetContext(schema *adapters.SchemaView, schemaID string) Model {
	m.schema = schema
	m.schemaID = schemaID
	return m
}

func (m Model) SetBaseDir(dir string) Model {
	m.baseDir = dir
	return m
}

// IsFocused reports whether the search input currently holds keyboard focus.
func (m Model) IsFocused() bool {
	return m.input.Focused()
}

func (m Model) Results() []components.ListItem {
	return m.results
}

func (m Model) StoreResults(items []components.ListItem) Model {
	m.results = items
	return m
}

func (m Model) Focus() (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.input, cmd = m.input.Focus()
	return m, cmd
}

func (m Model) Blur() Model {
	m.input = m.input.Blur()
	return m
}

func (m Model) View() string {
	return m.input.View()
}

func (m Model) HelpBindings() []key.Binding {
	return []key.Binding{
		m.keymap.SearchSubmit,
		m.keymap.SearchClear,
		m.keymap.NextGQLKind,
		m.keymap.PrevGQLKind,
		m.keymap.Quit,
	}
}

// HandleMsg handles key messages when the search input is focused.
// Only call this when on the Search tab and search is focused.
func (m Model) HandleMsg(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keymap.SearchSubmit):
			query := m.input.Value()
			if query != "" {
				m.input = m.input.Blur()
				return m, m.executeSearch(query)
			}
			return m, nil
		case key.Matches(msg, m.keymap.SearchClear):
			m.input = m.input.SetValue("")
			return m, nil
		default:
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

// executeSearch returns an async cmd that runs the search and emits ResultsReadyMsg.
// Schema is captured at call time so results convert against the correct schema version.
func (m Model) executeSearch(query string) tea.Cmd {
	if query == "" || m.schemaID == "" || m.baseDir == "" {
		return nil
	}
	baseDir := m.baseDir
	schemaID := m.schemaID
	schema := m.schema
	return func() tea.Msg {
		searcher := gosearch.NewSearcher(baseDir)
		defer searcher.Close()

		results, err := searcher.Search(schemaID, query, 50)
		if err != nil {
			return ResultsReadyMsg{Items: nil}
		}

		var items []components.ListItem
		if schema != nil {
			items = adapters.AdaptSearchResults(results, schema)
		}
		return ResultsReadyMsg{Items: items}
	}
}
