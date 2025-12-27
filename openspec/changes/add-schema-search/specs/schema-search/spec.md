# schema-search Specification

## Purpose
Enable efficient full-text search across GraphQL schema types, fields, and descriptions using Bleve indexing. Automatically manages indexes for all library schemas with background indexing and manual rebuild capabilities.

## ADDED Requirements

### Requirement: Index Creation and Management
The system SHALL create and maintain Bleve indexes for each schema in the library, storing them in `~/.config/gqlxp/schemas/<schema-id>.bleve/`.

#### Scenario: Create index for new schema
- **WHEN** a schema is added to the library
- **THEN** a Bleve index is created in the background containing all searchable content

#### Scenario: Index already exists
- **WHEN** a schema with an existing valid index is accessed
- **THEN** the existing index is reused without rebuilding

#### Scenario: Index directory location
- **WHEN** an index is created for schema with ID "github-api"
- **THEN** the index is stored at `~/.config/gqlxp/schemas/github-api.bleve/`

### Requirement: Document Structure and Indexing
The system SHALL index GraphQL types and fields with searchable fields for type category, name, description, and path.

#### Scenario: Index type definition
- **WHEN** indexing a GraphQL Object type named "User" with description "Represents a user account"
- **THEN** a document is created with type="Object", name="User", description="Represents a user account", path="User"

#### Scenario: Index field definition
- **WHEN** indexing a field "email" on type "User" with description "User's email address"
- **THEN** a document is created with type="Field", name="email", description="User's email address", path="User.email"

#### Scenario: Index nested field
- **WHEN** indexing a field "street" on nested type "Address" accessed via "User.address"
- **THEN** a document is created with path="User.address.street"

### Requirement: Background Indexing
The system SHALL index schemas in the background using goroutines to avoid blocking CLI operations, with progress indication.

#### Scenario: Non-blocking indexing
- **WHEN** a schema is added to the library
- **THEN** indexing starts in a background goroutine and the CLI returns control to the user

#### Scenario: Progress indication
- **WHEN** indexing is in progress
- **THEN** a progress indicator shows indexing status with schema name

#### Scenario: Indexing completion
- **WHEN** background indexing completes successfully
- **THEN** the index is marked as ready for search operations

### Requirement: Search Query Execution
The system SHALL execute full-text search queries across indexed schema content and return ranked results.

#### Scenario: Search by type name
- **WHEN** searching for "User"
- **THEN** all types and fields with "User" in their name or description are returned, ranked by relevance

#### Scenario: Search by description
- **WHEN** searching for "email address"
- **THEN** all types and fields with "email address" in their description are returned

#### Scenario: Multiple matches
- **WHEN** a search returns multiple results
- **THEN** results are ranked by Bleve's relevance scoring with highest scores first

#### Scenario: No matches
- **WHEN** a search query matches no documents
- **THEN** an empty result set is returned with appropriate user message

### Requirement: Search Blocking Until Index Ready
The system SHALL block search operations until the index is built and ready, providing clear feedback to users.

#### Scenario: Search with index ready
- **WHEN** a search is executed and the index is ready
- **THEN** results are returned immediately

#### Scenario: Search while indexing
- **WHEN** a search is executed while indexing is in progress
- **THEN** the search blocks until indexing completes, showing progress indicator

#### Scenario: Search timeout
- **WHEN** a search is blocked for more than 60 seconds waiting for index
- **THEN** an error is returned suggesting manual reindex or checking for indexing failures

### Requirement: Manual Index Rebuild
The system SHALL provide a manual reindex command for recovery from index corruption or errors.

#### Scenario: Reindex via flag
- **WHEN** the `--reindex` flag is provided to the search command
- **THEN** the index is deleted and rebuilt from the schema file

#### Scenario: Reindex specific schema
- **WHEN** reindexing is requested for a specific schema ID
- **THEN** only that schema's index is rebuilt, leaving other indexes unchanged

#### Scenario: Reindex all schemas
- **WHEN** reindexing is requested without a specific schema ID
- **THEN** all schemas in the library have their indexes rebuilt

### Requirement: Index Detection for Existing Schemas
The system SHALL detect missing indexes for existing schemas and create them automatically on first search.

#### Scenario: Missing index on search
- **WHEN** a search is executed for a schema without an index
- **THEN** the index is created automatically in the background before executing the search

#### Scenario: Corrupted index detection
- **WHEN** an index directory exists but cannot be opened by Bleve
- **THEN** the corrupted index is deleted and a new index is created

### Requirement: Search Result Presentation
The system SHALL present search results in an interactive format, allowing users to select and navigate to matches.

#### Scenario: Single result
- **WHEN** a search returns exactly one result
- **THEN** the TUI is opened directly to that type or field

#### Scenario: Multiple results selection
- **WHEN** a search returns multiple results
- **THEN** an interactive list selector is displayed showing all matches with their paths

#### Scenario: Result selection navigation
- **WHEN** a user selects a result from the list
- **THEN** the TUI opens to the selected type or field location

#### Scenario: Result display format
- **WHEN** results are displayed in the selector
- **THEN** each result shows the path (e.g., "Query.user.email") and a snippet of the description

### Requirement: Index Cleanup on Schema Removal
The system SHALL automatically remove indexes when their associated schemas are deleted from the library.

#### Scenario: Delete schema with index
- **WHEN** a schema is removed from the library
- **THEN** the associated `.bleve` index directory is also deleted

#### Scenario: Index deletion failure
- **WHEN** index deletion fails due to permissions or file locks
- **THEN** a warning is displayed but schema removal completes successfully

### Requirement: Error Handling and Recovery
The system SHALL handle indexing and search errors gracefully with clear error messages and recovery suggestions.

#### Scenario: Indexing failure
- **WHEN** indexing fails due to schema parsing errors
- **THEN** an error message is displayed and the user is prompted to check the schema or use `--reindex`

#### Scenario: Search on failed index
- **WHEN** a search is attempted but the index failed to build
- **THEN** an error is displayed with suggestion to run with `--reindex`

#### Scenario: Disk space error
- **WHEN** indexing fails due to insufficient disk space
- **THEN** a clear error message indicates disk space issue and suggests cleanup

#### Scenario: Permission error
- **WHEN** index creation fails due to file permissions
- **THEN** an error message indicates permission issue and shows the attempted path
