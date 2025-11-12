# Task 06: Panel Lifecycle Manager

**Priority:** Low
**Status:** Not Started
**Estimated Effort:** Small-Medium
**Dependencies:** None (complements Task 02)

## Problem Statement

Panel creation, sizing, and focus management are scattered across `mainModel`:

- Panel sizing logic in `sizePanels()` method
- Focus state management in `updatePanelFocusStates()` method
- Panel creation mixed with business logic
- Styling decisions embedded in model

This makes it difficult to:
- Test panel lifecycle independently
- Add panel animations or transitions
- Customize panel behavior
- Reuse panel management in other contexts

### Affected Files
- `tui/model.go` - Panel sizing and focus management
- `tui/components/panels.go` - Panel implementations

## Proposed Solution

Create a dedicated `PanelManager` to handle panel lifecycle:

```go
// tui/panels/manager.go
package panels

import (
    "github.com/tonysyu/gqlxp/tui/components"
    "github.com/tonysyu/gqlxp/tui/config"
)

// Manager handles panel lifecycle operations
type Manager struct {
    styles     config.Styles
    layoutCfg  config.LayoutConfig
}

func NewManager(styles config.Styles, layoutCfg config.LayoutConfig) *Manager {
    return &Manager{
        styles:    styles,
        layoutCfg: layoutCfg,
    }
}

// CreateListPanel creates a new list panel with standard configuration
func (m *Manager) CreateListPanel(
    items []components.ListItem,
    title string,
    opts ...PanelOption,
) *components.ListPanel {
    panel := components.NewListPanel(items, title)

    // Apply options
    for _, opt := range opts {
        opt(panel)
    }

    return panel
}

// CreateEmptyPanel creates an empty placeholder panel
func (m *Manager) CreateEmptyPanel() components.Panel {
    return components.NewStringPanel("")
}

// Size calculates and applies size to a panel
func (m *Manager) Size(panel components.Panel, spec SizeSpec) {
    width := m.calculateWidth(spec)
    height := m.calculateHeight(spec)
    panel.SetSize(width, height)
}

// SizePanels sizes multiple panels in a layout
func (m *Manager) SizePanels(panels []PanelSizeSpec) {
    for _, spec := range panels {
        m.Size(spec.Panel, spec.SizeSpec)
    }
}

// Focus sets a panel as focused
func (m *Manager) Focus(panel components.Panel) {
    if listPanel, ok := panel.(*components.ListPanel); ok {
        listPanel.SetFocused()
    }
}

// Blur sets a panel as blurred
func (m *Manager) Blur(panel components.Panel) {
    if listPanel, ok := panel.(*components.ListPanel); ok {
        listPanel.SetBlurred()
    }
}

// BlurAll blurs all panels in a slice
func (m *Manager) BlurAll(panels []components.Panel) {
    for _, panel := range panels {
        m.Blur(panel)
    }
}

// UpdateFocus updates focus state for a set of panels
func (m *Manager) UpdateFocus(panels []components.Panel, focusedIndex int) {
    m.BlurAll(panels)
    if focusedIndex >= 0 && focusedIndex < len(panels) {
        m.Focus(panels[focusedIndex])
    }
}

func (m *Manager) calculateWidth(spec SizeSpec) int {
    baseWidth := spec.TotalWidth / spec.ColumnsCount
    frameSize := 0

    if spec.IsFocused {
        frameSize = m.styles.FocusedPanel.GetHorizontalFrameSize()
    } else {
        frameSize = m.styles.BlurredPanel.GetHorizontalFrameSize()
    }

    return baseWidth - frameSize
}

func (m *Manager) calculateHeight(spec SizeSpec) int {
    height := spec.TotalHeight
    height -= m.layoutCfg.HelpHeight
    height -= m.layoutCfg.NavbarHeight
    height -= m.layoutCfg.BreadcrumbsHeight

    frameSize := 0
    if spec.IsFocused {
        frameSize = m.styles.FocusedPanel.GetVerticalFrameSize()
    } else {
        frameSize = m.styles.BlurredPanel.GetVerticalFrameSize()
    }

    return height - frameSize
}

// SizeSpec specifies how to size a panel
type SizeSpec struct {
    TotalWidth   int
    TotalHeight  int
    ColumnsCount int
    IsFocused    bool
}

// PanelSizeSpec combines a panel with its size spec
type PanelSizeSpec struct {
    Panel    components.Panel
    SizeSpec SizeSpec
}

// PanelOption is a functional option for customizing panels
type PanelOption func(*components.ListPanel)

// WithDescription sets panel description
func WithDescription(desc string) PanelOption {
    return func(p *components.ListPanel) {
        p.SetDescription(desc)
    }
}

// WithObjectType sets panel object type
func WithObjectType(item components.ListItem) PanelOption {
    return func(p *components.ListPanel) {
        p.SetObjectType(item)
    }
}
```

### Panel Factory

Create convenience factories for common panel types:

