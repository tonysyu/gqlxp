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
		item := newFieldDefItem(gql.NewFieldDefinition(field), &schema)
		panel, _ := item.Open()

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
		item := newFieldDefItem(gql.NewFieldDefinition(field), &schema)
		panel, _ := item.Open()

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
		item := newFieldDefItem(gql.NewFieldDefinition(field), &schema)
		panel, _ := item.Open()

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
	item := newTypeDefItem(gql.NewObjectDefinition(userObj), &schema)
	panel, ok := item.Open()

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
	item := newTypeDefItem(gql.NewInputObjectDefinition(inputObj), &schema)
	panel, ok := item.Open()

	is.True(ok)
	panel.SetSize(80, 40)

	content := testx.RenderMinimalPanel(panel)

	// FIXME: Render Default value for age
	assert.StringContains(content, testx.NormalizeView(`
		CreateUserInput

		name: String!
		email: String!
		age: Int
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
	item := newTypeDefItem(gql.NewEnumDefinition(enumObj), &schema)
	panel, ok := item.Open()

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
	item := newTypeDefItem(gql.NewScalarDefinition(scalarObj), &schema)
	panel, ok := item.Open()

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
	item := newTypeDefItem(gql.NewInterfaceDefinition(interfaceObj), &schema)
	panel, ok := item.Open()

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
	item := newTypeDefItem(gql.NewUnionDefinition(unionObj), &schema)
	panel, ok := item.Open()

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
	item := newFieldDefItem(gql.NewFieldDefinition(field), &schema)

	is.Equal(item.Title(), "simpleField: String")
	is.Equal(item.Description(), "") // No description
	is.Equal(item.FilterValue(), "simpleField")

	panel, ok := item.Open()
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
	item := newFieldDefItem(gql.NewFieldDefinition(field), &schema)
	panel, ok := item.Open()

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

func TestAdapterFunctions(t *testing.T) {
	is := is.New(t)

	schemaString := `
		type Query {
		  testField: String
		}

		type Mutation {
		  testMutation: String
		}

		type TestObject {
		  id: ID!
		}

		input TestInput {
		  name: String!
		}

		enum TestEnum {
		  VALUE_A
		  VALUE_B
		}

		scalar TestScalar

		interface TestInterface {
		  id: ID!
		}

		union TestUnion = TestObject

		directive @testDirective on FIELD_DEFINITION
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))

	// Test all adapter functions
	queryItems := AdaptFieldDefinitionsToItems(schema.GetSortedQueryFields(), &schema)
	is.Equal(len(queryItems), 1)

	mutationItems := AdaptFieldDefinitionsToItems(schema.GetSortedMutationFields(), &schema)
	is.Equal(len(mutationItems), 1)

	objectItems := AdaptObjectDefinitionsToItems(gql.WrapObjectDefinitions(gql.CollectAndSortMapValues(schema.Object)), &schema)
	is.Equal(len(objectItems), 1)

	inputItems := AdaptInputDefinitionsToItems(gql.WrapInputObjectDefinitions(gql.CollectAndSortMapValues(schema.Input)), &schema)
	is.Equal(len(inputItems), 1)

	enumItems := AdaptEnumDefinitionsToItems(gql.WrapEnumDefinitions(gql.CollectAndSortMapValues(schema.Enum)), &schema)
	is.Equal(len(enumItems), 1)

	scalarItems := AdaptScalarDefinitionsToItems(gql.WrapScalarDefinitions(gql.CollectAndSortMapValues(schema.Scalar)), &schema)
	is.Equal(len(scalarItems), 1)

	interfaceItems := AdaptInterfaceDefinitionsToItems(gql.WrapInterfaceDefinitions(gql.CollectAndSortMapValues(schema.Interface)), &schema)
	is.Equal(len(interfaceItems), 1)

	unionItems := AdaptUnionDefinitionsToItems(gql.WrapUnionDefinitions(gql.CollectAndSortMapValues(schema.Union)), &schema)
	is.Equal(len(unionItems), 1)

	directiveItems := AdaptDirectiveDefinitionsToItems(gql.WrapDirectiveDefinitions(gql.CollectAndSortMapValues(schema.Directive)))
	is.Equal(len(directiveItems), 1)
}

func TestEmptyAdapterInputs(t *testing.T) {
	is := is.New(t)

	// Test adapters with empty inputs
	emptyFieldItems := AdaptFieldDefinitionsToItems([]*gql.FieldDefinition{}, nil)
	is.Equal(len(emptyFieldItems), 0)

	emptyObjectItems := AdaptObjectDefinitionsToItems([]*gql.ObjectDefinition{}, nil)
	is.Equal(len(emptyObjectItems), 0)

	emptyInputItems := AdaptInputDefinitionsToItems([]*gql.InputObjectDefinition{}, nil)
	is.Equal(len(emptyInputItems), 0)

	emptyEnumItems := AdaptEnumDefinitionsToItems([]*gql.EnumDefinition{}, nil)
	is.Equal(len(emptyEnumItems), 0)

	emptyScalarItems := AdaptScalarDefinitionsToItems([]*gql.ScalarDefinition{}, nil)
	is.Equal(len(emptyScalarItems), 0)

	emptyInterfaceItems := AdaptInterfaceDefinitionsToItems([]*gql.InterfaceDefinition{}, nil)
	is.Equal(len(emptyInterfaceItems), 0)

	emptyUnionItems := AdaptUnionDefinitionsToItems([]*gql.UnionDefinition{}, nil)
	is.Equal(len(emptyUnionItems), 0)

	emptyDirectiveItems := AdaptDirectiveDefinitionsToItems([]*gql.DirectiveDefinition{})
	is.Equal(len(emptyDirectiveItems), 0)
}

func TestSimpleItemInterface(t *testing.T) {
	is := is.New(t)

	item := components.NewSimpleItem("Test Title", components.WithDescription("Test Description"))

	is.Equal(item.Title(), "Test Title")
	is.Equal(item.Description(), "Test Description")
	is.Equal(item.FilterValue(), "Test Title")

	// Simple items should not be openable
	panel, ok := item.Open()
	is.True(!ok)
	is.True(panel == nil)
}

func TestInputValueItemCreation(t *testing.T) {
	is := is.New(t)

	schemaString := `
		type Query {
		  testField(arg1: String!, arg2: Int, arg3: [String]): String
		}
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))
	field := schema.Query["testField"]

	// Test input value items creation
	items := adaptInputValueDefinitions(gql.WrapInputValueDefinitions(field.Arguments))
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

	item := newDirectiveDefinitionItem(gql.NewDirectiveDefinition(directive))
	is.Equal(item.Title(), "deprecated")
	is.Equal(item.Description(), "")

	// Directive items should not be openable (they're simple items)
	panel, ok := item.Open()
	is.True(!ok)
	is.True(panel == nil)
}

func TestInputValueDefinitionsEmpty(t *testing.T) {
	is := is.New(t)

	// Test with empty input value definitions
	items := adaptInputValueDefinitions([]*gql.InputValueDefinition{})
	is.Equal(len(items), 0)
}
