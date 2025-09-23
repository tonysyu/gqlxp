package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tonysyu/igq/gql"
)

func Start(schema gql.GraphQLSchema) (tea.Model, error) {
	p := tea.NewProgram(NewModel(schema))
	return p.Run()
}
