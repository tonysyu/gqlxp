package cli

import (
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
	is.True(strings.HasPrefix(errorLines[0], "query.graphql:"))                    // error identifies location
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
