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
	Open() (Panel, bool)
}

// Ensure that all item types implements ListItem interface
var _ ListItem = (*fieldItem)(nil)
var _ ListItem = (*typeDefItem)(nil)
var _ ListItem = (*argumentItem)(nil)
var _ ListItem = (*simpleItem)(nil)

func adaptFieldDefinitions(queryFields []*ast.FieldDefinition) []ListItem {
	adaptedItems := make([]ListItem, 0, len(queryFields))
	for _, f := range queryFields {
		adaptedItems = append(adaptedItems, newFieldDefItem(f))
	}
	return adaptedItems
}

func adaptObjectDefinitions(objects []*ast.ObjectDefinition) []ListItem {
	adaptedItems := make([]ListItem, 0, len(objects))
	for _, obj := range objects {
		adaptedItems = append(adaptedItems, newTypeDefItem(obj))
	}
	return adaptedItems
}

func adaptInputDefinitions(inputs []*ast.InputObjectDefinition) []ListItem {
	adaptedItems := make([]ListItem, 0, len(inputs))
	for _, input := range inputs {
		adaptedItems = append(adaptedItems, newTypeDefItem(input))
	}
	return adaptedItems
}

func adaptEnumDefinitions(enums []*ast.EnumDefinition) []ListItem {
	adaptedItems := make([]ListItem, 0, len(enums))
	for _, enum := range enums {
		adaptedItems = append(adaptedItems, newTypeDefItem(enum))
	}
	return adaptedItems
}

func adaptScalarDefinitions(scalars []*ast.ScalarDefinition) []ListItem {
	adaptedItems := make([]ListItem, 0, len(scalars))
	for _, scalar := range scalars {
		adaptedItems = append(adaptedItems, newTypeDefItem(scalar))
	}
	return adaptedItems
}

func adaptInterfaceDefinitions(interfaces []*ast.InterfaceDefinition) []ListItem {
	adaptedItems := make([]ListItem, 0, len(interfaces))
	for _, iface := range interfaces {
		adaptedItems = append(adaptedItems, newTypeDefItem(iface))
	}
	return adaptedItems
}

func adaptUnionDefinitions(unions []*ast.UnionDefinition) []ListItem {
	adaptedItems := make([]ListItem, 0, len(unions))
	for _, union := range unions {
		adaptedItems = append(adaptedItems, newTypeDefItem(union))
	}
	return adaptedItems
}

func adaptDirectiveDefinitions(directives []*ast.DirectiveDefinition) []ListItem {
	adaptedItems := make([]ListItem, 0, len(directives))
	for _, directive := range directives {
		adaptedItems = append(adaptedItems, newDirectiveDefinitionItem(directive))
	}
	return adaptedItems
}

func newFieldDefItem(gqlField *ast.FieldDefinition) ListItem {
	return fieldItem{
		gqlField: gqlField,
	}
}

// Adapter for DefaultItem interface used by charmbracelet/bubbles/list
// https://pkg.go.dev/github.com/charmbracelet/bubbles@v0.21.0/list#DefaultItem
type fieldItem struct {
	gqlField *ast.FieldDefinition
}

func (i fieldItem) Title() string {
	return i.gqlField.Name.Value
}

func (i fieldItem) FilterValue() string {
	return i.Title()
}

func (i fieldItem) Description() string {
	if desc := i.gqlField.GetDescription(); desc != nil {
		return desc.Value
	}
	return ""
}

// Implement ListItem interface
func (i fieldItem) Open() (Panel, bool) {
	// Create list items for the detail view
	var detailItems []ListItem

	// Add description as a header if available
	if desc := i.Description(); desc != "" {
		detailItems = append(detailItems, simpleItem{title: desc})
	}

	inputValueItems := adaptInputValueDefinitions(i.gqlField.Arguments)
	if len(inputValueItems) > 0 {
		detailItems = append(detailItems, newSectionHeader("Input Arguments"))
		detailItems = append(detailItems, inputValueItems...)
	}

	// Add result type section
	detailItems = append(detailItems, newSectionHeader("Result Type"))
	detailItems = append(detailItems, newTypeItem(i.gqlField.Type))
	return newListPanel(detailItems, i.Title()), true
}

