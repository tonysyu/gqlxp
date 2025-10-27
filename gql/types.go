package gql

import "github.com/vektah/gqlparser/v2/ast"

// Interface for GQL Type Definitions
type TypeDef interface {
	Name() string
	Description() string
}

var _ TypeDef = (*Field)(nil)
var _ TypeDef = (*Object)(nil)
var _ TypeDef = (*InputObject)(nil)
var _ TypeDef = (*Enum)(nil)
var _ TypeDef = (*Scalar)(nil)
var _ TypeDef = (*Interface)(nil)
var _ TypeDef = (*Union)(nil)
var _ TypeDef = (*Directive)(nil)

// gqlType interface for types that have a Name field (both ast and wrapped types)
type gqlType interface {
	*Field | *Object | *InputObject |
		*Enum | *Scalar | *Interface |
		*Union | *Directive
}

// Field represents GraphQL query fields, mutation fields, or object fields.
// See https://spec.graphql.org/October2021/#sec-Objects
type Field struct {
	name        string
	description string
	fieldType   *ast.Type
	arguments   ast.ArgumentDefinitionList
	// Keep reference to underlying ast for operations within gql package
	astField *ast.FieldDefinition
}

// newField creates a wrapper for ast.FieldDefinition
func newField(field *ast.FieldDefinition) *Field {
	if field == nil {
		return nil
	}
	return &Field{
		name:        field.Name,
		description: field.Description,
		fieldType:   field.Type,
		arguments:   field.Arguments,
		astField:    field,
	}
}

func (f *Field) Name() string        { return f.name }
func (f *Field) Description() string { return f.description }

// TypeString returns the string representation of the field's type
func (f *Field) TypeString() string {
	return getTypeString(f.fieldType)
}

// TypeName returns the named type (unwrapping List and NonNull)
func (f *Field) TypeName() string {
	return getNamedTypeName(f.fieldType)
}

// Signature returns the full field signature including arguments
func (f *Field) Signature() string {
	return getFieldString(f.astField)
}

// Arguments returns the field's arguments as wrapped Arguments
func (f *Field) Arguments() []*Argument {
	return wrapArguments(f.arguments)
}

// ResolveResultType returns the TypeDef for this field's return type.
// Returns an error if the type is a built-in scalar (String, Int, Boolean, etc.)
func (f *Field) ResolveResultType(schema *GraphQLSchema) (TypeDef, error) {
	typeName := getNamedTypeName(f.fieldType)
	return schema.NamedToTypeDef(typeName)
}

// wrapFields converts a slice of ast.FieldDefinition to wrapped Field types.
// This is useful for converting fields from Object, Interface, etc.
func wrapFields(astFields ast.FieldList) []*Field {
	fields := make([]*Field, 0, len(astFields))
	for _, field := range astFields {
		fields = append(fields, newField(field))
	}
	return fields
}

// Argument represents arguments to a GraphQL Object field.
// See https://spec.graphql.org/October2021/#sec-Field-Arguments
type Argument struct {
	name        string
	description string
	inputType   *ast.Type
	// Keep reference to underlying ast for operations within gql package
	astArg *ast.ArgumentDefinition
}

// newArgument creates a wrapper for ast.ArgumentDefinition
func newArgument(arg *ast.ArgumentDefinition) *Argument {
	if arg == nil {
		return nil
	}
	return &Argument{
		name:        arg.Name,
		description: arg.Description,
		inputType:   arg.Type,
		astArg:      arg,
	}
}

func (a *Argument) Name() string        { return a.name }
func (a *Argument) Description() string { return a.description }

// TypeString returns the string representation of the argument's type
func (a *Argument) TypeString() string {
	return getTypeString(a.inputType)
}

// Signature returns the full argument signature (name: Type = defaultValue)
func (a *Argument) Signature() string {
	return getArgumentString(a.astArg)
}

