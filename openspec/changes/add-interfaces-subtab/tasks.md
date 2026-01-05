# Tasks for add-interfaces-subtab

## Implementation Tasks

1. **Add helper to adapt interface names to list items**
   - [x] Add `AdaptInterfaces(interfaceNames []string, resolver TypeResolver) []ListItem` function in `tui/adapters/items.go`
   - [x] Convert interface name strings to typeDefItems using newNamedItem
   - [x] Return list items sorted alphabetically by interface name
   - [x] Add unit tests for AdaptInterfaces

2. **Modify Object type panel to include Interfaces tab**
   - [x] Update `typeDefItem.OpenPanel()` in `tui/adapters/items.go` for Object case
   - [x] Add "Interfaces" tab when `typeDef.Interfaces()` returns non-empty slice
   - [x] Use AdaptInterfaces to create interface list items
   - [x] Tab order: "Fields", "Interfaces"
   - [x] Add unit test verifying Object with interfaces creates both tabs

3. **Add integration test for navigation flow**
   - [x] Create test schema with Object implementing Interface
   - [x] Verify Object panel has Interfaces tab with correct content
   - [x] Verify selecting an interface navigates to the Interface panel
   - [x] Verify Interface panel displays its fields

4. **Run validation and tests**
   - [x] Run `just test` to ensure all tests pass
   - [x] Run `openspec validate add-interfaces-subtab --strict` to verify spec compliance
   - [x] Manually test TUI with real schema containing interface implementations

## Testing Strategy

- Unit tests for new adapter function (AdaptInterfaces)
- Unit tests for modified OpenPanel method (Object case)
- Integration test for navigation from Object to Interface
- Manual testing with GitHub schema (User implements Node)

## Dependencies

- Tasks 1-2 can be implemented in parallel
- Task 3 requires Tasks 1-2 to complete
- Task 4 validates entire implementation
