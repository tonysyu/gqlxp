package adapters

import (
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/tui/components"
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

	t.Run("Query field with no arguments shows description and result type", func(t *testing.T) {
		field := schema.Query["getAllPosts"]
		item := newFieldDefItem(field, &schema)
		panel, _ := item.OpenPanel()

		// Set a reasonable size for testing
		panel.SetSize(80, 40)

		content := testx.NormalizeView(panel.View())

		expected := testx.NormalizeView(`
			  getAllPosts
			  Return all posts
			  Result Type
			  │ [Post!]!
		`)

		assert.StringContains(content, expected)
		is.True(!strings.Contains(content, "Input Arguments")) // Should not have arguments section
	})

	t.Run("Query field with arguments shows all sections", func(t *testing.T) {
		field := schema.Query["getPostById"]
		item := newFieldDefItem(field, &schema)
		panel, _ := item.OpenPanel()

		// Set a reasonable size for testing
		panel.SetSize(80, 40)

		content := testx.RenderMinimalPanel(panel)
		assert.StringContains(content, testx.NormalizeView(`
			Result Type
		  │ Post

			Input Arguments
			id: ID!
		`))
	})

	t.Run("Mutation field with multiple arguments shows all sections", func(t *testing.T) {
		field := schema.Mutation["createPost"]
		item := newFieldDefItem(field, &schema)
		panel, _ := item.OpenPanel()

		// Set a reasonable size for testing
		panel.SetSize(80, 40)

		content := testx.RenderMinimalPanel(panel)
		assert.StringContains(content, testx.NormalizeView(`
			createPost
			Create a new post

			Result Type
		  │ Post!

			Input Arguments
			title: String!
			content: String!
			authorId: ID!
		`))

	})
}

func TestObjectDefinitionItemOpenPanel(t *testing.T) {
	is := is.New(t)
	assert := assert.New(t)

	schemaString := `
		type User {
		  id: ID!
		  name: String!
		  email: String
		  posts: [Post!]!
		}
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))

	userObj := schema.Object["User"]
	item := newTypeDefItem(userObj, &schema)
	panel, ok := item.OpenPanel()

	is.True(ok)
	panel.SetSize(80, 40)

	content := testx.RenderMinimalPanel(panel)

	assert.StringContains(content, testx.NormalizeView(`
		User

		id: ID!
		name: String!
		email: String
		posts: [Post!]!
	`))
}

func TestInputDefinitionItemOpenPanel(t *testing.T) {
	is := is.New(t)
	assert := assert.New(t)

	schemaString := `
		input CreateUserInput {
		  name: String!
		  email: String!
		  age: Int = 18
		}
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))

	inputObj := schema.Input["CreateUserInput"]
	item := newTypeDefItem(inputObj, &schema)
	panel, ok := item.OpenPanel()

	is.True(ok)
	panel.SetSize(80, 40)

	content := testx.RenderMinimalPanel(panel)

	assert.StringContains(content, testx.NormalizeView(`
		CreateUserInput

		name: String!
		email: String!
		age: Int = 18
	`))
}

func TestEnumDefinitionItemOpenPanel(t *testing.T) {
	is := is.New(t)
	assert := assert.New(t)

	schemaString := `
		enum Status {
		  ACTIVE
		  INACTIVE
		  PENDING
		}
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))

	enumObj := schema.Enum["Status"]
	item := newTypeDefItem(enumObj, &schema)
	panel, ok := item.OpenPanel()

	is.True(ok)
	panel.SetSize(80, 40)

	content := testx.RenderMinimalPanel(panel)
	assert.StringContains(content, testx.NormalizeView(`
		ACTIVE
		INACTIVE
		PENDING
	`))
}

func TestScalarDefinitionItemOpenPanel(t *testing.T) {
	is := is.New(t)

	schemaString := "scalar Date"

	schema, _ := gql.ParseSchema([]byte(schemaString))

	scalarObj := schema.Scalar["Date"]
	item := newTypeDefItem(scalarObj, &schema)
	panel, ok := item.OpenPanel()

	is.True(ok)
	panel.SetSize(80, 40)

	// Scalar types should have minimal content (just the name)
	content := panel.View()
	is.True(len(content) > 0)
}

func TestInterfaceDefinitionItemOpenPanel(t *testing.T) {
	is := is.New(t)
	assert := assert.New(t)

	schemaString := `
		interface Node {
		  id: ID!
		  createdAt: String
		}
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))

	interfaceObj := schema.Interface["Node"]
	item := newTypeDefItem(interfaceObj, &schema)
	panel, ok := item.OpenPanel()

	is.True(ok)
	panel.SetSize(80, 40)

	content := testx.RenderMinimalPanel(panel)

	assert.StringContains(content, testx.NormalizeView(`
		Node
		id: ID!
		createdAt: String
	`))
}

