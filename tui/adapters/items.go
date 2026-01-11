package adapters

import (
	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/gqlfmt"
	"github.com/tonysyu/gqlxp/tui/xplr/components"
	"github.com/tonysyu/gqlxp/utils/text"
)

// Ensure that all item types implements components.ListItem interface
var _ components.ListItem = (*fieldItem)(nil)
var _ components.ListItem = (*argumentItem)(nil)
var _ components.ListItem = (*typeDefItem)(nil)
var _ components.ListItem = (*directiveDefItem)(nil)
var _ components.ListItem = (*appliedDirectiveItem)(nil)

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
	return gqlfmt.GenerateFieldMarkdown(i.gqlField, i.resolver)
}

// OpenPanel displays arguments of field (if any) and the field's ObjectType
func (i fieldItem) OpenPanel() (*components.Panel, bool) {
	argumentItems := AdaptArguments(i.gqlField.Arguments(), i.resolver)
	resultTypeItem := newTypeDefItemFromField(i.gqlField, i.resolver)
	directiveItems := AdaptAppliedDirectives(i.gqlField.Directives(), i.resolver)

	panel := components.NewPanel([]components.ListItem{}, i.fieldName)
	panel.SetDescription(i.Description())

	// Create tabs for Result Type and Input Arguments
	var tabs []components.Tab
	tabs = append(tabs, components.Tab{
		Label:   "Type",
		Content: []components.ListItem{resultTypeItem},
	})
	if len(argumentItems) > 0 {
		tabs = append(tabs, components.Tab{
			Label:   "Inputs",
			Content: argumentItems,
		})
	}
	if len(directiveItems) > 0 {
		tabs = append(tabs, components.Tab{
			Label:   "Directives",
			Content: directiveItems,
		})
	}
	panel.SetTabs(tabs)

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
	resultTypeItem := newTypeDefItemFromArgument(i.gqlArgument, i.resolver)

	panel := components.NewPanel([]components.ListItem{}, i.argName)
	panel.SetDescription(i.Description())

	// Create a single tab for Result Type
	panel.SetTabs([]components.Tab{
		{
			Label:   "Type",
			Content: []components.ListItem{resultTypeItem},
		},
	})

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
	return gqlfmt.GenerateTypeDefMarkdown(i.typeDef, i.resolver)
}

