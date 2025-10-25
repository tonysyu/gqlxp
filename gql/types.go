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

// Field wraps ast.Field to avoid exposing gqlparser types outside gql package.
// This can represent query fields, mutation fields, or object fields.
type Field struct {
	name        string
	description string
	fieldType   *ast.Type
	arguments   ast.ArgumentDefinitionList
	// Keep reference to underlying ast for operations within gql package
	astField *ast.FieldDefinition
}

// NewField creates a wrapper for ast.FieldDefinition
func NewField(field *ast.FieldDefinition) *Field {
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

// Name returns the field name
func (f *Field) Name() string {
	return f.name
}

// Description returns the field description
func (f *Field) Description() string {
	return f.description
}

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

// Arguments returns the field's arguments as wrapped InputValues
func (f *Field) Arguments() []*InputValue {
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
		fields = append(fields, NewField(field))
	}
	return fields
}

// InputValue wraps ast.ArgumentDefinition to avoid exposing gqlparser types outside gql package.
// This represents input arguments or input object fields.
type InputValue struct {
	name        string
	description string
	inputType   *ast.Type
	// Keep reference to underlying ast for operations within gql package
	astArg   *ast.ArgumentDefinition
	astField *ast.FieldDefinition // For input object fields
}

// NewInputValue creates a wrapper for ast.ArgumentDefinition
func NewInputValue(arg *ast.ArgumentDefinition) *InputValue {
	if arg == nil {
		return nil
	}
	return &InputValue{
		name:        arg.Name,
		description: arg.Description,
		inputType:   arg.Type,
		astArg:      arg,
	}
}

// NewInputValueFromField creates a wrapper for input object field (ast.FieldDefinition)
func NewInputValueFromField(field *ast.FieldDefinition) *InputValue {
	if field == nil {
		return nil
	}
	return &InputValue{
		name:        field.Name,
		description: field.Description,
		inputType:   field.Type,
		astField:    field,
	}
}

// Name returns the input value name
func (i *InputValue) Name() string {
	return i.name
}

// Description returns the input value description
func (i *InputValue) Description() string {
	return i.description
}

// TypeString returns the string representation of the input value's type
func (i *InputValue) TypeString() string {
	return getTypeString(i.inputType)
}

// Signature returns the full input value signature (name: Type = defaultValue)
func (i *InputValue) Signature() string {
	if i.astArg != nil {
		return getArgumentString(i.astArg)
	}
	if i.astField != nil {
		return getInputFieldString(i.astField)
	}
	return i.name + ": " + getTypeString(i.inputType)
}

// wrapArguments converts a slice of ast.ArgumentDefinition to wrapped InputValue types
func wrapArguments(args ast.ArgumentDefinitionList) []*InputValue {
	inputValues := make([]*InputValue, 0, len(args))
	for _, arg := range args {
		inputValues = append(inputValues, NewInputValue(arg))
	}
	return inputValues
}

// wrapInputFields converts a slice of ast.FieldDefinition (for input objects) to wrapped InputValue types
func wrapInputFields(fields ast.FieldList) []*InputValue {
	inputValues := make([]*InputValue, 0, len(fields))
	for _, field := range fields {
		inputValues = append(inputValues, NewInputValueFromField(field))
	}
	return inputValues
}

// Object wraps ast.Definition (for Object types) to avoid exposing gqlparser types outside gql package.
type Object struct {
	name        string
	description string
	interfaces  []string
	fields      []*Field
	// Keep reference to underlying ast for operations within gql package
	astDef *ast.Definition
}

