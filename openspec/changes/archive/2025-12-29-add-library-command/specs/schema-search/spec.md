# schema-search Deltas

## MODIFIED Requirements

### Requirement: Manual Index Rebuild
The system SHALL provide manual reindex functionality through the library command instead of a search flag.

#### Scenario: Reindex specific schema
- **WHEN** reindexing is requested for a specific schema ID via `gqlxp library reindex <schema-id>`
- **THEN** only that schema's index is rebuilt, leaving other indexes unchanged

#### Scenario: Reindex all schemas
- **WHEN** reindexing is requested via `gqlxp library reindex --all`
- **THEN** all schemas in the library have their indexes rebuilt

### Requirement: Error Handling and Recovery
The system SHALL handle indexing and search errors gracefully with clear error messages and recovery suggestions.

#### Scenario: Indexing failure
- **WHEN** indexing fails due to schema parsing errors
- **THEN** an error message is displayed and the user is prompted to check the schema or use `gqlxp library reindex`

#### Scenario: Search on failed index
- **WHEN** a search is attempted but the index failed to build
- **THEN** an error is displayed with suggestion to run `gqlxp library reindex <schema-id>`

#### Scenario: Disk space error
- **WHEN** indexing fails due to insufficient disk space
- **THEN** a clear error message indicates disk space issue and suggests cleanup

#### Scenario: Permission error
- **WHEN** index creation fails due to file permissions
- **THEN** an error message indicates permission issue and shows the attempted path