// wrapArguments converts a slice of ast.ArgumentDefinition to wrapped Argument types
func wrapArguments(args ast.ArgumentDefinitionList) []*Argument {
	arguments := make([]*Argument, 0, len(args))
	for _, arg := range args {
		arguments = append(arguments, newArgument(arg))
	}
	return arguments
}

// InputField represents fields on InputObjects
// TODO: Evaluate whether InputField can be removed in favor of Field
type InputField struct {
	name        string
	description string
	inputType   *ast.Type
	// Keep reference to underlying ast for operations within gql package
	astField *ast.FieldDefinition
}

// newInputField creates a wrapper for input object field (ast.FieldDefinition)
func newInputField(field *ast.FieldDefinition) *InputField {
	if field == nil {
		return nil
	}
	return &InputField{
		name:        field.Name,
		description: field.Description,
		inputType:   field.Type,
		astField:    field,
	}
}

func (f *InputField) Name() string        { return f.name }
func (f *InputField) Description() string { return f.description }

// TypeString returns the string representation of the input field's type
func (f *InputField) TypeString() string {
	return getTypeString(f.inputType)
}

// Signature returns the full input field signature (name: Type = defaultValue)
func (f *InputField) Signature() string {
	return getInputFieldString(f.astField)
}

// wrapInputFields converts a slice of ast.FieldDefinition (for input objects) to wrapped InputField types
func wrapInputFields(fields ast.FieldList) []*InputField {
	inputFields := make([]*InputField, 0, len(fields))
	for _, field := range fields {
		inputFields = append(inputFields, newInputField(field))
	}
	return inputFields
}

// Object represents GraphQL objects.
// See https://spec.graphql.org/October2021/#sec-Objects
type Object struct {
	name        string
	description string
	interfaces  []string
	fields      []*Field
	// Keep reference to underlying ast for operations within gql package
	astDef *ast.Definition
}

// newObject creates a wrapper for ast.Definition (Object type)
func newObject(def *ast.Definition) *Object {
	if def == nil {
		return nil
	}
	return &Object{
		name:        def.Name,
		description: def.Description,
		interfaces:  def.Interfaces,
		fields:      wrapFields(def.Fields),
		astDef:      def,
	}
}

func (o *Object) Name() string        { return o.name }
func (o *Object) Description() string { return o.description }

// Interfaces returns the interfaces this object implements
func (o *Object) Interfaces() []string {
	return o.interfaces
}

// Fields returns the object's fields
func (o *Object) Fields() []*Field {
	return o.fields
}

// InputObject represents GraphQL input objects.
// See https://spec.graphql.org/October2021/#sec-Input-Objects
type InputObject struct {
	name        string
	description string
	fields      []*InputField
	// Keep reference to underlying ast for operations within gql package
	astDef *ast.Definition
}

// newInputObject creates a wrapper for ast.Definition (InputObject type)
func newInputObject(def *ast.Definition) *InputObject {
	if def == nil {
		return nil
	}
	return &InputObject{
		name:        def.Name,
		description: def.Description,
		fields:      wrapInputFields(def.Fields),
		astDef:      def,
	}
}

func (i *InputObject) Name() string        { return i.name }
func (i *InputObject) Description() string { return i.description }

// Fields returns the input object's fields
func (i *InputObject) Fields() []*InputField {
	return i.fields
}

// Enum represents GraphQL enums.
// See https://spec.graphql.org/October2021/#sec-Enums
type Enum struct {
	name        string
	description string
	values      []*EnumValue
	// Keep reference to underlying ast for operations within gql package
	astDef *ast.Definition
}

// newEnum creates a wrapper for ast.Definition (Enum type)
func newEnum(def *ast.Definition) *Enum {
	if def == nil {
		return nil
	}
	return &Enum{
		name:        def.Name,
		description: def.Description,
		values:      wrapEnumValues(def.EnumValues),
		astDef:      def,
	}
}

func (e *Enum) Name() string        { return e.name }
func (e *Enum) Description() string { return e.description }

