package adapters

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/gql"
)

func TestArgumentListCreation(t *testing.T) {
	is := is.New(t)

	schemaString := `
		type Query {
		  testField(arg1: String!, arg2: Int, arg3: [String]): String
		}
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))
	resolver := gql.NewSchemaResolver(&schema)
	field := schema.Query["testField"]

	// Test argument items creation
	items := adaptArguments(field.Arguments(), resolver)
	is.Equal(len(items), 3)

	// Test first argument
	item1 := items[0]
	is.Equal(item1.Title(), "arg1: String!")

	// Test second argument
	item2 := items[1]
	is.Equal(item2.Title(), "arg2: Int")

	// Test third argument
	item3 := items[2]
	is.Equal(item3.Title(), "arg3: [String]")
}
