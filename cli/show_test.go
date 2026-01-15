package cli

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/gqlfmt"
	"github.com/tonysyu/gqlxp/utils/testx/assert"
)

func TestGenerateMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		schema   string
		typeName string
		want     []string // substrings that should be present
		wantErr  bool
	}{
		{
			name: "Query field",
			schema: `
				type Query {
					"""Get user by ID"""
					getUser(id: ID!): User
				}
				type User { id: ID! }
			`,
			typeName: "Query.getUser",
			want:     []string{"# getUser", "getUser(id: ID!): User", "Get user by ID"},
		},
		{
			name: "Mutation field",
			schema: `
				type Query { placeholder: String }
				type Mutation {
					"""Create a new user"""
					createUser(name: String!): User!
				}
				type User { id: ID! }
			`,
			typeName: "Mutation.createUser",
			want:     []string{"# createUser", "createUser(name: String!): User!", "Create a new user"},
		},
		{
			name: "Object type",
			schema: `
				type Query { placeholder: String }
				"""A user in the system"""
				type User {
					"""Unique identifier"""
					id: ID!
					name: String!
				}
			`,
			typeName: "User",
			want:     []string{"# User", "A user in the system", "id: ID!", "name: String!", "Unique identifier"},
		},
		{
			name: "Object type with interface",
			schema: `
				type Query { placeholder: String }
				interface Node { id: ID! }
				"""A user in the system"""
				type User implements Node {
					id: ID!
					name: String!
				}
			`,
			typeName: "User",
			want:     []string{"# User", "A user in the system", "**Implements:** Node", "id: ID!", "name: String!"},
		},
		{
			name: "Input type",
			schema: `
				type Query { placeholder: String }
				"""Input for creating a user"""
				input CreateUserInput {
					"""User's name"""
					name: String!
					email: String!
				}
			`,
			typeName: "CreateUserInput",
			want:     []string{"# CreateUserInput", "Input for creating a user", "name: String!", "email: String!", "User's name"},
		},
		{
			name: "Enum type",
			schema: `
				type Query { placeholder: String }
				"""User role in the system"""
				enum Role {
					"""Administrator role"""
					ADMIN
					USER
				}
			`,
			typeName: "Role",
			want:     []string{"# Role", "User role in the system", "ADMIN", "USER", "Administrator role"},
		},
		{
			name: "Scalar type",
			schema: `
				type Query { placeholder: String }
				"""DateTime scalar type"""
				scalar DateTime
			`,
			typeName: "DateTime",
			want:     []string{"# DateTime", "DateTime scalar type", "_Scalar type_"},
		},
		{
			name: "Interface type",
			schema: `
				type Query { placeholder: String }
				"""Node interface"""
				interface Node {
					"""Unique identifier"""
					id: ID!
				}
			`,
			typeName: "Node",
			want:     []string{"# Node", "Node interface", "id: ID!", "Unique identifier"},
		},
		{
			name: "Union type",
			schema: `
				type Query { placeholder: String }
				"""Search result"""
				union SearchResult = User | Post
				type User { id: ID! }
				type Post { id: ID! }
			`,
			typeName: "SearchResult",
			want:     []string{"# SearchResult", "Search result", "**Union of:** User | Post"},
		},
		{
			name: "Directive",
			schema: `
				type Query { placeholder: String }
				"""Require authentication"""
				directive @auth(
					requires: String!
				) on FIELD_DEFINITION | OBJECT
			`,
			typeName: "@auth",
			want:     []string{"# @auth", "@auth(requires: String!)", "Require authentication", "**Locations:**", "FIELD_DEFINITION", "OBJECT"},
		},
		{
			name:     "Non-existent Query field",
			schema:   `type Query { getUser: String }`,
			typeName: "Query.nonExistent",
			wantErr:  true,
		},
		{
			name: "Non-existent Mutation field",
			schema: `
				type Query { placeholder: String }
				type Mutation { createUser: String }
			`,
			typeName: "Mutation.nonExistent",
			wantErr:  true,
		},
		{
			name:     "Non-existent type",
			schema:   `type Query { placeholder: String }`,
			typeName: "NonExistent",
			wantErr:  true,
		},
		{
			name:     "Non-existent directive",
			schema:   `type Query { placeholder: String }`,
			typeName: "@nonExistent",
			wantErr:  true,
		},
		{
			name: "Object.Field fallback to Object",
			schema: `
				type Query { placeholder: String }
				"""A user in the system"""
				type User {
					"""Unique identifier"""
					id: ID!
					name: String!
				}
			`,
			typeName: "User.name",
			want:     []string{"# User", "A user in the system", "id: ID!", "name: String!", "Unique identifier"},
		},
		{
			name:     "Non-existent Object.Field",
			schema:   `type Query { placeholder: String }`,
			typeName: "NonExistent.field",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			assert := assert.New(t)

			schema, err := gql.ParseSchema([]byte(tt.schema))
			if err != nil {
				t.Fatalf("Failed to parse schema: %v", err)
			}

			got, err := gqlfmt.GenerateMarkdown(schema, tt.typeName, gqlfmt.IncludeOptions{})
			is.Equal((err != nil), tt.wantErr) // generateMarkdown() error status

			if tt.wantErr {
				return
			}

			for _, want := range tt.want {
				assert.StringContains(got, want) // generateMarkdown() contains expected substring
			}
		})
	}
}

