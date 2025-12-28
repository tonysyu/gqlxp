# schema-library Specification

## Purpose
TBD - created by archiving change add-schema-library. Update Purpose after archive.
## Requirements
### Requirement: Config Directory Initialization
The system SHALL create a configuration directory structure on first use in the user's standard config location (`~/.config/gqlxp/` on macOS/Linux, `%APPDATA%\gqlxp\` on Windows).

#### Scenario: First-time initialization
- **WHEN** the library feature is used for the first time
- **THEN** the config directory is created with a `schemas/` subdirectory and empty `schemas/metadata.json` file

#### Scenario: Directory already exists
- **WHEN** the config directory already exists
- **THEN** the existing directory is used without modification

#### Scenario: Permission denied
- **WHEN** the config directory cannot be created due to permissions
- **THEN** an error message is displayed and the application exits gracefully

### Requirement: Schema Storage
The system SHALL store GraphQL schema files in the config directory with a unique identifier and trigger background indexing for search.

#### Scenario: Store new schema
- **WHEN** a user adds a schema with a unique ID
- **THEN** the schema file is saved as `schemas/<schema-id>.graphqls`

#### Scenario: Duplicate schema ID
- **WHEN** a user attempts to add a schema with an existing ID
- **THEN** an error is returned indicating the ID conflict

#### Scenario: Retrieve schema by ID
- **WHEN** a user requests a schema by its ID
- **THEN** the schema content is loaded from the stored file

#### Scenario: Trigger indexing on add
- **WHEN** a schema is successfully added to the library
- **THEN** background indexing is triggered to create a search index for the schema

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

### Requirement: Favorite Types Management
The system SHALL support marking type names as favorites for quick access within each schema.

#### Scenario: Add favorite type
- **WHEN** a user marks a type as favorite
- **THEN** the type name is added to the schema's favorites list in metadata

#### Scenario: Remove favorite type
- **WHEN** a user unmarks a favorite type
- **THEN** the type name is removed from the schema's favorites list

#### Scenario: List favorites
- **WHEN** a user requests the favorites for a schema
- **THEN** all favorited type names are returned

### Requirement: URL Pattern Configuration
The system SHALL support configurable URL patterns for opening web documentation of types and fields.

#### Scenario: Set type-specific URL pattern
- **WHEN** a user configures a URL pattern for a specific type (e.g., "Query", "Mutation")
- **THEN** the pattern is stored with template variables (`${type}`, `${field}`)

#### Scenario: Set wildcard URL pattern
- **WHEN** a user configures a wildcard URL pattern (`*`)
- **THEN** the pattern is used as fallback for types without specific patterns

#### Scenario: Resolve URL for type
- **WHEN** a URL is requested for a specific type and field
- **THEN** the appropriate pattern is selected and template variables are substituted

#### Scenario: No URL pattern configured
- **WHEN** a URL is requested but no pattern is configured
- **THEN** no URL is generated and an appropriate message is returned

### Requirement: Schema Listing
The system SHALL provide the ability to list all schemas in the library with their metadata.

#### Scenario: List all schemas
- **WHEN** a user requests the schema library list
- **THEN** all schema IDs and display names are returned

#### Scenario: Empty library
- **WHEN** the library contains no schemas
- **THEN** an empty list is returned

### Requirement: Schema Removal
The system SHALL support removing schemas, their associated metadata, and their search indexes from the library.

#### Scenario: Remove existing schema
- **WHEN** a user removes a schema by ID
- **THEN** the schema file is deleted and its entry is removed from `schemas/metadata.json`

#### Scenario: Remove non-existent schema
- **WHEN** a user attempts to remove a schema that doesn't exist
- **THEN** an error is returned indicating the schema was not found

#### Scenario: Remove search index
- **WHEN** a schema is removed from the library
- **THEN** the associated search index directory is also deleted if it exists

#### Scenario: Index deletion failure
- **WHEN** the search index cannot be deleted due to permissions or file locks
- **THEN** a warning is logged but schema removal completes successfully

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

### Requirement: Schema ID Validation
The system SHALL enforce schema ID format rules to ensure filesystem compatibility.

#### Scenario: Valid schema ID
- **WHEN** a schema ID contains only lowercase letters, numbers, and hyphens
- **THEN** the ID is accepted

#### Scenario: Invalid schema ID characters
- **WHEN** a schema ID contains uppercase letters, spaces, or special characters
- **THEN** an error is returned with guidance on valid ID format

#### Scenario: ID sanitization suggestion
- **WHEN** an invalid ID is provided
- **THEN** a sanitized version is suggested to the user

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
The system SHALL detect when schema files have changed, prompt for user action, and trigger re-indexing when content is updated.

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

#### Scenario: Trigger re-indexing on update
- **WHEN** a schema's content is updated in the library
- **THEN** background re-indexing is triggered to update the search index with new content

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

