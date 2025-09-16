package gql

import (
	"fmt"
	"log"
	"strings"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
)

func ParseSchema(schemaContent []byte) {
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

	// Find the Query type and extract its fields
	var queryFields []*ast.FieldDefinition
	for _, def := range doc.Definitions {
		switch typeDef := def.(type) {
		case *ast.ObjectDefinition:
			if typeDef.Name.Value == "Query" {
				queryFields = typeDef.Fields
				break
			}
		}
	}

	fmt.Printf("Found %d queries in the GitHub GraphQL schema:\n\n", len(queryFields))

	// Print all query fields
	for i, field := range queryFields {
		fmt.Printf("%2d. %-30s", i+1, field.Name.Value)

		// Print arguments if any
		if len(field.Arguments) > 0 {
			fmt.Print("(")
			argStrs := make([]string, len(field.Arguments))
			for j, arg := range field.Arguments {
				argStrs[j] = arg.Name.Value + ": " + getTypeString(arg.Type)
			}
			fmt.Print(strings.Join(argStrs, ", "))
			fmt.Print(")")
		}

		// Print return type
		fmt.Printf(" -> %s", getTypeString(field.Type))

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
}

// Helper function to convert AST type to string representation
func getTypeString(t ast.Type) string {
	switch typ := t.(type) {
	case *ast.Named:
		return typ.Name.Value
	case *ast.List:
		return "[" + getTypeString(typ.Type) + "]"
	case *ast.NonNull:
		return getTypeString(typ.Type) + "!"
	default:
		return "Unknown"
	}
}
