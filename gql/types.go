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