func TestUnionDefinitionItemOpenPanel(t *testing.T) {
	is := is.New(t)
	assert := assert.New(t)

	schemaString := `
		type User {
		  id: ID!
		}

		type Post {
		  id: ID!
		}

		union SearchResult = User | Post
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))

	unionObj := schema.Union["SearchResult"]
	item := newTypeDefItem(unionObj, &schema)
	panel, ok := item.OpenPanel()

	is.True(ok)
	panel.SetSize(80, 40)

	content := testx.RenderMinimalPanel(panel)

	assert.StringContains(content, testx.NormalizeView(`
		SearchResult
		User
		Post
	`))
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

	field := schema.Query["simpleField"]
	item := newFieldDefItem(field, &schema)

	is.Equal(item.Title(), "simpleField: String")
	is.Equal(item.Description(), "") // No description
	is.Equal(item.FilterValue(), "simpleField")

	panel, ok := item.OpenPanel()
	is.True(ok)
	panel.SetSize(80, 40)

	content := testx.NormalizeView(panel.View())

	expected := testx.NormalizeView(`
		  Result Type
		  │ String
	`)

	assert.StringContains(content, expected)
	is.True(!strings.Contains(content, "Input Arguments"))
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

	field := schema.Query["complexField"]
	item := newFieldDefItem(field, &schema)
	panel, ok := item.OpenPanel()

	is.True(ok)
	panel.SetSize(80, 40)

	content := testx.RenderMinimalPanel(panel)

	assert.StringContains(content, testx.NormalizeView(`
		Result Type
	  │ [String!]!

		Input Arguments
		id: ID!
		filters: FilterInput
		tags: [String!]!
		metadata: [String]
	`))
}

func TestSimpleItemInterface(t *testing.T) {
	is := is.New(t)

	item := components.NewSimpleItem("Test Title", components.WithDescription("Test Description"))

	is.Equal(item.Title(), "Test Title")
	is.Equal(item.Description(), "Test Description")
	is.Equal(item.FilterValue(), "Test Title")

	// Simple items should not be openable
	panel, ok := item.OpenPanel()
	is.True(!ok)
	is.True(panel == nil)
}

func TestArgumentListCreation(t *testing.T) {
	is := is.New(t)

	schemaString := `
		type Query {
		  testField(arg1: String!, arg2: Int, arg3: [String]): String
		}
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))
	field := schema.Query["testField"]

	// Test argument items creation
	items := adaptArguments(field.Arguments())
	is.Equal(len(items), 3)

	// Test first argument
	item1 := items[0].(components.SimpleItem)
	is.Equal(item1.Title(), "arg1: String!")

	// Test second argument
	item2 := items[1].(components.SimpleItem)
	is.Equal(item2.Title(), "arg2: Int")

	// Test third argument
	item3 := items[2].(components.SimpleItem)
	is.Equal(item3.Title(), "arg3: [String]")
}

func TestDirectiveDefinitionItemCreation(t *testing.T) {
	is := is.New(t)

	schemaString := `
		directive @deprecated(reason: String = "No longer supported") on FIELD_DEFINITION | ENUM_VALUE
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))
	directive := schema.Directive["deprecated"]

	item := newDirectiveDefinitionItem(directive)
	is.Equal(item.Title(), "deprecated")
	is.Equal(item.Description(), "")

	// Directive items should not be openable (they're simple items)
	panel, ok := item.OpenPanel()
	is.True(!ok)
	is.True(panel == nil)
}

func TestAdaptEmptyLists(t *testing.T) {
	is := is.New(t)

	// Test with empty arguments
	argItems := adaptArguments([]*gql.Argument{})
	is.Equal(len(argItems), 0)

	// Test with empty input fields
	fieldItems := adaptInputFields([]*gql.InputField{})
	is.Equal(len(fieldItems), 0)
}
