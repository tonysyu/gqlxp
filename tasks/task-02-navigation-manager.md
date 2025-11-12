# Task 02: Navigation State Manager

**Priority:** High
**Status:** Completed
**Estimated Effort:** Medium-Large
**Dependencies:** None

## Problem Statement

The `mainModel` in `tui/model.go` (419 lines) handles too many responsibilities:

- Panel stack management
- Stack position tracking
- Breadcrumb management
- GQL type selection and switching
- Panel focus state
- Sizing logic
- Rendering

This violates Single Responsibility Principle and makes the code:
- Hard to test individual behaviors
- Difficult to reason about state changes
- Challenging to add new navigation features (history, bookmarks, etc.)

### Affected Files
- `tui/model.go` - Main model with mixed responsibilities
- `tui/breadcrumbs.go` - Breadcrumb management
- Future navigation features

## Proposed Solution

Extract navigation concerns into dedicated, focused components:

### 1. PanelStack

Manages the stack of panels and current position:

```go
// tui/navigation/stack.go
package navigation

import "github.com/tonysyu/gqlxp/tui/components"

// PanelStack manages a stack of panels with navigation
type PanelStack struct {
    panels   []components.Panel
    position int
}

func NewPanelStack(initialCapacity int) *PanelStack {
    return &PanelStack{
        panels:   make([]components.Panel, initialCapacity),
        position: 0,
    }
}

// Current returns the currently focused panel
func (s *PanelStack) Current() components.Panel {
    if s.position >= 0 && s.position < len(s.panels) {
        return s.panels[s.position]
    }
    return nil
}

// Next returns the panel after the current position (right panel)
func (s *PanelStack) Next() components.Panel {
    nextPos := s.position + 1
    if nextPos < len(s.panels) {
        return s.panels[nextPos]
    }
    return nil
}

// MoveForward advances position if possible, returns success
func (s *PanelStack) MoveForward() bool {
    if s.position+1 < len(s.panels) {
        s.position++
        return true
    }
    return false
}

// MoveBackward moves position back if possible, returns success
func (s *PanelStack) MoveBackward() bool {
    if s.position > 0 {
        s.position--
        return true
    }
    return false
}

// Push adds a panel after current position, truncating rest
func (s *PanelStack) Push(panel components.Panel) {
    s.panels = s.panels[:s.position+1]
    s.panels = append(s.panels, panel)
}

// Replace replaces all panels with new set
func (s *PanelStack) Replace(panels []components.Panel) {
    s.panels = panels
    s.position = 0
}

// Position returns current position
func (s *PanelStack) Position() int {
    return s.position
}

// Len returns number of panels
func (s *PanelStack) Len() int {
    return len(s.panels)
}

// All returns all panels (for iteration)
func (s *PanelStack) All() []components.Panel {
    return s.panels
}
```

### 2. TypeSelector

Manages GQL type selection and cycling:

```go
// tui/navigation/type_selector.go
package navigation

type GQLType string

const (
    QueryType     GQLType = "Query"
    MutationType  GQLType = "Mutation"
    ObjectType    GQLType = "Object"
    InputType     GQLType = "Input"
    EnumType      GQLType = "Enum"
    ScalarType    GQLType = "Scalar"
    InterfaceType GQLType = "Interface"
    UnionType     GQLType = "Union"
    DirectiveType GQLType = "Directive"
)

// TypeSelector manages selection among available GQL types
type TypeSelector struct {
    types    []GQLType
    selected GQLType
}

func NewTypeSelector() *TypeSelector {
    types := []GQLType{
        QueryType, MutationType, ObjectType, InputType,
        EnumType, ScalarType, InterfaceType, UnionType, DirectiveType,
    }
    return &TypeSelector{
        types:    types,
        selected: QueryType,
    }
}

// Current returns currently selected type
func (ts *TypeSelector) Current() GQLType {
    return ts.selected
}

// Set changes selected type
func (ts *TypeSelector) Set(gqlType GQLType) {
    ts.selected = gqlType
}

// Next cycles to next type (with wraparound)
func (ts *TypeSelector) Next() GQLType {
    idx := ts.currentIndex()
    nextIdx := (idx + 1) % len(ts.types)
    ts.selected = ts.types[nextIdx]
    return ts.selected
}

// Previous cycles to previous type (with wraparound)
func (ts *TypeSelector) Previous() GQLType {
    idx := ts.currentIndex()
    prevIdx := (idx - 1 + len(ts.types)) % len(ts.types)
    ts.selected = ts.types[prevIdx]
    return ts.selected
}

// All returns all available types
func (ts *TypeSelector) All() []GQLType {
    return ts.types
}

func (ts *TypeSelector) currentIndex() int {
    for i, t := range ts.types {
        if t == ts.selected {
            return i
        }
    }
    return 0
}
```

### 3. NavigationManager

Coordinates all navigation state:

