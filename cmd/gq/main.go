package main

import (
	"fmt"
	"io/ioutil"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tonysyu/gq/gql"
	"github.com/tonysyu/gq/tui"
)

func main() {
	// Read schema file (assuming it's in the project root or data directory)
	schemaContent, err := ioutil.ReadFile("examples/github.graphqls")
	if err != nil {
		fmt.Printf("Error reading schema file: %v\n", err)
		os.Exit(1)
	}

	// Parse schema and get query fields
	queryFields := gql.ParseSchema(schemaContent)
	items := tui.AdaptGraphQLItems(queryFields)

	p := tea.NewProgram(tui.NewModel(items))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
