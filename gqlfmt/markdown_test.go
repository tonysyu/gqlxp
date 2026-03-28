package gqlfmt

import (
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/gql"
)

func TestGenerateMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		schema   string
		typeName string
		opts     IncludeOptions
		validate func(t *testing.T, md string)
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
			validate: func(t *testing.T, md string) {
				is := is.New(t)
				is.True(strings.Contains(md, "# getUser"))
				is.True(strings.Contains(md, "Get user by ID"))
			},
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
			validate: func(t *testing.T, md string) {
				is := is.New(t)
				is.True(strings.Contains(md, "# createUser"))
				is.True(strings.Contains(md, "Create a new user"))
			},
		},
		{
			name: "Directive",
			schema: `
				type Query { placeholder: String }
				"""Require authentication"""
				directive @auth(requires: String!) on FIELD_DEFINITION | OBJECT
			`,
			typeName: "@auth",
			validate: func(t *testing.T, md string) {
				is := is.New(t)
				is.True(strings.Contains(md, "# @auth"))
				is.True(strings.Contains(md, "Require authentication"))
				is.True(strings.Contains(md, "FIELD_DEFINITION"))
			},
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
			validate: func(t *testing.T, md string) {
				is := is.New(t)
				is.True(strings.Contains(md, "# User"))
				is.True(strings.Contains(md, "A user in the system"))
				is.True(strings.Contains(md, "id: ID!"))
			},
		},
		{
			name: "Object type implements interface",
			schema: `
				type Query { placeholder: String }
				interface Node { id: ID! }
				type User implements Node { id: ID! }
			`,
			typeName: "User",
			validate: func(t *testing.T, md string) {
				is := is.New(t)
				is.True(strings.Contains(md, "**Implements:**"))
				is.True(strings.Contains(md, "Node"))
			},
		},
		{
			name: "Enum type",
			schema: `
				type Query { placeholder: String }
				"""User role"""
				enum Role {
					"""Administrator role"""
					ADMIN
					USER
				}
			`,
			typeName: "Role",
			validate: func(t *testing.T, md string) {
				is := is.New(t)
				is.True(strings.Contains(md, "# Role"))
				is.True(strings.Contains(md, "ADMIN"))
				is.True(strings.Contains(md, "USER"))
			},
		},
		{
			name: "Scalar type",
			schema: `
				type Query { placeholder: String }
				"""A date-time value"""
				scalar DateTime
			`,
			typeName: "DateTime",
			validate: func(t *testing.T, md string) {
				is := is.New(t)
				is.True(strings.Contains(md, "# DateTime"))
				is.True(strings.Contains(md, "_Scalar type_"))
			},
		},
		{
			name: "Interface type",
			schema: `
				type Query { placeholder: String }
				"""Node interface"""
				interface Node { id: ID! }
			`,
			typeName: "Node",
			validate: func(t *testing.T, md string) {
				is := is.New(t)
				is.True(strings.Contains(md, "# Node"))
				is.True(strings.Contains(md, "id: ID!"))
			},
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
			validate: func(t *testing.T, md string) {
				is := is.New(t)
				is.True(strings.Contains(md, "# SearchResult"))
				is.True(strings.Contains(md, "**Union of:**"))
				is.True(strings.Contains(md, "User"))
				is.True(strings.Contains(md, "Post"))
			},
		},
		{
			name: "Input object type",
			schema: `
				type Query { placeholder: String }
				input CreateUserInput {
					name: String!
					email: String!
				}
			`,
			typeName: "CreateUserInput",
			validate: func(t *testing.T, md string) {
				is := is.New(t)
				is.True(strings.Contains(md, "# CreateUserInput"))
				is.True(strings.Contains(md, "name: String!"))
			},
		},
		{
			name: "Query field includes return type definition when requested",
			schema: `
				type Query {
					"""Get user by ID"""
					getUser(id: ID!): User
				}
				type User {
					"""Unique identifier"""
					id: ID!
					name: String!
				}
			`,
			typeName: "Query.getUser",
			opts:     IncludeOptions{ReturnType: true},
			validate: func(t *testing.T, md string) {
				is := is.New(t)
				is.True(strings.Contains(md, "## Return Type"))
				is.True(strings.Contains(md, "# User"))
				is.True(strings.Contains(md, "id: ID!"))
			},
		},
		{
			name: "Query field does not include return type by default",
			schema: `
				type Query {
					"""Get user by ID"""
					getUser(id: ID!): User
				}
				type User { id: ID! }
			`,
			typeName: "Query.getUser",
			validate: func(t *testing.T, md string) {
				is := is.New(t)
				is.True(!strings.Contains(md, "## Return Type"))
			},
		},
		{
			name: "Query field with scalar return type has no return type section",
			schema: `
				type Query {
					"""Get a greeting"""
					greet(name: String!): String
				}
			`,
			typeName: "Query.greet",
			opts:     IncludeOptions{ReturnType: true},
			validate: func(t *testing.T, md string) {
				is := is.New(t)
				is.True(!strings.Contains(md, "## Return Type"))
			},
		},
		{
			name: "Usages included when requested",
			schema: `
				type Query {
					getUser(id: ID!): User
				}
				type User { id: ID! }
			`,
			typeName: "User",
			opts:     IncludeOptions{Usages: true},
			validate: func(t *testing.T, md string) {
				is := is.New(t)
				is.True(strings.Contains(md, "## Usages"))
				is.True(strings.Contains(md, "Query.getUser"))
			},
		},
		{
			name:     "Non-existent Query field",
			schema:   `type Query { getUser: String }`,
			typeName: "Query.nonExistent",
			wantErr:  true,
		},
		{
			name:     "Non-existent Mutation field",
			schema:   `type Query { placeholder: String }`,
			typeName: "Mutation.nonExistent",
			wantErr:  true,
		},
		{
			name:     "Non-existent directive",
			schema:   `type Query { placeholder: String }`,
			typeName: "@nonExistent",
			wantErr:  true,
		},
		{
			name:     "Non-existent type",
			schema:   `type Query { placeholder: String }`,
			typeName: "NonExistent",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			schema, err := gql.ParseSchema([]byte(tt.schema))
			if err != nil {
				t.Fatalf("Failed to parse schema: %v", err)
			}

			got, err := GenerateMarkdown(schema, tt.typeName, tt.opts)
			is.Equal((err != nil), tt.wantErr)

			if tt.wantErr {
				return
			}

			if tt.validate != nil {
				tt.validate(t, got)
			}
		})
	}
}
