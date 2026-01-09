package adapters

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/tui/xplr/navigation"
)

func TestSchemaView_FindTypeCategory(t *testing.T) {
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
		wantCategory navigation.GQLType
		wantFound    bool
	}{
		{
			name:         "Query type",
			typeName:     "Query",
			wantCategory: navigation.QueryType,
			wantFound:    true,
		},
		{
			name:         "Mutation type",
			typeName:     "Mutation",
			wantCategory: navigation.MutationType,
			wantFound:    true,
		},
		{
			name:         "Object type",
			typeName:     "User",
			wantCategory: navigation.ObjectType,
			wantFound:    true,
		},
		{
			name:         "Input type",
			typeName:     "UserInput",
			wantCategory: navigation.InputType,
			wantFound:    true,
		},
		{
			name:         "Enum type",
			typeName:     "Status",
			wantCategory: navigation.EnumType,
			wantFound:    true,
		},
		{
			name:         "Scalar type",
			typeName:     "DateTime",
			wantCategory: navigation.ScalarType,
			wantFound:    true,
		},
		{
			name:         "Interface type",
			typeName:     "Node",
			wantCategory: navigation.InterfaceType,
			wantFound:    true,
		},
		{
			name:         "Union type",
			typeName:     "SearchResult",
			wantCategory: navigation.UnionType,
			wantFound:    true,
		},
		{
			name:         "Directive type",
			typeName:     "auth",
			wantCategory: navigation.DirectiveType,
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

			gotCategory, gotFound := schema.FindTypeCategory(tt.typeName)
			is.Equal(gotCategory, tt.wantCategory) // FindTypeCategory() category
			is.Equal(gotFound, tt.wantFound)       // FindTypeCategory() found
		})
	}
}
