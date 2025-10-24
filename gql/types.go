package gql

import "github.com/graphql-go/graphql/language/ast"

// FieldDefinition wraps ast.FieldDefinition to avoid exposing graphql-go types outside gql package.
// This can represent query fields, mutation fields, or object fields.
type FieldDefinition struct {
	name        string
	description string
	fieldType   ast.Type
	arguments   []*ast.InputValueDefinition
	// Keep reference to underlying ast for operations within gql package
	astField *ast.FieldDefinition
}

// NewFieldDefinition creates a wrapper for ast.FieldDefinition
func NewFieldDefinition(field *ast.FieldDefinition) *FieldDefinition {
	if field == nil {
		return nil
	}
	return &FieldDefinition{
		name:        field.Name.Value,
		description: getStringValue(field.GetDescription()),
		fieldType:   field.Type,
		arguments:   field.Arguments,
		astField:    field,
	}
}

// Name returns the field name
func (f *FieldDefinition) Name() string {
	return f.name
}

// Description returns the field description
func (f *FieldDefinition) Description() string {
	return f.description
}

// TypeString returns the string representation of the field's type
func (f *FieldDefinition) TypeString() string {
	return getTypeString(f.fieldType)
}

// TypeName returns the named type (unwrapping List and NonNull)
func (f *FieldDefinition) TypeName() string {
	named := getNamedFromType(f.fieldType)
	if named == nil {
		return ""
	}
	return named.Name.Value
}

// Signature returns the full field signature including arguments
func (f *FieldDefinition) Signature() string {
	return getFieldDefinitionString(f.astField)
}

// Arguments returns the field's arguments as wrapped InputValueDefinition types
func (f *FieldDefinition) Arguments() []*InputValueDefinition {
	return wrapInputValueDefinitions(f.arguments)
}

// ResolveResultType returns the NamedTypeDef for this field's return type.
// Returns an error if the type is a built-in scalar (String, Int, Boolean, etc.)
func (f *FieldDefinition) ResolveResultType(schema *GraphQLSchema) (NamedTypeDef, error) {
	named := getNamedFromType(f.fieldType)
	return schema.NamedToTypeDefinition(named)
}

// wrapFieldDefinitions converts a slice of ast.FieldDefinition to wrapped FieldDefinition types.
// This is useful for converting fields from ObjectDefinition, InterfaceDefinition, etc.
func wrapFieldDefinitions(astFields []*ast.FieldDefinition) []*FieldDefinition {
	fields := make([]*FieldDefinition, 0, len(astFields))
	for _, field := range astFields {
		fields = append(fields, NewFieldDefinition(field))
	}
	return fields
}

// InputValueDefinition wraps ast.InputValueDefinition to avoid exposing graphql-go types outside gql package.
// This represents input arguments or input object fields.
type InputValueDefinition struct {
	name        string
	description string
	inputType   ast.Type
	// Keep reference to underlying ast for operations within gql package
	astInputValue *ast.InputValueDefinition
}

// NewInputValueDefinition creates a wrapper for ast.InputValueDefinition
func NewInputValueDefinition(inputValue *ast.InputValueDefinition) *InputValueDefinition {
	if inputValue == nil {
		return nil
	}
	return &InputValueDefinition{
		name:          inputValue.Name.Value,
		description:   getStringValue(inputValue.Description),
		inputType:     inputValue.Type,
		astInputValue: inputValue,
	}
}

// Name returns the input value name
func (i *InputValueDefinition) Name() string {
	return i.name
}

// Description returns the input value description
func (i *InputValueDefinition) Description() string {
	return i.description
}

// TypeString returns the string representation of the input value's type
func (i *InputValueDefinition) TypeString() string {
	return getTypeString(i.inputType)
}

// Signature returns the full input value signature (name: Type = defaultValue)
func (i *InputValueDefinition) Signature() string {
	return getInputValueDefinitionString(i.astInputValue)
}

// wrapInputValueDefinitions converts a slice of ast.InputValueDefinition to wrapped types
func wrapInputValueDefinitions(astInputValues []*ast.InputValueDefinition) []*InputValueDefinition {
	inputValues := make([]*InputValueDefinition, 0, len(astInputValues))
	for _, iv := range astInputValues {
		inputValues = append(inputValues, NewInputValueDefinition(iv))
	}
	return inputValues
}

// ObjectDefinition wraps ast.ObjectDefinition to avoid exposing graphql-go types outside gql package.
type ObjectDefinition struct {
	name        string
	description string
	interfaces  []*Named
	fields      []*FieldDefinition
	// Keep reference to underlying ast for operations within gql package
	astObject *ast.ObjectDefinition
}

// NewObjectDefinition creates a wrapper for ast.ObjectDefinition
func NewObjectDefinition(obj *ast.ObjectDefinition) *ObjectDefinition {
	if obj == nil {
		return nil
	}
	return &ObjectDefinition{
		name:        obj.Name.Value,
		description: getStringValue(obj.GetDescription()),
		interfaces:  wrapNamed(obj.Interfaces),
		fields:      wrapFieldDefinitions(obj.Fields),
		astObject:   obj,
	}
}

