package library

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tonysyu/gqlxp/library"
)

func defaultCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "default [schema-id]",
		Short: "Set or show the default schema",
		Long:  "Sets or displays the default schema to use when no schema is specified.",
		Example: `  gqlxp library default           # Show current default
  gqlxp library default github    # Set default to 'github'
  gqlxp library default --clear   # Clear default setting`,
		RunE: func(cmd *cobra.Command, args []string) error {
			lib := library.NewLibrary()

			// Clear default if --clear is used
			clear, _ := cmd.Flags().GetBool("clear")
			if clear {
				if err := lib.SetDefaultSchema(""); err != nil {
					return fmt.Errorf("failed to clear default schema: %w", err)
				}
				fmt.Println("Default schema cleared")
				return nil
			}

			// No arguments - show current default
			if len(args) == 0 {
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
			schemaID := args[0]

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
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Flags().Bool("clear", false, "clear the default schema setting")

	return cmd
}
