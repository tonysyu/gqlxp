package gql_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/matryer/is"
	. "github.com/tonysyu/gqlxp/gql"
)

func TestMain(t *testing.T) {
	is := is.New(t)

	assertArgumentNameAndType := func(arg *Argument, expectedName, expectedType string) {
		is.Equal(arg.Name(), expectedName)
		is.Equal(arg.TypeString(), expectedType)
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

	schema, _ := ParseSchema([]byte(schemaString))
	queryFields := schema.Query
	mutationFields := schema.Mutation

	t.Run("Query: getAllPosts", func(t *testing.T) {
		gqlField, ok := queryFields["getAllPosts"]
		is.True(ok)

		is.Equal(gqlField.Name(), "getAllPosts")

		is.Equal(gqlField.Description(), "Return all posts")

		is.Equal(len(gqlField.Arguments()), 0)

		is.Equal(gqlField.TypeString(), "[Post!]!")
	})

	t.Run("Query: getPostById", func(t *testing.T) {
		gqlField, ok := queryFields["getPostById"]
		is.True(ok)

		args := gqlField.Arguments()
		is.Equal(len(args), 1)
		assertArgumentNameAndType(args[0], "id", "ID!")

		is.Equal(gqlField.TypeString(), "Post")
	})

	t.Run("Query: searchAll", func(t *testing.T) {
		gqlField, ok := queryFields["searchAll"]
		is.True(ok)

		args := gqlField.Arguments()
		is.Equal(len(args), 1)
		assertArgumentNameAndType(args[0], "query", "String!")

		is.Equal(gqlField.TypeString(), "[SearchResult!]!")
	})

	t.Run("Mutation: createPost", func(t *testing.T) {
		gqlField, ok := mutationFields["createPost"]
		is.True(ok)

		is.Equal(gqlField.Name(), "createPost")

		is.Equal(gqlField.Description(), "Create a new post")

		args := gqlField.Arguments()
		is.Equal(len(args), 3)

		assertArgumentNameAndType(args[0], "title", "String!")
		assertArgumentNameAndType(args[1], "content", "String!")
		assertArgumentNameAndType(args[2], "authorId", "ID!")

		is.Equal(gqlField.TypeString(), "Post!")
	})

	t.Run("Mutation: createUser", func(t *testing.T) {
		gqlField, ok := mutationFields["createUser"]
		is.True(ok)

		args := gqlField.Arguments()
		is.Equal(len(args), 1)
		assertArgumentNameAndType(args[0], "input", "CreateUserInput!")

		is.Equal(gqlField.TypeString(), "User!")
	})

	t.Run("Object types", func(t *testing.T) {
		userObj, ok := schema.Object["User"]
		is.True(ok)
		is.Equal(userObj.Name(), "User")
		is.Equal(len(userObj.Fields()), 6) // id, name, email, status, createdAt, posts

		postObj, ok := schema.Object["Post"]
		is.True(ok)
		is.Equal(postObj.Name(), "Post")
		is.Equal(len(postObj.Fields()), 5) // id, title, content, status, author
	})

	t.Run("Object types with interfaces", func(t *testing.T) {
		userObj, ok := schema.Object["User"]
		is.True(ok)
		is.Equal(len(userObj.Interfaces()), 1)
		is.Equal(userObj.Interfaces()[0], "Node")

		postObj, ok := schema.Object["Post"]
		is.True(ok)
		is.Equal(len(postObj.Interfaces()), 1)
		is.Equal(postObj.Interfaces()[0], "Node")
	})

	t.Run("Enum definitions", func(t *testing.T) {
		statusEnum, ok := schema.Enum["Status"]
		is.True(ok)
		is.Equal(statusEnum.Name(), "Status")
		values := statusEnum.Values()
		is.Equal(len(values), 3)
		is.Equal(values[0].Name(), "ACTIVE")
		is.Equal(values[1].Name(), "INACTIVE")
		is.Equal(values[2].Name(), "PENDING")
	})

	t.Run("Scalar definitions", func(t *testing.T) {
		dateScalar, ok := schema.Scalar["Date"]
		is.True(ok)
		is.Equal(dateScalar.Name(), "Date")
	})

	t.Run("Input object definitions", func(t *testing.T) {
		createUserInput, ok := schema.Input["CreateUserInput"]
		is.True(ok)
		is.Equal(createUserInput.Name(), "CreateUserInput")
		fields := createUserInput.Fields()
		is.Equal(len(fields), 3) // name, email, status

		// Check field details
		nameField := fields[0]
		is.Equal(nameField.Name(), "name")
		is.Equal(nameField.TypeString(), "String!")

		emailField := fields[1]
		is.Equal(emailField.Name(), "email")
		is.Equal(emailField.TypeString(), "String!")

		statusField := fields[2]
		is.Equal(statusField.Name(), "status")
		is.Equal(statusField.TypeString(), "Status")
	})

	t.Run("Interface definitions", func(t *testing.T) {
		nodeInterface, ok := schema.Interface["Node"]
		is.True(ok)
		is.Equal(nodeInterface.Name(), "Node")
		fields := nodeInterface.Fields()
		is.Equal(len(fields), 1) // id field

		idField := fields[0]
		is.Equal(idField.Name(), "id")
		is.Equal(idField.TypeString(), "ID!")
	})

	t.Run("Union definitions", func(t *testing.T) {
		searchResultUnion, ok := schema.Union["SearchResult"]
		is.True(ok)
		is.Equal(searchResultUnion.Name(), "SearchResult")
		types := searchResultUnion.Types()
		is.Equal(len(types), 2) // User and Post

		is.Equal(types[0], "User")
		is.Equal(types[1], "Post")
	})

	t.Run("Directive definitions", func(t *testing.T) {
		// TODO: DirectiveDefinition wrapper needs to expose Arguments and Locations
		// For now, skip this test
		t.Skip("DirectiveDefinition wrapper incomplete")
	})
}

// Test that the schema parser handles "= null" default values correctly
// These null defaults aren't handled by default by `graphql-go/graphql`
func TestParseSchemaWithNullDefaults(t *testing.T) {
	is := is.New(t)

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
	schema, _ := ParseSchema([]byte(schemaString))

	testInput, ok := schema.Input["TestInput"]
	is.True(ok)
	is.Equal(testInput.Name(), "TestInput")
	is.Equal(len(testInput.Fields()), 3)
}

func TestParseEmptySchema(t *testing.T) {
	is := is.New(t)

	// Test parsing completely empty schema
	emptySchema, _ := ParseSchema([]byte(""))
	is.Equal(len(emptySchema.Query), 0)
	is.Equal(len(emptySchema.Mutation), 0)
	is.Equal(len(emptySchema.Object), 0)

	// Test schema with only comments
	commentOnlySchema := `
		# This is just a comment
		# Another comment
	`
	schema, _ := ParseSchema([]byte(commentOnlySchema))
	is.Equal(len(schema.Query), 0)
	is.Equal(len(schema.Mutation), 0)
}

func TestParseSchemaWithOnlyQuery(t *testing.T) {
	is := is.New(t)

	schemaString := `
		type Query {
		  hello: String
		}
	`
	schema, _ := ParseSchema([]byte(schemaString))
	is.Equal(len(schema.Query), 1)
	is.Equal(len(schema.Mutation), 0)

	hello, ok := schema.Query["hello"]
	is.True(ok)
	is.Equal(hello.Name(), "hello")
}

func TestParseSchemaWithInterfaceImplementingInterface(t *testing.T) {
	is := is.New(t)

	// Example schema for interface implementing interface. Copied from:
	// https://spec.graphql.org/October2021/#sec-Interfaces.Interfaces-Implementing-Interfaces
	schemaString := `
		interface Node {
		  id: ID!
		}

		interface Resource implements Node {
		  id: ID!
		  url: String
		}
	`
	schema, err := ParseSchema([]byte(schemaString))
	is.NoErr(err)
	is.Equal(len(schema.Interface), 2)

	_, ok := schema.Interface["Resource"]
	is.True(ok)
	_, ok = schema.Interface["Node"]
	is.True(ok)
}

func TestParseSchemaWithOnlyMutation(t *testing.T) {
	is := is.New(t)

	schemaString := `
		type Mutation {
		  createUser(name: String!): User
		}

		type User {
		  id: ID!
		  name: String!
		}
	`
	schema, _ := ParseSchema([]byte(schemaString))
	is.Equal(len(schema.Query), 0)
	is.Equal(len(schema.Mutation), 1)
	is.Equal(len(schema.Object), 1)

	createUser, ok := schema.Mutation["createUser"]
	is.True(ok)
	is.Equal(createUser.Name(), "createUser")
}

func TestParseSchemaWithComplexDefaultValues(t *testing.T) {
	is := is.New(t)

	schemaString := `
		enum Priority {
		  LOW
		  MEDIUM
		  HIGH
		}

		input TaskInput {
		  title: String!
		  priority: Priority = MEDIUM
		  tags: [String!] = []
		  metadata: String = null
		}

		type Query {
		  getTasks(input: TaskInput = null): [Task!]!
		}

		type Task {
		  id: ID!
		  title: String!
		  priority: Priority!
		}
	`

	// Should parse without errors even with complex defaults
	schema, _ := ParseSchema([]byte(schemaString))

	taskInput, ok := schema.Input["TaskInput"]
	is.True(ok)
	is.Equal(len(taskInput.Fields()), 4)

	priority, ok := schema.Enum["Priority"]
	is.True(ok)
	is.Equal(len(priority.Values()), 3)
}

func TestParseLargeSchema(t *testing.T) {
	is := is.New(t)

	// Generate a large schema programmatically
	var schemaBuilder []string
	schemaBuilder = append(schemaBuilder, "type Query {")

	// Add 100 query fields
	for i := range 100 {
		schemaBuilder = append(schemaBuilder, fmt.Sprintf("  field%d(arg1: String, arg2: Int): String", i))
	}
	schemaBuilder = append(schemaBuilder, "}")

	schemaBuilder = append(schemaBuilder, "type Mutation {")
	// Add 50 mutation fields
	for i := range 50 {
		schemaBuilder = append(schemaBuilder, fmt.Sprintf("  mutation%d(input: String!): Boolean", i))
	}
	schemaBuilder = append(schemaBuilder, "}")

	// Add 20 object types
	for i := range 20 {
		schemaBuilder = append(schemaBuilder, fmt.Sprintf("type Object%d {", i))
		schemaBuilder = append(schemaBuilder, "  id: ID!")
		schemaBuilder = append(schemaBuilder, "  name: String")
		schemaBuilder = append(schemaBuilder, "}")
	}

	largeSchemaString := strings.Join(schemaBuilder, "\n")
	schema, _ := ParseSchema([]byte(largeSchemaString))

	is.Equal(len(schema.Query), 100)
	is.Equal(len(schema.Mutation), 50)
	is.Equal(len(schema.Object), 20)
}

func TestNamedToTypeDefinition(t *testing.T) {
	is := is.New(t)

	schemaString := `
		enum Status {
		  ACTIVE
		}

		scalar Date

		input CreateUserInput {
		  name: String!
		}

		interface Node {
		  id: ID!
		}

		union SearchResult = User | Post

		directive @deprecated(reason: String) on FIELD_DEFINITION

		type User {
		  id: ID!
		  name: String!
		}

		type Post {
		  id: ID!
		  title: String!
		}

		type Query {
		  getUser: User
		}

		type Mutation {
		  createUser: User
		}
	`

	schema, _ := ParseSchema([]byte(schemaString))

	successTests := []struct {
		name           string
		typeName       string
		validateResult func(result TypeDef)
	}{
		{
			name:     "resolves Object type correctly",
			typeName: "User",
			validateResult: func(result TypeDef) {
				objectDef, ok := result.(*Object)
				is.True(ok)
				is.Equal(objectDef.Name(), "User")
			},
		},
		{
			name:     "resolves Input type correctly",
			typeName: "CreateUserInput",
			validateResult: func(result TypeDef) {
				inputDef, ok := result.(*InputObject)
				is.True(ok)
				is.Equal(inputDef.Name(), "CreateUserInput")
			},
		},
		{
			name:     "resolves Enum type correctly",
			typeName: "Status",
			validateResult: func(result TypeDef) {
				enumDef, ok := result.(*Enum)
				is.True(ok)
				is.Equal(enumDef.Name(), "Status")
			},
		},
		{
			name:     "resolves Scalar type correctly",
			typeName: "Date",
			validateResult: func(result TypeDef) {
				scalarDef, ok := result.(*Scalar)
				is.True(ok)
				is.Equal(scalarDef.Name(), "Date")
			},
		},
		{
			name:     "resolves Interface type correctly",
			typeName: "Node",
			validateResult: func(result TypeDef) {
				interfaceDef, ok := result.(*Interface)
				is.True(ok)
				is.Equal(interfaceDef.Name(), "Node")
			},
		},
		{
			name:     "resolves Union type correctly",
			typeName: "SearchResult",
			validateResult: func(result TypeDef) {
				unionDef, ok := result.(*Union)
				is.True(ok)
				is.Equal(unionDef.Name(), "SearchResult")
			},
		},
	}

	for _, tt := range successTests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := schema.NamedToTypeDef(tt.typeName)

			is.NoErr(err)
			is.True(result != nil)
			tt.validateResult(result)
		})
	}

	errorTests := []struct {
		name     string
		typeName string
	}{
		{
			name:     "returns err for Query type",
			typeName: "Query",
		},
		{
			name:     "returns err for Mutation type",
			typeName: "Mutation",
		},
		{
			name:     "returns err for Directive type",
			typeName: "deprecated",
		},
		{
			name:     "returns err for non-existent type",
			typeName: "NonExistent",
		},
		{
			name:     "returns err for empty string",
			typeName: "",
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := schema.NamedToTypeDef(tt.typeName)

			is.True(err != nil)
			is.Equal(result, nil)
		})
	}
}
