package adapters

import (
	"strings"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/tui/components"
	"github.com/tonysyu/gqlxp/utils/text"
)

// Ensure that all item types implements components.ListItem interface
var _ components.ListItem = (*fieldItem)(nil)
var _ components.ListItem = (*typeDefItem)(nil)

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
		item := components.NewSimpleItem(
			node.Name.Value,
			components.WithDescription(gql.GetStringValue(node.Description)),
		)
		adaptedItems = append(adaptedItems, item)
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
	return text.GqlCode(text.JoinLines(fields...))
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

func (i fieldItem) Title() string       { return i.gqlField.Name.Value }
func (i fieldItem) FilterValue() string { return i.Title() }
func (i fieldItem) TypeName() string {
	resultType := gql.GetNamedFromType(i.gqlField.Type)
	return resultType.Name.Value
}

func (i fieldItem) Description() string {
	return gql.GetStringValue(i.gqlField.GetDescription())
}

func (i fieldItem) Details() string {
	return text.JoinParagraphs(
		text.H1(i.TypeName()),
		text.GqlCode(gql.GetFieldDefinitionString(i.gqlField)),
		i.Description(),
	)
}

// Implement components.ListItem interface
func (i fieldItem) Open() (components.Panel, bool) {
	// Create list items for the detail view
	var detailItems []components.ListItem

	inputValueItems := adaptInputValueDefinitions(i.gqlField.Arguments)
	if len(inputValueItems) > 0 {
		detailItems = append(detailItems, newSectionHeader("Input Arguments"))
		detailItems = append(detailItems, inputValueItems...)
	}

	// Add result type section
	detailItems = append(detailItems, newSectionHeader("Result Type"))
	detailItems = append(detailItems, newTypeItem(i.gqlField.Type, i.schema))
	panel := components.NewListPanel(detailItems, i.Title())

	// Add description as a header if available
	if desc := i.Description(); desc != "" {
		panel.SetDescription(desc)
	}
	return panel, true
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
	title    string
	typeName string
	typeDef  gql.NamedTypeDef
	schema   *gql.GraphQLSchema
}

func newTypeDefItem(typeDef gql.NamedTypeDef, schema *gql.GraphQLSchema) typeDefItem {
	title := typeDef.GetName().Value
	return typeDefItem{
		title:    title,
		typeName: title,
		typeDef:  typeDef,
		schema:   schema,
	}
}

func (i typeDefItem) Title() string       { return i.title }
func (i typeDefItem) FilterValue() string { return i.Title() }
func (i typeDefItem) TypeName() string    { return i.typeName }

func (i typeDefItem) Description() string {
	if desc := (i.typeDef).GetDescription(); desc != nil {
		return desc.Value
	}
	return ""
}

func (i typeDefItem) Details() string {
	parts := []string{text.H1(i.TypeName())}

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
			parts = append(parts, text.GqlCode(text.JoinLines(values...)))
		}
	case *ast.InputObjectDefinition:
		if len(typeDef.Fields) > 0 {
			var fields []string
			for _, field := range typeDef.Fields {
				fields = append(fields, gql.GetInputValueDefinitionString(field))
			}
			parts = append(parts, text.GqlCode(text.JoinLines(fields...)))
		}
	}

	return text.JoinParagraphs(parts...)
}

// Implement components.ListItem interface
func (i typeDefItem) Open() (components.Panel, bool) {
	// Create list items for the detail view
	var detailItems []components.ListItem

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

	panel := components.NewListPanel(detailItems, i.Title())
	// Add description as a header if available
	if desc := i.Description(); desc != "" {
		panel.SetDescription(desc)
	}
	return panel, true
}

func newSectionHeader(title string) components.SimpleItem {
	return components.NewSimpleItem("======== " + title + " ========")
}

func newNamedItem(node *ast.Named) components.SimpleItem {
	// TODO: This probably requires a reference to the schema to return full type when opening
	return components.NewSimpleItem(node.Name.Value)
}

func newTypeItem(t ast.Type, schema *gql.GraphQLSchema) components.ListItem {
	resultType, err := schema.NamedToTypeDefinition(gql.GetNamedFromType(t))
	if err != nil {
		// FIXME: Currently, this treats any error as a built-in type, but instead we should
		// check for _known_ built in types and handle errors intelligently.
		return components.NewSimpleItem(
			gql.GetTypeString(t),
			components.WithTypeName(gql.GetNamedFromType(t).Name.Value),
		)
	}
	return typeDefItem{
		title:    gql.GetTypeString(t),
		typeName: resultType.GetName().Value,
		typeDef:  resultType,
		schema:   schema,
	}
}

func newInputValueItem(inputValue *ast.InputValueDefinition) components.SimpleItem {
	// TODO: Update item to support proper Open and use custom display string
	return components.NewSimpleItem(
		gql.GetInputValueDefinitionString(inputValue),
		components.WithDescription(gql.GetStringValue(inputValue.Description)),
	)
}

func newDirectiveDefinitionItem(directive *ast.DirectiveDefinition) components.SimpleItem {
	return components.NewSimpleItem(
		directive.Name.Value,
		components.WithDescription(gql.GetStringValue(directive.Description)),
	)
}
