package gql

// Usage represents a single usage of a type within the schema
type Usage struct {
	ParentType string // Type containing this usage (e.g., "User", "Query")
	ParentKind string // "Object", "Interface", "Query", "Mutation"
	FieldName  string // Name of the field using this type
	Path       string // Full path like "Query.user" or "User.posts"
}

// buildUsageIndex builds the usage index for all types in the schema
func buildUsageIndex(schema *GraphQLSchema) {
	usages := make(map[string][]*Usage)

	// Index Query fields
	for fieldName, field := range schema.Query {
		recordFieldUsage(usages, field, "Query", "Query", fieldName)
	}

	// Index Mutation fields
	for fieldName, field := range schema.Mutation {
		recordFieldUsage(usages, field, "Mutation", "Mutation", fieldName)
	}

	// Index Object fields
	for objName, obj := range schema.Object {
		for _, field := range obj.Fields() {
			recordFieldUsage(usages, field, objName, "Object", field.Name())
		}
	}

	// Index Interface fields
	for ifName, iface := range schema.Interface {
		for _, field := range iface.Fields() {
			recordFieldUsage(usages, field, ifName, "Interface", field.Name())
		}
	}

	schema.Usages = usages
}

// recordFieldUsage extracts the type name from a field and records the usage
func recordFieldUsage(usages map[string][]*Usage, field *Field, parentType, parentKind, fieldName string) {
	typeName := field.ObjectTypeName()
	if typeName == "" {
		return
	}

	usage := &Usage{
		ParentType: parentType,
		ParentKind: parentKind,
		FieldName:  fieldName,
		Path:       parentType + "." + fieldName,
	}

	usages[typeName] = append(usages[typeName], usage)
}
