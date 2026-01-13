package gql

import "github.com/vektah/gqlparser/v2/ast"

// Interface for GQL Type Definitions
type TypeDef interface {
	Name() string
	Description() string
}

var _ TypeDef = (*Argument)(nil)
var _ TypeDef = (*DirectiveDef)(nil)
var _ TypeDef = (*Enum)(nil)
var _ TypeDef = (*Field)(nil)
var _ TypeDef = (*InputObject)(nil)
var _ TypeDef = (*Interface)(nil)
var _ TypeDef = (*Object)(nil)
var _ TypeDef = (*Scalar)(nil)
var _ TypeDef = (*Union)(nil)

type Callable interface {
	Name() string
	Arguments() []*Argument
}

var _ Callable = (*DirectiveDef)(nil)
var _ Callable = (*Field)(nil)

// gqlType interface for types that have a Name field (both ast and wrapped types)
type gqlType interface {
	*Field | *Object | *InputObject |
		*Enum | *Scalar | *Interface |
		*Union | *DirectiveDef
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

// ObjectTypeName returns the object type name (unwrapping List and NonNull)
// While TypeString would return [MyObject!]!, this just returns MyObject.
func (f *Field) ObjectTypeName() string {
	return getNamedTypeName(f.fieldType)
}

// Signature returns the full field signature including arguments.
// This can be as simple as:
//
//	"name: Type"
//
// Or can be as complex as:
//
//	"childConnection(first: Int = 10, after: ID!): ChildConnection!"
func (f *Field) Signature() string {
	return getFieldString(f)
}

// FormatSignature returns the field signature with optional multiline formatting.
// If maxWidth <= 0, always uses inline format (same as Signature()).
// If maxWidth > 0 and inline signature exceeds maxWidth, formats arguments on separate lines.
func (f *Field) FormatSignature(maxWidth int) string {
	return formatFieldStringWithWidth(f, maxWidth)
}

// Arguments returns the field's arguments as wrapped Arguments
func (f *Field) Arguments() []*Argument {
	return wrapArguments(f.arguments)
}

// Directives returns the field's applied directives
func (f *Field) Directives() []*AppliedDirective {
	if f.astField == nil {
		return nil
	}
	return wrapDirectives(f.astField.Directives)
}

// Default value if defined (empty-string if not)
func (f *Field) DefaultValue() string {
	return formatValue(f.astField.DefaultValue)
}

// ResolveObjectTypeDef returns the TypeDef for this field's type.
// Returns an error if the type is a built-in scalar (String, Int, Boolean, etc.)
func (f *Field) ResolveObjectTypeDef(schema *GraphQLSchema) (TypeDef, error) {
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
//
// See https://spec.graphql.org/October2021/#sec-Field-Arguments
//
// Argument is similar to Field, with the exception that it doesn't have an Arguments() method,
// since an Argument can't have sub-arguments.
type Argument struct {
	name        string
	description string
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
		astArg:      arg,
	}
}

func (a *Argument) Name() string        { return a.name }
func (a *Argument) Description() string { return a.description }

// TypeString returns the string representation of the argument's type
func (a *Argument) TypeString() string {
	return getTypeString(a.astArg.Type)
}

// ObjectTypeName returns the input type name (unwrapping List and NonNull)
// While TypeString would return [MyInput!]!, this just returns MyInput.
func (a *Argument) ObjectTypeName() string {
	return getNamedTypeName(a.astArg.Type)
}

// Signature returns the full argument signature (name: Type = defaultValue)
func (a *Argument) Signature() string {
	return getArgumentString(a)
}

// FormatSignature returns the argument signature (always single line for arguments)
// The maxWidth parameter is accepted for consistency with Field.FormatSignature but is ignored.
func (a *Argument) FormatSignature(maxWidth int) string {
	return getArgumentString(a)
}

// Default value if defined (empty-string if not)
func (a *Argument) DefaultValue() string {
	return formatValue(a.astArg.DefaultValue)
}

// Directives returns the argument's applied directives
func (a *Argument) Directives() []*AppliedDirective {
	if a.astArg == nil {
		return nil
	}
	return wrapDirectives(a.astArg.Directives)
}

// ResolveObjectTypeDef returns the TypeDef for this argument's input type.
// Returns an error if the type is a built-in scalar (String, Int, Boolean, etc.)
func (a *Argument) ResolveObjectTypeDef(schema *GraphQLSchema) (TypeDef, error) {
	return schema.NamedToTypeDef(a.ObjectTypeName())
}

// wrapArguments converts a slice of ast.ArgumentDefinition to wrapped Argument types
func wrapArguments(args ast.ArgumentDefinitionList) []*Argument {
	arguments := make([]*Argument, 0, len(args))
	for _, arg := range args {
		arguments = append(arguments, newArgument(arg))
	}
	return arguments
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

// Directives returns the object's applied directives
func (o *Object) Directives() []*AppliedDirective {
	if o.astDef == nil {
		return nil
	}
	return wrapDirectives(o.astDef.Directives)
}

// InputObject represents GraphQL input objects.
// See https://spec.graphql.org/October2021/#sec-Input-Objects
type InputObject struct {
	name        string
	description string
	fields      []*Field
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
		fields:      wrapFields(def.Fields),
		astDef:      def,
	}
}

func (i *InputObject) Name() string        { return i.name }
func (i *InputObject) Description() string { return i.description }

// Fields returns the input object's fields
func (i *InputObject) Fields() []*Field {
	return i.fields
}

// Directives returns the input object's applied directives
func (i *InputObject) Directives() []*AppliedDirective {
	if i.astDef == nil {
		return nil
	}
	return wrapDirectives(i.astDef.Directives)
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

// Directives returns the enum's applied directives
func (e *Enum) Directives() []*AppliedDirective {
	if e.astDef == nil {
		return nil
	}
	return wrapDirectives(e.astDef.Directives)
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

// Directives returns the scalar's applied directives
func (s *Scalar) Directives() []*AppliedDirective {
	if s.astDef == nil {
		return nil
	}
	return wrapDirectives(s.astDef.Directives)
}

// Interface represents GraphQL interface types.
// See https://spec.graphql.org/October2021/#sec-Interfaces
type Interface struct {
	name        string
	description string
	interfaces  []string
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
		interfaces:  def.Interfaces,
		fields:      wrapFields(def.Fields),
		astDef:      def,
	}
}

func (i *Interface) Name() string        { return i.name }
func (i *Interface) Description() string { return i.description }

// Interfaces returns the interfaces this interface implements
func (i *Interface) Interfaces() []string {
	return i.interfaces
}

// Fields returns the interface's fields
func (i *Interface) Fields() []*Field {
	return i.fields
}

// Directives returns the interface's applied directives
func (i *Interface) Directives() []*AppliedDirective {
	if i.astDef == nil {
		return nil
	}
	return wrapDirectives(i.astDef.Directives)
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

// Directives returns the union's applied directives
func (u *Union) Directives() []*AppliedDirective {
	if u.astDef == nil {
		return nil
	}
	return wrapDirectives(u.astDef.Directives)
}

// DirectiveDef represents GraphQL directive type definitions.
// See https://spec.graphql.org/October2021/#sec-Type-System.Directives
type DirectiveDef struct {
	name         string
	description  string
	arguments    ast.ArgumentDefinitionList
	locations    []string
	isRepeatable bool
	// Keep reference to underlying ast for operations within gql package
	astDirective *ast.DirectiveDefinition
}

// newDirectiveDef creates a wrapper for ast.DirectiveDefinition
func newDirectiveDef(directive *ast.DirectiveDefinition) *DirectiveDef {
	if directive == nil {
		return nil
	}
	// Convert DirectiveLocation to strings
	locations := make([]string, len(directive.Locations))
	for i, loc := range directive.Locations {
		locations[i] = string(loc)
	}
	return &DirectiveDef{
		name:         directive.Name,
		description:  directive.Description,
		arguments:    directive.Arguments,
		locations:    locations,
		isRepeatable: directive.IsRepeatable,
		astDirective: directive,
	}
}

func (d *DirectiveDef) Name() string        { return d.name }
func (d *DirectiveDef) Description() string { return d.description }

// Signature returns the directive signature including arguments.
// This can be as simple as:
//
//	"directiveName"
//
// Or can be as complex as:
//
//	"directiveWithArgs(count: Int = 10, special: Boolean!)"
func (d *DirectiveDef) Signature() string {
	return formatCallable(d, formatOpts{prefix: "@"})
}

// FormatSignature returns the field signature with optional multiline formatting.
// If maxWidth <= 0, always uses inline format (same as Signature()).
// If maxWidth > 0 and inline signature exceeds maxWidth, formats arguments on separate lines.
func (d *DirectiveDef) FormatSignature(maxWidth int) string {
	return formatCallable(d, formatOpts{prefix: "@", maxWidth: maxWidth})
}

// Arguments returns the directive's arguments as wrapped Arguments
func (d *DirectiveDef) Arguments() []*Argument {
	return wrapArguments(d.arguments)
}

// Locations returns the directive's valid locations
func (d *DirectiveDef) Locations() []string {
	return d.locations
}

// AppliedDirective represents a directive applied to a schema element (field, type, etc.).
type AppliedDirective struct {
	name             string
	appliedDirective *ast.Directive
}

func (d *AppliedDirective) Name() string        { return d.name }
func (d *AppliedDirective) Description() string { return "" }

// Signature returns the applied directive with its argument values.
// e.g., "@deprecated(reason: "Use newField")"
func (d *AppliedDirective) Signature() string {
	return formatAppliedDirective(d.appliedDirective)
}

// AppliedArguments returns the directive's applied argument values as a map
func (d *AppliedDirective) AppliedArguments() map[string]*ast.Value {
	if d.appliedDirective == nil {
		return nil
	}
	result := make(map[string]*ast.Value)
	for _, arg := range d.appliedDirective.Arguments {
		result[arg.Name] = arg.Value
	}
	return result
}

// FormattedArguments returns the directive's applied argument values as formatted strings
func (d *AppliedDirective) FormattedArguments() map[string]string {
	if d.appliedDirective == nil {
		return nil
	}
	result := make(map[string]string)
	for _, arg := range d.appliedDirective.Arguments {
		result[arg.Name] = formatValue(arg.Value)
	}
	return result
}

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

// Signature returns the enum value signature
func (e *EnumValue) Signature() string {
	return e.name
}

// Directives returns the enum value's applied directives
func (e *EnumValue) Directives() []*AppliedDirective {
	if e.astEnumValue == nil {
		return nil
	}
	return wrapDirectives(e.astEnumValue.Directives)
}

// wrapEnumValues converts a slice of ast.EnumValueDefinition to wrapped types
func wrapEnumValues(astEnumValues ast.EnumValueList) []*EnumValue {
	enumValues := make([]*EnumValue, 0, len(astEnumValues))
	for _, ev := range astEnumValues {
		enumValues = append(enumValues, newEnumValue(ev))
	}
	return enumValues
}

// wrapDirectives converts a slice of ast.Directive to wrapped AppliedDirective types
func wrapDirectives(astDirectives ast.DirectiveList) []*AppliedDirective {
	if len(astDirectives) == 0 {
		return nil
	}
	directives := make([]*AppliedDirective, 0, len(astDirectives))
	for _, dir := range astDirectives {
		directives = append(directives, newAppliedDirective(dir))
	}
	return directives
}

// newAppliedDirective creates a wrapper for applied ast.Directive (not DirectiveDefinition)
// Applied directives have arguments with values, not argument definitions
func newAppliedDirective(dir *ast.Directive) *AppliedDirective {
	if dir == nil {
		return nil
	}
	// Note: Applied directives (ast.Directive) are different from directive definitions
	// (ast.DirectiveDefinition). Applied directives have argument values, not definitions.
	return &AppliedDirective{
		name:             dir.Name,
		appliedDirective: dir,
	}
}
