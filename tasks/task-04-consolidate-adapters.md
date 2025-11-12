# Task 04: Consolidate Adapters with Generics

**Priority:** Medium
**Status:** Not Started
**Estimated Effort:** Small-Medium
**Dependencies:** Task 01 (Type Resolver) - recommended but not required

## Problem Statement

The `tui/adapters/items.go` file contains many repetitive adapter functions that follow the same pattern:

```go
func adaptObjectsToItems(objects []*gql.Object, schema *gql.GraphQLSchema) []components.ListItem {
    return adaptTypeDefsToItems(objects, schema)
}

func adaptInputObjectsToItems(inputs []*gql.InputObject, schema *gql.GraphQLSchema) []components.ListItem {
    return adaptTypeDefsToItems(inputs, schema)
}

func adaptEnumsToItems(enums []*gql.Enum, schema *gql.GraphQLSchema) []components.ListItem {
    return adaptTypeDefsToItems(enums, schema)
}

// ... and more
```

This creates:
- **Code duplication**: Same pattern repeated for each type
- **Maintenance burden**: Changes require updates to multiple functions
- **Boilerplate**: Each function is just a thin wrapper

### Affected Files
- `tui/adapters/items.go` - All `adapt*ToItems` functions

## Proposed Solution

Use Go generics to consolidate repetitive adapter functions:

```go
// tui/adapters/generic.go
package adapters

import (
    "github.com/tonysyu/gqlxp/gql"
    "github.com/tonysyu/gqlxp/tui/components"
)

// AdaptSlice converts a slice of items to ListItems using a factory function
func AdaptSlice[T any](
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
    return AdaptSlice(items, resolver, func(t T, r gql.TypeResolver) components.ListItem {
        return newTypeDefItem(t, r)
    })
}

// AdaptFields converts Field slices to ListItems
func AdaptFields(fields []*gql.Field, resolver gql.TypeResolver) []components.ListItem {
    return AdaptSlice(fields, resolver, newFieldItem)
}

// AdaptArguments converts Argument slices to ListItems
func AdaptArguments(args []*gql.Argument, resolver gql.TypeResolver) []components.ListItem {
    return AdaptSlice(args, resolver, newArgumentItem)
}

// AdaptDirectives converts Directive slices to ListItems
func AdaptDirectives(directives []*gql.Directive, resolver gql.TypeResolver) []components.ListItem {
    return AdaptSlice(directives, resolver, newDirectiveDefinitionItem)
}

// AdaptUnionTypes converts type name slices to ListItems
func AdaptUnionTypes(typeNames []string, resolver gql.TypeResolver) []components.ListItem {
    return AdaptSlice(typeNames, resolver, newNamedItem)
}

// AdaptEnumValues converts EnumValue slices to ListItems (no resolver needed)
func AdaptEnumValues(values []*gql.EnumValue) []components.ListItem {
    return AdaptSlice(values, nil, func(v *gql.EnumValue, _ gql.TypeResolver) components.ListItem {
        return components.NewSimpleItem(
            v.Name(),
            components.WithDescription(v.Description()),
        )
    })
}
```

### Update items.go

Simplify `items.go` to use generic adapters:

```go
// tui/adapters/items.go

// Remove these functions (now handled by AdaptTypeDefs):
// - adaptObjectsToItems
// - adaptInputObjectsToItems
// - adaptEnumsToItems
// - adaptScalarsToItems
// - adaptInterfacesToItems
// - adaptUnionsToItems

// Remove these (replaced by specific adapters):
// - adaptFieldsToItems (use AdaptFields)
// - adaptArgumentsToItems (use AdaptArguments)
// - adaptDirectivesToItems (use AdaptDirectives)
// - adaptUnionTypesToItems (use AdaptUnionTypes)
// - adaptEnumValuesToItems (use AdaptEnumValues)

// Keep only the item type definitions and factory functions:
// - fieldItem, newFieldItem
// - argumentItem, newArgumentItem
// - typeDefItem, newTypeDefItem
// - directiveItem, newDirectiveDefinitionItem
// - newNamedItem
```

### Update schema.go

Use consolidated adapters in `SchemaView`:

