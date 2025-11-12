# Task 05: Command Pattern for Actions

**Priority:** Low
**Status:** Not Started
**Estimated Effort:** Medium
**Dependencies:** Task 02 (Navigation Manager) - recommended

## Problem Statement

Currently, actions like navigation, panel opening, and type switching are tightly coupled to `mainModel.Update()`. This makes it difficult to:

- Test actions in isolation
- Implement undo/redo functionality
- Log or replay user actions
- Create keyboard macros or automation
- Separate UI event handling from business logic

### Affected Files
- `tui/model.go` - Update method with mixed action handling
- Future automation/scripting features

## Proposed Solution

Introduce a Command pattern to encapsulate actions as objects:

```go
// tui/commands/commands.go
package commands

import (
    "github.com/tonysyu/gqlxp/tui/components"
    "github.com/tonysyu/gqlxp/tui/navigation"
)

// Command represents an action that can be executed
type Command interface {
    // Execute performs the command and returns error if it fails
    Execute(ctx *Context) error

    // Undo reverses the command (optional, return error if not supported)
    Undo(ctx *Context) error

    // Description returns human-readable description
    Description() string

    // CanUndo returns whether this command supports undo
    CanUndo() bool
}

// Context provides access to application state for commands
type Context struct {
    Navigation  *navigation.NavigationManager
    Schema      SchemaProvider
    PanelSizer  PanelSizer
    // Add other state as needed
}

type SchemaProvider interface {
    LoadItemsForType(typeID string) []components.ListItem
}

type PanelSizer interface {
    SizePanels()
}

// NavigateForwardCommand moves forward in panel stack
type NavigateForwardCommand struct {
    executed bool
}

func NewNavigateForwardCommand() *NavigateForwardCommand {
    return &NavigateForwardCommand{}
}

func (c *NavigateForwardCommand) Execute(ctx *Context) error {
    // Add current item to breadcrumbs before moving
    c.executed = ctx.Navigation.NavigateForward(func() {
        // Add breadcrumb callback
    })
    if c.executed {
        ctx.PanelSizer.SizePanels()
    }
    return nil
}

func (c *NavigateForwardCommand) Undo(ctx *Context) error {
    if !c.executed {
        return nil
    }
    return NewNavigateBackwardCommand().Execute(ctx)
}

func (c *NavigateForwardCommand) Description() string {
    return "Navigate forward"
}

func (c *NavigateForwardCommand) CanUndo() bool {
    return true
}

// NavigateBackwardCommand moves backward in panel stack
type NavigateBackwardCommand struct {
    executed bool
}

func NewNavigateBackwardCommand() *NavigateBackwardCommand {
    return &NavigateBackwardCommand{}
}

func (c *NavigateBackwardCommand) Execute(ctx *Context) error {
    c.executed = ctx.Navigation.NavigateBackward()
    return nil
}

func (c *NavigateBackwardCommand) Undo(ctx *Context) error {
    if !c.executed {
        return nil
    }
    return NewNavigateForwardCommand().Execute(ctx)
}

func (c *NavigateBackwardCommand) Description() string {
    return "Navigate backward"
}

func (c *NavigateBackwardCommand) CanUndo() bool {
    return true
}

// OpenPanelCommand opens a new panel
type OpenPanelCommand struct {
    panel components.Panel
}

func NewOpenPanelCommand(panel components.Panel) *OpenPanelCommand {
    return &OpenPanelCommand{panel: panel}
}

func (c *OpenPanelCommand) Execute(ctx *Context) error {
    ctx.Navigation.OpenPanel(c.panel)
    ctx.PanelSizer.SizePanels()
    return nil
}

func (c *OpenPanelCommand) Undo(ctx *Context) error {
    // Would need to track previous state
    return NewNavigateBackwardCommand().Execute(ctx)
}

func (c *OpenPanelCommand) Description() string {
    return "Open panel"
}

func (c *OpenPanelCommand) CanUndo() bool {
    return true
}

// SwitchTypeCommand changes the selected GQL type
type SwitchTypeCommand struct {
    newTypeID   string
    prevTypeID  string
}

func NewSwitchTypeCommand(typeID string) *SwitchTypeCommand {
    return &SwitchTypeCommand{newTypeID: typeID}
}

func (c *SwitchTypeCommand) Execute(ctx *Context) error {
    c.prevTypeID = string(ctx.Navigation.CurrentType())
    ctx.Navigation.SwitchType(navigation.GQLType(c.newTypeID))

    // Load items for new type
    items := ctx.Schema.LoadItemsForType(c.newTypeID)
    // Create and set panel...

    return nil
}

func (c *SwitchTypeCommand) Undo(ctx *Context) error {
    return NewSwitchTypeCommand(c.prevTypeID).Execute(ctx)
}

func (c *SwitchTypeCommand) Description() string {
    return "Switch type to " + c.newTypeID
}

func (c *SwitchTypeCommand) CanUndo() bool {
    return true
}

// CycleTypeForwardCommand cycles to next type
type CycleTypeForwardCommand struct{}

func NewCycleTypeForwardCommand() *CycleTypeForwardCommand {
    return &CycleTypeForwardCommand{}
}

func (c *CycleTypeForwardCommand) Execute(ctx *Context) error {
    newType := ctx.Navigation.CycleTypeForward()
    return NewSwitchTypeCommand(string(newType)).Execute(ctx)
}

func (c *CycleTypeForwardCommand) Undo(ctx *Context) error {
    return NewCycleTypeBackwardCommand().Execute(ctx)
}

func (c *CycleTypeForwardCommand) Description() string {
    return "Cycle to next type"
}

func (c *CycleTypeForwardCommand) CanUndo() bool {
    return true
}

// CycleTypeBackwardCommand cycles to previous type
type CycleTypeBackwardCommand struct{}

func NewCycleTypeBackwardCommand() *CycleTypeBackwardCommand {
    return &CycleTypeBackwardCommand{}
}

func (c *CycleTypeBackwardCommand) Execute(ctx *Context) error {
    newType := ctx.Navigation.CycleTypeBackward()
    return NewSwitchTypeCommand(string(newType)).Execute(ctx)
}

func (c *CycleTypeBackwardCommand) Undo(ctx *Context) error {
    return NewCycleTypeForwardCommand().Execute(ctx)
}

func (c *CycleTypeBackwardCommand) Description() string {
    return "Cycle to previous type"
}

func (c *CycleTypeBackwardCommand) CanUndo() bool {
    return true
}

// ToggleOverlayCommand shows/hides overlay
type ToggleOverlayCommand struct {
    show bool
}

func NewToggleOverlayCommand(show bool) *ToggleOverlayCommand {
    return &ToggleOverlayCommand{show: show}
}

func (c *ToggleOverlayCommand) Execute(ctx *Context) error {
    // Would need overlay in context
    return nil
}

func (c *ToggleOverlayCommand) Undo(ctx *Context) error {
    return NewToggleOverlayCommand(!c.show).Execute(ctx)
}

func (c *ToggleOverlayCommand) Description() string {
    if c.show {
        return "Show overlay"
    }
    return "Hide overlay"
}

func (c *ToggleOverlayCommand) CanUndo() bool {
    return true
}
```

