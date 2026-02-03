package library

import (
	"context"
	"fmt"

	"github.com/tonysyu/gqlxp/library"
	"github.com/urfave/cli/v3"
)

func defaultCommand() *cli.Command {
	return &cli.Command{
		Name:      "default",
		Usage:     "Set or show the default schema",
		ArgsUsage: "[schema-id]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "clear",
				Usage: "clear the default schema setting",
			},
		},
		Description: `Sets or displays the default schema to use when no schema is specified.

Examples:
  gqlxp library default           # Show current default
  gqlxp library default github    # Set default to 'github'
  gqlxp library default --clear   # Clear default setting`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			lib := library.NewLibrary()

			// Clear default if --clear is used
			if cmd.Bool("clear") {
				if err := lib.SetDefaultSchema(""); err != nil {
					return fmt.Errorf("failed to clear default schema: %w", err)
				}
				fmt.Println("Default schema cleared")
				return nil
			}

			// No arguments - show current default
			if cmd.Args().Len() == 0 {
				defaultID, err := lib.GetDefaultSchema()
				if err != nil {
					return fmt.Errorf("failed to get default schema: %w", err)
				}

				if defaultID == "" {
					fmt.Println("No default schema set")
					return nil
				}

				schema, err := lib.Get(defaultID)
				if err != nil {
					return fmt.Errorf("failed to load default schema: %w", err)
				}

				fmt.Printf("Default schema: %s (%s)\n", defaultID, schema.Metadata.DisplayName)
				return nil
			}

			// Set default schema
			schemaID := cmd.Args().First()

			// Verify schema exists
			schema, err := lib.Get(schemaID)
			if err != nil {
				return schemaNotFoundError(lib, schemaID)
			}

			if err := lib.SetDefaultSchema(schemaID); err != nil {
				return fmt.Errorf("failed to set default schema: %w", err)
			}

			fmt.Printf("Default schema set to: %s (%s)\n", schemaID, schema.Metadata.DisplayName)
			return nil
		},
	}
}
