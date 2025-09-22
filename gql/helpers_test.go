package gql_test

import (
	"testing"

	"github.com/matryer/is"
	. "github.com/tonysyu/gq/gql"
)

func TestGetTypeString(t *testing.T) {
	is := is.New(t)

	schemaString := `
		type Query {
		  getString: String
		  getRequiredString: String!
		  getStringList: [String]
		  getRequiredStringList: [String]!
		  getListOfRequiredStrings: [String!]
		  getRequiredListOfRequiredStrings: [String!]!
		}
	`

	schema := ParseSchema([]byte(schemaString))

	testCases := []struct {
		fieldName    string
		expectedType string
	}{
		{"getString", "String"},
		{"getRequiredString", "String!"},
		{"getStringList", "[String]"},
		{"getRequiredStringList", "[String]!"},
		{"getListOfRequiredStrings", "[String!]"},
		{"getRequiredListOfRequiredStrings", "[String!]!"},
	}

	for _, tc := range testCases {
		t.Run(tc.fieldName, func(t *testing.T) {
			field, ok := schema.Query[tc.fieldName]
			is.True(ok)
			is.Equal(GetTypeString(field.Type), tc.expectedType)
		})
	}
}
