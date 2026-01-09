package gqlfmt

import (
	"fmt"
	"strings"

	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/utils/text"
)

// GenerateMarkdown generates markdown content for a given type name in the schema
func GenerateMarkdown(schema gql.GraphQLSchema, typeName string) (string, error) {
	resolver := gql.NewSchemaResolver(&schema)

	// Handle Query fields (Query.fieldName)
	if strings.HasPrefix(typeName, "Query.") {
		return generateQueryFieldMarkdown(schema, typeName, resolver)
	}

	// Handle Mutation fields (Mutation.fieldName)
	if strings.HasPrefix(typeName, "Mutation.") {
		return generateMutationFieldMarkdown(schema, typeName, resolver)
	}

	// Handle Directives (@directiveName)
	if strings.HasPrefix(typeName, "@") {
		return generateDirectiveMarkdown(schema, typeName, resolver)
	}

	// Handle regular types (Object, Input, Enum, Scalar, Interface, Union)
	return generateTypeMarkdown(schema, typeName, resolver)
}

// generateQueryFieldMarkdown generates markdown for a query field
func generateQueryFieldMarkdown(schema gql.GraphQLSchema, typeName string, resolver gql.TypeResolver) (string, error) {
	fieldName := strings.TrimPrefix(typeName, "Query.")
	field, ok := schema.Query[fieldName]
	if !ok {
		return "", fmt.Errorf("query field %q not found in schema", fieldName)
	}
	return GenerateFieldMarkdown(field, resolver), nil
}

// generateMutationFieldMarkdown generates markdown for a mutation field
func generateMutationFieldMarkdown(schema gql.GraphQLSchema, typeName string, resolver gql.TypeResolver) (string, error) {
	fieldName := strings.TrimPrefix(typeName, "Mutation.")
	field, ok := schema.Mutation[fieldName]
	if !ok {
		return "", fmt.Errorf("mutation field %q not found in schema", fieldName)
	}
	return GenerateFieldMarkdown(field, resolver), nil
}

// generateDirectiveMarkdown generates markdown for a directive
func generateDirectiveMarkdown(schema gql.GraphQLSchema, typeName string, resolver gql.TypeResolver) (string, error) {
	directiveName := strings.TrimPrefix(typeName, "@")
	directive, ok := schema.Directive[directiveName]
	if !ok {
		return "", fmt.Errorf("directive %q not found in schema", directiveName)
	}
	return GenerateDirectiveMarkdown(directive, resolver), nil
}

// generateTypeMarkdown generates markdown for a type definition
func generateTypeMarkdown(schema gql.GraphQLSchema, typeName string, resolver gql.TypeResolver) (string, error) {
	typeDef, err := schema.NamedToTypeDef(typeName)
	if err == nil {
		return GenerateTypeDefMarkdown(typeDef, resolver), nil
	}

	// If type not found and typeName contains a dot, try trimming the field part
	if !strings.Contains(typeName, ".") {
		return "", fmt.Errorf("type %q not found in schema: %w", typeName, err)
	}

	baseTypeName := typeName[:strings.LastIndex(typeName, ".")]
	typeDef, retryErr := schema.NamedToTypeDef(baseTypeName)
	if retryErr != nil {
		return "", fmt.Errorf("type %q not found in schema: %w", typeName, err)
	}

	return GenerateTypeDefMarkdown(typeDef, resolver), nil
}

// GenerateFieldMarkdown generates markdown for a GraphQL field
func GenerateFieldMarkdown(field *gql.Field, resolver gql.TypeResolver) string {
	parts := []string{
		text.H1(field.Name()),
		text.GqlCode(field.FormatSignature(80)),
		field.Description(),
	}
	return text.JoinParagraphs(parts...)
}

// GenerateDirectiveMarkdown generates markdown for a GraphQL directive
func GenerateDirectiveMarkdown(directive *gql.Directive, resolver gql.TypeResolver) string {
	parts := []string{
		text.H1("@" + directive.Name()),
		text.GqlCode(directive.FormatSignature(80)),
		directive.Description(),
	}
	if len(directive.Locations()) > 0 {
		locationList := []string{}
		for _, loc := range directive.Locations() {
			locationList = append(locationList, "- "+loc)
		}
		parts = append(parts, "**Locations:**\n"+text.JoinLines(locationList...))
	}
	return text.JoinParagraphs(parts...)
}

// GenerateTypeDefMarkdown generates markdown for a GraphQL type definition
func GenerateTypeDefMarkdown(typeDef gql.TypeDef, resolver gql.TypeResolver) string {
	parts := []string{text.H1(typeDef.Name())}

	// Add description if available
	if desc := typeDef.Description(); desc != "" {
		parts = append(parts, desc)
	}

	// Add type-specific details
	switch t := typeDef.(type) {
	case *gql.Object:
		if len(t.Interfaces()) > 0 {
			parts = append(parts, "**Implements:** "+strings.Join(t.Interfaces(), ", "))
		}
		fieldsWithDesc := FormatFieldDefinitionsWithDescriptions(t.Fields())
		if len(fieldsWithDesc) > 0 {
			parts = append(parts, fieldsWithDesc)
		}
	case *gql.Scalar:
		parts = append(parts, "_Scalar type_")
	case *gql.Interface:
		fieldsWithDesc := FormatFieldDefinitionsWithDescriptions(t.Fields())
		if len(fieldsWithDesc) > 0 {
			parts = append(parts, fieldsWithDesc)
		}
	case *gql.Union:
		if len(t.Types()) > 0 {
			parts = append(parts, "**Union of:** "+strings.Join(t.Types(), " | "))
		}
	case *gql.Enum:
		valuesWithDesc := FormatEnumValuesWithDescriptions(t.Values())
		if len(valuesWithDesc) > 0 {
			parts = append(parts, valuesWithDesc)
		}
	case *gql.InputObject:
		fieldsWithDesc := FormatFieldDefinitionsWithDescriptions(t.Fields())
		if len(fieldsWithDesc) > 0 {
			parts = append(parts, fieldsWithDesc)
		}
	}

	return text.JoinParagraphs(parts...)
}

// FormatFieldDefinitionsWithDescriptions formats field definitions with their descriptions
func FormatFieldDefinitionsWithDescriptions(fieldNodes []*gql.Field) string {
	if len(fieldNodes) == 0 {
		return ""
	}
	var parts []string
	for _, field := range fieldNodes {
		fieldParts := []string{}
		if desc := field.Description(); desc != "" {
			fieldParts = append(fieldParts, text.GqlDocString(desc))
		}
		fieldParts = append(fieldParts, field.Signature())
		parts = append(parts, text.JoinLines(fieldParts...))
	}
	return text.GqlCode(text.JoinParagraphs(parts...))
}

// FormatEnumValuesWithDescriptions formats enum values with their descriptions
func FormatEnumValuesWithDescriptions(enumValues []*gql.EnumValue) string {
	if len(enumValues) == 0 {
		return ""
	}
	var parts []string
	for _, val := range enumValues {
		valParts := []string{}
		if desc := val.Description(); desc != "" {
			valParts = append(valParts, text.GqlDocString(desc))
		}
		valParts = append(valParts, val.Name())
		parts = append(parts, text.JoinLines(valParts...))
	}
	return text.GqlCode(text.JoinParagraphs(parts...))
}
