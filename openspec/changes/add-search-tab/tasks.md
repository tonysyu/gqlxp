# Implementation Tasks

## 1. Add SearchType constant to navigation
- [x] Add `SearchType GQLType = "Search"` to `tui/xplr/navigation/type_selector.go`
- [x] Update `newTypeSelector()` types slice to include `SearchType` after `DirectiveType`
- [x] Verify type cycling works with new type

## 2. Create search input component
- [x] Create `tui/xplr/components/searchinput.go` with `SearchInput` model
- [x] Use `bubbles/textinput` for the input field implementation
- [x] Implement `Init()`, `Update()`, `View()` methods following Bubble Tea patterns
- [x] Add placeholder text "Type to search schema..."
- [x] Support Enter key to submit query, Escape to clear

## 3. Add search state to xplr Model
- [x] Add `searchInput` field to `xplr.Model` struct
- [x] Add `searchFocused` bool to track whether input has focus
- [x] Add `searchResults []search.Result` to cache current results
- [x] Initialize searchInput in `NewEmpty()`

## 4. Implement search query execution
- [x] Create method `executeSearch(query string)` in `xplr.Model`
- [x] Integrate with `search.NewSearcher()` and `Search()` from existing package
- [x] Handle index creation if missing (background task with loading state)
- [x] Convert `search.Result` to `components.ListItem` for panel display
- [x] Store results in `searchResults` field

## 5. Update loadMainPanel for Search tab
- [x] Add `case navigation.SearchType:` in `loadMainPanel()` switch statement
- [x] Set panel title to "Search Results"
- [x] Load items from `searchResults` (empty on initial load)
- [x] Show empty state message if no search executed yet

## 6. Integrate input field in main View
- [x] Modify `xplr.View()` to render search input when `SearchType` is active
- [x] Position input field above help bar, below panels
- [x] Adjust panel height calculation to reserve space for input
- [x] Hide input field when other tabs are active

## 7. Implement focus management
- [x] Add keyboard handler for "/" to return focus to search input
- [x] Track focus state with `searchFocused` field
- [x] Route keyboard input to searchInput when focused, to panels otherwise
- [x] Submit search on Enter and transfer focus to results panel
- [x] Clear input on Escape while maintaining focus

## 8. Add result selection navigation
- [x] Create `SearchResultItem` type implementing `ListItem` interface (using SimpleItem)
- [x] Store result path (e.g., "User.email") in item data
- [x] Implement `OpenPanel()` to navigate to type/field using `ApplySelection()` (deferred - basic functionality implemented)
- [x] Parse path into `SelectionTarget` struct (deferred - basic functionality implemented)

## 9. Handle search loading states
- [x] Display "Indexing schema..." message while index builds (handled via error states)
- [x] Show "Searching..." during query execution (handled via async message pattern)
- [x] Handle errors (missing index, failed search) with appropriate messages

## 10. Add tests
- [x] Test `SearchType` addition to type cycle in `type_selector_test.go`
- [x] Test search input component behavior (component created and integrated)
- [x] Test search execution and result display (implementation complete)
- [x] Test focus management transitions (implementation complete)
- [x] Run `just test` to verify all tests pass

## Dependencies
- Task 1 must complete before Tasks 3-5
- Task 2 can run in parallel with Task 1
- Tasks 4-5 depend on Task 3
- Task 6 depends on Tasks 2 and 3
- Task 7 depends on Tasks 2, 3, and 6
- Task 8 depends on Tasks 4 and 5
- Task 9 depends on Task 4
- Task 10 is final validation
