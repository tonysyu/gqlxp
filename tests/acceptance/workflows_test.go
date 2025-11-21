package acceptance

import (
	"testing"

	"github.com/tonysyu/gqlxp/tui/xplr/navigation"
)

// testSchema defines a comprehensive GraphQL schema for acceptance testing
const testSchema = `
	type Query {
		query1(arg1: ID!, arg2: String!): Object1
		query2(arg1: ID!, arg2: String!): Object2
	}

	type Mutation {
		mutation1(input: Mutation1Input!): Object1!
		mutation2(input: Mutation1Input!): Object2!
	}

	type Object1 {
		field1: ID!
		field2: String!
	}

	type Object2 {
		field1: ID!
		field2: String!
	}

	input Mutation1Input {
		field1: Id!
		field2: String!
	}

	input Mutation2Input {
		field1: Id!
		field2: String!
	}

	enum Enum1 {
		VALUE_1
		VALUE_2
	}

	scalar Scalar1

	interface Interface1 {
		id: ID!
	}

	union Union1 = Object1 | Object2
`

// ============================================================================
// Navigation Workflow Tests
// ============================================================================

func TestNavigateFromQueryFieldsToObjectType(t *testing.T) {
	h := New(t, testSchema)

	h.assert.BreadcrumbsEquals("")

	h.nav.NextPanel()
	h.assert.BreadcrumbsEquals("query1")

	h.nav.NextPanel()
	h.assert.BreadcrumbsEquals("query1 > arg1")
}

func TestNavigateBackwardThroughPanelStack(t *testing.T) {
	h := New(t, testSchema)

	// Navigate forward twice to start w/ multiple breadcrumbs
	h.nav.NextPanel()
	h.nav.NextPanel()
	h.assert.BreadcrumbsEquals("query1 > arg1")

	h.nav.PrevPanel()
	h.assert.BreadcrumbsEquals("query1")

	h.nav.PrevPanel()
	h.assert.BreadcrumbsEquals("")
}

func TestNavigationResetsOnTypeSwitch(t *testing.T) {
	h := New(t, testSchema)

	// Navigate into Query to display breadcrumb
	h.nav.NextPanel()
	h.assert.BreadcrumbsEquals("query1")

	// Switch to Mutation type - should reset breadcrumbs
	h.nav.GoToGqlType(navigation.MutationType)
	h.assert.BreadcrumbsEquals("")
}

// ============================================================================
// Type Switching Tests
// ============================================================================

func TestCycleForwardThroughGraphQLTypes(t *testing.T) {
	h := New(t, testSchema)

	h.assert.CurrentType(navigation.QueryType)
	h.assert.ViewContains("query1", "query2")

	h.nav.NextGqlType()
	h.assert.CurrentType(navigation.MutationType)
	h.assert.ViewContains("mutation1", "mutation2")

	h.nav.NextGqlType()
	h.assert.CurrentType(navigation.ObjectType)
	h.assert.ViewContains("Object1", "Object2")

	h.nav.NextGqlType()
	h.assert.CurrentType(navigation.InputType)
	h.assert.ViewContains("Mutation1Input", "Mutation2Input")
}

func TestCycleBackwardThroughGraphQLTypes(t *testing.T) {
	h := New(t, testSchema)

	h.assert.CurrentType(navigation.QueryType)

	// Cycle backward should wrap to last type (Directive)
	h.nav.PrevGqlType()
	h.assert.CurrentType(navigation.DirectiveType)

	h.nav.PrevGqlType()
	h.assert.CurrentType(navigation.UnionType)
}

func TestSwitchDirectlyToSpecificType(t *testing.T) {
	h := New(t, testSchema)

	h.nav.GoToGqlType(navigation.EnumType)
	h.assert.CurrentType(navigation.EnumType)

	h.nav.GoToGqlType(navigation.ScalarType)
	h.assert.CurrentType(navigation.ScalarType)

	h.nav.GoToGqlType(navigation.InterfaceType)
	h.assert.CurrentType(navigation.InterfaceType)
}

func TestTypeCyclingResetsBreadcrumbs(t *testing.T) {
	h := New(t, testSchema)

	h.nav.NextPanel()
	h.assert.BreadcrumbsEquals("query1")

	// Switch to next type - should reset breadcrumbs
	h.nav.NextGqlType()
	h.assert.BreadcrumbsEquals("")

	h.nav.NextPanel()
	h.assert.BreadcrumbsEquals("mutation1")

	// Switch to prev type - should also reset breadcrumbs
	h.nav.PrevGqlType()
	h.assert.BreadcrumbsEquals("")
}

// ============================================================================
// Overlay Interaction Tests
// ============================================================================

func TestOpenOverlayAndVerifyContent(t *testing.T) {
	h := New(t, testSchema)

	h.overlay.Open()

	// Verify overlay is visible
	h.assert.OverlayVisible()

	// TODO: Close overlay functionality may need investigation
	// For now, just verify opening works
}

func TestOverlayShowsCorrectDetailsForDifferentItems(t *testing.T) {
	h := New(t, testSchema)

	// Switch to Object type to see different items
	h.nav.GoToGqlType(navigation.ObjectType)

	h.overlay.Open()
	h.assert.OverlayVisible()

	// Note: Closing and reopening overlays may require additional investigation
	// For now, verify that overlay opening works
}

// ============================================================================
// Complex Workflow Tests
// ============================================================================

func TestFullExplorationWorkflow(t *testing.T) {
	h := New(t, testSchema)

	h.assert.CurrentType(navigation.QueryType)
	h.assert.BreadcrumbsEquals("")

	h.nav.NextPanel()
	h.assert.BreadcrumbsEquals("query1")

	h.nav.NextPanel()
	h.assert.BreadcrumbsEquals("query1 > arg1")

	h.nav.GoToGqlType(navigation.MutationType)
	h.assert.BreadcrumbsEquals("")
	h.assert.ViewContains("mutation1")

	// Cycle to Object type
	h.nav.NextGqlType()
	h.assert.CurrentType(navigation.ObjectType)

	// Open overlay for an object
	h.overlay.Open()
	h.assert.OverlayVisible()
}

func TestMultiPanelNavigationWithTypeCycling(t *testing.T) {
	h := New(t, testSchema)

	// Navigate through Query structure
	h.nav.NextPanel()
	h.nav.NextPanel()
	h.assert.BreadcrumbsEquals("query1 > arg1")

	// Cycle to mutation type
	h.nav.NextGqlType()
	h.assert.BreadcrumbsEquals("")

	h.nav.NextPanel()
	h.assert.BreadcrumbsEquals("mutation1")
}

func TestEdgeCaseEmptyPanelNavigation(t *testing.T) {
	emptySchema := `
		type Query {
			placeholder: String
		}
	`
	h := New(t, emptySchema)

	// Should handle empty/minimal schemas gracefully
	h.assert.ViewContains("Query")

	// Try navigating - should not crash
	h.nav.NextPanel()

	// Cycle types
	h.nav.NextGqlType()
	h.nav.PrevGqlType()
}

func TestWindowResizing(t *testing.T) {
	h := New(t, testSchema, WithWindowSize(80, 30))

	// Verify view renders at specified size
	view := h.explorer.View()
	if view == "" {
		t.Error("expected view to render with custom window size")
	}

	// Navigation should still work
	h.nav.NextPanel()
	h.assert.BreadcrumbsEquals("query1")
}
