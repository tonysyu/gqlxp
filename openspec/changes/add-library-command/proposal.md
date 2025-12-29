# add-library-command

## Summary
Add a new `library` CLI command that consolidates schema library management functionality, replacing the `config` command and `search --reindex` flag with a unified, intuitive interface.

## Motivation
The current CLI has schema library management functionality scattered across multiple commands:
- `config default-schema` for setting the default schema
- `search --reindex` for rebuilding search indexes
- No direct CLI support for listing, adding, or removing schemas from the library

This fragmentation makes the CLI harder to learn and use. Users must remember which command handles which library operation, and there's no obvious entry point for discovering library management capabilities.

## Proposal
Create a new `library` command that serves as the central interface for all library management operations:
- `library list` - List all schemas in the library
- `library add` - Add a schema to the library
- `library remove` - Remove a schema from the library
- `library default` - Set or show the default schema (replaces `config default-schema`)
- `library reindex` - Rebuild search indexes (replaces `search --reindex`)

The `config` command will be deprecated and removed, and the `--reindex` flag will be removed from the `search` command.

## Impact
### Breaking Changes
- The `config` command will be removed
- The `search --reindex` flag will be removed

### Migration Path
Users currently using:
- `gqlxp config default-schema [ID]` → `gqlxp library default [ID]`
- `gqlxp search --reindex SCHEMA QUERY` → `gqlxp library reindex SCHEMA` followed by `gqlxp search SCHEMA QUERY`

### Benefits
- Single, discoverable command for all library operations
- More intuitive command structure aligned with user mental models
- Clearer separation between library management and schema exploration/search
- Easier to extend with future library features

## Dependencies
None - this builds on existing library and search infrastructure.

## Alternatives Considered
1. **Keep current structure** - Rejected because it's confusing and doesn't scale well as library features grow
2. **Add library subcommands without removing config** - Rejected to avoid duplication and further confusion
3. **Use `schema` instead of `library`** - Rejected because "library" better represents the collection/management aspect

## Open Questions
None

## Related Specs
- `cli-interface` - Will be modified to add library command requirements
- `schema-library` - No changes needed (library package already has required functionality)
- `schema-search` - Will be modified to update reindex requirement
