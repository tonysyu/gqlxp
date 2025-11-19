# Tasks: Automatic Library Integration

## Phase 1: Library Enhancement

### 1. Add file hash support to library package
- [x] Add `FileHash string` field to `SchemaMetadata` struct in `library/types.go`
- [x] Implement `calculateFileHash(content []byte) string` function using SHA-256
- [x] Update `Library.Add()` to compute and store file hash
- [x] Update `Library.Add()` to convert source path to absolute path before storage
- [x] Ensure existing tests pass with new metadata field (default to empty string for backward compatibility)

### 2. Implement library lookup by file path
- [x] Add `FindByPath(absolutePath string) (*Schema, error)` method to `Library` interface
- [x] Implement path-based search in `FileLibrary.FindByPath()`
- [x] Add tests for path lookup with matching and non-matching cases
- [x] Handle case-sensitivity appropriately for the platform

### 3. Add library content update capability
- [x] Add `UpdateContent(id string, content []byte) error` method to `Library` interface
- [x] Implement in `FileLibrary` to update both schema file and hash in metadata
- [x] Preserve existing metadata (favorites, URL patterns, display name, created timestamp)
- [x] Update `UpdatedAt` timestamp on content changes
- [x] Add tests for update preserving metadata

## Phase 2: Interactive Prompting

### 4. Create prompting utilities
- [x] Create `internal/prompt` package for user input
- [x] Implement `PromptString(prompt string, defaultValue string) (string, error)` for basic input
- [x] Implement `PromptYesNo(prompt string) (bool, error)` for confirmation
- [x] Implement `PromptSchemaID(suggestedID string) (string, error)` with validation loop
- [x] Handle EOF and interrupt signals gracefully
- [x] Add tests for prompt functions (may require input mocking)

### 5. Build schema registration flow
- [x] Create `registerSchema(filePath string, content []byte) (string, error)` function in main
- [x] Use `PromptSchemaID()` with filename-based suggestion
- [x] Use `PromptString()` for display name with ID as default
- [x] Call `Library.Add()` with user-provided values
- [x] Handle validation errors with retry logic
- [x] Return schema ID for subsequent loading

## Phase 3: CLI Integration

### 6. Refactor main.go startup logic
- [x] Extract file loading logic into `loadSchemaFromFile(path string) ([]byte, error)`
- [x] Implement `resolveSchemaSource(filePath string) (schemaID string, content []byte, err error)` that:
  - [x] Normalizes path to absolute
  - [x] Calculates file hash
  - [x] Checks library for path match
  - [x] Handles no-match, full-match, and hash-mismatch cases
- [x] Replace existing file/library branching with new flow

### 7. Implement automatic library integration
- [x] Modify file path handler to call `resolveSchemaSource()`
- [x] Implement hash mismatch prompt: "Schema file changed. Update library? (y/n)"
- [x] On update: call `Library.UpdateContent()`, then proceed
- [x] On skip: load existing library schema by ID
- [x] On no match: call `registerSchema()`, then proceed
- [x] Ensure all paths lead to library-backed TUI session

### 8. Update no-args behavior
- [x] Check if library is empty when no arguments provided
- [x] If empty: display error "No schema file provided. Usage: gqlxp <schema-file>"
- [x] If not empty: call `tui.StartSchemaSelector()`
- [x] Remove requirement for `--library` flag to trigger selector

## Phase 4: Cleanup and Documentation

### 9. Update CLI help and usage
- [x] Simplify usage message to focus on file-first workflow
- [x] Document automatic library behavior
- [x] Update examples to reflect new interaction model
- [x] Removed `library` subcommands (functionality moved to TUI)

### 10. Add integration tests
- [ ] Test full flow: new file → prompts → library save → TUI start
- [ ] Test hash match: existing file unchanged → direct load
- [ ] Test hash mismatch: existing file changed → prompt → update/skip
- [ ] Test no args with empty library → error message
- [ ] Test no args with non-empty library → selector

### 11. Run validation and verification
- [x] Run `just test` to ensure all unit tests pass
- [x] Run `just verify` for linting and formatting
- [ ] Manually test interactive flows (cannot be fully automated)
- [ ] Test on sample schema files

### 12. Update documentation
- [x] Update README.md with new workflow description
- [x] Update docs/schema-library.md with automatic integration details
- [x] Add note about file path/hash tracking
- [x] Document that all schemas are now library-backed
