package cli

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/gql/introspection"
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

// resolveSchemaFromArgument resolves a schema argument to a Schema.
// The argument can be:
// 1. Empty string - use default schema from config
// 2. A schema ID that exists in the library
// 3. A file path (will be added to library if needed)
func resolveSchemaFromArgument(arg string) (*library.Schema, error) {
	lib := library.NewLibrary()

	var schemaID string

	// Empty argument - use default schema
	if arg == "" {
		defaultSchemaID, err := lib.GetDefaultSchema()
		if err != nil {
			return nil, fmt.Errorf("error getting default schema: %w", err)
		}
		if defaultSchemaID == "" {
			return nil, fmt.Errorf("no schema specified and no default schema set. Use 'gqlxp library default' to set one")
		}
		schemaID = defaultSchemaID
	} else {
		// First check if it's an existing schema ID
		if _, err := lib.Get(arg); err == nil {
			schemaID = arg
		} else {
			// Not a schema ID - try as file path
			resolvedID, _, err := resolveSchemaSource(arg)
			if err != nil {
				return nil, fmt.Errorf("invalid schema argument '%s': %w", arg, err)
			}
			schemaID = resolvedID
		}
	}

	// Load schema from library
	schema, err := lib.Get(schemaID)
	if err != nil {
		return nil, fmt.Errorf("failed to load schema '%s': %w", schemaID, err)
	}

	return schema, nil
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
  update   - Update a schema in the library
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
			updateCommand(),
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
				// Fetch schema from URL via introspection
				content, err = fetchSchemaFromURL(ctx, source, cmd.StringSlice("header"))
				if err != nil {
					return err
				}
				sourceInfo = source
			} else {
				// Load from file
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
			if _, err := gql.ParseSchema(content); err != nil {
				return fmt.Errorf("invalid GraphQL schema: %w", err)
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
			if _, err := gql.ParseSchema(content); err != nil {
				return fmt.Errorf("invalid GraphQL schema: %w", err)
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

// schemaSource represents the source of a schema (either URL or file path).
type schemaSource struct {
	URL      string
	FilePath string
}

// loadSchemaContent loads schema content from a file path or URL.
// Returns the content, the source info, and any error.
func loadSchemaContent(ctx context.Context, source string, headers []string) ([]byte, schemaSource, error) {
	if introspection.IsURL(source) {
		content, err := fetchSchemaFromURL(ctx, source, headers)
		if err != nil {
			return nil, schemaSource{}, err
		}
		return content, schemaSource{URL: source}, nil
	}

	// Load from file
	absPath, err := filepath.Abs(source)
	if err != nil {
		return nil, schemaSource{}, fmt.Errorf("failed to resolve absolute path: %w", err)
	}
	content, err := loadSchemaFromFile(absPath)
	if err != nil {
		return nil, schemaSource{}, err
	}
	return content, schemaSource{FilePath: absPath}, nil
}

// getSchemaID returns schema ID from flag or prompts user
func getSchemaID(cmd *cli.Command, source string) (string, error) {
	if flagID := cmd.String("id"); flagID != "" {
		if err := library.ValidateSchemaID(flagID); err != nil {
			return "", err
		}
		return flagID, nil
	}

	var suggested string
	if introspection.IsURL(source) {
		// Extract hostname from URL as suggested ID
		suggested = extractHostnameAsID(source)
	} else {
		// Generate suggested ID from filename
		basename := filepath.Base(source)
		ext := filepath.Ext(basename)
		suggested = strings.TrimSuffix(basename, ext)
	}
	suggested = library.SanitizeSchemaID(suggested)

	schemaID, err := PromptSchemaID(suggested)
	if err != nil {
		return "", fmt.Errorf("failed to get schema ID: %w", err)
	}

	return schemaID, nil
}

// extractHostnameAsID extracts the hostname from a URL and sanitizes it for use as an ID.
func extractHostnameAsID(urlStr string) string {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return "schema"
	}
	hostname := parsed.Hostname()
	// Remove common prefixes/suffixes
	hostname = strings.TrimPrefix(hostname, "api.")
	hostname = strings.TrimPrefix(hostname, "www.")
	// Take just the first part if it's a subdomain
	parts := strings.Split(hostname, ".")
	if len(parts) > 0 {
		return parts[0]
	}
	return hostname
}

// fetchSchemaFromURL fetches a GraphQL schema via introspection from the given URL.
func fetchSchemaFromURL(ctx context.Context, endpoint string, headers []string) ([]byte, error) {
	opts := introspection.DefaultClientOptions()

	// Parse and add custom headers
	if len(headers) > 0 {
		customHeaders, err := introspection.ParseHeaders(headers)
		if err != nil {
			return nil, fmt.Errorf("failed to parse headers: %w", err)
		}
		for k, v := range customHeaders {
			opts.Headers[k] = v
		}
	}

	fmt.Printf("Fetching schema from %s...\n", endpoint)

	resp, err := introspection.FetchSchema(ctx, endpoint, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch schema: %w", err)
	}

	sdl, err := introspection.ToSDL(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to SDL: %w", err)
	}

	return sdl, nil
}

// getDisplayName returns display name from flag or prompts user
func getDisplayName(cmd *cli.Command, defaultName string) (string, error) {
	if flagName := cmd.String("name"); flagName != "" {
		return flagName, nil
	}

	displayName, err := PromptString("Enter display name", defaultName)
	if err != nil {
		return "", fmt.Errorf("failed to get display name: %w", err)
	}

	return displayName, nil
}

// confirmSchemaRemoval prompts user for confirmation unless --force is used
func confirmSchemaRemoval(cmd *cli.Command, schemaID string, schema *library.Schema) error {
	if cmd.Bool("force") {
		return nil
	}

	confirm, err := PromptYesNo(fmt.Sprintf("Remove schema '%s' (%s)?", schemaID, schema.Metadata.DisplayName))
	if err != nil {
		return fmt.Errorf("failed to get confirmation: %w", err)
	}

	if !confirm {
		fmt.Println("Cancelled")
		return fmt.Errorf("removal cancelled")
	}

	return nil
}

// schemaNotFoundError returns an error with the schema ID and lists available schemas.
func schemaNotFoundError(lib library.Library, schemaID string) error {
	schemas, err := lib.List()
	if err != nil || len(schemas) == 0 {
		return fmt.Errorf("schema '%s' not found in library", schemaID)
	}

	var ids []string
	for _, s := range schemas {
		ids = append(ids, s.ID)
	}
	return fmt.Errorf("schema '%s' not found in library. Available: %s", schemaID, strings.Join(ids, ", "))
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
				return schemaNotFoundError(lib, schemaID)
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
		return schemaNotFoundError(lib, schemaID)
	}

	fmt.Printf("Reindexing '%s'...\n", schemaID)

	// Parse schema
	parsedSchema, err := gql.ParseSchema(schema.Content)
	if err != nil {
		return fmt.Errorf("failed to parse schema: %w", err)
	}

	// Remove old index (ignore errors - index might not exist)
	_ = indexer.Remove(schemaID)

	// Create new index
	if err := indexer.Index(schemaID, &parsedSchema); err != nil {
		return fmt.Errorf("failed to index schema: %w", err)
	}

	fmt.Printf("Index rebuilt successfully for '%s'\n", schemaID)
	return nil
}
