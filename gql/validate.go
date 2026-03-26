package gql

import (
	"errors"
	"fmt"

	gqlparser "github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vektah/gqlparser/v2/parser"
	"github.com/vektah/gqlparser/v2/validator"
)

// ValidationError represents a single validation error with location info.
type ValidationError struct {
	Line    int
	Column  int
	Message string
}

// ValidateOperation validates a GraphQL operation document against a schema.
// Returns validation errors (empty slice if valid) and any fatal error loading the schema.
func ValidateOperation(schemaContent []byte, operationContent string) ([]ValidationError, error) {
	astSchema, err := gqlparser.LoadSchema(&ast.Source{
		Name:  "schema.graphql",
		Input: string(schemaContent),
	})
	if err != nil {
		return nil, fmt.Errorf("error loading schema: %w", err)
	}

	doc, parseErr := parser.ParseQuery(&ast.Source{Input: operationContent})
	if parseErr != nil {
		var gqlErr *gqlerror.Error
		if errors.As(parseErr, &gqlErr) {
			return toValidationErrors(gqlerror.List{gqlErr}), nil
		}
		return []ValidationError{{Message: parseErr.Error()}}, nil
	}

	return toValidationErrors(validator.Validate(astSchema, doc)), nil
}

func toValidationErrors(errs gqlerror.List) []ValidationError {
	result := make([]ValidationError, 0, len(errs))
	for _, err := range errs {
		ve := ValidationError{Message: err.Message}
		if len(err.Locations) > 0 {
			ve.Line = err.Locations[0].Line
			ve.Column = err.Locations[0].Column
		}
		result = append(result, ve)
	}
	return result
}
