package gql_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/gql"
)

func TestBuildUsageIndex_FieldReturnTypes(t *testing.T) {
	is := is.New(t)

	schemaContent := []byte(`
		type User {
			id: ID!
			name: String!
		}

		type Post {
			author: User!
			title: String!
		}

		type Query {
			user(id: ID!): User
			users: [User!]!
			post(id: ID!): Post
		}
	`)

	schema, err := gql.ParseSchema(schemaContent)
	is.NoErr(err)

	// Test User usages
	userUsages := schema.Usages["User"]
	is.True(userUsages != nil)
	is.Equal(len(userUsages), 3) // Query.user, Query.users, Post.author

	// Collect paths and verify all usages
	paths := make([]string, len(userUsages))
	for i, usage := range userUsages {
		paths[i] = usage.Path
	}

	// Verify all three paths exist
	is.True(contains(paths, "Query.user"))
	is.True(contains(paths, "Query.users"))
	is.True(contains(paths, "Post.author"))

	// Test Post usages
	postUsages := schema.Usages["Post"]
	is.True(postUsages != nil)
	is.Equal(len(postUsages), 1) // Query.post

	is.Equal(postUsages[0].ParentType, "Query")
	is.Equal(postUsages[0].FieldName, "post")
	is.Equal(postUsages[0].Path, "Query.post")
}

func TestBuildUsageIndex_MutationFields(t *testing.T) {
	is := is.New(t)

	schemaContent := []byte(`
		type User {
			id: ID!
			name: String!
		}

		type Mutation {
			createUser(name: String!): User!
			updateUser(id: ID!, name: String!): User
		}
	`)

	schema, err := gql.ParseSchema(schemaContent)
	is.NoErr(err)

	userUsages := schema.Usages["User"]
	is.True(userUsages != nil)
	is.Equal(len(userUsages), 2) // Mutation.createUser, Mutation.updateUser

	// Collect paths (order not guaranteed due to map iteration)
	paths := make([]string, len(userUsages))
	for i, usage := range userUsages {
		paths[i] = usage.Path
		// Verify all have correct ParentType and ParentKind
		is.Equal(usage.ParentType, "Mutation")
		is.Equal(usage.ParentKind, "Mutation")
	}

	// Verify both paths exist
	is.True(contains(paths, "Mutation.createUser"))
	is.True(contains(paths, "Mutation.updateUser"))
}

func TestBuildUsageIndex_NestedFields(t *testing.T) {
	is := is.New(t)

	schemaContent := []byte(`
		type User {
			id: ID!
			name: String!
			friends: [User!]!
		}

		type Query {
			user(id: ID!): User
		}
	`)

	schema, err := gql.ParseSchema(schemaContent)
	is.NoErr(err)

	userUsages := schema.Usages["User"]
	is.True(userUsages != nil)
	is.Equal(len(userUsages), 2) // Query.user, User.friends (self-reference)

	// Collect paths (order not guaranteed due to map iteration)
	paths := make([]string, len(userUsages))
	for i, usage := range userUsages {
		paths[i] = usage.Path
	}

	// Verify both paths exist
	is.True(contains(paths, "Query.user"))
	is.True(contains(paths, "User.friends"))
}

func TestBuildUsageIndex_WrappedTypes(t *testing.T) {
	is := is.New(t)

	schemaContent := []byte(`
		type User {
			id: ID!
		}

		type Query {
			user: User
			users: [User]
			requiredUser: User!
			requiredUsers: [User!]!
		}
	`)

	schema, err := gql.ParseSchema(schemaContent)
	is.NoErr(err)

	userUsages := schema.Usages["User"]
	is.True(userUsages != nil)
	is.Equal(len(userUsages), 4) // All four variants should be tracked

	paths := make([]string, len(userUsages))
	for i, usage := range userUsages {
		paths[i] = usage.Path
	}

	// All should reference User despite different wrapping
	is.True(contains(paths, "Query.user"))
	is.True(contains(paths, "Query.users"))
	is.True(contains(paths, "Query.requiredUser"))
	is.True(contains(paths, "Query.requiredUsers"))
}

func TestBuildUsageIndex_NoUsages(t *testing.T) {
	is := is.New(t)

	schemaContent := []byte(`
		type User {
			id: ID!
			name: String!
		}

		type Post {
			title: String!
		}

		type Query {
			post: Post
		}
	`)

	schema, err := gql.ParseSchema(schemaContent)
	is.NoErr(err)

	// User has no usages
	userUsages := schema.Usages["User"]
	is.Equal(len(userUsages), 0)

	// Post has one usage
	postUsages := schema.Usages["Post"]
	is.True(postUsages != nil)
	is.Equal(len(postUsages), 1)
}

func TestBuildUsageIndex_InterfaceFields(t *testing.T) {
	is := is.New(t)

	schemaContent := []byte(`
		type User {
			id: ID!
		}

		interface Node {
			id: ID!
			owner: User
		}

		type Query {
			node(id: ID!): Node
		}
	`)

	schema, err := gql.ParseSchema(schemaContent)
	is.NoErr(err)

	// User should be used by Interface field
	userUsages := schema.Usages["User"]
	is.True(userUsages != nil)
	is.Equal(len(userUsages), 1)

	is.Equal(userUsages[0].ParentType, "Node")
	is.Equal(userUsages[0].ParentKind, "Interface")
	is.Equal(userUsages[0].FieldName, "owner")
	is.Equal(userUsages[0].Path, "Node.owner")
}

// Helper function
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
