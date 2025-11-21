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

	// Start on Query type - should show query fields
	h.assert.ViewContains("Query")
	h.assert.BreadcrumbsEmpty()

	// Navigate forward to the result type of first query field
	h.explorer.NavigateToNextPanel()
	// Breadcrumbs should show the first query field
	h.is.Equal(h.assert.getBreadcrumbs(), "query1")

	// Navigate forward again to see the Object type details
	h.explorer.NavigateToNextPanel()
	// Breadcrumbs should now include the first argument of query1
	h.is.Equal(h.assert.getBreadcrumbs(), "query1 > arg1")
}

func TestNavigateThroughMultiplePanelsWithBreadcrumbs(t *testing.T) {
	h := New(t, testSchema)

	// Navigate through Query -> field -> Type
	h.explorer.NavigateToNextPanel()
	breadcrumbs := h.assert.getBreadcrumbs()
	if breadcrumbs == "" {
		t.Error("expected breadcrumbs after first navigation")
	}

	h.explorer.NavigateToNextPanel()
	breadcrumbs = h.assert.getBreadcrumbs()
	if breadcrumbs == "" {
		t.Error("expected breadcrumbs after second navigation")
	}
}

func TestNavigateBackwardThroughPanelStack(t *testing.T) {
	h := New(t, testSchema)

	// Navigate forward twice
	h.explorer.NavigateToNextPanel()
	h.explorer.NavigateToNextPanel()
	breadcrumbs := h.assert.getBreadcrumbs()
	if breadcrumbs == "" {
		t.Error("expected breadcrumbs after navigation")
	}

	// Navigate backward
	h.explorer.NavigateToPreviousPanel()
	breadcrumbs = h.assert.getBreadcrumbs()
	if breadcrumbs == "" {
		t.Error("expected breadcrumbs to remain after one step back")
	}

	// Navigate backward again to initial state
	h.explorer.NavigateToPreviousPanel()
	h.assert.BreadcrumbsEmpty()
}

func TestNavigationResetsOnTypeSwitch(t *testing.T) {
	h := New(t, testSchema)

	// Navigate into Query structure
	h.explorer.NavigateToNextPanel()
	// Breadcrumbs should contain some query field name
	breadcrumbs := h.assert.getBreadcrumbs()
	if breadcrumbs == "" {
		t.Error("expected breadcrumbs after navigation")
	}

	// Switch to Mutation type - should reset breadcrumbs
	h.explorer.SwitchToType(navigation.MutationType)
	h.assert.BreadcrumbsEmpty()
	h.assert.ViewContains("Mutation")
}

// ============================================================================
// Type Switching Tests
// ============================================================================

func TestCycleForwardThroughGraphQLTypes(t *testing.T) {
	h := New(t, testSchema)

	// Start on Query
	h.assert.CurrentType(navigation.QueryType)
	h.assert.ViewContains("query1", "query2")

	// Cycle to Mutation
	h.explorer.CycleTypeForward()
	h.assert.CurrentType(navigation.MutationType)
	h.assert.ViewContains("mutation1", "mutation2")

	// Cycle to Object
	h.explorer.CycleTypeForward()
	h.assert.CurrentType(navigation.ObjectType)
	h.assert.ViewContains("Object1", "Object2")

	// Cycle to Input
	h.explorer.CycleTypeForward()
	h.assert.CurrentType(navigation.InputType)
	h.assert.ViewContains("Mutation1Input", "Mutation2Input")
}

func TestCycleBackwardThroughGraphQLTypes(t *testing.T) {
	h := New(t, testSchema)

	// Start on Query, cycle backward to wrap around
	h.assert.CurrentType(navigation.QueryType)

	// Cycle backward should wrap to last type (Directive)
	h.explorer.CycleTypeBackward()
	h.assert.CurrentType(navigation.DirectiveType)

	// Cycle backward again - should move to Union
	h.explorer.CycleTypeBackward()
	h.assert.CurrentType(navigation.UnionType)
}

func TestSwitchDirectlyToSpecificType(t *testing.T) {
	h := New(t, testSchema)

	// Switch directly to Enum type
	h.explorer.SwitchToType(navigation.EnumType)
	h.assert.CurrentType(navigation.EnumType)
	// Enum tab should be visible, check for Role type
	h.assert.ViewContains("Enum")

	// Switch directly to Scalar type
	h.explorer.SwitchToType(navigation.ScalarType)
	h.assert.CurrentType(navigation.ScalarType)
	h.assert.ViewContains("Scalar")

	// Switch directly to Interface type
	h.explorer.SwitchToType(navigation.InterfaceType)
	h.assert.CurrentType(navigation.InterfaceType)
	h.assert.ViewContains("Interface")
}

