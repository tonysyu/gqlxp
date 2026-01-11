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
	// Create list items for the detail view
	var detailItems []components.ListItem
	var tabs []components.Tab
	var directiveItems []components.ListItem

	switch typeDef := (i.typeDef).(type) {
	case *gql.Object:
		detailItems = append(detailItems, adaptFields(typeDef.Fields(), i.resolver)...)
		tabs = append(tabs, components.Tab{
			Label:   "Fields",
			Content: detailItems,
		})
		// Add Interfaces tab if the object implements any interfaces
		if interfaces := typeDef.Interfaces(); len(interfaces) > 0 {
			tabs = append(tabs, components.Tab{
				Label:   "Interfaces",
				Content: adaptInterfaces(interfaces, i.resolver),
			})
		}
		// Add Usages tab if the type is used elsewhere
		if usages, _ := i.resolver.ResolveUsages(typeDef.Name()); len(usages) > 0 {
			tabs = append(tabs, components.Tab{
				Label:   "Usages",
				Content: adaptUsages(usages, i.resolver),
			})
		}
		directiveItems = adaptAppliedDirectives(typeDef.Directives(), i.resolver)
	case *gql.Scalar:
		// Add Usages tab if the type is used elsewhere
		if usages, _ := i.resolver.ResolveUsages(typeDef.Name()); len(usages) > 0 {
			tabs = append(tabs, components.Tab{
				Label:   "Usages",
				Content: adaptUsages(usages, i.resolver),
			})
		}
		directiveItems = adaptAppliedDirectives(typeDef.Directives(), i.resolver)
	case *gql.Interface:
		detailItems = append(detailItems, adaptFields(typeDef.Fields(), i.resolver)...)
		tabs = append(tabs, components.Tab{
			Label:   "Fields",
			Content: detailItems,
		})
		// Add Interfaces tab if the interface implements any interfaces
		if interfaces := typeDef.Interfaces(); len(interfaces) > 0 {
			tabs = append(tabs, components.Tab{
				Label:   "Interfaces",
				Content: adaptInterfaces(interfaces, i.resolver),
			})
		}
		// Add Usages tab if the type is used elsewhere
		if usages, _ := i.resolver.ResolveUsages(typeDef.Name()); len(usages) > 0 {
			tabs = append(tabs, components.Tab{
				Label:   "Usages",
				Content: adaptUsages(usages, i.resolver),
			})
		}
		directiveItems = adaptAppliedDirectives(typeDef.Directives(), i.resolver)
	case *gql.Union:
		detailItems = append(detailItems, adaptUnionTypes(typeDef.Types(), i.resolver)...)
		// Add Usages tab if the type is used elsewhere
		if usages, _ := i.resolver.ResolveUsages(typeDef.Name()); len(usages) > 0 {
			tabs = append(tabs, components.Tab{
				Label:   "Usages",
				Content: adaptUsages(usages, i.resolver),
			})
		}
		directiveItems = adaptAppliedDirectives(typeDef.Directives(), i.resolver)
	case *gql.Enum:
		detailItems = append(detailItems, adaptEnumValues(typeDef.Values())...)
		// Add Usages tab if the type is used elsewhere
		if usages, _ := i.resolver.ResolveUsages(typeDef.Name()); len(usages) > 0 {
			tabs = append(tabs, components.Tab{
				Label:   "Usages",
				Content: adaptUsages(usages, i.resolver),
			})
		}
		directiveItems = adaptAppliedDirectives(typeDef.Directives(), i.resolver)
	case *gql.InputObject:
		detailItems = append(detailItems, adaptFields(typeDef.Fields(), i.resolver)...)
		// Only add tabs if we have usages or directives
		if usages, _ := i.resolver.ResolveUsages(typeDef.Name()); len(usages) > 0 {
			tabs = append(tabs, components.Tab{
				Label:   "Fields",
				Content: detailItems,
			})
			tabs = append(tabs, components.Tab{
				Label:   "Usages",
				Content: adaptUsages(usages, i.resolver),
			})
		}
		directiveItems = adaptAppliedDirectives(typeDef.Directives(), i.resolver)
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
