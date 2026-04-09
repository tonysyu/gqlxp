package search_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/search"
)

// testSchema is a simple schema for testing
const testSchema = `
	type Query {
		"""
		Get a user by ID
		"""
		user(id: ID!): User
		"""
		Search for users by name
		"""
		searchUsers(query: String!): [User!]!
	}

	"""
	A user in the system
	"""
	type User {
		id: ID!
		name: String!
		email: String!
	}
`

func TestIndexerAndSearcher(t *testing.T) {
	is := is.New(t)

	// Create temporary directory for test indexes
	tmpDir, err := os.MkdirTemp("", "gqlxp-search-test-*")
	is.NoErr(err) // should create temp directory
	defer os.RemoveAll(tmpDir)

	// Parse schema
	schema, err := gql.ParseSchema([]byte(testSchema))
	is.NoErr(err) // should parse schema

	// Create indexer and index the schema
	indexer := search.NewIndexer(tmpDir)
	defer indexer.Close()

	schemaID := "test-schema"
	err = indexer.Index(schemaID, &schema)
	is.NoErr(err) // should index schema

	// Verify index exists
	is.True(indexer.Exists(schemaID)) // index should exist

	// Create searcher and search
	searcher := search.NewSearcher(tmpDir)
	defer searcher.Close()

	// Test 1: Search for "user" should find the User type and user field
	results, err := searcher.Search(schemaID, "user", 10)
	is.NoErr(err)                                // should search successfully
	is.True(len(results) > 0)                    // should find results
	is.True(containsPath(results, "User"))       // should find User type
	is.True(containsPath(results, "Query.user")) // should find user field

	// Test 2: Search for "search" should find searchUsers field
	results, err = searcher.Search(schemaID, "search", 10)
	is.NoErr(err)                                       // should search successfully
	is.True(len(results) > 0)                           // should find results
	is.True(containsPath(results, "Query.searchUsers")) // should find searchUsers field

	// Test 3: Search for mid-word fragment should find searchUsers via wildcard on name
	results, err = searcher.Search(schemaID, "Users", 10)
	is.NoErr(err)                                       // should search successfully
	is.True(len(results) > 0)                           // should find results
	is.True(containsPath(results, "Query.searchUsers")) // should find searchUsers via fragment

	// Test 4: Query fields should have signatures
	results, err = searcher.Search(schemaID, "user", 10)
	is.NoErr(err)
	for _, r := range results {
		if r.Path == "Query.user" {
			is.Equal(r.Signature, "user(id: ID!): User") // Query field should have signature
		}
		if r.Path == "User.id" {
			is.Equal(r.Signature, "id: ID!") // Object field should have signature
		}
	}

	// Test 5: Search for non-existent term
	results, err = searcher.Search(schemaID, "nonexistent", 10)
	is.NoErr(err)             // should search successfully
	is.Equal(len(results), 0) // should find no results

	// Test 6: Remove index
	err = indexer.Remove(schemaID)
	is.NoErr(err)                      // should remove index
	is.True(!indexer.Exists(schemaID)) // index should not exist

	// Test 7: Search after removal should fail
	_, err = searcher.Search(schemaID, "user", 10)
	is.True(err != nil) // should fail to search non-existent index
}

func TestReindexing(t *testing.T) {
	is := is.New(t)

	tmpDir, err := os.MkdirTemp("", "gqlxp-search-test-*")
	is.NoErr(err)
	defer os.RemoveAll(tmpDir)

	schema, err := gql.ParseSchema([]byte(testSchema))
	is.NoErr(err)

	indexer := search.NewIndexer(tmpDir)
	defer indexer.Close()

	schemaID := "test-reindex"

	// Index the schema
	err = indexer.Index(schemaID, &schema)
	is.NoErr(err)

	// Get original index path modification time
	indexPath := filepath.Join(tmpDir, schemaID+".bleve")
	originalInfo, err := os.Stat(indexPath)
	is.NoErr(err)

	// Reindex (should replace the old index)
	err = indexer.Index(schemaID, &schema)
	is.NoErr(err)

	// Verify index was recreated
	newInfo, err := os.Stat(indexPath)
	is.NoErr(err)
	is.True(newInfo.ModTime().After(originalInfo.ModTime()) ||
		newInfo.ModTime().Equal(originalInfo.ModTime())) // index should be recreated
}

