package tui

import (
	"strings"
	"testing"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/matryer/is"
	"github.com/tonysyu/igq/gql"
)

func TestQueryAndMutationItemOpenPanel(t *testing.T) {
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

	schema, _ := gql.ParseSchema([]byte(schemaString))

	t.Run("Query field with no arguments shows description and result type", func(t *testing.T) {
		field := schema.Query["getAllPosts"]
		item := newFieldDefItem(field)
		panel, _ := item.Open()

		// Set a reasonable size for testing
		panel.SetSize(80, 40)

		content := panel.View()

		is.True(strings.Contains(content, "Return all posts"))
		is.True(strings.Contains(content, "======== Result Type ========"))
		is.True(strings.Contains(content, "[Post!]!"))
		is.True(!strings.Contains(content, "======== Input Arguments ========")) // Should not have arguments section
	})

	t.Run("Query field with arguments shows all sections", func(t *testing.T) {
		field := schema.Query["getPostById"]
		item := newFieldDefItem(field)
		panel, _ := item.Open()

		// Set a reasonable size for testing
		panel.SetSize(80, 40)

		content := panel.View()

		is.True(strings.Contains(content, "======== Input Arguments ========"))
		is.True(strings.Contains(content, "id: ID!"))
		is.True(strings.Contains(content, "======== Result Type ========"))
		is.True(strings.Contains(content, "Post"))
	})

	t.Run("Mutation field with multiple arguments shows all sections", func(t *testing.T) {
		field := schema.Mutation["createPost"]
		item := newFieldDefItem(field)
		panel, _ := item.Open()

		// Set a reasonable size for testing
		panel.SetSize(80, 40)

		content := panel.View()

		is.True(strings.Contains(content, "Create a new post"))
		is.True(strings.Contains(content, "======== Input Arguments ========"))
		is.True(strings.Contains(content, "title: String!"))
		is.True(strings.Contains(content, "content: String!"))
		is.True(strings.Contains(content, "authorId: ID!"))
		is.True(strings.Contains(content, "======== Result Type ========"))
		is.True(strings.Contains(content, "Post"))
	})
}

func TestObjectDefinitionItemOpenPanel(t *testing.T) {
	is := is.New(t)

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
	item := newTypeDefItem(userObj)
	panel, ok := item.Open()

	is.True(ok)
	panel.SetSize(80, 40)

	content := panel.View()

	// Object panels show field names, not their types (types are shown when opening individual fields)
	is.True(strings.Contains(content, "id"))
	is.True(strings.Contains(content, "name"))
	is.True(strings.Contains(content, "email"))
	is.True(strings.Contains(content, "posts"))
	is.True(strings.Contains(content, "4 items")) // Should show 4 fields
}

func TestInputDefinitionItemOpenPanel(t *testing.T) {
	is := is.New(t)

	schemaString := `
		input CreateUserInput {
		  name: String!
		  email: String!
		  age: Int = 18
		}
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))

	inputObj := schema.Input["CreateUserInput"]
	item := newTypeDefItem(inputObj)
	panel, ok := item.Open()

	is.True(ok)
	panel.SetSize(80, 40)

	content := panel.View()
	is.True(strings.Contains(content, "name: String!"))
	is.True(strings.Contains(content, "email: String!"))
	is.True(strings.Contains(content, "age: Int"))
}

func TestEnumDefinitionItemOpenPanel(t *testing.T) {
	is := is.New(t)

	schemaString := `
		enum Status {
		  ACTIVE
		  INACTIVE
		  PENDING
		}
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))

	enumObj := schema.Enum["Status"]
	item := newTypeDefItem(enumObj)
	panel, ok := item.Open()

	is.True(ok)
	panel.SetSize(80, 40)

	content := panel.View()
	is.True(strings.Contains(content, "ACTIVE"))
	is.True(strings.Contains(content, "INACTIVE"))
	is.True(strings.Contains(content, "PENDING"))
}

