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
		recordArgumentDirectiveUsages(usages, field.Arguments(), "Query", "Query", fieldName)
		recordDirectiveUsages(usages, field.Directives(), "Query", "Query", fieldName)
	}

	// Index Mutation fields (return types and argument types)
	for fieldName, field := range schema.Mutation {
		recordFieldUsage(usages, field, "Mutation", "Mutation", fieldName)
		recordArgumentUsages(usages, field.Arguments(), "Mutation", "Mutation", fieldName)
		recordArgumentDirectiveUsages(usages, field.Arguments(), "Mutation", "Mutation", fieldName)
		recordDirectiveUsages(usages, field.Directives(), "Mutation", "Mutation", fieldName)
	}

	// Index Object fields (return types and argument types)
	for objName, obj := range schema.Object {
		// Track directives on the object type itself
		recordDirectiveUsages(usages, obj.Directives(), objName, "Object", "")

		for _, field := range obj.Fields() {
			recordFieldUsage(usages, field, objName, "Object", field.Name())
			recordArgumentUsages(usages, field.Arguments(), objName, "Object", field.Name())
			recordArgumentDirectiveUsages(usages, field.Arguments(), objName, "Object", field.Name())
			recordDirectiveUsages(usages, field.Directives(), objName, "Object", field.Name())
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
		// Track directives on the interface type itself
		recordDirectiveUsages(usages, iface.Directives(), ifName, "Interface", "")

		for _, field := range iface.Fields() {
			recordFieldUsage(usages, field, ifName, "Interface", field.Name())
			recordArgumentUsages(usages, field.Arguments(), ifName, "Interface", field.Name())
			recordArgumentDirectiveUsages(usages, field.Arguments(), ifName, "Interface", field.Name())
			recordDirectiveUsages(usages, field.Directives(), ifName, "Interface", field.Name())
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
		// Track directives on the input type itself
		recordDirectiveUsages(usages, input.Directives(), inputName, "Input", "")

		for _, field := range input.Fields() {
			recordFieldUsage(usages, field, inputName, "Input", field.Name())
			recordDirectiveUsages(usages, field.Directives(), inputName, "Input", field.Name())
		}
	}

	// Index Union member types
	for unionName, union := range schema.Union {
		// Track directives on the union type itself
		recordDirectiveUsages(usages, union.Directives(), unionName, "Union", "")

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

	// Index Enum types and values
	for enumName, enum := range schema.Enum {
		// Track directives on the enum type itself
		recordDirectiveUsages(usages, enum.Directives(), enumName, "Enum", "")

		// Track directives on enum values
		for _, enumValue := range enum.Values() {
			recordDirectiveUsages(usages, enumValue.Directives(), enumName, "Enum", enumValue.Name())
		}
	}

	// Index Scalar types
	for scalarName, scalar := range schema.Scalar {
		// Track directives on the scalar type itself
		recordDirectiveUsages(usages, scalar.Directives(), scalarName, "Scalar", "")
	}

	// Index Directive argument types
	for directiveName, directive := range schema.Directive {
		recordArgumentUsages(usages, directive.Arguments(), directiveName, "Directive", "")
		recordArgumentDirectiveUsages(usages, directive.Arguments(), directiveName, "Directive", "")
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

// recordDirectiveUsages records usages for all directives applied to a type or field
func recordDirectiveUsages(usages map[string][]*Usage, directives []*AppliedDirective, parentType, parentKind, fieldName string) {
	for _, directive := range directives {
		directiveName := directive.Name()
		if directiveName == "" {
			continue
		}

		path := parentType
		if fieldName != "" {
			path += "." + fieldName
		}

		usage := &Usage{
			ParentType: parentType,
			ParentKind: parentKind,
			FieldName:  fieldName,
			Path:       path,
		}

		usages[directiveName] = append(usages[directiveName], usage)
	}
}

// recordArgumentDirectiveUsages records directive usages for all arguments
func recordArgumentDirectiveUsages(usages map[string][]*Usage, arguments []*Argument, parentType, parentKind, fieldName string) {
	for _, arg := range arguments {
		argDirectives := arg.Directives()
		if len(argDirectives) == 0 {
			continue
		}

		// Build the path for the argument
		path := parentType
		if fieldName != "" {
			path += "." + fieldName
		}
		path += "(" + arg.Name() + ")"

		// Record directive usages on this argument
		for _, directive := range argDirectives {
			directiveName := directive.Name()
			if directiveName == "" {
				continue
			}

			usage := &Usage{
				ParentType: parentType,
				ParentKind: parentKind,
				FieldName:  fieldName,
				Path:       path,
			}

			usages[directiveName] = append(usages[directiveName], usage)
		}
	}
}
