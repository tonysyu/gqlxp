package library

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tonysyu/gqlxp/cli/prompt"
	"github.com/tonysyu/gqlxp/gql/introspection"
	"github.com/tonysyu/gqlxp/library"
)

func addCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <schema-file-or-url>",
		Short: "Add a schema to the library from a file or URL",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			source := args[0]
			lib := library.NewLibrary()

			var content []byte
			var sourceInfo string
			var err error

			headers, _ := cmd.Flags().GetStringArray("header")

			if introspection.IsURL(source) {
				content, err = fetchSchemaFromURL(ctx, source, headers)
				if err != nil {
					return err
				}
				sourceInfo = source
			} else {
				absPath, err := filepath.Abs(source)
				if err != nil {
					return fmt.Errorf("failed to resolve absolute path: %w", err)
				}
				content, err = loadSchemaFromFile(absPath)
				if err != nil {
					return err
				}
				sourceInfo = absPath
			}

			// Validate it's a valid GraphQL schema
			if err := ValidateSchema(content); err != nil {
				return err
			}

			// Get schema ID (from flag or prompt)
			schemaID, err := getSchemaID(cmd, source)
			if err != nil {
				return err
			}

			// Get display name (from flag or prompt)
			displayName, err := getDisplayName(cmd, schemaID)
			if err != nil {
				return err
			}

			// Add to library
			err = lib.AddFromContent(schemaID, displayName, content, sourceInfo)
			if err != nil {
				if !errors.Is(err, library.ErrSchemaExists) {
					return err
				}

				// Schema already exists — ask user whether to overwrite
				overwrite, promptErr := prompt.YesNo(fmt.Sprintf("Schema '%s' already exists. Overwrite?", schemaID))
				if promptErr != nil {
					return promptErr
				}
				if !overwrite {
					return nil
				}

				if err := lib.UpdateContent(schemaID, content); err != nil {
					return err
				}

				// Update display name if it differs from existing
				existing, err := lib.Get(schemaID)
				if err != nil {
					return err
				}
				if existing.Metadata.DisplayName != displayName {
					existing.Metadata.DisplayName = displayName
					if err := lib.UpdateMetadata(schemaID, existing.Metadata); err != nil {
						return err
					}
				}

				fmt.Printf("Updated schema '%s' (%s) in library\n", schemaID, displayName)
				return nil
			}

			fmt.Printf("Added schema '%s' (%s) to library\n", schemaID, displayName)
			return nil
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Flags().String("id", "", "schema ID (lowercase letters, numbers, hyphens)")
	cmd.Flags().String("name", "", "display name for the schema")
	cmd.Flags().StringArrayP("header", "H", nil, "HTTP header for URL requests (e.g., 'Authorization: Bearer token')")

	return cmd
}
