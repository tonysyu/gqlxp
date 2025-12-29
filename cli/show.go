package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/library"
	"github.com/tonysyu/gqlxp/utils/terminal"
	"github.com/tonysyu/gqlxp/utils/text"
	"github.com/urfave/cli/v3"
)

// showCommand creates the show subcommand
func showCommand() *cli.Command {
	return &cli.Command{
		Name:      "show",
		Usage:     "Show a GraphQL type definition in the terminal",
		ArgsUsage: "[schema-file] <type-name>",
		Description: `Shows the details of a GraphQL type directly to the terminal.

The schema-file argument is optional if a default schema has been set.
Use 'gqlxp library default' to set the default schema.

The type-name can be:
- A Query field name (prefix with "Query.")
- A Mutation field name (prefix with "Mutation.")
- A type name (Object, Input, Enum, Scalar, Interface, Union)
- A directive name (prefix with "@")

Examples:
  gqlxp show examples/github.graphqls User
  gqlxp show User  # Uses default schema
  gqlxp show Query.getUser
  gqlxp show Mutation.createUser
  gqlxp show @auth`,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-pager",
				Usage: "disable pager and show directly to stdout",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() < 1 || cmd.Args().Len() > 2 {
				return fmt.Errorf("requires 1 or 2 arguments: [schema-file] <type-name>")
			}

			var schemaArg, typeName string
			noPager := cmd.Bool("no-pager")

			// Parse arguments
			if cmd.Args().Len() == 2 {
				// Two arguments: schema and type
				schemaArg = cmd.Args().Get(0)
				typeName = cmd.Args().Get(1)
			} else {
				// One argument: type only (use default schema)
				typeName = cmd.Args().Get(0)
			}

			return printType(schemaArg, typeName, noPager)
		},
	}
}

func printType(schemaArg, typeName string, noPager bool) error {
	lib := library.NewLibrary()
	var content []byte

	// Determine which schema to use
	var schemaID string
	if schemaArg == "" {
		// No schema specified - use default schema
		defaultSchemaID, err := lib.GetDefaultSchema()
		if err != nil {
			return fmt.Errorf("error getting default schema: %w", err)
		}
		if defaultSchemaID == "" {
			return fmt.Errorf("no schema specified and no default schema set. Use 'gqlxp library default' to set one")
		}
		schemaID = defaultSchemaID
	} else {
		// Schema specified - resolve as either ID or file path
		resolvedID, err := resolveSchemaArgument(schemaArg)
		if err != nil {
			return err
		}
		schemaID = resolvedID
	}

	// Load schema content from library
	libSchema, err := lib.Get(schemaID)
	if err != nil {
		return fmt.Errorf("error loading schema: %w", err)
	}
	content = libSchema.Content

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

	// Render markdown using terminal renderer
	renderer, err := terminal.NewMarkdownRenderer()
	if err != nil {
		// Fallback to plain text if renderer fails
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
	if terminal.ShouldUsePager(rendered, noPager) {
		return terminal.ShowInPager(rendered)
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
