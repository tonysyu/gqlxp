# Change: Add Schema Library with Persistent Storage and Metadata

## Why
Currently, gqlxp only supports exploring GraphQL schemas provided via command-line file paths. Users must remember file locations and cannot persist preferences or metadata about schemas. Adding a schema library with persistent storage enables users to manage a collection of schemas with associated metadata like favorited types and URL patterns for web views.

## What Changes
- Store GraphQL schemas in `~/.config/gqlxp/schemas/`
- Store schema metadata in `~/.config/gqlxp/metadata/`
- Support schema-specific metadata:
  - Display names for schemas
  - Favorited type names for quick access
  - URL patterns to open web documentation for selected types
- Add commands/UI for managing schema library (add, list, remove, update metadata)
- Support both file-path mode (current) and library mode (new)

## Impact
- Affected specs: `schema-library` (new capability)
- Affected code:
  - `cmd/gqlxp/main.go` - Add library mode support
  - New package `library` - Schema and metadata storage/retrieval
  - `tui/` - UI components for library management and favorites
  - Configuration file handling for `~/.config/gqlxp/`
