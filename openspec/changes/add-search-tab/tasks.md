# Implementation Tasks

## 1. Add SearchType constant to navigation
- Add `SearchType GQLType = "Search"` to `tui/xplr/navigation/type_selector.go`
- Update `newTypeSelector()` types slice to include `SearchType` after `DirectiveType`
- Verify type cycling works with new type

## 2. Create search input component
- Create `tui/xplr/components/searchinput.go` with `SearchInput` model
- Use `bubbles/textinput` for the input field implementation
- Implement `Init()`, `Update()`, `View()` methods following Bubble Tea patterns
- Add placeholder text "Type to search schema..."
- Support Enter key to submit query, Escape to clear

## 3. Add search state to xplr Model
- Add `searchInput` field to `xplr.Model` struct
- Add `searchFocused` bool to track whether input has focus
- Add `searchResults []search.Result` to cache current results
- Initialize searchInput in `NewEmpty()`

## 4. Implement search query execution
- Create method `executeSearch(query string)` in `xplr.Model`
- Integrate with `search.NewSearcher()` and `Search()` from existing package
- Handle index creation if missing (background task with loading state)
- Convert `search.Result` to `components.ListItem` for panel display
- Store results in `searchResults` field

## 5. Update loadMainPanel for Search tab
- Add `case navigation.SearchType:` in `loadMainPanel()` switch statement
- Set panel title to "Search Results"
- Load items from `searchResults` (empty on initial load)
- Show empty state message if no search executed yet

## 6. Integrate input field in main View
- Modify `xplr.View()` to render search input when `SearchType` is active
- Position input field above help bar, below panels
- Adjust panel height calculation to reserve space for input
- Hide input field when other tabs are active

## 7. Implement focus management
- Add keyboard handler for "/" to return focus to search input
- Track focus state with `searchFocused` field
- Route keyboard input to searchInput when focused, to panels otherwise
- Submit search on Enter and transfer focus to results panel
- Clear input on Escape while maintaining focus

## 8. Add result selection navigation
- Create `SearchResultItem` type implementing `ListItem` interface
- Store result path (e.g., "User.email") in item data
- Implement `OpenPanel()` to navigate to type/field using `ApplySelection()`
- Parse path into `SelectionTarget` struct

## 9. Handle search loading states
- Display "Indexing schema..." message while index builds
- Show "Searching..." during query execution
- Handle errors (missing index, failed search) with appropriate messages

## 10. Add tests
- Test `SearchType` addition to type cycle in `type_selector_test.go`
- Test search input component behavior
- Test search execution and result display
- Test focus management transitions
- Run `just test` to verify all tests pass

## Dependencies
- Task 1 must complete before Tasks 3-5
- Task 2 can run in parallel with Task 1
- Tasks 4-5 depend on Task 3
- Task 6 depends on Tasks 2 and 3
- Task 7 depends on Tasks 2, 3, and 6
- Task 8 depends on Tasks 4 and 5
- Task 9 depends on Task 4
- Task 10 is final validation
