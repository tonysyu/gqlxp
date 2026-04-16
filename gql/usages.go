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

	schema.Walk(SchemaVisitor{
		VisitField: func(ctx VisitContext, name string, field *Field) {
			recordFieldUsage(usages, field, ctx.Kind, ctx.Kind, name)
			recordArgumentUsages(usages, field.Arguments(), ctx.Kind, ctx.Kind, name)
			recordArgumentDirectiveUsages(usages, field.Arguments(), ctx.Kind, ctx.Kind, name)
			recordDirectiveUsages(usages, field.Directives(), ctx.Kind, ctx.Kind, name)
		},
		VisitObject: func(ctx VisitContext, name string, obj *Object) {
			recordDirectiveUsages(usages, obj.Directives(), name, "Object", "")
			for _, ifaceName := range obj.Interfaces() {
				usages[ifaceName] = append(usages[ifaceName], &Usage{
					ParentType: name,
					ParentKind: "Object",
					Path:       name,
				})
			}
		},
		VisitObjectField: func(ctx VisitContext, field *Field) {
			recordFieldUsage(usages, field, ctx.ParentName, "Object", field.Name())
			recordArgumentUsages(usages, field.Arguments(), ctx.ParentName, "Object", field.Name())
			recordArgumentDirectiveUsages(usages, field.Arguments(), ctx.ParentName, "Object", field.Name())
			recordDirectiveUsages(usages, field.Directives(), ctx.ParentName, "Object", field.Name())
		},
		VisitInterface: func(ctx VisitContext, name string, iface *Interface) {
			recordDirectiveUsages(usages, iface.Directives(), name, "Interface", "")
			for _, parentIfaceName := range iface.Interfaces() {
				usages[parentIfaceName] = append(usages[parentIfaceName], &Usage{
					ParentType: name,
					ParentKind: "Interface",
					Path:       name,
				})
			}
		},
		VisitInterfaceField: func(ctx VisitContext, field *Field) {
			recordFieldUsage(usages, field, ctx.ParentName, "Interface", field.Name())
			recordArgumentUsages(usages, field.Arguments(), ctx.ParentName, "Interface", field.Name())
			recordArgumentDirectiveUsages(usages, field.Arguments(), ctx.ParentName, "Interface", field.Name())
			recordDirectiveUsages(usages, field.Directives(), ctx.ParentName, "Interface", field.Name())
		},
		VisitInput: func(ctx VisitContext, name string, input *InputObject) {
			recordDirectiveUsages(usages, input.Directives(), name, "Input", "")
		},
		VisitInputField: func(ctx VisitContext, field *Field) {
			recordFieldUsage(usages, field, ctx.ParentName, "Input", field.Name())
			recordDirectiveUsages(usages, field.Directives(), ctx.ParentName, "Input", field.Name())
		},
		VisitEnum: func(ctx VisitContext, name string, enum *Enum) {
			recordDirectiveUsages(usages, enum.Directives(), name, "Enum", "")
		},
		VisitEnumValue: func(ctx VisitContext, value *EnumValue) {
			recordDirectiveUsages(usages, value.Directives(), ctx.ParentName, "Enum", value.Name())
		},
		VisitScalar: func(ctx VisitContext, name string, scalar *Scalar) {
			recordDirectiveUsages(usages, scalar.Directives(), name, "Scalar", "")
		},
		VisitUnion: func(ctx VisitContext, name string, union *Union) {
			recordDirectiveUsages(usages, union.Directives(), name, "Union", "")
			for _, memberType := range union.Types() {
				usages[memberType] = append(usages[memberType], &Usage{
					ParentType: name,
					ParentKind: "Union",
					Path:       name,
				})
			}
		},
		VisitDirective: func(ctx VisitContext, name string, directive *DirectiveDef) {
			recordArgumentUsages(usages, directive.Arguments(), name, "Directive", "")
			recordArgumentDirectiveUsages(usages, directive.Arguments(), name, "Directive", "")
		},
	})

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
