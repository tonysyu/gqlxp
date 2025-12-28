package cli

import (
	"context"
	"fmt"

	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/library"
	"github.com/tonysyu/gqlxp/search"
	"github.com/tonysyu/gqlxp/tui"
	"github.com/tonysyu/gqlxp/tui/adapters"
	"github.com/urfave/cli/v3"
)

// searchCommand creates the search subcommand
func searchCommand() *cli.Command {
	return &cli.Command{
		Name:      "search",
		Usage:     "Search for types and fields in a GraphQL schema",
		ArgsUsage: "[schema-file] <query>",
		Description: `Searches for types and fields matching the given query.

The schema-file argument is optional if a default schema has been set.
Use 'gqlxp config default-schema' to set the default schema.

Examples:
  gqlxp search examples/github.graphqls user
  gqlxp search user  # Uses default schema
  gqlxp search "mutation"`,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "reindex",
				Usage: "rebuild the search index before searching",
			},
			&cli.IntFlag{
				Name:  "limit",
				Usage: "maximum number of results to return",
				Value: 10,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() < 1 || cmd.Args().Len() > 2 {
				return fmt.Errorf("requires 1 or 2 arguments: [schema-file] <query>")
			}

			var schemaArg, query string
			reindex := cmd.Bool("reindex")
			limit := cmd.Int("limit")

			// Parse arguments (similar pattern to print command)
			if cmd.Args().Len() == 1 {
				query = cmd.Args().First()
			} else {
				schemaArg = cmd.Args().First()
				query = cmd.Args().Get(1)
			}

			// Get library and schema
			lib := library.NewLibrary()

			var schemaID string
			var content []byte
			var err error

			// Resolve schema source
			if schemaArg != "" {
				schemaID, content, err = resolveSchemaSource(schemaArg)
			} else {
				schemaID, content, err = resolveDefaultSchema(lib)
			}
			if err != nil {
				return err
			}

			// Get schemas directory for indexing
			schemasDir, err := library.GetSchemasDir()
			if err != nil {
				return fmt.Errorf("failed to get schemas directory: %w", err)
			}

			indexer := search.NewIndexer(schemasDir)
			defer indexer.Close()

			// Reindex if requested or index doesn't exist
			if reindex || !indexer.Exists(schemaID) {
				schema, err := gql.ParseSchema(content)
				if err != nil {
					return fmt.Errorf("failed to parse schema: %w", err)
				}

				fmt.Printf("Indexing schema '%s'...\n", schemaID)
				if err := indexer.Index(schemaID, &schema); err != nil {
					return fmt.Errorf("failed to index schema: %w", err)
				}
			}

			// Search
			searcher := search.NewSearcher(schemasDir)
			defer searcher.Close()

			results, err := searcher.Search(schemaID, query, limit)
			if err != nil {
				return fmt.Errorf("search failed: %w (try using --reindex)", err)
			}

			// Display results
			if len(results) == 0 {
				fmt.Printf("No results found for query: %q\n", query)
				return nil
			}

			// If single result, open TUI directly at that location
			if len(results) == 1 {
				return openSchemaAtPath(schemaID, content, results[0].Path)
			}

			// Multiple results - show list and let user choose
			fmt.Printf("Found %d results for %q:\n\n", len(results), query)
			for i, result := range results {
				fmt.Printf("%d. %s (%s)\n", i+1, result.Path, result.Type)
				if result.Description != "" {
					fmt.Printf("   %s\n", result.Description)
				}
			}

			// For now, just list results. In the future, add interactive selection
			return nil
		},
	}
}

// openSchemaAtPath opens the TUI at a specific path in the schema
func openSchemaAtPath(schemaID string, content []byte, path string) error {
	schema, err := adapters.ParseSchema(content)
	if err != nil {
		return fmt.Errorf("failed to parse schema: %w", err)
	}

	lib := library.NewLibrary()
	libSchema, err := lib.Get(schemaID)
	if err != nil {
		return fmt.Errorf("failed to get schema metadata: %w", err)
	}

	_, err = tui.StartWithLibraryData(schema, schemaID, libSchema.Metadata)
	return err
}

// resolveDefaultSchema gets the default schema from library config
func resolveDefaultSchema(lib library.Library) (string, []byte, error) {
	defaultID, err := lib.GetDefaultSchema()
	if err != nil {
		return "", nil, fmt.Errorf("failed to get default schema: %w", err)
	}

	if defaultID == "" {
		return "", nil, fmt.Errorf("no schema-file specified and no default schema set. Use 'gqlxp config default-schema <id>' to set a default")
	}

	schema, err := lib.Get(defaultID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to load default schema '%s': %w", defaultID, err)
	}

	return defaultID, schema.Content, nil
}
