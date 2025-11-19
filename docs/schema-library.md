# Schema Library

## Overview

Schema library provides persistent storage for GraphQL schemas with automatic integration, file hash tracking, and metadata like favorites and URL patterns. Schemas are automatically saved to the library when loaded from files.

## Storage Structure

```
~/.config/gqlxp/
└── schemas/
    ├── metadata.json
    ├── github-api.graphqls
    └── shopify-api.graphqls
```

**Directory Locations:**
- macOS/Linux: `~/.config/gqlxp/schemas/`
- Windows: `%APPDATA%\gqlxp\schemas\`

## Metadata Schema

All schema metadata is stored in `schemas/metadata.json` with schema-id as top-level keys:

```json
{
  "github-api": {
    "displayName": "GitHub GraphQL API",
    "sourceFile": "/path/to/original/github.graphqls",
    "fileHash": "a3f2c8b1...",
    "favorites": ["Query", "Repository", "User"],
    "urlPatterns": {
      "Query": "https://docs.github.com/graphql/reference/queries#${field}",
      "Mutation": "https://docs.github.com/graphql/reference/mutations#${field}",
      "*": "https://docs.github.com/graphql/reference/objects#${type}"
    },
    "createdAt": "2025-11-13T10:30:00Z",
    "updatedAt": "2025-11-13T10:30:00Z"
  }
}
```

**Fields:**
- `displayName`: Human-readable schema name
- `sourceFile`: Absolute path to original file
- `fileHash`: SHA-256 hash of schema content (for change detection)
- `favorites`: List of favorited type names
- `urlPatterns`: URL templates for documentation links
  - Type-specific patterns (e.g., "Query", "Mutation")
  - Wildcard pattern (`*`) as fallback
  - Template variables: `${type}`, `${field}`
- `createdAt`: Schema creation timestamp
- `updatedAt`: Last metadata update timestamp

## Schema ID Format

Schema IDs must:
- Contain only lowercase letters, numbers, and hyphens
- Be valid filenames across platforms
- Be unique within library

**Examples:**
- Valid: `github-api`, `shopify-v2`, `internal-api`
- Invalid: `GitHub_API`, `api v2`, `my@api`

Use `SanitizeSchemaID()` to convert invalid IDs.

## Architecture

### Package Structure

```
library/
├── types.go       # Core data types
├── config.go      # Config directory resolution
├── library.go     # Library interface and implementation
└── library_test.go
```

### Core Types

**SchemaMetadata**: Metadata for stored schema
**Schema**: Schema content with metadata
**SchemaInfo**: Basic schema info for listing

### Interface

```go
type Library interface {
    Add(id, displayName, sourcePath string) error
    Get(id string) (*Schema, error)
    List() ([]SchemaInfo, error)
    Remove(id string) error
    UpdateMetadata(id string, metadata SchemaMetadata) error
    AddFavorite(id, typeName string) error
    RemoveFavorite(id, typeName string) error
    SetURLPattern(id, typePattern, urlPattern string) error
    FindByPath(absolutePath string) (*Schema, error)
    UpdateContent(id string, content []byte) error
}
```

### Implementation Details

**File-based storage**: Simple, portable, no external dependencies
**Atomic writes**: Metadata updates use temp file + rename
**Single metadata file**: Easier to manage than per-schema files for small libraries
**Config directory**: Follows XDG Base Directory specification

## CLI Usage

```sh
# Load schema file (automatically saved to library on first use)
gqlxp examples/github.graphqls
# Prompts for schema ID and display name on first use
# Detects changes and prompts for update on subsequent uses

# Open library selector (when no arguments provided)
gqlxp

# Note: Schemas can be removed via the TUI selector interface
```

## TUI Features

All schemas are now library-backed with access to:

**Schema Selector** (`gqlxp` with no args)
- Interactive list of all schemas in library
- Filter/search by schema ID or display name
- Enter to select and open schema
- Delete key to remove schemas from library

**Favorites** (Press `f`)
- Mark types as favorites for quick identification
- Favorited types show ★ indicator in type lists
- Favorites persist across sessions in `metadata.json`
- Toggle on/off by pressing `f` on selected type

## Automatic Library Integration

**Workflow:**
1. Load schema file with `gqlxp <file-path>`
2. System checks library for existing entry by file path
3. If new: prompts for schema ID and display name, saves to library
4. If exists: compares file hash
   - Hash match: loads existing library version
   - Hash mismatch: prompts to update library or use existing version
5. All TUI sessions are library-backed with full metadata support
