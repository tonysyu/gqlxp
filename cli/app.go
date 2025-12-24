package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tonysyu/gqlxp/library"
	"github.com/tonysyu/gqlxp/tui"
	"github.com/tonysyu/gqlxp/tui/adapters"
	"github.com/urfave/cli/v3"
)

// NewApp creates and configures the CLI application.
func NewApp() *cli.Command {
	return &cli.Command{
		Name:  "gqlxp",
		Usage: "Explore GraphQL schemas interactively",
		Description: `Schema files are automatically saved to your library on first use.
When loading a previously imported file, you'll be prompted to update
if changes are detected.

Use the TUI interface to manage library schemas (remove, view, etc).`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log-file",
				Aliases: []string{"l"},
				Usage:   "Enable debug logging to `FILE`",
				Sources: cli.EnvVars("GQLXP_LOGFILE"),
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			setupLogging(cmd.String("log-file"))

			// No arguments - check library
			if cmd.Args().Len() == 0 {
				return openLibrarySelector()
			}

			// Load schema from file
			schemaFile := cmd.Args().First()
			return loadAndStartFromFile(schemaFile)
		},
		Commands: []*cli.Command{
			{
				Name:    "library",
				Aliases: []string{"lib"},
				Usage:   "Open the library schema selector",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					setupLogging(cmd.String("log-file"))
					return openLibrarySelector()
				},
			},
			printCommand(),
			configCommand(),
		},
	}
}

func setupLogging(logFile string) {
	if logFile != "" {
		f, err := tea.LogToFile(logFile, "debug")
		if err != nil {
			abort(fmt.Sprintf("Error opening log file: %v", err))
		}
		// Note: We can't defer here as this isn't main, but the log file
		// will be closed when the program exits
		_ = f
	}
}

func openLibrarySelector() error {
	lib := library.NewLibrary()
	schemas, err := lib.List()
	if err != nil {
		return fmt.Errorf("error checking library: %w", err)
	}

	if len(schemas) == 0 {
		abort("No schemas in library. Usage: gqlxp <schema-file>")
	}

	// Library has schemas - open selector
	if _, err := tui.StartSchemaSelector(); err != nil {
		return fmt.Errorf("error starting library selector: %w", err)
	}
	return nil
}

func loadAndStartFromFile(schemaFile string) error {
	// Resolve schema source through library (automatic integration)
	schemaID, content, err := resolveSchemaSource(schemaFile)
	if err != nil {
		return fmt.Errorf("error resolving schema: %w", err)
	}

	// Parse schema
	schema, err := adapters.ParseSchema(content)
	if err != nil {
		return fmt.Errorf("error parsing schema: %w", err)
	}

	// Get library metadata
	lib := library.NewLibrary()
	libSchema, err := lib.Get(schemaID)
	if err != nil {
		return fmt.Errorf("error loading schema metadata: %w", err)
	}

	// Start with library data
	if _, err := tui.StartWithLibraryData(schema, schemaID, libSchema.Metadata); err != nil {
		return fmt.Errorf("error starting tui: %w", err)
	}
	return nil
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
	update, err := PromptYesNo("Update library")
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
	schemaID, err := PromptSchemaID(suggested)
	if err != nil {
		return "", fmt.Errorf("failed to get schema ID: %w", err)
	}

	// Prompt for display name
	displayName, err := PromptString("Enter display name", schemaID)
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