### Command History

Add command history for undo/redo:

```go
// tui/commands/history.go
package commands

// History tracks executed commands for undo/redo
type History struct {
    commands []Command
    position int
    maxSize  int
}

func NewHistory(maxSize int) *History {
    return &History{
        commands: make([]Command, 0, maxSize),
        position: -1,
        maxSize:  maxSize,
    }
}

// Execute runs command and adds to history
func (h *History) Execute(cmd Command, ctx *Context) error {
    if err := cmd.Execute(ctx); err != nil {
        return err
    }

    // Truncate history after current position
    h.commands = h.commands[:h.position+1]

    // Add command
    h.commands = append(h.commands, cmd)
    h.position++

    // Trim if exceeds max size
    if len(h.commands) > h.maxSize {
        h.commands = h.commands[1:]
        h.position--
    }

    return nil
}

// Undo reverses last command
func (h *History) Undo(ctx *Context) error {
    if !h.CanUndo() {
        return nil
    }

    cmd := h.commands[h.position]
    if !cmd.CanUndo() {
        return nil // Skip non-undoable commands
    }

    if err := cmd.Undo(ctx); err != nil {
        return err
    }

    h.position--
    return nil
}

// Redo re-executes undone command
func (h *History) Redo(ctx *Context) error {
    if !h.CanRedo() {
        return nil
    }

    h.position++
    cmd := h.commands[h.position]
    return cmd.Execute(ctx)
}

// CanUndo returns whether undo is possible
func (h *History) CanUndo() bool {
    return h.position >= 0
}

// CanRedo returns whether redo is possible
func (h *History) CanRedo() bool {
    return h.position < len(h.commands)-1
}

// Clear removes all history
func (h *History) Clear() {
    h.commands = h.commands[:0]
    h.position = -1
}

// GetHistory returns list of command descriptions
func (h *History) GetHistory() []string {
    descriptions := make([]string, len(h.commands))
    for i, cmd := range h.commands {
        descriptions[i] = cmd.Description()
    }
    return descriptions
}
```

