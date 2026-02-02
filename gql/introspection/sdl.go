package introspection

import (
	"fmt"
	"sort"
	"strings"
)

// builtInTypes are types that should be excluded from SDL output.
var builtInTypes = map[string]bool{
	"String":              true,
	"Int":                 true,
	"Float":               true,
	"Boolean":             true,
	"ID":                  true,
	"__Schema":            true,
	"__Type":              true,
	"__TypeKind":          true,
	"__Field":             true,
	"__InputValue":        true,
	"__EnumValue":         true,
	"__Directive":         true,
	"__DirectiveLocation": true,
}

// builtInDirectives are directives that should be excluded from SDL output.
var builtInDirectives = map[string]bool{
	"skip":        true,
	"include":     true,
	"deprecated":  true,
	"specifiedBy": true,
}

// ToSDL converts an introspection response to SDL format.
func ToSDL(resp *Response) ([]byte, error) {
	if resp == nil || resp.Data == nil {
		return nil, fmt.Errorf("invalid introspection response")
	}

	var sb strings.Builder

	// Write schema definition if non-standard root types
	writeSchemaDefinition(&sb, &resp.Data.Schema)

	// Sort types by name for consistent output
	types := make([]FullType, 0, len(resp.Data.Schema.Types))
	for _, t := range resp.Data.Schema.Types {
		if !builtInTypes[t.Name] {
			types = append(types, t)
		}
	}
	sort.Slice(types, func(i, j int) bool {
		return types[i].Name < types[j].Name
	})

	// Write types
	for _, t := range types {
		writeType(&sb, &t)
	}

	// Write directives
	for _, d := range resp.Data.Schema.Directives {
		if !builtInDirectives[d.Name] {
			writeDirective(&sb, &d)
		}
	}

	return []byte(sb.String()), nil
}

func writeSchemaDefinition(sb *strings.Builder, schema *Schema) {
	// Only write schema definition if root types have non-standard names
	queryName := ""
	mutationName := ""
	subscriptionName := ""

	if schema.QueryType != nil {
		queryName = schema.QueryType.Name
	}
	if schema.MutationType != nil {
		mutationName = schema.MutationType.Name
	}
	if schema.SubscriptionType != nil {
		subscriptionName = schema.SubscriptionType.Name
	}

	// Check if all root types have standard names
	if queryName == "Query" || queryName == "" {
		if mutationName == "Mutation" || mutationName == "" {
			if subscriptionName == "Subscription" || subscriptionName == "" {
				return
			}
		}
	}

	sb.WriteString("schema {\n")
	if queryName != "" {
		sb.WriteString(fmt.Sprintf("  query: %s\n", queryName))
	}
	if mutationName != "" {
		sb.WriteString(fmt.Sprintf("  mutation: %s\n", mutationName))
	}
	if subscriptionName != "" {
		sb.WriteString(fmt.Sprintf("  subscription: %s\n", subscriptionName))
	}
	sb.WriteString("}\n\n")
}

func writeType(sb *strings.Builder, t *FullType) {
	switch t.Kind {
	case "OBJECT":
		writeObjectType(sb, t)
	case "INTERFACE":
		writeInterfaceType(sb, t)
	case "UNION":
		writeUnionType(sb, t)
	case "ENUM":
		writeEnumType(sb, t)
	case "INPUT_OBJECT":
		writeInputObjectType(sb, t)
	case "SCALAR":
		writeScalarType(sb, t)
	}
}

func writeDescription(sb *strings.Builder, desc *string, indent string) {
	if desc == nil || *desc == "" {
		return
	}
	// Use block string for multi-line descriptions
	if strings.Contains(*desc, "\n") {
		sb.WriteString(fmt.Sprintf("%s\"\"\"\n", indent))
		for _, line := range strings.Split(*desc, "\n") {
			sb.WriteString(fmt.Sprintf("%s%s\n", indent, line))
		}
		sb.WriteString(fmt.Sprintf("%s\"\"\"\n", indent))
	} else {
		sb.WriteString(fmt.Sprintf("%s\"%s\"\n", indent, escapeString(*desc)))
	}
}

func escapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}

func writeObjectType(sb *strings.Builder, t *FullType) {
	writeDescription(sb, t.Description, "")
	sb.WriteString(fmt.Sprintf("type %s", t.Name))

	if len(t.Interfaces) > 0 {
		var ifaces []string
		for _, iface := range t.Interfaces {
			if iface.Name != nil {
				ifaces = append(ifaces, *iface.Name)
			}
		}
		if len(ifaces) > 0 {
			sb.WriteString(" implements ")
			sb.WriteString(strings.Join(ifaces, " & "))
		}
	}

	sb.WriteString(" {\n")
	for _, f := range t.Fields {
		writeField(sb, &f)
	}
	sb.WriteString("}\n\n")
}

