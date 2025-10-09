package adapters

import (
	"github.com/tonysyu/igq/gql"
	"github.com/tonysyu/igq/tui/components"
)

type SchemaView struct {
	schema gql.GraphQLSchema
}

func NewSchemaView(schema gql.GraphQLSchema) SchemaView {
	return SchemaView{
		schema: schema,
	}
}

func (p *SchemaView) GetQueryItems() []components.ListItem {
	return AdaptFieldDefinitionsToItems(gql.CollectAndSortMapValues(p.schema.Query), &p.schema)
}

func (p *SchemaView) GetMutationItems() []components.ListItem {
	return AdaptFieldDefinitionsToItems(gql.CollectAndSortMapValues(p.schema.Mutation), &p.schema)
}

func (p *SchemaView) GetObjectItems() []components.ListItem {
	return AdaptObjectDefinitionsToItems(gql.CollectAndSortMapValues(p.schema.Object), &p.schema)
}

func (p *SchemaView) GetInputItems() []components.ListItem {
	return AdaptInputDefinitionsToItems(gql.CollectAndSortMapValues(p.schema.Input), &p.schema)
}

func (p *SchemaView) GetEnumItems() []components.ListItem {
	return AdaptEnumDefinitionsToItems(gql.CollectAndSortMapValues(p.schema.Enum), &p.schema)
}

func (p *SchemaView) GetScalarItems() []components.ListItem {
	return AdaptScalarDefinitionsToItems(gql.CollectAndSortMapValues(p.schema.Scalar), &p.schema)
}

func (p *SchemaView) GetInterfaceItems() []components.ListItem {
	return AdaptInterfaceDefinitionsToItems(gql.CollectAndSortMapValues(p.schema.Interface), &p.schema)
}

func (p *SchemaView) GetUnionItems() []components.ListItem {
	return AdaptUnionDefinitionsToItems(gql.CollectAndSortMapValues(p.schema.Union), &p.schema)
}

func (p *SchemaView) GetDirectiveItems() []components.ListItem {
	return AdaptDirectiveDefinitionsToItems(gql.CollectAndSortMapValues(p.schema.Directive))
}

// GetSchema returns the underlying GraphQL schema (primarily for testing)
func (p *SchemaView) GetSchema() *gql.GraphQLSchema {
	return &p.schema
}
