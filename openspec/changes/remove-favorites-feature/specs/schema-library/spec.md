# schema-library Specification Delta

## REMOVED Requirements

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

## MODIFIED Requirements

### Requirement: Schema Metadata Persistence
The system SHALL store schema metadata in a single JSON file (`schemas/metadata.json`) with schema-id as top-level keys, supporting display names, ~~favorites,~~ URL patterns, **file paths, and file hashes**.

#### Scenario: Update metadata
- **WHEN** metadata is modified (display name, ~~favorites,~~ URL patterns)
- **THEN** the `schemas/metadata.json` file is atomically updated with the new values **and UpdatedAt timestamp is set to current time**

### Requirement: Schema Update Detection
The system SHALL detect when schema files have changed, prompt for user action, and trigger re-indexing when content is updated.

#### Scenario: Preserve metadata on update
- **WHEN** a schema is updated due to file changes
- **THEN** existing metadata (~~favorites,~~ URL patterns, display name) is preserved
