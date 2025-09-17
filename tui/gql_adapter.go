package tui

import (
	"github.com/graphql-go/graphql/language/ast"
)

func AdaptGraphQLItems(queryFields map[string]*ast.FieldDefinition) []item {
	adaptedItems := make([]item, 0, len(queryFields))
	for _, f := range queryFields {
		adaptedItems = append(adaptedItems, newItem(f))
	}

	return adaptedItems
}

func newItem(gqlField *ast.FieldDefinition) item {
	return item{
		gqlField: gqlField,
	}
}

type item struct {
	gqlField *ast.FieldDefinition
}

func (i item) Title() string {
	return i.gqlField.Name.Value
}
