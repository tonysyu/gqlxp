package adapters

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/tui/xplr/navigation"
)

func TestSchemaView_FindKind(t *testing.T) {
	schemaContent := `
		type Query {
			user: User
		}

		type Mutation {
			createUser: User
		}

		type User {
			id: ID!
			name: String!
		}

		input UserInput {
			name: String!
		}

		enum Status {
			ACTIVE
			INACTIVE
		}

		scalar DateTime

		interface Node {
			id: ID!
		}

		union SearchResult = User

		directive @auth on FIELD_DEFINITION
	`

	schema, err := ParseSchemaString(schemaContent)
	if err != nil {
		t.Fatalf("ParseSchemaString() error = %v", err)
	}

	tests := []struct {
		name         string
		typeName     string
		wantCategory navigation.GQLKind
		wantFound    bool
	}{
		{
			name:         "Query type",
			typeName:     "Query",
			wantCategory: navigation.QueryKind,
			wantFound:    true,
		},
		{
			name:         "Mutation type",
			typeName:     "Mutation",
			wantCategory: navigation.MutationKind,
			wantFound:    true,
		},
		{
			name:         "Object type",
			typeName:     "User",
			wantCategory: navigation.ObjectKind,
			wantFound:    true,
		},
		{
			name:         "Input type",
			typeName:     "UserInput",
			wantCategory: navigation.InputKind,
			wantFound:    true,
		},
		{
			name:         "Enum type",
			typeName:     "Status",
			wantCategory: navigation.EnumKind,
			wantFound:    true,
		},
		{
			name:         "Scalar type",
			typeName:     "DateTime",
			wantCategory: navigation.ScalarKind,
			wantFound:    true,
		},
		{
			name:         "Interface type",
			typeName:     "Node",
			wantCategory: navigation.InterfaceKind,
			wantFound:    true,
		},
		{
			name:         "Union type",
			typeName:     "SearchResult",
			wantCategory: navigation.UnionKind,
			wantFound:    true,
		},
		{
			name:         "Directive type",
			typeName:     "auth",
			wantCategory: navigation.DirectiveKind,
			wantFound:    true,
		},
		{
			name:         "Non-existent type",
			typeName:     "NonExistent",
			wantCategory: "",
			wantFound:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			gotCategory, gotFound := schema.FindKind(tt.typeName)
			is.Equal(gotCategory, tt.wantCategory) // FindKind() category
			is.Equal(gotFound, tt.wantFound)       // FindKind() found
		})
	}
}