// Name returns the object name
func (o *ObjectDefinition) Name() string {
	return o.name
}

// Description returns the object description
func (o *ObjectDefinition) Description() string {
	return o.description
}

// Interfaces returns the interfaces this object implements
func (o *ObjectDefinition) Interfaces() []*Named {
	return o.interfaces
}

// Fields returns the object's fields
func (o *ObjectDefinition) Fields() []*FieldDefinition {
	return o.fields
}

// GetName returns the underlying AST name (for NamedTypeDef interface)
func (o *ObjectDefinition) GetName() *ast.Name {
	return o.astObject.GetName()
}

// GetDescription returns the underlying AST description (for NamedTypeDef interface)
func (o *ObjectDefinition) GetDescription() *ast.StringValue {
	return o.astObject.GetDescription()
}

// InputObjectDefinition wraps ast.InputObjectDefinition to avoid exposing graphql-go types outside gql package.
type InputObjectDefinition struct {
	name        string
	description string
	fields      []*InputValueDefinition
	// Keep reference to underlying ast for operations within gql package
	astInputObject *ast.InputObjectDefinition
}

// NewInputObjectDefinition creates a wrapper for ast.InputObjectDefinition
func NewInputObjectDefinition(input *ast.InputObjectDefinition) *InputObjectDefinition {
	if input == nil {
		return nil
	}
	return &InputObjectDefinition{
		name:           input.Name.Value,
		description:    getStringValue(input.GetDescription()),
		fields:         wrapInputValueDefinitions(input.Fields),
		astInputObject: input,
	}
}

// Name returns the input object name
func (i *InputObjectDefinition) Name() string {
	return i.name
}

// Description returns the input object description
func (i *InputObjectDefinition) Description() string {
	return i.description
}

// Fields returns the input object's fields
func (i *InputObjectDefinition) Fields() []*InputValueDefinition {
	return i.fields
}

// GetName returns the underlying AST name (for NamedTypeDef interface)
func (i *InputObjectDefinition) GetName() *ast.Name {
	return i.astInputObject.GetName()
}

// GetDescription returns the underlying AST description (for NamedTypeDef interface)
func (i *InputObjectDefinition) GetDescription() *ast.StringValue {
	return i.astInputObject.GetDescription()
}

// EnumDefinition wraps ast.EnumDefinition to avoid exposing graphql-go types outside gql package.
type EnumDefinition struct {
	name        string
	description string
	values      []*EnumValueDefinition
	// Keep reference to underlying ast for operations within gql package
	astEnum *ast.EnumDefinition
}

// NewEnumDefinition creates a wrapper for ast.EnumDefinition
func NewEnumDefinition(enum *ast.EnumDefinition) *EnumDefinition {
	if enum == nil {
		return nil
	}
	return &EnumDefinition{
		name:        enum.Name.Value,
		description: getStringValue(enum.GetDescription()),
		values:      wrapEnumValueDefinitions(enum.Values),
		astEnum:     enum,
	}
}

// Name returns the enum name
func (e *EnumDefinition) Name() string {
	return e.name
}

// Description returns the enum description
func (e *EnumDefinition) Description() string {
	return e.description
}

// Values returns the enum values
func (e *EnumDefinition) Values() []*EnumValueDefinition {
	return e.values
}

// GetName returns the underlying AST name (for NamedTypeDef interface)
func (e *EnumDefinition) GetName() *ast.Name {
	return e.astEnum.GetName()
}

// GetDescription returns the underlying AST description (for NamedTypeDef interface)
func (e *EnumDefinition) GetDescription() *ast.StringValue {
	return e.astEnum.GetDescription()
}

// ScalarDefinition wraps ast.ScalarDefinition to avoid exposing graphql-go types outside gql package.
type ScalarDefinition struct {
	name        string
	description string
	// Keep reference to underlying ast for operations within gql package
	astScalar *ast.ScalarDefinition
}

// NewScalarDefinition creates a wrapper for ast.ScalarDefinition
func NewScalarDefinition(scalar *ast.ScalarDefinition) *ScalarDefinition {
	if scalar == nil {
		return nil
	}
	return &ScalarDefinition{
		name:        scalar.Name.Value,
		description: getStringValue(scalar.GetDescription()),
		astScalar:   scalar,
	}
}

// Name returns the scalar name
func (s *ScalarDefinition) Name() string {
	return s.name
}

// Description returns the scalar description
func (s *ScalarDefinition) Description() string {
	return s.description
}

// GetName returns the underlying AST name (for NamedTypeDef interface)
func (s *ScalarDefinition) GetName() *ast.Name {
	return s.astScalar.GetName()
}

// GetDescription returns the underlying AST description (for NamedTypeDef interface)
func (s *ScalarDefinition) GetDescription() *ast.StringValue {
	return s.astScalar.GetDescription()
}

