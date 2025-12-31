package search

import (
	"fmt"
	"path/filepath"

	"github.com/blevesearch/bleve/v2"
)

// BleveSearcher implements Searcher using Bleve
type BleveSearcher struct {
	baseDir string // Base directory for storing indexes
}

// NewSearcher creates a new BleveSearcher with the given base directory
func NewSearcher(baseDir string) *BleveSearcher {
	return &BleveSearcher{baseDir: baseDir}
}

// Search finds matching types and fields in a schema
func (b *BleveSearcher) Search(schemaID string, query string, limit int) ([]SearchResult, error) {
	indexPath := b.getIndexPath(schemaID)

	// Open the index
	index, err := bleve.Open(indexPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open index: %w", err)
	}
	defer index.Close()

	// Create separate match queries for each field with different boost values
	// Priority order: name (highest), type, description, path (lowest)
	searchRequest := newSearchRequest(query)
	searchRequest.Size = limit

	// Execute search
	searchResults, err := index.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	// Convert results to SearchResult slice
	results := make([]SearchResult, 0, len(searchResults.Hits))
	for _, hit := range searchResults.Hits {
		result := SearchResult{
			Score: hit.Score,
		}

		// Extract fields from hit
		if typeVal, ok := hit.Fields["type"].(string); ok {
			result.Type = typeVal
		}
		if nameVal, ok := hit.Fields["name"].(string); ok {
			result.Name = nameVal
		}
		if descVal, ok := hit.Fields["description"].(string); ok {
			result.Description = descVal
		}
		if pathVal, ok := hit.Fields["path"].(string); ok {
			result.Path = pathVal
		}

		results = append(results, result)
	}

	return results, nil
}

func newSearchRequest(query string) *bleve.SearchRequest {
	// QueryStringQuery provides flexible, user-defined search configuration.
	// See https://blevesearch.com/docs/Query-String-Query/
	queryStringQuery := bleve.NewQueryStringQuery(query)
	searchRequest := bleve.NewSearchRequest(queryStringQuery)
	searchRequest.Fields = []string{"type", "name", "description", "path"}
	return searchRequest
}

// Close closes the searcher (no-op for BleveSearcher)
func (b *BleveSearcher) Close() error {
	return nil
}

// getIndexPath returns the path to the index directory for a schema
func (b *BleveSearcher) getIndexPath(schemaID string) string {
	return filepath.Join(b.baseDir, schemaID+".bleve")
}
