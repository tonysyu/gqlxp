package tui

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tonysyu/gq/gql"
)

func TestNewModel(t *testing.T) {
	is := is.New(t)

	// Create a basic schema for testing
	schemaString := `
		type Query {
			getAllPosts: [Post!]!
			getPostById(id: ID!): Post
		}

		type Mutation {
			createPost(title: String!, content: String!): Post!
		}

		type Post {
			id: ID!
			title: String!
		}
	`

	schema := gql.ParseSchema([]byte(schemaString))
	model := NewModel(schema)

	// Test initial state
	is.Equal(len(model.panels), intialPanels)
	is.Equal(model.focus, 0)
	is.Equal(model.fieldType, QueryType)
	is.Equal(len(model.schema.Query), 2)    // getAllPosts, getPostById
	is.Equal(len(model.schema.Mutation), 1) // createPost

	// Test that first panel is properly initialized with Query fields
	firstPanel := model.panels[0]
	// The first panel should be a list panel with Query fields
	is.True(firstPanel != nil)

	// Test keybindings are properly set
	is.True(model.keymap.next.Enabled())
	is.True(model.keymap.prev.Enabled())
	is.True(model.keymap.quit.Enabled())
	is.True(model.keymap.toggle.Enabled())
}
