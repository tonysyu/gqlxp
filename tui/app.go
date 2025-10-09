package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tonysyu/igq/tui/adapters"
)

func Start(schema adapters.SchemaView) (tea.Model, error) {
	p := tea.NewProgram(newModel(schema))
	return p.Run()
}
