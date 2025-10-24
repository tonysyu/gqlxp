package adapters

import (
	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/tui/components"
)

type SchemaView struct {
	schema gql.GraphQLSchema
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
	return adaptFieldDefinitionsToItems(p.schema.GetSortedQueryFields(), &p.schema)
}

func (p *SchemaView) GetMutationItems() []components.ListItem {
	return adaptFieldDefinitionsToItems(p.schema.GetSortedMutationFields(), &p.schema)
}

func (p *SchemaView) GetObjectItems() []components.ListItem {
	return adaptObjectDefinitionsToItems(gql.CollectAndSortMapValues(p.schema.Object), &p.schema)
}

func (p *SchemaView) GetInputItems() []components.ListItem {
	return adaptInputDefinitionsToItems(gql.CollectAndSortMapValues(p.schema.Input), &p.schema)
}

func (p *SchemaView) GetEnumItems() []components.ListItem {
	return adaptEnumDefinitionsToItems(gql.CollectAndSortMapValues(p.schema.Enum), &p.schema)
}

func (p *SchemaView) GetScalarItems() []components.ListItem {
	return adaptScalarDefinitionsToItems(gql.CollectAndSortMapValues(p.schema.Scalar), &p.schema)
}

func (p *SchemaView) GetInterfaceItems() []components.ListItem {
	return adaptInterfaceDefinitionsToItems(gql.CollectAndSortMapValues(p.schema.Interface), &p.schema)
}

func (p *SchemaView) GetUnionItems() []components.ListItem {
	return adaptUnionDefinitionsToItems(gql.CollectAndSortMapValues(p.schema.Union), &p.schema)
}

func (p *SchemaView) GetDirectiveItems() []components.ListItem {
	return adaptDirectiveDefinitionsToItems(gql.CollectAndSortMapValues(p.schema.Directive))
}