```go
// tui/navigation/manager.go
package navigation

import "github.com/tonysyu/gqlxp/tui/components"

// NavigationManager coordinates panel stack, breadcrumbs, and type selection
type NavigationManager struct {
    stack        *PanelStack
    breadcrumbs  BreadcrumbTracker // interface to existing breadcrumbsModel
    typeSelector *TypeSelector
}

func NewNavigationManager(visiblePanels int, breadcrumbs BreadcrumbTracker) *NavigationManager {
    return &NavigationManager{
        stack:        NewPanelStack(visiblePanels),
        breadcrumbs:  breadcrumbs,
        typeSelector: NewTypeSelector(),
    }
}

// NavigateForward moves forward in panel stack
func (nm *NavigationManager) NavigateForward(addToBreadcrumbs func()) bool {
    if nm.stack.MoveForward() {
        addToBreadcrumbs()
        return true
    }
    return false
}

// NavigateBackward moves backward in panel stack
func (nm *NavigationManager) NavigateBackward() bool {
    if nm.stack.MoveBackward() {
        nm.breadcrumbs.Pop()
        return true
    }
    return false
}

// OpenPanel pushes new panel onto stack
func (nm *NavigationManager) OpenPanel(panel components.Panel) {
    nm.stack.Push(panel)
}

// SwitchType changes selected GQL type and resets navigation
func (nm *NavigationManager) SwitchType(gqlType GQLType) {
    nm.typeSelector.Set(gqlType)
    nm.breadcrumbs.Reset()
}

// CycleTypeForward moves to next GQL type
func (nm *NavigationManager) CycleTypeForward() GQLType {
    return nm.typeSelector.Next()
}

// CycleTypeBackward moves to previous GQL type
func (nm *NavigationManager) CycleTypeBackward() GQLType {
    return nm.typeSelector.Previous()
}

// Stack returns the panel stack (for rendering)
func (nm *NavigationManager) Stack() *PanelStack {
    return nm.stack
}

// CurrentType returns currently selected GQL type
func (nm *NavigationManager) CurrentType() GQLType {
    return nm.typeSelector.Current()
}

// AllTypes returns all available GQL types
func (nm *NavigationManager) AllTypes() []GQLType {
    return nm.typeSelector.All()
}
```

### 4. Update mainModel

Simplify `mainModel` to use `NavigationManager`:

```go
// tui/model.go
type mainModel struct {
    schema     adapters.SchemaView
    navigation *navigation.NavigationManager
    overlay    overlayModel
    // ... styling, sizing, etc.
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // ... overlay handling ...

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, m.keymap.NextPanel):
            if m.navigation.NavigateForward(m.addCurrentItemToBreadcrumbs) {
                m.updatePanelFocusStates()
                return m, m.openSelectedItemInCurrentPanel()
            }
        case key.Matches(msg, m.keymap.PrevPanel):
            if m.navigation.NavigateBackward() {
                m.updatePanelFocusStates()
            }
        case key.Matches(msg, m.keymap.ToggleGQLType):
            m.navigation.CycleTypeForward()
            m.resetAndLoadMainPanel()
        // ... etc
        }
    case components.OpenPanelMsg:
        m.navigation.OpenPanel(msg.Panel)
        m.sizePanels()
    }
    // ...
}
```

## Benefits

1. **Separation of concerns**: Each component has single, clear responsibility
2. **Testability**: Can test navigation logic independent of UI
3. **Maintainability**: Easier to understand and modify navigation behavior
4. **Extensibility**: Easy to add features like:
   - Navigation history
   - Bookmarks
   - Multiple stacks
   - Undo/redo
5. **Reduced complexity**: `mainModel` focuses on orchestration, not mechanics

## Implementation Steps

1. Create `tui/navigation/` package directory
2. Implement `PanelStack` with tests
3. Implement `TypeSelector` with tests
4. Define `BreadcrumbTracker` interface for existing breadcrumbs
5. Implement `NavigationManager` with tests
6. Update `mainModel` to use `NavigationManager`
7. Refactor navigation-related methods in `mainModel`
8. Run full test suite: `just test`
9. Update architecture documentation

## Testing Strategy

```go
// tui/navigation/stack_test.go
func TestPanelStack_MoveForward(t *testing.T) {
    stack := NewPanelStack(2)
    stack.panels = []components.Panel{
        components.NewStringPanel("1"),
        components.NewStringPanel("2"),
    }

    moved := stack.MoveForward()
    assert.True(t, moved)
    assert.Equal(t, 1, stack.Position())

    moved = stack.MoveForward()
    assert.False(t, moved) // Can't move past end
}

// tui/navigation/type_selector_test.go
func TestTypeSelector_Next(t *testing.T) {
    ts := NewTypeSelector()
    assert.Equal(t, QueryType, ts.Current())

    ts.Next()
    assert.Equal(t, MutationType, ts.Current())

    // Test wraparound
    for i := 0; i < 10; i++ {
        ts.Next()
    }
    assert.Equal(t, MutationType, ts.Current()) // Wrapped around
}
```

## Potential Issues

- **Migration complexity**: Requires careful refactoring of `mainModel`
- **Breaking changes**: May affect how panels are created/managed
- **Test updates**: Existing tests may need updates

## Future Enhancements

1. **History stack**: Track navigation history for back/forward
2. **State persistence**: Save/restore navigation state
3. **Multiple workspaces**: Manage multiple independent panel stacks
4. **Navigation shortcuts**: Jump to specific panels/types

## Related Tasks

- **Task 06** (Panel Lifecycle): Can integrate with navigation manager
- **Task 05** (Command Pattern): Commands could operate on navigation manager
- **Task 03** (Type Registry): TypeSelector could use registry instead of hardcoded list
