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
	return AdaptFieldDefinitionsToItems(p.schema.GetSortedQueryFields(), &p.schema)
}

func (p *SchemaView) GetMutationItems() []components.ListItem {
	return AdaptFieldDefinitionsToItems(p.schema.GetSortedMutationFields(), &p.schema)
}

func (p *SchemaView) GetObjectItems() []components.ListItem {
	return AdaptObjectDefinitionsToItems(gql.WrapObjectDefinitions(gql.CollectAndSortMapValues(p.schema.Object)), &p.schema)
}

func (p *SchemaView) GetInputItems() []components.ListItem {
	return AdaptInputDefinitionsToItems(gql.WrapInputObjectDefinitions(gql.CollectAndSortMapValues(p.schema.Input)), &p.schema)
}

func (p *SchemaView) GetEnumItems() []components.ListItem {
	return AdaptEnumDefinitionsToItems(gql.WrapEnumDefinitions(gql.CollectAndSortMapValues(p.schema.Enum)), &p.schema)
}

func (p *SchemaView) GetScalarItems() []components.ListItem {
	return AdaptScalarDefinitionsToItems(gql.WrapScalarDefinitions(gql.CollectAndSortMapValues(p.schema.Scalar)), &p.schema)
}

func (p *SchemaView) GetInterfaceItems() []components.ListItem {
	return AdaptInterfaceDefinitionsToItems(gql.WrapInterfaceDefinitions(gql.CollectAndSortMapValues(p.schema.Interface)), &p.schema)
}

func (p *SchemaView) GetUnionItems() []components.ListItem {
	return AdaptUnionDefinitionsToItems(gql.WrapUnionDefinitions(gql.CollectAndSortMapValues(p.schema.Union)), &p.schema)
}

func (p *SchemaView) GetDirectiveItems() []components.ListItem {
	return AdaptDirectiveDefinitionsToItems(gql.WrapDirectiveDefinitions(gql.CollectAndSortMapValues(p.schema.Directive)))
}
