package cli

import (
	"context"
	"fmt"

	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/gqlfmt"
	"github.com/tonysyu/gqlxp/utils/terminal"
	"github.com/urfave/cli/v3"
)

// showCommand creates the show subcommand
func showCommand() *cli.Command {
	return &cli.Command{
		Name:      "show",
		Usage:     "Show a GraphQL type definition in the terminal",
		ArgsUsage: "<type-name>",
		Description: `Shows the full definition of a GraphQL type or field.

Uses default schema when --schema is not specified.
Use 'gqlxp library default' to set the default schema.

For AI/programmatic use, add --json --no-pager for machine-readable output.

Type-name formats:
  User                 Object, Input, Enum, Scalar, Interface, or Union type
  Query.getUser        A query field
  Mutation.createUser  A mutation field
  @auth                A directive

--include options (comma-separated):
  usages       List all types and fields that reference this type
  return-type  Show the full definition of the return type (Query/Mutation only)

Examples:
  gqlxp show User                                     # Uses default schema
  gqlxp show -s github User --json --no-pager         # JSON output for AI use
  gqlxp show -s github Query.getUser --include return-type
  gqlxp show -s examples/github.graphqls User         # Uses specific file
  gqlxp show @auth`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "schema",
				Aliases: []string{"s"},
				Usage:   "Schema file path or library ID to use",
			},
			&cli.BoolFlag{
				Name:  "no-pager",
				Usage: "disable pager; use for non-interactive/AI use",
			},
			&cli.BoolFlag{
				Name:  "json",
				Usage: "output as JSON (recommended for AI/programmatic use)",
			},
			&cli.StringFlag{
				Name:  "include",
				Usage: "comma-separated additional sections: usages (where this type is referenced), return-type (full return type definition; Query/Mutation only)",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			jsonOutput := cmd.Bool("json")
			if cmd.Args().Len() != 1 {
				return handleError(fmt.Errorf("requires exactly 1 argument: <type-name>"), jsonOutput)
			}

			typeName := cmd.Args().First()
			noPager := cmd.Bool("no-pager")
			include := cmd.String("include")

			// Get schema (empty string for default when no flag specified)
			schemaArg := cmd.String("schema")

			return handleError(printType(schemaArg, typeName, noPager, jsonOutput, include), jsonOutput)
		},
	}
}

func printType(schemaArg, typeName string, noPager bool, jsonOutput bool, include string) error {
	// Resolve schema argument (path, ID, or default)
	schema, err := resolveSchemaFromArgument(schemaArg)
	if err != nil {
		return err
	}

	// Parse schema
	parsedSchema, err := gql.ParseSchema(schema.Content)
	if err != nil {
		return fmt.Errorf("error parsing schema: %w", err)
	}

	// Parse include options
	opts := gqlfmt.ParseIncludeOptions(include)

	// Handle JSON output
	if jsonOutput {
		jsonStr, err := gqlfmt.GenerateJSON(parsedSchema, typeName, opts)
		if err != nil {
			return err
		}
		fmt.Println(jsonStr)
		return nil
	}

	// Generate markdown content based on type name
	markdown, err := gqlfmt.GenerateMarkdown(parsedSchema, typeName, opts)
	if err != nil {
		return err
	}

	// Render markdown using terminal renderer
	renderer, _ := terminal.NewMarkdownRenderer()
	rendered := terminal.RenderMarkdownOrPlain(renderer, markdown)

	// Use pager if content is long enough and not disabled
	if terminal.ShouldUsePager(rendered, noPager) {
		return terminal.ShowInPager(rendered)
	}

	fmt.Print(rendered)
	return nil
}