func TestScalarDefinitionItemOpenPanel(t *testing.T) {
	is := is.New(t)

	schemaString := "scalar Date"

	schema, _ := gql.ParseSchema([]byte(schemaString))

	scalarObj := schema.Scalar["Date"]
	item := newTypeDefItem(scalarObj)
	panel, ok := item.Open()

	is.True(ok)
	panel.SetSize(80, 40)

	// Scalar types should have minimal content (just the name)
	content := panel.View()
	is.True(len(content) > 0)
}

func TestInterfaceDefinitionItemOpenPanel(t *testing.T) {
	is := is.New(t)

	schemaString := `
		interface Node {
		  id: ID!
		  createdAt: String
		}
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))

	interfaceObj := schema.Interface["Node"]
	item := newTypeDefItem(interfaceObj)
	panel, ok := item.Open()

	is.True(ok)
	panel.SetSize(80, 40)

	content := panel.View()
	// Interface panels show field names, not their types (types are shown when opening individual fields)
	is.True(strings.Contains(content, "id"))
	is.True(strings.Contains(content, "createdAt"))
	is.True(strings.Contains(content, "2 items")) // Should show 2 fields
}

func TestUnionDefinitionItemOpenPanel(t *testing.T) {
	is := is.New(t)

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
	item := newTypeDefItem(unionObj)
	panel, ok := item.Open()

	is.True(ok)
	panel.SetSize(80, 40)

	content := panel.View()
	is.True(strings.Contains(content, "User"))
	is.True(strings.Contains(content, "Post"))
}

func TestFieldDefinitionWithoutDescription(t *testing.T) {
	is := is.New(t)

	schemaString := `
		type Query {
		  simpleField: String
		}
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))

	field := schema.Query["simpleField"]
	item := newFieldDefItem(field)

	is.Equal(item.Title(), "simpleField")
	is.Equal(item.Description(), "") // No description
	is.Equal(item.FilterValue(), "simpleField")

	panel, ok := item.Open()
	is.True(ok)
	panel.SetSize(80, 40)

	content := panel.View()
	is.True(strings.Contains(content, "======== Result Type ========"))
	is.True(strings.Contains(content, "String"))
	is.True(!strings.Contains(content, "======== Input Arguments ========"))
}

func TestFieldDefinitionWithComplexArguments(t *testing.T) {
	is := is.New(t)

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
	item := newFieldDefItem(field)
	panel, ok := item.Open()

	is.True(ok)
	panel.SetSize(80, 40)

	content := panel.View()
	is.True(strings.Contains(content, "======== Input Arguments ========"))
	is.True(strings.Contains(content, "id: ID!"))
	is.True(strings.Contains(content, "filters: FilterInput"))
	is.True(strings.Contains(content, "tags: [String!]!"))
	is.True(strings.Contains(content, "metadata: [String]"))
	is.True(strings.Contains(content, "======== Result Type ========"))
	is.True(strings.Contains(content, "[String!]!"))
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
	queryItems := adaptFieldDefinitions(gql.CollectAndSortMapValues(schema.Query))
	is.Equal(len(queryItems), 1)

	mutationItems := adaptFieldDefinitions(gql.CollectAndSortMapValues(schema.Mutation))
	is.Equal(len(mutationItems), 1)

	objectItems := adaptObjectDefinitions(gql.CollectAndSortMapValues(schema.Object))
	is.Equal(len(objectItems), 1)

	inputItems := adaptInputDefinitions(gql.CollectAndSortMapValues(schema.Input))
	is.Equal(len(inputItems), 1)

	enumItems := adaptEnumDefinitions(gql.CollectAndSortMapValues(schema.Enum))
	is.Equal(len(enumItems), 1)

	scalarItems := adaptScalarDefinitions(gql.CollectAndSortMapValues(schema.Scalar))
	is.Equal(len(scalarItems), 1)

	interfaceItems := adaptInterfaceDefinitions(gql.CollectAndSortMapValues(schema.Interface))
	is.Equal(len(interfaceItems), 1)

	unionItems := adaptUnionDefinitions(gql.CollectAndSortMapValues(schema.Union))
	is.Equal(len(unionItems), 1)

	directiveItems := adaptDirectiveDefinitions(gql.CollectAndSortMapValues(schema.Directive))
	is.Equal(len(directiveItems), 1)
}

