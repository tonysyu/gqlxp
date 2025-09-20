package gql_test

import (
	"testing"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/matryer/is"
	. "github.com/tonysyu/gq/gql"
)

func TestMain(t *testing.T) {
	is := is.New(t)

	assertArgumentNameAndType := func(arg *ast.InputValueDefinition, expectedName, expectedType string) {
		is.Equal(arg.Name.Value, expectedName)
		is.Equal(GetTypeString(arg.Type), expectedType)
	}

	schemaString := `
		type User {
		  id: ID!
		  name: String!
		  email: String!
		  posts: [Post!]!
		}

		type Post {
		  id: ID!
		  title: String!
		  content: String!
		  author: User!
		}

		type Query {
		  """
		  Return all posts
		  """
		  getAllPosts: [Post!]!
		  getPostById(id: ID!): Post
		}

		type Mutation {
		  """
		  Create a new post
		  """
		  createPost(title: String!, content: String!, authorId: ID!): Post!
		}
	`

	schema := ParseSchema([]byte(schemaString))
	queryFields := schema.Query
	mutationFields := schema.Mutation

	t.Run("Query: getAllPosts", func(t *testing.T) {
		gqlField, ok := queryFields["getAllPosts"]
		is.True(ok)

		is.Equal(gqlField.Name.Value, "getAllPosts")
		is.Equal(gqlField.Kind, "FieldDefinition")

		is.Equal(gqlField.Description.Value, "Return all posts")

		is.Equal(len(gqlField.Arguments), 0)

		is.Equal(GetTypeString(gqlField.Type), "[Post!]!")
	})
	t.Run("Query: getPostById", func(t *testing.T) {
		gqlField, ok := queryFields["getPostById"]
		is.True(ok)

		is.Equal(len(gqlField.Arguments), 1)
		assertArgumentNameAndType(gqlField.Arguments[0], "id", "ID!")

		is.Equal(GetTypeString(gqlField.Type), "Post")
	})

	t.Run("Mutation: createPost", func(t *testing.T) {
		gqlField, ok := mutationFields["createPost"]
		is.True(ok)

		is.Equal(gqlField.Name.Value, "createPost")
		is.Equal(gqlField.Kind, "FieldDefinition")

		is.Equal(gqlField.Description.Value, "Create a new post")

		is.Equal(len(gqlField.Arguments), 3)

		assertArgumentNameAndType(gqlField.Arguments[0], "title", "String!")
		assertArgumentNameAndType(gqlField.Arguments[1], "content", "String!")
		assertArgumentNameAndType(gqlField.Arguments[2], "authorId", "ID!")

		is.Equal(GetTypeString(gqlField.Type), "Post!")
	})
}
