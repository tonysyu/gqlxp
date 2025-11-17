# Tasks: Fix Favorites Behavior

## Implementation Tasks

- [x] **Modify favorite toggle logic to use context-aware name selection**
   - Update `toggleFavoriteForSelectedItem()` in `tui/xplr/model.go` to check if current panel is a top-level GQL type panel
   - For top-level panels, use `RefName()` instead of `TypeName()` when storing favorites
   - For non-top-level panels, continue using `TypeName()`
   - **Validation**: Manual testing - toggle favorites in Query panel stores field names, toggle in type detail panel stores type names

- [x] **Preserve selected item after favorite toggle**
   - Capture selected item's identifier (RefName or TypeName) before toggling favorite
   - After receiving `FavoriteToggledMsg` and refreshing panel, restore selection to matching item
   - **Validation**: Manual testing - selected item remains highlighted after pressing 'f' key

- [x] **Add favorite indicator to panel titles**
   - Modify panel creation/update logic to check if panel's type/field name is in favorites list
   - When favorited, prefix panel title with "★ "
   - **Validation**: Manual testing - favorite a Query field, navigate to its detail panel, verify "★" appears in panel title

- [x] **Run full test suite**
   - Execute `just test` to ensure no regressions
   - **Validation**: All tests pass

- [x] **Manual integration testing**
   - Test favoriting Query fields, Mutation fields, and other types
   - Verify selection preservation after toggling
   - Verify panel title indicators appear correctly
   - Test unfavoriting to ensure indicators are removed
   - **Validation**: All scenarios work as expected

## Notes
- Top-level panels are those displaying GQL type categories: Query, Mutation, Object, Input, Enum, Scalar, Interface, Union, Directive
- The implementation should check if `nav.TypeSelector().Current()` matches one of these types to determine context