### Update mainModel

Integrate commands into `mainModel`:

```go
// tui/model.go
type mainModel struct {
    // ... existing fields
    cmdHistory *commands.History
    cmdContext *commands.Context
}

func newModel(schema adapters.SchemaView) mainModel {
    m := mainModel{
        // ... existing initialization
        cmdHistory: commands.NewHistory(100),
    }

    m.cmdContext = &commands.Context{
        Navigation: m.navigation,
        Schema:     &m.schema,
        PanelSizer: &m,
    }

    return m
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // ... overlay handling ...

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, m.keymap.NextPanel):
            m.cmdHistory.Execute(commands.NewNavigateForwardCommand(), m.cmdContext)
        case key.Matches(msg, m.keymap.PrevPanel):
            m.cmdHistory.Execute(commands.NewNavigateBackwardCommand(), m.cmdContext)
        case key.Matches(msg, m.keymap.ToggleGQLType):
            m.cmdHistory.Execute(commands.NewCycleTypeForwardCommand(), m.cmdContext)
        case key.Matches(msg, m.keymap.ReverseToggleGQLType):
            m.cmdHistory.Execute(commands.NewCycleTypeBackwardCommand(), m.cmdContext)
        case key.Matches(msg, m.keymap.Undo):
            m.cmdHistory.Undo(m.cmdContext)
        case key.Matches(msg, m.keymap.Redo):
            m.cmdHistory.Redo(m.cmdContext)
        }
    case components.OpenPanelMsg:
        m.cmdHistory.Execute(commands.NewOpenPanelCommand(msg.Panel), m.cmdContext)
    }

    // ...
}
```

## Benefits

1. **Testability**: Commands testable in isolation
2. **Undo/Redo**: Free undo/redo functionality
3. **Logging**: Easy to log all actions
4. **Macros**: Can record and replay command sequences
5. **Separation**: Clear boundary between UI and logic
6. **Debugging**: Command history helps debugging

## Implementation Steps

1. Create `tui/commands/` package
2. Define `Command` interface and `Context`
3. Implement core commands (navigate, switch type, open panel)
4. Implement `History` for undo/redo
5. Add tests for commands and history
6. Update `mainModel` to use commands
7. Add undo/redo keybindings
8. Run tests: `just test`
9. Update documentation

## Testing Strategy

```go
// tui/commands/commands_test.go
func TestNavigateForwardCommand(t *testing.T) {
    ctx := newTestContext()
    cmd := NewNavigateForwardCommand()

    err := cmd.Execute(ctx)
    assert.NoError(t, err)
    assert.Equal(t, 1, ctx.Navigation.Stack().Position())
}

// tui/commands/history_test.go
func TestHistory_UndoRedo(t *testing.T) {
    ctx := newTestContext()
    history := NewHistory(10)

    // Execute commands
    history.Execute(NewNavigateForwardCommand(), ctx)
    history.Execute(NewNavigateForwardCommand(), ctx)
    assert.Equal(t, 2, ctx.Navigation.Stack().Position())

    // Undo
    history.Undo(ctx)
    assert.Equal(t, 1, ctx.Navigation.Stack().Position())

    // Redo
    history.Redo(ctx)
    assert.Equal(t, 2, ctx.Navigation.Stack().Position())
}
```

## Potential Issues

- **Complexity**: Adds abstraction layer
- **State management**: Commands need access to app state
- **Memory**: History can grow large (limit with maxSize)

## Future Enhancements

1. **Macro recording**: Record command sequences
2. **Scripting**: Execute commands from scripts
3. **Remote control**: Accept commands over network
4. **Command palette**: Show available commands in UI
5. **Keyboard shortcuts**: Bind arbitrary commands to keys
6. **Command composition**: Combine commands into composite commands

## Related Tasks

- **Task 02** (Navigation Manager): Commands operate on navigation manager
- **Task 10** (Event Bus): Commands could publish events
- **Task 09** (Configuration): Load command bindings from config