func TestGenerateFieldMarkdown(t *testing.T) {
	assert := assert.New(t)

	schema := `
		type Query {
			"""Get user by ID"""
			getUser(id: ID!): User
		}
		type User { id: ID! }
	`

	parsedSchema, err := gql.ParseSchema([]byte(schema))
	if err != nil {
		t.Fatalf("Failed to parse schema: %v", err)
	}

	field := parsedSchema.Query["getUser"]
	if field == nil {
		t.Fatal("Expected getUser field to exist")
	}

	markdown := gqlfmt.GenerateFieldMarkdown(field, nil)

	expectedSubstrings := []string{
		"# getUser",
		"getUser(id: ID!): User",
		"Get user by ID",
	}

	for _, expected := range expectedSubstrings {
		assert.StringContains(markdown, expected) // markdown contains expected substring
	}
}

func TestGenerateTypeDefMarkdown(t *testing.T) {
	assert := assert.New(t)

	schema := `
		type Query { placeholder: String }
		"""A user in the system"""
		type User {
			"""Unique identifier"""
			id: ID!
			name: String!
		}
	`

	parsedSchema, err := gql.ParseSchema([]byte(schema))
	if err != nil {
		t.Fatalf("Failed to parse schema: %v", err)
	}

	typeDef := parsedSchema.Object["User"]
	if typeDef == nil {
		t.Fatal("Expected User type to exist")
	}

	markdown := gqlfmt.GenerateTypeDefMarkdown(typeDef, nil)

	expectedSubstrings := []string{
		"# User",
		"A user in the system",
		"id: ID!",
		"name: String!",
		"Unique identifier",
	}

	for _, expected := range expectedSubstrings {
		assert.StringContains(markdown, expected) // markdown contains expected substring
	}
}

func TestGenerateJSON(t *testing.T) {
	tests := []struct {
		name     string
		schema   string
		typeName string
		want     []string // substrings that should be present in JSON
		wantErr  bool
	}{
		{
			name: "Query field JSON",
			schema: `
				type Query {
					"""Get user by ID"""
					getUser(id: ID!): User
				}
				type User { id: ID! }
			`,
			typeName: "Query.getUser",
			want:     []string{`"name": "getUser"`, `"type": "User"`, `"description": "Get user by ID"`},
		},
		{
			name: "Object type JSON",
			schema: `
				type Query { placeholder: String }
				"""A user in the system"""
				type User {
					id: ID!
					name: String!
				}
			`,
			typeName: "User",
			want:     []string{`"name": "User"`, `"kind": "Object"`, `"description": "A user in the system"`},
		},
		{
			name: "Enum type JSON",
			schema: `
				type Query { placeholder: String }
				enum Role {
					ADMIN
					USER
				}
			`,
			typeName: "Role",
			want:     []string{`"name": "Role"`, `"kind": "Enum"`, `"values"`},
		},
		{
			name:     "Non-existent type JSON",
			schema:   `type Query { placeholder: String }`,
			typeName: "NonExistent",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			assert := assert.New(t)

			schema, err := gql.ParseSchema([]byte(tt.schema))
			if err != nil {
				t.Fatalf("Failed to parse schema: %v", err)
			}

			got, err := gqlfmt.GenerateJSON(schema, tt.typeName, gqlfmt.IncludeOptions{})
			is.Equal((err != nil), tt.wantErr) // generateJSON() error status

			if tt.wantErr {
				return
			}

			for _, want := range tt.want {
				assert.StringContains(got, want) // generateJSON() contains expected substring
			}
		})
	}
}
