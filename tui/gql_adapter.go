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

func adaptFieldDefinitions(queryFields map[string]*ast.FieldDefinition) []item {
	adaptedItems := make([]item, 0, len(queryFields))
	for _, f := range queryFields {
		adaptedItems = append(adaptedItems, newItem(f))
	}

	return adaptedItems
}

func adaptObjectDefinitions(objects map[string]*ast.ObjectDefinition) []list.Item {
	adaptedItems := make([]list.Item, 0, len(objects))
	for _, obj := range objects {
		adaptedItems = append(adaptedItems, objectDefinitionItem{obj})
	}
	return adaptedItems
}

func adaptInputDefinitions(inputs map[string]*ast.InputObjectDefinition) []list.Item {
	adaptedItems := make([]list.Item, 0, len(inputs))
	for _, input := range inputs {
		adaptedItems = append(adaptedItems, inputDefinitionItem{input})
	}
	return adaptedItems
}

func adaptEnumDefinitions(enums map[string]*ast.EnumDefinition) []list.Item {
	adaptedItems := make([]list.Item, 0, len(enums))
	for _, enum := range enums {
		adaptedItems = append(adaptedItems, enumDefinitionItem{enum})
	}
	return adaptedItems
}

func adaptScalarDefinitions(scalars map[string]*ast.ScalarDefinition) []list.Item {
	adaptedItems := make([]list.Item, 0, len(scalars))
	for _, scalar := range scalars {
		adaptedItems = append(adaptedItems, scalarDefinitionItem{scalar})
	}
	return adaptedItems
}

func adaptInterfaceDefinitions(interfaces map[string]*ast.InterfaceDefinition) []list.Item {
	adaptedItems := make([]list.Item, 0, len(interfaces))
	for _, iface := range interfaces {
		adaptedItems = append(adaptedItems, interfaceDefinitionItem{iface})
	}
	return adaptedItems
}

func adaptUnionDefinitions(unions map[string]*ast.UnionDefinition) []list.Item {
	adaptedItems := make([]list.Item, 0, len(unions))
	for _, union := range unions {
		adaptedItems = append(adaptedItems, unionDefinitionItem{union})
	}
	return adaptedItems
}

func adaptDirectiveDefinitions(directives map[string]*ast.DirectiveDefinition) []list.Item {
	adaptedItems := make([]list.Item, 0, len(directives))
	for _, directive := range directives {
		adaptedItems = append(adaptedItems, directiveDefinitionItem{directive})
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
	return ""
}

func (ai argumentItem) FilterValue() string { return ai.name }

// resultTypeItem displays result type information
type resultTypeItem struct {
	typeName string
}

func (rti resultTypeItem) Title() string       { return rti.typeName }
func (rti resultTypeItem) Description() string { return "" }
func (rti resultTypeItem) FilterValue() string { return rti.typeName }

// objectDefinitionItem displays object type information
type objectDefinitionItem struct {
	obj *ast.ObjectDefinition
}

func (odi objectDefinitionItem) Title() string { return odi.obj.Name.Value }
func (odi objectDefinitionItem) Description() string {
	if odi.obj.Description != nil {
		return odi.obj.Description.Value
	}
	return ""
}
func (odi objectDefinitionItem) FilterValue() string { return odi.obj.Name.Value }

// inputDefinitionItem displays input type information
type inputDefinitionItem struct {
	input *ast.InputObjectDefinition
}

func (idi inputDefinitionItem) Title() string { return idi.input.Name.Value }
func (idi inputDefinitionItem) Description() string {
	if idi.input.Description != nil {
		return idi.input.Description.Value
	}
	return ""
}
func (idi inputDefinitionItem) FilterValue() string { return idi.input.Name.Value }

// enumDefinitionItem displays enum type information
type enumDefinitionItem struct {
	enum *ast.EnumDefinition
}

func (edi enumDefinitionItem) Title() string { return edi.enum.Name.Value }
func (edi enumDefinitionItem) Description() string {
	if edi.enum.Description != nil {
		return edi.enum.Description.Value
	}
	return ""
}
func (edi enumDefinitionItem) FilterValue() string { return edi.enum.Name.Value }

// scalarDefinitionItem displays scalar type information
type scalarDefinitionItem struct {
	scalar *ast.ScalarDefinition
}

func (sdi scalarDefinitionItem) Title() string { return sdi.scalar.Name.Value }
func (sdi scalarDefinitionItem) Description() string {
	if sdi.scalar.Description != nil {
		return sdi.scalar.Description.Value
	}
	return ""
}
func (sdi scalarDefinitionItem) FilterValue() string { return sdi.scalar.Name.Value }

// interfaceDefinitionItem displays interface type information
type interfaceDefinitionItem struct {
	iface *ast.InterfaceDefinition
}

func (idi interfaceDefinitionItem) Title() string { return idi.iface.Name.Value }
func (idi interfaceDefinitionItem) Description() string {
	if idi.iface.Description != nil {
		return idi.iface.Description.Value
	}
	return ""
}
func (idi interfaceDefinitionItem) FilterValue() string { return idi.iface.Name.Value }

// unionDefinitionItem displays union type information
type unionDefinitionItem struct {
	union *ast.UnionDefinition
}

func (udi unionDefinitionItem) Title() string { return udi.union.Name.Value }
func (udi unionDefinitionItem) Description() string {
	if udi.union.Description != nil {
		return udi.union.Description.Value
	}
	return ""
}
func (udi unionDefinitionItem) FilterValue() string { return udi.union.Name.Value }

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
