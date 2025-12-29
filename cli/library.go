package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/library"
	"github.com/tonysyu/gqlxp/search"
	"github.com/urfave/cli/v3"
)

func loadSchemaFromFile(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}
	return content, nil
}

func resolveSchemaSource(filePath string) (schemaID string, content []byte, err error) {
	lib := library.NewLibrary()

	// Normalize to absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Load file content
	content, err = loadSchemaFromFile(absPath)
	if err != nil {
		return "", nil, err
	}

	// Calculate file hash
	fileHash := library.CalculateFileHash(content)
	existingSchema, err := lib.FindByPath(absPath)

	// No match - register new schema
	if err != nil {
		schemaID, err := registerSchema(absPath, content)
		return schemaID, content, err
	}

	// Hash matches - use existing schema
	if existingSchema.Metadata.FileHash == fileHash {
		return existingSchema.ID, existingSchema.Content, nil
	}

	// Hash mismatch - handle update workflow
	return handleSchemaUpdate(lib, existingSchema, content)
}

func handleSchemaUpdate(lib library.Library, existingSchema *library.Schema, newContent []byte) (string, []byte, error) {
	fmt.Printf("Schema file has changed since last import.\n")
	update, err := PromptYesNo("Update library")
	if err != nil {
		return "", nil, fmt.Errorf("failed to get user input: %w", err)
	}

	// User chose not to update
	if !update {
		fmt.Printf("Using existing library version\n")
		return existingSchema.ID, existingSchema.Content, nil
	}

	// Update library content
	if err := lib.UpdateContent(existingSchema.ID, newContent); err != nil {
		return "", nil, fmt.Errorf("failed to update library: %w", err)
	}

	fmt.Printf("Library schema '%s' updated\n", existingSchema.ID)
	return existingSchema.ID, newContent, nil
}

func registerSchema(filePath string, content []byte) (string, error) {
	lib := library.NewLibrary()

	// Generate suggested ID from filename
	basename := filepath.Base(filePath)
	ext := filepath.Ext(basename)
	suggested := strings.TrimSuffix(basename, ext)
	suggested = library.SanitizeSchemaID(suggested)

	// Prompt for schema ID
	schemaID, err := PromptSchemaID(suggested)
	if err != nil {
		return "", fmt.Errorf("failed to get schema ID: %w", err)
	}

	// Prompt for display name
	displayName, err := PromptString("Enter display name", schemaID)
	if err != nil {
		return "", fmt.Errorf("failed to get display name: %w", err)
	}

	// Add to library
	if err := lib.Add(schemaID, displayName, filePath); err != nil {
		return "", fmt.Errorf("failed to add schema to library: %w", err)
	}

	fmt.Printf("Schema '%s' added to library\n", schemaID)
	return schemaID, nil
}

// libraryCommand creates the library subcommand
func libraryCommand() *cli.Command {
	return &cli.Command{
		Name:  "library",
		Usage: "Manage schema library",
		Description: `Centralized interface for managing the schema library.

Available subcommands:
  list     - List all schemas in the library
  add      - Add a schema to the library
  remove   - Remove a schema from the library
  default  - Set or show the default schema
  reindex  - Rebuild search indexes`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Default action is to list schemas
			return listCommand().Run(ctx, []string{})
		},
		Commands: []*cli.Command{
			listCommand(),
			addCommand(),
			removeCommand(),
			defaultCommand(),
			reindexCommand(),
		},
	}
}

func listCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List all schemas in the library",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			lib := library.NewLibrary()
			schemas, err := lib.List()
			if err != nil {
				return fmt.Errorf("failed to list schemas: %w", err)
			}

			if len(schemas) == 0 {
				fmt.Println("No schemas in library. Add one with: gqlxp library add <schema-file>")
				return nil
			}

			// Get default schema to mark it
			defaultID, _ := lib.GetDefaultSchema()

			fmt.Println("Schemas in library:")
			for _, schema := range schemas {
				marker := " "
				if schema.ID == defaultID {
					marker = "*"
				}
				fmt.Printf("%s %s (%s)\n", marker, schema.ID, schema.DisplayName)
			}

			return nil
		},
	}
}

