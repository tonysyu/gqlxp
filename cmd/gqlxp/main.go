package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tonysyu/gqlxp/library"
	"github.com/tonysyu/gqlxp/tui"
	"github.com/tonysyu/gqlxp/tui/adapters"
)

func main() {
	if logFile := os.Getenv("GQLXP_LOGFILE"); logFile != "" {
		f, err := tea.LogToFile(logFile, "debug")
		if err != nil {
			abort(fmt.Sprintf("Error opening log file: %v", err))
		}
		defer f.Close()
	}

	if len(os.Args) < 2 {
		showUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	// Handle library subcommands
	if command == "library" {
		handleLibraryCommand()
		return
	}

	// Handle direct file mode (existing behavior)
	if !strings.HasPrefix(command, "--") {
		schemaFile := command
		loadAndStartFromFile(schemaFile)
		return
	}

	// Handle flags
	if command == "--library" {
		if len(os.Args) == 2 {
			// Library selection mode
			if _, err := tui.StartSchemaSelector(); err != nil {
				abort(fmt.Sprintf("Error starting library selector: %v", err))
			}
		} else if len(os.Args) == 3 {
			// Direct library schema load
			schemaID := os.Args[2]
			loadAndStartFromLibrary(schemaID)
		} else {
			showUsage()
			os.Exit(1)
		}
		return
	}

	showUsage()
	os.Exit(1)
}

func showUsage() {
	usage := `Usage:
  gqlxp <schema-file>              Load schema from file path
  gqlxp --library                  Select schema from library
  gqlxp --library <schema-id>      Load schema from library
  gqlxp library add <id> <file>    Add schema to library
  gqlxp library list               List schemas in library
  gqlxp library remove <id>        Remove schema from library

Examples:
  gqlxp schema.graphqls            # Explore schema from file
  gqlxp --library github-api       # Explore schema from library
  gqlxp library add github-api github-schema.graphqls  # Add to library
`
	fmt.Print(usage)
}

func loadAndStartFromFile(schemaFile string) {
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

func loadAndStartFromLibrary(schemaID string) {
	lib := library.NewLibrary()

	schema, err := lib.Get(schemaID)
	if err != nil {
		abort(fmt.Sprintf("Error loading schema from library: %v", err))
	}

	parsedSchema, err := adapters.ParseSchema(schema.Content)
	if err != nil {
		abort(fmt.Sprintf("Error parsing schema: %v", err))
	}

	if _, err := tui.StartWithLibraryData(parsedSchema, schemaID, schema.Metadata); err != nil {
		abort(fmt.Sprintf("Error starting tui: %v", err))
	}
}

func handleLibraryCommand() {
	if len(os.Args) < 3 {
		showUsage()
		os.Exit(1)
	}

	subcommand := os.Args[2]
	lib := library.NewLibrary()

	switch subcommand {
	case "add":
		if len(os.Args) != 5 {
			abort("Usage: gqlxp library add <id> <file-path>")
		}
		schemaID := os.Args[3]
		filePath := os.Args[4]

		// Use schema ID as display name by default
		displayName := schemaID

		if err := lib.Add(schemaID, displayName, filePath); err != nil {
			abort(fmt.Sprintf("Error adding schema: %v", err))
		}

		fmt.Printf("Schema '%s' added to library\n", schemaID)

	case "list":
		schemas, err := lib.List()
		if err != nil {
			abort(fmt.Sprintf("Error listing schemas: %v", err))
		}

		if len(schemas) == 0 {
			fmt.Println("No schemas in library")
			return
		}

		fmt.Println("Schemas in library:")
		for _, schema := range schemas {
			fmt.Printf("  %s - %s\n", schema.ID, schema.DisplayName)
		}

	case "remove":
		if len(os.Args) != 4 {
			abort("Usage: gqlxp library remove <id>")
		}
		schemaID := os.Args[3]

		if err := lib.Remove(schemaID); err != nil {
			abort(fmt.Sprintf("Error removing schema: %v", err))
		}

		fmt.Printf("Schema '%s' removed from library\n", schemaID)

	default:
		abort(fmt.Sprintf("Unknown library subcommand: %s", subcommand))
	}
}

func abort(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