// InterfaceDefinition wraps ast.InterfaceDefinition to avoid exposing graphql-go types outside gql package.
type InterfaceDefinition struct {
	name        string
	description string
	fields      []*FieldDefinition
	// Keep reference to underlying ast for operations within gql package
	astInterface *ast.InterfaceDefinition
}

// NewInterfaceDefinition creates a wrapper for ast.InterfaceDefinition
func NewInterfaceDefinition(iface *ast.InterfaceDefinition) *InterfaceDefinition {
	if iface == nil {
		return nil
	}
	return &InterfaceDefinition{
		name:         iface.Name.Value,
		description:  getStringValue(iface.GetDescription()),
		fields:       wrapFieldDefinitions(iface.Fields),
		astInterface: iface,
	}
}

// Name returns the interface name
func (i *InterfaceDefinition) Name() string {
	return i.name
}

// Description returns the interface description
func (i *InterfaceDefinition) Description() string {
	return i.description
}

// Fields returns the interface's fields
func (i *InterfaceDefinition) Fields() []*FieldDefinition {
	return i.fields
}

// GetName returns the underlying AST name (for NamedTypeDef interface)
func (i *InterfaceDefinition) GetName() *ast.Name {
	return i.astInterface.GetName()
}

// GetDescription returns the underlying AST description (for NamedTypeDef interface)
func (i *InterfaceDefinition) GetDescription() *ast.StringValue {
	return i.astInterface.GetDescription()
}

// UnionDefinition wraps ast.UnionDefinition to avoid exposing graphql-go types outside gql package.
type UnionDefinition struct {
	name        string
	description string
	types       []*Named
	// Keep reference to underlying ast for operations within gql package
	astUnion *ast.UnionDefinition
}

// NewUnionDefinition creates a wrapper for ast.UnionDefinition
func NewUnionDefinition(union *ast.UnionDefinition) *UnionDefinition {
	if union == nil {
		return nil
	}
	return &UnionDefinition{
		name:        union.Name.Value,
		description: getStringValue(union.GetDescription()),
		types:       wrapNamed(union.Types),
		astUnion:    union,
	}
}

// Name returns the union name
func (u *UnionDefinition) Name() string {
	return u.name
}

// Description returns the union description
func (u *UnionDefinition) Description() string {
	return u.description
}

// Types returns the union's member types
func (u *UnionDefinition) Types() []*Named {
	return u.types
}

// GetName returns the underlying AST name (for NamedTypeDef interface)
func (u *UnionDefinition) GetName() *ast.Name {
	return u.astUnion.GetName()
}

// GetDescription returns the underlying AST description (for NamedTypeDef interface)
func (u *UnionDefinition) GetDescription() *ast.StringValue {
	return u.astUnion.GetDescription()
}

// DirectiveDefinition wraps ast.DirectiveDefinition to avoid exposing graphql-go types outside gql package.
type DirectiveDefinition struct {
	name        string
	description string
	// Keep reference to underlying ast for operations within gql package
	astDirective *ast.DirectiveDefinition
}

// NewDirectiveDefinition creates a wrapper for ast.DirectiveDefinition
func NewDirectiveDefinition(directive *ast.DirectiveDefinition) *DirectiveDefinition {
	if directive == nil {
		return nil
	}
	return &DirectiveDefinition{
		name:         directive.Name.Value,
		description:  getStringValue(directive.Description),
		astDirective: directive,
	}
}

// Name returns the directive name
func (d *DirectiveDefinition) Name() string {
	return d.name
}

// Description returns the directive description
func (d *DirectiveDefinition) Description() string {
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

// EnumValueDefinition wraps ast.EnumValueDefinition to avoid exposing graphql-go types outside gql package.
type EnumValueDefinition struct {
	name        string
	description string
	// Keep reference to underlying ast for operations within gql package
	astEnumValue *ast.EnumValueDefinition
}

// NewEnumValueDefinition creates a wrapper for ast.EnumValueDefinition
func NewEnumValueDefinition(enumValue *ast.EnumValueDefinition) *EnumValueDefinition {
	if enumValue == nil {
		return nil
	}
	return &EnumValueDefinition{
		name:         enumValue.Name.Value,
		description:  getStringValue(enumValue.Description),
		astEnumValue: enumValue,
	}
}

// Name returns the enum value name
func (e *EnumValueDefinition) Name() string {
	return e.name
}

// Description returns the enum value description
func (e *EnumValueDefinition) Description() string {
	return e.description
}

// wrapEnumValueDefinitions converts a slice of ast.EnumValueDefinition to wrapped types
func wrapEnumValueDefinitions(astEnumValues []*ast.EnumValueDefinition) []*EnumValueDefinition {
	enumValues := make([]*EnumValueDefinition, 0, len(astEnumValues))
	for _, ev := range astEnumValues {
		enumValues = append(enumValues, NewEnumValueDefinition(ev))
	}
	return enumValues
}