// OpenPanel displays list of fields on type (if any)
func (i typeDefItem) OpenPanel() (*components.Panel, bool) {
	// Create list items for the detail view
	var detailItems []components.ListItem
	var tabs []components.Tab
	var directiveItems []components.ListItem

	switch typeDef := (i.typeDef).(type) {
	case *gql.Object:
		detailItems = append(detailItems, AdaptFields(typeDef.Fields(), i.resolver)...)
		tabs = append(tabs, components.Tab{
			Label:   "Fields",
			Content: detailItems,
		})
		// Add Interfaces tab if the object implements any interfaces
		if interfaces := typeDef.Interfaces(); len(interfaces) > 0 {
			tabs = append(tabs, components.Tab{
				Label:   "Interfaces",
				Content: AdaptInterfaces(interfaces, i.resolver),
			})
		}
		// Add Usages tab if the type is used elsewhere
		if usages, _ := i.resolver.ResolveUsages(typeDef.Name()); len(usages) > 0 {
			tabs = append(tabs, components.Tab{
				Label:   "Usages",
				Content: AdaptUsages(usages, i.resolver),
			})
		}
		directiveItems = AdaptAppliedDirectives(typeDef.Directives(), i.resolver)
	case *gql.Scalar:
		// Add Usages tab if the type is used elsewhere
		if usages, _ := i.resolver.ResolveUsages(typeDef.Name()); len(usages) > 0 {
			tabs = append(tabs, components.Tab{
				Label:   "Usages",
				Content: AdaptUsages(usages, i.resolver),
			})
		}
		directiveItems = AdaptAppliedDirectives(typeDef.Directives(), i.resolver)
	case *gql.Interface:
		detailItems = append(detailItems, AdaptFields(typeDef.Fields(), i.resolver)...)
		tabs = append(tabs, components.Tab{
			Label:   "Fields",
			Content: detailItems,
		})
		// Add Interfaces tab if the interface implements any interfaces
		if interfaces := typeDef.Interfaces(); len(interfaces) > 0 {
			tabs = append(tabs, components.Tab{
				Label:   "Interfaces",
				Content: AdaptInterfaces(interfaces, i.resolver),
			})
		}
		// Add Usages tab if the type is used elsewhere
		if usages, _ := i.resolver.ResolveUsages(typeDef.Name()); len(usages) > 0 {
			tabs = append(tabs, components.Tab{
				Label:   "Usages",
				Content: AdaptUsages(usages, i.resolver),
			})
		}
		directiveItems = AdaptAppliedDirectives(typeDef.Directives(), i.resolver)
	case *gql.Union:
		detailItems = append(detailItems, AdaptUnionTypes(typeDef.Types(), i.resolver)...)
		// Add Usages tab if the type is used elsewhere
		if usages, _ := i.resolver.ResolveUsages(typeDef.Name()); len(usages) > 0 {
			tabs = append(tabs, components.Tab{
				Label:   "Usages",
				Content: AdaptUsages(usages, i.resolver),
			})
		}
		directiveItems = AdaptAppliedDirectives(typeDef.Directives(), i.resolver)
	case *gql.Enum:
		detailItems = append(detailItems, AdaptEnumValues(typeDef.Values())...)
		// Add Usages tab if the type is used elsewhere
		if usages, _ := i.resolver.ResolveUsages(typeDef.Name()); len(usages) > 0 {
			tabs = append(tabs, components.Tab{
				Label:   "Usages",
				Content: AdaptUsages(usages, i.resolver),
			})
		}
		directiveItems = AdaptAppliedDirectives(typeDef.Directives(), i.resolver)
	case *gql.InputObject:
		detailItems = append(detailItems, AdaptFields(typeDef.Fields(), i.resolver)...)
		// Only add tabs if we have usages or directives
		if usages, _ := i.resolver.ResolveUsages(typeDef.Name()); len(usages) > 0 {
			tabs = append(tabs, components.Tab{
				Label:   "Fields",
				Content: detailItems,
			})
			tabs = append(tabs, components.Tab{
				Label:   "Usages",
				Content: AdaptUsages(usages, i.resolver),
			})
		}
		directiveItems = AdaptAppliedDirectives(typeDef.Directives(), i.resolver)
	}

	// Add Directives tab if the type has directives
	if len(directiveItems) > 0 {
		tabs = append(tabs, components.Tab{
			Label:   "Directives",
			Content: directiveItems,
		})
	}

	panel := components.NewPanel(detailItems, i.Title())
	// Add description as a header if available
	if desc := i.Description(); desc != "" {
		panel.SetDescription(desc)
	}
	// Set tabs if any were created
	if len(tabs) > 0 {
		panel.SetTabs(tabs)
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

// Adapter/delegate for gql.DirectiveDef to support ListItem interface (for schema directives)
type directiveDefItem struct {
	gqlDirective  *gql.DirectiveDef
	resolver      gql.TypeResolver
	directiveName string
}

func newDirectiveDefItem(directive *gql.DirectiveDef, resolver gql.TypeResolver) components.ListItem {
	return directiveDefItem{
		gqlDirective:  directive,
		resolver:      resolver,
		directiveName: directive.Name(),
	}
}

func (i directiveDefItem) Title() string       { return i.gqlDirective.Signature() }
func (i directiveDefItem) FilterValue() string { return i.directiveName }
func (i directiveDefItem) TypeName() string    { return "@" + i.directiveName }
func (i directiveDefItem) RefName() string     { return i.directiveName }

func (i directiveDefItem) Description() string {
	return i.gqlDirective.Description()
}

func (i directiveDefItem) Details() string {
	return gqlfmt.GenerateDirectiveMarkdown(i.gqlDirective, i.resolver)
}

// OpenPanel displays arguments of directive (if any)
func (i directiveDefItem) OpenPanel() (*components.Panel, bool) {
	argumentItems := AdaptArguments(i.gqlDirective.Arguments(), i.resolver)

	panel := components.NewPanel([]components.ListItem{}, "@"+i.directiveName)
	panel.SetDescription(i.Description())

	var tabs []components.Tab
	if len(argumentItems) > 0 {
		tabs = append(tabs, components.Tab{
			Label:   "Arguments",
			Content: argumentItems,
		})
	}

	// Add Usages tab if the directive is used elsewhere
	if usages, _ := i.resolver.ResolveUsages(i.directiveName); len(usages) > 0 {
		tabs = append(tabs, components.Tab{
			Label:   "Usages",
			Content: AdaptUsages(usages, i.resolver),
		})
	}

	if len(tabs) > 0 {
		panel.SetTabs(tabs)
	}

	return panel, true
}

// Adapter/delegate for gql.AppliedDirective to support ListItem interface (for directives on fields/types)
type appliedDirectiveItem struct {
	gqlDirective  *gql.AppliedDirective
	resolver      gql.TypeResolver
	directiveName string
}

func newAppliedDirectiveItem(directive *gql.AppliedDirective, resolver gql.TypeResolver) components.ListItem {
	return appliedDirectiveItem{
		gqlDirective:  directive,
		resolver:      resolver,
		directiveName: directive.Name(),
	}
}

func (i appliedDirectiveItem) Title() string       { return i.gqlDirective.Signature() }
func (i appliedDirectiveItem) FilterValue() string { return i.directiveName }
func (i appliedDirectiveItem) TypeName() string    { return "@" + i.directiveName }
func (i appliedDirectiveItem) RefName() string     { return i.directiveName }

func (i appliedDirectiveItem) Description() string {
	return ""
}

func (i appliedDirectiveItem) Details() string {
	return text.JoinParagraphs(
		text.H1("@"+i.directiveName),
		text.GqlCode(i.gqlDirective.Signature()),
	)
}

// OpenPanel displays arguments from the directive definition
func (i appliedDirectiveItem) OpenPanel() (*components.Panel, bool) {
	// Resolve the directive name to its definition
	directiveDef, err := i.resolver.ResolveDirective(i.directiveName)
	if err != nil {
		// If we can't resolve the directive, return no panel
		return nil, false
	}

	// Get the arguments from the directive definition
	argumentItems := AdaptArguments(directiveDef.Arguments(), i.resolver)

	panel := components.NewPanel(argumentItems, "@"+i.directiveName)
	panel.SetDescription(directiveDef.Description())

	return panel, true
}
