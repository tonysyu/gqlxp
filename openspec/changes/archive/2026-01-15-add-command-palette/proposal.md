# Proposal: Add Command Palette

## Overview
Add a Command Palette to the TUI explorer using `charmbracelet/bubbles/list` that provides a searchable list of all available keyboard shortcuts. Users can trigger actions by selecting commands from the palette instead of memorizing keybindings.

## Why
Currently, users must:
- Remember keyboard shortcuts across different contexts (Main, Panel, Overlay)
- Use help view to discover available commands
- Execute commands only via keybindings

A Command Palette provides:
- Discoverability of all available commands in one place
- Keyboard-driven command execution without memorizing shortcuts
- Context-aware display showing which commands are currently available
- Common UX pattern familiar from editors (VS Code, Sublime Text)

## What Changes
- Add new `tui/cmdpalette` package with Command Palette model using `bubbles/list`
- Create `commandItem` type representing a keymap command with context, description, keys, and enabled state
- Add Ctrl+P keybinding to trigger palette from any context
- Integrate palette into `xplr.Model` with message interception similar to overlay pattern
- Collect all commands from MainKeymaps, PanelKeymaps, OverlayKeymaps, and GlobalKeymaps
- Display inactive commands with dimmed styling based on current context
- Execute commands by sending KeyMsg when user selects and presses Enter
- Close palette with Space/q/Esc without executing

## Impact
- Affected specs: `tui-architecture` (new TUI component), `interface-navigation` (new keymap entry)
- Affected code:
  - `tui/cmdpalette/model.go` - New command palette component
  - `tui/xplr/model.go` - Integration and message handling
  - `tui/config/keys.go` - Add palette trigger keybinding
  - `tui/xplr/view.go` - Render palette overlay

## Scope
- Add Command Palette model using `bubbles/list` component
- Trigger palette with Ctrl+P keybinding (global, works in all contexts)
- Display all keymaps from Main, Panel, Overlay, and Global contexts
- Show inactive commands with dimmed styling when not applicable to current context
- Execute focused command on Enter and close palette
- Close palette on Space/q/Esc without executing
- Format command titles as "[Context]: [description]" (e.g., "Main: next panel", "Panel: next tab")
- Format command descriptions as the key combination (e.g., "]/tab", "⇧+L")

## Out of Scope
- Command history or recently used commands
- Custom command aliases or user-defined commands
- Fuzzy search beyond built-in list filtering
- Command categories or grouping beyond context prefix

## Related Specs
- `tui-architecture` - TUI structure and model patterns
- `interface-navigation` - Keymap definitions and help display

## Risks and Mitigation
- **Risk**: Ctrl+P may conflict with terminal or shell shortcuts
  - **Mitigation**: Document the keybinding; users can remap if needed
- **Risk**: Command palette may intercept keys needed for other views
  - **Mitigation**: Command palette intercepts all messages when active, similar to overlay pattern
- **Risk**: Determining command availability context may be complex
  - **Mitigation**: Start with simple rules (Main context = xplr active, Panel = panel focused, Overlay = overlay active)