func addCommand() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Usage:     "Add a schema to the library",
		ArgsUsage: "<schema-file>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "id",
				Usage: "schema ID (lowercase letters, numbers, hyphens)",
			},
			&cli.StringFlag{
				Name:  "name",
				Usage: "display name for the schema",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 1 {
				return fmt.Errorf("requires exactly 1 argument: <schema-file>")
			}

			filePath := cmd.Args().First()
			lib := library.NewLibrary()

			// Normalize to absolute path
			absPath, err := filepath.Abs(filePath)
			if err != nil {
				return fmt.Errorf("failed to resolve absolute path: %w", err)
			}

			// Load file content to validate it exists and is readable
			content, err := loadSchemaFromFile(absPath)
			if err != nil {
				return err
			}

			// Validate it's a valid GraphQL schema
			if _, err := gql.ParseSchema(content); err != nil {
				return fmt.Errorf("invalid GraphQL schema: %w", err)
			}

			// Get or prompt for schema ID
			var schemaID string
			if cmd.String("id") != "" {
				schemaID = cmd.String("id")
				if err := library.ValidateSchemaID(schemaID); err != nil {
					return err
				}
			} else {
				// Generate suggested ID from filename
				basename := filepath.Base(filePath)
				ext := filepath.Ext(basename)
				suggested := strings.TrimSuffix(basename, ext)
				suggested = library.SanitizeSchemaID(suggested)

				schemaID, err = PromptSchemaID(suggested)
				if err != nil {
					return fmt.Errorf("failed to get schema ID: %w", err)
				}
			}

			// Get or prompt for display name
			var displayName string
			if cmd.String("name") != "" {
				displayName = cmd.String("name")
			} else {
				displayName, err = PromptString("Enter display name", schemaID)
				if err != nil {
					return fmt.Errorf("failed to get display name: %w", err)
				}
			}

			// Add to library
			if err := lib.Add(schemaID, displayName, absPath); err != nil {
				return err
			}

			fmt.Printf("Added schema '%s' (%s) to library\n", schemaID, displayName)
			return nil
		},
	}
}

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
				return fmt.Errorf("schema '%s' not found in library", schemaID)
			}

			// Confirm removal unless --force is used
			if !cmd.Bool("force") {
				confirm, err := PromptYesNo(fmt.Sprintf("Remove schema '%s' (%s)?", schemaID, schema.Metadata.DisplayName))
				if err != nil {
					return fmt.Errorf("failed to get confirmation: %w", err)
				}
				if !confirm {
					fmt.Println("Cancelled")
					return nil
				}
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
				return fmt.Errorf("schema '%s' not found in library", schemaID)
			}

			if err := lib.SetDefaultSchema(schemaID); err != nil {
				return fmt.Errorf("failed to set default schema: %w", err)
			}

			fmt.Printf("Default schema set to: %s (%s)\n", schemaID, schema.Metadata.DisplayName)
			return nil
		},
	}
}

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

			// Get schemas directory for indexing
			schemasDir, err := library.GetSchemasDir()
			if err != nil {
				return fmt.Errorf("failed to get schemas directory: %w", err)
			}

			indexer := search.NewIndexer(schemasDir)
			defer indexer.Close()

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
					if err := reindexSchema(lib, indexer, schemaInfo.ID); err != nil {
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
			return reindexSchema(lib, indexer, schemaID)
		},
	}
}

func reindexSchema(lib library.Library, indexer search.Indexer, schemaID string) error {
	// Get schema
	schema, err := lib.Get(schemaID)
	if err != nil {
		return fmt.Errorf("schema '%s' not found in library", schemaID)
	}

	fmt.Printf("Reindexing '%s'...\n", schemaID)

	// Parse schema
	parsedSchema, err := gql.ParseSchema(schema.Content)
	if err != nil {
		return fmt.Errorf("failed to parse schema: %w", err)
	}

	// Remove old index
	if err := indexer.Remove(schemaID); err != nil {
		// Ignore errors - index might not exist
	}

	// Create new index
	if err := indexer.Index(schemaID, &parsedSchema); err != nil {
		return fmt.Errorf("failed to index schema: %w", err)
	}

	fmt.Printf("Index rebuilt successfully for '%s'\n", schemaID)
	return nil
}
