package library

import (
	"context"
	"fmt"

	"github.com/tonysyu/gqlxp/library"
	"github.com/urfave/cli/v3"
)

func removeCommand() *cli.Command {
	return &cli.Command{
		Name:      "remove",
		Usage:     "Remove a schema from the library",
		ArgsUsage: "<schema-id>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "force",
				Usage: "skip confirmation prompt",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 1 {
				return fmt.Errorf("requires exactly 1 argument: <schema-id>")
			}

			schemaID := cmd.Args().First()
			lib := library.NewLibrary()

			// Verify schema exists
			schema, err := lib.Get(schemaID)
			if err != nil {
				return schemaNotFoundError(lib, schemaID)
			}

			// Confirm removal unless --force is used
			if err := confirmSchemaRemoval(cmd, schemaID, schema); err != nil {
				return err
			}

			// Check if this is the default schema
			defaultID, _ := lib.GetDefaultSchema()
			isDefault := schemaID == defaultID

			// Remove schema
			if err := lib.Remove(schemaID); err != nil {
				return fmt.Errorf("failed to remove schema: %w", err)
			}

			fmt.Printf("Removed schema '%s' from library\n", schemaID)

			// Clear default if necessary
			if isDefault {
				if err := lib.SetDefaultSchema(""); err != nil {
					fmt.Printf("Warning: failed to clear default schema setting: %v\n", err)
				} else {
					fmt.Println("Default schema setting cleared")
				}
			}

			return nil
		},
	}
}

func confirmSchemaRemoval(cmd *cli.Command, schemaID string, schema *library.Schema) error {
	if cmd.Bool("force") {
		return nil
	}

	confirm, err := promptYesNo(fmt.Sprintf("Remove schema '%s' (%s)?", schemaID, schema.Metadata.DisplayName))
	if err != nil {
		return fmt.Errorf("failed to get confirmation: %w", err)
	}

	if !confirm {
		fmt.Println("Cancelled")
		return fmt.Errorf("removal cancelled")
	}

	return nil
}
