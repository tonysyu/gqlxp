# Capability: Schema Library

## ADDED Requirements

### Requirement: Config Directory Initialization
The system SHALL create a configuration directory structure on first use in the user's standard config location (`~/.config/gqlxp/` on macOS/Linux, `%APPDATA%\gqlxp\` on Windows).

#### Scenario: First-time initialization
- **WHEN** the library feature is used for the first time
- **THEN** the config directory structure is created with `schemas/` and `metadata/` subdirectories

#### Scenario: Directory already exists
- **WHEN** the config directory already exists
- **THEN** the existing directory is used without modification

#### Scenario: Permission denied
- **WHEN** the config directory cannot be created due to permissions
- **THEN** an error message is displayed and the application exits gracefully

### Requirement: Schema Storage
The system SHALL store GraphQL schema files in the config directory with a unique identifier.

#### Scenario: Store new schema
- **WHEN** a user adds a schema with a unique ID
- **THEN** the schema file is saved as `schemas/<schema-id>.graphqls`

#### Scenario: Duplicate schema ID
- **WHEN** a user attempts to add a schema with an existing ID
- **THEN** an error is returned indicating the ID conflict

#### Scenario: Retrieve schema by ID
- **WHEN** a user requests a schema by its ID
- **THEN** the schema content is loaded from the stored file

### Requirement: Schema Metadata Persistence
The system SHALL store schema metadata as JSON files with support for display names, favorites, and URL patterns.

#### Scenario: Create metadata for new schema
- **WHEN** a schema is added to the library
- **THEN** a metadata file is created at `metadata/<schema-id>.json` with default values

#### Scenario: Update metadata
- **WHEN** metadata is modified (display name, favorites, URL patterns)
- **THEN** the metadata file is atomically updated with the new values

#### Scenario: Load metadata
- **WHEN** a schema is retrieved from the library
- **THEN** the associated metadata is loaded and returned with the schema

#### Scenario: Invalid metadata JSON
- **WHEN** a metadata file contains invalid JSON
- **THEN** an error is logged and default metadata is used

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
The system SHALL support removing schemas and their associated metadata from the library.

#### Scenario: Remove existing schema
- **WHEN** a user removes a schema by ID
- **THEN** both the schema file and metadata file are deleted

#### Scenario: Remove non-existent schema
- **WHEN** a user attempts to remove a schema that doesn't exist
- **THEN** an error is returned indicating the schema was not found

### Requirement: Backward Compatibility
The system SHALL maintain existing file-path mode behavior when library features are not used.

#### Scenario: Direct file path mode
- **WHEN** the application is invoked with a file path argument
- **THEN** the schema is loaded directly from the file without accessing the library

#### Scenario: Library mode flag
- **WHEN** the application is invoked with a library mode flag
- **THEN** the schema is loaded from the library instead of a file path

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
