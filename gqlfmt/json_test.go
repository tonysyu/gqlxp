package gqlfmt

import (
	"encoding/json"
	"testing"

	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/utils/testx/assert"
)

func TestGenerateJSON(t *testing.T) {
	tests := []struct {
		name     string
		schema   string
		typeName string
		validate func(t *testing.T, jsonStr string)
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
			validate: func(t *testing.T, jsonStr string) {
				var result JSONField
				if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
					t.Fatalf("Failed to parse JSON: %v", err)
				}
				is := is.New(t)
				is.Equal(result.Name, "getUser")
				is.Equal(result.Kind, "Query")
				is.Equal(result.Type, "User")
				is.Equal(result.Description, "Get user by ID")
				is.Equal(len(result.Arguments), 1)
				is.Equal(result.Arguments[0].Name, "id")
				is.Equal(result.Arguments[0].Type, "ID!")
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
			validate: func(t *testing.T, jsonStr string) {
				var result JSONField
				if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
					t.Fatalf("Failed to parse JSON: %v", err)
				}
				is := is.New(t)
				is.Equal(result.Name, "createUser")
				is.Equal(result.Kind, "Mutation")
				is.Equal(result.Type, "User!")
				is.Equal(result.Description, "Create a new user")
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
			validate: func(t *testing.T, jsonStr string) {
				var result JSONTypeDef
				if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
					t.Fatalf("Failed to parse JSON: %v", err)
				}
				is := is.New(t)
				is.Equal(result.Name, "User")
				is.Equal(result.Kind, "Object")
				is.Equal(result.Description, "A user in the system")
				is.Equal(len(result.Fields), 2)
				is.Equal(result.Fields[0].Name, "id")
				is.Equal(result.Fields[0].Description, "Unique identifier")
			},
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
			validate: func(t *testing.T, jsonStr string) {
				var result JSONTypeDef
				if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
					t.Fatalf("Failed to parse JSON: %v", err)
				}
				is := is.New(t)
				is.Equal(result.Name, "User")
				is.Equal(result.Kind, "Object")
				is.Equal(len(result.Interfaces), 1)
				is.Equal(result.Interfaces[0], "Node")
			},
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
			validate: func(t *testing.T, jsonStr string) {
				var result JSONTypeDef
				if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
					t.Fatalf("Failed to parse JSON: %v", err)
				}
				is := is.New(t)
				is.Equal(result.Name, "Role")
				is.Equal(result.Kind, "Enum")
				is.Equal(result.Description, "User role in the system")
				is.Equal(len(result.Values), 2)
				is.Equal(result.Values[0].Name, "ADMIN")
				is.Equal(result.Values[0].Description, "Administrator role")
			},
		},
		{
			name: "Scalar type",
			schema: `
				type Query { placeholder: String }
				"""DateTime scalar type"""
				scalar DateTime
			`,
			typeName: "DateTime",
			validate: func(t *testing.T, jsonStr string) {
				var result JSONTypeDef
				if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
					t.Fatalf("Failed to parse JSON: %v", err)
				}
				is := is.New(t)
				is.Equal(result.Name, "DateTime")
				is.Equal(result.Kind, "Scalar")
				is.Equal(result.Description, "DateTime scalar type")
			},
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
			validate: func(t *testing.T, jsonStr string) {
				var result JSONTypeDef
				if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
					t.Fatalf("Failed to parse JSON: %v", err)
				}
				is := is.New(t)
				is.Equal(result.Name, "Node")
				is.Equal(result.Kind, "Interface")
				is.Equal(len(result.Fields), 1)
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
			validate: func(t *testing.T, jsonStr string) {
				var result JSONTypeDef
				if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
					t.Fatalf("Failed to parse JSON: %v", err)
				}
				is := is.New(t)
				is.Equal(result.Name, "SearchResult")
				is.Equal(result.Kind, "Union")
				is.Equal(len(result.Types), 2)
			},
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
			validate: func(t *testing.T, jsonStr string) {
				var result JSONDirective
				if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
					t.Fatalf("Failed to parse JSON: %v", err)
				}
				is := is.New(t)
				is.Equal(result.Name, "auth")
				is.Equal(result.Kind, "Directive")
				is.Equal(result.Description, "Require authentication")
				is.Equal(len(result.Locations), 2)
				is.Equal(len(result.Arguments), 1)
			},
		},
		{
			name:     "Non-existent Query field",
			schema:   `type Query { getUser: String }`,
			typeName: "Query.nonExistent",
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

			got, err := GenerateJSON(schema, tt.typeName, IncludeOptions{})
			is.Equal((err != nil), tt.wantErr) // GenerateJSON() error status

			if tt.wantErr {
				return
			}

			// Validate it's valid JSON
			var jsonCheck interface{}
			err = json.Unmarshal([]byte(got), &jsonCheck)
			is.NoErr(err) // JSON should be valid

			// Run custom validation if provided
			if tt.validate != nil {
				tt.validate(t, got)
			}
		})
	}
}

func TestJSONPrettyPrint(t *testing.T) {
	assert := assert.New(t)

	schema := `
		type Query {
			getUser(id: ID!): User
		}
		type User { id: ID! }
	`

	parsedSchema, err := gql.ParseSchema([]byte(schema))
	if err != nil {
		t.Fatalf("Failed to parse schema: %v", err)
	}

	jsonStr, err := GenerateJSON(parsedSchema, "Query.getUser", IncludeOptions{})
	if err != nil {
		t.Fatalf("Failed to generate JSON: %v", err)
	}

	// Check for 2-space indentation
	assert.StringContains(jsonStr, "\n  \"name\"") // JSON should be pretty-printed with 2-space indent
}
