package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/tonysyu/gqlxp/gql"
)

func validateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate [<operation-file>]",
		Short: "Validate a GraphQL operation against a schema",
		Long: `Validates a GraphQL operation against a schema.

Uses default schema when --schema is not specified.
Use 'gqlxp library default' to set the default schema.

Reads from a file argument if provided, or from stdin if omitted.
Exits with code 0 if valid, code 1 if there are errors.

JSON output format: {"valid": true|false, "errors": [{"line": N, "column": N, "message": "..."}]}`,
		Example: `  gqlxp validate examples/queries/github-user.graphql
  gqlxp validate -s github examples/queries/github-user.graphql
  gqlxp validate --json examples/queries/github-user.graphql
  cat examples/queries/github-user.graphql | gqlxp validate`,
		RunE: func(cmd *cobra.Command, args []string) error {
			schemaArg, _ := cmd.Flags().GetString("schema")
			jsonOutput, _ := cmd.Flags().GetBool("json")
			aiMode, _ := cmd.Flags().GetBool("ai")
			if aiMode {
				jsonOutput = true
				os.Setenv("NO_COLOR", "1")
			}
			var filePath string
			if len(args) > 0 {
				filePath = args[0]
			}
			return handleError(runValidateCommand(schemaArg, filePath, jsonOutput), jsonOutput)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Flags().Bool("json", false, "output results as JSON (recommended for AI/programmatic use)")
	cmd.Flags().Bool("ai", false, "AI/programmatic mode: JSON output, no pager, no color")

	return cmd
}

func runValidateCommand(schemaArg, filePath string, jsonOutput bool) error {
	schema, err := LoadSchema(schemaArg)
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
	return printJSON(result)
}

func printJSON(v any) error {
	jsonBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal to JSON: %w", err)
	}
	fmt.Println(string(jsonBytes))
	return nil
}

func printJSONError(err error) {
	type jsonErr struct {
		Error string `json:"error"`
	}
	_ = printJSON(jsonErr{Error: err.Error()})
}

func handleError(err error, jsonOutput bool) error {
	if err == nil {
		return nil
	}
	if jsonOutput {
		printJSONError(err)
		return nil
	}
	return err
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
