package library

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tonysyu/gqlxp/cli/prompt"
	"github.com/tonysyu/gqlxp/library"
)

func removeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <schema-id>",
		Short: "Remove a schema from the library",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			schemaID := args[0]
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
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Flags().Bool("force", false, "skip confirmation prompt")

	return cmd
}

func confirmSchemaRemoval(cmd *cobra.Command, schemaID string, schema *library.Schema) error {
	force, _ := cmd.Flags().GetBool("force")
	if force {
		return nil
	}

	confirm, err := prompt.YesNo(fmt.Sprintf("Remove schema '%s' (%s)?", schemaID, schema.Metadata.DisplayName))
	if err != nil {
		return fmt.Errorf("failed to get confirmation: %w", err)
	}

	if !confirm {
		fmt.Println("Cancelled")
		return fmt.Errorf("removal cancelled")
	}

	return nil
}
