# Tasks: Automatic Library Integration

## Phase 1: Library Enhancement

### 1. Add file hash support to library package
- Add `FileHash string` field to `SchemaMetadata` struct in `library/types.go`
- Implement `calculateFileHash(content []byte) string` function using SHA-256
- Update `Library.Add()` to compute and store file hash
- Update `Library.Add()` to convert source path to absolute path before storage
- Ensure existing tests pass with new metadata field (default to empty string for backward compatibility)

### 2. Implement library lookup by file path
- Add `FindByPath(absolutePath string) (*Schema, error)` method to `Library` interface
- Implement path-based search in `FileLibrary.FindByPath()`
- Add tests for path lookup with matching and non-matching cases
- Handle case-sensitivity appropriately for the platform

### 3. Add library content update capability
- Add `UpdateContent(id string, content []byte) error` method to `Library` interface
- Implement in `FileLibrary` to update both schema file and hash in metadata
- Preserve existing metadata (favorites, URL patterns, display name, created timestamp)
- Update `UpdatedAt` timestamp on content changes
- Add tests for update preserving metadata

## Phase 2: Interactive Prompting

### 4. Create prompting utilities
- Create `internal/prompt` package for user input
- Implement `PromptString(prompt string, defaultValue string) (string, error)` for basic input
- Implement `PromptYesNo(prompt string) (bool, error)` for confirmation
- Implement `PromptSchemaID(suggestedID string) (string, error)` with validation loop
- Handle EOF and interrupt signals gracefully
- Add tests for prompt functions (may require input mocking)

### 5. Build schema registration flow
- Create `registerSchema(filePath string, content []byte) (string, error)` function in main
- Use `PromptSchemaID()` with filename-based suggestion
- Use `PromptString()` for display name with ID as default
- Call `Library.Add()` with user-provided values
- Handle validation errors with retry logic
- Return schema ID for subsequent loading

## Phase 3: CLI Integration

### 6. Refactor main.go startup logic
- Extract file loading logic into `loadSchemaFromFile(path string) ([]byte, error)`
- Implement `resolveSchemaSource(filePath string) (schemaID string, content []byte, err error)` that:
  - Normalizes path to absolute
  - Calculates file hash
  - Checks library for path match
  - Handles no-match, full-match, and hash-mismatch cases
- Replace existing file/library branching with new flow

### 7. Implement automatic library integration
- Modify file path handler to call `resolveSchemaSource()`
- Implement hash mismatch prompt: "Schema file changed. Update library? (y/n)"
- On update: call `Library.UpdateContent()`, then proceed
- On skip: load existing library schema by ID
- On no match: call `registerSchema()`, then proceed
- Ensure all paths lead to library-backed TUI session

### 8. Update no-args behavior
- Check if library is empty when no arguments provided
- If empty: display error "No schema file provided. Usage: gqlxp <schema-file>"
- If not empty: call `tui.StartSchemaSelector()`
- Remove requirement for `--library` flag to trigger selector

## Phase 4: Cleanup and Documentation

### 9. Update CLI help and usage
- Simplify usage message to focus on file-first workflow
- Document automatic library behavior
- Update examples to reflect new interaction model
- Consider deprecation notice for `library add` subcommand (keep functional)

### 10. Add integration tests
- Test full flow: new file → prompts → library save → TUI start
- Test hash match: existing file unchanged → direct load
- Test hash mismatch: existing file changed → prompt → update/skip
- Test no args with empty library → error message
- Test no args with non-empty library → selector

### 11. Run validation and verification
- Run `just test` to ensure all unit tests pass
- Run `just verify` for linting and formatting
- Manually test interactive flows (cannot be fully automated)
- Test on sample schema files

### 12. Update documentation
- Update README.md with new workflow description
- Update docs/schema-library.md with automatic integration details
- Add note about file path/hash tracking
- Document that all schemas are now library-backed
