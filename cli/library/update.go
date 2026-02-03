package library

import (
	"context"
	"fmt"

	"github.com/tonysyu/gqlxp/library"
	"github.com/urfave/cli/v3"
)

func updateCommand() *cli.Command {
	return &cli.Command{
		Name:      "update",
		Usage:     "Update a schema in the library",
		ArgsUsage: "[schema-file-or-url]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "id",
				Usage:    "schema ID to update (required)",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:    "header",
				Aliases: []string{"H"},
				Usage:   "HTTP header for URL requests (e.g., 'Authorization: Bearer token')",
			},
		},
		Description: `Updates a schema in the library with new content.

If no schema source is provided, attempts to re-fetch from the original URL.
If the schema has no stored URL, you must provide a file path or URL.

Examples:
  gqlxp library update --id github                    # Re-fetch from original URL
  gqlxp library update --id github ./schema.graphqls  # Update from file
  gqlxp library update --id github https://api.example.com/graphql  # Update from URL`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			schemaID := cmd.String("id")
			lib := library.NewLibrary()

			// Verify schema exists
			existingSchema, err := lib.Get(schemaID)
			if err != nil {
				return schemaNotFoundError(lib, schemaID)
			}

			var content []byte
			var newSource schemaSource

			if cmd.Args().Len() > 0 {
				// Source provided - load from file or URL
				source := cmd.Args().First()
				content, newSource, err = loadSchemaContent(ctx, source, cmd.StringSlice("header"))
				if err != nil {
					return err
				}
			} else {
				// No source provided - try to re-fetch from stored URL
				if existingSchema.Metadata.SourceURL == "" {
					return fmt.Errorf("schema '%s' has no stored URL. Provide a file path or URL to update from", schemaID)
				}

				content, err = fetchSchemaFromURL(ctx, existingSchema.Metadata.SourceURL, cmd.StringSlice("header"))
				if err != nil {
					return fmt.Errorf("failed to fetch from stored URL: %w", err)
				}
			}

			// Validate it's a valid GraphQL schema before making any changes
			if err := validateSchema(content); err != nil {
				return err
			}

			// Calculate hash to check if content changed
			newHash := library.CalculateFileHash(content)

			if newHash == existingSchema.Metadata.FileHash {
				// Content unchanged - just update timestamp
				if err := lib.UpdateMetadata(schemaID, existingSchema.Metadata); err != nil {
					return fmt.Errorf("failed to update metadata: %w", err)
				}
				fmt.Printf("Schema '%s' is already up to date (timestamp updated)\n", schemaID)
				return nil
			}

			// Content changed - update the schema
			if err := lib.UpdateContent(schemaID, content); err != nil {
				return fmt.Errorf("failed to update schema: %w", err)
			}

			// Update source info if a new source was provided
			if newSource.URL != "" || newSource.FilePath != "" {
				schema, err := lib.Get(schemaID)
				if err != nil {
					return fmt.Errorf("failed to get updated schema: %w", err)
				}
				if newSource.URL != "" {
					schema.Metadata.SourceURL = newSource.URL
					schema.Metadata.SourceFile = "" // Clear file path when updating from URL
				} else if newSource.FilePath != "" {
					schema.Metadata.SourceFile = newSource.FilePath
					schema.Metadata.SourceURL = "" // Clear URL when updating from file
				}
				if err := lib.UpdateMetadata(schemaID, schema.Metadata); err != nil {
					return fmt.Errorf("failed to update source info: %w", err)
				}
			}

			fmt.Printf("Schema '%s' updated successfully\n", schemaID)
			return nil
		},
	}
}
