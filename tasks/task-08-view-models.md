# Task 08: Separate View Models from Domain Models

**Priority:** Medium
**Status:** Not Started
**Estimated Effort:** Medium-Large
**Dependencies:** Task 01 (Type Resolver) - recommended

## Problem Statement

Currently, adapter items (`fieldItem`, `argumentItem`, `typeDefItem`) mix presentation and domain concerns:

- Formatting logic (signatures, details) embedded in items
- UI-specific methods (`Title()`, `Details()`) on domain wrappers
- Hard to change presentation without touching domain layer
- Testing requires understanding both UI and domain
- Can't easily support multiple presentation styles

### Affected Files
- `tui/adapters/items.go` - Adapter items mixing concerns
- `gql/types.go` - Domain types

## Proposed Solution

Introduce ViewModels that separate presentation logic from domain models:

```go
// tui/viewmodels/field.go
package viewmodels

import (
    "github.com/tonysyu/gqlxp/gql"
    "github.com/tonysyu/gqlxp/tui/components"
    "github.com/tonysyu/gqlxp/utils/text"
)

// FieldViewModel provides presentation logic for Fields
type FieldViewModel struct {
    field    *gql.Field
    resolver gql.TypeResolver
}

func NewFieldViewModel(field *gql.Field, resolver gql.TypeResolver) *FieldViewModel {
    return &FieldViewModel{
        field:    field,
        resolver: resolver,
    }
}

// Domain returns the underlying domain object
func (vm *FieldViewModel) Domain() *gql.Field {
    return vm.field
}

// DisplayTitle returns the title for list display
func (vm *FieldViewModel) DisplayTitle() string {
    return vm.field.Signature()
}

// DisplayDescription returns the description for list display
func (vm *FieldViewModel) DisplayDescription() string {
    return vm.field.Description()
}

// FilterValue returns the value to use for filtering
func (vm *FieldViewModel) FilterValue() string {
    return vm.field.Name()
}

// RefName returns the reference name for breadcrumbs
func (vm *FieldViewModel) RefName() string {
    return vm.field.Name()
}

// TypeName returns the type name for navigation
func (vm *FieldViewModel) TypeName() string {
    return vm.field.ObjectTypeName()
}

// Details returns formatted details for overlay display
func (vm *FieldViewModel) Details() string {
    return text.JoinParagraphs(
        text.H1(vm.RefName()),
        text.GqlCode(vm.field.FormatSignature(80)),
        vm.DisplayDescription(),
    )
}

// CanOpen returns whether this field can open a detail panel
func (vm *FieldViewModel) CanOpen() bool {
    return true
}

// OpenPanel creates a panel for this field
func (vm *FieldViewModel) OpenPanel() (components.Panel, bool) {
    // Get arguments as view models
    args := vm.field.Arguments()
    argViewModels := make([]*ArgumentViewModel, len(args))
    for i, arg := range args {
        argViewModels[i] = NewArgumentViewModel(arg, vm.resolver)
    }

    // Convert to list items
    items := make([]components.ListItem, len(argViewModels))
    for i, argVM := range argViewModels {
        items[i] = NewListItemFromViewModel(argVM)
    }

    panel := components.NewListPanel(items, vm.field.Name())
    panel.SetDescription(vm.DisplayDescription())

    // Set result type
    resultTypeVM := vm.ResultTypeViewModel()
    if resultTypeVM != nil {
        panel.SetObjectType(NewListItemFromViewModel(resultTypeVM))
    }

    return panel, true
}

// ResultTypeViewModel returns view model for the result type
func (vm *FieldViewModel) ResultTypeViewModel() ViewModel {
    resultType, err := vm.resolver.ResolveFieldType(vm.field)

    if err != nil {
        if gql.IsBuiltInType(err) {
            return &SimpleViewModel{
                title:    vm.field.TypeString(),
                typeName: vm.field.ObjectTypeName(),
            }
        }
        // Handle other errors
        return nil
    }

    return NewTypeDefViewModel(resultType, vm.resolver)
}
```

