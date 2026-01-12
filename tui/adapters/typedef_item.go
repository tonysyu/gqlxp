package adapters

import (
	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/gqlfmt"
	"github.com/tonysyu/gqlxp/tui/xplr/components"
)

// Adapter/delegate for gql.TypeDef to support ListItem interface
//
// gql.TypeDef is an interface for all of the following:
// - gql.Enum
// - gql.InputObject
// - gql.Interface
// - gql.Object
// - gql.Scalar
// - gql.Union
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
	var tabs []components.Tab

	switch typeDef := (i.typeDef).(type) {
	case *gql.Object:
		tabs = append(tabs, newFieldsTab(typeDef.Fields(), i.resolver))
		if interfaces := typeDef.Interfaces(); len(interfaces) > 0 {
			tabs = append(tabs, newInterfacesTab(interfaces, i.resolver))
		}
		if usages, _ := i.resolver.ResolveUsages(typeDef.Name()); len(usages) > 0 {
			tabs = append(tabs, newUsagesTab(usages, i.resolver))
		}
		if len(typeDef.Directives()) > 0 {
			tabs = append(tabs, newDirectivesTab(typeDef.Directives(), i.resolver))
		}
	case *gql.Scalar:
		// Note that Scalars have no "default" data (e.g. fields, values, etc.)
		if usages, _ := i.resolver.ResolveUsages(typeDef.Name()); len(usages) > 0 {
			tabs = append(tabs, newUsagesTab(usages, i.resolver))
		}
		if len(typeDef.Directives()) > 0 {
			tabs = append(tabs, newDirectivesTab(typeDef.Directives(), i.resolver))
		}
	case *gql.Interface:
		tabs = append(tabs, newFieldsTab(typeDef.Fields(), i.resolver))
		if interfaces := typeDef.Interfaces(); len(interfaces) > 0 {
			tabs = append(tabs, newInterfacesTab(interfaces, i.resolver))
		}
		if usages, _ := i.resolver.ResolveUsages(typeDef.Name()); len(usages) > 0 {
			tabs = append(tabs, newUsagesTab(usages, i.resolver))
		}
		if len(typeDef.Directives()) > 0 {
			tabs = append(tabs, newDirectivesTab(typeDef.Directives(), i.resolver))
		}
	case *gql.Union:
		tabs = append(tabs, components.Tab{
			Label:   "Types",
			Content: adaptUnionTypes(typeDef.Types(), i.resolver),
		})
		if usages, _ := i.resolver.ResolveUsages(typeDef.Name()); len(usages) > 0 {
			tabs = append(tabs, newUsagesTab(usages, i.resolver))
		}
		if len(typeDef.Directives()) > 0 {
			tabs = append(tabs, newDirectivesTab(typeDef.Directives(), i.resolver))
		}
	case *gql.Enum:
		tabs = append(tabs, components.Tab{
			Label:   "Values",
			Content: adaptEnumValues(typeDef.Values()),
		})
		if usages, _ := i.resolver.ResolveUsages(typeDef.Name()); len(usages) > 0 {
			tabs = append(tabs, newUsagesTab(usages, i.resolver))
		}
		if len(typeDef.Directives()) > 0 {
			tabs = append(tabs, newDirectivesTab(typeDef.Directives(), i.resolver))
		}
	case *gql.InputObject:
		tabs = append(tabs, newFieldsTab(typeDef.Fields(), i.resolver))
		if usages, _ := i.resolver.ResolveUsages(typeDef.Name()); len(usages) > 0 {
			tabs = append(tabs, newUsagesTab(usages, i.resolver))
		}
		if len(typeDef.Directives()) > 0 {
			tabs = append(tabs, newDirectivesTab(typeDef.Directives(), i.resolver))
		}
	}

	// Pass empty list of items since all the default types have SubTabs
	panel := components.NewPanel([]components.ListItem{}, i.Title())

	// Add description as a header if available
	if desc := i.Description(); desc != "" {
		panel.SetDescription(desc)
	}

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
