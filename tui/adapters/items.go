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

func adaptFieldDefinitionsToItems(queryFields []*gql.Field, schema *gql.GraphQLSchema) []components.ListItem {
	adaptedItems := make([]components.ListItem, 0, len(queryFields))
	for _, f := range queryFields {
		adaptedItems = append(adaptedItems, newFieldDefItem(f, schema))
	}
	return adaptedItems
}

func adaptTypeDefsToItems[T gql.TypeDef](typeDefs []T, schema *gql.GraphQLSchema) []components.ListItem {
	adaptedItems := make([]components.ListItem, 0, len(typeDefs))
	for _, td := range typeDefs {
		adaptedItems = append(adaptedItems, newTypeDefItem(td, schema))
	}
	return adaptedItems
}

func adaptObjectDefinitionsToItems(objects []*gql.Object, schema *gql.GraphQLSchema) []components.ListItem {
	return adaptTypeDefsToItems(objects, schema)
}

func adaptInputDefinitionsToItems(inputs []*gql.InputObject, schema *gql.GraphQLSchema) []components.ListItem {
	return adaptTypeDefsToItems(inputs, schema)
}

func adaptEnumDefinitionsToItems(enums []*gql.Enum, schema *gql.GraphQLSchema) []components.ListItem {
	return adaptTypeDefsToItems(enums, schema)
}

func adaptScalarDefinitionsToItems(scalars []*gql.Scalar, schema *gql.GraphQLSchema) []components.ListItem {
	return adaptTypeDefsToItems(scalars, schema)
}

func adaptInterfaceDefinitionsToItems(interfaces []*gql.Interface, schema *gql.GraphQLSchema) []components.ListItem {
	return adaptTypeDefsToItems(interfaces, schema)
}

func adaptUnionDefinitionsToItems(unions []*gql.Union, schema *gql.GraphQLSchema) []components.ListItem {
	return adaptTypeDefsToItems(unions, schema)
}

func adaptDirectiveDefinitionsToItems(directives []*gql.Directive) []components.ListItem {
	adaptedItems := make([]components.ListItem, 0, len(directives))
	for _, directive := range directives {
		adaptedItems = append(adaptedItems, newDirectiveDefinitionItem(directive))
	}
	return adaptedItems
}

func adaptNamedToItems(typeNames []string) []components.ListItem {
	adaptedItems := make([]components.ListItem, 0, len(typeNames))
	for _, typeName := range typeNames {
		adaptedItems = append(adaptedItems, newNamedItem(typeName))
	}
	return adaptedItems
}