### ViewModel Interface

Define common interface for all view models:

```go
// tui/viewmodels/viewmodel.go
package viewmodels

import "github.com/tonysyu/gqlxp/tui/components"

// ViewModel provides presentation logic for domain objects
type ViewModel interface {
    // Display methods
    DisplayTitle() string
    DisplayDescription() string
    FilterValue() string
    RefName() string
    TypeName() string
    Details() string

    // Navigation
    CanOpen() bool
    OpenPanel() (components.Panel, bool)
}

// Ensure implementations satisfy interface
var _ ViewModel = (*FieldViewModel)(nil)
var _ ViewModel = (*ArgumentViewModel)(nil)
var _ ViewModel = (*TypeDefViewModel)(nil)
var _ ViewModel = (*DirectiveViewModel)(nil)
var _ ViewModel = (*SimpleViewModel)(nil)
```

### Additional ViewModels

```go
// tui/viewmodels/argument.go
type ArgumentViewModel struct {
    argument *gql.Argument
    resolver gql.TypeResolver
}

func NewArgumentViewModel(arg *gql.Argument, resolver gql.TypeResolver) *ArgumentViewModel {
    return &ArgumentViewModel{
        argument: arg,
        resolver: resolver,
    }
}

// ... implement ViewModel interface

// tui/viewmodels/typedef.go
type TypeDefViewModel struct {
    typeDef  gql.TypeDef
    resolver gql.TypeResolver
}

func NewTypeDefViewModel(typeDef gql.TypeDef, resolver gql.TypeResolver) *TypeDefViewModel {
    return &TypeDefViewModel{
        typeDef:  typeDef,
        resolver: resolver,
    }
}

// ... implement ViewModel interface

// tui/viewmodels/directive.go
type DirectiveViewModel struct {
    directive *gql.Directive
    resolver  gql.TypeResolver
}

func NewDirectiveViewModel(directive *gql.Directive, resolver gql.TypeResolver) *DirectiveViewModel {
    return &DirectiveViewModel{
        directive: directive,
        resolver:  resolver,
    }
}

// ... implement ViewModel interface

// tui/viewmodels/simple.go
type SimpleViewModel struct {
    title       string
    description string
    typeName    string
}

func NewSimpleViewModel(title string, opts ...SimpleOption) *SimpleViewModel {
    vm := &SimpleViewModel{title: title}
    for _, opt := range opts {
        opt(vm)
    }
    return vm
}

// ... implement ViewModel interface
```

### Adapter to ListItem

Create adapter from ViewModel to ListItem:

```go
// tui/viewmodels/listitem.go
package viewmodels

import (
    "github.com/charmbracelet/bubbles/list"
    "github.com/tonysyu/gqlxp/tui/components"
)

// ListItemAdapter adapts a ViewModel to components.ListItem
type ListItemAdapter struct {
    viewModel ViewModel
}

func NewListItemFromViewModel(vm ViewModel) components.ListItem {
    return &ListItemAdapter{viewModel: vm}
}

func (a *ListItemAdapter) Title() string {
    return a.viewModel.DisplayTitle()
}

func (a *ListItemAdapter) Description() string {
    return a.viewModel.DisplayDescription()
}

func (a *ListItemAdapter) FilterValue() string {
    return a.viewModel.FilterValue()
}

func (a *ListItemAdapter) RefName() string {
    return a.viewModel.RefName()
}

func (a *ListItemAdapter) TypeName() string {
    return a.viewModel.TypeName()
}

func (a *ListItemAdapter) Details() string {
    return a.viewModel.Details()
}

func (a *ListItemAdapter) OpenPanel() (components.Panel, bool) {
    if !a.viewModel.CanOpen() {
        return nil, false
    }
    return a.viewModel.OpenPanel()
}

// Ensure it implements required interfaces
var _ list.Item = (*ListItemAdapter)(nil)
var _ components.ListItem = (*ListItemAdapter)(nil)
```

