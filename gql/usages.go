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

	// Index Query fields (return types and argument types)
	for fieldName, field := range schema.Query {
		recordFieldUsage(usages, field, "Query", "Query", fieldName)
		recordArgumentUsages(usages, field.Arguments(), "Query", "Query", fieldName)
	}

	// Index Mutation fields (return types and argument types)
	for fieldName, field := range schema.Mutation {
		recordFieldUsage(usages, field, "Mutation", "Mutation", fieldName)
		recordArgumentUsages(usages, field.Arguments(), "Mutation", "Mutation", fieldName)
	}

	// Index Object fields (return types and argument types)
	for objName, obj := range schema.Object {
		for _, field := range obj.Fields() {
			recordFieldUsage(usages, field, objName, "Object", field.Name())
			recordArgumentUsages(usages, field.Arguments(), objName, "Object", field.Name())
		}
		// Track interface implementations
		for _, ifaceName := range obj.Interfaces() {
			usage := &Usage{
				ParentType: objName,
				ParentKind: "Object",
				FieldName:  "",
				Path:       objName,
			}
			usages[ifaceName] = append(usages[ifaceName], usage)
		}
	}

	// Index Interface fields (return types and argument types)
	for ifName, iface := range schema.Interface {
		for _, field := range iface.Fields() {
			recordFieldUsage(usages, field, ifName, "Interface", field.Name())
			recordArgumentUsages(usages, field.Arguments(), ifName, "Interface", field.Name())
		}
		// Track interface-to-interface implementations
		for _, parentIfaceName := range iface.Interfaces() {
			usage := &Usage{
				ParentType: ifName,
				ParentKind: "Interface",
				FieldName:  "",
				Path:       ifName,
			}
			usages[parentIfaceName] = append(usages[parentIfaceName], usage)
		}
	}

	// Index InputObject field types
	for inputName, input := range schema.Input {
		for _, field := range input.Fields() {
			recordFieldUsage(usages, field, inputName, "Input", field.Name())
		}
	}

	// Index Union member types
	for unionName, union := range schema.Union {
		for _, memberType := range union.Types() {
			usage := &Usage{
				ParentType: unionName,
				ParentKind: "Union",
				FieldName:  "",
				Path:       unionName,
			}
			usages[memberType] = append(usages[memberType], usage)
		}
	}

	// Index Directive argument types
	for directiveName, directive := range schema.Directive {
		recordArgumentUsages(usages, directive.Arguments(), directiveName, "Directive", "")
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

// recordArgumentUsages records usages for all argument types in a field or directive
func recordArgumentUsages(usages map[string][]*Usage, arguments []*Argument, parentType, parentKind, fieldName string) {
	for _, arg := range arguments {
		typeName := arg.ObjectTypeName()
		if typeName == "" {
			continue
		}

		path := parentType
		if fieldName != "" {
			path += "." + fieldName
		}
		path += "(" + arg.Name() + ": " + typeName + ")"

		usage := &Usage{
			ParentType: parentType,
			ParentKind: parentKind,
			FieldName:  fieldName,
			Path:       path,
		}

		usages[typeName] = append(usages[typeName], usage)
	}
}
