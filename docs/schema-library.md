# Schema Library

## Overview

Schema library provides persistent storage for GraphQL schemas with associated metadata like favorites and URL patterns.

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
- `sourceFile`: Original file path (for reference)
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
}
```

### Implementation Details

**File-based storage**: Simple, portable, no external dependencies
**Atomic writes**: Metadata updates use temp file + rename
**Single metadata file**: Easier to manage than per-schema files for small libraries
**Config directory**: Follows XDG Base Directory specification

## CLI Usage

```sh
# Add schema
gqlxp library add <schema-id> <file-path>

# List schemas
gqlxp library list

# Load from library
gqlxp --library <schema-id>

# Remove schema
gqlxp library remove <schema-id>
```

## Future Enhancements

**Not yet implemented:**
- Interactive TUI schema selector
- Favorites visual indicators
- URL opening for selected types
- Display name customization via CLI
- URL pattern management via CLI
