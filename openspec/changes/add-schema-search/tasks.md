# Implementation Tasks

## 1. Package Structure and Dependencies
- [x] 1.1 Create `search/` package directory
- [x] 1.2 Add `github.com/blevesearch/bleve/v2` dependency to `go.mod`
- [x] 1.3 Create `search/types.go` with core interfaces (Indexer, SearchResult)
- [x] 1.4 Update `tests/fitness/hierarchy_test.go` to include `search` in PACKAGE_HIERARCHY
- [x] 1.5 Update `tests/fitness/dependency_test.go` to restrict bleve imports to `search` package only
- [x] 1.6 Run `just test tests/fitness` to verify fitness tests pass

## 2. Core Search Implementation
- [x] 2.1 Create `search/indexer.go` with BleveIndexer implementation
- [x] 2.2 Implement index creation with document mapping (type, name, description, path fields)
- [x] 2.3 Implement schema parsing to extract indexable documents from gql.GraphQLSchema
- [x] 2.4 Create `search/searcher.go` with Search() function
- [x] 2.5 Implement result ranking and formatting
- [x] 2.6 Add unit tests for indexer in `search/indexer_test.go`
- [x] 2.7 Add unit tests for searcher in `search/searcher_test.go`

## 3. Background Indexing
- [x] 3.1 Create `search/background.go` with goroutine-based indexing
- [x] 3.2 Implement progress reporting channel/callback mechanism
- [x] 3.3 Add index readiness detection (check if index exists and is valid)
- [x] 3.4 Implement blocking search with timeout while index builds
- [x] 3.5 Add error handling for indexing failures
- [x] 3.6 Add unit tests for background indexing behavior

## 4. Library Integration
- [x] 4.1 Add `indexSchema()` helper function in `library/library.go`
- [x] 4.2 Modify `library.Add()` to trigger background indexing after schema storage
- [x] 4.3 Modify `library.UpdateContent()` to trigger re-indexing after content update
- [x] 4.4 Modify `library.Remove()` to delete index directory when schema removed
- [x] 4.5 Add `library.GetIndexPath(schemaID)` helper for index location
- [x] 4.6 Implement detection of missing indexes for existing schemas
- [x] 4.7 Add integration tests for library + search interaction

## 5. CLI Search Command
- [x] 5.1 Add `searchCommand()` to `cli/app.go`
- [x] 5.2 Implement query parsing and validation
- [x] 5.3 Implement progress indicator display during indexing
- [x] 5.4 Implement interactive result selector using bubbles list (for multiple matches)
- [x] 5.5 Implement direct TUI launch for single match
- [x] 5.6 Add `--reindex` flag for manual index rebuilding
- [x] 5.7 Add error messaging for indexing failures with reindex suggestion
- [ ] 5.8 Add acceptance test for search command workflow

## 6. Documentation
- [x] 6.1 Update `README.md` with search command examples
- [x] 6.2 Update `docs/architecture.md` to include `search` package description
- [x] 6.3 Create `docs/search.md` with detailed search usage and index management
- [x] 6.4 Add search command to `docs/index.md`
- [x] 6.5 Update `CLAUDE.md` if any new AI-specific guidance needed

## 7. Validation and Polish
- [x] 7.1 Run `just test` to ensure all tests pass
- [ ] 7.2 Run `just verify` to check linting and formatting
- [x] 7.3 Test search with small schema (examples/github.graphqls)
- [ ] 7.4 Test search with large schema (if available)
- [ ] 7.5 Test `--reindex` flag manually
- [ ] 7.6 Verify index cleanup on schema removal
- [ ] 7.7 Check disk space usage of indexes