// NewObject creates a wrapper for ast.Definition (Object type)
func NewObject(def *ast.Definition) *Object {
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

// Name returns the object name
func (o *Object) Name() string {
	return o.name
}

// Description returns the object description
func (o *Object) Description() string {
	return o.description
}

// Interfaces returns the interfaces this object implements
func (o *Object) Interfaces() []string {
	return o.interfaces
}

// Fields returns the object's fields
func (o *Object) Fields() []*Field {
	return o.fields
}

// InputObject wraps ast.Definition (for InputObject types) to avoid exposing gqlparser types outside gql package.
type InputObject struct {
	name        string
	description string
	fields      []*InputValue
	// Keep reference to underlying ast for operations within gql package
	astDef *ast.Definition
}

// NewInputObject creates a wrapper for ast.Definition (InputObject type)
func NewInputObject(def *ast.Definition) *InputObject {
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

// Name returns the input object name
func (i *InputObject) Name() string {
	return i.name
}

// Description returns the input object description
func (i *InputObject) Description() string {
	return i.description
}

// Fields returns the input object's fields
func (i *InputObject) Fields() []*InputValue {
	return i.fields
}

// Enum wraps ast.Definition (for Enum types) to avoid exposing gqlparser types outside gql package.
type Enum struct {
	name        string
	description string
	values      []*EnumValue
	// Keep reference to underlying ast for operations within gql package
	astDef *ast.Definition
}

// NewEnum creates a wrapper for ast.Definition (Enum type)
func NewEnum(def *ast.Definition) *Enum {
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

// Name returns the enum name
func (e *Enum) Name() string {
	return e.name
}

// Description returns the enum description
func (e *Enum) Description() string {
	return e.description
}

// Values returns the enum values
func (e *Enum) Values() []*EnumValue {
	return e.values
}

// Scalar wraps ast.Definition (for Scalar types) to avoid exposing gqlparser types outside gql package.
type Scalar struct {
	name        string
	description string
	// Keep reference to underlying ast for operations within gql package
	astDef *ast.Definition
}

// NewScalar creates a wrapper for ast.Definition (Scalar type)
func NewScalar(def *ast.Definition) *Scalar {
	if def == nil {
		return nil
	}
	return &Scalar{
		name:        def.Name,
		description: def.Description,
		astDef:      def,
	}
}

// Name returns the scalar name
func (s *Scalar) Name() string {
	return s.name
}

// Description returns the scalar description
func (s *Scalar) Description() string {
	return s.description
}

// Interface wraps ast.Definition (for Interface types) to avoid exposing gqlparser types outside gql package.
type Interface struct {
	name        string
	description string
	fields      []*Field
	// Keep reference to underlying ast for operations within gql package
	astDef *ast.Definition
}

// NewInterface creates a wrapper for ast.Definition (Interface type)
func NewInterface(def *ast.Definition) *Interface {
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

// Name returns the interface name
func (i *Interface) Name() string {
	return i.name
}

// Description returns the interface description
func (i *Interface) Description() string {
	return i.description
}

// Fields returns the interface's fields
func (i *Interface) Fields() []*Field {
	return i.fields
}

// Union wraps ast.Definition (for Union types) to avoid exposing gqlparser types outside gql package.
type Union struct {
	name        string
	description string
	types       []string
	// Keep reference to underlying ast for operations within gql package
	astDef *ast.Definition
}

// NewUnion creates a wrapper for ast.Definition (Union type)
func NewUnion(def *ast.Definition) *Union {
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

// Name returns the union name
func (u *Union) Name() string {
	return u.name
}

// Description returns the union description
func (u *Union) Description() string {
	return u.description
}

// Types returns the union's member types as type names
func (u *Union) Types() []string {
	return u.types
}

// Directive wraps ast.DirectiveDefinition to avoid exposing gqlparser types outside gql package.
type Directive struct {
	name        string
	description string
	// Keep reference to underlying ast for operations within gql package
	astDirective *ast.DirectiveDefinition
}

// NewDirective creates a wrapper for ast.DirectiveDefinition
func NewDirective(directive *ast.DirectiveDefinition) *Directive {
	if directive == nil {
		return nil
	}
	return &Directive{
		name:         directive.Name,
		description:  directive.Description,
		astDirective: directive,
	}
}

// Name returns the directive name
func (d *Directive) Name() string {
	return d.name
}

// Description returns the directive description
func (d *Directive) Description() string {
	return d.description
}

// EnumValue wraps ast.EnumValueDefinition to avoid exposing gqlparser types outside gql package.
type EnumValue struct {
	name        string
	description string
	// Keep reference to underlying ast for operations within gql package
	astEnumValue *ast.EnumValueDefinition
}

// NewEnumValue creates a wrapper for ast.EnumValueDefinition
func NewEnumValue(enumValue *ast.EnumValueDefinition) *EnumValue {
	if enumValue == nil {
		return nil
	}
	return &EnumValue{
		name:         enumValue.Name,
		description:  enumValue.Description,
		astEnumValue: enumValue,
	}
}

// Name returns the enum value name
func (e *EnumValue) Name() string {
	return e.name
}

// Description returns the enum value description
func (e *EnumValue) Description() string {
	return e.description
}

// wrapEnumValues converts a slice of ast.EnumValueDefinition to wrapped types
func wrapEnumValues(astEnumValues ast.EnumValueList) []*EnumValue {
	enumValues := make([]*EnumValue, 0, len(astEnumValues))
	for _, ev := range astEnumValues {
		enumValues = append(enumValues, NewEnumValue(ev))
	}
	return enumValues
}
