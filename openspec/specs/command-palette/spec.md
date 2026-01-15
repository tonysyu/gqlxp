# command-palette Specification

## Purpose
TBD - created by archiving change add-command-palette. Update Purpose after archive.
## Requirements
### Requirement: Command Palette Trigger
The TUI SHALL provide a command palette that can be triggered via keyboard shortcut from any context.

#### Scenario: Open palette with Ctrl+P
- **WHEN** the user presses Ctrl+P in any TUI context (main xplr view, overlay, search)
- **THEN** the command palette SHALL open and display a list of available commands
- **AND** the palette SHALL intercept all keyboard input until closed

#### Scenario: Palette appears as overlay
- **WHEN** the command palette is opened
- **THEN** it SHALL appear as a centered overlay on top of the current view
- **AND** SHALL use similar styling to the details overlay

#### Scenario: Palette closed by default
- **WHEN** the TUI starts or returns from executing a command
- **THEN** the command palette SHALL be hidden
- **AND** SHALL not intercept any messages

### Requirement: Command List Display
The command palette SHALL display all available keyboard shortcuts as a searchable list.

#### Scenario: Display all keymap commands
- **WHEN** the command palette is opened
- **THEN** it SHALL display commands from Main, Panel, Overlay, and Global keymaps
- **AND** SHALL use `charmbracelet/bubbles/list` component for rendering

#### Scenario: Command title format
- **WHEN** displaying a command in the list
- **THEN** the title SHALL be formatted as "[Context]: [description]"
- **AND** Context SHALL be one of: "Main", "Panel", "Overlay", "Global"
- **AND** description SHALL be the help text from the key binding (e.g., "next panel", "prev tab")

#### Scenario: Command description shows keys
- **WHEN** displaying a command in the list
- **THEN** the description SHALL show the key combination (e.g., "]/tab", "⇧+L", "⌃+c")
- **AND** SHALL use the same format as the help display

#### Scenario: List supports filtering
- **WHEN** the user types in the command palette
- **THEN** the list SHALL filter commands based on title and key text
- **AND** SHALL use the built-in bubbles/list filtering

### Requirement: Command Availability State
The command palette SHALL indicate which commands are available in the current context.

#### Scenario: Inactive commands shown with dimmed style
- **WHEN** a command is not applicable in the current context
- **THEN** it SHALL still appear in the list
- **AND** SHALL be rendered with dimmed/faint styling
- **AND** SHALL remain selectable but visually distinct

#### Scenario: Search commands inactive outside Search tab
- **WHEN** the command palette is opened and current tab is not Search
- **THEN** search-related commands (SearchFocus, SearchSubmit, SearchClear) SHALL be marked inactive
- **AND** SHALL be rendered with dimmed styling

#### Scenario: Overlay commands inactive when no overlay
- **WHEN** the command palette is opened and overlay is not active
- **THEN** overlay-specific commands (Overlay.Close) SHALL be marked inactive
- **AND** SHALL be rendered with dimmed styling

#### Scenario: Panel commands inactive in overlay
- **WHEN** the command palette is opened while overlay is active
- **THEN** panel-specific commands (NextTab, PrevTab) SHALL be marked inactive
- **AND** SHALL be rendered with dimmed styling

#### Scenario: All commands active in main context
- **WHEN** the command palette is opened in main xplr view with no overlay
- **THEN** all non-search commands SHALL be marked active
- **AND** SHALL be rendered with normal styling

### Requirement: Command Execution
The command palette SHALL execute the selected command when the user presses Enter.

#### Scenario: Execute command on Enter
- **WHEN** the user selects a command and presses Enter
- **THEN** the command palette SHALL close
- **AND** the selected command's keyboard shortcut SHALL be executed
- **AND** the action SHALL be the same as if the user pressed the actual key combination

#### Scenario: Disabled commands can be selected
- **WHEN** the user selects an inactive/disabled command and presses Enter
- **THEN** the command palette SHALL close
- **AND** no action SHALL be executed (command is no-op in current context)

#### Scenario: Execute sends KeyMsg
- **WHEN** executing a command
- **THEN** the palette SHALL send a `tea.KeyMsg` with the keys from the binding
- **AND** the message SHALL be processed by the main xplr model's Update method

### Requirement: Palette Dismissal
The command palette SHALL close without executing a command when certain keys are pressed.

#### Scenario: Close with Space
- **WHEN** the user presses Space in the command palette
- **THEN** the palette SHALL close
- **AND** no command SHALL be executed

#### Scenario: Close with q
- **WHEN** the user presses "q" in the command palette
- **THEN** the palette SHALL close
- **AND** no command SHALL be executed

#### Scenario: Close with Escape
- **WHEN** the user presses Escape in the command palette
- **THEN** the palette SHALL close
- **AND** no command SHALL be executed

#### Scenario: Palette intercepts close keys
- **WHEN** the palette is open and close keys are pressed
- **THEN** they SHALL close the palette instead of being passed to underlying views
- **AND** the palette SHALL return intercepted=true for these messages

### Requirement: Message Interception
The command palette SHALL intercept all messages when active, similar to the overlay pattern.

#### Scenario: Intercept all input when active
- **WHEN** the command palette is active
- **THEN** all keyboard messages SHALL be intercepted
- **AND** SHALL not be passed to the main xplr model
- **AND** the Update method SHALL return intercepted=true

#### Scenario: Pass through when inactive
- **WHEN** the command palette is not active
- **THEN** all messages SHALL pass through to the main model
- **AND** the Update method SHALL return intercepted=false

#### Scenario: Handle window size updates
- **WHEN** a WindowSizeMsg is received while palette is active
- **THEN** the palette SHALL update its dimensions
- **AND** SHALL re-center the overlay
- **AND** SHALL mark the message as intercepted

