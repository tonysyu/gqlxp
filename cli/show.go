package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tonysyu/gqlxp/gqlfmt"
	"github.com/tonysyu/gqlxp/utils/terminal"
)

// showCommand creates the show subcommand
func showCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <type-name>",
		Short: "Show a GraphQL type definition in the terminal",
		Args:  cobra.ExactArgs(1),
		Long: `Shows the full definition of a GraphQL type or field.

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
  return-type  Show the full definition of the return type (Query/Mutation only)`,
		Example: `  gqlxp show User                                     # Uses default schema
  gqlxp show -s github User --json --no-pager         # JSON output for AI use
  gqlxp show -s github Query.getUser --include return-type
  gqlxp show -s examples/github.graphqls User         # Uses specific file
  gqlxp show @auth`,
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput, _ := cmd.Flags().GetBool("json")
			typeName := args[0]
			noPager, _ := cmd.Flags().GetBool("no-pager")
			include, _ := cmd.Flags().GetString("include")
			schemaArg, _ := cmd.Flags().GetString("schema")

			return handleError(printType(schemaArg, typeName, noPager, jsonOutput, include), jsonOutput)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Flags().Bool("no-pager", false, "disable pager; use for non-interactive/AI use")
	cmd.Flags().Bool("json", false, "output as JSON (recommended for AI/programmatic use)")
	cmd.Flags().String("include", "", "comma-separated additional sections: usages (where this type is referenced), return-type (full return type definition; Query/Mutation only)")

	return cmd
}

func printType(schemaArg, typeName string, noPager bool, jsonOutput bool, include string) error {
	schema, err := LoadSchema(schemaArg)
	if err != nil {
		return err
	}

	// Parse include options
	opts := gqlfmt.ParseIncludeOptions(include)

	// Handle JSON output
	if jsonOutput {
		jsonStr, err := gqlfmt.GenerateJSON(schema.GQLSchema, typeName, opts)
		if err != nil {
			return err
		}
		fmt.Println(jsonStr)
		return nil
	}

	// Generate markdown content based on type name
	markdown, err := gqlfmt.GenerateMarkdown(schema.GQLSchema, typeName, opts)
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
