# Task 03: Type Registry

**Priority:** Medium
**Status:** Not Started
**Estimated Effort:** Small-Medium
**Dependencies:** None (complements Task 02)

## Problem Statement

The current GQL type system uses:
- String-based `gqlType` enum constants
- Hardcoded `availableGQLTypes` slice
- Switch statements to map types to loader functions

This approach:
- Isn't extensible (can't add custom categories without code changes)
- Duplicates type information across code
- Makes it hard to customize available types
- Couples type definition to loading logic

### Affected Files
- `tui/model.go` - Type constants and switch statements
- `tui/adapters/schema.go` - Loading methods

## Proposed Solution

Create a type registry that treats GQL type categories as first-class, configurable objects:

```go
// tui/registry/types.go
package registry

import "github.com/tonysyu/gqlxp/tui/components"

// GQLTypeCategory represents a category of GraphQL types that can be displayed
type GQLTypeCategory struct {
    // ID is the unique identifier for this category
    ID string

    // DisplayName is shown in the UI navbar
    DisplayName string

    // Description provides additional context (optional)
    Description string

    // LoadItems loads items for this category from a schema
    LoadItems func(schema SchemaProvider) []components.ListItem

    // Shortcut key for jumping to this type (optional)
    Shortcut string
}

// SchemaProvider is the interface categories use to load items
// This allows registry to be independent of concrete schema implementation
type SchemaProvider interface {
    GetQueryItems() []components.ListItem
    GetMutationItems() []components.ListItem
    GetObjectItems() []components.ListItem
    GetInputItems() []components.ListItem
    GetEnumItems() []components.ListItem
    GetScalarItems() []components.ListItem
    GetInterfaceItems() []components.ListItem
    GetUnionItems() []components.ListItem
    GetDirectiveItems() []components.ListItem
}

// TypeRegistry manages available GQL type categories
type TypeRegistry struct {
    categories []GQLTypeCategory
    current    int
}

// NewDefaultTypeRegistry creates a registry with standard GraphQL types
func NewDefaultTypeRegistry() *TypeRegistry {
    categories := []GQLTypeCategory{
        {
            ID:          "query",
            DisplayName: "Query",
            Description: "GraphQL Query fields",
            LoadItems:   func(s SchemaProvider) []components.ListItem { return s.GetQueryItems() },
            Shortcut:    "q",
        },
        {
            ID:          "mutation",
            DisplayName: "Mutation",
            Description: "GraphQL Mutation fields",
            LoadItems:   func(s SchemaProvider) []components.ListItem { return s.GetMutationItems() },
            Shortcut:    "m",
        },
        {
            ID:          "object",
            DisplayName: "Object",
            Description: "GraphQL Object types",
            LoadItems:   func(s SchemaProvider) []components.ListItem { return s.GetObjectItems() },
            Shortcut:    "o",
        },
        {
            ID:          "input",
            DisplayName: "Input",
            Description: "GraphQL Input types",
            LoadItems:   func(s SchemaProvider) []components.ListItem { return s.GetInputItems() },
            Shortcut:    "i",
        },
        {
            ID:          "enum",
            DisplayName: "Enum",
            Description: "GraphQL Enum types",
            LoadItems:   func(s SchemaProvider) []components.ListItem { return s.GetEnumItems() },
            Shortcut:    "e",
        },
        {
            ID:          "scalar",
            DisplayName: "Scalar",
            Description: "GraphQL Scalar types",
            LoadItems:   func(s SchemaProvider) []components.ListItem { return s.GetScalarItems() },
            Shortcut:    "s",
        },
        {
            ID:          "interface",
            DisplayName: "Interface",
            Description: "GraphQL Interface types",
            LoadItems:   func(s SchemaProvider) []components.ListItem { return s.GetInterfaceItems() },
            Shortcut:    "I",
        },
        {
            ID:          "union",
            DisplayName: "Union",
            Description: "GraphQL Union types",
            LoadItems:   func(s SchemaProvider) []components.ListItem { return s.GetUnionItems() },
            Shortcut:    "u",
        },
        {
            ID:          "directive",
            DisplayName: "Directive",
            Description: "GraphQL Directive types",
            LoadItems:   func(s SchemaProvider) []components.ListItem { return s.GetDirectiveItems() },
            Shortcut:    "d",
        },
    }

    return &TypeRegistry{
        categories: categories,
        current:    0,
    }
}

// Current returns the currently selected category
func (r *TypeRegistry) Current() GQLTypeCategory {
    return r.categories[r.current]
}

// SetCurrent sets the current category by ID
func (r *TypeRegistry) SetCurrent(id string) bool {
    for i, cat := range r.categories {
        if cat.ID == id {
            r.current = i
            return true
        }
    }
    return false
}

// Next moves to the next category (with wraparound)
func (r *TypeRegistry) Next() GQLTypeCategory {
    r.current = (r.current + 1) % len(r.categories)
    return r.categories[r.current]
}

// Previous moves to the previous category (with wraparound)
func (r *TypeRegistry) Previous() GQLTypeCategory {
    r.current = (r.current - 1 + len(r.categories)) % len(r.categories)
    return r.categories[r.current]
}

// All returns all categories
func (r *TypeRegistry) All() []GQLTypeCategory {
    return r.categories
}

// Add registers a new category
func (r *TypeRegistry) Add(category GQLTypeCategory) {
    r.categories = append(r.categories, category)
}

// Remove unregisters a category by ID
func (r *TypeRegistry) Remove(id string) bool {
    for i, cat := range r.categories {
        if cat.ID == id {
            r.categories = append(r.categories[:i], r.categories[i+1:]...)
            if r.current >= len(r.categories) {
                r.current = len(r.categories) - 1
            }
            return true
        }
    }
    return false
}

// FindByShortcut returns category with matching shortcut
func (r *TypeRegistry) FindByShortcut(shortcut string) (GQLTypeCategory, bool) {
    for _, cat := range r.categories {
        if cat.Shortcut == shortcut {
            return cat, true
        }
    }
    return GQLTypeCategory{}, false
}
```

