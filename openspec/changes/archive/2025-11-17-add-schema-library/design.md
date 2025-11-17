# Design: Schema Library with Persistent Storage

## Context
gqlxp currently operates as a single-session tool that loads schemas from file paths. Users want to maintain a library of schemas with persistent metadata (favorites, URL patterns) to enhance their workflow. This requires file-based storage in the user's config directory following XDG conventions.

## Goals / Non-Goals

### Goals
- Store schemas and metadata in `~/.config/gqlxp/`
- Support schema-specific metadata (display names, favorites, URL patterns)
- Maintain backward compatibility with existing file-path mode
- Simple file-based storage (no database required)
- Cross-platform support (macOS, Linux, Windows)

### Non-Goals
- Cloud synchronization or remote storage
- Versioning or history tracking of schemas
- Real-time schema fetching from GraphQL endpoints
- Multi-user or shared library support

## Decisions

### Decision 1: File-Based Storage Format
**What**: Use file system with single JSON metadata file
- Schemas stored as `.graphqls` files in `~/.config/gqlxp/schemas/<schema-id>.graphqls`
- All metadata stored in single JSON file at `~/.config/gqlxp/schemas/metadata.json` with schema-id as top-level keys

**Why**: Simple, portable, human-readable, no external dependencies; easier to manage than multiple files for small libraries

**Alternatives considered**:
- Per-schema metadata files: More file operations; harder to enumerate all schemas
- SQLite database: Overkill for simple key-value storage; adds dependency

### Decision 2: Schema ID Format
**What**: Use sanitized schema name as ID (lowercase, alphanumeric + hyphens)

**Why**: Human-readable filenames, URL-safe, easy to reference

**Alternatives considered**:
- UUID: Less readable, harder to debug
- Hash of schema content: Not stable across metadata-only updates

### Decision 3: Metadata Schema
**What**: Single JSON file (`~/.config/gqlxp/schemas/metadata.json`) with schema-id as top-level keys:
```json
{
  "github-api": {
    "displayName": "GitHub GraphQL API",
    "sourceFile": "/path/to/original/github.graphqls",
    "favorites": ["Query", "Repository", "User"],
    "urlPatterns": {
      "Query": "https://docs.github.com/graphql/reference/queries#${field}",
      "Mutation": "https://docs.github.com/graphql/reference/mutations#${field}",
      "*": "https://docs.github.com/graphql/reference/objects#${type}"
    },
    "createdAt": "2025-11-13T10:30:00Z",
    "updatedAt": "2025-11-13T10:30:00Z"
  },
  "shopify-api": {
    "displayName": "Shopify Storefront API",
    "sourceFile": "/path/to/shopify.graphqls",
    "favorites": ["Product", "Collection"],
    "urlPatterns": {
      "*": "https://shopify.dev/docs/api/storefront#${type}"
    },
    "createdAt": "2025-11-13T11:00:00Z",
    "updatedAt": "2025-11-13T11:00:00Z"
  }
}
```

**Why**:
- Single file easier to manage for small libraries
- Schema-id as key enables quick lookup
- URL patterns use simple template syntax (`${type}`, `${field}`)
- Wildcards (`*`) provide fallback patterns
- Timestamps for debugging and potential future sorting

### Decision 4: Config Directory Resolution
**What**: Use standard config directory per platform:
- macOS/Linux: `$HOME/.config/gqlxp/`
- Windows: `%APPDATA%\gqlxp\`
- Fallback: `$HOME/.gqlxp/` if config dir unavailable

**Why**: Follows XDG Base Directory specification and platform conventions

**Implementation**: Use Go standard library or minimal helper for cross-platform paths

### Decision 5: Backward Compatibility
**What**: Maintain existing CLI behavior:
- `gqlxp <file-path>` - Direct file mode (current behavior)
- `gqlxp --library` - Library selection mode (new)
- `gqlxp --library <schema-id>` - Direct library schema load (new)

**Why**: Preserves existing workflows; gradual adoption of library features

## Risks / Trade-offs

### Risk: Config Directory Permissions
**Mitigation**: Check permissions on startup; provide clear error messages; graceful fallback to read-only mode

### Risk: File Conflicts from Concurrent Access
**Mitigation**: Single-user assumption; advisory file locking for future enhancement

### Risk: Metadata Corruption
**Mitigation**: Atomic writes with temp file + rename; JSON validation on load

### Trade-off: Simple Storage vs. Query Performance
**Decision**: Accept linear scan for small libraries (<100 schemas)
**Rationale**: Simplicity over premature optimization; add indexing if needed

## Migration Plan

### Phase 1: Core Library (MVP)
1. Implement storage layer with schema and metadata persistence
2. Add library commands (add, list, remove)
3. Basic CLI integration with `--library` flag

### Phase 2: Metadata Features
1. Implement favorites tracking
2. Implement URL pattern storage
3. Add metadata update commands

### Phase 3: TUI Integration (Future)
1. Schema selector screen
2. Favorites indicators
3. URL opening on hotkey

### Rollback
- No migration of existing data required
- Library feature is additive; can be disabled by not using flags
- Remove `~/.config/gqlxp/` directory to reset

## Open Questions
- Should URL opening be automatic or manual (hotkey)?
- Should favorites affect default sorting or just visual indicators?
- Should we support importing schemas from URLs in the future?
