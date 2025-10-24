package gql

import "github.com/graphql-go/graphql/language/ast"

// Field wraps ast.Field to avoid exposing graphql-go types outside gql package.
// This can represent query fields, mutation fields, or object fields.
type Field struct {
	name        string
	description string
	fieldType   ast.Type
	arguments   []*ast.InputValueDefinition
	// Keep reference to underlying ast for operations within gql package
	astField *ast.FieldDefinition
}

// NewField creates a wrapper for ast.FieldDefinition
func NewField(field *ast.FieldDefinition) *Field {
	if field == nil {
		return nil
	}
	return &Field{
		name:        field.Name.Value,
		description: getStringValue(field.GetDescription()),
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
	named := getNamedFromType(f.fieldType)
	if named == nil {
		return ""
	}
	return named.Name.Value
}

// Signature returns the full field signature including arguments
func (f *Field) Signature() string {
	return getFieldString(f.astField)
}

// Arguments returns the field's arguments as wrapped InputValues
func (f *Field) Arguments() []*InputValue {
	return wrapInputValues(f.arguments)
}

// ResolveResultType returns the NamedTypeDef for this field's return type.
// Returns an error if the type is a built-in scalar (String, Int, Boolean, etc.)
func (f *Field) ResolveResultType(schema *GraphQLSchema) (NamedTypeDef, error) {
	named := getNamedFromType(f.fieldType)
	return schema.NamedToTypeDef(named)
}

// wrapFields converts a slice of ast.FieldDefinition to wrapped Field types.
// This is useful for converting fields from Object, Interface, etc.
func wrapFields(astFields []*ast.FieldDefinition) []*Field {
	fields := make([]*Field, 0, len(astFields))
	for _, field := range astFields {
		fields = append(fields, NewField(field))
	}
	return fields
}

// InputValue wraps ast.InputValue to avoid exposing graphql-go types outside gql package.
// This represents input arguments or input object fields.
type InputValue struct {
	name        string
	description string
	inputType   ast.Type
	// Keep reference to underlying ast for operations within gql package
	astInputValue *ast.InputValueDefinition
}

// NewInputValue creates a wrapper for ast.InputValueDefinition
func NewInputValue(inputValue *ast.InputValueDefinition) *InputValue {
	if inputValue == nil {
		return nil
	}
	return &InputValue{
		name:          inputValue.Name.Value,
		description:   getStringValue(inputValue.Description),
		inputType:     inputValue.Type,
		astInputValue: inputValue,
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
	return getInputValueString(i.astInputValue)
}

// wrapInputValues converts a slice of ast.InputValueDefinition to wrapped types
func wrapInputValues(astInputValues []*ast.InputValueDefinition) []*InputValue {
	inputValues := make([]*InputValue, 0, len(astInputValues))
	for _, iv := range astInputValues {
		inputValues = append(inputValues, NewInputValue(iv))
	}
	return inputValues
}

// Object wraps ast.Object to avoid exposing graphql-go types outside gql package.
type Object struct {
	name        string
	description string
	interfaces  []*Named
	fields      []*Field
	// Keep reference to underlying ast for operations within gql package
	astObject *ast.ObjectDefinition
}

// NewObject creates a wrapper for ast.ObjectDefinition
func NewObject(obj *ast.ObjectDefinition) *Object {
	if obj == nil {
		return nil
	}
	return &Object{
		name:        obj.Name.Value,
		description: getStringValue(obj.GetDescription()),
		interfaces:  wrapNamed(obj.Interfaces),
		fields:      wrapFields(obj.Fields),
		astObject:   obj,
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
func (o *Object) Interfaces() []*Named {
	return o.interfaces
}

// Fields returns the object's fields
func (o *Object) Fields() []*Field {
	return o.fields
}

// GetName returns the underlying AST name (for NamedTypeDef interface)
func (o *Object) GetName() *ast.Name {
	return o.astObject.GetName()
}

// GetDescription returns the underlying AST description (for NamedTypeDef interface)
func (o *Object) GetDescription() *ast.StringValue {
	return o.astObject.GetDescription()
}

// InputObject wraps ast.InputObject to avoid exposing graphql-go types outside gql package.
type InputObject struct {
	name        string
	description string
	fields      []*InputValue
	// Keep reference to underlying ast for operations within gql package
	astInputObject *ast.InputObjectDefinition
}

// NewInputObject creates a wrapper for ast.InputObjectDefinition
func NewInputObject(input *ast.InputObjectDefinition) *InputObject {
	if input == nil {
		return nil
	}
	return &InputObject{
		name:           input.Name.Value,
		description:    getStringValue(input.GetDescription()),
		fields:         wrapInputValues(input.Fields),
		astInputObject: input,
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

// GetName returns the underlying AST name (for NamedTypeDef interface)
func (i *InputObject) GetName() *ast.Name {
	return i.astInputObject.GetName()
}

// GetDescription returns the underlying AST description (for NamedTypeDef interface)
func (i *InputObject) GetDescription() *ast.StringValue {
	return i.astInputObject.GetDescription()
}

// Enum wraps ast.Enum to avoid exposing graphql-go types outside gql package.
type Enum struct {
	name        string
	description string
	values      []*EnumValue
	// Keep reference to underlying ast for operations within gql package
	astEnum *ast.EnumDefinition
}

// NewEnum creates a wrapper for ast.EnumDefinition
func NewEnum(enum *ast.EnumDefinition) *Enum {
	if enum == nil {
		return nil
	}
	return &Enum{
		name:        enum.Name.Value,
		description: getStringValue(enum.GetDescription()),
		values:      wrapEnumValues(enum.Values),
		astEnum:     enum,
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

// GetName returns the underlying AST name (for NamedTypeDef interface)
func (e *Enum) GetName() *ast.Name {
	return e.astEnum.GetName()
}

// GetDescription returns the underlying AST description (for NamedTypeDef interface)
func (e *Enum) GetDescription() *ast.StringValue {
	return e.astEnum.GetDescription()
}

// Scalar wraps ast.Scalar to avoid exposing graphql-go types outside gql package.
type Scalar struct {
	name        string
	description string
	// Keep reference to underlying ast for operations within gql package
	astScalar *ast.ScalarDefinition
}

// NewScalar creates a wrapper for ast.ScalarDefinition
func NewScalar(scalar *ast.ScalarDefinition) *Scalar {
	if scalar == nil {
		return nil
	}
	return &Scalar{
		name:        scalar.Name.Value,
		description: getStringValue(scalar.GetDescription()),
		astScalar:   scalar,
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

// GetName returns the underlying AST name (for NamedTypeDef interface)
func (s *Scalar) GetName() *ast.Name {
	return s.astScalar.GetName()
}

// GetDescription returns the underlying AST description (for NamedTypeDef interface)
func (s *Scalar) GetDescription() *ast.StringValue {
	return s.astScalar.GetDescription()
}

// Interface wraps ast.Interface to avoid exposing graphql-go types outside gql package.
type Interface struct {
	name        string
	description string
	fields      []*Field
	// Keep reference to underlying ast for operations within gql package
	astInterface *ast.InterfaceDefinition
}

// NewInterface creates a wrapper for ast.InterfaceDefinition
func NewInterface(iface *ast.InterfaceDefinition) *Interface {
	if iface == nil {
		return nil
	}
	return &Interface{
		name:         iface.Name.Value,
		description:  getStringValue(iface.GetDescription()),
		fields:       wrapFields(iface.Fields),
		astInterface: iface,
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

// GetName returns the underlying AST name (for NamedTypeDef interface)
func (i *Interface) GetName() *ast.Name {
	return i.astInterface.GetName()
}

// GetDescription returns the underlying AST description (for NamedTypeDef interface)
func (i *Interface) GetDescription() *ast.StringValue {
	return i.astInterface.GetDescription()
}

// Union wraps ast.Union to avoid exposing graphql-go types outside gql package.
type Union struct {
	name        string
	description string
	types       []*Named
	// Keep reference to underlying ast for operations within gql package
	astUnion *ast.UnionDefinition
}

// NewUnion creates a wrapper for ast.UnionDefinition
func NewUnion(union *ast.UnionDefinition) *Union {
	if union == nil {
		return nil
	}
	return &Union{
		name:        union.Name.Value,
		description: getStringValue(union.GetDescription()),
		types:       wrapNamed(union.Types),
		astUnion:    union,
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

// Types returns the union's member types
func (u *Union) Types() []*Named {
	return u.types
}

// GetName returns the underlying AST name (for NamedTypeDef interface)
func (u *Union) GetName() *ast.Name {
	return u.astUnion.GetName()
}

// GetDescription returns the underlying AST description (for NamedTypeDef interface)
func (u *Union) GetDescription() *ast.StringValue {
	return u.astUnion.GetDescription()
}

// Directive wraps ast.Directive to avoid exposing graphql-go types outside gql package.
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
		name:         directive.Name.Value,
		description:  getStringValue(directive.Description),
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

// Named wraps ast.Named to avoid exposing graphql-go types outside gql package.
type Named struct {
	name string
	// Keep reference to underlying ast for operations within gql package
	astNamed *ast.Named
}

// NewNamed creates a wrapper for ast.Named
func NewNamed(named *ast.Named) *Named {
	if named == nil {
		return nil
	}
	return &Named{
		name:     named.Name.Value,
		astNamed: named,
	}
}

// Name returns the name value
func (n *Named) Name() string {
	return n.name
}

// wrapNamed converts a slice of ast.Named to wrapped types
func wrapNamed(astNamed []*ast.Named) []*Named {
	named := make([]*Named, 0, len(astNamed))
	for _, n := range astNamed {
		named = append(named, NewNamed(n))
	}
	return named
}

// EnumValue wraps ast.EnumValue to avoid exposing graphql-go types outside gql package.
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
		name:         enumValue.Name.Value,
		description:  getStringValue(enumValue.Description),
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
func wrapEnumValues(astEnumValues []*ast.EnumValueDefinition) []*EnumValue {
	enumValues := make([]*EnumValue, 0, len(astEnumValues))
	for _, ev := range astEnumValues {
		enumValues = append(enumValues, NewEnumValue(ev))
	}
	return enumValues
}
