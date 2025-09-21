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

	// Comprehensive schema that includes all GraphQL definition types
	schemaString := `
		enum Status {
		  ACTIVE
		  INACTIVE
		  PENDING
		}

		scalar Date

		input CreateUserInput {
		  name: String!
		  email: String!
		  status: Status = ACTIVE
		}

		interface Node {
		  id: ID!
		}

		union SearchResult = User | Post

		directive @deprecated(reason: String = "No longer supported") on FIELD_DEFINITION | ENUM_VALUE

		type User implements Node {
		  id: ID!
		  name: String!
		  email: String!
		  status: Status!
		  createdAt: Date!
		  posts: [Post!]!
		}

		type Post implements Node {
		  id: ID!
		  title: String!
		  content: String!
		  status: Status!
		  author: User!
		}

		type Query {
		  """
		  Return all posts
		  """
		  getAllPosts: [Post!]!
		  getPostById(id: ID!): Post
		  searchAll(query: String!): [SearchResult!]!
		}

		type Mutation {
		  """
		  Create a new post
		  """
		  createPost(title: String!, content: String!, authorId: ID!): Post!
		  createUser(input: CreateUserInput!): User!
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

	t.Run("Query: searchAll", func(t *testing.T) {
		gqlField, ok := queryFields["searchAll"]
		is.True(ok)

		is.Equal(len(gqlField.Arguments), 1)
		assertArgumentNameAndType(gqlField.Arguments[0], "query", "String!")

		is.Equal(GetTypeString(gqlField.Type), "[SearchResult!]!")
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

	t.Run("Mutation: createUser", func(t *testing.T) {
		gqlField, ok := mutationFields["createUser"]
		is.True(ok)

		is.Equal(len(gqlField.Arguments), 1)
		assertArgumentNameAndType(gqlField.Arguments[0], "input", "CreateUserInput!")

		is.Equal(GetTypeString(gqlField.Type), "User!")
	})

	t.Run("Object types", func(t *testing.T) {
		userObj, ok := schema.Object["User"]
		is.True(ok)
		is.Equal(userObj.Name.Value, "User")
		is.Equal(userObj.Kind, "ObjectDefinition")
		is.Equal(len(userObj.Fields), 6) // id, name, email, status, createdAt, posts

		postObj, ok := schema.Object["Post"]
		is.True(ok)
		is.Equal(postObj.Name.Value, "Post")
		is.Equal(len(postObj.Fields), 5) // id, title, content, status, author
	})

	t.Run("Object types with interfaces", func(t *testing.T) {
		userObj, ok := schema.Object["User"]
		is.True(ok)
		is.Equal(len(userObj.Interfaces), 1)
		is.Equal(userObj.Interfaces[0].Name.Value, "Node")

		postObj, ok := schema.Object["Post"]
		is.True(ok)
		is.Equal(len(postObj.Interfaces), 1)
		is.Equal(postObj.Interfaces[0].Name.Value, "Node")
	})

	t.Run("Enum definitions", func(t *testing.T) {
		statusEnum, ok := schema.Enum["Status"]
		is.True(ok)
		is.Equal(statusEnum.Name.Value, "Status")
		is.Equal(statusEnum.Kind, "EnumDefinition")
		is.Equal(len(statusEnum.Values), 3)
		is.Equal(statusEnum.Values[0].Name.Value, "ACTIVE")
		is.Equal(statusEnum.Values[1].Name.Value, "INACTIVE")
		is.Equal(statusEnum.Values[2].Name.Value, "PENDING")
	})

	t.Run("Scalar definitions", func(t *testing.T) {
		dateScalar, ok := schema.Scalar["Date"]
		is.True(ok)
		is.Equal(dateScalar.Name.Value, "Date")
		is.Equal(dateScalar.Kind, "ScalarDefinition")
	})

	t.Run("Input object definitions", func(t *testing.T) {
		createUserInput, ok := schema.Input["CreateUserInput"]
		is.True(ok)
		is.Equal(createUserInput.Name.Value, "CreateUserInput")
		is.Equal(createUserInput.Kind, "InputObjectDefinition")
		is.Equal(len(createUserInput.Fields), 3) // name, email, status

		// Check field details
		nameField := createUserInput.Fields[0]
		is.Equal(nameField.Name.Value, "name")
		is.Equal(GetTypeString(nameField.Type), "String!")

		emailField := createUserInput.Fields[1]
		is.Equal(emailField.Name.Value, "email")
		is.Equal(GetTypeString(emailField.Type), "String!")

		statusField := createUserInput.Fields[2]
		is.Equal(statusField.Name.Value, "status")
		is.Equal(GetTypeString(statusField.Type), "Status")
	})

	t.Run("Interface definitions", func(t *testing.T) {
		nodeInterface, ok := schema.Interface["Node"]
		is.True(ok)
		is.Equal(nodeInterface.Name.Value, "Node")
		is.Equal(nodeInterface.Kind, "InterfaceDefinition")
		is.Equal(len(nodeInterface.Fields), 1) // id field

		idField := nodeInterface.Fields[0]
		is.Equal(idField.Name.Value, "id")
		is.Equal(GetTypeString(idField.Type), "ID!")
	})

	t.Run("Union definitions", func(t *testing.T) {
		searchResultUnion, ok := schema.Union["SearchResult"]
		is.True(ok)
		is.Equal(searchResultUnion.Name.Value, "SearchResult")
		is.Equal(searchResultUnion.Kind, "UnionDefinition")
		is.Equal(len(searchResultUnion.Types), 2) // User and Post

		is.Equal(searchResultUnion.Types[0].Name.Value, "User")
		is.Equal(searchResultUnion.Types[1].Name.Value, "Post")
	})

	t.Run("Directive definitions", func(t *testing.T) {
		deprecatedDirective, ok := schema.Directive["deprecated"]
		is.True(ok)
		is.Equal(deprecatedDirective.Name.Value, "deprecated")
		is.Equal(deprecatedDirective.Kind, "DirectiveDefinition")

		// Check directive arguments
		is.Equal(len(deprecatedDirective.Arguments), 1)
		reasonArg := deprecatedDirective.Arguments[0]
		is.Equal(reasonArg.Name.Value, "reason")
		is.Equal(GetTypeString(reasonArg.Type), "String")

		// Check directive locations
		is.True(len(deprecatedDirective.Locations) >= 2) // FIELD_DEFINITION and ENUM_VALUE
	})
}

func TestGetTypeString(t *testing.T) {
	is := is.New(t)

	schemaString := `
		type Query {
		  getString: String
		  getRequiredString: String!
		  getStringList: [String]
		  getRequiredStringList: [String]!
		  getListOfRequiredStrings: [String!]
		  getRequiredListOfRequiredStrings: [String!]!
		}
	`

	schema := ParseSchema([]byte(schemaString))

	testCases := []struct {
		fieldName    string
		expectedType string
	}{
		{"getString", "String"},
		{"getRequiredString", "String!"},
		{"getStringList", "[String]"},
		{"getRequiredStringList", "[String]!"},
		{"getListOfRequiredStrings", "[String!]"},
		{"getRequiredListOfRequiredStrings", "[String!]!"},
	}

	for _, tc := range testCases {
		t.Run(tc.fieldName, func(t *testing.T) {
			field, ok := schema.Query[tc.fieldName]
			is.True(ok)
			is.Equal(GetTypeString(field.Type), tc.expectedType)
		})
	}
}

func TestParseSchemaWithNullDefaults(t *testing.T) {
	is := is.New(t)

	// Test that the schema parser handles "= null" default values correctly
	schemaString := `
		input TestInput {
		  name: String!
		  description: String = null
		  count: Int = null
		}

		type Query {
		  test(input: TestInput!): String
		}
	`

	// This should not panic or fail to parse
	schema := ParseSchema([]byte(schemaString))

	testInput, ok := schema.Input["TestInput"]
	is.True(ok)
	is.Equal(testInput.Name.Value, "TestInput")
	is.Equal(len(testInput.Fields), 3)
}
