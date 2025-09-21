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
var _ ListItem = (*directiveDefinitionItem)(nil)
var _ ListItem = (*resultTypeItem)(nil)
var _ ListItem = (*sectionHeaderItem)(nil)

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
		adaptedItems = append(adaptedItems, directiveDefinitionItem{directive})
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
		detailItems = append(detailItems, descriptionItem{content: desc})
	}

	inputValueItems := adaptInputValueDefinitions(i.gqlField.Arguments)
	if len(inputValueItems) > 0 {
		detailItems = append(detailItems, sectionHeaderItem{title: "Input Arguments"})
		detailItems = append(detailItems, inputValueItems...)
	}

	// Add result type section
	detailItems = append(detailItems, sectionHeaderItem{title: "Result Type"})
	detailItems = append(detailItems, resultTypeItem{
		typeName: gql.GetTypeString(i.gqlField.Type),
	})

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
		detailItems = append(detailItems, descriptionItem{content: desc})
	}

	// TODO: Update definitions with details
	switch typeDef := (i.typeDef).(type) {
	case *ast.ObjectDefinition:
		detailItems = append(detailItems, adaptFieldDefinitions(typeDef.Fields)...)
	case *ast.ScalarDefinition:
	case *ast.InterfaceDefinition:
	case *ast.UnionDefinition:
	case *ast.EnumDefinition:
		// pass
	case *ast.InputObjectDefinition:
		detailItems = append(detailItems, adaptInputValueDefinitions(typeDef.Fields)...)
	}

	return newListPanel(detailItems, i.Title()), true
}

// descriptionItem displays field description
type descriptionItem struct {
	content string
}

func (di descriptionItem) Title() string       { return di.content }
func (di descriptionItem) Description() string { return "" }
func (di descriptionItem) FilterValue() string { return di.content }
func (di descriptionItem) Open() (Panel, bool) { return nil, false }

// sectionHeaderItem displays section headers
type sectionHeaderItem struct {
	title string
}

func (shi sectionHeaderItem) Title() string       { return "======== " + shi.title + " ========" }
func (shi sectionHeaderItem) Description() string { return "" }
func (shi sectionHeaderItem) FilterValue() string { return "" }
func (shi sectionHeaderItem) Open() (Panel, bool) { return nil, false }

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

// resultTypeItem displays result type information
type resultTypeItem struct {
	typeName string
}

func (rti resultTypeItem) Title() string       { return rti.typeName }
func (rti resultTypeItem) Description() string { return "" }
func (rti resultTypeItem) FilterValue() string { return rti.typeName }
func (rti resultTypeItem) Open() (Panel, bool) { return nil, false }

// directiveDefinitionItem displays directive type information
type directiveDefinitionItem struct {
	directive *ast.DirectiveDefinition
}

func (ddi directiveDefinitionItem) Title() string { return ddi.directive.Name.Value }
func (ddi directiveDefinitionItem) Description() string {
	if ddi.directive.Description != nil {
		return ddi.directive.Description.Value
	}
	return ""
}
func (ddi directiveDefinitionItem) FilterValue() string { return ddi.directive.Name.Value }
func (ddi directiveDefinitionItem) Open() (Panel, bool) { return nil, false }
