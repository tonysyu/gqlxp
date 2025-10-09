package main

import (
	"fmt"
	"os"

	"github.com/tonysyu/igq/tui/adapters"
	"github.com/tonysyu/igq/tui"
)

func main() {
	if len(os.Args) < 2 {
		abort("Usage: igq <schema-file>")
	}

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
