package cli

import (
	"context"
	"fmt"

	"github.com/tonysyu/gqlxp/library"
	"github.com/urfave/cli/v3"
)

// configCommand creates the config subcommand
func configCommand() *cli.Command {
	return &cli.Command{
		Name:  "config",
		Usage: "Manage gqlxp configuration",
		Commands: []*cli.Command{
			{
				Name:      "default-schema",
				Aliases:   []string{"s"},
				Usage:     "Set or show the default schema",
				ArgsUsage: "[schema-id-or-path]",
				Description: `Sets the default schema to use when no schema is specified.

The argument can be:
- A schema ID from the library (e.g., "github")
- A file path to a schema (will be added to library if needed)
- Omit the argument to show the current default schema

Examples:
  gqlxp config default-schema github
  gqlxp config default-schema ./schema.graphqls
  gqlxp config default-schema  # Show current default`,
				Action: func(ctx context.Context, cmd *cli.Command) error {
					lib := library.NewLibrary()

					// No arguments - show current default
					if cmd.Args().Len() == 0 {
						defaultSchema, err := lib.GetDefaultSchema()
						if err != nil {
							return fmt.Errorf("error getting default schema: %w", err)
						}

						if defaultSchema == "" {
							fmt.Println("No default schema set")
							return nil
						}

						schema, err := lib.Get(defaultSchema)
						if err != nil {
							return fmt.Errorf("error loading default schema: %w", err)
						}

						fmt.Printf("Default schema: %s (%s)\n", defaultSchema, schema.Metadata.DisplayName)
						return nil
					}

					// Set default schema
					schemaArg := cmd.Args().Get(0)
					schemaID, err := resolveSchemaArgument(schemaArg)
					if err != nil {
						return err
					}

					if err := lib.SetDefaultSchema(schemaID); err != nil {
						return fmt.Errorf("error setting default schema: %w", err)
					}

					schema, err := lib.Get(schemaID)
					if err != nil {
						return fmt.Errorf("error loading schema: %w", err)
					}

					fmt.Printf("Default schema set to: %s (%s)\n", schemaID, schema.Metadata.DisplayName)
					return nil
				},
			},
		},
	}
}

// resolveSchemaArgument resolves a schema argument to a schema ID.
// The argument can be either:
// 1. A schema ID that exists in the library
// 2. A file path (will be added to library if needed)
func resolveSchemaArgument(arg string) (string, error) {
	lib := library.NewLibrary()

	// First check if it's an existing schema ID
	if _, err := lib.Get(arg); err == nil {
		return arg, nil
	}

	// Not a schema ID - try as file path
	schemaID, _, err := resolveSchemaSource(arg)
	if err != nil {
		return "", fmt.Errorf("invalid schema argument '%s': %w", arg, err)
	}

	return schemaID, nil
}
