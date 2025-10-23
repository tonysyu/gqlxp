package gql

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
)

// GraphQLSchema represents the GraphQL schema with Query and Mutation field definitions
type GraphQLSchema struct {
	Query      map[string]*ast.FieldDefinition
	Mutation   map[string]*ast.FieldDefinition
	Object     map[string]*ast.ObjectDefinition
	Input      map[string]*ast.InputObjectDefinition
	Enum       map[string]*ast.EnumDefinition
	Scalar     map[string]*ast.ScalarDefinition
	Interface  map[string]*ast.InterfaceDefinition
	Union      map[string]*ast.UnionDefinition
	Directive  map[string]*ast.DirectiveDefinition
	nameToType map[string]string
}

func buildGraphQLTypes(doc *ast.Document) GraphQLSchema {
	gqlSchema := GraphQLSchema{
		Query:      make(map[string]*ast.FieldDefinition),
		Mutation:   make(map[string]*ast.FieldDefinition),
		Object:     make(map[string]*ast.ObjectDefinition),
		Input:      make(map[string]*ast.InputObjectDefinition),
		Enum:       make(map[string]*ast.EnumDefinition),
		Scalar:     make(map[string]*ast.ScalarDefinition),
		Interface:  make(map[string]*ast.InterfaceDefinition),
		Union:      make(map[string]*ast.UnionDefinition),
		Directive:  make(map[string]*ast.DirectiveDefinition),
		nameToType: make(map[string]string),
	}

	for _, def := range doc.Definitions {
		switch typeDef := def.(type) {
		case *ast.ObjectDefinition:
			if typeDef.Name.Value == "Query" {
				for _, field := range typeDef.Fields {
					gqlSchema.Query[field.Name.Value] = field
				}
				gqlSchema.nameToType["Query"] = "Query"
			} else if typeDef.Name.Value == "Mutation" {
				for _, field := range typeDef.Fields {
					gqlSchema.Mutation[field.Name.Value] = field
				}
				gqlSchema.nameToType["Mutation"] = "Mutation"
			} else {
				gqlSchema.Object[typeDef.Name.Value] = typeDef
				gqlSchema.nameToType[typeDef.Name.Value] = "Object"
			}
		case *ast.InputObjectDefinition:
			gqlSchema.Input[typeDef.Name.Value] = typeDef
			gqlSchema.nameToType[typeDef.Name.Value] = "Input"
		case *ast.EnumDefinition:
			gqlSchema.Enum[typeDef.Name.Value] = typeDef
			gqlSchema.nameToType[typeDef.Name.Value] = "Enum"
		case *ast.ScalarDefinition:
			gqlSchema.Scalar[typeDef.Name.Value] = typeDef
			gqlSchema.nameToType[typeDef.Name.Value] = "Scalar"
		case *ast.InterfaceDefinition:
			gqlSchema.Interface[typeDef.Name.Value] = typeDef
			gqlSchema.nameToType[typeDef.Name.Value] = "Interface"
		case *ast.UnionDefinition:
			gqlSchema.Union[typeDef.Name.Value] = typeDef
			gqlSchema.nameToType[typeDef.Name.Value] = "Union"
		case *ast.DirectiveDefinition:
			gqlSchema.Directive[typeDef.Name.Value] = typeDef
			gqlSchema.nameToType[typeDef.Name.Value] = "Directive"
		case *ast.InputValueDefinition:
		// Ignore: Not sure what to do w/ input values right now
		default:
			fmt.Printf("Unknown type: %#v\n", typeDef)
		}
	}

	return gqlSchema
}

func ParseSchema(schemaContent []byte) (GraphQLSchema, error) {
	// Clean up the schema content to remove problematic syntax
	// Nullable values are null by default, and explicit defaults results in parsing error
	cleanedSchema := strings.ReplaceAll(string(schemaContent), " = null", "")

	// Parse the schema
	doc, err := parser.Parse(parser.ParseParams{
		Source: cleanedSchema,
	})
	if err != nil {
		return GraphQLSchema{}, err
	}

	gqlSchema := buildGraphQLTypes(doc)
	return gqlSchema, nil
}

// GetSortedQueryFields returns all query fields as wrapped FieldDefinitions, sorted by name.
func (s *GraphQLSchema) GetSortedQueryFields() []*FieldDefinition {
	astFields := CollectAndSortMapValues(s.Query)
	fields := make([]*FieldDefinition, 0, len(astFields))
	for _, field := range astFields {
		fields = append(fields, NewFieldDefinition(field))
	}
	return fields
}

// GetSortedMutationFields returns all mutation fields as wrapped FieldDefinitions, sorted by name.
func (s *GraphQLSchema) GetSortedMutationFields() []*FieldDefinition {
	astFields := CollectAndSortMapValues(s.Mutation)
	fields := make([]*FieldDefinition, 0, len(astFields))
	for _, field := range astFields {
		fields = append(fields, NewFieldDefinition(field))
	}
	return fields
}

// NamedToTypeDefinition resolves a Named type to its actual type definition.
// Returns (nil, nil) for nil input or special types (Query, Mutation, Directive) that don't have type definitions.
// Returns (nil, error) when the type name is not found in the schema.
func (s *GraphQLSchema) NamedToTypeDefinition(named *ast.Named) (NamedTypeDef, error) {
	if named == nil {
		return nil, fmt.Errorf("nil not supported for NamedToTypeDefinition")
	}

	typeName := named.Name.Value
	typeCategory, ok := s.nameToType[typeName]
	if !ok {
		return nil, fmt.Errorf("type %q not found in schema", typeName)
	}

	switch typeCategory {
	case "Query":
		return nil, fmt.Errorf("query type not supported")
	case "Mutation":
		return nil, fmt.Errorf("mutation type not supported")
	case "Object":
		return s.Object[typeName], nil
	case "Input":
		return s.Input[typeName], nil
	case "Enum":
		return s.Enum[typeName], nil
	case "Scalar":
		return s.Scalar[typeName], nil
	case "Interface":
		return s.Interface[typeName], nil
	case "Union":
		return s.Union[typeName], nil
	case "Directive":
		// Directive definitions don't implement NamedTypeDef interface
		return nil, fmt.Errorf("directive type not supported")
	default:
		return nil, fmt.Errorf("unknown type category %q for type %q", typeCategory, typeName)
	}
}
