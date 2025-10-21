package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tonysyu/gqlxp/tui"
	"github.com/tonysyu/gqlxp/tui/adapters"
)

func main() {
	if len(os.Args) < 2 {
		abort("Usage: gqlxp <schema-file>")
	}

	if logFile := os.Getenv("GQLXP_LOGFILE"); logFile != "" {
		f, err := tea.LogToFile(logFile, "debug")
		if err != nil {
			abort(fmt.Sprintf("Error opening log file: %v", err))
		}
		defer f.Close()
	}

	// TODO: Store cache of paths and add UI for choosing path if not defined.
	schemaFile := os.Args[1]
	schemaContent, err := os.ReadFile(schemaFile)
	if err != nil {
		abort(fmt.Sprintf("Error reading schema file '%s': %v\n", schemaFile, err))
	}

	schema, err := adapters.ParseSchema(schemaContent)
	if err != nil {
		abort(fmt.Sprintf("Error parsing schema: %v", err))
	}

	if _, err := tui.Start(schema); err != nil {
		abort(fmt.Sprintf("Error starting tui: %v", err))
	}
}

func abort(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
