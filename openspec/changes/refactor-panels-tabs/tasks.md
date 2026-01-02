# Implementation Tasks

## 1. Core Panel Tab Infrastructure
- [ ] 1.1 Add Tab struct with Label and Content fields
- [ ] 1.2 Add tab state fields to Panel struct (tabs []Tab, activeTab int)
- [ ] 1.3 Remove `resultType`, `focusOnResultType`, and `lastSelectedIndex` fields
- [ ] 1.4 Implement tab rendering with active/inactive styling (based on bubbletea tabs example)
- [ ] 1.5 Add `SetTabs([]Tab)` method to Panel to configure tabs
- [ ] 1.6 Update Panel.View() to render tabs and active tab content
- [ ] 1.7 Write unit tests for tab rendering and state management

## 2. Tab Navigation
- [ ] 2.1 Add Shift-H/Shift-L keybinding handling in Panel.Update()
- [ ] 2.2 Implement tab switching logic (prev/next with boundary checks)
- [ ] 2.3 Update content area when active tab changes
- [ ] 2.4 Reset list selection to first item when switching tabs
- [ ] 2.5 Write tests for tab navigation and wrapping behavior

## 3. Adapter Updates
- [ ] 3.1 Update `fieldItem.OpenPanel()` to create Tab objects for "Result Type" and "Input Arguments"
- [ ] 3.2 Update `argumentItem.OpenPanel()` to create Tab object for "Result Type"
- [ ] 3.3 Update `typeDefItem.OpenPanel()` to use `SetTabs()` or keep single content display
- [ ] 3.4 Write tests for adapter tab configuration

## 4. Global Keybinding Integration
- [ ] 4.1 Add Shift-H and Shift-L keybindings to `tui/xplr/model.go`
- [ ] 4.2 Route tab navigation keys to focused panel
- [ ] 4.3 Update help text to document Shift-H/Shift-L navigation
- [ ] 4.4 Write integration tests for keybindings

## 5. Test Updates and Validation
- [ ] 5.1 Update `panels_test.go` to remove result type special-case tests
- [ ] 5.2 Add tests for tab-based navigation scenarios
- [ ] 5.3 Update acceptance tests if they rely on old result type behavior
- [ ] 5.4 Run `just test` to verify all tests pass
- [ ] 5.5 Manual testing of tab navigation in TUI

## 6. Cleanup
- [ ] 6.1 Remove deprecated `SetObjectType()` method
- [ ] 6.2 Remove result type section rendering code from Panel.View()
- [ ] 6.3 Update documentation if needed (per CLAUDE.md, put in docs/* not CLAUDE.md)
- [ ] 6.4 Final verification with `just verify`
