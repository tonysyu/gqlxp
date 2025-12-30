package xplr

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/tui/adapters"
	"github.com/tonysyu/gqlxp/tui/xplr/navigation"
)

const testSchema = `
	type Query {
		user: User
	}

	type User {
		id: ID!
		name: String!
	}
`

func TestApplySelection_TypeOnly(t *testing.T) {
	is := is.New(t)
	schema, err := adapters.ParseSchemaString(testSchema)
	is.NoErr(err)

	m := New(schema)
	m, _ = m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	target := SelectionTarget{
		TypeName: "User",
	}
	m.ApplySelection(target)

	is.Equal(m.CurrentType(), string(navigation.ObjectType))
}

func TestApplySelection_FieldSelection(t *testing.T) {
	is := is.New(t)
	schema, err := adapters.ParseSchemaString(testSchema)
	is.NoErr(err)

	m := New(schema)
	m, _ = m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	target := SelectionTarget{
		TypeName:  "Query",
		FieldName: "user",
	}

	// Debug: Check stack before
	t.Logf("Stack length before: %d, position: %d", m.nav.Stack().Len(), m.nav.Stack().Position())

	m.ApplySelection(target)

	// Debug: Check stack after
	t.Logf("Stack length after: %d, position: %d", m.nav.Stack().Len(), m.nav.Stack().Position())

	is.Equal(m.CurrentType(), string(navigation.QueryType))

	// Check breadcrumbs
	breadcrumbs := m.nav.Breadcrumbs()
	t.Logf("Breadcrumbs: %v", breadcrumbs)
	if len(breadcrumbs) == 0 {
		t.Error("Expected breadcrumbs to be set, got empty")
	}
}
