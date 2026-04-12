package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tonysyu/gqlxp/gqlfmt"
)

func generateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate <Query.fieldName|Mutation.fieldName>",
		Short: "Generate a skeleton GraphQL operation",
		Args:  cobra.ExactArgs(1),
		Long: `Scaffolds a skeleton GraphQL operation for a Query or Mutation field.
Output is a complete GraphQL operation document printed to stdout.

Uses default schema when --schema is not specified.
Use 'gqlxp library default' to set the default schema.

--depth controls how many levels of nested object fields are expanded (default: 1).
Use 'gqlxp show <type>' to inspect type definitions before generating.`,
		Example: `  gqlxp generate Query.getUser
  gqlxp generate -s examples/github.graphqls Query.repository
  gqlxp generate --depth 2 Query.getUser      # Expand nested fields 2 levels deep
  gqlxp generate Mutation.createUser`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fieldPath := args[0]
			schemaArg, _ := cmd.Flags().GetString("schema")
			depth, _ := cmd.Flags().GetInt("depth")
			includeDeprecated, _ := cmd.Flags().GetBool("include-deprecated")

			schema, err := LoadSchema(schemaArg)
			if err != nil {
				return err
			}

			opts := gqlfmt.GenerateOptions{
				Depth:             depth,
				IncludeDeprecated: includeDeprecated,
			}

			operation, err := gqlfmt.GenerateOperation(schema.GQLSchema, fieldPath, opts)
			if err != nil {
				return err
			}

			fmt.Println(operation)
			return nil
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Flags().Int("depth", 1, "levels of nested object fields to expand")
	cmd.Flags().Bool("include-deprecated", false, "include deprecated fields in the generated operation")

	return cmd
}
