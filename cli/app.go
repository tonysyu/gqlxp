package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tonysyu/gqlxp/library"
	"github.com/tonysyu/gqlxp/tui"
	"github.com/tonysyu/gqlxp/tui/adapters"
)

// appCommand creates the app subcommand for launching the TUI.
func appCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app",
		Short: "Launch the GraphQL schema explorer TUI",
		Long: `Opens the interactive TUI for exploring GraphQL schemas.

With no schema flag, opens the library selector to choose from saved schemas.
Use --schema/-s to open a specific schema file or library ID.`,
		Example: `  gqlxp app                             # Open library selector
  gqlxp app -s examples/github.graphqls # Open specific schema file
  gqlxp app -s github-api               # Open schema from library`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeTUICommand(cmd)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Flags().StringP("log-file", "l", "", "Enable debug logging to `FILE`")
	cmd.Flags().String("select", "", "Pre-select TYPE or TYPE.FIELD in TUI")

	return cmd
}

// executeTUICommand is the shared logic for launching the TUI,
// used by both the root command and app subcommand.
func executeTUICommand(cmd *cobra.Command) error {
	logFile, _ := cmd.Flags().GetString("log-file")
	setupLogging(logFile)

	schemaArg, _ := cmd.Flags().GetString("schema")
	selectTarget, _ := cmd.Flags().GetString("select")

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
	schema, err := LoadSchema(schemaFile)
	if err != nil {
		return fmt.Errorf("error resolving schema: %w", err)
	}

	schemaView := adapters.NewSchemaView(schema.GQLSchema)

	// Ensure search index exists before launching TUI
	lib := library.NewLibrary()
	_ = lib.EnsureIndex(schema.ID, &schema.GQLSchema)

	// Load metadata for TUI
	libSchema, err := lib.Get(schema.ID)
	if err != nil {
		return fmt.Errorf("failed to load schema metadata: %w", err)
	}

	// Start with library data
	if selectTarget != "" {
		// Parse selection target and start with selection
		typeName, fieldName := parseSelectionTarget(selectTarget)
		target := tui.SelectionTarget{
			TypeName:  typeName,
			FieldName: fieldName,
		}
		if _, err := tui.StartWithSelection(schemaView, schema.ID, libSchema.Metadata, target); err != nil {
			return fmt.Errorf("error starting tui: %w", err)
		}
	} else {
		// Start normally without selection
		if _, err := tui.StartWithLibraryData(schemaView, schema.ID, libSchema.Metadata); err != nil {
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
