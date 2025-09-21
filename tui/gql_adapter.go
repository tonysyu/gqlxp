package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/tonysyu/gq/gql"
)

// list item that can be "opened" to provide additional information about the item.
// The opened data is represented as a Panel instance that can be rendered to users.
type ListItem interface {
	list.DefaultItem

	// Open Panel to show additional information.
	Open() Panel
}

// Ensure that all item types implements ListItem interface
var _ ListItem = (*item)(nil)

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
	// Create list items for the detail view
	var detailItems []list.Item

	// Add description as a header if available
	if desc := i.Description(); desc != "" {
		detailItems = append(detailItems, descriptionItem{content: desc})
	}

	// Add arguments section if any
	if len(i.gqlField.Arguments) > 0 {
		detailItems = append(detailItems, sectionHeaderItem{title: "Input Arguments"})
		for _, arg := range i.gqlField.Arguments {
			argDesc := ""
			if arg.Description != nil && arg.Description.Value != "" {
				argDesc = arg.Description.Value
			}
			detailItems = append(detailItems, argumentItem{
				name:        arg.Name.Value,
				argType:     gql.GetTypeString(arg.Type),
				description: argDesc,
			})
		}
	}

	// Add result type section
	detailItems = append(detailItems, sectionHeaderItem{title: "Result Type"})
	detailItems = append(detailItems, resultTypeItem{
		typeName: gql.GetTypeString(i.gqlField.Type),
	})

	return newListPanel(detailItems, i.Title()+" Details")
}

// descriptionItem displays field description
type descriptionItem struct {
	content string
}

func (di descriptionItem) Title() string       { return di.content }
func (di descriptionItem) Description() string { return "" }
func (di descriptionItem) FilterValue() string { return di.content }

// sectionHeaderItem displays section headers
type sectionHeaderItem struct {
	title string
}

func (shi sectionHeaderItem) Title() string       { return "======== " + shi.title + " ========" }
func (shi sectionHeaderItem) Description() string { return "" }
func (shi sectionHeaderItem) FilterValue() string { return "" }

// argumentItem displays argument information
type argumentItem struct {
	name        string
	argType     string
	description string
}

func (ai argumentItem) Title() string {
	return fmt.Sprintf("• %s: %s", ai.name, ai.argType)
}

func (ai argumentItem) Description() string {
	if ai.description != "" {
		return ai.description
	}
	return "No description available"
}

func (ai argumentItem) FilterValue() string { return ai.name }

// resultTypeItem displays result type information
type resultTypeItem struct {
	typeName string
}

func (rti resultTypeItem) Title() string       { return rti.typeName }
func (rti resultTypeItem) Description() string { return "Return type for this field" }
func (rti resultTypeItem) FilterValue() string { return rti.typeName }
