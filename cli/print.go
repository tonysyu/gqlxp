package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/utils/text"
	"github.com/urfave/cli/v3"
)

// printCommand creates the print subcommand
func printCommand() *cli.Command {
	return &cli.Command{
		Name:      "print",
		Usage:     "Print a GraphQL type definition to the terminal",
		ArgsUsage: "<schema-file> <type-name>",
		Description: `Prints the details of a GraphQL type directly to the terminal.

The type-name can be:
- A Query field name (prefix with "Query.")
- A Mutation field name (prefix with "Mutation.")
- A type name (Object, Input, Enum, Scalar, Interface, Union)
- A directive name (prefix with "@")

Examples:
  gqlxp print schema.graphqls User
  gqlxp print schema.graphqls Query.getUser
  gqlxp print schema.graphqls Mutation.createUser
  gqlxp print schema.graphqls @auth`,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-pager",
				Usage: "disable pager and print directly to stdout",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 2 {
				return fmt.Errorf("requires exactly 2 arguments: <schema-file> <type-name>")
			}

			schemaFile := cmd.Args().Get(0)
			typeName := cmd.Args().Get(1)
			noPager := cmd.Bool("no-pager")

			return printType(schemaFile, typeName, noPager)
		},
	}
}

func printType(schemaFile, typeName string, noPager bool) error {
	// Load schema from file
	content, err := loadSchemaFromFile(schemaFile)
	if err != nil {
		return err
	}

	// Parse schema
	schema, err := gql.ParseSchema(content)
	if err != nil {
		return fmt.Errorf("error parsing schema: %w", err)
	}

	// Generate markdown content based on type name
	markdown, err := generateMarkdown(schema, typeName)
	if err != nil {
		return err
	}

	// Render markdown using glamour
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
	)
	if err != nil {
		// Fallback to plain text if glamour fails
		fmt.Println(markdown)
		return nil
	}

	rendered, err := renderer.Render(markdown)
	if err != nil {
		// Fallback to plain text if rendering fails
		fmt.Println(markdown)
		return nil
	}

	// Use pager if content is long enough and not disabled
	if shouldUsePager(rendered, noPager) {
		return showInPager(rendered)
	}

	fmt.Print(rendered)
	return nil
}

func generateMarkdown(schema gql.GraphQLSchema, typeName string) (string, error) {
	resolver := gql.NewSchemaResolver(&schema)

	// Handle Query fields (Query.fieldName)
	if strings.HasPrefix(typeName, "Query.") {
		fieldName := strings.TrimPrefix(typeName, "Query.")
		field, ok := schema.Query[fieldName]
		if !ok {
			return "", fmt.Errorf("query field %q not found in schema", fieldName)
		}
		return generateFieldMarkdown(field, resolver), nil
	}

	// Handle Mutation fields (Mutation.fieldName)
	if strings.HasPrefix(typeName, "Mutation.") {
		fieldName := strings.TrimPrefix(typeName, "Mutation.")
		field, ok := schema.Mutation[fieldName]
		if !ok {
			return "", fmt.Errorf("mutation field %q not found in schema", fieldName)
		}
		return generateFieldMarkdown(field, resolver), nil
	}

	// Handle Directives (@directiveName)
	if strings.HasPrefix(typeName, "@") {
		directiveName := strings.TrimPrefix(typeName, "@")
		directive, ok := schema.Directive[directiveName]
		if !ok {
			return "", fmt.Errorf("directive %q not found in schema", directiveName)
		}
		return generateDirectiveMarkdown(directive, resolver), nil
	}

	// Handle regular types (Object, Input, Enum, Scalar, Interface, Union)
	typeDef, err := schema.NamedToTypeDef(typeName)
	if err != nil {
		return "", fmt.Errorf("type %q not found in schema: %w", typeName, err)
	}

	return generateTypeDefMarkdown(typeDef, resolver), nil
}

func generateFieldMarkdown(field *gql.Field, resolver gql.TypeResolver) string {
	parts := []string{
		text.H1(field.Name()),
		text.GqlCode(field.FormatSignature(80)),
		field.Description(),
	}
	return text.JoinParagraphs(parts...)
}

func generateDirectiveMarkdown(directive *gql.Directive, resolver gql.TypeResolver) string {
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

func generateTypeDefMarkdown(typeDef gql.TypeDef, resolver gql.TypeResolver) string {
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
		fieldsWithDesc := formatFieldDefinitionsWithDescriptions(t.Fields())
		if len(fieldsWithDesc) > 0 {
			parts = append(parts, fieldsWithDesc)
		}
	case *gql.Scalar:
		parts = append(parts, "_Scalar type_")
	case *gql.Interface:
		fieldsWithDesc := formatFieldDefinitionsWithDescriptions(t.Fields())
		if len(fieldsWithDesc) > 0 {
			parts = append(parts, fieldsWithDesc)
		}
	case *gql.Union:
		if len(t.Types()) > 0 {
			parts = append(parts, "**Union of:** "+strings.Join(t.Types(), " | "))
		}
	case *gql.Enum:
		valuesWithDesc := formatEnumValuesWithDescriptions(t.Values())
		if len(valuesWithDesc) > 0 {
			parts = append(parts, valuesWithDesc)
		}
	case *gql.InputObject:
		fieldsWithDesc := formatFieldDefinitionsWithDescriptions(t.Fields())
		if len(fieldsWithDesc) > 0 {
			parts = append(parts, fieldsWithDesc)
		}
	}

	return text.JoinParagraphs(parts...)
}

func formatFieldDefinitionsWithDescriptions(fieldNodes []*gql.Field) string {
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

func formatEnumValuesWithDescriptions(enumValues []*gql.EnumValue) string {
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
