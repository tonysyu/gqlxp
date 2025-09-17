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

	// Extract field names for the TUI
	fieldNames := make([]string, 0, len(queryFields))
	for fieldName := range queryFields {
		fieldNames = append(fieldNames, fieldName)
	}

	p := tea.NewProgram(tui.NewModel(fieldNames))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