```go
// tui/adapters/schema.go
func (p *SchemaView) GetQueryItems() []components.ListItem {
    return AdaptFields(gql.CollectAndSortMapValues(p.schema.Query), p.resolver)
}

func (p *SchemaView) GetMutationItems() []components.ListItem {
    return AdaptFields(gql.CollectAndSortMapValues(p.schema.Mutation), p.resolver)
}

func (p *SchemaView) GetObjectItems() []components.ListItem {
    return AdaptTypeDefs(gql.CollectAndSortMapValues(p.schema.Object), p.resolver)
}

func (p *SchemaView) GetInputItems() []components.ListItem {
    return AdaptTypeDefs(gql.CollectAndSortMapValues(p.schema.Input), p.resolver)
}

func (p *SchemaView) GetEnumItems() []components.ListItem {
    return AdaptTypeDefs(gql.CollectAndSortMapValues(p.schema.Enum), p.resolver)
}

func (p *SchemaView) GetScalarItems() []components.ListItem {
    return AdaptTypeDefs(gql.CollectAndSortMapValues(p.schema.Scalar), p.resolver)
}

func (p *SchemaView) GetInterfaceItems() []components.ListItem {
    return AdaptTypeDefs(gql.CollectAndSortMapValues(p.schema.Interface), p.resolver)
}

func (p *SchemaView) GetUnionItems() []components.ListItem {
    return AdaptTypeDefs(gql.CollectAndSortMapValues(p.schema.Union), p.resolver)
}

func (p *SchemaView) GetDirectiveItems() []components.ListItem {
    return AdaptDirectives(gql.CollectAndSortMapValues(p.schema.Directive), p.resolver)
}
```

### Update item implementations

Update calls within item types:

```go
// In fieldItem.OpenPanel():
func (i fieldItem) OpenPanel() (components.Panel, bool) {
    argumentItems := AdaptArguments(i.gqlField.Arguments(), i.resolver)
    // ... rest
}

// In typeDefItem.OpenPanel():
func (i typeDefItem) OpenPanel() (components.Panel, bool) {
    var detailItems []components.ListItem

    switch typeDef := (i.typeDef).(type) {
    case *gql.Object:
        detailItems = AdaptFields(typeDef.Fields(), i.resolver)
    case *gql.Interface:
        detailItems = AdaptFields(typeDef.Fields(), i.resolver)
    case *gql.Union:
        detailItems = AdaptUnionTypes(typeDef.Types(), i.resolver)
    case *gql.Enum:
        detailItems = AdaptEnumValues(typeDef.Values())
    case *gql.InputObject:
        detailItems = AdaptFields(typeDef.Fields(), i.resolver)
    }
    // ... rest
}
```

## Benefits

1. **DRY**: Eliminates ~50-100 lines of boilerplate code
2. **Type safety**: Generics provide compile-time type checking
3. **Maintainability**: Changes to adaptation logic in one place
4. **Clarity**: Intent is clearer with descriptive function names
5. **Extensibility**: Easy to add new adapter types

## Implementation Steps

1. Create `tui/adapters/generic.go`
2. Implement `AdaptSlice` and specialized adapters
3. Add tests for generic adapters
4. Update `schema.go` to use new adapters
5. Update item implementations to use new adapters
6. Remove old adapter functions from `items.go`
7. Run tests: `just test`
8. Verify no performance regression

## Testing Strategy

```go
// tui/adapters/generic_test.go
func TestAdaptSlice(t *testing.T) {
    fields := []*gql.Field{
        testField("field1"),
        testField("field2"),
    }
    resolver := &mockResolver{}

    items := AdaptSlice(fields, resolver, newFieldItem)

    assert.Equal(t, 2, len(items))
    assert.Equal(t, "field1", items[0].FilterValue())
    assert.Equal(t, "field2", items[1].FilterValue())
}

func TestAdaptTypeDefs(t *testing.T) {
    objects := []*gql.Object{
        testObject("User"),
        testObject("Post"),
    }
    resolver := &mockResolver{}

    items := AdaptTypeDefs(objects, resolver)

    assert.Equal(t, 2, len(items))
    assert.Equal(t, "User", items[0].FilterValue())
}

func TestAdaptEnumValues(t *testing.T) {
    values := []*gql.EnumValue{
        testEnumValue("ACTIVE"),
        testEnumValue("INACTIVE"),
    }

    items := AdaptEnumValues(values)

    assert.Equal(t, 2, len(items))
    assert.Equal(t, "ACTIVE", items[0].FilterValue())
}
```

## Potential Issues

- **Go version**: Requires Go 1.18+ for generics
- **Complexity**: Generic code can be harder to understand initially
- **Performance**: Negligible - compiler optimizes generics well

## Future Enhancements

1. **Parallel adaptation**: Use goroutines for large slices
2. **Filtered adaptation**: Add predicate parameter to filter items
3. **Mapped adaptation**: Transform items during adaptation
4. **Batched adaptation**: Process items in batches for very large schemas

## Related Tasks

- **Task 01** (Type Resolver): Generic adapters work better with resolver interface
- **Task 08** (View Models): Could create generic ViewModelAdapter
