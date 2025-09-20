package gql_test

import (
	"testing"

	"github.com/matryer/is"
	. "github.com/tonysyu/gq/gql"
)

func TestMain(t *testing.T) {
	is := is.New(t)
	schema := `
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
	`

	types := ParseSchema([]byte(schema))
	queryFields := types["Query"]

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
}