func TestEmptyAdapterInputs(t *testing.T) {
	is := is.New(t)

	// Test adapters with empty inputs
	emptyFieldItems := adaptFieldDefinitions([]*ast.FieldDefinition{})
	is.Equal(len(emptyFieldItems), 0)

	emptyObjectItems := adaptObjectDefinitions([]*ast.ObjectDefinition{})
	is.Equal(len(emptyObjectItems), 0)

	emptyInputItems := adaptInputDefinitions([]*ast.InputObjectDefinition{})
	is.Equal(len(emptyInputItems), 0)

	emptyEnumItems := adaptEnumDefinitions([]*ast.EnumDefinition{})
	is.Equal(len(emptyEnumItems), 0)

	emptyScalarItems := adaptScalarDefinitions([]*ast.ScalarDefinition{})
	is.Equal(len(emptyScalarItems), 0)

	emptyInterfaceItems := adaptInterfaceDefinitions([]*ast.InterfaceDefinition{})
	is.Equal(len(emptyInterfaceItems), 0)

	emptyUnionItems := adaptUnionDefinitions([]*ast.UnionDefinition{})
	is.Equal(len(emptyUnionItems), 0)

	emptyDirectiveItems := adaptDirectiveDefinitions([]*ast.DirectiveDefinition{})
	is.Equal(len(emptyDirectiveItems), 0)
}

func TestSimpleItemInterface(t *testing.T) {
	is := is.New(t)

	item := simpleItem{
		title:       "Test Title",
		description: "Test Description",
	}

	is.Equal(item.Title(), "Test Title")
	is.Equal(item.Description(), "Test Description")
	is.Equal(item.FilterValue(), "Test Title")

	// Simple items should not be openable
	panel, ok := item.Open()
	is.True(!ok)
	is.True(panel == nil)
}

func TestSectionHeaderCreation(t *testing.T) {
	is := is.New(t)

	header := newSectionHeader("Test Section")
	is.Equal(header.Title(), "======== Test Section ========")
	is.Equal(header.Description(), "")
	is.Equal(header.FilterValue(), "======== Test Section ========")

	// Section headers should not be openable
	panel, ok := header.Open()
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
	items := adaptInputValueDefinitions(field.Arguments)
	is.Equal(len(items), 3)

	// Test first argument
	item1 := items[0].(simpleItem)
	is.Equal(item1.Title(), "arg1: String!")

	// Test second argument
	item2 := items[1].(simpleItem)
	is.Equal(item2.Title(), "arg2: Int")

	// Test third argument
	item3 := items[2].(simpleItem)
	is.Equal(item3.Title(), "arg3: [String]")
}

func TestTypeItemCreation(t *testing.T) {
	is := is.New(t)

	schemaString := `
		type Query {
		  simpleField: String
		  listField: [String!]!
		  complexField: [User]
		}

		type User {
		  id: ID!
		}
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))

	// Test type items for different field types
	simpleField := schema.Query["simpleField"]
	typeItem1 := newTypeItem(simpleField.Type)
	is.Equal(typeItem1.Title(), "String")

	listField := schema.Query["listField"]
	typeItem2 := newTypeItem(listField.Type)
	is.Equal(typeItem2.Title(), "[String!]!")

	complexField := schema.Query["complexField"]
	typeItem3 := newTypeItem(complexField.Type)
	is.Equal(typeItem3.Title(), "[User]")
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
	panel, ok := item.Open()
	is.True(!ok)
	is.True(panel == nil)
}

func TestInputValueDefinitionsEmpty(t *testing.T) {
	is := is.New(t)

	// Test with empty input value definitions
	items := adaptInputValueDefinitions([]*ast.InputValueDefinition{})
	is.Equal(len(items), 0)
}
