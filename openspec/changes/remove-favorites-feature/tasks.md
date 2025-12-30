# Tasks: Remove Favorites Feature

## Implementation Tasks

1. **Remove TUI favorites keybinding and toggle logic**
   - Remove `ToggleFavorite` field from `keymap` struct in `tui/xplr/model.go`
   - Remove `ToggleFavorite` key binding initialization
   - Remove `toggleFavoriteForSelectedItem()` method
   - Remove `toggleFavorite()` method
   - Remove case statement handling `ToggleFavorite` key in `Update()`
   - Remove `ToggleFavorite` from help view in `tui/xplr/view.go`
   - **Verify**: Code compiles without errors

2. **Remove favoritable item wrapper**
   - Delete `tui/xplr/favoritable_item.go` entirely
   - Remove `wrapItemsWithFavorites` calls in `tui/xplr/model.go`
   - **Verify**: Code compiles, TUI displays items without star indicators

3. **Remove Favorites field from TUI models**
   - Remove `Favorites []string` from `InitMsg` in `tui/xplr/model.go`
   - Remove `Favorites []string` from `Options` in `tui/xplr/model.go`
   - Remove `Favorites []string` from `Model` in `tui/xplr/model.go`
   - Remove favorites initialization in `New()` function
   - Remove `Favorites` field from `TUIOptions` in `tui/model.go`
   - Remove favorites assignment from library metadata in `tui/app.go`
   - **Verify**: TUI starts without errors

4. **Remove library favorites methods**
   - Remove `AddFavorite()` method from `library/library.go`
   - Remove `RemoveFavorite()` method from `library/library.go`
   - Remove `AddFavorite` and `RemoveFavorite` from `Library` interface
   - **Verify**: Library compiles without errors

5. **Remove Favorites field from library metadata**
   - Remove `Favorites []string` field from `SchemaMetadata` in `library/types.go`
   - Remove `Favorites: []string{}` initialization in `Add()` method
   - Remove `Favorites: []string{}` initialization in `Update()` method
   - **Verify**: Library continues to load existing metadata.json files

6. **Run full test suite**
   - Run `just test` to verify all tests pass
   - Fix any test failures related to favorites removal
   - Update `tui/model_test.go` to remove `Favorites` field from test data
   - **Verify**: All tests pass

7. **Update documentation (if needed)**
   - Check if `docs/schema-library.md` mentions favorites
   - Remove any favorites documentation if present
   - **Verify**: Documentation is consistent with code

## Validation

- [x] Code compiles without errors
- [x] All tests pass (`just test`)
- [x] TUI starts and functions normally
- [x] Library loads existing schemas without errors
- [x] No references to favorites remain in codebase (verify with grep)
