package adapters

import (
	"strings"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/tonysyu/igq/gql"
	"github.com/tonysyu/igq/tui/components"
)

// Ensure that all item types implements components.ListItem interface
var _ components.ListItem = (*fieldItem)(nil)
var _ components.ListItem = (*typeDefItem)(nil)
var _ components.ListItem = (*components.SimpleItem)(nil)

func AdaptFieldDefinitionsToItems(queryFields []*ast.FieldDefinition, schema *gql.GraphQLSchema) []components.ListItem {
	adaptedItems := make([]components.ListItem, 0, len(queryFields))
	for _, f := range queryFields {
		adaptedItems = append(adaptedItems, newFieldDefItem(f, schema))
	}
	return adaptedItems
}

func adaptTypeDefsToItems[T gql.NamedTypeDef](typeDefs []T, schema *gql.GraphQLSchema) []components.ListItem {
	adaptedItems := make([]components.ListItem, 0, len(typeDefs))
	for _, td := range typeDefs {
		adaptedItems = append(adaptedItems, newTypeDefItem(td, schema))
	}
	return adaptedItems
}

func AdaptObjectDefinitionsToItems(objects []*ast.ObjectDefinition, schema *gql.GraphQLSchema) []components.ListItem {
	return adaptTypeDefsToItems(objects, schema)
}

func AdaptInputDefinitionsToItems(inputs []*ast.InputObjectDefinition, schema *gql.GraphQLSchema) []components.ListItem {
	return adaptTypeDefsToItems(inputs, schema)
}

func AdaptEnumDefinitionsToItems(enums []*ast.EnumDefinition, schema *gql.GraphQLSchema) []components.ListItem {
	return adaptTypeDefsToItems(enums, schema)
}

func AdaptScalarDefinitionsToItems(scalars []*ast.ScalarDefinition, schema *gql.GraphQLSchema) []components.ListItem {
	return adaptTypeDefsToItems(scalars, schema)
}

func AdaptInterfaceDefinitionsToItems(interfaces []*ast.InterfaceDefinition, schema *gql.GraphQLSchema) []components.ListItem {
	return adaptTypeDefsToItems(interfaces, schema)
}

func AdaptUnionDefinitionsToItems(unions []*ast.UnionDefinition, schema *gql.GraphQLSchema) []components.ListItem {
	return adaptTypeDefsToItems(unions, schema)
}

func AdaptDirectiveDefinitionsToItems(directives []*ast.DirectiveDefinition) []components.ListItem {
	adaptedItems := make([]components.ListItem, 0, len(directives))
	for _, directive := range directives {
		adaptedItems = append(adaptedItems, newDirectiveDefinitionItem(directive))
	}
	return adaptedItems
}

func adaptNamedToItems(namedNodes []*ast.Named) []components.ListItem {
	adaptedItems := make([]components.ListItem, 0, len(namedNodes))
	for _, node := range namedNodes {
		adaptedItems = append(adaptedItems, newNamedItem(node))
	}
	return adaptedItems
}

func adaptEnumValueDefinitionsToItems(enumNodes []*ast.EnumValueDefinition) []components.ListItem {
	adaptedItems := make([]components.ListItem, 0, len(enumNodes))
	for _, node := range enumNodes {
		adaptedItems = append(adaptedItems, components.NewSimpleItem(
			node.Name.Value,
			gql.GetStringValue(node.Description),
		))
	}
	return adaptedItems
}

func formatFieldDefinitionsToCodeBlock(fieldNodes []*ast.FieldDefinition) string {
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
	schema   *gql.GraphQLSchema
}

