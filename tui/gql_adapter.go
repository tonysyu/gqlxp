package tui

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/tonysyu/gq/gql"
)

func adaptGraphQLItems(queryFields map[string]*ast.FieldDefinition) []item {
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
	if desc := i.gqlField.GetDescription(); desc != nil {
		return desc.Value
	}
	return ""
}

// Implement ListItem interface
func (i item) Open() Panel {
	var content strings.Builder

	// Add description if available
	if desc := i.Description(); desc != "" {
		content.WriteString(desc)
		content.WriteString("\n\n")
	}

	// Add input arguments section
	if len(i.gqlField.Arguments) > 0 {
		content.WriteString("Input Arguments:\n")
		for _, arg := range i.gqlField.Arguments {
			argName := arg.Name.Value
			argType := gql.GetTypeString(arg.Type)
			content.WriteString(fmt.Sprintf("  • %s: %s", argName, argType))

			// Add argument description if available
			if arg.Description != nil && arg.Description.Value != "" {
				content.WriteString(fmt.Sprintf(" - %s", arg.Description.Value))
			}
			content.WriteString("\n")
		}
		content.WriteString("\n")
	}

	// Add result type section
	content.WriteString("Result Type:\n")
	resultType := gql.GetTypeString(i.gqlField.Type)
	content.WriteString(fmt.Sprintf("  %s", resultType))

	return newStringPanel(content.String())
}
