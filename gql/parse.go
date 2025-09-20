package gql

import (
	"fmt"
	"log"
	"strings"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
)

// GraphQLSchema represents the mapping of GraphQL type names to their field definitions
type GraphQLSchema map[string]map[string]*ast.FieldDefinition

func buildGraphQLTypes(doc *ast.Document) GraphQLSchema {
	gqlSchema := make(GraphQLSchema)

	for _, def := range doc.Definitions {
		switch typeDef := def.(type) {
		case *ast.ObjectDefinition:
			if typeDef.Name.Value == "Query" {
				queryMap := make(map[string]*ast.FieldDefinition)
				for _, field := range typeDef.Fields {
					queryMap[field.Name.Value] = field
				}
				gqlSchema["Query"] = queryMap
				break
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

func printQueries(queryFields map[string]*ast.FieldDefinition) {
	fmt.Printf("Found %d queries in the GitHub GraphQL schema:\n\n", len(queryFields))

	// Print all query fields
	i := 1
	for fieldName, field := range queryFields {
		printFieldDefinition(i, fieldName, field)
		i++
	}
}

func printFieldDefinition(index int, fieldName string, field *ast.FieldDefinition) {
	fmt.Printf("%2d. %-30s", index, fieldName)

	// Print arguments if any
	if len(field.Arguments) > 0 {
		fmt.Print("(")
		argStrs := make([]string, len(field.Arguments))
		for j, arg := range field.Arguments {
			argStrs[j] = arg.Name.Value + ": " + GetTypeString(arg.Type)
		}
		fmt.Print(strings.Join(argStrs, ", "))
		fmt.Print(")")
	}

	// Print return type
	fmt.Printf(" -> %s", GetTypeString(field.Type))

	// Print description if available
	if field.Description != nil {
		description := strings.ReplaceAll(field.Description.Value, "\n", " ")
		if len(description) > 80 {
			description = description[:77] + "..."
		}
		fmt.Printf("\n    %s", description)
	}

	fmt.Println()
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
