# Codebase Simplification TODO

## Priority Simplifications

### 3. Simplify Keymap Construction
**Issue**: Uses reflection to build `globalKeyBinds` from keymap struct (`tui/xplr/model.go:120-125`)

```go
// Current approach (lines 120-125)
v := reflect.ValueOf(m.keymap)
m.globalKeyBinds = make([]key.Binding, v.NumField())
for i := range v.NumField() {
    m.globalKeyBinds[i] = v.Field(i).Interface().(key.Binding)
}
```

**Fix**: Replace with explicit slice construction:
```go
m.globalKeyBinds = []key.Binding{
    m.keymap.NextPanel,
    m.keymap.PrevPanel,
    m.keymap.Quit,
    m.keymap.ToggleGQLType,
    m.keymap.ReverseToggleGQLType,
    m.keymap.ToggleOverlay,
    m.keymap.ToggleFavorite,
}
```

**Benefits**:
- Clearer and more maintainable
- No reflection overhead
- Compile-time safety

**Estimated Impact**: More readable, same line count

---

### 4. Split Large model.go File
**Issue**: `tui/xplr/model.go` is 614 lines and handles multiple responsibilities:
- State management
- Update logic (Update method)
- Panel management
- Favorites logic
- Rendering (View, renderGQLTypeNavbar, renderBreadcrumbs)

**Fix**: Extract rendering functions to separate file `tui/xplr/view.go`:
- Move `View()` method
- Move `renderGQLTypeNavbar()` method
- Move `renderBreadcrumbs()` method

**Benefits**:
- Better separation of concerns
- Easier to navigate and maintain
- Follows single responsibility principle

**Estimated Impact**: Split 614-line file into ~450 + ~150 lines

---

### 5. Reduce Getter/Setter Boilerplate (DISCUSS)
**Issue**: Many simple getter/setter methods in `tui/xplr/model.go:150-178`:
- `SetSchemaID()` / `GetSchemaID()`
- `SetFavorites()` / `GetFavorites()`
- `SetHasLibraryData()` / `HasLibraryData()`
- `Width()` / `Height()`

**Note**: Project coding guidelines say "Use accessor methods instead of exposing struct fields"

**Options**:
1. Keep as-is (follows current guidelines)
2. Make select fields public where they don't need validation
3. Keep getters, remove setters and set fields internally only

**Question**: What's your preference?

---

### 6. Address Existing TODOs/FIXMEs

#### High Priority
- `tui/xplr/components/panels.go:42` - FIXME: Clear old panels if current item can't be opened
  ```go
  // FIXME: This should also clear old panels if current item can't be opened
  ```

#### Medium Priority
- `tui/adapters/items.go:248` - Better error handling for field result types
- `tui/adapters/items.go:267` - Better error handling for argument types
  ```go
  // TODO: Currently, this treats any error as a built-in type, but instead we should
  // check for _known_ built in types and handle errors intelligently.
  ```

#### Lower Priority
- `tui/libselect/model.go:117` - Show errors in UI when loading schema fails
- `tui/libselect/model.go:123` - Show errors in UI when starting explorer fails
- `gql/parse_test.go:233` - DirectiveDefinition wrapper needs to expose Arguments and Locations

---

### 7. Simplify Favorites Wrapper Pattern (CONSIDER)
**Issue**: Every item gets wrapped/unwrapped when handling favorites

**Current flow**:
1. Items wrapped with `wrapItemsWithFavorites()` to add star indicator
2. When favorites change, `refreshPanelsWithFavorites()` must:
   - Unwrap all items to get originals
   - Re-wrap items with updated favorites
3. Wrapper delegates all methods to wrapped item

**Files involved**:
- `tui/xplr/favoritable_item.go` - 80 lines of wrapper code
- `tui/xplr/model.go:416-496` - Wrapping/unwrapping logic

**Alternative approach**: Store favorite state alongside items rather than wrapping
- Could use a `map[string]bool` for favorites lookup during rendering
- Title rendering checks favorites map instead of using wrapper

**Trade-off**: Less abstraction overhead but couples rendering to favorites state

**Question**: Is this worth exploring?

---

## Lower Priority Optimizations

### 8. Simplify Panel Result Type Handling
**Issue**: Special "virtual item at top" logic in `tui/xplr/components/panels.go` adds complexity:
- `resultType` field for virtual item
- `focusOnResultType` boolean flag
- Special navigation logic in `Update()` (lines 91-117)

**Current behavior**: Displays result type above argument list with separate focus handling

**Consider**: Whether this pattern could be simplified (e.g., result type as first real list item)

---

### 9. Message Routing Logic
**Issue**: `shouldFocusedPanelReceiveMessage()` in `tui/xplr/model.go:312-329` filters messages based on global keybinds

**Current approach**: Check if key matches any global binding to prevent panel from receiving it

**Consider**: Whether different architecture could make this cleaner (e.g., explicit message routing)

---

## Metrics

**Total Go files**: 50
**Key file sizes**:
- `tui/xplr/model.go`: 614 lines
- `tui/xplr/components/panels.go`: 275 lines
- `tui/xplr/navigation/manager.go`: 125 lines

**Test coverage**: 19 test files

---

## Recommendation

Start with **items 1-4** (high impact, low risk):
1. Remove duplicate GQLType definitions
2. Remove unused code
3. Simplify keymap construction
4. Split model.go file

Then assess whether to continue with remaining items based on impact and team preferences.
