# Task 07: Error Handling Strategy

**Priority:** High
**Status:** Not Started
**Estimated Effort:** Small
**Dependencies:** Task 01 (Type Resolver) - recommended

## Problem Statement

There are TODOs in the codebase regarding built-in type detection:

```go
// tui/adapters/items.go:326, 345
// TODO: Currently, this treats any error as a built-in type, but instead we should
// check for _known_ built in types and handle errors intelligently.
```

Currently, when `ResolveObjectTypeDef` returns an error, the code assumes it's a built-in scalar type. This is fragile because:

- Any resolution error is treated as a built-in type
- No distinction between "built-in type" and "actual error"
- Could hide real errors in type resolution
- Not clear what types are considered built-in

### Affected Files
- `gql/types.go` - Type resolution methods
- `tui/adapters/items.go` - Error handling in type resolution
- `gql/parse.go` - Schema parsing

## Proposed Solution

### 1. Define Built-In Types

Create explicit built-in type registry:

```go
// gql/builtins.go
package gql

// BuiltInTypes manages GraphQL built-in scalar types
type BuiltInTypes struct {
    types map[string]struct{}
}

// NewBuiltInTypes creates a registry of standard GraphQL built-in types
func NewBuiltInTypes() *BuiltInTypes {
    return &BuiltInTypes{
        types: map[string]struct{}{
            "String":  {},
            "Int":     {},
            "Float":   {},
            "Boolean": {},
            "ID":      {},
        },
    }
}

// IsBuiltIn returns whether a type name is a built-in scalar
func (b *BuiltInTypes) IsBuiltIn(typeName string) bool {
    _, exists := b.types[typeName]
    return exists
}

// Add registers a custom scalar as built-in (for custom scalar types)
func (b *BuiltInTypes) Add(typeName string) {
    b.types[typeName] = struct{}{}
}

// Remove unregisters a type name
func (b *BuiltInTypes) Remove(typeName string) {
    delete(b.types, typeName)
}

// All returns all registered built-in type names
func (b *BuiltInTypes) All() []string {
    names := make([]string, 0, len(b.types))
    for name := range b.types {
        names = append(names, name)
    }
    return names
}

// Standard built-in types instance
var StandardBuiltIns = NewBuiltInTypes()
```

### 2. Create Typed Errors

Define specific error types for different failure modes:

```go
// gql/errors.go
package gql

import "fmt"

// TypeNotFoundError indicates a type definition wasn't found in the schema
type TypeNotFoundError struct {
    TypeName string
}

func (e *TypeNotFoundError) Error() string {
    return fmt.Sprintf("type not found: %s", e.TypeName)
}

// IsTypeNotFound checks if an error is TypeNotFoundError
func IsTypeNotFound(err error) bool {
    _, ok := err.(*TypeNotFoundError)
    return ok
}

// BuiltInTypeError indicates the type is a built-in scalar
type BuiltInTypeError struct {
    TypeName string
}

func (e *BuiltInTypeError) Error() string {
    return fmt.Sprintf("built-in type: %s", e.TypeName)
}

// IsBuiltInType checks if an error is BuiltInTypeError
func IsBuiltInType(err error) bool {
    _, ok := err.(*BuiltInTypeError)
    return ok
}
```

### 3. Update Type Resolution

Improve error handling in type resolution:

```go
// gql/parse.go
// Update NamedToTypeDef to check built-ins explicitly
func (s *GraphQLSchema) NamedToTypeDef(typeName string) (TypeDef, error) {
    // Check if it's a built-in type first
    if StandardBuiltIns.IsBuiltIn(typeName) {
        return nil, &BuiltInTypeError{TypeName: typeName}
    }

    // Try to find in schema
    if obj, ok := s.Object[typeName]; ok {
        return obj, nil
    }
    if input, ok := s.Input[typeName]; ok {
        return input, nil
    }
    if enum, ok := s.Enum[typeName]; ok {
        return enum, nil
    }
    if scalar, ok := s.Scalar[typeName]; ok {
        return scalar, nil
    }
    if iface, ok := s.Interface[typeName]; ok {
        return iface, nil
    }
    if union, ok := s.Union[typeName]; ok {
        return union, nil
    }

    // Type not found
    return nil, &TypeNotFoundError{TypeName: typeName}
}
```

### 4. Update Adapter Error Handling

Handle errors intelligently in adapters:

