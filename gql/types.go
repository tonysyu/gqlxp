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
		description: GetStringValue(field.GetDescription()),
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
	return GetTypeString(f.fieldType)
}

// TypeName returns the named type (unwrapping List and NonNull)
func (f *FieldDefinition) TypeName() string {
	named := GetNamedFromType(f.fieldType)
	if named == nil {
		return ""
	}
	return named.Name.Value
}

// Signature returns the full field signature including arguments
func (f *FieldDefinition) Signature() string {
	return GetFieldDefinitionString(f.astField)
}

// HasArguments returns true if the field has arguments
func (f *FieldDefinition) HasArguments() bool {
	return len(f.arguments) > 0
}

// ArgumentCount returns the number of arguments
func (f *FieldDefinition) ArgumentCount() int {
	return len(f.arguments)
}

// Arguments returns the field's arguments as wrapped InputValueDefinition types
func (f *FieldDefinition) Arguments() []*InputValueDefinition {
	return WrapInputValueDefinitions(f.arguments)
}

// ResolveResultType returns the NamedTypeDef for this field's return type.
// Returns an error if the type is a built-in scalar (String, Int, Boolean, etc.)
func (f *FieldDefinition) ResolveResultType(schema *GraphQLSchema) (NamedTypeDef, error) {
	named := GetNamedFromType(f.fieldType)
	return schema.NamedToTypeDefinition(named)
}

// QueryDefinition is an alias for FieldDefinition to make the intent clearer when working with Query fields.
// Use schema.GetQueryField() or schema.GetQueryFields() to obtain QueryDefinition instances.
type QueryDefinition = FieldDefinition

// WrapFieldDefinitions converts a slice of ast.FieldDefinition to wrapped FieldDefinition types.
// This is useful for converting fields from ObjectDefinition, InterfaceDefinition, etc.
func WrapFieldDefinitions(astFields []*ast.FieldDefinition) []*FieldDefinition {
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
		description:   GetStringValue(inputValue.Description),
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
	return GetTypeString(i.inputType)
}

// Signature returns the full input value signature (name: Type = defaultValue)
func (i *InputValueDefinition) Signature() string {
	return GetInputValueDefinitionString(i.astInputValue)
}

