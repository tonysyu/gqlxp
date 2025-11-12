# Task 09: Configuration Layer

**Priority:** Low
**Status:** Not Started
**Estimated Effort:** Small-Medium
**Dependencies:** None

## Problem Statement

Configuration values are scattered throughout the codebase:

- Magic numbers in code (e.g., `maxDescriptionHeight = 5`)
- Constants in `tui/config/config.go` (e.g., `VisiblePanelCount = 2`)
- Hardcoded sizes (e.g., `HelpHeight`, `NavbarHeight`)
- No way for users to customize behavior
- Difficult to find and change configuration values

### Affected Files
- `tui/config/config.go` - Partial configuration
- `tui/components/panels.go` - Magic numbers
- `tui/model.go` - Layout calculations
- Various files with hardcoded values

## Proposed Solution

Create a comprehensive configuration layer:

### 1. Layout Configuration

```go
// tui/config/layout.go
package config

// LayoutConfig defines layout and sizing configuration
type LayoutConfig struct {
    // Panel configuration
    VisiblePanelCount    int
    MaxPanelStackSize    int
    MaxDescriptionHeight int

    // Header/footer heights
    HelpHeight        int
    NavbarHeight      int
    BreadcrumbsHeight int

    // Panel sizing
    PanelMinWidth  int
    PanelMinHeight int
}

// DefaultLayoutConfig returns sensible defaults
func DefaultLayoutConfig() LayoutConfig {
    return LayoutConfig{
        VisiblePanelCount:    2,
        MaxPanelStackSize:    20,
        MaxDescriptionHeight: 5,
        HelpHeight:           1,
        NavbarHeight:         1,
        BreadcrumbsHeight:    1,
        PanelMinWidth:        40,
        PanelMinHeight:       10,
    }
}

// Validate checks if configuration is valid
func (c LayoutConfig) Validate() error {
    if c.VisiblePanelCount < 1 {
        return fmt.Errorf("VisiblePanelCount must be >= 1, got %d", c.VisiblePanelCount)
    }
    if c.MaxDescriptionHeight < 0 {
        return fmt.Errorf("MaxDescriptionHeight must be >= 0, got %d", c.MaxDescriptionHeight)
    }
    return nil
}
```

### 2. Behavior Configuration

```go
// tui/config/behavior.go
package config

// BehaviorConfig defines application behavior configuration
type BehaviorConfig struct {
    // Auto-open behavior
    AutoOpenFirstItem    bool
    AutoOpenOnNavigation bool

    // Filtering
    EnableFiltering     bool
    CaseSensitiveFilter bool

    // Navigation
    WrapAroundNavigation bool // Wrap to start when reaching end

    // History
    EnableHistory    bool
    MaxHistorySize   int
    EnableUndoRedo   bool

    // Performance
    LazyLoadPanels   bool
    CacheRenderedPanels bool
}

func DefaultBehaviorConfig() BehaviorConfig {
    return BehaviorConfig{
        AutoOpenFirstItem:    true,
        AutoOpenOnNavigation: true,
        EnableFiltering:      true,
        CaseSensitiveFilter:  false,
        WrapAroundNavigation: true,
        EnableHistory:        false,
        MaxHistorySize:       100,
        EnableUndoRedo:       false,
        LazyLoadPanels:       false,
        CacheRenderedPanels:  false,
    }
}

func (c BehaviorConfig) Validate() error {
    if c.MaxHistorySize < 0 {
        return fmt.Errorf("MaxHistorySize must be >= 0, got %d", c.MaxHistorySize)
    }
    return nil
}
```

### 3. Display Configuration

