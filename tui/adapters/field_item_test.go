package adapters

import (
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/utils/testx"
	"github.com/tonysyu/gqlxp/utils/testx/assert"
)

func TestQueryAndMutationItemOpenPanel(t *testing.T) {
	is := is.New(t)
	assert := assert.New(t)

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

	schema, _ := gql.ParseSchema([]byte(schemaString))
	resolver := gql.NewSchemaResolver(&schema)

	t.Run("Query field with no arguments shows description and result type", func(t *testing.T) {
		field := schema.Query["getAllPosts"]
		item := newFieldItem(field, resolver)
		panel, _ := item.OpenPanel()

		// Set a reasonable size for testing
		panel.SetSize(80, 40)

		content := testx.NormalizeView(panel.View())

		expected := testx.NormalizeView(`
			  getAllPosts
			  Return all posts
			  Type
			  [Post!]!
		`)

		assert.StringContains(content, expected)
		is.True(!strings.Contains(content, "Inputs")) // Should not have arguments section
	})

	t.Run("Query field with arguments shows all sections", func(t *testing.T) {
		field := schema.Query["getPostById"]
		item := newFieldItem(field, resolver)
		panel, _ := item.OpenPanel()

		// Set a reasonable size for testing
		panel.SetSize(80, 40)

		// Verify Result Type tab (default active tab)
		content := renderMinimalPanel(panel)
		assert.StringContains(content, testx.NormalizeView(`
			Type    Inputs
			Post
		`))

		panel = nextPanelTab(panel) // Switch to Input Arguments tab
		content = renderMinimalPanel(panel)
		assert.StringContains(content, testx.NormalizeView(`
			Type    Inputs
			id: ID!
		`))
	})

	t.Run("Mutation field with multiple arguments shows all sections", func(t *testing.T) {
		field := schema.Mutation["createPost"]
		item := newFieldItem(field, resolver)
		panel, _ := item.OpenPanel()

		// Set a reasonable size for testing
		panel.SetSize(80, 40)

		// Verify Result Type tab (default active tab)
		content := renderMinimalPanel(panel)
		assert.StringContains(content, testx.NormalizeView(`
			createPost
			Create a new post

			Type    Inputs
			Post!
		`))

		panel = nextPanelTab(panel) // Switch to Input Arguments tab
		content = renderMinimalPanel(panel)
		assert.StringContains(content, testx.NormalizeView(`
			createPost
			Create a new post

			Type    Inputs
			title: String!
			content: String!
			authorId: ID!
		`))
	})
}

func TestFieldDefinitionWithoutDescription(t *testing.T) {
	is := is.New(t)
	assert := assert.New(t)

	schemaString := `
		type Query {
		  simpleField: String
		}
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))
	resolver := gql.NewSchemaResolver(&schema)

	field := schema.Query["simpleField"]
	item := newFieldItem(field, resolver)

	is.Equal(item.Title(), "simpleField: String")
	is.Equal(item.Description(), "") // No description
	is.Equal(item.FilterValue(), "simpleField")

	panel, ok := item.OpenPanel()
	is.True(ok)
	panel.SetSize(80, 40)

	content := testx.NormalizeView(panel.View())

	expected := testx.NormalizeView(`
		Type
		String
	`)

	assert.StringContains(content, expected)
	is.True(!strings.Contains(content, "Inputs"))
}

func TestFieldDefinitionWithComplexArguments(t *testing.T) {
	is := is.New(t)
	assert := assert.New(t)

	schemaString := `
		input FilterInput {
		  search: String
		  limit: Int
		}

		type Query {
		  complexField(
		    id: ID!
		    filters: FilterInput
		    tags: [String!]!
		    metadata: [String]
		  ): [String!]!
		}
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))
	resolver := gql.NewSchemaResolver(&schema)

	field := schema.Query["complexField"]
	item := newFieldItem(field, resolver)
	panel, ok := item.OpenPanel()

	is.True(ok)
	panel.SetSize(80, 40)

	// Verify Result Type tab (default active tab)
	content := renderMinimalPanel(panel)
	assert.StringContains(content, testx.NormalizeView(`
		Type    Inputs
		[String!]!
	`))

	panel = nextPanelTab(panel) // Switch to Input Arguments tab
	content = renderMinimalPanel(panel)
	assert.StringContains(content, testx.NormalizeView(`
		Type    Inputs
		id: ID!
		filters: FilterInput
		tags: [String!]!
		metadata: [String]
	`))
}