### Update mainModel

Simplify type handling in `mainModel`:

```go
// tui/model.go
type mainModel struct {
    schema       adapters.SchemaView
    typeRegistry *registry.TypeRegistry
    // ... rest
}

func newModel(schema adapters.SchemaView) mainModel {
    m := mainModel{
        schema:       schema,
        typeRegistry: registry.NewDefaultTypeRegistry(),
        // ... rest
    }
    m.resetAndLoadMainPanel()
    return m
}

func (m *mainModel) loadMainPanel() {
    category := m.typeRegistry.Current()
    items := category.LoadItems(&m.schema)
    title := category.DisplayName + " Fields" // or use category title

    m.panelStack[0] = components.NewListPanel(items, title)
    m.stackPosition = 0
    m.updatePanelFocusStates()

    // Auto-open first item...
}

func (m *mainModel) renderGQLTypeNavbar() string {
    var tabs []string

    current := m.typeRegistry.Current()
    for _, category := range m.typeRegistry.All() {
        var style lipgloss.Style
        if category.ID == current.ID {
            style = m.Styles.ActiveTab
        } else {
            style = m.Styles.InactiveTab
        }
        tabs = append(tabs, style.Render(category.DisplayName))
    }

    navbar := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
    return m.Styles.Navbar.Render(navbar)
}
```

## Benefits

1. **Extensibility**: Easy to add custom categories without code changes
2. **Self-documenting**: Categories describe themselves with metadata
3. **Flexibility**: Users could customize which types to show
4. **DRY**: Type information in one place
5. **Testability**: Easy to test with custom registries
6. **Future features**: Enables shortcuts, filtering, grouping

## Implementation Steps

1. Create `tui/registry/types.go` package
2. Implement `GQLTypeCategory` and `TypeRegistry`
3. Add tests for registry operations
4. Update `mainModel` to use registry
5. Remove old `gqlType` constants and `availableGQLTypes`
6. Update switch statements to use registry
7. Run tests: `just test`
8. Update documentation

## Testing Strategy

```go
// tui/registry/types_test.go
func TestTypeRegistry_Navigation(t *testing.T) {
    reg := NewDefaultTypeRegistry()

    assert.Equal(t, "query", reg.Current().ID)

    reg.Next()
    assert.Equal(t, "mutation", reg.Current().ID)

    reg.Previous()
    assert.Equal(t, "query", reg.Current().ID)
}

func TestTypeRegistry_CustomCategory(t *testing.T) {
    reg := NewDefaultTypeRegistry()

    custom := GQLTypeCategory{
        ID:          "custom",
        DisplayName: "Custom",
        LoadItems: func(s SchemaProvider) []components.ListItem {
            return []components.ListItem{}
        },
    }

    reg.Add(custom)
    reg.SetCurrent("custom")

    assert.Equal(t, "custom", reg.Current().ID)
}

func TestTypeRegistry_Shortcuts(t *testing.T) {
    reg := NewDefaultTypeRegistry()

    cat, found := reg.FindByShortcut("q")
    assert.True(t, found)
    assert.Equal(t, "query", cat.ID)
}
```

## Potential Issues

- **Breaking changes**: Removes `gqlType` enum
- **Configuration complexity**: If exposing to users, need UI for customization

## Future Enhancements

1. **User configuration**: Load custom categories from config file
2. **Keyboard shortcuts**: Jump to type by shortcut key
3. **Type filtering**: Show/hide certain categories
4. **Type grouping**: Organize categories into groups
5. **Lazy loading**: Only load items when category selected
6. **Plugins**: Allow plugins to register custom categories

## Related Tasks

- **Task 02** (Navigation Manager): `TypeSelector` could use registry
- **Task 09** (Configuration Layer): Could load registry from config
- **Task 05** (Command Pattern): Commands could operate on registry