// Values returns the enum values
func (e *Enum) Values() []*EnumValue {
	return e.values
}

// Scalar represents GraphQL scalar types.
// See https://spec.graphql.org/October2021/#sec-Scalars
type Scalar struct {
	name        string
	description string
	// Keep reference to underlying ast for operations within gql package
	astDef *ast.Definition
}

// newScalar creates a wrapper for ast.Definition (Scalar type)
func newScalar(def *ast.Definition) *Scalar {
	if def == nil {
		return nil
	}
	return &Scalar{
		name:        def.Name,
		description: def.Description,
		astDef:      def,
	}
}

func (s *Scalar) Name() string        { return s.name }
func (s *Scalar) Description() string { return s.description }

// Interface represents GraphQL interface types.
// See https://spec.graphql.org/October2021/#sec-Interfaces
type Interface struct {
	name        string
	description string
	fields      []*Field
	// Keep reference to underlying ast for operations within gql package
	astDef *ast.Definition
}

// newInterface creates a wrapper for ast.Definition (Interface type)
func newInterface(def *ast.Definition) *Interface {
	if def == nil {
		return nil
	}
	return &Interface{
		name:        def.Name,
		description: def.Description,
		fields:      wrapFields(def.Fields),
		astDef:      def,
	}
}

func (i *Interface) Name() string        { return i.name }
func (i *Interface) Description() string { return i.description }

// Fields returns the interface's fields
func (i *Interface) Fields() []*Field {
	return i.fields
}

// Union represents GraphQL union types.
// See https://spec.graphql.org/October2021/#sec-Unions
type Union struct {
	name        string
	description string
	types       []string
	// Keep reference to underlying ast for operations within gql package
	astDef *ast.Definition
}

// newUnion creates a wrapper for ast.Definition (Union type)
func newUnion(def *ast.Definition) *Union {
	if def == nil {
		return nil
	}
	return &Union{
		name:        def.Name,
		description: def.Description,
		types:       def.Types,
		astDef:      def,
	}
}

func (u *Union) Name() string        { return u.name }
func (u *Union) Description() string { return u.description }

// Types returns the union's member types as type names
func (u *Union) Types() []string {
	return u.types
}

// Directive represents GraphQL directive types.
// See https://spec.graphql.org/October2021/#sec-Type-System.Directives
type Directive struct {
	name        string
	description string
	// Keep reference to underlying ast for operations within gql package
	astDirective *ast.DirectiveDefinition
}

// newDirective creates a wrapper for ast.DirectiveDefinition
func newDirective(directive *ast.DirectiveDefinition) *Directive {
	if directive == nil {
		return nil
	}
	return &Directive{
		name:         directive.Name,
		description:  directive.Description,
		astDirective: directive,
	}
}

func (d *Directive) Name() string        { return d.name }
func (d *Directive) Description() string { return d.description }

// EnumValue represents values on GraphQL enum types.
// See https://spec.graphql.org/October2021/#sec-Enums
type EnumValue struct {
	name        string
	description string
	// Keep reference to underlying ast for operations within gql package
	astEnumValue *ast.EnumValueDefinition
}

// newEnumValue creates a wrapper for ast.EnumValueDefinition
func newEnumValue(enumValue *ast.EnumValueDefinition) *EnumValue {
	if enumValue == nil {
		return nil
	}
	return &EnumValue{
		name:         enumValue.Name,
		description:  enumValue.Description,
		astEnumValue: enumValue,
	}
}

func (e *EnumValue) Name() string        { return e.name }
func (e *EnumValue) Description() string { return e.description }

// wrapEnumValues converts a slice of ast.EnumValueDefinition to wrapped types
func wrapEnumValues(astEnumValues ast.EnumValueList) []*EnumValue {
	enumValues := make([]*EnumValue, 0, len(astEnumValues))
	for _, ev := range astEnumValues {
		enumValues = append(enumValues, newEnumValue(ev))
	}
	return enumValues
}
