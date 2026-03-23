package library

import "github.com/tonysyu/gqlxp/search"

// NewLibraryWithIndexer creates a Library with an injected indexer, for testing.
func NewLibraryWithIndexer(indexer search.Indexer) Library {
	return &FileLibrary{indexer: indexer}
}
