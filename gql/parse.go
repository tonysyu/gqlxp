package gql

import (
	"log"
	"strings"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
)

// GraphQLSchema represents the GraphQL schema with Query and Mutation field definitions
type GraphQLSchema struct {
	Query    map[string]*ast.FieldDefinition
	Mutation map[string]*ast.FieldDefinition
}

func buildGraphQLTypes(doc *ast.Document) GraphQLSchema {
	gqlSchema := GraphQLSchema{
		Query:    make(map[string]*ast.FieldDefinition),
		Mutation: make(map[string]*ast.FieldDefinition),
	}

	for _, def := range doc.Definitions {
		switch typeDef := def.(type) {
		case *ast.ObjectDefinition:
			if typeDef.Name.Value == "Query" {
				for _, field := range typeDef.Fields {
					gqlSchema.Query[field.Name.Value] = field
				}
			} else if typeDef.Name.Value == "Mutation" {
				for _, field := range typeDef.Fields {
					gqlSchema.Mutation[field.Name.Value] = field
				}
			}
		}
	}

	return gqlSchema
}

func ParseSchema(schemaContent []byte) GraphQLSchema {
	// Clean up the schema content to remove problematic syntax
	// Nullable values are null by default, and explicit defaults results in parsing error
	cleanedSchema := strings.ReplaceAll(string(schemaContent), " = null", "")

	// Parse the schema
	doc, err := parser.Parse(parser.ParseParams{
		Source: cleanedSchema,
	})
	if err != nil {
		log.Fatalf("Failed to parse schema: %v", err)
	}

	gqlSchema := buildGraphQLTypes(doc)
	return gqlSchema
}

// Helper function to convert AST type to string representation
func GetTypeString(t ast.Type) string {
	switch typ := t.(type) {
	case *ast.Named:
		return typ.Name.Value
	case *ast.List:
		return "[" + GetTypeString(typ.Type) + "]"
	case *ast.NonNull:
		return GetTypeString(typ.Type) + "!"
	default:
		return "Unknown"
	}
}
