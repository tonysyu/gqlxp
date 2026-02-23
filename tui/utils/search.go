package utils

import (
	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/library"
	"github.com/tonysyu/gqlxp/search"
)

// EnsureSearchIndexForSchema creates the search index for a schema if it doesn't already exist.
func EnsureSearchIndexForSchema(schemaID string, schema *gql.GraphQLSchema) {
	schemasDir, err := library.GetSchemasDir()
	if err != nil {
		return
	}
	indexer := search.NewIndexer(schemasDir)
	defer indexer.Close()
	if !indexer.Exists(schemaID) {
		_ = indexer.Index(schemaID, schema)
	}
}