func newFieldDefItem(gqlField *ast.FieldDefinition, schema *gql.GraphQLSchema) components.ListItem {
	return fieldItem{
		gqlField: gqlField,
		schema:   schema,
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

// Implement components.ListItem interface
func (i fieldItem) Open() (components.Panel, bool) {
	// Create list items for the detail view
	var detailItems []components.ListItem

	// Add description as a header if available
	if desc := i.Description(); desc != "" {
		detailItems = append(detailItems, components.NewSimpleItem(desc, ""))
	}

	inputValueItems := adaptInputValueDefinitions(i.gqlField.Arguments)
	if len(inputValueItems) > 0 {
		detailItems = append(detailItems, newSectionHeader("Input Arguments"))
		detailItems = append(detailItems, inputValueItems...)
	}

	// Add result type section
	detailItems = append(detailItems, newSectionHeader("Result Type"))
	// TODO: Use NamedToTypeDefinition and newTypeDefItem
	// resultType, err := i.schema.NamedToTypeDefinition(gql.GetNamedFromType(i.gqlField.Type))
	// if err != nil {
	// 	detailItems = append(detailItems, newTypeItem(i.gqlField.Type))
	// } else {
	// 	detailItems = append(detailItems, newTypeDefItem(resultType, i.schema))
	// }
	detailItems = append(detailItems, newTypeItem(i.gqlField.Type))
	// detailItems = append(detailItems, newTypeItem(i.gqlField.Type))
	return components.NewListPanel(detailItems, i.Title()), true
}

// Create an array of ListItem instances given InputValueDefinition. This is used for
// `ast.FieldDefinition.Arguments` and `ast.InputObjectDefinition.Fields`
func adaptInputValueDefinitions(inputValues []*ast.InputValueDefinition) []components.ListItem {
	var items []components.ListItem
	if len(inputValues) > 0 {
		for _, arg := range inputValues {
			items = append(items, newInputValueItem(arg))
		}
	}
	return items

}

// Adapter/delegate for ast.TypeDefinition to support ListItem interface
type typeDefItem struct {
	typeDef gql.NamedTypeDef
	schema  *gql.GraphQLSchema
}

func newTypeDefItem(typeDef gql.NamedTypeDef, schema *gql.GraphQLSchema) typeDefItem {
	return typeDefItem{
		typeDef: typeDef,
		schema:  schema,
	}
}

func (i typeDefItem) Title() string {
	return i.typeDef.GetName().Value
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
		codeBlock := formatFieldDefinitionsToCodeBlock(typeDef.Fields)
		if len(codeBlock) > 0 {
			parts = append(parts, codeBlock)
		}
	case *ast.ScalarDefinition:
		parts = append(parts, "_Scalar type_")
	case *ast.InterfaceDefinition:
		codeBlock := formatFieldDefinitionsToCodeBlock(typeDef.Fields)
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

// Implement components.ListItem interface
func (i typeDefItem) Open() (components.Panel, bool) {
	// Create list items for the detail view
	var detailItems []components.ListItem

	// Add description as a header if available
	if desc := i.Description(); desc != "" {
		detailItems = append(detailItems, components.NewSimpleItem(desc, ""))
	}

	switch typeDef := (i.typeDef).(type) {
	case *ast.ObjectDefinition:
		detailItems = append(detailItems, AdaptFieldDefinitionsToItems(typeDef.Fields, i.schema)...)
	case *ast.ScalarDefinition:
		// No details needed
	case *ast.InterfaceDefinition:
		detailItems = append(detailItems, AdaptFieldDefinitionsToItems(typeDef.Fields, i.schema)...)
	case *ast.UnionDefinition:
		detailItems = append(detailItems, adaptNamedToItems(typeDef.Types)...)
	case *ast.EnumDefinition:
		detailItems = append(detailItems, adaptEnumValueDefinitionsToItems(typeDef.Values)...)
	case *ast.InputObjectDefinition:
		detailItems = append(detailItems, adaptInputValueDefinitions(typeDef.Fields)...)
	}

	return components.NewListPanel(detailItems, i.Title()), true
}

func newSectionHeader(title string) components.SimpleItem {
	return components.NewSimpleItem("======== "+title+" ========", "")
}

func newNamedItem(node *ast.Named) components.SimpleItem {
	// TODO: This probably requires a reference to the schema to return full type when opening
	return components.NewSimpleItem(node.Name.Value, "")
}

func newTypeItem(t ast.Type) components.SimpleItem {
	// TODO: This probably requires a reference to the schema to return full type when opening
	return components.NewSimpleItem(gql.GetTypeString(t), "")
}

func newInputValueItem(inputValue *ast.InputValueDefinition) components.SimpleItem {
	// TODO: Update item to support proper Open and use custom display string
	return components.NewSimpleItem(
		gql.GetInputValueDefinitionString(inputValue),
		gql.GetStringValue(inputValue.Description),
	)
}

func newDirectiveDefinitionItem(directive *ast.DirectiveDefinition) components.SimpleItem {
	return components.NewSimpleItem(
		directive.Name.Value,
		gql.GetStringValue(directive.Description),
	)
}
