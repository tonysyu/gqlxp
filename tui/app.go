package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tonysyu/gqlxp/library"
	"github.com/tonysyu/gqlxp/tui/adapters"
)

func Start(schema adapters.SchemaView) (tea.Model, error) {
	p := tea.NewProgram(newModel(schema))
	return p.Run()
}

// StartWithLibraryData starts the TUI with library metadata for favorites
func StartWithLibraryData(schema adapters.SchemaView, schemaID string, metadata library.SchemaMetadata) (tea.Model, error) {
	m := newModel(schema)
	m.schemaID = schemaID
	m.favorites = metadata.Favorites
	m.hasLibraryData = true
	p := tea.NewProgram(m)
	return p.Run()
}
