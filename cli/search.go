package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/library"
	"github.com/tonysyu/gqlxp/search"
	"github.com/tonysyu/gqlxp/utils/terminal"
	"github.com/urfave/cli/v3"
)

var (
	headerStyle = lipgloss.NewStyle().Foreground(terminal.ColorDimMagenta)
	codeStyle   = lipgloss.NewStyle().Foreground(terminal.ColorDimIndigo)
)

// searchCommand creates the search subcommand
func searchCommand() *cli.Command {
	return &cli.Command{
		Name:      "search",
		Usage:     "Search for types and fields in a GraphQL schema",
		ArgsUsage: "<query>",
		Description: `Searches for types and fields matching the given query.

Uses default schema when --schema is not specified.
Use 'gqlxp library default' to set the default schema.

Examples:
  gqlxp search user                      # Uses default schema
  gqlxp search -s examples/github.graphqls user # Uses specific file
  gqlxp search -s github-api user        # Uses library ID
  gqlxp search -s github-api "mutation"`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "schema",
				Aliases: []string{"s"},
				Usage:   "Schema file path or library ID to search",
			},
			&cli.IntFlag{
				Name:  "limit",
				Usage: "maximum number of results to return",
				Value: 30,
			},
			&cli.BoolFlag{
				Name:  "no-pager",
				Usage: "disable pager and show directly to stdout",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 1 {
				return fmt.Errorf("requires exactly 1 argument: <query>")
			}

			query := cmd.Args().First()
			limit := cmd.Int("limit")
			noPager := cmd.Bool("no-pager")

			// Get schema (empty string for default when no flag specified)
			schemaArg := cmd.String("schema")

			// Resolve schema argument (path, ID, or default)
			schema, err := resolveSchemaFromArgument(schemaArg)
			if err != nil {
				return err
			}

			schemaID := schema.ID
			content := schema.Content

			// Get schemas directory for indexing
			schemasDir, err := library.GetSchemasDir()
			if err != nil {
				return fmt.Errorf("failed to get schemas directory: %w", err)
			}

			indexer := search.NewIndexer(schemasDir)
			defer indexer.Close()

			// Create index if it doesn't exist
			if !indexer.Exists(schemaID) {
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
				return fmt.Errorf("search failed: %w (try using 'gqlxp library reindex %s')", err, schemaID)
			}

			// Display results
			if len(results) == 0 {
				fmt.Printf("No results found for query: %q\n", query)
				return nil
			}

			var maxLimitInfo string
			if len(results) == limit {
				maxLimitInfo = fmt.Sprintf(" (increase search %s for more)", codeStyle.Render("--limit N"))
			}
			// Multiple results - show list and let user choose
			var output strings.Builder
			fmt.Fprintf(&output, "Found %d results for %q%s:\n\n", len(results), query, maxLimitInfo)

			for i, result := range results {
				// Highlight the type in pink
				fmt.Fprintf(&output, "%d. %s %s\n", i+1, headerStyle.Render(result.Path), "("+result.Type+")")
				if result.Description != "" {
					fmt.Fprintf(&output, "   %s\n", result.Description)
				}
				showCmd := formatShowCommand(schemaID, result)
				// Highlight the show command in purple
				fmt.Fprintf(&output, "   More info: %s\n", codeStyle.Render(showCmd))
			}

			// Use pager if content is long enough and not disabled
			rendered := output.String()
			if terminal.ShouldUsePager(rendered, noPager) {
				return terminal.ShowInPager(rendered)
			}

			fmt.Print(rendered)
			return nil
		},
	}
}

// formatShowCommand generates a show command for a search result
// Query and Mutation fields use the operation.field format (e.g., "Query.user")
// Other types use the parent type/directive name
func formatShowCommand(schemaID string, result search.SearchResult) string {
	path := result.Path
	showPath := ""

	// Determine what to show based on the path
	if strings.HasPrefix(path, "Query.") || strings.HasPrefix(path, "Mutation.") {
		// Only include operation and field (e.g., "Query.user", not "Query.user.name")
		parts := strings.Split(path, ".")
		if len(parts) >= 2 {
			showPath = parts[0] + "." + parts[1]
		} else {
			showPath = path
		}
	} else {
		// For other types, use the first part (type name or directive)
		parts := strings.Split(path, ".")
		showPath = parts[0]
	}

	// Use new flag-based syntax
	return fmt.Sprintf("gqlxp show --schema %s %s", schemaID, showPath)
}
