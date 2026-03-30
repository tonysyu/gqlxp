package cli

import (
	"context"
	"fmt"

	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/gqlfmt"
	"github.com/urfave/cli/v3"
)

func generateCommand() *cli.Command {
	return &cli.Command{
		Name:      "generate",
		Usage:     "Generate a skeleton GraphQL operation",
		ArgsUsage: "<Query.fieldName|Mutation.fieldName>",
		Description: `Scaffolds a skeleton GraphQL operation for a Query or Mutation field.
Output is a complete GraphQL operation document printed to stdout.

Uses default schema when --schema is not specified.
Use 'gqlxp library default' to set the default schema.

--depth controls how many levels of nested object fields are expanded (default: 1).
Use 'gqlxp show <type>' to inspect type definitions before generating.

Examples:
  gqlxp generate Query.getUser
  gqlxp generate -s examples/github.graphqls Query.repository
  gqlxp generate --depth 2 Query.getUser      # Expand nested fields 2 levels deep
  gqlxp generate Mutation.createUser`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "schema",
				Aliases: []string{"s"},
				Usage:   "Schema file path or library ID to use",
			},
			&cli.IntFlag{
				Name:  "depth",
				Usage: "levels of nested object fields to expand (default: 1)",
				Value: 1,
			},
			&cli.BoolFlag{
				Name:  "include-deprecated",
				Usage: "include deprecated fields in the generated operation",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 1 {
				return fmt.Errorf("requires exactly 1 argument: <Query.fieldName|Mutation.fieldName>")
			}

			fieldPath := cmd.Args().First()
			schemaArg := cmd.String("schema")
			depth := int(cmd.Int("depth"))
			includeDeprecated := cmd.Bool("include-deprecated")

			schema, err := resolveSchemaFromArgument(schemaArg)
			if err != nil {
				return err
			}

			parsedSchema, err := gql.ParseSchema(schema.Content)
			if err != nil {
				return fmt.Errorf("error parsing schema: %w", err)
			}

			opts := gqlfmt.GenerateOptions{
				Depth:             depth,
				IncludeDeprecated: includeDeprecated,
			}

			operation, err := gqlfmt.GenerateOperation(parsedSchema, fieldPath, opts)
			if err != nil {
				return err
			}

			fmt.Println(operation)
			return nil
		},
	}
}
