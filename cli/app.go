package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/tonysyu/gqlxp/library"
	"github.com/tonysyu/gqlxp/tui"
	"github.com/tonysyu/gqlxp/tui/adapters"
	"github.com/tonysyu/gqlxp/tui/utils"
	"github.com/urfave/cli/v3"
)

// appCommand creates the app subcommand for launching the TUI.
func appCommand() *cli.Command {
	return &cli.Command{
		Name:  "app",
		Usage: "Launch the GraphQL schema explorer TUI",
		Description: `Opens the interactive TUI for exploring GraphQL schemas.

With no schema flag, opens the library selector to choose from saved schemas.
Use --schema/-s to open a specific schema file or library ID.

Examples:
  gqlxp app                             # Open library selector
  gqlxp app -s examples/github.graphqls # Open specific schema file
  gqlxp app -s github-api               # Open schema from library`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "schema",
				Aliases: []string{"s"},
				Usage:   "Schema file path or library ID to open",
			},
			&cli.StringFlag{
				Name:    "log-file",
				Aliases: []string{"l"},
				Usage:   "Enable debug logging to `FILE`",
			},
			&cli.StringFlag{
				Name:  "select",
				Usage: "Pre-select TYPE or TYPE.FIELD in TUI",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return executeTUICommand(ctx, cmd)
		},
	}
}

// executeTUICommand is the shared logic for launching the TUI,
// used by both the root command and app subcommand.
func executeTUICommand(ctx context.Context, cmd *cli.Command) error {
	setupLogging(cmd.String("log-file"))

	schemaArg := cmd.String("schema")
	selectTarget := cmd.String("select")

	// No schema specified - open library selector
	if schemaArg == "" {
		return openLibrarySelector()
	}

	// Load schema from file or library
	return loadAndStartFromFile(schemaArg, selectTarget)
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

func loadAndStartFromFile(schemaFile, selectTarget string) error {
	// Resolve schema argument (path, ID, or default)
	libSchema, err := resolveSchemaFromArgument(schemaFile)
	if err != nil {
		return fmt.Errorf("error resolving schema: %w", err)
	}

	// Parse schema
	schema, err := adapters.ParseSchema(libSchema.Content)
	if err != nil {
		return fmt.Errorf("error parsing schema: %w", err)
	}

	// Ensure search index exists before launching TUI
	utils.EnsureSearchIndexForSchema(libSchema.ID, schema.Schema())

	// Start with library data
	if selectTarget != "" {
		// Parse selection target and start with selection
		typeName, fieldName := parseSelectionTarget(selectTarget)
		target := tui.SelectionTarget{
			TypeName:  typeName,
			FieldName: fieldName,
		}
		if _, err := tui.StartWithSelection(schema, libSchema.ID, libSchema.Metadata, target); err != nil {
			return fmt.Errorf("error starting tui: %w", err)
		}
	} else {
		// Start normally without selection
		if _, err := tui.StartWithLibraryData(schema, libSchema.ID, libSchema.Metadata); err != nil {
			return fmt.Errorf("error starting tui: %w", err)
		}
	}
	return nil
}

// parseSelectionTarget parses a selection target string into type and field names.
// Format: "TypeName" or "TypeName.fieldName"
// Returns typeName and fieldName (empty string if no field specified).
func parseSelectionTarget(target string) (typeName, fieldName string) {
	if target == "" {
		return "", ""
	}
	parts := strings.SplitN(target, ".", 2)
	typeName = parts[0]
	if len(parts) > 1 {
		fieldName = parts[1]
	}
	return typeName, fieldName
}
