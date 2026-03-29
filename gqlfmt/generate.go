package gqlfmt

import (
	"fmt"
	"strings"

	"github.com/tonysyu/gqlxp/gql"
)

// GenerateOptions controls operation generation behavior.
type GenerateOptions struct {
	Depth             int
	IncludeDeprecated bool
}

// GenerateOperation scaffolds a skeleton GraphQL operation for a Query or Mutation field.
func GenerateOperation(schema gql.GraphQLSchema, fieldPath string, opts GenerateOptions) (string, error) {
	var field *gql.Field
	var operationType string

	if strings.HasPrefix(fieldPath, "Query.") {
		fieldName := strings.TrimPrefix(fieldPath, "Query.")
		f, ok := schema.Query[fieldName]
		if !ok {
			return "", fmt.Errorf("query field %q not found", fieldName)
		}
		field = f
		operationType = "query"
	} else if strings.HasPrefix(fieldPath, "Mutation.") {
		fieldName := strings.TrimPrefix(fieldPath, "Mutation.")
		f, ok := schema.Mutation[fieldName]
		if !ok {
			return "", fmt.Errorf("mutation field %q not found", fieldName)
		}
		field = f
		operationType = "mutation"
	} else {
		return "", fmt.Errorf("field path must start with Query. or Mutation., got %q", fieldPath)
	}

	operationName := toPascalCase(field.Name())
	vars := collectVariables(field)
	selectionSet := buildSelectionSet(schema, field.ObjectTypeName(), opts.Depth, opts.IncludeDeprecated, "    ")

	return formatOperation(operationType, operationName, vars, field, selectionSet), nil
}

// toPascalCase converts camelCase to PascalCase.
func toPascalCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// collectVariables returns variable declarations for all non-null arguments.
func collectVariables(field *gql.Field) []string {
	var vars []string
	for _, arg := range field.Arguments() {
		if strings.HasSuffix(arg.TypeString(), "!") {
			vars = append(vars, fmt.Sprintf("$%s: %s", arg.Name(), arg.TypeString()))
		}
	}
	return vars
}

// isFieldDeprecated returns true if the field has an @deprecated directive.
func isFieldDeprecated(field *gql.Field) bool {
	for _, d := range field.Directives() {
		if d.Name() == "deprecated" {
			return true
		}
	}
	return false
}

// buildSelectionSet generates a { ... } selection set string for the given type.
// indent is the indentation for fields inside the selection set.
// Returns "" if the type has no selectable fields (e.g., scalars, enums).
func buildSelectionSet(schema gql.GraphQLSchema, typeName string, depth int, includeDeprecated bool, indent string) string {
	typeDef, err := schema.NamedToTypeDef(typeName)
	if err != nil {
		return "" // scalar or unknown type
	}

	closingIndent := strings.TrimSuffix(indent, "  ")

	switch t := typeDef.(type) {
	case *gql.Object:
		return buildFieldsSelection(schema, t.Fields(), depth, includeDeprecated, indent, closingIndent)
	case *gql.Interface:
		return buildFieldsSelection(schema, t.Fields(), depth, includeDeprecated, indent, closingIndent)
	case *gql.Union:
		return buildUnionSelection(schema, t, depth, includeDeprecated, indent, closingIndent)
	default:
		return "" // enum, scalar: no selection set
	}
}

// buildFieldsSelection builds a selection set from a list of fields.
func buildFieldsSelection(schema gql.GraphQLSchema, fields []*gql.Field, depth int, includeDeprecated bool, indent, closingIndent string) string {
	var lines []string
	for _, f := range fields {
		if !includeDeprecated && isFieldDeprecated(f) {
			continue
		}
		childTypeName := f.ObjectTypeName()
		childTypeDef, err := schema.NamedToTypeDef(childTypeName)
		if err != nil {
			// Scalar or unknown: include as leaf
			lines = append(lines, indent+f.Name())
			continue
		}
		switch childTypeDef.(type) {
		case *gql.Object, *gql.Interface, *gql.Union:
			if depth > 0 {
				childSelection := buildSelectionSet(schema, childTypeName, depth-1, includeDeprecated, indent+"  ")
				if childSelection != "" {
					lines = append(lines, fmt.Sprintf("%s%s %s", indent, f.Name(), childSelection))
				} else {
					lines = append(lines, indent+f.Name())
				}
			} else {
				lines = append(lines, fmt.Sprintf("%s# %s (%s)", indent, f.Name(), childTypeName))
			}
		default:
			// Enum: include as leaf
			lines = append(lines, indent+f.Name())
		}
	}
	if len(lines) == 0 {
		return ""
	}
	return "{\n" + strings.Join(lines, "\n") + "\n" + closingIndent + "}"
}

// buildUnionSelection builds inline fragments for each union member type.
func buildUnionSelection(schema gql.GraphQLSchema, union *gql.Union, depth int, includeDeprecated bool, indent, closingIndent string) string {
	var lines []string
	for _, memberType := range union.Types() {
		memberSelection := buildSelectionSet(schema, memberType, depth, includeDeprecated, indent+"  ")
		if memberSelection != "" {
			lines = append(lines, fmt.Sprintf("%s... on %s %s", indent, memberType, memberSelection))
		} else {
			lines = append(lines, fmt.Sprintf("%s... on %s { }", indent, memberType))
		}
	}
	if len(lines) == 0 {
		return ""
	}
	return "{\n" + strings.Join(lines, "\n") + "\n" + closingIndent + "}"
}

// formatOperation assembles the complete GraphQL operation string.
func formatOperation(operationType, operationName string, vars []string, field *gql.Field, selectionSet string) string {
	var header string
	if len(vars) > 0 {
		header = fmt.Sprintf("%s %s(%s)", operationType, operationName, strings.Join(vars, ", "))
	} else {
		header = fmt.Sprintf("%s %s", operationType, operationName)
	}

	fieldLine := buildFieldCall(field)

	if selectionSet != "" {
		return fmt.Sprintf("%s {\n  %s %s\n}", header, fieldLine, selectionSet)
	}
	return fmt.Sprintf("%s {\n  %s\n}", header, fieldLine)
}

// buildFieldCall formats the field invocation with its required argument references.
func buildFieldCall(field *gql.Field) string {
	var argParts []string
	for _, arg := range field.Arguments() {
		if strings.HasSuffix(arg.TypeString(), "!") {
			argParts = append(argParts, fmt.Sprintf("%s: $%s", arg.Name(), arg.Name()))
		}
	}
	if len(argParts) > 0 {
		return fmt.Sprintf("%s(%s)", field.Name(), strings.Join(argParts, ", "))
	}
	return field.Name()
}
