package gql_test

import (
	"testing"

	"github.com/matryer/is"
	. "github.com/tonysyu/gqlxp/gql"
)

func TestSchemaResolver_ResolveType(t *testing.T) {
	is := is.New(t)
	schema, _ := ParseSchema([]byte(`
		type User {
			id: ID!
			name: String!
		}

		type Post {
			title: String!
		}

		union SearchResult = User | Post

		type Query {
			getUser(id: ID!): User
			search(query: String!): [SearchResult!]!
		}
	`))

	resolver := NewSchemaResolver(&schema)

	t.Run("ResolveType resolves Object type", func(t *testing.T) {
		typeDef, err := resolver.ResolveType("User")
		is.NoErr(err)
		is.Equal(typeDef.Name(), "User")

		_, ok := typeDef.(*Object)
		is.True(ok) // Expected Object type
	})

	t.Run("ResolveType resolves Union type", func(t *testing.T) {
		typeDef, err := resolver.ResolveType("SearchResult")
		is.NoErr(err)
		is.Equal(typeDef.Name(), "SearchResult")

		_, ok := typeDef.(*Union)
		is.True(ok) // Expected Union type
	})

	t.Run("ResolveType returns error for built-in scalar", func(t *testing.T) {
		_, err := resolver.ResolveType("String")
		is.True(err != nil) // Expected error for built-in scalar
	})

	t.Run("ResolveType returns error for non-existent type", func(t *testing.T) {
		_, err := resolver.ResolveType("NonExistent")
		is.True(err != nil) // Expected error for non-existent type
	})
}

func TestSchemaResolver_ResolveFieldType(t *testing.T) {
	is := is.New(t)
	schema, _ := ParseSchema([]byte(`
		type User {
			id: ID!
			name: String!
		}

		union SearchResult = User

		type Query {
			getUser(id: ID!): User
			search(query: String!): [SearchResult!]!
		}
	`))

	resolver := NewSchemaResolver(&schema)

	t.Run("ResolveFieldType resolves field returning Object", func(t *testing.T) {
		getUser := schema.Query["getUser"]
		typeDef, err := resolver.ResolveFieldType(getUser)
		is.NoErr(err)
		is.Equal(typeDef.Name(), "User")

		_, ok := typeDef.(*Object)
		is.True(ok) // Expected Object type
	})

	t.Run("ResolveFieldType resolves field returning Union in list", func(t *testing.T) {
		search := schema.Query["search"]
		typeDef, err := resolver.ResolveFieldType(search)
		is.NoErr(err)
		is.Equal(typeDef.Name(), "SearchResult")

		_, ok := typeDef.(*Union)
		is.True(ok) // Expected Union type
	})
}

func TestSchemaResolver_ResolveArgumentType(t *testing.T) {
	is := is.New(t)
	schema, _ := ParseSchema([]byte(`
		input UserInput {
			name: String!
		}

		type User {
			id: ID!
		}

		type Query {
			createUser(input: UserInput!): User
			getUser(id: ID!): User
		}
	`))

	resolver := NewSchemaResolver(&schema)

	t.Run("ResolveArgumentType resolves InputObject argument", func(t *testing.T) {
		createUser := schema.Query["createUser"]
		args := createUser.Arguments()
		is.Equal(len(args), 1)

		typeDef, err := resolver.ResolveArgumentType(args[0])
		is.NoErr(err)
		is.Equal(typeDef.Name(), "UserInput")

		_, ok := typeDef.(*InputObject)
		is.True(ok) // Expected InputObject type
	})

	t.Run("ResolveArgumentType returns error for built-in scalar argument", func(t *testing.T) {
		getUser := schema.Query["getUser"]
		args := getUser.Arguments()
		is.Equal(len(args), 1)

		_, err := resolver.ResolveArgumentType(args[0])
		is.True(err != nil) // Expected error for built-in scalar
	})
}
