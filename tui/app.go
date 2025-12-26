package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tonysyu/gqlxp/library"
	"github.com/tonysyu/gqlxp/tui/adapters"
)

func Start(schema adapters.SchemaView) (tea.Model, error) {
	m := newModelWithXplr(schema)
	p := tea.NewProgram(m)
	return p.Run()
}

// StartWithLibraryData starts the TUI with library metadata for favorites
func StartWithLibraryData(schema adapters.SchemaView, schemaID string, metadata library.SchemaMetadata) (tea.Model, error) {
	m := newModelWithXplrAndLibrary(schema, schemaID, metadata)
	p := tea.NewProgram(m)
	return p.Run()
}

// StartSchemaSelector starts the schema selector TUI
func StartSchemaSelector() (tea.Model, error) {
	m, err := newModelWithLibselect()
	if err != nil {
		return nil, err
	}
	p := tea.NewProgram(m)
	return p.Run()
}

func SetupLogging(logFile string) error {
	if logFile != "" {
		f, err := tea.LogToFile(logFile, "debug")
		if err != nil {
			return err
		}
		// Note: We can't defer here as this isn't main, but the log file
		// will be closed when the program exits
		_ = f
	}
	return nil
}
