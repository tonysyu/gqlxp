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

// Adapter for DefaultItem interface used by charmbracelet/bubbles/list
// https://pkg.go.dev/github.com/charmbracelet/bubbles@v0.21.0/list#DefaultItem
type item struct {
	gqlField *ast.FieldDefinition
}

func (i item) Title() string {
	return i.gqlField.Name.Value
}

func (i item) FilterValue() string {
	return i.Title()
}

func (i item) Description() string {
	return i.gqlField.GetDescription().Value
}

// Implement InteractiveListItem interface
func (i item) Open() Panel {
	return newStringPanel(i.Description())
}