func adaptEnumValueDefinitionsToItems(enumNodes []*gql.EnumValue) []components.ListItem {
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

func formatFieldDefinitionsToCodeBlock(fieldNodes []*gql.Field) string {
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
	gqlField  *gql.Field
	schema    *gql.GraphQLSchema
	fieldName string
}

func newFieldDefItem(gqlField *gql.Field, schema *gql.GraphQLSchema) components.ListItem {
	return fieldItem{
		gqlField:  gqlField,
		schema:    schema,
		fieldName: gqlField.Name(),
	}
}

func (i fieldItem) Title() string       { return i.gqlField.Signature() }
func (i fieldItem) FilterValue() string { return i.fieldName }
func (i fieldItem) TypeName() string    { return i.gqlField.ObjectTypeName() }

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

// OpenPanel displays arguments of field (if any) and the field's ObjectType
func (i fieldItem) OpenPanel() (components.Panel, bool) {
	argumentItems := adaptArguments(i.gqlField.Arguments())

	panel := components.NewListPanel(argumentItems, i.fieldName)
	panel.SetDescription(i.Description())
	// Set result type as virtual item at top
	panel.SetObjectType(newFieldTypeItem(i.gqlField, i.schema))

	return panel, true
}

// Create an array of ListItem instances for field arguments
func adaptArguments(arguments []*gql.Argument) []components.ListItem {
	var items []components.ListItem
	if len(arguments) > 0 {
		for _, arg := range arguments {
			items = append(items, newArgumentItem(arg))
		}
	}
	return items
}

// Create an array of ListItem instances for input object fields
func adaptInputFields(fields []*gql.InputField) []components.ListItem {
	var items []components.ListItem
	if len(fields) > 0 {
		for _, field := range fields {
			items = append(items, newInputFieldItem(field))
		}
	}
	return items
}

// Adapter/delegate for gql.NamedTypeDef to support ListItem interface
type typeDefItem struct {
	title    string
	typeName string
	typeDef  gql.TypeDef
	schema   *gql.GraphQLSchema
}

func newTypeDefItem(typeDef gql.TypeDef, schema *gql.GraphQLSchema) typeDefItem {
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
	case *gql.Object:
		if len(typeDef.Interfaces()) > 0 {
			parts = append(parts, "**Implements:** "+strings.Join(typeDef.Interfaces(), ", "))
		}
		codeBlock := formatFieldDefinitionsToCodeBlock(typeDef.Fields())
		if len(codeBlock) > 0 {
			parts = append(parts, codeBlock)
		}
	case *gql.Scalar:
		parts = append(parts, "_Scalar type_")
	case *gql.Interface:
		codeBlock := formatFieldDefinitionsToCodeBlock(typeDef.Fields())
		if len(codeBlock) > 0 {
			parts = append(parts, codeBlock)
		}
	case *gql.Union:
		if len(typeDef.Types()) > 0 {
			parts = append(parts, "**Union of:** "+strings.Join(typeDef.Types(), " | "))
		}
	case *gql.Enum:
		if len(typeDef.Values()) > 0 {
			var values []string
			for _, val := range typeDef.Values() {
				values = append(values, val.Name())
			}
			parts = append(parts, text.GqlCode(text.JoinLines(values...)))
		}
	case *gql.InputObject:
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
func (i typeDefItem) OpenPanel() (components.Panel, bool) {
	// Create list items for the detail view
	var detailItems []components.ListItem

	switch typeDef := (i.typeDef).(type) {
	case *gql.Object:
		detailItems = append(detailItems, adaptFieldDefinitionsToItems(typeDef.Fields(), i.schema)...)
	case *gql.Scalar:
		// No details needed
	case *gql.Interface:
		detailItems = append(detailItems, adaptFieldDefinitionsToItems(typeDef.Fields(), i.schema)...)
	case *gql.Union:
		detailItems = append(detailItems, adaptNamedToItems(typeDef.Types())...)
	case *gql.Enum:
		detailItems = append(detailItems, adaptEnumValueDefinitionsToItems(typeDef.Values())...)
	case *gql.InputObject:
		detailItems = append(detailItems, adaptInputFields(typeDef.Fields())...)
	}

	panel := components.NewListPanel(detailItems, i.Title())
	// Add description as a header if available
	if desc := i.Description(); desc != "" {
		panel.SetDescription(desc)
	}
	return panel, true
}

func newNamedItem(typeName string) components.SimpleItem {
	// TODO: This probably requires a reference to the schema to return full type when opening
	return components.NewSimpleItem(typeName)
}

// newFieldTypeItem creates a list item for a field's result type
func newFieldTypeItem(field *gql.Field, schema *gql.GraphQLSchema) components.ListItem {
	resultType, err := field.ResolveObjectTypeDef(schema)
	if err != nil {
		// FIXME: Currently, this treats any error as a built-in type, but instead we should
		// check for _known_ built in types and handle errors intelligently.
		return components.NewSimpleItem(
			field.TypeString(),
			components.WithTypeName(field.ObjectTypeName()),
		)
	}
	return typeDefItem{
		title:    field.TypeString(),
		typeName: resultType.Name(),
		typeDef:  resultType,
		schema:   schema,
	}
}

func newArgumentItem(argument *gql.Argument) components.SimpleItem {
	// TODO: Update item to support proper Open and use custom display string
	return components.NewSimpleItem(
		argument.Signature(),
		components.WithDescription(argument.Description()),
	)
}

func newInputFieldItem(field *gql.InputField) components.SimpleItem {
	// TODO: Update item to support proper Open and use custom display string
	return components.NewSimpleItem(
		field.Signature(),
		components.WithDescription(field.Description()),
	)
}

func newDirectiveDefinitionItem(directive *gql.Directive) components.SimpleItem {
	return components.NewSimpleItem(
		directive.Name(),
		components.WithDescription(directive.Description()),
	)
}