### Update Adapters

Simplify adapters to create ViewModels:

```go
// tui/adapters/schema.go
func (sv *SchemaView) GetQueryItems() []components.ListItem {
    fields := gql.CollectAndSortMapValues(sv.schema.Query)
    return sv.fieldsToListItems(fields)
}

func (sv *SchemaView) fieldsToListItems(fields []*gql.Field) []components.ListItem {
    items := make([]components.ListItem, len(fields))
    for i, field := range fields {
        vm := viewmodels.NewFieldViewModel(field, sv.resolver)
        items[i] = viewmodels.NewListItemFromViewModel(vm)
    }
    return items
}

func (sv *SchemaView) GetObjectItems() []components.ListItem {
    objects := gql.CollectAndSortMapValues(sv.schema.Object)
    return sv.typeDefsToListItems(objects)
}

func (sv *SchemaView) typeDefsToListItems(typeDefs []gql.TypeDef) []components.ListItem {
    items := make([]components.ListItem, len(typeDefs))
    for i, typeDef := range typeDefs {
        vm := viewmodels.NewTypeDefViewModel(typeDef, sv.resolver)
        items[i] = viewmodels.NewListItemFromViewModel(vm)
    }
    return items
}
```

## Benefits

1. **Separation of concerns**: Presentation logic separate from domain
2. **Testability**: ViewModels easily testable without UI
3. **Flexibility**: Can support multiple presentation styles
4. **Reusability**: ViewModels reusable across different UI components
5. **Maintainability**: Changes to presentation don't affect domain
6. **Type safety**: Clear interfaces between layers

## Implementation Steps

1. Create `tui/viewmodels/` package
2. Define `ViewModel` interface
3. Implement `FieldViewModel`, `ArgumentViewModel`, etc.
4. Implement `ListItemAdapter`
5. Add tests for all view models
6. Update adapters to use view models
7. Remove old adapter items from `items.go`
8. Run tests: `just test`
9. Update documentation

## Testing Strategy

```go
// tui/viewmodels/field_test.go
func TestFieldViewModel_DisplayTitle(t *testing.T) {
    field := testField("search", "String")
    vm := NewFieldViewModel(field, &mockResolver{})

    title := vm.DisplayTitle()

    assert.Equal(t, "search: String", title)
}

func TestFieldViewModel_Details(t *testing.T) {
    field := testFieldWithDescription("search", "String", "Performs search")
    vm := NewFieldViewModel(field, &mockResolver{})

    details := vm.Details()

    assert.Contains(t, details, "search")
    assert.Contains(t, details, "Performs search")
}

func TestFieldViewModel_OpenPanel(t *testing.T) {
    field := testFieldWithArgs("search", "String", []gql.Argument{...})
    vm := NewFieldViewModel(field, &mockResolver{})

    panel, ok := vm.OpenPanel()

    assert.True(t, ok)
    assert.NotNil(t, panel)
}

// tui/viewmodels/listitem_test.go
func TestListItemAdapter_Title(t *testing.T) {
    vm := &SimpleViewModel{title: "Test"}
    item := NewListItemFromViewModel(vm)

    assert.Equal(t, "Test", item.Title())
}
```

## Potential Issues

- **Abstraction overhead**: Adds another layer
- **Migration effort**: Need to update all adapter code
- **Complexity**: More types to understand

## Future Enhancements

1. **Multiple presentation styles**: Compact, detailed, etc.
2. **Theming**: ViewModels could support themes
3. **Localization**: ViewModels could handle translations
4. **Caching**: Cache formatted strings in ViewModels
5. **Validation**: ViewModels could validate before display

## Related Tasks

- **Task 01** (Type Resolver): ViewModels use resolver for type resolution
- **Task 04** (Consolidate Adapters): ViewModels simplify adapter logic
- **Task 07** (Error Handling): ViewModels handle presentation of errors
