package adapters

import (
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
	return AdaptFields(gql.CollectAndSortMapValues(p.schema.Query), p.resolver)
}

func (p *SchemaView) GetMutationItems() []components.ListItem {
	return AdaptFields(gql.CollectAndSortMapValues(p.schema.Mutation), p.resolver)
}

func (p *SchemaView) GetObjectItems() []components.ListItem {
	return AdaptTypeDefs(gql.CollectAndSortMapValues(p.schema.Object), p.resolver)
}

func (p *SchemaView) GetInputItems() []components.ListItem {
	return AdaptTypeDefs(gql.CollectAndSortMapValues(p.schema.Input), p.resolver)
}

func (p *SchemaView) GetEnumItems() []components.ListItem {
	return AdaptTypeDefs(gql.CollectAndSortMapValues(p.schema.Enum), p.resolver)
}

func (p *SchemaView) GetScalarItems() []components.ListItem {
	return AdaptTypeDefs(gql.CollectAndSortMapValues(p.schema.Scalar), p.resolver)
}

func (p *SchemaView) GetInterfaceItems() []components.ListItem {
	return AdaptTypeDefs(gql.CollectAndSortMapValues(p.schema.Interface), p.resolver)
}

func (p *SchemaView) GetUnionItems() []components.ListItem {
	return AdaptTypeDefs(gql.CollectAndSortMapValues(p.schema.Union), p.resolver)
}

func (p *SchemaView) GetDirectiveItems() []components.ListItem {
	return AdaptDirectives(gql.CollectAndSortMapValues(p.schema.Directive), p.resolver)
}

// FindTypeCategory returns the GQL type category for the given type name
// Returns (category, true) if found, ("", false) if not found
func (p *SchemaView) FindTypeCategory(typeName string) (navigation.GQLType, bool) {
	if category, ok := p.schema.NameToType[typeName]; ok {
		return navigation.GQLType(category), true
	}
	return "", false
}
