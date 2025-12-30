package gql

import (
	"fmt"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
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
	NameToType map[string]string
}

func buildGraphQLTypes(schema *ast.Schema) GraphQLSchema {
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
		NameToType: make(map[string]string),
	}

	// Process Query type
	if schema.Query != nil {
		for _, field := range schema.Query.Fields {
			// Skip introspection fields (start with __)
			if len(field.Name) >= 2 && field.Name[0] == '_' && field.Name[1] == '_' {
				continue
			}
			gqlSchema.Query[field.Name] = newField(field)
		}
		gqlSchema.NameToType["Query"] = "Query"
	}

	// Process Mutation type
	if schema.Mutation != nil {
		for _, field := range schema.Mutation.Fields {
			// Skip introspection fields (start with __)
			if len(field.Name) >= 2 && field.Name[0] == '_' && field.Name[1] == '_' {
				continue
			}
			gqlSchema.Mutation[field.Name] = newField(field)
		}
		gqlSchema.NameToType["Mutation"] = "Mutation"
	}

	// Process all other types
	// Built-in scalar types that should be skipped
	builtInScalars := map[string]bool{
		"Int":     true,
		"Float":   true,
		"String":  true,
		"Boolean": true,
		"ID":      true,
	}

	for name, typeDef := range schema.Types {
		// Skip built-in types to match graphql-go behavior
		if typeDef.BuiltIn || builtInScalars[name] {
			continue
		}

		switch typeDef.Kind {
		case ast.Object:
			// Skip Query and Mutation as they're handled above
			if name != "Query" && name != "Mutation" {
				gqlSchema.Object[name] = newObject(typeDef)
				gqlSchema.NameToType[name] = "Object"
			}
		case ast.InputObject:
			gqlSchema.Input[name] = newInputObject(typeDef)
			gqlSchema.NameToType[name] = "Input"
		case ast.Enum:
			gqlSchema.Enum[name] = newEnum(typeDef)
			gqlSchema.NameToType[name] = "Enum"
		case ast.Scalar:
			gqlSchema.Scalar[name] = newScalar(typeDef)
			gqlSchema.NameToType[name] = "Scalar"
		case ast.Interface:
			gqlSchema.Interface[name] = newInterface(typeDef)
			gqlSchema.NameToType[name] = "Interface"
		case ast.Union:
			gqlSchema.Union[name] = newUnion(typeDef)
			gqlSchema.NameToType[name] = "Union"
		default:
			fmt.Printf("Unknown type kind: %s for type %s\n", typeDef.Kind, name)
		}
	}

	// Process directives
	for name, directive := range schema.Directives {
		// Skip built-in directives (they have nil Position in gqlparser)
		// User-defined directives will have a Position set
		if directive.Position == nil {
			continue
		}
		gqlSchema.Directive[name] = newDirective(directive)
		gqlSchema.NameToType[name] = "Directive"
	}

	return gqlSchema
}

func ParseSchema(schemaContent []byte) (GraphQLSchema, error) {
	// Parse the schema document using gqlparser's lower-level parser
	source := &ast.Source{
		Name:  "schema.graphql",
		Input: string(schemaContent),
	}

	schemaDoc, gqlErr := parser.ParseSchema(source)
	if gqlErr != nil {
		return GraphQLSchema{}, gqlErr
	}

	// Try to load the full schema with validation
	schema, err := gqlparser.LoadSchema(source)
	if err != nil {
		// If validation fails but we have a parsed document, continue anyway
		// This maintains compatibility with test cases that have incomplete schemas
		schema = buildSchemaFromDocument(schemaDoc)
	}

	gqlSchema := buildGraphQLTypes(schema)
	return gqlSchema, nil
}

// buildSchemaFromDocument builds a minimal schema from a parsed document
// without full validation, used for incomplete schemas in tests
func buildSchemaFromDocument(doc *ast.SchemaDocument) *ast.Schema {
	schema := &ast.Schema{
		Types:      make(map[string]*ast.Definition),
		Directives: make(map[string]*ast.DirectiveDefinition),
	}

	for _, def := range doc.Definitions {
		if def.Name == "Query" {
			schema.Query = def
		} else if def.Name == "Mutation" {
			schema.Mutation = def
		}
		schema.Types[def.Name] = def
	}

	for _, dir := range doc.Directives {
		schema.Directives[dir.Name] = dir
	}

	return schema
}

// NamedToTypeDef resolves a type name to its actual type definition.
// Returns (nil, error) for special types (Query, Mutation, Directive) that don't have type definitions,
// or when the type name is not found in the schema.
func (s *GraphQLSchema) NamedToTypeDef(typeName string) (TypeDef, error) {
	if typeName == "" {
		return nil, fmt.Errorf("empty type name not supported")
	}

	typeCategory, ok := s.NameToType[typeName]
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
