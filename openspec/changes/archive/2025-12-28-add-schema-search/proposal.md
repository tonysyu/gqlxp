# Change: Add Schema Search

## Why
Large GraphQL schemas are difficult to navigate when looking for specific types, fields, or documentation. Users need efficient search capabilities to quickly find types and fields by name or description across their schema library without manually scanning through the TUI.

## What Changes
- New `search` package positioned between `library` and `gql` in the package hierarchy
- New `gqlxp search <query>` CLI command for searching schema contents
- Integration with Bleve v2 (github.com/blevesearch/bleve/v2) for full-text indexing
- Automatic background indexing when schemas are added or updated
- Manual `--reindex` flag for rebuilding indexes
- Interactive result selector when multiple matches are found
- Index storage in `~/.config/gqlxp/schemas/<schema-id>.bleve/`
- Progress indicators for indexing operations
- Automatic detection and indexing of existing schemas without indexes

## Impact
- **Affected specs**: schema-library, schema-search (new), code-quality
- **Affected code**:
  - `library/library.go` - add indexing triggers
  - `cli/app.go` - add search command
  - New `search/` package
  - `tests/fitness/hierarchy_test.go` - update package hierarchy
  - `tests/fitness/dependency_test.go` - add bleve dependency restriction
- **New dependency**: `github.com/blevesearch/bleve/v2` (restricted to `search` package only)
- **Breaking changes**: None - purely additive feature
