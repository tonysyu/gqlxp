package cli

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/tonysyu/gqlxp/gql"
	"github.com/urfave/cli/v3"
)

func parseCommand() *cli.Command {
	return &cli.Command{
		Name:      "parse",
		Usage:     "Validate a GraphQL operation against a schema",
		ArgsUsage: "[<operation-file>]",
		Description: `Validates a GraphQL operation against a schema.

Uses default schema when --schema is not specified.
Use 'gqlxp library default' to set the default schema.

Reads from a file argument if provided, or from stdin if omitted.
Exits with code 0 if valid, code 1 if there are errors.

Examples:
  gqlxp parse query.graphql
  gqlxp parse -s my-schema query.graphql
  cat query.graphql | gqlxp parse`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "schema",
				Aliases: []string{"s"},
				Usage:   "Schema file path or library ID to use",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			schemaArg := cmd.String("schema")
			var filePath string
			if cmd.Args().Len() > 0 {
				filePath = cmd.Args().First()
			}
			return runParseCommand(schemaArg, filePath)
		},
	}
}

func runParseCommand(schemaArg, filePath string) error {
	schema, err := resolveSchemaFromArgument(schemaArg)
	if err != nil {
		return err
	}

	var operationContent string
	var sourceName string
	if filePath != "" {
		data, readErr := os.ReadFile(filePath)
		if readErr != nil {
			return fmt.Errorf("error reading file: %w", readErr)
		}
		operationContent = string(data)
		sourceName = filePath
	} else {
		data, readErr := io.ReadAll(os.Stdin)
		if readErr != nil {
			return fmt.Errorf("error reading stdin: %w", readErr)
		}
		operationContent = string(data)
		sourceName = "<stdin>"
	}

	errorLines := validateOperation(schema.Content, operationContent, sourceName)
	for _, line := range errorLines {
		fmt.Println(line)
	}
	if len(errorLines) > 0 {
		os.Exit(1)
	}
	return nil
}

func validateOperation(schemaContent []byte, operationContent, sourceName string) []string {
	errs, err := gql.ValidateOperation(schemaContent, operationContent)
	if err != nil {
		return []string{fmt.Sprintf("%s: %s", sourceName, err.Error())}
	}

	lines := make([]string, 0, len(errs))
	for _, ve := range errs {
		if ve.Line > 0 {
			lines = append(lines, fmt.Sprintf("%s:%d:%d: %s", sourceName, ve.Line, ve.Column, ve.Message))
		} else {
			lines = append(lines, fmt.Sprintf("%s: %s", sourceName, ve.Message))
		}
	}
	return lines
}
