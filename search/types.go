package search

import "github.com/tonysyu/gqlxp/gql"

// SearchResult represents a single search result with ranking information
type SearchResult struct {
	Type        string  // Type of result (Object, Field, Enum, etc.)
	Name        string  // Name of the type or field
	Path        string  // Full path (e.g., "Query.user.name")
	Description string  // Description text
	Score       float64 // Relevance score from Bleve
}

// Indexer manages schema indexing operations
type Indexer interface {
	// Index creates or updates the index for a schema
	Index(schemaID string, schema *gql.GraphQLSchema) error

	// Remove deletes the index for a schema
	Remove(schemaID string) error

	// Exists checks if an index exists for a schema
	Exists(schemaID string) bool

	// Close closes the indexer and releases resources
	Close() error
}

// Searcher performs search operations on indexed schemas
type Searcher interface {
	// Search finds matching types and fields in a schema
	Search(schemaID string, query string, limit int) ([]SearchResult, error)

	// Close closes the searcher and releases resources
	Close() error
}
