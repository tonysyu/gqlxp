package gql

import (
	"fmt"
	"log"
	"strings"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
)

// GraphQLSchema represents the GraphQL schema with Query and Mutation field definitions
type GraphQLSchema struct {
	Query     map[string]*ast.FieldDefinition
	Mutation  map[string]*ast.FieldDefinition
	Object    map[string]*ast.ObjectDefinition
	Input     map[string]*ast.InputObjectDefinition
	Enum      map[string]*ast.EnumDefinition
	Scalar    map[string]*ast.ScalarDefinition
	Interface map[string]*ast.InterfaceDefinition
	Union     map[string]*ast.UnionDefinition
	Directive map[string]*ast.DirectiveDefinition
}

func buildGraphQLTypes(doc *ast.Document) GraphQLSchema {
	gqlSchema := GraphQLSchema{
		Query:     make(map[string]*ast.FieldDefinition),
		Mutation:  make(map[string]*ast.FieldDefinition),
		Object:    make(map[string]*ast.ObjectDefinition),
		Input:     make(map[string]*ast.InputObjectDefinition),
		Enum:      make(map[string]*ast.EnumDefinition),
		Scalar:    make(map[string]*ast.ScalarDefinition),
		Interface: make(map[string]*ast.InterfaceDefinition),
		Union:     make(map[string]*ast.UnionDefinition),
		Directive: make(map[string]*ast.DirectiveDefinition),
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
			} else {
				gqlSchema.Object[typeDef.Name.Value] = typeDef
			}
		case *ast.InputObjectDefinition:
			gqlSchema.Input[typeDef.Name.Value] = typeDef
		case *ast.EnumDefinition:
			gqlSchema.Enum[typeDef.Name.Value] = typeDef
		case *ast.ScalarDefinition:
			gqlSchema.Scalar[typeDef.Name.Value] = typeDef
		case *ast.InterfaceDefinition:
			gqlSchema.Interface[typeDef.Name.Value] = typeDef
		case *ast.UnionDefinition:
			gqlSchema.Union[typeDef.Name.Value] = typeDef
		case *ast.DirectiveDefinition:
			gqlSchema.Directive[typeDef.Name.Value] = typeDef
		case *ast.InputValueDefinition:
		// Ignore: Not sure what to do w/ input values right now
		default:
			fmt.Printf("Unknown type: %#v\n", typeDef)
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
