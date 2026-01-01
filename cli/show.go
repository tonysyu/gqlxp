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
		Description: `Shows the details of a GraphQL type directly to the terminal.

Uses default schema when --schema is not specified.
Use 'gqlxp library default' to set the default schema.

The type-name can be:
- A Query field name (prefix with "Query.")
- A Mutation field name (prefix with "Mutation.")
- A type name (Object, Input, Enum, Scalar, Interface, Union)
- A directive name (prefix with "@")

Examples:
  gqlxp show User                        # Uses default schema
  gqlxp show -s examples/github.graphqls User # Uses specific file
  gqlxp show -s github-api User          # Uses library ID
  gqlxp show -s github-api Query.getUser
  gqlxp show Mutation.createUser
  gqlxp show @auth`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "schema",
				Aliases: []string{"s"},
				Usage:   "Schema file path or library ID to use",
			},
			&cli.BoolFlag{
				Name:  "no-pager",
				Usage: "disable pager and show directly to stdout",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 1 {
				return fmt.Errorf("requires exactly 1 argument: <type-name>")
			}

			typeName := cmd.Args().First()
			noPager := cmd.Bool("no-pager")

			// Get schema (empty string for default when no flag specified)
			schemaArg := cmd.String("schema")

			return printType(schemaArg, typeName, noPager)
		},
	}
}

func printType(schemaArg, typeName string, noPager bool) error {
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

	// Generate markdown content based on type name
	markdown, err := gqlfmt.GenerateMarkdown(parsedSchema, typeName)
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
