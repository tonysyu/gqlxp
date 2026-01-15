package gqlfmt

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tonysyu/gqlxp/gql"
)

// GenerateJSON generates JSON output for a given type name in the schema
func GenerateJSON(schema gql.GraphQLSchema, typeName string, opts IncludeOptions) (string, error) {
	resolver := gql.NewSchemaResolver(&schema)

	// Handle Query fields (Query.fieldName)
	if strings.HasPrefix(typeName, "Query.") {
		return generateQueryFieldJSON(schema, typeName, resolver, opts)
	}

	// Handle Mutation fields (Mutation.fieldName)
	if strings.HasPrefix(typeName, "Mutation.") {
		return generateMutationFieldJSON(schema, typeName, resolver, opts)
	}

	// Handle Directives (@directiveName)
	if strings.HasPrefix(typeName, "@") {
		return generateDirectiveJSON(schema, typeName, resolver, opts)
	}

	// Handle regular types (Object, Input, Enum, Scalar, Interface, Union)
	return generateTypeJSON(schema, typeName, resolver, opts)
}

// generateQueryFieldJSON generates JSON for a query field
func generateQueryFieldJSON(schema gql.GraphQLSchema, typeName string, resolver gql.TypeResolver, opts IncludeOptions) (string, error) {
	fieldName := strings.TrimPrefix(typeName, "Query.")
	field, ok := schema.Query[fieldName]
	if !ok {
		return "", fmt.Errorf("query field %q not found in schema", fieldName)
	}
	result := convertFieldToJSONWithKind(field, "Query")
	if opts.Usages {
		result.Usages = getUsagesJSON(resolver, field.ObjectTypeName())
	}
	return marshalJSON(result), nil
}

// generateMutationFieldJSON generates JSON for a mutation field
func generateMutationFieldJSON(schema gql.GraphQLSchema, typeName string, resolver gql.TypeResolver, opts IncludeOptions) (string, error) {
	fieldName := strings.TrimPrefix(typeName, "Mutation.")
	field, ok := schema.Mutation[fieldName]
	if !ok {
		return "", fmt.Errorf("mutation field %q not found in schema", fieldName)
	}
	result := convertFieldToJSONWithKind(field, "Mutation")
	if opts.Usages {
		result.Usages = getUsagesJSON(resolver, field.ObjectTypeName())
	}
	return marshalJSON(result), nil
}

// generateDirectiveJSON generates JSON for a directive
func generateDirectiveJSON(schema gql.GraphQLSchema, typeName string, resolver gql.TypeResolver, opts IncludeOptions) (string, error) {
	directiveName := strings.TrimPrefix(typeName, "@")
	directive, ok := schema.Directive[directiveName]
	if !ok {
		return "", fmt.Errorf("directive %q not found in schema", directiveName)
	}
	jsonDir := convertDirectiveToJSON(directive)
	jsonDir.Kind = "Directive"
	if opts.Usages {
		jsonDir.Usages = getUsagesJSON(resolver, directiveName)
	}
	return marshalJSON(jsonDir), nil
}

// generateTypeJSON generates JSON for a type definition
func generateTypeJSON(schema gql.GraphQLSchema, typeName string, resolver gql.TypeResolver, opts IncludeOptions) (string, error) {
	typeDef, err := schema.NamedToTypeDef(typeName)
	if err != nil {
		return "", fmt.Errorf("type %q not found in schema: %w", typeName, err)
	}
	result := convertTypeDefToJSON(typeDef)
	if opts.Usages {
		result.Usages = getUsagesJSON(resolver, typeDef.Name())
	}
	return marshalJSON(result), nil
}

// getUsagesJSON retrieves usages for a type and converts them to JSON format
func getUsagesJSON(resolver gql.TypeResolver, typeName string) []JSONUsage {
	usages, _ := resolver.ResolveUsages(typeName)
	if len(usages) == 0 {
		return nil
	}
	result := make([]JSONUsage, 0, len(usages))
	for _, u := range usages {
		result = append(result, JSONUsage{
			Path:       u.Path,
			ParentType: u.ParentType,
			ParentKind: u.ParentKind,
			FieldName:  u.FieldName,
		})
	}
	return result
}

