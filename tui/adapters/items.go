package adapters

import (
	"strings"

	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/tui/xplr/components"
	"github.com/tonysyu/gqlxp/utils/text"
)

// Ensure that all item types implements components.ListItem interface
var _ components.ListItem = (*fieldItem)(nil)
var _ components.ListItem = (*argumentItem)(nil)
var _ components.ListItem = (*typeDefItem)(nil)
var _ components.ListItem = (*directiveItem)(nil)

func formatFieldDefinitionsWithDescriptions(fieldNodes []*gql.Field) string {
	if len(fieldNodes) == 0 {
		return ""
	}
	var parts []string
	for _, field := range fieldNodes {
		fieldParts := []string{}
		if desc := field.Description(); desc != "" {
			fieldParts = append(fieldParts, text.GqlDocString(desc))
		}
		fieldParts = append(fieldParts, field.Signature())
		parts = append(parts, text.JoinLines(fieldParts...))
	}
	return text.GqlCode(text.JoinParagraphs(parts...))
}

func formatEnumValuesWithDescriptions(enumValues []*gql.EnumValue) string {
	if len(enumValues) == 0 {
		return ""
	}
	var parts []string
	for _, val := range enumValues {
		valParts := []string{}
		if desc := val.Description(); desc != "" {
			valParts = append(valParts, text.GqlDocString(desc))
		}
		valParts = append(valParts, val.Name())
		parts = append(parts, text.JoinLines(valParts...))
	}
	return text.GqlCode(text.JoinParagraphs(parts...))
}

// Adapter/delegate for gql.FieldDefinition to support ListItem interface
type fieldItem struct {
	gqlField  *gql.Field
	resolver  gql.TypeResolver
	fieldName string
}

func newFieldItem(gqlField *gql.Field, resolver gql.TypeResolver) components.ListItem {
	return fieldItem{
		gqlField:  gqlField,
		resolver:  resolver,
		fieldName: gqlField.Name(),
	}
}

func (i fieldItem) Title() string       { return i.gqlField.Signature() }
func (i fieldItem) FilterValue() string { return i.fieldName }
func (i fieldItem) TypeName() string    { return i.gqlField.ObjectTypeName() }
func (i fieldItem) RefName() string     { return i.gqlField.Name() }

func (i fieldItem) Description() string {
	return i.gqlField.Description()
}

func (i fieldItem) Details() string {
	return text.JoinParagraphs(
		text.H1(i.RefName()),
		text.GqlCode(i.gqlField.FormatSignature(80)),
		i.Description(),
	)
}

// OpenPanel displays arguments of field (if any) and the field's ObjectType
func (i fieldItem) OpenPanel() (*components.Panel, bool) {
	argumentItems := AdaptArguments(i.gqlField.Arguments(), i.resolver)

	panel := components.NewPanel(argumentItems, i.fieldName)
	panel.SetDescription(i.Description())
	// Set result type as virtual item at top
	panel.SetObjectType(newTypeDefItemFromField(i.gqlField, i.resolver))

	return panel, true
}

// Adapter/delegate for gql.Argument to support ListItem interface
type argumentItem struct {
	gqlArgument *gql.Argument
	resolver    gql.TypeResolver
	argName     string
}

func newArgumentItem(gqlArgument *gql.Argument, resolver gql.TypeResolver) components.ListItem {
	return argumentItem{
		gqlArgument: gqlArgument,
		resolver:    resolver,
		argName:     gqlArgument.Name(),
	}
}

func (i argumentItem) Title() string       { return i.gqlArgument.Signature() }
func (i argumentItem) FilterValue() string { return i.argName }
func (i argumentItem) TypeName() string    { return i.gqlArgument.ObjectTypeName() }
func (i argumentItem) RefName() string     { return i.gqlArgument.Name() }

func (i argumentItem) Description() string {
	return i.gqlArgument.Description()
}

func (i argumentItem) Details() string {
	return text.JoinParagraphs(
		text.H1(i.argName),
		text.GqlCode(i.gqlArgument.FormatSignature(80)),
		i.Description(),
	)
}

// OpenPanel displays the argument's type definition
func (i argumentItem) OpenPanel() (*components.Panel, bool) {
	// Create an empty panel for the argument (similar to how fieldItem creates a panel for arguments)
	panel := components.NewPanel([]components.ListItem{}, i.argName)
	panel.SetDescription(i.Description())

	// Set the argument's type as the object type at the top
	panel.SetObjectType(newTypeDefItemFromArgument(i.gqlArgument, i.resolver))

	return panel, true
}

// Adapter/delegate for gql.NamedTypeDef to support ListItem interface
type typeDefItem struct {
	title    string
	typeName string
	typeDef  gql.TypeDef
	resolver gql.TypeResolver
}

func newTypeDefItem(typeDef gql.TypeDef, resolver gql.TypeResolver) typeDefItem {
	title := typeDef.Name()
	return typeDefItem{
		title:    title,
		typeName: title,
		typeDef:  typeDef,
		resolver: resolver,
	}
}

