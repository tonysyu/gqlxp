# schema-library Specification Deltas

## MODIFIED Requirements

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
