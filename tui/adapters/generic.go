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

// AdaptTypeDefs is a specialized version for types implementing gql.TypeDef
func AdaptTypeDefs[T gql.TypeDef](
	items []T,
	resolver gql.TypeResolver,
) []components.ListItem {
	return adaptSlice(items, resolver, func(t T, r gql.TypeResolver) components.ListItem {
		return newTypeDefItem(t, r)
	})
}

// AdaptFields converts Field slices to ListItems
func AdaptFields(fields []*gql.Field, resolver gql.TypeResolver) []components.ListItem {
	return adaptSlice(fields, resolver, newFieldItem)
}

// AdaptArguments converts Argument slices to ListItems
func AdaptArguments(args []*gql.Argument, resolver gql.TypeResolver) []components.ListItem {
	return adaptSlice(args, resolver, newArgumentItem)
}

// AdaptDirectives converts Directive slices to ListItems
func AdaptDirectives(directives []*gql.Directive, resolver gql.TypeResolver) []components.ListItem {
	return adaptSlice(directives, resolver, newDirectiveDefinitionItem)
}

// AdaptUnionTypes converts type name slices to ListItems
func AdaptUnionTypes(typeNames []string, resolver gql.TypeResolver) []components.ListItem {
	return adaptSlice(typeNames, resolver, newNamedItem)
}

// AdaptEnumValues converts EnumValue slices to ListItems (no resolver needed)
func AdaptEnumValues(values []*gql.EnumValue) []components.ListItem {
	return adaptSlice(values, nil, func(v *gql.EnumValue, _ gql.TypeResolver) components.ListItem {
		return components.NewSimpleItem(
			v.Name(),
			components.WithDescription(v.Description()),
		)
	})
}
