package search

import (
	"fmt"
	"path/filepath"

	"github.com/blevesearch/bleve/v2"
)

const (
	// Query scoring boosts for each document field that matches to prioritize some fields.
	// A value of 1 provides no boost.
	BOOST_NAME_FIELD = 2.0
	BOOST_TYPE_FIELD = 1.5
	BOOST_PATH_FIELD = 1.5
	BOOST_DESC_FIELD = 1.0
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
	nameQuery := bleve.NewMatchQuery(query)
	nameQuery.SetField("name")
	nameQuery.SetBoost(BOOST_NAME_FIELD)

	typeQuery := bleve.NewMatchQuery(query)
	typeQuery.SetField("type")
	typeQuery.SetBoost(BOOST_TYPE_FIELD)

	descQuery := bleve.NewMatchQuery(query)
	descQuery.SetField("description")
	descQuery.SetBoost(BOOST_DESC_FIELD)

	pathQuery := bleve.NewMatchQuery(query)
	pathQuery.SetField("path")
	pathQuery.SetBoost(BOOST_PATH_FIELD)

	// Combine queries using a BooleanQuery with should clauses
	// Documents matching any of these fields will be returned, with higher scores for higher-priority fields
	boolQuery := bleve.NewBooleanQuery()
	boolQuery.AddShould(nameQuery)
	boolQuery.AddShould(typeQuery)
	boolQuery.AddShould(descQuery)
	boolQuery.AddShould(pathQuery)
	boolQuery.SetMinShould(1)

	searchRequest := bleve.NewSearchRequest(boolQuery)
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