```go
// tui/adapters/items.go
func newTypeDefItemFromField(field *gql.Field, resolver gql.TypeResolver) components.ListItem {
    resultType, err := resolver.ResolveFieldType(field)

    if err != nil {
        // Check if it's a built-in type
        if gql.IsBuiltInType(err) {
            return components.NewSimpleItem(
                field.TypeString(),
                components.WithTypeName(field.ObjectTypeName()),
            )
        }

        // Check if type wasn't found
        if gql.IsTypeNotFound(err) {
            // Log warning or handle missing type
            return components.NewSimpleItem(
                field.TypeString()+" (not found)",
                components.WithTypeName(field.ObjectTypeName()),
            )
        }

        // Unknown error - log and create simple item
        // Could add error logging here
        return components.NewSimpleItem(
            field.TypeString(),
            components.WithTypeName(field.ObjectTypeName()),
        )
    }

    // Successfully resolved
    return typeDefItem{
        title:    field.TypeString(),
        typeName: resultType.Name(),
        typeDef:  resultType,
        resolver: resolver,
    }
}

func newTypeDefItemFromArgument(argument *gql.Argument, resolver gql.TypeResolver) components.ListItem {
    resultType, err := resolver.ResolveArgumentType(argument)

    if err != nil {
        if gql.IsBuiltInType(err) {
            return components.NewSimpleItem(
                argument.TypeString(),
                components.WithTypeName(argument.ObjectTypeName()),
            )
        }

        if gql.IsTypeNotFound(err) {
            return components.NewSimpleItem(
                argument.TypeString()+" (not found)",
                components.WithTypeName(argument.ObjectTypeName()),
            )
        }

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
```

### 5. Support Custom Scalars

Allow registering custom scalars:

```go
// gql/parse.go
// Update ParseSchema to register custom scalars
func ParseSchema(schemaContent []byte) (GraphQLSchema, error) {
    // ... existing parsing ...

    schema := GraphQLSchema{
        Query:     queryFields,
        Mutation:  mutationFields,
        Object:    objects,
        Input:     inputs,
        Enum:      enums,
        Scalar:    scalars,
        Interface: interfaces,
        Union:     unions,
        Directive: directives,
    }

    // Register custom scalars as built-ins
    for name := range scalars {
        if !StandardBuiltIns.IsBuiltIn(name) {
            StandardBuiltIns.Add(name)
        }
    }

    return schema, nil
}
```

## Benefits

1. **Clarity**: Explicit handling of built-in types vs errors
2. **Correctness**: Won't hide real errors as built-in types
3. **Extensibility**: Easy to add custom scalars
4. **Debugging**: Clear error messages for missing types
5. **Maintainability**: Centralized built-in type management
6. **Type safety**: Type-specific errors

## Implementation Steps

1. Create `gql/builtins.go` with `BuiltInTypes`
2. Create `gql/errors.go` with typed errors
3. Add tests for built-in type registry and errors
4. Update `NamedToTypeDef` to return specific errors
5. Update adapter error handling to check error types
6. Update `ParseSchema` to register custom scalars
7. Remove TODOs from `items.go`
8. Run tests: `just test`
9. Update documentation

## Testing Strategy

```go
// gql/builtins_test.go
func TestBuiltInTypes_IsBuiltIn(t *testing.T) {
    builtIns := NewBuiltInTypes()

    assert.True(t, builtIns.IsBuiltIn("String"))
    assert.True(t, builtIns.IsBuiltIn("Int"))
    assert.False(t, builtIns.IsBuiltIn("CustomType"))
}

func TestBuiltInTypes_AddCustom(t *testing.T) {
    builtIns := NewBuiltInTypes()

    builtIns.Add("DateTime")
    assert.True(t, builtIns.IsBuiltIn("DateTime"))
}

// gql/errors_test.go
func TestTypeNotFoundError(t *testing.T) {
    err := &TypeNotFoundError{TypeName: "Missing"}

    assert.True(t, IsTypeNotFound(err))
    assert.False(t, IsBuiltInType(err))
    assert.Equal(t, "type not found: Missing", err.Error())
}

// gql/parse_test.go
func TestNamedToTypeDef_BuiltInType(t *testing.T) {
    schema := parseTestSchema(t)

    _, err := schema.NamedToTypeDef("String")

    assert.True(t, IsBuiltInType(err))
}

func TestNamedToTypeDef_TypeNotFound(t *testing.T) {
    schema := parseTestSchema(t)

    _, err := schema.NamedToTypeDef("NonExistent")

    assert.True(t, IsTypeNotFound(err))
}

// tui/adapters/items_test.go
func TestNewTypeDefItemFromField_BuiltInType(t *testing.T) {
    field := testFieldWithType("String")
    resolver := &mockResolver{}

    item := newTypeDefItemFromField(field, resolver)

    assert.Equal(t, "String", item.TypeName())
}

func TestNewTypeDefItemFromField_NotFound(t *testing.T) {
    field := testFieldWithType("Missing")
    resolver := &mockResolver{}

    item := newTypeDefItemFromField(field, resolver)

    assert.Contains(t, item.Title(), "not found")
}
```

## Potential Issues

- **Breaking changes**: Error return values change
- **Migration**: Need to update all error handling sites
- **Backward compatibility**: Ensure existing behavior maintained

## Future Enhancements

1. **Error reporting UI**: Show type resolution errors in UI
2. **Schema validation**: Validate schema for missing types
3. **Type suggestions**: Suggest similar type names for typos
4. **Custom scalar configuration**: Load custom scalars from config
5. **Error recovery**: Attempt to recover from schema errors

## Related Tasks

- **Task 01** (Type Resolver): Resolver can use better error handling
- **Task 08** (View Models): ViewModels can handle errors in presentation layer
