package library

import (
	"context"
	"fmt"

	"github.com/tonysyu/gqlxp/library"
	"github.com/urfave/cli/v3"
)

func reindexCommand() *cli.Command {
	return &cli.Command{
		Name:      "reindex",
		Usage:     "Rebuild search indexes",
		ArgsUsage: "[schema-id]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "all",
				Usage: "reindex all schemas in the library",
			},
		},
		Description: `Rebuilds search indexes for schemas in the library.

Examples:
  gqlxp library reindex github    # Reindex 'github' schema
  gqlxp library reindex --all     # Reindex all schemas`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			lib := library.NewLibrary()

			// Reindex all schemas
			if cmd.Bool("all") {
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
			if cmd.Args().Len() != 1 {
				return fmt.Errorf("requires exactly 1 argument: <schema-id>, or use --all flag")
			}

			schemaID := cmd.Args().First()
			return reindexSchema(lib, schemaID)
		},
	}
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
