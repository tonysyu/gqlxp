# Task 01: Type Resolver Service

**Priority:** High
**Status:** Not Started
**Estimated Effort:** Medium
**Dependencies:** None

## Problem Statement

Currently, every adapter item (`fieldItem`, `argumentItem`, `typeDefItem`, `directiveItem`) holds a direct reference to `*gql.GraphQLSchema`. This creates:

- **Tight coupling**: UI layer directly depends on entire schema structure
- **Memory overhead**: Each item duplicates the schema pointer
- **Testing difficulty**: Hard to mock type resolution behavior
- **Inflexibility**: Can't add caching or alternate resolution strategies

### Affected Files
- `tui/adapters/items.go` - All adapter item types
- `tui/adapters/schema.go` - Schema view methods

## Proposed Solution

Introduce a `TypeResolver` interface that encapsulates type resolution logic:

```go
// gql/resolver.go
package gql

// TypeResolver provides methods for resolving GraphQL type definitions
type TypeResolver interface {
    // ResolveType resolves a type name to its definition
    ResolveType(typeName string) (TypeDef, error)

    // ResolveFieldType resolves a field's result type
    ResolveFieldType(field *Field) (TypeDef, error)

    // ResolveArgumentType resolves an argument's input type
    ResolveArgumentType(arg *Argument) (TypeDef, error)
}

// SchemaResolver implements TypeResolver using a GraphQLSchema
type SchemaResolver struct {
    schema *GraphQLSchema
}

func NewSchemaResolver(schema *GraphQLSchema) *SchemaResolver {
    return &SchemaResolver{schema: schema}
}

func (r *SchemaResolver) ResolveType(typeName string) (TypeDef, error) {
    return r.schema.NamedToTypeDef(typeName)
}

func (r *SchemaResolver) ResolveFieldType(field *Field) (TypeDef, error) {
    return field.ResolveObjectTypeDef(r.schema)
}

func (r *SchemaResolver) ResolveArgumentType(arg *Argument) (TypeDef, error) {
    return arg.ResolveObjectTypeDef(r.schema)
}
```

### Update Adapter Items

Change all adapter items to use `TypeResolver` instead of `*GraphQLSchema`:

```go
// tui/adapters/items.go
type fieldItem struct {
    gqlField  *gql.Field
    resolver  gql.TypeResolver  // Changed from schema *gql.GraphQLSchema
    fieldName string
}

func newFieldItem(gqlField *gql.Field, resolver gql.TypeResolver) components.ListItem {
    return fieldItem{
        gqlField:  gqlField,
        resolver:  resolver,
        fieldName: gqlField.Name(),
    }
}

func (i fieldItem) OpenPanel() (components.Panel, bool) {
    argumentItems := adaptArgumentsToItems(i.gqlField.Arguments(), i.resolver)

    panel := components.NewListPanel(argumentItems, i.fieldName)
    panel.SetDescription(i.Description())
    panel.SetObjectType(newTypeDefItemFromField(i.gqlField, i.resolver))

    return panel, true
}
```

### Update SchemaView

Update `SchemaView` to provide resolver:

```go
// tui/adapters/schema.go
type SchemaView struct {
    schema   gql.GraphQLSchema
    resolver gql.TypeResolver
}

func NewSchemaView(schema gql.GraphQLSchema) SchemaView {
    return SchemaView{
        schema:   schema,
        resolver: gql.NewSchemaResolver(&schema),
    }
}

func (p *SchemaView) GetQueryItems() []components.ListItem {
    return adaptFieldsToItems(gql.CollectAndSortMapValues(p.schema.Query), p.resolver)
}
```

## Benefits

1. **Decoupling**: UI layer depends on interface, not concrete schema
2. **Testability**: Easy to create mock resolvers for testing
3. **Extensibility**: Can add caching, lazy loading, or remote resolution
4. **Memory efficiency**: Interface reference is smaller than schema pointer
5. **Cleaner abstractions**: Type resolution is explicit and documented

## Implementation Steps

1. Create `gql/resolver.go` with `TypeResolver` interface and `SchemaResolver` implementation
2. Add tests for `SchemaResolver` in `gql/resolver_test.go`
3. Update all adapter item types to use `TypeResolver` instead of `*GraphQLSchema`
4. Update all adapter factory functions to accept `TypeResolver` parameter
5. Update `SchemaView` to create and provide `TypeResolver`
6. Run tests to ensure no regressions: `just test`
7. Update documentation if needed

## Testing Strategy

```go
// gql/resolver_test.go
func TestSchemaResolver_ResolveType(t *testing.T) {
    schema := parseTestSchema(t)
    resolver := NewSchemaResolver(&schema)

    typeDef, err := resolver.ResolveType("User")
    assert.NoError(t, err)
    assert.Equal(t, "User", typeDef.Name())
}

// tui/adapters/items_test.go - with mock resolver
type mockResolver struct {
    types map[string]gql.TypeDef
}

func (m *mockResolver) ResolveType(name string) (gql.TypeDef, error) {
    if td, ok := m.types[name]; ok {
        return td, nil
    }
    return nil, fmt.Errorf("type not found: %s", name)
}
```

## Potential Issues

- **Breaking change**: All adapter creation code needs updating
- **Migration effort**: Need to update all call sites
- **Performance**: Minimal - interface dispatch is negligible

## Future Enhancements

After implementing basic resolver:

1. **Caching resolver**: Wrap `SchemaResolver` with caching layer
2. **Lazy resolver**: Load types on-demand for large schemas
3. **Remote resolver**: Fetch types from GraphQL introspection endpoint
4. **Composite resolver**: Chain multiple resolvers (local + remote)

## Related Tasks

- **Task 07** (Error Handling): Complements resolver with built-in type detection
- **Task 04** (Consolidate Adapters): Can leverage resolver for generic adapters
- **Task 08** (View Models): ViewModels can use resolver for presentation logic
