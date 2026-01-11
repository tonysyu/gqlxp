package gql

import "fmt"

// TypeResolver provides methods for resolving GraphQL type definitions
type TypeResolver interface {
	// ResolveType resolves a type name to its definition
	ResolveType(typeName string) (TypeDef, error)

	// ResolveFieldType resolves a field's result type
	ResolveFieldType(field *Field) (TypeDef, error)

	// ResolveArgumentType resolves an argument's input type
	ResolveArgumentType(arg *Argument) (TypeDef, error)

	// ResolveDirective resolves a directive name to its definition
	ResolveDirective(directiveName string) (*DirectiveDef, error)

	// ResolveUsages returns all usages of a given type name
	ResolveUsages(typeName string) ([]*Usage, error)

	// ResolveQueryOrMutationField resolves a field from Query or Mutation
	ResolveQueryOrMutationField(parentType, fieldName string) (*Field, error)
}

// SchemaResolver implements TypeResolver using a GraphQLSchema
type SchemaResolver struct {
	schema *GraphQLSchema
}

// NewSchemaResolver creates a new SchemaResolver that uses the provided schema
func NewSchemaResolver(schema *GraphQLSchema) *SchemaResolver {
	return &SchemaResolver{schema: schema}
}

// ResolveType resolves a type name to its definition
func (r *SchemaResolver) ResolveType(typeName string) (TypeDef, error) {
	return r.schema.NamedToTypeDef(typeName)
}

// ResolveFieldType resolves a field's result type
func (r *SchemaResolver) ResolveFieldType(field *Field) (TypeDef, error) {
	return field.ResolveObjectTypeDef(r.schema)
}

// ResolveArgumentType resolves an argument's input type
func (r *SchemaResolver) ResolveArgumentType(arg *Argument) (TypeDef, error) {
	return arg.ResolveObjectTypeDef(r.schema)
}

// ResolveDirective resolves a directive name to its definition
func (r *SchemaResolver) ResolveDirective(directiveName string) (*DirectiveDef, error) {
	directive, ok := r.schema.Directive[directiveName]
	if !ok {
		return nil, fmt.Errorf("directive %q not found in schema", directiveName)
	}
	return directive, nil
}

// ResolveUsages returns all usages of a given type name
func (r *SchemaResolver) ResolveUsages(typeName string) ([]*Usage, error) {
	usages, ok := r.schema.Usages[typeName]
	if !ok {
		return []*Usage{}, nil
	}
	return usages, nil
}

// ResolveQueryOrMutationField resolves a field from Query or Mutation
func (r *SchemaResolver) ResolveQueryOrMutationField(parentType, fieldName string) (*Field, error) {
	var field *Field
	var ok bool

	switch parentType {
	case "Query":
		field, ok = r.schema.Query[fieldName]
	case "Mutation":
		field, ok = r.schema.Mutation[fieldName]
	default:
		return nil, fmt.Errorf("%q is not Query or Mutation", parentType)
	}

	if !ok {
		return nil, fmt.Errorf("field %q not found in %s", fieldName, parentType)
	}
	return field, nil
}
