package adapters

import (
	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/tui/components"
)

type SchemaView struct {
	schema gql.GraphQLSchema
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
		schema: schema,
	}
}

func (p *SchemaView) GetQueryItems() []components.ListItem {
	return adaptFieldsToItems(gql.CollectAndSortMapValues(p.schema.Query), &p.schema)
}

func (p *SchemaView) GetMutationItems() []components.ListItem {
	return adaptFieldsToItems(gql.CollectAndSortMapValues(p.schema.Mutation), &p.schema)
}

func (p *SchemaView) GetObjectItems() []components.ListItem {
	return adaptObjectsToItems(gql.CollectAndSortMapValues(p.schema.Object), &p.schema)
}

func (p *SchemaView) GetInputItems() []components.ListItem {
	return adaptInputObjectsToItems(gql.CollectAndSortMapValues(p.schema.Input), &p.schema)
}

func (p *SchemaView) GetEnumItems() []components.ListItem {
	return adaptEnumsToItems(gql.CollectAndSortMapValues(p.schema.Enum), &p.schema)
}

func (p *SchemaView) GetScalarItems() []components.ListItem {
	return adaptScalarsToItems(gql.CollectAndSortMapValues(p.schema.Scalar), &p.schema)
}

func (p *SchemaView) GetInterfaceItems() []components.ListItem {
	return adaptInterfacesToItems(gql.CollectAndSortMapValues(p.schema.Interface), &p.schema)
}

func (p *SchemaView) GetUnionItems() []components.ListItem {
	return adaptUnionsToItems(gql.CollectAndSortMapValues(p.schema.Union), &p.schema)
}

func (p *SchemaView) GetDirectiveItems() []components.ListItem {
	return adaptDirectivesToItems(gql.CollectAndSortMapValues(p.schema.Directive), &p.schema)
}
