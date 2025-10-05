package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/tonysyu/igq/gql"
)

// list item that can be "opened" to provide additional information about the item.
// The opened data is represented as a Panel instance that can be rendered to users.
type ListItem interface {
	list.DefaultItem

	// Open Panel to show additional information.
	Open() (Panel, bool)

	// Details returns markdown-formatted details for the item.
	Details() string
}

// Ensure that all item types implements ListItem interface
var _ ListItem = (*fieldItem)(nil)
var _ ListItem = (*typeDefItem)(nil)
var _ ListItem = (*simpleItem)(nil)

func adaptFieldDefinitionsToItems(queryFields []*ast.FieldDefinition) []ListItem {
	adaptedItems := make([]ListItem, 0, len(queryFields))
	for _, f := range queryFields {
		adaptedItems = append(adaptedItems, newFieldDefItem(f))
	}
	return adaptedItems
}

func adaptObjectDefinitionsToItems(objects []*ast.ObjectDefinition) []ListItem {
	adaptedItems := make([]ListItem, 0, len(objects))
	for _, obj := range objects {
		adaptedItems = append(adaptedItems, newTypeDefItem(obj))
	}
	return adaptedItems
}

func adaptInputDefinitionsToItems(inputs []*ast.InputObjectDefinition) []ListItem {
	adaptedItems := make([]ListItem, 0, len(inputs))
	for _, input := range inputs {
		adaptedItems = append(adaptedItems, newTypeDefItem(input))
	}
	return adaptedItems
}

func adaptEnumDefinitionsToItems(enums []*ast.EnumDefinition) []ListItem {
	adaptedItems := make([]ListItem, 0, len(enums))
	for _, enum := range enums {
		adaptedItems = append(adaptedItems, newTypeDefItem(enum))
	}
	return adaptedItems
}

func adaptScalarDefinitionsToItems(scalars []*ast.ScalarDefinition) []ListItem {
	adaptedItems := make([]ListItem, 0, len(scalars))
	for _, scalar := range scalars {
		adaptedItems = append(adaptedItems, newTypeDefItem(scalar))
	}
	return adaptedItems
}

func adaptInterfaceDefinitionsToItems(interfaces []*ast.InterfaceDefinition) []ListItem {
	adaptedItems := make([]ListItem, 0, len(interfaces))
	for _, iface := range interfaces {
		adaptedItems = append(adaptedItems, newTypeDefItem(iface))
	}
	return adaptedItems
}

func adaptUnionDefinitionsToItems(unions []*ast.UnionDefinition) []ListItem {
	adaptedItems := make([]ListItem, 0, len(unions))
	for _, union := range unions {
		adaptedItems = append(adaptedItems, newTypeDefItem(union))
	}
	return adaptedItems
}

func adaptDirectiveDefinitionsToItems(directives []*ast.DirectiveDefinition) []ListItem {
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

func adaptEnumValueDefinitionsToItems(enumNodes []*ast.EnumValueDefinition) []ListItem {
	adaptedItems := make([]ListItem, 0, len(enumNodes))
	for _, node := range enumNodes {
		adaptedItems = append(adaptedItems, simpleItem{
			title:       node.Name.Value,
			description: gql.GetStringValue(node.Description),
		})
	}
	return adaptedItems
}

func adaptFieldDefinitionsToCodeBlock(fieldNodes []*ast.FieldDefinition) string {
	if len(fieldNodes) == 0 {
		return ""
	}
	var fields []string
	for _, field := range fieldNodes {
		fields = append(fields, gql.GetFieldDefinitionString(field))
	}
	return gqlCode(joinLines(fields...))
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
	return gql.GetStringValue(i.gqlField.GetDescription())
}

func (i fieldItem) Details() string {
	return joinParagraphs(
		h1(i.Title()),
		gqlCode(gql.GetFieldDefinitionString(i.gqlField)),
		i.Description(),
	)
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
	// TODO: Can this reuse gql.GetTypeName? The fact that `i.typeDef` is the interface seems to
	// cause problems using it directly.
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

func (i typeDefItem) Details() string {
	parts := []string{h1(i.Title())}

	// Add description if available
	if desc := i.Description(); desc != "" {
		parts = append(parts, desc)
	}

	// Add type-specific details
	switch typeDef := (i.typeDef).(type) {
	case *ast.ObjectDefinition:
		if len(typeDef.Interfaces) > 0 {
			interfaceNames := make([]string, len(typeDef.Interfaces))
			for i, iface := range typeDef.Interfaces {
				interfaceNames[i] = iface.Name.Value
			}
			parts = append(parts, "**Implements:** "+strings.Join(interfaceNames, ", "))
		}
		codeBlock := adaptFieldDefinitionsToCodeBlock(typeDef.Fields)
		if len(codeBlock) > 0 {
			parts = append(parts, codeBlock)
		}
	case *ast.ScalarDefinition:
		parts = append(parts, "_Scalar type_")
	case *ast.InterfaceDefinition:
		codeBlock := adaptFieldDefinitionsToCodeBlock(typeDef.Fields)
		if len(codeBlock) > 0 {
			parts = append(parts, codeBlock)
		}
	case *ast.UnionDefinition:
		if len(typeDef.Types) > 0 {
			typeNames := make([]string, len(typeDef.Types))
			for i, t := range typeDef.Types {
				typeNames[i] = t.Name.Value
			}
			parts = append(parts, "**Union of:** "+strings.Join(typeNames, " | "))
		}
	case *ast.EnumDefinition:
		if len(typeDef.Values) > 0 {
			var values []string
			for _, val := range typeDef.Values {
				values = append(values, val.Name.Value)
			}
			parts = append(parts, gqlCode(joinLines(values...)))
		}
	case *ast.InputObjectDefinition:
		if len(typeDef.Fields) > 0 {
			var fields []string
			for _, field := range typeDef.Fields {
				fields = append(fields, gql.GetInputValueDefinitionString(field))
			}
			parts = append(parts, gqlCode(joinLines(fields...)))
		}
	}

	return joinParagraphs(parts...)
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
		detailItems = append(detailItems, adaptFieldDefinitionsToItems(typeDef.Fields)...)
	case *ast.ScalarDefinition:
		// No details needed
	case *ast.InterfaceDefinition:
		detailItems = append(detailItems, adaptFieldDefinitionsToItems(typeDef.Fields)...)
	case *ast.UnionDefinition:
		detailItems = append(detailItems, adaptNamedItems(typeDef.Types)...)
	case *ast.EnumDefinition:
		detailItems = append(detailItems, adaptEnumValueDefinitionsToItems(typeDef.Values)...)
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
func (di simpleItem) Details() string {
	if di.description != "" {
		return joinParagraphs(h1(di.Title()), di.Description())
	}
	return h1(di.Title())
}
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
	return simpleItem{
		title:       gql.GetInputValueDefinitionString(inputValue),
		description: gql.GetStringValue(inputValue.Description),
	}
}

func newDirectiveDefinitionItem(directive *ast.DirectiveDefinition) simpleItem {
	return simpleItem{
		title:       directive.Name.Value,
		description: gql.GetStringValue(directive.Description),
	}
}
