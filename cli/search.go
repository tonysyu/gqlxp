package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/tonysyu/gqlxp/docs"
	"github.com/tonysyu/gqlxp/library"
	"github.com/tonysyu/gqlxp/search"
	"github.com/tonysyu/gqlxp/utils/terminal"
	"github.com/tonysyu/gqlxp/utils/text"
	"github.com/urfave/cli/v3"
)

var (
	headerStyle = lipgloss.NewStyle().Foreground(terminal.ColorDimMagenta)
	codeStyle   = lipgloss.NewStyle().Foreground(terminal.ColorDimIndigo)
)

// canonicalSearchKinds maps lowercase kind names to their canonical form.
var canonicalSearchKinds = map[string]string{
	"query":          "Query",
	"mutation":       "Mutation",
	"object":         "Object",
	"input":          "Input",
	"enum":           "Enum",
	"scalar":         "Scalar",
	"interface":      "Interface",
	"union":          "Union",
	"directive":      "Directive",
	"objectfield":    "ObjectField",
	"inputfield":     "InputField",
	"interfacefield": "InterfaceField",
}

// validSearchFields are the valid field names in search queries.
var validSearchFields = []string{"kind", "name", "description", "path", "usage", "implements"}

// searchFieldPattern matches fieldname: patterns in bleve query strings.
var searchFieldPattern = regexp.MustCompile(`\b([a-zA-Z][a-zA-Z0-9]*):\S`)

// extractQueryFields returns unique field names found in a bleve query string.
func extractQueryFields(query string) []string {
	seen := map[string]bool{}
	var fields []string
	for _, m := range searchFieldPattern.FindAllStringSubmatch(query, -1) {
		f := m[1]
		if !seen[f] {
			seen[f] = true
			fields = append(fields, f)
		}
	}
	return fields
}

// unknownFieldWarnings returns warning messages for unrecognized field names in query.
func unknownFieldWarnings(query string) []string {
	validSet := make(map[string]bool, len(validSearchFields))
	for _, f := range validSearchFields {
		validSet[f] = true
	}
	var warnings []string
	for _, field := range extractQueryFields(query) {
		if validSet[field] {
			continue
		}
		msg := fmt.Sprintf("⚠️ unknown search field %q", field)
		if suggestions := text.SuggestSimilarWords(field, validSearchFields); len(suggestions) > 0 {
			msg += fmt.Sprintf(" — did you mean: %s?", strings.Join(suggestions, " or "))
		}
		warnings = append(warnings, msg)
	}
	return warnings
}

// searchCommand creates the search subcommand
func searchCommand() *cli.Command {
	return &cli.Command{
		Name:      "search",
		Usage:     "Search for types and fields in a GraphQL schema",
		ArgsUsage: "[query]",
		Description: `Searches for types and fields matching the given query.

Uses default schema when --schema is not specified.
Use 'gqlxp library default' to set the default schema.

For AI/programmatic use, add --json --no-pager for machine-readable output.
JSON output: [{"path":"Type.field","kind":"Query|Object|...","description":"...","signature":"..."}]

Query syntax:
  Plain keyword   Matches names and descriptions
  kind:<Kind>     Filter by kind (e.g., kind:Query, kind:Object)
  Combined        "+kind:Query user" filters to Query kind matching "user"

Examples:
  gqlxp search user                                  # Uses default schema
  gqlxp search -s github user --json --no-pager      # JSON output for AI use
  gqlxp search -s github --kind Query                # List all queries
  gqlxp search -s examples/github.graphqls user      # Uses specific file`,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "syntax",
				Usage: "show search syntax documentation and exit",
			},
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
				Usage: "disable pager; use for non-interactive/AI use",
			},
			&cli.BoolFlag{
				Name:  "json",
				Usage: "output results as JSON (recommended for AI/programmatic use)",
			},
			&cli.StringFlag{
				Name:  "kind",
				Usage: "filter by document kind: Query, Mutation, Object, Input, Enum, Scalar, Interface, Union, Directive, ObjectField, InputField, InterfaceField",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Bool("syntax") {
				fmt.Print(docs.SearchSyntax)
				return nil
			}

			jsonOutput := cmd.Bool("json")
			return handleError(runSearchCommand(cmd, jsonOutput), jsonOutput)
		},
	}
}

