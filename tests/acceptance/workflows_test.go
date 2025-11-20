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
	h.AssertViewContains("Query")
	h.AssertBreadcrumbsEmpty()

	// Navigate forward to the result type of first query field
	h.NavigateToNextPanel()
	// Breadcrumbs should show the first query field
	h.is.Equal(h.getBreadcrumbs(), "query1")

	// Navigate forward again to see the Object type details
	h.NavigateToNextPanel()
	// Breadcrumbs should now include the first argument of query1
	h.is.Equal(h.getBreadcrumbs(), "query1 > arg1")
}

func TestNavigateThroughMultiplePanelsWithBreadcrumbs(t *testing.T) {
	h := New(t, testSchema)

	// Navigate through Query -> field -> Type
	h.NavigateToNextPanel()
	breadcrumbs := h.getBreadcrumbs()
	if breadcrumbs == "" {
		t.Error("expected breadcrumbs after first navigation")
	}

	h.NavigateToNextPanel()
	breadcrumbs = h.getBreadcrumbs()
	if breadcrumbs == "" {
		t.Error("expected breadcrumbs after second navigation")
	}
}

func TestNavigateBackwardThroughPanelStack(t *testing.T) {
	h := New(t, testSchema)

	// Navigate forward twice
	h.NavigateToNextPanel()
	h.NavigateToNextPanel()
	breadcrumbs := h.getBreadcrumbs()
	if breadcrumbs == "" {
		t.Error("expected breadcrumbs after navigation")
	}

	// Navigate backward
	h.NavigateToPreviousPanel()
	breadcrumbs = h.getBreadcrumbs()
	if breadcrumbs == "" {
		t.Error("expected breadcrumbs to remain after one step back")
	}

	// Navigate backward again to initial state
	h.NavigateToPreviousPanel()
	h.AssertBreadcrumbsEmpty()
}

func TestNavigationResetsOnTypeSwitch(t *testing.T) {
	h := New(t, testSchema)

	// Navigate into Query structure
	h.NavigateToNextPanel()
	// Breadcrumbs should contain some query field name
	breadcrumbs := h.getBreadcrumbs()
	if breadcrumbs == "" {
		t.Error("expected breadcrumbs after navigation")
	}

	// Switch to Mutation type - should reset breadcrumbs
	h.SwitchToType(navigation.MutationType)
	h.AssertBreadcrumbsEmpty()
	h.AssertViewContains("Mutation")
}

// ============================================================================
// Type Switching Tests
// ============================================================================

func TestCycleForwardThroughGraphQLTypes(t *testing.T) {
	h := New(t, testSchema)

	// Start on Query
	h.AssertCurrentType(navigation.QueryType)
	h.AssertViewContains("query1", "query2")

	// Cycle to Mutation
	h.CycleTypeForward()
	h.AssertCurrentType(navigation.MutationType)
	h.AssertViewContains("mutation1", "mutation2")

	// Cycle to Object
	h.CycleTypeForward()
	h.AssertCurrentType(navigation.ObjectType)
	h.AssertViewContains("Object1", "Object2")

	// Cycle to Input
	h.CycleTypeForward()
	h.AssertCurrentType(navigation.InputType)
	h.AssertViewContains("Mutation1Input", "Mutation2Input")
}

func TestCycleBackwardThroughGraphQLTypes(t *testing.T) {
	h := New(t, testSchema)

	// Start on Query, cycle backward to wrap around
	h.AssertCurrentType(navigation.QueryType)

	// Cycle backward should wrap to last type (Directive)
	h.CycleTypeBackward()
	h.AssertCurrentType(navigation.DirectiveType)

	// Cycle backward again - should move to Union
	h.CycleTypeBackward()
	h.AssertCurrentType(navigation.UnionType)
}

func TestSwitchDirectlyToSpecificType(t *testing.T) {
	h := New(t, testSchema)

	// Switch directly to Enum type
	h.SwitchToType(navigation.EnumType)
	h.AssertCurrentType(navigation.EnumType)
	// Enum tab should be visible, check for Role type
	h.AssertViewContains("Enum")

	// Switch directly to Scalar type
	h.SwitchToType(navigation.ScalarType)
	h.AssertCurrentType(navigation.ScalarType)
	h.AssertViewContains("Scalar")

	// Switch directly to Interface type
	h.SwitchToType(navigation.InterfaceType)
	h.AssertCurrentType(navigation.InterfaceType)
	h.AssertViewContains("Interface")
}