```go
// tui/config/display.go
package config

// DisplayConfig defines display and formatting configuration
type DisplayConfig struct {
    // Signature formatting
    MaxSignatureWidth int // Signatures longer than this use multiline format

    // Overlay
    OverlayWidth  int // Percentage of screen width (0-100)
    OverlayHeight int // Percentage of screen height (0-100)

    // List display
    ShowItemDescriptions bool
    ShowLineNumbers      bool

    // Detail display
    ShowTypeDetails      bool
    ShowFieldSignatures  bool
}

func DefaultDisplayConfig() DisplayConfig {
    return DisplayConfig{
        MaxSignatureWidth:    80,
        OverlayWidth:         80,
        OverlayHeight:        80,
        ShowItemDescriptions: true,
        ShowLineNumbers:      false,
        ShowTypeDetails:      true,
        ShowFieldSignatures:  true,
    }
}

func (c DisplayConfig) Validate() error {
    if c.OverlayWidth < 10 || c.OverlayWidth > 100 {
        return fmt.Errorf("OverlayWidth must be 10-100, got %d", c.OverlayWidth)
    }
    if c.OverlayHeight < 10 || c.OverlayHeight > 100 {
        return fmt.Errorf("OverlayHeight must be 10-100, got %d", c.OverlayHeight)
    }
    return nil
}
```

### 4. Keybinding Configuration

```go
// tui/config/keybindings.go
package config

// KeybindingConfig defines customizable keybindings
type KeybindingConfig struct {
    // Navigation
    NextPanel     []string
    PrevPanel     []string
    NextType      []string
    PrevType      []string

    // Actions
    ToggleOverlay []string
    Quit          []string
    Undo          []string
    Redo          []string

    // Search
    StartFilter   []string
    ClearFilter   []string
}

func DefaultKeybindingConfig() KeybindingConfig {
    return KeybindingConfig{
        NextPanel:     []string{"tab", "]"},
        PrevPanel:     []string{"shift+tab", "["},
        NextType:      []string{"ctrl+t", "}"},
        PrevType:      []string{"ctrl+r", "{"},
        ToggleOverlay: []string{" "},
        Quit:          []string{"ctrl+c", "ctrl+d"},
        Undo:          []string{"ctrl+z"},
        Redo:          []string{"ctrl+y"},
        StartFilter:   []string{"/"},
        ClearFilter:   []string{"esc"},
    }
}
```

### 5. Main Configuration

Combine all configuration:

```go
// tui/config/config.go
package config

import (
    "encoding/json"
    "os"
    "path/filepath"
)

// Config is the main application configuration
type Config struct {
    Layout      LayoutConfig
    Behavior    BehaviorConfig
    Display     DisplayConfig
    Keybindings KeybindingConfig
    Styles      Styles // Existing styles config
}

// DefaultConfig returns default configuration
func DefaultConfig() Config {
    return Config{
        Layout:      DefaultLayoutConfig(),
        Behavior:    DefaultBehaviorConfig(),
        Display:     DefaultDisplayConfig(),
        Keybindings: DefaultKeybindingConfig(),
        Styles:      DefaultStyles(),
    }
}

// Validate validates all configuration sections
func (c Config) Validate() error {
    if err := c.Layout.Validate(); err != nil {
        return fmt.Errorf("layout config: %w", err)
    }
    if err := c.Behavior.Validate(); err != nil {
        return fmt.Errorf("behavior config: %w", err)
    }
    if err := c.Display.Validate(); err != nil {
        return fmt.Errorf("display config: %w", err)
    }
    return nil
}

// LoadConfig loads configuration from file
func LoadConfig(path string) (Config, error) {
    // Try to load from file
    data, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            // File doesn't exist, use defaults
            return DefaultConfig(), nil
        }
        return Config{}, fmt.Errorf("reading config: %w", err)
    }

    var config Config
    if err := json.Unmarshal(data, &config); err != nil {
        return Config{}, fmt.Errorf("parsing config: %w", err)
    }

    if err := config.Validate(); err != nil {
        return Config{}, fmt.Errorf("invalid config: %w", err)
    }

    return config, nil
}

// SaveConfig saves configuration to file
func (c Config) SaveConfig(path string) error {
    data, err := json.MarshalIndent(c, "", "  ")
    if err != nil {
        return fmt.Errorf("marshaling config: %w", err)
    }

    // Ensure directory exists
    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("creating config directory: %w", err)
    }

    if err := os.WriteFile(path, data, 0644); err != nil {
        return fmt.Errorf("writing config: %w", err)
    }

    return nil
}

// ConfigPath returns the default config file path
func ConfigPath() (string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    return filepath.Join(homeDir, ".config", "gqlxp", "config.json"), nil
}
```

