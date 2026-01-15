# Implementation Tasks

## 1. Create command palette model structure
- [ ] Create `tui/cmdpalette/model.go` with `Model` struct
- [ ] Add `list.Model` field for the bubbles list component
- [ ] Add `active` bool field to track visibility state
- [ ] Add `width`, `height` fields for sizing
- [ ] Add `keymap` field with palette-specific keybindings (Close)
- [ ] Implement `New()` constructor initializing the list component

## 2. Define command palette item types
- [ ] Create `commandItem` struct implementing `list.Item` interface
- [ ] Add fields: `title` (string), `keys` (string), `binding` (key.Binding), `enabled` (bool)
- [ ] Implement `Title()` to return formatted "[Context]: [description]"
- [ ] Implement `Description()` to return key combination string
- [ ] Implement `FilterValue()` to return searchable text (title + keys)

## 3. Build command items from keymaps
- [ ] Create `buildCommandItems()` function taking all keymap structs as parameters
- [ ] Extract commands from `MainKeymaps`: NextPanel, PrevPanel, NextGQLType, PrevGQLType, ToggleOverlay, SearchFocus, SearchSubmit, SearchClear
- [ ] Extract commands from `PanelKeymaps`: NextTab, PrevTab
- [ ] Extract commands from `OverlayKeymaps`: Close
- [ ] Extract commands from `GlobalKeymaps`: Quit
- [ ] Format titles with context prefix (e.g., "Main: next panel")
- [ ] Return slice of `commandItem`

## 4. Implement command availability logic
- [ ] Create `updateCommandAvailability()` method on `Model`
- [ ] Accept current state parameters: `overlayActive`, `searchFocused`, `currentGQLType`
- [ ] Set `enabled=false` for SearchFocus/SearchSubmit/SearchClear when not on Search tab
- [ ] Set `enabled=false` for Overlay.Close when overlay not active
- [ ] Set `enabled=false` for Panel keymaps when overlay is active
- [ ] Update list items with new enabled states

## 5. Implement Model methods
- [ ] Implement `Init() tea.Cmd` returning nil
- [ ] Implement `Update(msg tea.Msg) (Model, tea.Cmd, bool)` with intercepted flag
- [ ] Handle `tea.KeyMsg` for Close keybinding (Space, q, Esc) to hide palette
- [ ] Handle `tea.KeyMsg` for Select keybinding (Enter) to execute command and hide palette
- [ ] Handle `tea.WindowSizeMsg` to update dimensions
- [ ] Forward other messages to list.Update() when palette is active
- [ ] Return `intercepted=true` when palette is active, `false` otherwise

## 6. Implement command execution
- [ ] Create `executeCommand(item commandItem) tea.Cmd` method
- [ ] Return tea.Cmd that sends the original KeyMsg for the command's binding
- [ ] Use `key.Binding.Keys()` to get the key string and construct KeyMsg
- [ ] Handle case where command is disabled (no-op)

## 7. Implement View rendering
- [ ] Implement `View() string` method
- [ ] Return empty string when palette not active
- [ ] Apply custom delegate with dimmed style for disabled items
- [ ] Render list with title "Command Palette"
- [ ] Center palette overlay on screen (similar to overlay.Model)
- [ ] Use `Styles.Overlay` for consistent styling

## 8. Integrate palette into xplr Model
- [ ] Add `cmdpalette` field to `xplr.Model` struct
- [ ] Initialize in `NewEmpty()` by calling cmdpalette.New()
- [ ] Add Ctrl+P keybinding to `MainKeymaps` (or create new field for global palette key)
- [ ] Handle Ctrl+P in `xplr.Update()` to show palette
- [ ] Call `palette.updateCommandAvailability()` before showing
- [ ] Update palette in `xplr.Update()` when active, capture intercepted flag
- [ ] Skip other Update logic when palette intercepts messages

## 9. Update xplr View to render palette
- [ ] Modify `xplr.View()` to render palette overlay when active
- [ ] Layer palette view on top of main content (similar to overlay)
- [ ] Ensure palette appears above all other UI elements

## 10. Handle command palette keybindings
- [ ] Add palette toggle key to help display when appropriate
- [ ] Ensure Ctrl+P works from all xplr contexts (main, search, overlay)
- [ ] Test that Space/q close palette without executing command
- [ ] Test that Enter executes focused command and closes palette
- [ ] Test that Esc closes palette without executing

## 11. Style disabled commands
- [ ] Create custom list delegate with `Render()` implementation
- [ ] Apply dimmed/faint styling to disabled items
- [ ] Use normal styling for enabled items
- [ ] Ensure disabled items are still selectable but show as inactive

## 12. Add tests
- [ ] Test `commandItem` implements list.Item interface correctly
- [ ] Test `buildCommandItems()` returns all expected commands
- [ ] Test `updateCommandAvailability()` correctly enables/disables based on state
- [ ] Test command execution sends correct KeyMsg
- [ ] Test palette intercepts messages when active
- [ ] Run `just test` to verify all tests pass

## Dependencies
- Tasks 1-2 can run in parallel
- Task 3 depends on Task 2
- Task 4 depends on Tasks 1 and 3
- Task 5 depends on Tasks 1, 3, 4
- Task 6 depends on Tasks 2, 3
- Task 7 depends on Tasks 1, 5, 11
- Task 8 depends on Tasks 1, 4, 5
- Task 9 depends on Tasks 7, 8
- Task 10 depends on Tasks 8, 9
- Task 11 can be done in parallel with Tasks 5-7
- Task 12 is final validation
