# Design: Schema Search

## Context
GraphQL schemas can contain hundreds or thousands of types, fields, and enums. Users need efficient search to find specific types or fields without manually navigating the TUI. The search feature must handle large schemas, provide fast results, and integrate seamlessly with the existing library system.

**Constraints:**
- No external services (terminal-only application)
- Must work offline with local schema files
- Indexing should not block CLI operations
- Indexes must be rebuildable in case of corruption

**Stakeholders:**
- Users exploring large schemas
- Library system (triggers indexing on schema changes)

## Goals / Non-Goals

**Goals:**
- Sub-second search response time for typical queries
- Search across type names, field names, and descriptions
- Non-blocking user experience during indexing
- Automatic index management (create, update, rebuild)
- Reliable error recovery (manual reindex option)

**Non-Goals:**
- Fuzzy search or typo correction (future enhancement)
- Search across multiple schemas simultaneously
- Real-time indexing (batch indexing is sufficient)
- Index versioning or migration strategies

## Decisions

### Decision: Use Bleve v2 for indexing
**Rationale:**
- Pure Go implementation (no C dependencies, easy cross-platform builds)
- Actively maintained by Couchbase
- Rich feature set (text analysis, stemming, scoring)
- No external services required (embedded database)
- Proven in production use

**Alternatives considered:**
- `bluge`: Newer fork of Bleve but less mature
- Custom inverted index: Significant development time, reinventing the wheel
- Ripgrep integration: No ranking, poor UX for partial matches

### Decision: Index storage location
**Path:** `~/.config/gqlxp/schemas/<schema-id>.bleve/`

**Rationale:**
- Co-located with schema files for easy management
- One index per schema (simple deletion, no shared state)
- Standard config directory location
- Easy cleanup when schema is removed

**Alternatives considered:**
- Single shared index: Complex updates, harder to manage
- In-memory only: Lost on restart, slow startup for large schemas

### Decision: Background indexing with blocking search
**Strategy:**
- Indexing runs in goroutine with progress indicator
- Search operations block until index is ready
- Failed indexing prompts user to manually reindex

**Rationale:**
- Simple UX (users know index is building)
- Avoids partial/stale results
- Progress feedback prevents "hanging" perception
- Manual recovery for edge cases

**Alternatives considered:**
- Non-blocking search with partial results: Confusing UX, inconsistent results
- Synchronous indexing: Blocks CLI, poor UX for large schemas
- Background with retry: Complex, could loop infinitely on persistent errors

### Decision: Document structure
**Fields indexed:**
- `type` (Object, Enum, etc.)
- `name` (type/field name)
- `description` (type/field description)
- `path` (e.g., "Query.user.name")
- `schemaID` (for filtering)

**Rationale:**
- Covers user's search requirements (names + descriptions)
- Path enables precise result presentation
- Type filtering for future enhancements

### Decision: Package hierarchy placement
**Hierarchy:** `cmd → cli → tui → library → search → gql → utils`

**Rationale:**
- `search` needs `gql` to parse schemas and extract searchable data
- `library` needs `search` to trigger indexing on schema changes
- Maintains clean dependency flow (no cycles)

**Alternatives considered:**
- `search` at same level as `library`: Would require both to import shared indexing interface
- `search` below `gql`: Library couldn't trigger indexing without importing `gql` directly

## Risks / Trade-offs

### Risk: Index corruption
**Mitigation:** Manual `--reindex` flag rebuilds from source schema

### Risk: Disk space usage
**Impact:** ~1-5% of schema size for indexes
**Mitigation:** Indexes are regenerable, can be deleted and rebuilt

### Risk: Indexing performance for large schemas
**Impact:** May take seconds for very large schemas (>10MB)
**Mitigation:**
- Background operation with progress indicator
- Users can continue using TUI while indexing
- One-time cost per schema update

### Trade-off: Blocking search vs. partial results
**Decision:** Block search until index ready
**Rationale:** Simpler UX, predictable behavior, avoids confusion about incomplete results

## Migration Plan

**Automatic migration:**
1. On first search, detect missing indexes for existing schemas
2. Prompt user or auto-index (TBD in implementation)
3. Show progress for each schema being indexed

**Rollback:**
- Indexes are optional; deleting `.bleve/` directories reverts to non-search state
- No schema data is modified

**Performance:**
- Initial indexing happens once per schema
- Subsequent searches are fast (sub-second)

## Open Questions

1. Should missing indexes be automatically created on first search, or require explicit user action?
   - **Recommendation:** Auto-create with progress indicator for better UX

2. Should search results be ranked by relevance or alphabetically?
   - **Recommendation:** Use Bleve's relevance scoring for better UX

3. Should there be a maximum index size limit or warning?
   - **Recommendation:** No limit initially; add warning if user reports issues
