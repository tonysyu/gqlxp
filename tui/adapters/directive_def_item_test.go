package adapters

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/gql"
)

func TestDirectiveDefinitionItemCreation(t *testing.T) {
	is := is.New(t)

	schemaString := `
		directive @deprecated(reason: String = "No longer supported") on FIELD_DEFINITION | ENUM_VALUE
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))
	resolver := gql.NewSchemaResolver(&schema)
	directive := schema.Directive["deprecated"]

	item := newDirectiveDefItem(directive, resolver)
	is.Equal(item.Title(), "@deprecated(reason: String = \"No longer supported\")")
	is.Equal(item.Description(), "")

	// Directive items are now openable and show their arguments
	panel, ok := item.OpenPanel()
	is.True(ok)
	is.True(panel != nil)
}
