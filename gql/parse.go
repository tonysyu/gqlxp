package gql

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
)

// GraphQLSchema represents the GraphQL schema with Query and Mutation field definitions
type GraphQLSchema struct {
	Query      map[string]*Field
	Mutation   map[string]*Field
	Object     map[string]*Object
	Input      map[string]*InputObject
	Enum       map[string]*Enum
	Scalar     map[string]*Scalar
	Interface  map[string]*Interface
	Union      map[string]*Union
	Directive  map[string]*Directive
	nameToType map[string]string
}

func buildGraphQLTypes(doc *ast.Document) GraphQLSchema {
	gqlSchema := GraphQLSchema{
		Query:      make(map[string]*Field),
		Mutation:   make(map[string]*Field),
		Object:     make(map[string]*Object),
		Input:      make(map[string]*InputObject),
		Enum:       make(map[string]*Enum),
		Scalar:     make(map[string]*Scalar),
		Interface:  make(map[string]*Interface),
		Union:      make(map[string]*Union),
		Directive:  make(map[string]*Directive),
		nameToType: make(map[string]string),
	}

	for _, def := range doc.Definitions {
		switch typeDef := def.(type) {
		case *ast.ObjectDefinition:
			if typeDef.Name.Value == "Query" {
				for _, field := range typeDef.Fields {
					gqlSchema.Query[field.Name.Value] = NewField(field)
				}
				gqlSchema.nameToType["Query"] = "Query"
			} else if typeDef.Name.Value == "Mutation" {
				for _, field := range typeDef.Fields {
					gqlSchema.Mutation[field.Name.Value] = NewField(field)
				}
				gqlSchema.nameToType["Mutation"] = "Mutation"
			} else {
				gqlSchema.Object[typeDef.Name.Value] = NewObject(typeDef)
				gqlSchema.nameToType[typeDef.Name.Value] = "Object"
			}
		case *ast.InputObjectDefinition:
			gqlSchema.Input[typeDef.Name.Value] = NewInputObject(typeDef)
			gqlSchema.nameToType[typeDef.Name.Value] = "Input"
		case *ast.EnumDefinition:
			gqlSchema.Enum[typeDef.Name.Value] = NewEnum(typeDef)
			gqlSchema.nameToType[typeDef.Name.Value] = "Enum"
		case *ast.ScalarDefinition:
			gqlSchema.Scalar[typeDef.Name.Value] = NewScalar(typeDef)
			gqlSchema.nameToType[typeDef.Name.Value] = "Scalar"
		case *ast.InterfaceDefinition:
			gqlSchema.Interface[typeDef.Name.Value] = NewInterface(typeDef)
			gqlSchema.nameToType[typeDef.Name.Value] = "Interface"
		case *ast.UnionDefinition:
			gqlSchema.Union[typeDef.Name.Value] = NewUnion(typeDef)
			gqlSchema.nameToType[typeDef.Name.Value] = "Union"
		case *ast.DirectiveDefinition:
			gqlSchema.Directive[typeDef.Name.Value] = NewDirective(typeDef)
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

// GetSortedQueryFields returns all query fields as wrapped Fields, sorted by name.
func (s *GraphQLSchema) GetSortedQueryFields() []*Field {
	return CollectAndSortMapValues(s.Query)
}

// GetSortedMutationFields returns all mutation fields as wrapped Fields, sorted by name.
func (s *GraphQLSchema) GetSortedMutationFields() []*Field {
	return CollectAndSortMapValues(s.Mutation)
}

// NamedToTypeDef resolves a Named type to its actual type definition.
// Returns (nil, nil) for nil input or special types (Query, Mutation, Directive) that don't have type definitions.
// Returns (nil, error) when the type name is not found in the schema.
func (s *GraphQLSchema) NamedToTypeDef(named *ast.Named) (NamedTypeDef, error) {
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
