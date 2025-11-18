# Design: Automatic Library Integration

## Architecture Overview

### Current Flow
```
gqlxp <file>           → Parse file → Start TUI (no library)
gqlxp --library        → Select from library → Start TUI
gqlxp library add ...  → Add to library (separate command)
```

### Proposed Flow
```
gqlxp                  → Check library → Select or Error
gqlxp <file>           → Check hash/path → Prompt if needed → Save to library → Start TUI
```

## Key Components

### 1. File Hash Storage
**Decision**: Use SHA-256 hash of file contents
**Rationale**: Standard, fast, collision-resistant, available in Go stdlib

**Storage**: Add to `SchemaMetadata` struct:
```go
type SchemaMetadata struct {
    DisplayName string
    SourceFile  string
    FileHash    string  // NEW: SHA-256 hash of schema content
    // ... existing fields
}
```

### 2. File Path Matching
**Decision**: Store absolute file paths
**Rationale**: Unambiguous, works with symlinks, handles relative paths correctly

**Normalization**: Convert all paths to absolute during storage and comparison

### 3. Interactive Prompting
**Decision**: Use standard input/output before TUI starts
**Rationale**: Simple, no additional dependencies, clear separation of concerns

**Prompts**:
1. Schema ID: "Enter schema ID (lowercase letters, numbers, hyphens): "
2. Display name: "Enter display name [default: <id>]: "
3. Hash mismatch: "Schema file has changed. Update library? (y/n): "

### 4. Match Detection Logic
```
Load file → Calculate hash → Query library by path

Match cases:
1. No path match → Prompt for ID + display name → Save → Continue
2. Path + hash match → Load from library → Continue
3. Path match, hash mismatch → Prompt update/load → Act accordingly
```

### 5. CLI Simplification
**Remove**: All `library` subcommands (`add`, `list`, `remove`) - replaced by automatic behavior and TUI-based management
**Change**: `gqlxp` with no args → open selector if library not empty (selector provides list/remove functionality)

## Trade-offs

### Automatic vs Explicit
- **Pro**: Reduced friction, simpler mental model, single workflow
- **Con**: Less control, might prompt unexpectedly
- **Mitigation**: Clear prompts, ability to skip updates

### Hash Storage Overhead
- **Pro**: Minimal storage cost (~64 bytes per schema)
- **Con**: Requires rehashing on every file load
- **Mitigation**: SHA-256 is fast enough for schema files (<1ms typically)

### Breaking Changes
- **Removed**: All `library` CLI subcommands (`add`, `list`, `remove`)
- **Changed**: CLI interface automatically manages library storage
- **Impact**: Users relying on `library` subcommands must use TUI selector for list/remove operations
- **Mitigation**: Document in changelog, TUI selector provides full library management capabilities

## Implementation Phases

### Phase 1: Metadata Enhancement
- Add `FileHash` to `SchemaMetadata`
- Implement hash calculation function
- Update `Library.Add()` to store hash and absolute path

### Phase 2: Match Detection
- Implement path/hash lookup function
- Add prompting utilities for user input
- Create interactive flow for schema registration

### Phase 3: CLI Integration
- Refactor `main.go` startup flow
- Implement automatic save-on-load behavior
- Handle no-args case with selector

### Phase 4: Cleanup
- Remove `library` subcommand handling from `main.go`
- Update documentation to remove references to `library` commands
- Add tests for new flows
- Ensure TUI selector supports schema removal functionality
