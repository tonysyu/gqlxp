package cli

import (
	"context"
	"fmt"

	"github.com/tonysyu/gqlxp/library"
	"github.com/tonysyu/gqlxp/tui"
	"github.com/tonysyu/gqlxp/tui/adapters"
	"github.com/urfave/cli/v3"
)

// appCommand creates the app subcommand for launching the TUI.
func appCommand() *cli.Command {
	return &cli.Command{
		Name:      "app",
		Usage:     "Launch the GraphQL schema explorer TUI",
		ArgsUsage: "[schema-file]",
		Description: `Opens the interactive TUI for exploring GraphQL schemas.

With no arguments, opens the library selector to choose from saved schemas.
With a schema file argument, opens the TUI with that schema loaded.

Examples:
  gqlxp app                          # Open library selector
  gqlxp app examples/github.graphqls # Open specific schema`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return executeTUICommand(ctx, cmd)
		},
	}
}

// executeTUICommand is the shared logic for launching the TUI,
// used by both the root command and app subcommand.
func executeTUICommand(ctx context.Context, cmd *cli.Command) error {
	setupLogging(cmd.String("log-file"))

	// No arguments - open library selector
	if cmd.Args().Len() == 0 {
		return openLibrarySelector()
	}

	// Load schema from file
	schemaFile := cmd.Args().First()
	return loadAndStartFromFile(schemaFile)
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