```go
// tui/panels/factory.go
package panels

import (
    "github.com/tonysyu/gqlxp/gql"
    "github.com/tonysyu/gqlxp/tui/components"
)

// Factory creates commonly-used panel types
type Factory struct {
    manager  *Manager
    resolver gql.TypeResolver
}

func NewFactory(manager *Manager, resolver gql.TypeResolver) *Factory {
    return &Factory{
        manager:  manager,
        resolver: resolver,
    }
}

// CreateFieldPanel creates a panel for displaying a field's arguments and result type
func (f *Factory) CreateFieldPanel(field *gql.Field) *components.ListPanel {
    argumentItems := adaptArgumentsToItems(field.Arguments(), f.resolver)

    panel := f.manager.CreateListPanel(
        argumentItems,
        field.Name(),
        WithDescription(field.Description()),
        WithObjectType(newTypeDefItemFromField(field, f.resolver)),
    )

    return panel
}

// CreateTypeDefPanel creates a panel for displaying a type definition
func (f *Factory) CreateTypeDefPanel(typeDef gql.TypeDef) *components.ListPanel {
    var items []components.ListItem

    switch td := typeDef.(type) {
    case *gql.Object:
        items = adaptFieldsToItems(td.Fields(), f.resolver)
    case *gql.Interface:
        items = adaptFieldsToItems(td.Fields(), f.resolver)
    case *gql.Union:
        items = adaptUnionTypesToItems(td.Types(), f.resolver)
    case *gql.Enum:
        items = adaptEnumValuesToItems(td.Values())
    case *gql.InputObject:
        items = adaptFieldsToItems(td.Fields(), f.resolver)
    }

    panel := f.manager.CreateListPanel(
        items,
        typeDef.Name(),
        WithDescription(typeDef.Description()),
    )

    return panel
}
```

### Update mainModel

Simplify `mainModel` using panel manager:

```go
// tui/model.go
type mainModel struct {
    schema       adapters.SchemaView
    panelManager *panels.Manager
    panelFactory *panels.Factory
    // ... rest
}

func newModel(schema adapters.SchemaView) mainModel {
    styles := config.DefaultStyles()
    layoutCfg := config.DefaultLayoutConfig()

    panelMgr := panels.NewManager(styles, layoutCfg)
    resolver := gql.NewSchemaResolver(&schema.schema)
    panelFactory := panels.NewFactory(panelMgr, resolver)

    m := mainModel{
        schema:       schema,
        panelManager: panelMgr,
        panelFactory: panelFactory,
        styles:       styles,
        // ... rest
    }

    return m
}

func (m *mainModel) sizePanels() {
    specs := []panels.PanelSizeSpec{
        {
            Panel: m.panelStack[m.stackPosition],
            SizeSpec: panels.SizeSpec{
                TotalWidth:   m.width,
                TotalHeight:  m.height,
                ColumnsCount: config.VisiblePanelCount,
                IsFocused:    true,
            },
        },
    }

    if len(m.panelStack) > m.stackPosition+1 {
        specs = append(specs, panels.PanelSizeSpec{
            Panel: m.panelStack[m.stackPosition+1],
            SizeSpec: panels.SizeSpec{
                TotalWidth:   m.width,
                TotalHeight:  m.height,
                ColumnsCount: config.VisiblePanelCount,
                IsFocused:    false,
            },
        })
    }

    m.panelManager.SizePanels(specs)
}

func (m *mainModel) updatePanelFocusStates() {
    m.panelManager.UpdateFocus(m.panelStack, m.stackPosition)
}

func (m *mainModel) resetAndLoadMainPanel() {
    // Create empty panels
    m.panelStack = make([]components.Panel, config.VisiblePanelCount)
    for i := range config.VisiblePanelCount {
        m.panelStack[i] = m.panelManager.CreateEmptyPanel()
    }

    m.stackPosition = 0
    m.breadcrumbs.Reset()
    m.loadMainPanel()
}
```

## Benefits

1. **Separation of concerns**: Panel operations isolated
2. **Testability**: Easy to test panel lifecycle independently
3. **Reusability**: Panel manager can be used in other contexts
4. **Extensibility**: Easy to add animations, transitions
5. **Consistency**: Panel creation follows consistent patterns
6. **Maintainability**: Panel logic in one place

## Implementation Steps

1. Create `tui/panels/` package
2. Implement `Manager` with sizing and focus methods
3. Implement `Factory` for common panel types
4. Add tests for manager and factory
5. Update `mainModel` to use panel manager
6. Remove panel lifecycle code from `mainModel`
7. Run tests: `just test`
8. Update documentation

## Testing Strategy

```go
// tui/panels/manager_test.go
func TestManager_CreateListPanel(t *testing.T) {
    mgr := NewManager(config.DefaultStyles(), config.DefaultLayoutConfig())

    items := []components.ListItem{
        components.NewSimpleItem("item1"),
    }

    panel := mgr.CreateListPanel(items, "Test Panel")

    assert.NotNil(t, panel)
    assert.Equal(t, "Test Panel", panel.Title())
}

func TestManager_Size(t *testing.T) {
    mgr := NewManager(config.DefaultStyles(), config.DefaultLayoutConfig())
    panel := mgr.CreateEmptyPanel()

    spec := panels.SizeSpec{
        TotalWidth:   100,
        TotalHeight:  50,
        ColumnsCount: 2,
        IsFocused:    true,
    }

    mgr.Size(panel, spec)

    // Verify panel was sized (would need to expose size or verify via render)
}

func TestManager_UpdateFocus(t *testing.T) {
    mgr := NewManager(config.DefaultStyles(), config.DefaultLayoutConfig())

    panels := []components.Panel{
        components.NewListPanel([]components.ListItem{}, "Panel 1"),
        components.NewListPanel([]components.ListItem{}, "Panel 2"),
    }

    mgr.UpdateFocus(panels, 1)

    // Verify panel 1 is focused (would check internal state)
}
```

## Potential Issues

- **Abstraction overhead**: Adds layer between model and panels
- **Limited benefits**: May not provide enough value for complexity

## Future Enhancements

1. **Panel animations**: Fade in/out, slide transitions
2. **Panel caching**: Cache rendered panels for performance
3. **Panel templates**: Predefined panel layouts
4. **Panel pooling**: Reuse panel instances
5. **Custom panel types**: Easy registration of new panel types

## Related Tasks

- **Task 02** (Navigation Manager): Manager could integrate with navigation
- **Task 05** (Command Pattern): Commands could use panel manager
- **Task 09** (Configuration Layer): Load panel styling from config