// WrapInputValueDefinitions converts a slice of ast.InputValueDefinition to wrapped types
func WrapInputValueDefinitions(astInputValues []*ast.InputValueDefinition) []*InputValueDefinition {
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
		description: GetStringValue(obj.GetDescription()),
		interfaces:  WrapNamed(obj.Interfaces),
		fields:      WrapFieldDefinitions(obj.Fields),
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

// GetKind returns the underlying AST kind (for ast.TypeDefinition interface)
func (o *ObjectDefinition) GetKind() string {
	return o.astObject.GetKind()
}

// GetLoc returns the underlying AST location (for ast.Node interface)
func (o *ObjectDefinition) GetLoc() *ast.Location {
	return o.astObject.GetLoc()
}

// GetOperation returns the underlying AST operation (for ast.TypeDefinition interface)
func (o *ObjectDefinition) GetOperation() string {
	return o.astObject.GetOperation()
}

// GetVariableDefinitions returns the underlying AST variable definitions (for ast.TypeDefinition interface)
func (o *ObjectDefinition) GetVariableDefinitions() []*ast.VariableDefinition {
	return o.astObject.GetVariableDefinitions()
}

// GetSelectionSet returns the underlying AST selection set (for ast.TypeDefinition interface)
func (o *ObjectDefinition) GetSelectionSet() *ast.SelectionSet {
	return o.astObject.GetSelectionSet()
}

// WrapObjectDefinitions converts a slice of ast.ObjectDefinition to wrapped types
func WrapObjectDefinitions(astObjects []*ast.ObjectDefinition) []*ObjectDefinition {
	objects := make([]*ObjectDefinition, 0, len(astObjects))
	for _, obj := range astObjects {
		objects = append(objects, NewObjectDefinition(obj))
	}
	return objects
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
		description:    GetStringValue(input.GetDescription()),
		fields:         WrapInputValueDefinitions(input.Fields),
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

// GetKind returns the underlying AST kind (for ast.TypeDefinition interface)
func (i *InputObjectDefinition) GetKind() string {
	return i.astInputObject.GetKind()
}

// GetLoc returns the underlying AST location (for ast.Node interface)
func (i *InputObjectDefinition) GetLoc() *ast.Location {
	return i.astInputObject.GetLoc()
}

// GetOperation returns the underlying AST operation (for ast.TypeDefinition interface)
func (i *InputObjectDefinition) GetOperation() string {
	return i.astInputObject.GetOperation()
}

// GetVariableDefinitions returns the underlying AST variable definitions (for ast.TypeDefinition interface)
func (i *InputObjectDefinition) GetVariableDefinitions() []*ast.VariableDefinition {
	return i.astInputObject.GetVariableDefinitions()
}

// GetSelectionSet returns the underlying AST selection set (for ast.TypeDefinition interface)
func (i *InputObjectDefinition) GetSelectionSet() *ast.SelectionSet {
	return i.astInputObject.GetSelectionSet()
}

// WrapInputObjectDefinitions converts a slice of ast.InputObjectDefinition to wrapped types
func WrapInputObjectDefinitions(astInputs []*ast.InputObjectDefinition) []*InputObjectDefinition {
	inputs := make([]*InputObjectDefinition, 0, len(astInputs))
	for _, input := range astInputs {
		inputs = append(inputs, NewInputObjectDefinition(input))
	}
	return inputs
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
		description: GetStringValue(enum.GetDescription()),
		values:      WrapEnumValueDefinitions(enum.Values),
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

// GetKind returns the underlying AST kind (for ast.TypeDefinition interface)
func (e *EnumDefinition) GetKind() string {
	return e.astEnum.GetKind()
}

// GetLoc returns the underlying AST location (for ast.Node interface)
func (e *EnumDefinition) GetLoc() *ast.Location {
	return e.astEnum.GetLoc()
}

// GetOperation returns the underlying AST operation (for ast.TypeDefinition interface)
func (e *EnumDefinition) GetOperation() string {
	return e.astEnum.GetOperation()
}

// GetVariableDefinitions returns the underlying AST variable definitions (for ast.TypeDefinition interface)
func (e *EnumDefinition) GetVariableDefinitions() []*ast.VariableDefinition {
	return e.astEnum.GetVariableDefinitions()
}

// GetSelectionSet returns the underlying AST selection set (for ast.TypeDefinition interface)
func (e *EnumDefinition) GetSelectionSet() *ast.SelectionSet {
	return e.astEnum.GetSelectionSet()
}

// WrapEnumDefinitions converts a slice of ast.EnumDefinition to wrapped types
func WrapEnumDefinitions(astEnums []*ast.EnumDefinition) []*EnumDefinition {
	enums := make([]*EnumDefinition, 0, len(astEnums))
	for _, enum := range astEnums {
		enums = append(enums, NewEnumDefinition(enum))
	}
	return enums
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
		description: GetStringValue(scalar.GetDescription()),
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

// GetKind returns the underlying AST kind (for ast.TypeDefinition interface)
func (s *ScalarDefinition) GetKind() string {
	return s.astScalar.GetKind()
}

// GetLoc returns the underlying AST location (for ast.Node interface)
func (s *ScalarDefinition) GetLoc() *ast.Location {
	return s.astScalar.GetLoc()
}

// GetOperation returns the underlying AST operation (for ast.TypeDefinition interface)
func (s *ScalarDefinition) GetOperation() string {
	return s.astScalar.GetOperation()
}

// GetVariableDefinitions returns the underlying AST variable definitions (for ast.TypeDefinition interface)
func (s *ScalarDefinition) GetVariableDefinitions() []*ast.VariableDefinition {
	return s.astScalar.GetVariableDefinitions()
}

// GetSelectionSet returns the underlying AST selection set (for ast.TypeDefinition interface)
func (s *ScalarDefinition) GetSelectionSet() *ast.SelectionSet {
	return s.astScalar.GetSelectionSet()
}

// WrapScalarDefinitions converts a slice of ast.ScalarDefinition to wrapped types
func WrapScalarDefinitions(astScalars []*ast.ScalarDefinition) []*ScalarDefinition {
	scalars := make([]*ScalarDefinition, 0, len(astScalars))
	for _, scalar := range astScalars {
		scalars = append(scalars, NewScalarDefinition(scalar))
	}
	return scalars
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
		description:  GetStringValue(iface.GetDescription()),
		fields:       WrapFieldDefinitions(iface.Fields),
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

// GetKind returns the underlying AST kind (for ast.TypeDefinition interface)
func (i *InterfaceDefinition) GetKind() string {
	return i.astInterface.GetKind()
}

// GetLoc returns the underlying AST location (for ast.Node interface)
func (i *InterfaceDefinition) GetLoc() *ast.Location {
	return i.astInterface.GetLoc()
}

// GetOperation returns the underlying AST operation (for ast.TypeDefinition interface)
func (i *InterfaceDefinition) GetOperation() string {
	return i.astInterface.GetOperation()
}

// GetVariableDefinitions returns the underlying AST variable definitions (for ast.TypeDefinition interface)
func (i *InterfaceDefinition) GetVariableDefinitions() []*ast.VariableDefinition {
	return i.astInterface.GetVariableDefinitions()
}

// GetSelectionSet returns the underlying AST selection set (for ast.TypeDefinition interface)
func (i *InterfaceDefinition) GetSelectionSet() *ast.SelectionSet {
	return i.astInterface.GetSelectionSet()
}

// WrapInterfaceDefinitions converts a slice of ast.InterfaceDefinition to wrapped types
func WrapInterfaceDefinitions(astInterfaces []*ast.InterfaceDefinition) []*InterfaceDefinition {
	interfaces := make([]*InterfaceDefinition, 0, len(astInterfaces))
	for _, iface := range astInterfaces {
		interfaces = append(interfaces, NewInterfaceDefinition(iface))
	}
	return interfaces
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
		description: GetStringValue(union.GetDescription()),
		types:       WrapNamed(union.Types),
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

// GetKind returns the underlying AST kind (for ast.TypeDefinition interface)
func (u *UnionDefinition) GetKind() string {
	return u.astUnion.GetKind()
}

// GetLoc returns the underlying AST location (for ast.Node interface)
func (u *UnionDefinition) GetLoc() *ast.Location {
	return u.astUnion.GetLoc()
}

// GetOperation returns the underlying AST operation (for ast.TypeDefinition interface)
func (u *UnionDefinition) GetOperation() string {
	return u.astUnion.GetOperation()
}

// GetVariableDefinitions returns the underlying AST variable definitions (for ast.TypeDefinition interface)
func (u *UnionDefinition) GetVariableDefinitions() []*ast.VariableDefinition {
	return u.astUnion.GetVariableDefinitions()
}

// GetSelectionSet returns the underlying AST selection set (for ast.TypeDefinition interface)
func (u *UnionDefinition) GetSelectionSet() *ast.SelectionSet {
	return u.astUnion.GetSelectionSet()
}

// WrapUnionDefinitions converts a slice of ast.UnionDefinition to wrapped types
func WrapUnionDefinitions(astUnions []*ast.UnionDefinition) []*UnionDefinition {
	unions := make([]*UnionDefinition, 0, len(astUnions))
	for _, union := range astUnions {
		unions = append(unions, NewUnionDefinition(union))
	}
	return unions
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
		description:  GetStringValue(directive.Description),
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

// WrapDirectiveDefinitions converts a slice of ast.DirectiveDefinition to wrapped types
func WrapDirectiveDefinitions(astDirectives []*ast.DirectiveDefinition) []*DirectiveDefinition {
	directives := make([]*DirectiveDefinition, 0, len(astDirectives))
	for _, directive := range astDirectives {
		directives = append(directives, NewDirectiveDefinition(directive))
	}
	return directives
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

// WrapNamed converts a slice of ast.Named to wrapped types
func WrapNamed(astNamed []*ast.Named) []*Named {
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
		description:  GetStringValue(enumValue.Description),
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

// WrapEnumValueDefinitions converts a slice of ast.EnumValueDefinition to wrapped types
func WrapEnumValueDefinitions(astEnumValues []*ast.EnumValueDefinition) []*EnumValueDefinition {
	enumValues := make([]*EnumValueDefinition, 0, len(astEnumValues))
	for _, ev := range astEnumValues {
		enumValues = append(enumValues, NewEnumValueDefinition(ev))
	}
	return enumValues
}