// JSONField represents a field in JSON format
type JSONField struct {
	Name        string               `json:"name"`
	Kind        string               `json:"kind,omitempty"`
	Type        string               `json:"type"`
	Description string               `json:"description,omitempty"`
	Arguments   []JSONArgument       `json:"arguments,omitempty"`
	Directives  []JSONDirectiveUsage `json:"directives,omitempty"`
	Usages      []JSONUsage          `json:"usages,omitempty"`
}

// JSONArgument represents an argument in JSON format
type JSONArgument struct {
	Name         string               `json:"name"`
	Type         string               `json:"type"`
	Description  string               `json:"description,omitempty"`
	DefaultValue string               `json:"defaultValue,omitempty"`
	Directives   []JSONDirectiveUsage `json:"directives,omitempty"`
}

// JSONDirectiveUsage represents a directive usage in JSON format
type JSONDirectiveUsage struct {
	Name      string         `json:"name"`
	Arguments []JSONArgValue `json:"arguments,omitempty"`
}

// JSONArgValue represents a directive argument value
type JSONArgValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// JSONUsage represents a type usage in JSON format
type JSONUsage struct {
	Path       string `json:"path"`
	ParentType string `json:"parentType"`
	ParentKind string `json:"parentKind"`
	FieldName  string `json:"fieldName"`
}

// JSONEnumValue represents an enum value in JSON format
type JSONEnumValue struct {
	Name        string               `json:"name"`
	Description string               `json:"description,omitempty"`
	Directives  []JSONDirectiveUsage `json:"directives,omitempty"`
}

// JSONTypeDef represents a type definition in JSON format
type JSONTypeDef struct {
	Name        string               `json:"name"`
	Kind        string               `json:"kind"`
	Description string               `json:"description,omitempty"`
	Fields      []JSONField          `json:"fields,omitempty"`
	Values      []JSONEnumValue      `json:"values,omitempty"`
	Types       []string             `json:"types,omitempty"`
	Interfaces  []string             `json:"interfaces,omitempty"`
	Directives  []JSONDirectiveUsage `json:"directives,omitempty"`
	Usages      []JSONUsage          `json:"usages,omitempty"`
}

// JSONDirective represents a directive definition in JSON format
type JSONDirective struct {
	Name        string         `json:"name"`
	Kind        string         `json:"kind"`
	Description string         `json:"description,omitempty"`
	Locations   []string       `json:"locations,omitempty"`
	Arguments   []JSONArgument `json:"arguments,omitempty"`
	Usages      []JSONUsage    `json:"usages,omitempty"`
}

// convertFieldToJSON converts a gql.Field to JSON format
func convertFieldToJSON(field *gql.Field) JSONField {
	return convertFieldToJSONWithKind(field, "")
}

// convertFieldToJSONWithKind converts a gql.Field to JSON format with a specific kind
func convertFieldToJSONWithKind(field *gql.Field, kind string) JSONField {
	return JSONField{
		Name:        field.Name(),
		Kind:        kind,
		Type:        field.TypeString(),
		Description: field.Description(),
		Arguments:   convertArgumentsToJSON(field.Arguments()),
		Directives:  convertDirectivesToJSON(field.Directives()),
	}
}

// convertArgumentsToJSON converts a slice of gql.Argument to JSON format
func convertArgumentsToJSON(args []*gql.Argument) []JSONArgument {
	if len(args) == 0 {
		return nil
	}
	result := make([]JSONArgument, 0, len(args))
	for _, arg := range args {
		jsonArg := JSONArgument{
			Name:        arg.Name(),
			Type:        arg.TypeString(),
			Description: arg.Description(),
			Directives:  convertDirectivesToJSON(arg.Directives()),
		}
		if defaultVal := arg.DefaultValue(); defaultVal != "" {
			jsonArg.DefaultValue = defaultVal
		}
		result = append(result, jsonArg)
	}
	return result
}

