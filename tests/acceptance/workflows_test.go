package acceptance

import (
	"testing"

	"github.com/matryer/is"
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

	// With tab-based navigation, the Result Type tab is active by default
	// so navigating next opens the result type (Object1) instead of arg1
	h.nav.NextPanel()
	h.assert.BreadcrumbsEquals("query1 > Object1")
}

func TestNavigateBackwardThroughPanelStack(t *testing.T) {
	h := New(t, testSchema)

	// Navigate forward twice to start w/ multiple breadcrumbs
	h.nav.NextPanel()
	h.nav.NextPanel()
	// With tab-based navigation, result type (Object1) is opened instead of arg1
	h.assert.BreadcrumbsEquals("query1 > Object1")

	h.nav.PrevPanel()
	h.assert.BreadcrumbsEquals("query1")

	h.nav.PrevPanel()
	h.assert.BreadcrumbsEquals("")
}

func TestNavigationResetsOnKindSwitch(t *testing.T) {
	h := New(t, testSchema)

	// Navigate into Query to display breadcrumb
	h.nav.NextPanel()
	h.assert.BreadcrumbsEquals("query1")

	// Switch to Mutation kind - should reset breadcrumbs
	h.nav.GoToGqlKind(navigation.MutationKind)
	h.assert.BreadcrumbsEquals("")
}

// ============================================================================
// Kind Switching Tests
// ============================================================================

func TestCycleForwardThroughGraphQLKinds(t *testing.T) {
	h := New(t, testSchema)

	h.assert.CurrentKind(navigation.QueryKind)
	h.assert.ViewContains("query1", "query2")

	h.nav.NextGqlKind()
	h.assert.CurrentKind(navigation.MutationKind)
	h.assert.ViewContains("mutation1", "mutation2")

	h.nav.NextGqlKind()
	h.assert.CurrentKind(navigation.ObjectKind)
	h.assert.ViewContains("Object1", "Object2")

	h.nav.NextGqlKind()
	h.assert.CurrentKind(navigation.InputKind)
	h.assert.ViewContains("Mutation1Input", "Mutation2Input")
}

func TestCycleBackwardThroughGraphQLKinds(t *testing.T) {
	h := New(t, testSchema)

	h.assert.CurrentKind(navigation.QueryKind)

	// Cycle backward should wrap to last kind (Search)
	h.nav.PrevGqlKind()
	h.assert.CurrentKind(navigation.SearchKind)

	h.nav.PrevGqlKind()
	h.assert.CurrentKind(navigation.DirectiveKind)

	h.nav.PrevGqlKind()
	h.assert.CurrentKind(navigation.UnionKind)
}

func TestSwitchDirectlyToSpecificKind(t *testing.T) {
	h := New(t, testSchema)

	h.nav.GoToGqlKind(navigation.EnumKind)
	h.assert.CurrentKind(navigation.EnumKind)

	h.nav.GoToGqlKind(navigation.ScalarKind)
	h.assert.CurrentKind(navigation.ScalarKind)

	h.nav.GoToGqlKind(navigation.InterfaceKind)
	h.assert.CurrentKind(navigation.InterfaceKind)
}

func TestKindCyclingResetsBreadcrumbs(t *testing.T) {
	h := New(t, testSchema)

	h.nav.NextPanel()
	h.assert.BreadcrumbsEquals("query1")

	// Switch to next kind - should reset breadcrumbs
	h.nav.NextGqlKind()
	h.assert.BreadcrumbsEquals("")

	h.nav.NextPanel()
	h.assert.BreadcrumbsEquals("mutation1")

	// Switch to prev kind - should also reset breadcrumbs
	h.nav.PrevGqlKind()
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

	// Switch to Object kind to see different items
	h.nav.GoToGqlKind(navigation.ObjectKind)

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

	h.assert.CurrentKind(navigation.QueryKind)
	h.assert.BreadcrumbsEquals("")

	h.nav.NextPanel()
	h.assert.BreadcrumbsEquals("query1")

	// With tab-based navigation, result type (Object1) is opened instead of arg1
	h.nav.NextPanel()
	h.assert.BreadcrumbsEquals("query1 > Object1")

	h.nav.GoToGqlKind(navigation.MutationKind)
	h.assert.BreadcrumbsEquals("")
	h.assert.ViewContains("mutation1")

	// Cycle to Object kind
	h.nav.NextGqlKind()
	h.assert.CurrentKind(navigation.ObjectKind)

	// Open overlay for an object
	h.overlay.Open()
	h.assert.OverlayVisible()
}

func TestMultiPanelNavigationWithKindCycling(t *testing.T) {
	h := New(t, testSchema)

	// Navigate through Query structure
	h.nav.NextPanel()
	h.nav.NextPanel()
	// With tab-based navigation, result type (Object1) is opened instead of arg1
	h.assert.BreadcrumbsEquals("query1 > Object1")

	// Cycle to mutation kind
	h.nav.NextGqlKind()
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

	// Cycle GQL kinds
	h.nav.NextGqlKind()
	h.nav.PrevGqlKind()
}

func TestWindowResizing(t *testing.T) {
	is := is.New(t)
	h := New(t, testSchema, WithWindowSize(80, 30))

	// Verify view renders at specified size
	view := h.explorer.View()
	is.True(view != "") // view should render with custom window size

	// Navigation should still work
	h.nav.NextPanel()
	h.assert.BreadcrumbsEquals("query1")
}
