package gql_test

import (
	"testing"

	"github.com/matryer/is"
	. "github.com/tonysyu/gq/gql"
)

func TestMain(t *testing.T) {
	is := is.New(t)
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
		idArg := gqlField.Arguments[0]
		is.Equal(idArg.Name.Value, "id")

		is.Equal(GetTypeString(gqlField.Type), "Post")
	})

	t.Run("Mutation: createPost", func(t *testing.T) {
		gqlField, ok := mutationFields["createPost"]
		is.True(ok)

		is.Equal(gqlField.Name.Value, "createPost")
		is.Equal(gqlField.Kind, "FieldDefinition")

		is.Equal(gqlField.Description.Value, "Create a new post")

		is.Equal(len(gqlField.Arguments), 3)

		titleArg := gqlField.Arguments[0]
		is.Equal(titleArg.Name.Value, "title")
		is.Equal(GetTypeString(titleArg.Type), "String!")

		contentArg := gqlField.Arguments[1]
		is.Equal(contentArg.Name.Value, "content")
		is.Equal(GetTypeString(contentArg.Type), "String!")

		authorIdArg := gqlField.Arguments[2]
		is.Equal(authorIdArg.Name.Value, "authorId")
		is.Equal(GetTypeString(authorIdArg.Type), "ID!")

		is.Equal(GetTypeString(gqlField.Type), "Post!")
	})
}
