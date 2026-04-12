package library

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tonysyu/gqlxp/library"
)

func reindexCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reindex [schema-id]",
		Short: "Rebuild search indexes",
		Long:  "Rebuilds search indexes for schemas in the library.",
		Example: `  gqlxp library reindex github    # Reindex 'github' schema
  gqlxp library reindex --all     # Reindex all schemas`,
		RunE: func(cmd *cobra.Command, args []string) error {
			lib := library.NewLibrary()

			// Reindex all schemas
			all, _ := cmd.Flags().GetBool("all")
			if all {
				schemas, err := lib.List()
				if err != nil {
					return fmt.Errorf("failed to list schemas: %w", err)
				}

				if len(schemas) == 0 {
					fmt.Println("No schemas in library to reindex")
					return nil
				}

				for _, schemaInfo := range schemas {
					if err := reindexSchema(lib, schemaInfo.ID); err != nil {
						fmt.Printf("Error reindexing '%s': %v\n", schemaInfo.ID, err)
						continue
					}
				}

				fmt.Printf("Reindexed %d schema(s)\n", len(schemas))
				return nil
			}

			// Reindex single schema
			if len(args) != 1 {
				return fmt.Errorf("requires exactly 1 argument: <schema-id>, or use --all flag")
			}

			schemaID := args[0]
			return reindexSchema(lib, schemaID)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Flags().Bool("all", false, "reindex all schemas in the library")

	return cmd
}

func reindexSchema(lib library.Library, schemaID string) error {
	// Verify schema exists for user-friendly error message
	if _, err := lib.Get(schemaID); err != nil {
		return schemaNotFoundError(lib, schemaID)
	}

	fmt.Printf("Reindexing '%s'...\n", schemaID)

	if err := lib.Reindex(schemaID); err != nil {
		return fmt.Errorf("failed to index schema: %w", err)
	}

	fmt.Printf("Index rebuilt successfully for '%s'\n", schemaID)
	return nil
}