// convertDirectivesToJSON converts a slice of gql.AppliedDirective to JSON format
func convertDirectivesToJSON(directives []*gql.AppliedDirective) []JSONDirectiveUsage {
	if len(directives) == 0 {
		return nil
	}
	result := make([]JSONDirectiveUsage, 0, len(directives))
	for _, dir := range directives {
		result = append(result, JSONDirectiveUsage{
			Name:      dir.Name(),
			Arguments: convertDirectiveArgsToJSON(dir.FormattedArguments()),
		})
	}
	return result
}

// convertDirectiveArgsToJSON converts directive arguments to JSON format
func convertDirectiveArgsToJSON(args map[string]string) []JSONArgValue {
	if len(args) == 0 {
		return nil
	}
	result := make([]JSONArgValue, 0, len(args))
	for name, value := range args {
		result = append(result, JSONArgValue{
			Name:  name,
			Value: value,
		})
	}
	return result
}

// convertEnumValuesToJSON converts a slice of gql.EnumValue to JSON format
func convertEnumValuesToJSON(values []*gql.EnumValue) []JSONEnumValue {
	if len(values) == 0 {
		return nil
	}
	result := make([]JSONEnumValue, 0, len(values))
	for _, val := range values {
		result = append(result, JSONEnumValue{
			Name:        val.Name(),
			Description: val.Description(),
			Directives:  convertDirectivesToJSON(val.Directives()),
		})
	}
	return result
}

// convertDirectiveToJSON converts a gql.DirectiveDef to JSON format
func convertDirectiveToJSON(directive *gql.DirectiveDef) JSONDirective {
	return JSONDirective{
		Name:        directive.Name(),
		Description: directive.Description(),
		Locations:   directive.Locations(),
		Arguments:   convertArgumentsToJSON(directive.Arguments()),
	}
}

// convertTypeDefToJSON converts a gql.TypeDef to JSON format
func convertTypeDefToJSON(typeDef gql.TypeDef) JSONTypeDef {
	result := JSONTypeDef{
		Name:        typeDef.Name(),
		Description: typeDef.Description(),
	}

	// Add type-specific details and determine kind
	switch t := typeDef.(type) {
	case *gql.Object:
		result.Kind = "Object"
		result.Fields = convertFieldsToJSON(t.Fields())
		result.Interfaces = t.Interfaces()
		result.Directives = convertDirectivesToJSON(t.Directives())
	case *gql.Interface:
		result.Kind = "Interface"
		result.Fields = convertFieldsToJSON(t.Fields())
		result.Directives = convertDirectivesToJSON(t.Directives())
	case *gql.InputObject:
		result.Kind = "Input"
		result.Fields = convertFieldsToJSON(t.Fields())
		result.Directives = convertDirectivesToJSON(t.Directives())
	case *gql.Enum:
		result.Kind = "Enum"
		result.Values = convertEnumValuesToJSON(t.Values())
		result.Directives = convertDirectivesToJSON(t.Directives())
	case *gql.Union:
		result.Kind = "Union"
		result.Types = t.Types()
		result.Directives = convertDirectivesToJSON(t.Directives())
	case *gql.Scalar:
		result.Kind = "Scalar"
		result.Directives = convertDirectivesToJSON(t.Directives())
	}

	return result
}

// convertFieldsToJSON converts a slice of gql.Field to JSON format
func convertFieldsToJSON(fields []*gql.Field) []JSONField {
	if len(fields) == 0 {
		return nil
	}
	result := make([]JSONField, 0, len(fields))
	for _, field := range fields {
		result = append(result, convertFieldToJSON(field))
	}
	return result
}

// marshalJSON marshals an object to pretty-printed JSON with 2-space indentation
func marshalJSON(v interface{}) string {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "failed to marshal JSON: %v"}`, err)
	}
	return string(bytes)
}