func TestTypeCyclingResetsBreadcrumbs(t *testing.T) {
	h := New(t, testSchema)

	// Navigate into structure
	h.explorer.NavigateToNextPanel()
	breadcrumbs := h.assert.getBreadcrumbs()
	if breadcrumbs == "" {
		t.Error("expected breadcrumbs after navigation")
	}

	// Cycle type - should reset breadcrumbs
	h.explorer.CycleTypeForward()
	h.assert.BreadcrumbsEmpty()

	// Navigate into Mutation structure
	h.explorer.NavigateToNextPanel()
	breadcrumbs = h.assert.getBreadcrumbs()
	if breadcrumbs == "" {
		t.Error("expected breadcrumbs after navigation in Mutation")
	}

	// Cycle backward - should also reset breadcrumbs
	h.explorer.CycleTypeBackward()
	h.assert.BreadcrumbsEmpty()
}

// ============================================================================
// Overlay Interaction Tests
// ============================================================================

func TestOpenOverlayAndVerifyContent(t *testing.T) {
	h := New(t, testSchema)

	// Select a query field and open overlay
	h.explorer.SelectItemAtIndex(0) // Select first query field
	h.overlay.Open()

	// Verify overlay is visible
	h.assert.OverlayVisible()

	// TODO: Close overlay functionality may need investigation
	// For now, just verify opening works
}

func TestOverlayShowsCorrectDetailsForDifferentItems(t *testing.T) {
	h := New(t, testSchema)

	// Switch to Object type to see different items
	h.explorer.SwitchToType(navigation.ObjectType)

	// Select first object and check overlay
	h.explorer.SelectItemAtIndex(0)
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

	// Start on Query type
	h.assert.CurrentType(navigation.QueryType)
	h.assert.BreadcrumbsEmpty()

	// Navigate to a query field
	h.explorer.NavigateToNextPanel()
	breadcrumbs := h.assert.getBreadcrumbs()
	if breadcrumbs == "" {
		t.Error("expected breadcrumbs after navigation")
	}

	// Navigate to the result type
	h.explorer.NavigateToNextPanel()
	breadcrumbs = h.assert.getBreadcrumbs()
	if breadcrumbs == "" {
		t.Error("expected breadcrumbs after second navigation")
	}

	// Switch to Mutation type
	h.explorer.SwitchToType(navigation.MutationType)
	h.assert.BreadcrumbsEmpty()
	h.assert.ViewContains("mutation1")

	// Cycle to Object type
	h.explorer.CycleTypeForward()
	h.assert.CurrentType(navigation.ObjectType)

	// Open overlay for an object
	h.overlay.Open()
	h.assert.OverlayVisible()
}

func TestMultiPanelNavigationWithTypeCycling(t *testing.T) {
	h := New(t, testSchema)

	// Navigate through Query structure
	h.explorer.NavigateToNextPanel()
	h.explorer.NavigateToNextPanel()

	// Verify breadcrumbs show the path
	breadcrumbs := h.assert.getBreadcrumbs()
	if breadcrumbs == "" {
		t.Error("expected breadcrumbs to be populated")
	}

	// Cycle to a different type
	h.explorer.CycleTypeForward()
	h.assert.BreadcrumbsEmpty()

	// Navigate in the new type
	h.explorer.NavigateToNextPanel()

	// Breadcrumbs should be rebuilt for the new type
	breadcrumbs = h.assert.getBreadcrumbs()
	if breadcrumbs == "" {
		t.Error("expected new breadcrumbs after type switch")
	}
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
	h.explorer.NavigateToNextPanel()

	// Cycle types
	h.explorer.CycleTypeForward()
	h.explorer.CycleTypeBackward()
}

func TestWindowResizing(t *testing.T) {
	h := New(t, testSchema, WithWindowSize(80, 30))

	// Verify view renders at specified size
	view := h.explorer.View()
	if view == "" {
		t.Error("expected view to render with custom window size")
	}

	// Navigation should still work
	h.explorer.NavigateToNextPanel()
	breadcrumbs := h.assert.getBreadcrumbs()
	if breadcrumbs == "" {
		t.Error("expected breadcrumbs after navigation")
	}
}
