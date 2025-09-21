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

func adaptNamedItems(namedNodes []*ast.Named) []ListItem {
	adaptedItems := make([]ListItem, 0, len(namedNodes))
	for _, node := range namedNodes {
		adaptedItems = append(adaptedItems, newNamedItem(node))
	}
	return adaptedItems
}

func adaptEnumValueDefinitions(enumNodes []*ast.EnumValueDefinition) []ListItem {
	adaptedItems := make([]ListItem, 0, len(enumNodes))
	for _, node := range enumNodes {
		adaptedItems = append(adaptedItems, simpleItem{
			title:       node.Name.Value,
			description: node.Description.Value,
		})
	}
	return adaptedItems
}

// Adapter/delegate for ast.FieldDefinition to support ListItem interface
type fieldItem struct {
	gqlField *ast.FieldDefinition
}

func newFieldDefItem(gqlField *ast.FieldDefinition) ListItem {
	return fieldItem{
		gqlField: gqlField,
	}
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
	var items []ListItem
	if len(inputValues) > 0 {
		for _, arg := range inputValues {
			items = append(items, newInputValueItem(arg))
		}
	}
	return items

}

// Adapter/delegate for ast.TypeDefinition to support ListItem interface
type typeDefItem struct {
	typeDef ast.TypeDefinition
}

func newTypeDefItem(typeDef ast.TypeDefinition) typeDefItem {
	return typeDefItem{
		typeDef: typeDef,
	}
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

	switch typeDef := (i.typeDef).(type) {
	case *ast.ObjectDefinition:
		detailItems = append(detailItems, adaptFieldDefinitions(typeDef.Fields)...)
	case *ast.ScalarDefinition:
	case *ast.InterfaceDefinition:
		detailItems = append(detailItems, adaptFieldDefinitions(typeDef.Fields)...)
	case *ast.UnionDefinition:
		detailItems = append(detailItems, adaptNamedItems(typeDef.Types)...)
	case *ast.EnumDefinition:
		detailItems = append(detailItems, adaptEnumValueDefinitions(typeDef.Values)...)
	case *ast.InputObjectDefinition:
		detailItems = append(detailItems, adaptInputValueDefinitions(typeDef.Fields)...)
	}

	return newListPanel(detailItems, i.Title()), true
}

// ListItem interface with arbitrary title and description and no-op Open() function.
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

func newNamedItem(node *ast.Named) simpleItem {
	// TODO: This probably requires a reference to the schema to return full type when opening
	return simpleItem{title: node.Name.Value}
}

func newTypeItem(t ast.Type) simpleItem {
	// TODO: This probably requires a reference to the schema to return full type when opening
	return simpleItem{title: gql.GetTypeString(t)}
}

func newInputValueItem(inputValue *ast.InputValueDefinition) simpleItem {
	// TODO: Update item to support proper Open and use custom display string
	fieldName := inputValue.Name.Value
	fieldType := gql.GetTypeString(inputValue.Type)
	description := ""
	if inputValue.Description != nil {
		description = inputValue.Description.Value
	}
	return simpleItem{
		title:       fmt.Sprintf("%s: %s", fieldName, fieldType),
		description: description,
	}
}

func newDirectiveDefinitionItem(directive *ast.DirectiveDefinition) simpleItem {
	return simpleItem{
		title:       directive.Name.Value,
		description: directive.Description.Value,
	}
}
