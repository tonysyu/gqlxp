package library

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tonysyu/gqlxp/library"
)

func updateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [schema-file-or-url]",
		Short: "Update a schema in the library",
		Long: `Updates a schema in the library with new content.

If no schema source is provided, attempts to re-fetch from the original URL.
If the schema has no stored URL, you must provide a file path or URL.`,
		Example: `  gqlxp library update --id github                    # Re-fetch from original URL
  gqlxp library update --id github ./schema.graphqls  # Update from file
  gqlxp library update --id github https://api.example.com/graphql  # Update from URL`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			schemaID, _ := cmd.Flags().GetString("id")
			headers, _ := cmd.Flags().GetStringArray("header")
			lib := library.NewLibrary()

			// Verify schema exists
			existingSchema, err := lib.Get(schemaID)
			if err != nil {
				return schemaNotFoundError(lib, schemaID)
			}

			var content []byte
			var newSource schemaSource

			if len(args) > 0 {
				// Source provided - load from file or URL
				source := args[0]
				content, newSource, err = LoadSchemaContent(ctx, source, headers)
				if err != nil {
					return err
				}
			} else {
				// No source provided - try to re-fetch from stored URL
				if existingSchema.Metadata.SourceURL == "" {
					return fmt.Errorf("schema '%s' has no stored URL. Provide a file path or URL to update from", schemaID)
				}

				content, err = fetchSchemaFromURL(ctx, existingSchema.Metadata.SourceURL, headers)
				if err != nil {
					return fmt.Errorf("failed to fetch from stored URL: %w", err)
				}
			}

			// Validate it's a valid GraphQL schema before making any changes
			if err := ValidateSchema(content); err != nil {
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
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Flags().String("id", "", "schema ID to update (required)")
	cmd.Flags().StringArrayP("header", "H", nil, "HTTP header for URL requests (e.g., 'Authorization: Bearer token')")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}
