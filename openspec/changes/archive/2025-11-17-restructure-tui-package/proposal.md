# Restructure TUI Package

## Problem

The `tui` package currently mixes three distinct concerns in a flat structure:
- Schema library selection (schema_selector.go)
- Schema exploration interface (model.go, favoritable_item.go, app.go)
- Overlay display (overlay.go)

This creates organizational confusion and makes it unclear which code handles which UI mode. The current `tui/model.go` is specifically the schema explorer model, but its name suggests it's the top-level TUI model.

## Goals

1. Separate the three UI modes into distinct subpackages with clear responsibilities
2. Create a true top-level `tui/model.go` that delegates to the appropriate submode
3. Maintain backward compatibility with existing entry points (`tui.Start`, `tui.StartWithLibraryData`, `tui.StartSchemaSelector`)
4. Keep truly shared utilities (`adapters/`, `config/`) at the `tui` level
5. Move explorer-specific packages (`components/`, `navigation/`) under `xplr/`

## Non-Goals

- Changing the public API signatures for `tui.Start*` functions
- Moving or restructuring `adapters/` or `config/` packages (these remain shared)
- Adding new features or capabilities beyond the refactoring

## Proposed Solution

Restructure the `tui` package into three subpackages:

```
tui/
├── model.go              # Top-level model that delegates to submodes
├── app.go                # Entry points (Start, StartWithLibraryData, StartSchemaSelector)
├── adapters/             # GraphQL AST to UI component conversion (shared)
├── config/               # Styles and configuration (shared)
├── libselect/            # Library selection mode
│   └── model.go          # Schema selector (from schema_selector.go)
├── xplr/                 # Schema exploration mode
│   ├── model.go          # Main explorer model (from current model.go)
│   ├── favoritable_item.go
│   ├── components/       # UI components specific to explorer
│   └── navigation/       # Navigation state management for explorer
└── overlay/              # Overlay display mode
    └── model.go          # Overlay model (from overlay.go)
```

Shared packages at `tui/` level:
- `adapters/` - GraphQL AST to UI component conversion (used by xplr and potentially libselect)
- `config/` - Styles and configuration (used by all submodes)

Explorer-specific packages under `tui/xplr/`:
- `components/` - Panels and list items for schema navigation
- `navigation/` - Panel stack, breadcrumbs, and type selection state

The top-level `tui/model.go` will implement a state machine that:
1. Starts in the appropriate mode based on entry point
2. Handles transitions between modes (e.g., library selection → schema exploration)
3. Delegates all UI logic to the active submode

## Open Questions

None - this is a straightforward refactoring with clear boundaries.

## Impact

- All files in `tui/` (except `app.go`, `adapters/`, and `config/`) will be moved to subpackages
- `components/` and `navigation/` will move under `xplr/` as they're explorer-specific
- Tests will be moved alongside their corresponding implementation files
- Import paths will change for internal TUI code, but public API remains unchanged
- No impact on `cmd/gqlxp` - entry points remain the same
