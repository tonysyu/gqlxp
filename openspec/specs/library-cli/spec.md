# library-cli Specification

## Purpose
TBD - created by archiving change add-library-command. Update Purpose after archive.
## Requirements
### Requirement: Library Command Structure
The system SHALL provide a `library` command with subcommands for all library management operations.

#### Scenario: Library command help
- **WHEN** user runs `gqlxp library --help`
- **THEN** help text is displayed showing available subcommands: list, add, remove, default, reindex

#### Scenario: Library command without subcommand
- **WHEN** user runs `gqlxp library` with no subcommand
- **THEN** the default behavior is to list all schemas (equivalent to `gqlxp library list`)

### Requirement: List Schemas
The system SHALL provide a `library list` command that displays all schemas in the library.

#### Scenario: List schemas with entries
- **WHEN** user runs `gqlxp library list` and the library contains schemas
- **THEN** each schema is displayed with its ID, display name, and source file path

#### Scenario: List empty library
- **WHEN** user runs `gqlxp library list` and the library is empty
- **THEN** a message "No schemas in library. Add one with: gqlxp library add <schema-file>" is displayed

#### Scenario: List with default schema indicator
- **WHEN** user runs `gqlxp library list` and a default schema is set
- **THEN** the default schema is indicated with a marker (e.g., "* github (GitHub API)")

### Requirement: Add Schema to Library
The system SHALL provide a `library add` command that adds a schema file to the library with interactive prompts.

#### Scenario: Add schema with ID and name
- **WHEN** user runs `gqlxp library add --id github --name "GitHub API" schema.graphqls`
- **THEN** the schema is added to the library with the specified ID and display name

#### Scenario: Add schema with interactive prompts
- **WHEN** user runs `gqlxp library add schema.graphqls` without --id flag
- **THEN** the user is prompted for schema ID and display name

#### Scenario: Add duplicate schema ID
- **WHEN** user attempts to add a schema with an ID that already exists
- **THEN** an error "Schema 'github' already exists in library" is displayed and the operation fails

#### Scenario: Add schema with invalid ID
- **WHEN** user provides an invalid schema ID (e.g., with uppercase or spaces)
- **THEN** an error is displayed with a suggested sanitized ID

### Requirement: Remove Schema from Library
The system SHALL provide a `library remove` command that removes a schema and its index from the library.

#### Scenario: Remove existing schema
- **WHEN** user runs `gqlxp library remove github`
- **THEN** the schema file, metadata entry, and search index are deleted

#### Scenario: Remove with confirmation prompt
- **WHEN** user runs `gqlxp library remove github` without --force flag
- **THEN** a confirmation prompt "Remove schema 'github' (GitHub API)? [y/N]" is displayed

#### Scenario: Remove with force flag
- **WHEN** user runs `gqlxp library remove --force github`
- **THEN** the schema is removed without confirmation prompt

#### Scenario: Remove non-existent schema
- **WHEN** user runs `gqlxp library remove nonexistent`
- **THEN** an error "Schema 'nonexistent' not found in library" is displayed

#### Scenario: Remove default schema
- **WHEN** user removes the schema that is set as default
- **THEN** the default schema setting is cleared and a message indicates this

### Requirement: Default Schema Configuration
The system SHALL provide a `library default` command that sets or displays the default schema.

#### Scenario: Set default schema
- **WHEN** user runs `gqlxp library default github`
- **THEN** the default schema is set to 'github' and a confirmation message is displayed

#### Scenario: Show current default
- **WHEN** user runs `gqlxp library default` with no argument
- **THEN** the current default schema ID and display name are displayed

#### Scenario: Clear default schema
- **WHEN** user runs `gqlxp library default --clear`
- **THEN** the default schema setting is removed

#### Scenario: Set non-existent schema as default
- **WHEN** user runs `gqlxp library default nonexistent`
- **THEN** an error "Schema 'nonexistent' not found in library" is displayed

### Requirement: Schema Reindexing
The system SHALL provide a `library reindex` command that rebuilds search indexes for schemas.

#### Scenario: Reindex specific schema
- **WHEN** user runs `gqlxp library reindex github`
- **THEN** the search index for schema 'github' is deleted and rebuilt

#### Scenario: Reindex all schemas
- **WHEN** user runs `gqlxp library reindex --all`
- **THEN** all schemas in the library have their search indexes rebuilt

#### Scenario: Reindex with progress indication
- **WHEN** reindexing is in progress
- **THEN** progress messages are displayed: "Reindexing 'github'...", "Index rebuilt successfully"

#### Scenario: Reindex non-existent schema
- **WHEN** user runs `gqlxp library reindex nonexistent`
- **THEN** an error "Schema 'nonexistent' not found in library" is displayed