func (i typeDefItem) Title() string       { return i.title }
func (i typeDefItem) FilterValue() string { return i.Title() }
func (i typeDefItem) TypeName() string    { return i.typeName }
func (i typeDefItem) RefName() string     { return i.typeDef.Name() }
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
		fieldsWithDesc := formatFieldDefinitionsWithDescriptions(typeDef.Fields())
		if len(fieldsWithDesc) > 0 {
			parts = append(parts, fieldsWithDesc)
		}
	case *gql.Scalar:
		parts = append(parts, "_Scalar type_")
	case *gql.Interface:
		fieldsWithDesc := formatFieldDefinitionsWithDescriptions(typeDef.Fields())
		if len(fieldsWithDesc) > 0 {
			parts = append(parts, fieldsWithDesc)
		}
	case *gql.Union:
		if len(typeDef.Types()) > 0 {
			parts = append(parts, "**Union of:** "+strings.Join(typeDef.Types(), " | "))
		}
	case *gql.Enum:
		valuesWithDesc := formatEnumValuesWithDescriptions(typeDef.Values())
		if len(valuesWithDesc) > 0 {
			parts = append(parts, valuesWithDesc)
		}
	case *gql.InputObject:
		fieldsWithDesc := formatFieldDefinitionsWithDescriptions(typeDef.Fields())
		if len(fieldsWithDesc) > 0 {
			parts = append(parts, fieldsWithDesc)
		}
	}

	return text.JoinParagraphs(parts...)
}

// OpenPanel displays list of fields on type (if any)
func (i typeDefItem) OpenPanel() (*components.Panel, bool) {
	// Create list items for the detail view
	var detailItems []components.ListItem

	switch typeDef := (i.typeDef).(type) {
	case *gql.Object:
		detailItems = append(detailItems, AdaptFields(typeDef.Fields(), i.resolver)...)
	case *gql.Scalar:
		// No details needed
	case *gql.Interface:
		detailItems = append(detailItems, AdaptFields(typeDef.Fields(), i.resolver)...)
	case *gql.Union:
		detailItems = append(detailItems, AdaptUnionTypes(typeDef.Types(), i.resolver)...)
	case *gql.Enum:
		detailItems = append(detailItems, AdaptEnumValues(typeDef.Values())...)
	case *gql.InputObject:
		detailItems = append(detailItems, AdaptFields(typeDef.Fields(), i.resolver)...)
	}

	panel := components.NewPanel(detailItems, i.Title())
	// Add description as a header if available
	if desc := i.Description(); desc != "" {
		panel.SetDescription(desc)
	}
	return panel, true
}

func newNamedItem(typeName string, resolver gql.TypeResolver) components.ListItem {
	// Try to resolve the type name to a full TypeDef
	typeDef, err := resolver.ResolveType(typeName)
	if err != nil {
		// Type not found or is a primitive - fallback to simple item
		return components.NewSimpleItem(typeName)
	}
	// Create a typeDefItem for resolvable types
	return newTypeDefItem(typeDef, resolver)
}

// newTypeDefItemFromField creates a list item for a field's result type
func newTypeDefItemFromField(field *gql.Field, resolver gql.TypeResolver) components.ListItem {
	resultType, err := resolver.ResolveFieldType(field)
	if err != nil {
		// TODO: Currently, this treats any error as a built-in type, but instead we should
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
		resolver: resolver,
	}
}

// newTypeDefItemFromArgument creates a list item for an argument's type
func newTypeDefItemFromArgument(argument *gql.Argument, resolver gql.TypeResolver) components.ListItem {
	resultType, err := resolver.ResolveArgumentType(argument)
	if err != nil {
		// TODO: Currently, this treats any error as a built-in type, but instead we should
		// check for _known_ built in types and handle errors intelligently.
		return components.NewSimpleItem(
			argument.TypeString(),
			components.WithTypeName(argument.ObjectTypeName()),
		)
	}
	return typeDefItem{
		title:    argument.TypeString(),
		typeName: resultType.Name(),
		typeDef:  resultType,
		resolver: resolver,
	}
}

// Adapter/delegate for gql.Directive to support ListItem interface
type directiveItem struct {
	gqlDirective  *gql.Directive
	resolver      gql.TypeResolver
	directiveName string
}

func newDirectiveDefinitionItem(directive *gql.Directive, resolver gql.TypeResolver) components.ListItem {
	return directiveItem{
		gqlDirective:  directive,
		resolver:      resolver,
		directiveName: directive.Name(),
	}
}

func (i directiveItem) Title() string       { return i.gqlDirective.Signature() }
func (i directiveItem) FilterValue() string { return i.directiveName }
func (i directiveItem) TypeName() string    { return "@" + i.directiveName }
func (i directiveItem) RefName() string     { return i.directiveName }

func (i directiveItem) Description() string {
	return i.gqlDirective.Description()
}

func (i directiveItem) Details() string {
	parts := []string{
		text.H1(i.TypeName()),
		text.GqlCode(i.gqlDirective.FormatSignature(80)),
		i.Description(),
	}
	if len(i.gqlDirective.Locations()) > 0 {
		locationList := []string{}
		for _, loc := range i.gqlDirective.Locations() {
			locationList = append(locationList, "- "+loc)
		}
		parts = append(parts, "**Locations:**\n"+text.JoinLines(locationList...))
	}
	return text.JoinParagraphs(parts...)
}

// OpenPanel displays arguments of directive (if any)
func (i directiveItem) OpenPanel() (*components.Panel, bool) {
	argumentItems := AdaptArguments(i.gqlDirective.Arguments(), i.resolver)

	panel := components.NewPanel(argumentItems, "@"+i.directiveName)
	panel.SetDescription(i.Description())

	return panel, true
}
