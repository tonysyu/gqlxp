package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tonysyu/gqlxp/internal/prompt"
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
		// No arguments - check library
		lib := library.NewLibrary()
		schemas, err := lib.List()
		if err != nil {
			abort(fmt.Sprintf("Error checking library: %v", err))
		}

		if len(schemas) == 0 {
			abort("No schema file provided. Usage: gqlxp <schema-file>")
		}

		// Library has schemas - open selector
		if _, err := tui.StartSchemaSelector(); err != nil {
			abort(fmt.Sprintf("Error starting library selector: %v", err))
		}
		return
	}

	command := os.Args[1]

	// Handle direct file mode
	if !strings.HasPrefix(command, "--") {
		schemaFile := command
		loadAndStartFromFile(schemaFile)
		return
	}

	showUsage()
	os.Exit(1)
}

func showUsage() {
	usage := `Usage:
  gqlxp <schema-file>              Load and explore schema from file
  gqlxp                            Select from library (if not empty)

Schema files are automatically saved to your library on first use.
When loading a previously imported file, you'll be prompted to update
if changes are detected.

Use the TUI interface to manage library schemas (remove, view, etc).

Examples:
  gqlxp schema.graphqls            # Load schema (prompts for library details on first use)
  gqlxp                            # Open library selector
`
	fmt.Print(usage)
}

func loadAndStartFromFile(schemaFile string) {
	// Resolve schema source through library (automatic integration)
	schemaID, content, err := resolveSchemaSource(schemaFile)
	if err != nil {
		abort(fmt.Sprintf("Error resolving schema: %v", err))
	}

	// Parse schema
	schema, err := adapters.ParseSchema(content)
	if err != nil {
		abort(fmt.Sprintf("Error parsing schema: %v", err))
	}

	// Get library metadata
	lib := library.NewLibrary()
	libSchema, err := lib.Get(schemaID)
	if err != nil {
		abort(fmt.Sprintf("Error loading schema metadata: %v", err))
	}

	// Start with library data
	if _, err := tui.StartWithLibraryData(schema, schemaID, libSchema.Metadata); err != nil {
		abort(fmt.Sprintf("Error starting tui: %v", err))
	}
}

func loadSchemaFromFile(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}
	return content, nil
}

func resolveSchemaSource(filePath string) (schemaID string, content []byte, err error) {
	lib := library.NewLibrary()

	// Normalize to absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Load file content
	content, err = loadSchemaFromFile(absPath)
	if err != nil {
		return "", nil, err
	}

	// Calculate file hash
	fileHash := library.CalculateFileHash(content)
	existingSchema, err := lib.FindByPath(absPath)

	// No match - register new schema
	if err != nil {
		schemaID, err := registerSchema(absPath, content)
		return schemaID, content, err
	}

	// Hash matches - use existing schema
	if existingSchema.Metadata.FileHash == fileHash {
		return existingSchema.ID, existingSchema.Content, nil
	}

	// Hash mismatch - handle update workflow
	return handleSchemaUpdate(lib, existingSchema, content)
}

func handleSchemaUpdate(lib library.Library, existingSchema *library.Schema, newContent []byte) (string, []byte, error) {
	fmt.Printf("Schema file has changed since last import.\n")
	update, err := prompt.PromptYesNo("Update library")
	if err != nil {
		return "", nil, fmt.Errorf("failed to get user input: %w", err)
	}

	// User chose not to update
	if !update {
		fmt.Printf("Using existing library version\n")
		return existingSchema.ID, existingSchema.Content, nil
	}

	// Update library content
	if err := lib.UpdateContent(existingSchema.ID, newContent); err != nil {
		return "", nil, fmt.Errorf("failed to update library: %w", err)
	}

	fmt.Printf("Library schema '%s' updated\n", existingSchema.ID)
	return existingSchema.ID, newContent, nil
}

func registerSchema(filePath string, content []byte) (string, error) {
	lib := library.NewLibrary()

	// Generate suggested ID from filename
	basename := filepath.Base(filePath)
	ext := filepath.Ext(basename)
	suggested := strings.TrimSuffix(basename, ext)
	suggested = library.SanitizeSchemaID(suggested)

	// Prompt for schema ID
	schemaID, err := prompt.PromptSchemaID(suggested)
	if err != nil {
		return "", fmt.Errorf("failed to get schema ID: %w", err)
	}

	// Prompt for display name
	displayName, err := prompt.PromptString("Enter display name", schemaID)
	if err != nil {
		return "", fmt.Errorf("failed to get display name: %w", err)
	}

	// Add to library
	if err := lib.Add(schemaID, displayName, filePath); err != nil {
		return "", fmt.Errorf("failed to add schema to library: %w", err)
	}

	fmt.Printf("Schema '%s' added to library\n", schemaID)
	return schemaID, nil
}

func abort(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