func writeInterfaceType(sb *strings.Builder, t *FullType) {
	writeDescription(sb, t.Description, "")
	sb.WriteString(fmt.Sprintf("interface %s", t.Name))

	if len(t.Interfaces) > 0 {
		var ifaces []string
		for _, iface := range t.Interfaces {
			if iface.Name != nil {
				ifaces = append(ifaces, *iface.Name)
			}
		}
		if len(ifaces) > 0 {
			sb.WriteString(" implements ")
			sb.WriteString(strings.Join(ifaces, " & "))
		}
	}

	sb.WriteString(" {\n")
	for _, f := range t.Fields {
		writeField(sb, &f)
	}
	sb.WriteString("}\n\n")
}

func writeUnionType(sb *strings.Builder, t *FullType) {
	writeDescription(sb, t.Description, "")
	sb.WriteString(fmt.Sprintf("union %s = ", t.Name))

	var members []string
	for _, pt := range t.PossibleTypes {
		if pt.Name != nil {
			members = append(members, *pt.Name)
		}
	}
	sb.WriteString(strings.Join(members, " | "))
	sb.WriteString("\n\n")
}

func writeEnumType(sb *strings.Builder, t *FullType) {
	writeDescription(sb, t.Description, "")
	sb.WriteString(fmt.Sprintf("enum %s {\n", t.Name))

	for _, ev := range t.EnumValues {
		writeDescription(sb, ev.Description, "  ")
		sb.WriteString(fmt.Sprintf("  %s", ev.Name))
		if ev.IsDeprecated {
			writeDeprecation(sb, ev.DeprecationReason)
		}
		sb.WriteString("\n")
	}

	sb.WriteString("}\n\n")
}

func writeInputObjectType(sb *strings.Builder, t *FullType) {
	writeDescription(sb, t.Description, "")
	sb.WriteString(fmt.Sprintf("input %s {\n", t.Name))

	for _, f := range t.InputFields {
		writeInputValue(sb, &f, "  ", false)
		sb.WriteString("\n")
	}

	sb.WriteString("}\n\n")
}

func writeScalarType(sb *strings.Builder, t *FullType) {
	writeDescription(sb, t.Description, "")
	sb.WriteString(fmt.Sprintf("scalar %s\n\n", t.Name))
}

func writeField(sb *strings.Builder, f *Field) {
	writeDescription(sb, f.Description, "  ")
	sb.WriteString(fmt.Sprintf("  %s", f.Name))

	if len(f.Args) > 0 {
		sb.WriteString("(")
		for i, arg := range f.Args {
			if i > 0 {
				sb.WriteString(", ")
			}
			writeInputValue(sb, &arg, "", true)
		}
		sb.WriteString(")")
	}

	sb.WriteString(": ")
	sb.WriteString(formatTypeRef(&f.Type))

	if f.IsDeprecated {
		writeDeprecation(sb, f.DeprecationReason)
	}

	sb.WriteString("\n")
}

func writeInputValue(sb *strings.Builder, iv *InputValue, indent string, inline bool) {
	if !inline {
		writeDescription(sb, iv.Description, indent)
	}
	sb.WriteString(fmt.Sprintf("%s%s: %s", indent, iv.Name, formatTypeRef(&iv.Type)))

	if iv.DefaultValue != nil {
		sb.WriteString(fmt.Sprintf(" = %s", *iv.DefaultValue))
	}
}

func writeDeprecation(sb *strings.Builder, reason *string) {
	if reason != nil && *reason != "" && *reason != "No longer supported" {
		sb.WriteString(fmt.Sprintf(" @deprecated(reason: \"%s\")", escapeString(*reason)))
	} else {
		sb.WriteString(" @deprecated")
	}
}

func writeDirective(sb *strings.Builder, d *Directive) {
	writeDescription(sb, d.Description, "")
	sb.WriteString(fmt.Sprintf("directive @%s", d.Name))

	if len(d.Args) > 0 {
		sb.WriteString("(")
		for i, arg := range d.Args {
			if i > 0 {
				sb.WriteString(", ")
			}
			writeInputValue(sb, &arg, "", true)
		}
		sb.WriteString(")")
	}

	if len(d.Locations) > 0 {
		sb.WriteString(" on ")
		sb.WriteString(strings.Join(d.Locations, " | "))
	}

	sb.WriteString("\n\n")
}

func formatTypeRef(tr *TypeRef) string {
	if tr == nil {
		return ""
	}

	switch tr.Kind {
	case "NON_NULL":
		return formatTypeRef(tr.OfType) + "!"
	case "LIST":
		return "[" + formatTypeRef(tr.OfType) + "]"
	default:
		if tr.Name != nil {
			return *tr.Name
		}
		return ""
	}
}
