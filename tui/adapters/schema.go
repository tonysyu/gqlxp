package adapters

import (
	"fmt"

	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/tui/xplr/components"
	"github.com/tonysyu/gqlxp/tui/xplr/navigation"
)

type SchemaView struct {
	schema   gql.GraphQLSchema
	resolver gql.TypeResolver
}

func ParseSchemaString(schemaContent string) (SchemaView, error) {
	return ParseSchema([]byte(schemaContent))
}

func ParseSchema(schemaContent []byte) (SchemaView, error) {
	schema, err := gql.ParseSchema(schemaContent)
	if err != nil {
		return SchemaView{}, err
	}
	return NewSchemaView(schema), nil
}

func NewSchemaView(schema gql.GraphQLSchema) SchemaView {
	return SchemaView{
		schema:   schema,
		resolver: gql.NewSchemaResolver(&schema),
	}
}

func (p *SchemaView) GetQueryItems() []components.ListItem {
	return adaptFields(gql.CollectAndSortMapValues(p.schema.Query), p.resolver)
}

func (p *SchemaView) GetMutationItems() []components.ListItem {
	return adaptFields(gql.CollectAndSortMapValues(p.schema.Mutation), p.resolver)
}

func (p *SchemaView) GetObjectItems() []components.ListItem {
	return adaptTypeDefs(gql.CollectAndSortMapValues(p.schema.Object), p.resolver)
}

func (p *SchemaView) GetInputItems() []components.ListItem {
	return adaptTypeDefs(gql.CollectAndSortMapValues(p.schema.Input), p.resolver)
}

func (p *SchemaView) GetEnumItems() []components.ListItem {
	return adaptTypeDefs(gql.CollectAndSortMapValues(p.schema.Enum), p.resolver)
}

func (p *SchemaView) GetScalarItems() []components.ListItem {
	return adaptTypeDefs(gql.CollectAndSortMapValues(p.schema.Scalar), p.resolver)
}

func (p *SchemaView) GetInterfaceItems() []components.ListItem {
	return adaptTypeDefs(gql.CollectAndSortMapValues(p.schema.Interface), p.resolver)
}

func (p *SchemaView) GetUnionItems() []components.ListItem {
	return adaptTypeDefs(gql.CollectAndSortMapValues(p.schema.Union), p.resolver)
}

func (p *SchemaView) GetDirectiveItems() []components.ListItem {
	return adaptDirectiveDefs(gql.CollectAndSortMapValues(p.schema.Directive), p.resolver)
}

// Schema returns the underlying GraphQL schema
func (p *SchemaView) Schema() *gql.GraphQLSchema {
	return &p.schema
}

// FindKind returns the GQL kind for the given type name.
// Returns (kind, true) if found, ("", false) if not found
func (p *SchemaView) FindKind(typeName string) (navigation.GQLKind, bool) {
	if category, ok := p.schema.NameToKind[typeName]; ok {
		return navigation.GQLKind(category), true
	}
	return "", false
}

// ResolveField looks up a field by search result kind, parent type name, and field name.
// Supported kinds: "Query", "Mutation", "ObjectField", "InputField", "InterfaceField"
func (p *SchemaView) ResolveField(kind, typeName, fieldName string) (*gql.Field, error) {
	switch kind {
	case "Query":
		if f, ok := p.schema.Query[fieldName]; ok {
			return f, nil
		}
		return nil, fmt.Errorf("query field %q not found", fieldName)
	case "Mutation":
		if f, ok := p.schema.Mutation[fieldName]; ok {
			return f, nil
		}
		return nil, fmt.Errorf("mutation field %q not found", fieldName)
	case "ObjectField":
		obj, ok := p.schema.Object[typeName]
		if !ok {
			return nil, fmt.Errorf("object type %q not found", typeName)
		}
		return findFieldByName(obj.Fields(), fieldName)
	case "InputField":
		inp, ok := p.schema.Input[typeName]
		if !ok {
			return nil, fmt.Errorf("input type %q not found", typeName)
		}
		return findFieldByName(inp.Fields(), fieldName)
	case "InterfaceField":
		iface, ok := p.schema.Interface[typeName]
		if !ok {
			return nil, fmt.Errorf("interface type %q not found", typeName)
		}
		return findFieldByName(iface.Fields(), fieldName)
	default:
		return nil, fmt.Errorf("kind %q does not have resolvable fields", kind)
	}
}

// findFieldByName searches for a field by name in a slice of fields.
func findFieldByName(fields []*gql.Field, name string) (*gql.Field, error) {
	for _, field := range fields {
		if field.Name() == name {
			return field, nil
		}
	}
	return nil, fmt.Errorf("field %q not found", name)
}
