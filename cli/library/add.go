package library

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/tonysyu/gqlxp/gql/introspection"
	"github.com/tonysyu/gqlxp/library"
	"github.com/urfave/cli/v3"
)

func addCommand() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Usage:     "Add a schema to the library from a file or URL",
		ArgsUsage: "<schema-file-or-url>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "id",
				Usage: "schema ID (lowercase letters, numbers, hyphens)",
			},
			&cli.StringFlag{
				Name:  "name",
				Usage: "display name for the schema",
			},
			&cli.StringSliceFlag{
				Name:    "header",
				Aliases: []string{"H"},
				Usage:   "HTTP header for URL requests (e.g., 'Authorization: Bearer token')",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 1 {
				return fmt.Errorf("requires exactly 1 argument: <schema-file-or-url>")
			}

			source := cmd.Args().First()
			lib := library.NewLibrary()

			var content []byte
			var sourceInfo string
			var err error

			if introspection.IsURL(source) {
				content, err = fetchSchemaFromURL(ctx, source, cmd.StringSlice("header"))
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
			if err := validateSchema(content); err != nil {
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
			if err := lib.AddFromContent(schemaID, displayName, content, sourceInfo); err != nil {
				return err
			}

			fmt.Printf("Added schema '%s' (%s) to library\n", schemaID, displayName)
			return nil
		},
	}
}
