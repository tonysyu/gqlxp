# Tasks for add-library-command

## Implementation Tasks

### 1. Create library CLI command structure
**Status:** pending
**Validation:** Run `gqlxp library --help` to verify command structure

- Create `cli/library.go` with main library command
- Add library command to CLI app commands in `cli/main.go`
- Implement help text showing all subcommands
- Set default action to list schemas when no subcommand provided

### 2. Implement library list subcommand
**Status:** pending
**Validation:** Run `gqlxp library list` with populated and empty library

- Create `listCommand()` function in `cli/library.go`
- Call `library.List()` to retrieve schemas
- Format output showing ID, display name, and source file path
- Indicate default schema with marker (e.g., "*")
- Display helpful message when library is empty

### 3. Implement library add subcommand
**Status:** pending
**Validation:** Run `gqlxp library add` with various flag combinations

- Create `addCommand()` function in `cli/library.go`
- Add `--id` and `--name` flags for non-interactive usage
- Implement interactive prompts for ID and name when flags not provided
- Validate schema ID format and show suggestions for invalid IDs
- Call `library.Add()` with validated inputs
- Display success message with schema details

### 4. Implement library remove subcommand
**Status:** pending
**Validation:** Run `gqlxp library remove` with confirmation and force flag

- Create `removeCommand()` function in `cli/library.go`
- Add `--force` flag to skip confirmation
- Implement confirmation prompt showing schema ID and display name
- Call `library.Remove()` to delete schema, metadata, and index
- Clear default schema setting if removing the default
- Display appropriate success or error messages

### 5. Implement library default subcommand
**Status:** pending
**Validation:** Run `gqlxp library default` to show/set/clear default schema

- Create `defaultCommand()` function in `cli/library.go`
- Migrate code from `cli/config.go` default-schema subcommand
- Add `--clear` flag to remove default schema setting
- Display current default when no argument provided
- Validate schema exists before setting as default
- Display confirmation messages

### 6. Implement library reindex subcommand
**Status:** pending
**Validation:** Run `gqlxp library reindex` for single and all schemas

- Create `reindexCommand()` function in `cli/library.go`
- Add `--all` flag to reindex all library schemas
- Parse schema content and create indexer
- Delete existing index and rebuild from schema
- Display progress messages during indexing
- Handle errors gracefully with clear messages

### 7. Remove --reindex flag from search command
**Status:** pending
**Validation:** Run `gqlxp search --reindex` to verify flag is gone

- Remove `--reindex` flag definition from `searchCommand()` in `cli/search.go`
- Remove reindex logic from search command action
- Keep auto-index creation when index doesn't exist
- Update command description to remove reindex references

### 8. Remove config command
**Status:** pending
**Validation:** Run `gqlxp config` to verify helpful error message

- Remove `configCommand()` from `cli/main.go` commands list
- Delete `cli/config.go` file
- Optionally add migration hint if user tries old command (can detect in error handling)

### 9. Update tests
**Status:** pending
**Validation:** Run `just test` to verify all tests pass

- Remove tests for config command
- Remove tests for search --reindex flag
- Add tests for library list command
- Add tests for library add command
- Add tests for library remove command
- Add tests for library default command
- Add tests for library reindex command

### 10. Update documentation
**Status:** pending
**Validation:** Review README and docs for accuracy

- Update `README.md` CLI examples to use library command
- Update any references from `config` to `library` command
- Add library command examples to relevant documentation
- Document migration from old commands to new ones

## Verification

After all tasks:
- Run `just test` to ensure all tests pass
- Run `just build` to ensure clean build
- Manually test each library subcommand
- Test migration scenarios from old commands