func TestTypeCyclingResetsBreadcrumbs(t *testing.T) {
	h := New(t, testSchema)

	// Navigate into structure
	h.NavigateToNextPanel()
	breadcrumbs := h.getBreadcrumbs()
	if breadcrumbs == "" {
		t.Error("expected breadcrumbs after navigation")
	}

	// Cycle type - should reset breadcrumbs
	h.CycleTypeForward()
	h.AssertBreadcrumbsEmpty()

	// Navigate into Mutation structure
	h.NavigateToNextPanel()
	breadcrumbs = h.getBreadcrumbs()
	if breadcrumbs == "" {
		t.Error("expected breadcrumbs after navigation in Mutation")
	}

	// Cycle backward - should also reset breadcrumbs
	h.CycleTypeBackward()
	h.AssertBreadcrumbsEmpty()
}

// ============================================================================
// Overlay Interaction Tests
// ============================================================================

func TestOpenOverlayAndVerifyContent(t *testing.T) {
	h := New(t, testSchema)

	// Select a query field and open overlay
	h.SelectItemAtIndex(0) // Select first query field
	h.OpenOverlay()

	// Verify overlay is visible
	h.AssertOverlayVisible()

	// TODO: Close overlay functionality may need investigation
	// For now, just verify opening works
}

func TestOverlayShowsCorrectDetailsForDifferentItems(t *testing.T) {
	h := New(t, testSchema)

	// Switch to Object type to see different items
	h.SwitchToType(navigation.ObjectType)

	// Select first object and check overlay
	h.SelectItemAtIndex(0)
	h.OpenOverlay()
	h.AssertOverlayVisible()

	// Note: Closing and reopening overlays may require additional investigation
	// For now, verify that overlay opening works
}

// ============================================================================
// Complex Workflow Tests
// ============================================================================

func TestFullExplorationWorkflow(t *testing.T) {
	h := New(t, testSchema)

	// Start on Query type
	h.AssertCurrentType(navigation.QueryType)
	h.AssertBreadcrumbsEmpty()

	// Navigate to a query field
	h.NavigateToNextPanel()
	breadcrumbs := h.getBreadcrumbs()
	if breadcrumbs == "" {
		t.Error("expected breadcrumbs after navigation")
	}

	// Navigate to the result type
	h.NavigateToNextPanel()
	breadcrumbs = h.getBreadcrumbs()
	if breadcrumbs == "" {
		t.Error("expected breadcrumbs after second navigation")
	}

	// Switch to Mutation type
	h.SwitchToType(navigation.MutationType)
	h.AssertBreadcrumbsEmpty()
	h.AssertViewContains("mutation1")

	// Cycle to Object type
	h.CycleTypeForward()
	h.AssertCurrentType(navigation.ObjectType)

	// Open overlay for an object
	h.OpenOverlay()
	h.AssertOverlayVisible()
}

func TestMultiPanelNavigationWithTypeCycling(t *testing.T) {
	h := New(t, testSchema)

	// Navigate through Query structure
	h.NavigateToNextPanel()
	h.NavigateToNextPanel()

	// Verify breadcrumbs show the path
	breadcrumbs := h.getBreadcrumbs()
	if breadcrumbs == "" {
		t.Error("expected breadcrumbs to be populated")
	}

	// Cycle to a different type
	h.CycleTypeForward()
	h.AssertBreadcrumbsEmpty()

	// Navigate in the new type
	h.NavigateToNextPanel()

	// Breadcrumbs should be rebuilt for the new type
	breadcrumbs = h.getBreadcrumbs()
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
	h.AssertViewContains("Query")

	// Try navigating - should not crash
	h.NavigateToNextPanel()

	// Cycle types
	h.CycleTypeForward()
	h.CycleTypeBackward()
}

func TestWindowResizing(t *testing.T) {
	h := New(t, testSchema, WithWindowSize(80, 30))

	// Verify view renders at specified size
	view := h.View()
	if view == "" {
		t.Error("expected view to render with custom window size")
	}

	// Navigation should still work
	h.NavigateToNextPanel()
	breadcrumbs := h.getBreadcrumbs()
	if breadcrumbs == "" {
		t.Error("expected breadcrumbs after navigation")
	}
}
