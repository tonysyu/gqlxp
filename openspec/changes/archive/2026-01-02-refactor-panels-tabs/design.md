# Design: Tab-Based Panel Navigation

## Context
Currently, panels special-case the display of a GraphQL field's result type using dedicated fields (`resultType`, `focusOnResultType`) and custom cursor navigation logic. This pattern is difficult to extend for additional relationships like:
- Implementations (for interfaces)
- Implemented interfaces (for objects)
- Back-references (types that reference this type)
- Union member types

The bubbletea tabs example provides a proven pattern for tabbed interfaces that we can adapt.

## Goals / Non-Goals

**Goals:**
- Replace special-case result type handling with general-purpose tab system
- Enable future addition of new relationship tabs without architectural changes
- Use Shift-H/Shift-L for tab navigation (matches the user's preference)
- Maintain current panel auto-opening behavior for selected items

**Non-Goals:**
- Implementing new relationship types (Implementations, Interfaces, etc.) in this change
- Changing the visual styling beyond what's needed for tabs
- Adding tab navigation at the top-level model (tabs are panel-specific)

## Decisions

### Decision: Tab Structure
Each panel will maintain:
```go
type Tab struct {
    Label   string
    Content []list.Item
}

type Panel struct {
    tabs          []Tab              // Tabs with labels and content
    activeTab     int                // Currently selected tab index
    // ... existing fields for list, styling, etc.
}
```

**Rationale**:
- Type-safe: Label and content are guaranteed to stay together
- Cannot desynchronize: Impossible to have mismatched arrays
- Idiomatic Go: Encapsulates related data in a struct
- Easy to extend: Can add fields to Tab (e.g., `Enabled bool`, `Icon string`) without refactoring
- Maintainable: Adding/removing tabs is a single operation

**Alternatives considered**:
- Parallel arrays (tabs []string, tabContent [][]list.Item): Risk of desynchronization, harder to extend
- Single content array with filtering: Would break list state when switching tabs

### Decision: Keybindings
- **Shift-H**: Previous tab
- **Shift-L**: Next tab
- No wrapping (staying at first/last tab when at boundaries)

**Rationale**:
- User specifically requested Shift-H/Shift-L
- Matches vim-style navigation (h=left, l=right)
- No wrapping prevents accidental tab jumps

**Alternatives considered**:
- Tab/Shift-Tab: Already used for panel navigation
- Ctrl-H/Ctrl-L: Less intuitive, Shift is easier to discover

### Decision: Single Tab Handling
When a panel has only one "tab" (e.g., a type with only fields, no result type), the tab bar will be hidden and only the content displayed.

**Rationale**:
- Cleaner UI when tabs aren't needed
- Reduces visual clutter
- Existing behavior for types without result types is maintained

**Alternatives considered**:
- Always show tabs for consistency: Creates unnecessary UI noise

### Decision: Tab Content Updates
When switching tabs:
1. Update `activeTab` index
2. Call `ListModel.SetItems(tabs[activeTab].Content)`
3. Reset selection to first item
4. Trigger auto-open for first item (maintains current behavior)

**Rationale**:
- Reusing existing list view logic keeps the implementation simple
- Auto-opening on tab switch provides immediate context
- Resetting to first item is intuitive (similar to navigating to a new panel)

### Decision: Migration from SetObjectType
Replace:
```go
panel.SetObjectType(resultTypeItem)
```

With:
```go
panel.SetTabs([]Tab{
    {Label: "Type", Content: []list.Item{resultTypeItem}},
    {Label: "Inputs", Content: argumentItems},
})
```

For types with no result type (e.g., directives), use single tab or no tabs:
```go
panel.SetTabs([]Tab{
    {Label: "Inputs", Content: argumentItems},
})
// Or just use existing SetItems() with no tabs
```

**Rationale**:
- Clear, explicit API for setting tab data
- Type-safe: Cannot pass mismatched labels and content
- Allows different items per adapter without conditional logic
- Backward compatible for single-content panels (via SetItems)

## Risks / Trade-offs

### Risk: Increased complexity for simple cases
**Impact**: Types with only fields now need tab handling code even if tabs aren't shown.
**Mitigation**: Make `SetTabs()` optional. If not called, panel works as before with just `SetItems()`.

### Risk: Tab navigation conflicts with other keybindings
**Impact**: Shift-H/Shift-L might conflict with future features.
**Mitigation**: Document in help system. Shift+key combinations are currently underutilized in the app.

### Trade-off: Memory overhead for tab content
**Impact**: Storing duplicate item arrays for each tab uses more memory.
**Mitigation**: In practice, panels have <100 items typically. Memory impact is negligible for TUI usage.

## Migration Plan

### Phase 1: Add tab infrastructure (non-breaking)
1. Add tab fields to Panel struct
2. Add SetTabs() method (coexists with SetObjectType)
3. Update View() to check for tabs and render accordingly

### Phase 2: Update adapters
1. Migrate fieldItem.OpenPanel() to use SetTabs()
2. Migrate argumentItem.OpenPanel() to use SetTabs()
3. Keep typeDefItem.OpenPanel() using SetItems() (no tabs needed)

### Phase 3: Remove legacy code
1. Remove SetObjectType() method
2. Remove resultType, focusOnResultType fields
3. Remove result type navigation logic from Update()

### Rollback Plan
If issues arise:
1. Revert adapter changes (Phase 2)
2. Keep tab infrastructure but unused
3. Investigate and fix before proceeding

## Open Questions
None - design is ready for implementation.
