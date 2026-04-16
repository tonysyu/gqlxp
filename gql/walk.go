package gql

// VisitContext holds traversal position passed to SchemaVisitor callbacks.
type VisitContext struct {
	// Kind is the top-level collection: "Query", "Mutation", "Object",
	// "Interface", "Input", "Enum", "Scalar", "Union", or "Directive".
	Kind string

	// ParentName is the enclosing type name when visiting sub-elements
	// (fields, enum values). Empty when visiting a top-level type.
	ParentName string
}

// SchemaVisitor holds per-kind callbacks for schema traversal.
// All fields are optional; nil callbacks are silently skipped.
type SchemaVisitor struct {
	// Top-level type callbacks
	VisitField     func(ctx VisitContext, name string, f *Field)
	VisitObject    func(ctx VisitContext, name string, obj *Object)
	VisitInterface func(ctx VisitContext, name string, iface *Interface)
	VisitInput     func(ctx VisitContext, name string, input *InputObject)
	VisitEnum      func(ctx VisitContext, name string, enum *Enum)
	VisitScalar    func(ctx VisitContext, name string, scalar *Scalar)
	VisitUnion     func(ctx VisitContext, name string, union *Union)
	VisitDirective func(ctx VisitContext, name string, directive *DirectiveDef)

	// Sub-element callbacks — ctx.ParentName is always set when these fire
	VisitObjectField    func(ctx VisitContext, field *Field)
	VisitInterfaceField func(ctx VisitContext, field *Field)
	VisitInputField     func(ctx VisitContext, field *Field)
	VisitEnumValue      func(ctx VisitContext, value *EnumValue)
}

// Walk traverses every node in the schema, invoking matching callbacks on v.
// Types are visited alphabetically by name within each kind.
// Kind order: Query, Mutation, Object, Interface, Input, Enum, Scalar, Union, Directive.
func (s *GraphQLSchema) Walk(v SchemaVisitor) {
	// Query fields
	if v.VisitField != nil {
		ctx := VisitContext{Kind: "Query"}
		for _, field := range CollectAndSortMapValues(s.Query) {
			v.VisitField(ctx, field.Name(), field)
		}
	}

	// Mutation fields
	if v.VisitField != nil {
		ctx := VisitContext{Kind: "Mutation"}
		for _, field := range CollectAndSortMapValues(s.Mutation) {
			v.VisitField(ctx, field.Name(), field)
		}
	}

	// Objects (always iterate; check callbacks per call)
	for _, obj := range CollectAndSortMapValues(s.Object) {
		name := obj.Name()
		if v.VisitObject != nil {
			v.VisitObject(VisitContext{Kind: "Object"}, name, obj)
		}
		if v.VisitObjectField != nil {
			ctx := VisitContext{Kind: "Object", ParentName: name}
			for _, field := range obj.Fields() {
				v.VisitObjectField(ctx, field)
			}
		}
	}

	// Interfaces
	for _, iface := range CollectAndSortMapValues(s.Interface) {
		name := iface.Name()
		if v.VisitInterface != nil {
			v.VisitInterface(VisitContext{Kind: "Interface"}, name, iface)
		}
		if v.VisitInterfaceField != nil {
			ctx := VisitContext{Kind: "Interface", ParentName: name}
			for _, field := range iface.Fields() {
				v.VisitInterfaceField(ctx, field)
			}
		}
	}

	// Input objects
	for _, input := range CollectAndSortMapValues(s.Input) {
		name := input.Name()
		if v.VisitInput != nil {
			v.VisitInput(VisitContext{Kind: "Input"}, name, input)
		}
		if v.VisitInputField != nil {
			ctx := VisitContext{Kind: "Input", ParentName: name}
			for _, field := range input.Fields() {
				v.VisitInputField(ctx, field)
			}
		}
	}

	// Enums
	for _, enum := range CollectAndSortMapValues(s.Enum) {
		name := enum.Name()
		if v.VisitEnum != nil {
			v.VisitEnum(VisitContext{Kind: "Enum"}, name, enum)
		}
		if v.VisitEnumValue != nil {
			ctx := VisitContext{Kind: "Enum", ParentName: name}
			for _, value := range enum.Values() {
				v.VisitEnumValue(ctx, value)
			}
		}
	}

	// Scalars
	if v.VisitScalar != nil {
		for _, scalar := range CollectAndSortMapValues(s.Scalar) {
			v.VisitScalar(VisitContext{Kind: "Scalar"}, scalar.Name(), scalar)
		}
	}

	// Unions
	if v.VisitUnion != nil {
		for _, union := range CollectAndSortMapValues(s.Union) {
			v.VisitUnion(VisitContext{Kind: "Union"}, union.Name(), union)
		}
	}

	// Directives
	if v.VisitDirective != nil {
		for _, directive := range CollectAndSortMapValues(s.Directive) {
			v.VisitDirective(VisitContext{Kind: "Directive"}, directive.Name(), directive)
		}
	}
}
