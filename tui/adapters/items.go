package adapters

import (
	"strings"

	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/tui/components"
	"github.com/tonysyu/gqlxp/utils/text"
)

// Ensure that all item types implements components.ListItem interface
var _ components.ListItem = (*fieldItem)(nil)
var _ components.ListItem = (*typeDefItem)(nil)

func adaptFieldDefinitionsToItems(queryFields []*gql.FieldDefinition, schema *gql.GraphQLSchema) []components.ListItem {
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

func adaptObjectDefinitionsToItems(objects []*gql.ObjectDefinition, schema *gql.GraphQLSchema) []components.ListItem {
	return adaptTypeDefsToItems(objects, schema)
}

func adaptInputDefinitionsToItems(inputs []*gql.InputObjectDefinition, schema *gql.GraphQLSchema) []components.ListItem {
	return adaptTypeDefsToItems(inputs, schema)
}

func adaptEnumDefinitionsToItems(enums []*gql.EnumDefinition, schema *gql.GraphQLSchema) []components.ListItem {
	return adaptTypeDefsToItems(enums, schema)
}

func adaptScalarDefinitionsToItems(scalars []*gql.ScalarDefinition, schema *gql.GraphQLSchema) []components.ListItem {
	return adaptTypeDefsToItems(scalars, schema)
}

func adaptInterfaceDefinitionsToItems(interfaces []*gql.InterfaceDefinition, schema *gql.GraphQLSchema) []components.ListItem {
	return adaptTypeDefsToItems(interfaces, schema)
}

func adaptUnionDefinitionsToItems(unions []*gql.UnionDefinition, schema *gql.GraphQLSchema) []components.ListItem {
	return adaptTypeDefsToItems(unions, schema)
}

func adaptDirectiveDefinitionsToItems(directives []*gql.DirectiveDefinition) []components.ListItem {
	adaptedItems := make([]components.ListItem, 0, len(directives))
	for _, directive := range directives {
		adaptedItems = append(adaptedItems, newDirectiveDefinitionItem(directive))
	}
	return adaptedItems
}

func adaptNamedToItems(namedNodes []*gql.Named) []components.ListItem {
	adaptedItems := make([]components.ListItem, 0, len(namedNodes))
	for _, node := range namedNodes {
		adaptedItems = append(adaptedItems, newNamedItem(node))
	}
	return adaptedItems
}

func adaptEnumValueDefinitionsToItems(enumNodes []*gql.EnumValueDefinition) []components.ListItem {
	adaptedItems := make([]components.ListItem, 0, len(enumNodes))
	for _, node := range enumNodes {
		item := components.NewSimpleItem(
			node.Name(),
			components.WithDescription(node.Description()),
		)
		adaptedItems = append(adaptedItems, item)
	}
	return adaptedItems
}

func formatFieldDefinitionsToCodeBlock(fieldNodes []*gql.FieldDefinition) string {
	if len(fieldNodes) == 0 {
		return ""
	}
	var fields []string
	for _, field := range fieldNodes {
		fields = append(fields, field.Signature())
	}
	return text.GqlCode(text.JoinLines(fields...))
}

// Adapter/delegate for gql.FieldDefinition to support ListItem interface
type fieldItem struct {
	gqlField  *gql.FieldDefinition
	schema    *gql.GraphQLSchema
	fieldName string
}

func newFieldDefItem(gqlField *gql.FieldDefinition, schema *gql.GraphQLSchema) components.ListItem {
	return fieldItem{
		gqlField:  gqlField,
		schema:    schema,
		fieldName: gqlField.Name(),
	}
}

func (i fieldItem) Title() string       { return i.gqlField.Signature() }
func (i fieldItem) FilterValue() string { return i.fieldName }
func (i fieldItem) TypeName() string    { return i.gqlField.TypeName() }

func (i fieldItem) Description() string {
	return i.gqlField.Description()
}

func (i fieldItem) Details() string {
	return text.JoinParagraphs(
		text.H1(i.TypeName()),
		text.GqlCode(i.gqlField.Signature()),
		i.Description(),
	)
}

// Implement components.ListItem interface
func (i fieldItem) Open() (components.Panel, bool) {
	// Only add actual argument items to list (no section headers)
	inputValueItems := adaptInputValueDefinitions(i.gqlField.Arguments())

	panel := components.NewListPanel(inputValueItems, i.fieldName)

	// Add description as a header if available
	if desc := i.Description(); desc != "" {
		panel.SetDescription(desc)
	}

	// Set result type as virtual item at top
	panel.SetResultType(newFieldTypeItem(i.gqlField, i.schema))

	return panel, true
}

// Create an array of ListItem instances given InputValueDefinition. This is used for
// field arguments and input object fields.
func adaptInputValueDefinitions(inputValues []*gql.InputValueDefinition) []components.ListItem {
	var items []components.ListItem
	if len(inputValues) > 0 {
		for _, arg := range inputValues {
			items = append(items, newInputValueItem(arg))
		}
	}
	return items

}

// Adapter/delegate for gql.NamedTypeDef to support ListItem interface
type typeDefItem struct {
	title    string
	typeName string
	typeDef  gql.NamedTypeDef
	schema   *gql.GraphQLSchema
}

func newTypeDefItem(typeDef gql.NamedTypeDef, schema *gql.GraphQLSchema) typeDefItem {
	title := typeDef.Name()
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
func (i typeDefItem) Description() string { return i.typeDef.Description() }

func (i typeDefItem) Details() string {
	parts := []string{text.H1(i.TypeName())}

	// Add description if available
	if desc := i.Description(); desc != "" {
		parts = append(parts, desc)
	}

	// Add type-specific details
	switch typeDef := (i.typeDef).(type) {
	case *gql.ObjectDefinition:
		if len(typeDef.Interfaces()) > 0 {
			interfaceNames := make([]string, len(typeDef.Interfaces()))
			for i, iface := range typeDef.Interfaces() {
				interfaceNames[i] = iface.Name()
			}
			parts = append(parts, "**Implements:** "+strings.Join(interfaceNames, ", "))
		}
		codeBlock := formatFieldDefinitionsToCodeBlock(typeDef.Fields())
		if len(codeBlock) > 0 {
			parts = append(parts, codeBlock)
		}
	case *gql.ScalarDefinition:
		parts = append(parts, "_Scalar type_")
	case *gql.InterfaceDefinition:
		codeBlock := formatFieldDefinitionsToCodeBlock(typeDef.Fields())
		if len(codeBlock) > 0 {
			parts = append(parts, codeBlock)
		}
	case *gql.UnionDefinition:
		if len(typeDef.Types()) > 0 {
			typeNames := make([]string, len(typeDef.Types()))
			for i, t := range typeDef.Types() {
				typeNames[i] = t.Name()
			}
			parts = append(parts, "**Union of:** "+strings.Join(typeNames, " | "))
		}
	case *gql.EnumDefinition:
		if len(typeDef.Values()) > 0 {
			var values []string
			for _, val := range typeDef.Values() {
				values = append(values, val.Name())
			}
			parts = append(parts, text.GqlCode(text.JoinLines(values...)))
		}
	case *gql.InputObjectDefinition:
		if len(typeDef.Fields()) > 0 {
			var fields []string
			for _, field := range typeDef.Fields() {
				fields = append(fields, field.Signature())
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
	case *gql.ObjectDefinition:
		detailItems = append(detailItems, adaptFieldDefinitionsToItems(typeDef.Fields(), i.schema)...)
	case *gql.ScalarDefinition:
		// No details needed
	case *gql.InterfaceDefinition:
		detailItems = append(detailItems, adaptFieldDefinitionsToItems(typeDef.Fields(), i.schema)...)
	case *gql.UnionDefinition:
		detailItems = append(detailItems, adaptNamedToItems(typeDef.Types())...)
	case *gql.EnumDefinition:
		detailItems = append(detailItems, adaptEnumValueDefinitionsToItems(typeDef.Values())...)
	case *gql.InputObjectDefinition:
		detailItems = append(detailItems, adaptInputValueDefinitions(typeDef.Fields())...)
	}

	panel := components.NewListPanel(detailItems, i.Title())
	// Add description as a header if available
	if desc := i.Description(); desc != "" {
		panel.SetDescription(desc)
	}
	return panel, true
}

func newNamedItem(node *gql.Named) components.SimpleItem {
	// TODO: This probably requires a reference to the schema to return full type when opening
	return components.NewSimpleItem(node.Name())
}

// newFieldTypeItem creates a list item for a field's result type
func newFieldTypeItem(field *gql.FieldDefinition, schema *gql.GraphQLSchema) components.ListItem {
	resultType, err := field.ResolveResultType(schema)
	if err != nil {
		// FIXME: Currently, this treats any error as a built-in type, but instead we should
		// check for _known_ built in types and handle errors intelligently.
		return components.NewSimpleItem(
			field.TypeString(),
			components.WithTypeName(field.TypeName()),
		)
	}
	return typeDefItem{
		title:    field.TypeString(),
		typeName: resultType.Name(),
		typeDef:  resultType,
		schema:   schema,
	}
}

func newInputValueItem(inputValue *gql.InputValueDefinition) components.SimpleItem {
	// TODO: Update item to support proper Open and use custom display string
	return components.NewSimpleItem(
		inputValue.Signature(),
		components.WithDescription(inputValue.Description()),
	)
}

func newDirectiveDefinitionItem(directive *gql.DirectiveDefinition) components.SimpleItem {
	return components.NewSimpleItem(
		directive.Name(),
		components.WithDescription(directive.Description()),
	)
}