### 6. Update Application

Use configuration throughout app:

```go
// cmd/gqlxp/main.go
func main() {
    // ... existing code ...

    // Load configuration
    configPath, err := config.ConfigPath()
    if err != nil {
        // Use defaults
        cfg = config.DefaultConfig()
    } else {
        cfg, err = config.LoadConfig(configPath)
        if err != nil {
            fmt.Printf("Warning: failed to load config: %v\n", err)
            cfg = config.DefaultConfig()
        }
    }

    // Start TUI with config
    if _, err := tui.StartWithConfig(schema, cfg); err != nil {
        abort(fmt.Sprintf("Error starting tui: %v", err))
    }
}

// tui/app.go
func StartWithConfig(schema adapters.SchemaView, cfg config.Config) (tea.Model, error) {
    model := newModelWithConfig(schema, cfg)
    p := tea.NewProgram(model, tea.WithAltScreen())
    return p.Run()
}

// tui/model.go
func newModelWithConfig(schema adapters.SchemaView, cfg config.Config) mainModel {
    m := mainModel{
        schema:  schema,
        config:  cfg,
        Styles:  cfg.Styles,
        // ... use cfg.Layout, cfg.Behavior, etc.
    }
    // ...
}
```

## Benefits

1. **Centralization**: All config in one place
2. **Customization**: Users can customize behavior
3. **Validation**: Config validated on load
4. **Defaults**: Sensible defaults provided
5. **Persistence**: Save/load user preferences
6. **Documentation**: Config file self-documenting

## Implementation Steps

1. Create configuration types in `tui/config/`
2. Implement `LoadConfig` and `SaveConfig`
3. Add configuration validation
4. Update application to use config
5. Replace magic numbers with config values
6. Add tests for configuration
7. Create example config file
8. Update documentation
9. Run tests: `just test`

## Testing Strategy

```go
// tui/config/layout_test.go
func TestLayoutConfig_Validate(t *testing.T) {
    cfg := DefaultLayoutConfig()
    assert.NoError(t, cfg.Validate())

    cfg.VisiblePanelCount = 0
    assert.Error(t, cfg.Validate())
}

// tui/config/config_test.go
func TestLoadConfig(t *testing.T) {
    // Create temp config file
    tmpfile := createTempConfigFile(t, `{
        "Layout": {"VisiblePanelCount": 3},
        "Behavior": {"AutoOpenFirstItem": false}
    }`)
    defer os.Remove(tmpfile)

    cfg, err := LoadConfig(tmpfile)

    assert.NoError(t, err)
    assert.Equal(t, 3, cfg.Layout.VisiblePanelCount)
    assert.False(t, cfg.Behavior.AutoOpenFirstItem)
}

func TestLoadConfig_NotFound(t *testing.T) {
    cfg, err := LoadConfig("/nonexistent/config.json")

    assert.NoError(t, err) // Should use defaults
    assert.Equal(t, DefaultConfig(), cfg)
}
```

## Potential Issues

- **Configuration complexity**: Too many options can be overwhelming
- **Backward compatibility**: Config format changes need migration
- **Validation**: Need comprehensive validation

## Future Enhancements

1. **Config UI**: TUI for editing configuration
2. **Profiles**: Multiple named configuration profiles
3. **Environment overrides**: Override config with env vars
4. **Schema-specific config**: Different config per schema
5. **Hot reload**: Reload config without restart
6. **Config templates**: Predefined configuration templates

## Related Tasks

- **Task 02** (Navigation Manager): Use config for navigation behavior
- **Task 03** (Type Registry): Load type visibility from config
- **Task 05** (Command Pattern): Load keybindings from config
- **Task 06** (Panel Lifecycle): Use config for panel sizing