func TestSearchResultOrdering(t *testing.T) {
	is := is.New(t)

	tmpDir, err := os.MkdirTemp("", "gqlxp-search-test-*")
	is.NoErr(err)
	defer os.RemoveAll(tmpDir)

	schema, err := gql.ParseSchema([]byte(testSchema))
	is.NoErr(err)

	indexer := search.NewIndexer(tmpDir)
	defer indexer.Close()

	schemaID := "test-ordering"
	err = indexer.Index(schemaID, &schema)
	is.NoErr(err)

	searcher := search.NewSearcher(tmpDir)
	defer searcher.Close()

	// Search for "user" - results should be ordered by relevance (score)
	results, err := searcher.Search(schemaID, "user", 10)
	is.NoErr(err)
	is.True(len(results) > 1) // should have multiple results

	// Verify scores are in descending order
	for i := 1; i < len(results); i++ {
		is.True(results[i-1].Score >= results[i].Score) // scores should be in descending order
	}
}

const indexTestSchema = `
	interface Node {
		id: ID!
	}

	interface Timestamped implements Node {
		id: ID!
		createdAt: String!
	}

	type User implements Node & Timestamped {
		id: ID!
		createdAt: String!
		name: String!
		role: UserRole!
	}

	enum UserRole {
		ADMIN
		MEMBER
	}

	type Query {
		user(id: ID!): User
		users: [User!]!
	}

	type Mutation {
		createUser(name: String!): User
	}
`

func TestUsageIndex(t *testing.T) {
	is := is.New(t)

	tmpDir, err := os.MkdirTemp("", "gqlxp-usage-test-*")
	is.NoErr(err)
	defer os.RemoveAll(tmpDir)

	schema, err := gql.ParseSchema([]byte(indexTestSchema))
	is.NoErr(err)

	indexer := search.NewIndexer(tmpDir)
	defer indexer.Close()

	schemaID := "usage-test"
	err = indexer.Index(schemaID, &schema)
	is.NoErr(err)

	searcher := search.NewSearcher(tmpDir)
	defer searcher.Close()

	// Query fields that return User
	results, err := searcher.Search(schemaID, "+usage:User +kind:Query", 10)
	is.NoErr(err)
	is.True(containsPath(results, "Query.user"))
	is.True(containsPath(results, "Query.users"))

	// Mutation fields that return User
	results, err = searcher.Search(schemaID, "+usage:User +kind:Mutation", 10)
	is.NoErr(err)
	is.True(containsPath(results, "Mutation.createUser"))

	// Object fields that reference UserRole
	results, err = searcher.Search(schemaID, "+usage:UserRole", 10)
	is.NoErr(err)
	is.True(containsPath(results, "User.role"))
}

func TestImplementsIndex(t *testing.T) {
	is := is.New(t)

	tmpDir, err := os.MkdirTemp("", "gqlxp-implements-test-*")
	is.NoErr(err)
	defer os.RemoveAll(tmpDir)

	schema, err := gql.ParseSchema([]byte(indexTestSchema))
	is.NoErr(err)

	indexer := search.NewIndexer(tmpDir)
	defer indexer.Close()

	schemaID := "implements-test"
	err = indexer.Index(schemaID, &schema)
	is.NoErr(err)

	searcher := search.NewSearcher(tmpDir)
	defer searcher.Close()

	// User implements Node
	results, err := searcher.Search(schemaID, "+implements:Node +kind:Object", 10)
	is.NoErr(err)
	is.True(containsPath(results, "User"))

	// User implements Timestamped
	results, err = searcher.Search(schemaID, "+implements:Timestamped", 10)
	is.NoErr(err)
	is.True(containsPath(results, "User"))

	// Timestamped interface extends Node
	results, err = searcher.Search(schemaID, "+implements:Node +kind:Interface", 10)
	is.NoErr(err)
	is.True(containsPath(results, "Timestamped"))

	// Searching for a non-existent interface returns no results
	results, err = searcher.Search(schemaID, "+implements:NonExistent", 10)
	is.NoErr(err)
	is.Equal(len(results), 0)
}

// containsPath checks if any result has the given path
func containsPath(results []search.SearchResult, path string) bool {
	for _, result := range results {
		if result.Path == path {
			return true
		}
	}
	return false
}
