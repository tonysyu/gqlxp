package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tonysyu/gqlxp/library"
	"github.com/tonysyu/gqlxp/tui/adapters"
	"github.com/tonysyu/gqlxp/tui/libselect"
	"github.com/tonysyu/gqlxp/tui/xplr"
)

func Start(schema adapters.SchemaView) (tea.Model, error) {
	p := tea.NewProgram(xplr.New(schema))
	return p.Run()
}

// StartWithLibraryData starts the TUI with library metadata for favorites
func StartWithLibraryData(schema adapters.SchemaView, schemaID string, metadata library.SchemaMetadata) (tea.Model, error) {
	m := xplr.New(schema)
	m.SetSchemaID(schemaID)
	m.SetFavorites(metadata.Favorites)
	m.SetHasLibraryData(true)
	p := tea.NewProgram(m)
	return p.Run()
}

// StartSchemaSelector starts the schema selector TUI
func StartSchemaSelector() (tea.Model, error) {
	return libselect.Start()
}
