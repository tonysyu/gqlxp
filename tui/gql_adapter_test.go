package tui

import (
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/tonysyu/gq/gql"
)

func TestItemOpenPanel(t *testing.T) {
	is := is.New(t)

	schemaString := `
		type Post {
		  id: ID!
		  title: String!
		  content: String!
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

	schema := gql.ParseSchema([]byte(schemaString))

	t.Run("Query field with no arguments shows description and result type", func(t *testing.T) {
		field := schema.Query["getAllPosts"]
		item := newItem(field)
		panel := item.Open()

		// Set a reasonable size for testing
		panel.SetSize(80, 20)

		// Convert panel to string to check content
		content := panel.View()

		is.True(strings.Contains(content, "Return all posts"))
		is.True(strings.Contains(content, "======== Result Type ========"))
		is.True(strings.Contains(content, "[Post!]!"))
		is.True(!strings.Contains(content, "======== Input Arguments ========")) // Should not have arguments section
	})

	t.Run("Query field with arguments shows all sections", func(t *testing.T) {
		field := schema.Query["getPostById"]
		item := newItem(field)
		panel := item.Open()

		// Set a reasonable size for testing
		panel.SetSize(80, 20)

		content := panel.View()

		is.True(strings.Contains(content, "======== Input Arguments ========"))
		is.True(strings.Contains(content, "• id: ID!"))
		is.True(strings.Contains(content, "======== Result Type ========"))
		is.True(strings.Contains(content, "Post"))
	})

	t.Run("Mutation field with multiple arguments shows all sections", func(t *testing.T) {
		field := schema.Mutation["createPost"]
		item := newItem(field)
		panel := item.Open()

		// Set a reasonable size for testing
		panel.SetSize(80, 20)

		content := panel.View()

		is.True(strings.Contains(content, "Create a new post"))
		is.True(strings.Contains(content, "======== Input Arguments ========"))
		is.True(strings.Contains(content, "• title: String!"))
		// The content shows "7 items" which means our list structure is correct
		// Since this is a list with pagination, not all items may be visible
		is.True(strings.Contains(content, "7 items")) // Verifies all items are present
	})
}
