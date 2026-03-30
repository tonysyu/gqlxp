package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/tonysyu/gqlxp/gql"
	"github.com/urfave/cli/v3"
)

func validateCommand() *cli.Command {
	return &cli.Command{
		Name:      "validate",
		Usage:     "Validate a GraphQL operation against a schema",
		ArgsUsage: "[<operation-file>]",
		Description: `Validates a GraphQL operation against a schema.

Uses default schema when --schema is not specified.
Use 'gqlxp library default' to set the default schema.

Reads from a file argument if provided, or from stdin if omitted.
Exits with code 0 if valid, code 1 if there are errors.

JSON output format: {"valid": true|false, "errors": [{"line": N, "column": N, "message": "..."}]}

Examples:
  gqlxp validate examples/queries/github-user.graphql
  gqlxp validate -s github examples/queries/github-user.graphql
  gqlxp validate --json examples/queries/github-user.graphql
  cat examples/queries/github-user.graphql | gqlxp validate`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "schema",
				Aliases: []string{"s"},
				Usage:   "Schema file path or library ID to use",
			},
			&cli.BoolFlag{
				Name:  "json",
				Usage: "output results as JSON (recommended for AI/programmatic use)",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			schemaArg := cmd.String("schema")
			jsonOutput := cmd.Bool("json")
			var filePath string
			if cmd.Args().Len() > 0 {
				filePath = cmd.Args().First()
			}
			return runValidateCommand(schemaArg, filePath, jsonOutput)
		},
	}
}

func runValidateCommand(schemaArg, filePath string, jsonOutput bool) error {
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

	if jsonOutput {
		return printValidationResultJSON(schema.Content, operationContent, sourceName)
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

type validationResult struct {
	Valid  bool            `json:"valid"`
	Errors []validationErr `json:"errors,omitempty"`
}

type validationErr struct {
	Line    int    `json:"line,omitempty"`
	Column  int    `json:"column,omitempty"`
	Message string `json:"message"`
}

func buildValidationResult(schemaContent []byte, operationContent, sourceName string) validationResult {
	errs, err := gql.ValidateOperation(schemaContent, operationContent)
	if err != nil {
		return validationResult{
			Valid:  false,
			Errors: []validationErr{{Message: fmt.Sprintf("%s: %s", sourceName, err.Error())}},
		}
	}

	result := validationResult{Valid: len(errs) == 0}
	for _, ve := range errs {
		result.Errors = append(result.Errors, validationErr{
			Line:    ve.Line,
			Column:  ve.Column,
			Message: ve.Message,
		})
	}
	return result
}

func printValidationResultJSON(schemaContent []byte, operationContent, sourceName string) error {
	result := buildValidationResult(schemaContent, operationContent, sourceName)
	if err := printJSON(result); err != nil {
		return err
	}
	if !result.Valid {
		os.Exit(1)
	}
	return nil
}

func printJSON(v any) error {
	jsonBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal to JSON: %w", err)
	}
	fmt.Println(string(jsonBytes))
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
