package gql

// TypeResolver provides methods for resolving GraphQL type definitions
type TypeResolver interface {
	// ResolveType resolves a type name to its definition
	ResolveType(typeName string) (TypeDef, error)

	// ResolveFieldType resolves a field's result type
	ResolveFieldType(field *Field) (TypeDef, error)

	// ResolveArgumentType resolves an argument's input type
	ResolveArgumentType(arg *Argument) (TypeDef, error)
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