// Create an array of ListItem instances given InputValueDefinition. This is used for
// `ast.FieldDefinition.Arguments` and `ast.InputObjectDefinition.Fields`
func adaptInputValueDefinitions(inputValues []*ast.InputValueDefinition) []ListItem {
	// Add arguments section if any
	var items []ListItem
	if len(inputValues) > 0 {
		for _, arg := range inputValues {
			argDesc := ""
			if arg.Description != nil && arg.Description.Value != "" {
				argDesc = arg.Description.Value
			}
			// TODO: Update argumentItem to support proper Open and use custom display string
			items = append(items, argumentItem{
				name:        arg.Name.Value,
				argType:     gql.GetTypeString(arg.Type),
				description: argDesc,
			})
		}
	}
	return items

}

func newTypeDefItem(typeDef ast.TypeDefinition) typeDefItem {
	return typeDefItem{
		typeDef: typeDef,
	}
}

// Adapter for DefaultItem interface used by charmbracelet/bubbles/list
// https://pkg.go.dev/github.com/charmbracelet/bubbles@v0.21.0/list#DefaultItem
type typeDefItem struct {
	typeDef ast.TypeDefinition
}

func (i typeDefItem) Title() string {
	switch typeDef := (i.typeDef).(type) {
	case *ast.ScalarDefinition:
		return typeDef.Name.Value
	case *ast.ObjectDefinition:
		return typeDef.Name.Value
	case *ast.InterfaceDefinition:
		return typeDef.Name.Value
	case *ast.UnionDefinition:
		return typeDef.Name.Value
	case *ast.EnumDefinition:
		return typeDef.Name.Value
	case *ast.InputObjectDefinition:
		return typeDef.Name.Value
	}
	return "UNKNOWN"
}

func (i typeDefItem) FilterValue() string {
	return i.Title()
}

func (i typeDefItem) Description() string {
	if desc := (i.typeDef).GetDescription(); desc != nil {
		return desc.Value
	}
	return ""
}

// Implement ListItem interface
func (i typeDefItem) Open() (Panel, bool) {
	// Create list items for the detail view
	var detailItems []ListItem

	// Add description as a header if available
	if desc := i.Description(); desc != "" {
		detailItems = append(detailItems, simpleItem{title: desc})
	}

	// TODO: Update definitions with details
	switch typeDef := (i.typeDef).(type) {
	case *ast.ObjectDefinition:
		detailItems = append(detailItems, adaptFieldDefinitions(typeDef.Fields)...)
	case *ast.ScalarDefinition:
	case *ast.InterfaceDefinition:
		detailItems = append(detailItems, adaptFieldDefinitions(typeDef.Fields)...)
	case *ast.UnionDefinition:
		// TODO: This probably requires a reference to the schema
	case *ast.EnumDefinition:
		// TODO: This probably requires a reference to the schema
	case *ast.InputObjectDefinition:
		detailItems = append(detailItems, adaptInputValueDefinitions(typeDef.Fields)...)
	}

	return newListPanel(detailItems, i.Title()), true
}

// simpleItem displays title and description and has a no-op Open() function
type simpleItem struct {
	title, description string
}

func (di simpleItem) Title() string       { return di.title }
func (di simpleItem) Description() string { return di.description }
func (di simpleItem) FilterValue() string { return di.title }
func (di simpleItem) Open() (Panel, bool) { return nil, false }

func newSectionHeader(title string) simpleItem {
	return simpleItem{title: "======== " + title + " ========"}
}

// argumentItem displays argument information
type argumentItem struct {
	name        string
	argType     string
	description string
}

func (ai argumentItem) Title() string {
	return fmt.Sprintf("%s: %s", ai.name, ai.argType)
}
func (ai argumentItem) Description() string {
	if ai.description != "" {
		return ai.description
	}
	return ""
}
func (ai argumentItem) FilterValue() string { return ai.name }
func (ai argumentItem) Open() (Panel, bool) { return nil, false }

func newTypeItem(t ast.Type) simpleItem {
	// TODO: This probably requires a reference to the schema to return full type when opening
	return simpleItem{title: gql.GetTypeString(t)}
}

func newDirectiveDefinitionItem(directive *ast.DirectiveDefinition) simpleItem {
	return simpleItem{
		title:       directive.Name.Value,
		description: directive.Description.Value,
	}
}
