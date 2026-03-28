package cli

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/matryer/is"
)

const parseTestSchema = `
type Query {
	user(id: ID!): User
	users: [User!]!
}
type User {
	id: ID!
	name: String!
}
`

func TestValidateOperation_ValidOperation(t *testing.T) {
	is := is.New(t)

	operation := `query { user(id: "1") { id name } }`
	errorLines := validateOperation([]byte(parseTestSchema), operation, "query.graphql")

	is.Equal(len(errorLines), 0) // valid operation: no errors, exits with code 0
}

func TestValidateOperation_SyntaxError(t *testing.T) {
	is := is.New(t)

	operation := `query { user(id: "1") { id name ` // missing closing braces
	errorLines := validateOperation([]byte(parseTestSchema), operation, "query.graphql")

	is.True(len(errorLines) > 0)                                // syntax error should produce errors, exits with code 1
	is.True(strings.HasPrefix(errorLines[0], "query.graphql:")) // format: source:line:col: message
	parts := strings.SplitN(errorLines[0], ":", 4)
	is.True(len(parts) >= 4) // source:line:col: message
}

func TestValidateOperation_UnknownField(t *testing.T) {
	is := is.New(t)

	operation := `query { user(id: "1") { id nonExistentField } }`
	errorLines := validateOperation([]byte(parseTestSchema), operation, "query.graphql")

	is.True(len(errorLines) > 0)                                                  // unknown field produces error, exits with code 1
	is.True(strings.HasPrefix(errorLines[0], "query.graphql:"))                   // error identifies location
	is.True(strings.Contains(strings.Join(errorLines, "\n"), "nonExistentField")) // error identifies the unknown field
}

func TestValidateOperation_MultipleErrors(t *testing.T) {
	is := is.New(t)

	operation := `query { user(id: "1") { id badField1 badField2 } }`
	errorLines := validateOperation([]byte(parseTestSchema), operation, "query.graphql")

	is.True(len(errorLines) > 1) // multiple errors: all reported, one per line
}

func TestValidateOperation_StdinSource(t *testing.T) {
	is := is.New(t)

	operation := `query { user(id: "1") { nonExistentField } }`
	errorLines := validateOperation([]byte(parseTestSchema), operation, "<stdin>")

	is.True(len(errorLines) > 0)
	is.True(strings.HasPrefix(errorLines[0], "<stdin>:")) // stdin input uses <stdin> as source in error messages
}

func TestValidationResultJSON_ValidOperation(t *testing.T) {
	is := is.New(t)

	result := buildValidationResult([]byte(parseTestSchema), `query { user(id: "1") { id name } }`, "query.graphql")

	is.True(result.Valid)
	is.Equal(len(result.Errors), 0)

	// Ensure it marshals to valid JSON
	jsonBytes, err := json.Marshal(result)
	is.NoErr(err)
	is.True(strings.Contains(string(jsonBytes), `"valid":true`))
	is.True(!strings.Contains(string(jsonBytes), `"errors"`)) // omitempty: no errors field when valid
}

func TestValidationResultJSON_InvalidOperation(t *testing.T) {
	is := is.New(t)

	result := buildValidationResult([]byte(parseTestSchema), `query { user(id: "1") { badField } }`, "query.graphql")

	is.True(!result.Valid)
	is.True(len(result.Errors) > 0)
	is.True(result.Errors[0].Line > 0)
	is.True(result.Errors[0].Column > 0)
	is.True(strings.Contains(result.Errors[0].Message, "badField"))
}

func TestValidationResultJSON_ParseError(t *testing.T) {
	is := is.New(t)

	result := buildValidationResult([]byte(parseTestSchema), `query { user(id: "1") { id `, "query.graphql")

	is.True(!result.Valid)
	is.True(len(result.Errors) > 0)
}
