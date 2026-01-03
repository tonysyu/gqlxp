# Tasks for add-interfaces-subtab

## Implementation Tasks

1. **Add helper to adapt interface names to list items**
   - Add `AdaptInterfaces(interfaceNames []string, resolver TypeResolver) []ListItem` function in `tui/adapters/items.go`
   - Convert interface name strings to typeDefItems using newNamedItem
   - Return list items sorted alphabetically by interface name
   - Add unit tests for AdaptInterfaces

2. **Modify Object type panel to include Interfaces tab**
   - Update `typeDefItem.OpenPanel()` in `tui/adapters/items.go` for Object case
   - Add "Interfaces" tab when `typeDef.Interfaces()` returns non-empty slice
   - Use AdaptInterfaces to create interface list items
   - Tab order: "Fields", "Interfaces"
   - Add unit test verifying Object with interfaces creates both tabs

3. **Add integration test for navigation flow**
   - Create test schema with Object implementing Interface
   - Verify Object panel has Interfaces tab with correct content
   - Verify selecting an interface navigates to the Interface panel
   - Verify Interface panel displays its fields

4. **Run validation and tests**
   - Run `just test` to ensure all tests pass
   - Run `openspec validate add-interfaces-subtab --strict` to verify spec compliance
   - Manually test TUI with real schema containing interface implementations

## Testing Strategy

- Unit tests for new adapter function (AdaptInterfaces)
- Unit tests for modified OpenPanel method (Object case)
- Integration test for navigation from Object to Interface
- Manual testing with GitHub schema (User implements Node)

## Dependencies

- Tasks 1-2 can be implemented in parallel
- Task 3 requires Tasks 1-2 to complete
- Task 4 validates entire implementation
