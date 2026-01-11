package adapters

import (
	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/tui/xplr/components"
)

// adaptSlice converts a slice of items to ListItems using a factory function
func adaptSlice[T any](
	items []T,
	resolver gql.TypeResolver,
	factory func(T, gql.TypeResolver) components.ListItem,
) []components.ListItem {
	result := make([]components.ListItem, 0, len(items))
	for _, item := range items {
		result = append(result, factory(item, resolver))
	}
	return result
}

// adaptTypeDefs is a specialized version for types implementing gql.TypeDef
func adaptTypeDefs[T gql.TypeDef](
	items []T,
	resolver gql.TypeResolver,
) []components.ListItem {
	return adaptSlice(items, resolver, func(t T, r gql.TypeResolver) components.ListItem {
		return newTypeDefItem(t, r)
	})
}

// adaptFields converts Field slices to ListItems
func adaptFields(fields []*gql.Field, resolver gql.TypeResolver) []components.ListItem {
	return adaptSlice(fields, resolver, newFieldItem)
}

// adaptArguments converts Argument slices to ListItems
func adaptArguments(args []*gql.Argument, resolver gql.TypeResolver) []components.ListItem {
	return adaptSlice(args, resolver, newArgumentItem)
}

// adaptDirectiveDefs converts DirectiveDef slices to ListItems (for schema directives)
func adaptDirectiveDefs(directives []*gql.DirectiveDef, resolver gql.TypeResolver) []components.ListItem {
	return adaptSlice(directives, resolver, newDirectiveDefItem)
}

// adaptAppliedDirectives converts AppliedDirective slices to ListItems (for directives on fields/types)
func adaptAppliedDirectives(directives []*gql.AppliedDirective, resolver gql.TypeResolver) []components.ListItem {
	return adaptSlice(directives, resolver, newAppliedDirectiveItem)
}

// adaptUnionTypes converts type name slices to ListItems
func adaptUnionTypes(typeNames []string, resolver gql.TypeResolver) []components.ListItem {
	return adaptSlice(typeNames, resolver, newNamedItem)
}

// adaptInterfaces converts interface name slices to ListItems
func adaptInterfaces(interfaceNames []string, resolver gql.TypeResolver) []components.ListItem {
	return adaptSlice(interfaceNames, resolver, newNamedItem)
}

// adaptEnumValues converts EnumValue slices to ListItems (no resolver needed)
func adaptEnumValues(values []*gql.EnumValue) []components.ListItem {
	return adaptSlice(values, nil, func(v *gql.EnumValue, _ gql.TypeResolver) components.ListItem {
		return components.NewSimpleItem(
			v.Name(),
			components.WithDescription(v.Description()),
		)
	})
}

// adaptUsages converts Usage slices to ListItems
func adaptUsages(usages []*gql.Usage, resolver gql.TypeResolver) []components.ListItem {
	return adaptSlice(usages, resolver, newUsageItem)
}
