# schema-library Specification Delta

## MODIFIED Requirements

### Requirement: Schema Metadata Persistence
The system SHALL store schema metadata in a single JSON file (`schemas/metadata.json`) with schema-id as top-level keys, supporting display names, favorites, URL patterns, **file paths, and file hashes**.

#### Scenario: Create metadata for new schema
- **WHEN** a schema is added to the library
- **THEN** an entry is added to `schemas/metadata.json` with the schema-id as key, **absolute source file path, SHA-256 hash of content**, and default metadata values

#### Scenario: Update metadata
- **WHEN** metadata is modified (display name, favorites, URL patterns)
- **THEN** the `schemas/metadata.json` file is atomically updated with the new values **and UpdatedAt timestamp is set to current time**

#### Scenario: Load metadata
- **WHEN** a schema is retrieved from the library
- **THEN** the associated metadata is loaded from `schemas/metadata.json` using the schema-id as key **including file path and hash**

#### Scenario: Query by file path
- **WHEN** searching for an existing schema by source file path
- **THEN** the library returns matching schema ID and metadata if an absolute path match exists

#### Scenario: Query by file hash
- **WHEN** searching for an existing schema by content hash
- **THEN** the library returns matching schema ID and metadata if a hash match exists

### Requirement: Backward Compatibility
The system SHALL **require all schema exploration to use the library**, while maintaining intuitive file-based workflows through automatic library integration.

#### Scenario: Direct file path with library match
- **WHEN** the application is invoked with a file path that matches an existing library entry (same path and hash)
- **THEN** the schema is loaded from the library without prompting

#### Scenario: Direct file path with hash mismatch
- **WHEN** the application is invoked with a file path that matches an existing library entry path but with different hash
- **THEN** the user is prompted to update the library or load the existing library version

#### Scenario: Direct file path without library match
- **WHEN** the application is invoked with a file path that is not in the library
- **THEN** the user is prompted for schema ID and display name, then the schema is saved to the library before starting the TUI

#### Scenario: No arguments with non-empty library
- **WHEN** the application is invoked without arguments and the library contains schemas
- **THEN** the schema selector interface is displayed

#### Scenario: No arguments with empty library
- **WHEN** the application is invoked without arguments and the library is empty
- **THEN** an error message instructs the user to provide a schema file path

## ADDED Requirements

### Requirement: File Hash Management
The system SHALL compute and store SHA-256 hashes of schema file contents for change detection.

#### Scenario: Calculate hash on add
- **WHEN** a schema is added to the library from a file
- **THEN** the SHA-256 hash of the file contents is calculated and stored in metadata

#### Scenario: Calculate hash on load
- **WHEN** a schema file is provided as input
- **THEN** the SHA-256 hash of the file contents is calculated for comparison with library entries

#### Scenario: Hash comparison
- **WHEN** comparing a file to a library entry
- **THEN** hashes are compared to detect if file contents have changed

### Requirement: Interactive Schema Registration
The system SHALL prompt users for required information when adding schemas to the library.

#### Scenario: Prompt for schema ID
- **WHEN** a schema file is not in the library
- **THEN** the user is prompted to enter a schema ID with format validation and suggestions

#### Scenario: Prompt for display name
- **WHEN** a schema file is being added to the library
- **THEN** the user is prompted to enter a display name with the schema ID as the default

#### Scenario: Schema ID validation on input
- **WHEN** a user enters an invalid schema ID during registration
- **THEN** an error is shown with a sanitized suggestion and the prompt is repeated

#### Scenario: Accept default display name
- **WHEN** a user presses enter without entering a display name
- **THEN** the schema ID is used as the display name

### Requirement: Schema Update Detection
The system SHALL detect when schema files have changed and prompt for user action.

#### Scenario: Prompt on hash mismatch
- **WHEN** a schema file path matches a library entry but the hash differs
- **THEN** the user is prompted to either update the library or load the existing library version

#### Scenario: Update library on user confirmation
- **WHEN** a user chooses to update the library for a changed file
- **THEN** the schema content and hash are updated in the library and the TUI is started with the updated schema

#### Scenario: Load existing on user decline
- **WHEN** a user chooses to load the existing library version for a changed file
- **THEN** the schema is loaded from the library without updating and the TUI is started

#### Scenario: Preserve metadata on update
- **WHEN** a schema is updated due to file changes
- **THEN** existing metadata (favorites, URL patterns, display name) is preserved

### Requirement: Library Path Lookup
The system SHALL support efficient lookup of schemas by source file path.

#### Scenario: Absolute path normalization
- **WHEN** a file path is provided (relative or absolute)
- **THEN** the path is converted to absolute form for comparison and storage

#### Scenario: Find schema by path
- **WHEN** querying the library for a schema by file path
- **THEN** all library entries are searched for a matching source file path

#### Scenario: No path match found
- **WHEN** no library entry matches the provided file path
- **THEN** the query returns no match indicating the schema is new

## REMOVED Requirements

### Requirement: Backward Compatibility - Library mode flag
**REMOVED**: The `--library` flag scenario is deprecated in favor of automatic library integration.

**Rationale**: All schema operations now use the library, making explicit library mode flags redundant. The no-arguments case automatically opens the selector when the library is non-empty.