func runSearchCommand(cmd *cli.Command, jsonOutput bool) error {
	kindFilter := cmd.String("kind")
	hasKindFilter := kindFilter != ""

	if hasKindFilter {
		if cmd.Args().Len() > 1 {
			return fmt.Errorf("requires at most 1 argument when --kind is set: [query]")
		}
	} else {
		if cmd.Args().Len() != 1 {
			return fmt.Errorf("requires exactly 1 argument: <query>")
		}
	}

	query := cmd.Args().First()

	if hasKindFilter {
		var err error
		query, err = applyKindFilter(kindFilter, query)
		if err != nil {
			return err
		}
	}
	limit := cmd.Int("limit")
	noPager := cmd.Bool("no-pager")

	// Get schema (empty string for default when no flag specified)
	schemaArg := cmd.String("schema")

	// Resolve schema argument (path, ID, or default)
	schema, err := LoadSchema(schemaArg)
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

	// Create index if it doesn't exist
	if !indexer.Exists(schema.ID) {
		fmt.Printf("Indexing schema '%s'...\n", schema.ID)
		if err := indexer.Index(schema.ID, &schema.GQLSchema); err != nil {
			return fmt.Errorf("failed to index schema: %w", err)
		}
	}

	// Search
	searcher := search.NewSearcher(schemasDir)
	defer searcher.Close()

	results, err := searcher.Search(schema.ID, query, limit)
	if err != nil {
		return fmt.Errorf("search failed: %w (try using 'gqlxp library reindex %s')", err, schema.ID)
	}

	// Handle JSON output
	if jsonOutput {
		err := printSearchResultsJSON(results)
		for _, w := range unknownFieldWarnings(query) {
			fmt.Fprintln(os.Stderr, w)
		}
		return err
	}

	// Display results
	if len(results) == 0 {
		fmt.Printf("No results found for query: %q\n", query)
		for _, w := range unknownFieldWarnings(query) {
			fmt.Println(w)
		}
		return nil
	}

	var maxLimitInfo string
	if len(results) == limit {
		maxLimitInfo = fmt.Sprintf(" (increase search %s for more)", codeStyle.Render("--limit N"))
	}
	// Multiple results - show list and let user choose
	var output strings.Builder

	// Show command suggestions in header
	pathArg := headerStyle.Render("<object>.<field>")
	fmt.Fprintf(&output, "To display more info about a result, run: \n\t%s %s\n",
		codeStyle.Render("gqlxp show --schema "+schema.ID), pathArg)
	fmt.Fprintf(&output, "To open result in TUI app, run: \n\t%s --select %s\n\n",
		codeStyle.Render("gqlxp app --schema "+schema.ID), pathArg)

	fmt.Fprintf(&output, "Found %d results for %q%s:\n\n", len(results), query, maxLimitInfo)
	for i, result := range results {
		// Highlight the type in pink
		fmt.Fprintf(&output, "%d. %s %s\n", i+1, headerStyle.Render(result.Path), "("+result.Kind+")")
		if result.Description != "" {
			fmt.Fprintf(&output, "   %s\n", result.Description)
		}
	}

	// Use pager if content is long enough and not disabled
	rendered := output.String()
	if terminal.ShouldUsePager(rendered, noPager) {
		if err := terminal.ShowInPager(rendered); err != nil {
			return err
		}
		for _, w := range unknownFieldWarnings(query) {
			fmt.Println(w)
		}
		return nil
	}

	fmt.Print(rendered)
	for _, w := range unknownFieldWarnings(query) {
		fmt.Println(w)
	}
	return nil
}

// applyKindFilter validates kindFilter and prepends a +kind:<Canonical> clause to query.
func applyKindFilter(kindFilter, query string) (string, error) {
	if strings.Contains(query, "kind:") {
		return "", fmt.Errorf("cannot use --kind flag when query already contains a kind: filter")
	}
	canonical, ok := canonicalSearchKinds[strings.ToLower(kindFilter)]
	if !ok {
		return "", fmt.Errorf("invalid --kind value %q; valid values: Query, Mutation, Object, Input, Enum, Scalar, Interface, Union, Directive, ObjectField, InputField, InterfaceField", kindFilter)
	}
	kindClause := "+kind:" + canonical
	if query != "" {
		return kindClause + " " + query, nil
	}
	return kindClause, nil
}

// printSearchResultsJSON outputs search results as pretty-printed JSON
func printSearchResultsJSON(results []search.SearchResult) error {
	jsonBytes, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results to JSON: %w", err)
	}
	fmt.Println(string(jsonBytes))
	return nil
}
