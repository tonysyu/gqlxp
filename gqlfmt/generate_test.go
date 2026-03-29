package gqlfmt

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/gql"
)

func mustParseSchema(t *testing.T, schema string) gql.GraphQLSchema {
	t.Helper()
	parsed, err := gql.ParseSchema([]byte(schema))
	if err != nil {
		t.Fatalf("Failed to parse schema: %v", err)
	}
	return parsed
}

func TestGenerateOperation(t *testing.T) {
	tests := []struct {
		name     string
		schema   string
		field    string
		opts     GenerateOptions
		expected string
		wantErr  bool
	}{
		{
			name: "scalar return type",
			schema: `
				type Query { version: String }
			`,
			field: "Query.version",
			opts:  GenerateOptions{Depth: 1},
			expected: `query Version {
  version
}`,
		},
		{
			name: "query with non-null arg and object return",
			schema: `
				type Query { getUser(id: ID!): User }
				type User { id: ID!, name: String }
			`,
			field: "Query.getUser",
			opts:  GenerateOptions{Depth: 1},
			expected: `query GetUser($id: ID!) {
  getUser(id: $id) {
    id
    name
  }
}`,
		},
		{
			name: "nullable arg excluded from variables",
			schema: `
				type Query { search(q: String): String }
			`,
			field: "Query.search",
			opts:  GenerateOptions{Depth: 1},
			expected: `query Search {
  search
}`,
		},
		{
			name: "mutation",
			schema: `
				type Query { placeholder: String }
				type Mutation { createUser(name: String!): User }
				type User { id: ID! }
			`,
			field: "Mutation.createUser",
			opts:  GenerateOptions{Depth: 1},
			expected: `mutation CreateUser($name: String!) {
  createUser(name: $name) {
    id
  }
}`,
		},
		{
			name: "depth 0 omits object fields with comment",
			schema: `
				type Query { getUser(id: ID!): User }
				type User { id: ID!, profile: Profile }
				type Profile { bio: String }
			`,
			field: "Query.getUser",
			opts:  GenerateOptions{Depth: 0},
			expected: `query GetUser($id: ID!) {
  getUser(id: $id) {
    id
    # profile (Profile)
  }
}`,
		},
		{
			name: "depth 1 expands nested object",
			schema: `
				type Query { getUser(id: ID!): User }
				type User { id: ID!, profile: Profile }
				type Profile { bio: String, avatar: String }
			`,
			field: "Query.getUser",
			opts:  GenerateOptions{Depth: 1},
			expected: `query GetUser($id: ID!) {
  getUser(id: $id) {
    id
    profile {
      bio
      avatar
    }
  }
}`,
		},
		{
			name: "depth 1 comments out doubly-nested object fields",
			schema: `
				type Query { getUser(id: ID!): User }
				type User { id: ID!, profile: Profile }
				type Profile { bio: String, avatar: Image }
				type Image { url: String }
			`,
			field: "Query.getUser",
			opts:  GenerateOptions{Depth: 1},
			expected: `query GetUser($id: ID!) {
  getUser(id: $id) {
    id
    profile {
      bio
      # avatar (Image)
    }
  }
}`,
		},
		{
			name: "deprecated fields excluded by default",
			schema: `
				type Query { getUser: User }
				type User {
					id: ID!
					oldField: String @deprecated(reason: "use newField")
					newField: String
				}
			`,
			field: "Query.getUser",
			opts:  GenerateOptions{Depth: 1},
			expected: `query GetUser {
  getUser {
    id
    newField
  }
}`,
		},
		{
			name: "deprecated fields included with flag",
			schema: `
				type Query { getUser: User }
				type User {
					id: ID!
					oldField: String @deprecated(reason: "use newField")
					newField: String
				}
			`,
			field: "Query.getUser",
			opts:  GenerateOptions{Depth: 1, IncludeDeprecated: true},
			expected: `query GetUser {
  getUser {
    id
    oldField
    newField
  }
}`,
		},
		{
			name:    "invalid field path",
			schema:  `type Query { placeholder: String }`,
			field:   "User.getUser",
			opts:    GenerateOptions{Depth: 1},
			wantErr: true,
		},
		{
			name:    "unknown query field",
			schema:  `type Query { placeholder: String }`,
			field:   "Query.notFound",
			opts:    GenerateOptions{Depth: 1},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			schema := mustParseSchema(t, tt.schema)
			result, err := GenerateOperation(schema, tt.field, tt.opts)
			if tt.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)
			is.Equal(result, tt.expected)
		})
	}
}
