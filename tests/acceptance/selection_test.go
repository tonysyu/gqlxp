package acceptance

import (
	"testing"

	"github.com/tonysyu/gqlxp/tui/xplr/navigation"
)

// ============================================================================
// Selection Target Tests
// ============================================================================

func TestSelectTypeByName(t *testing.T) {
	h := New(t, testSchema, WithSelection("Object1", ""))

	// Should switch to Object kind
	h.assert.CurrentKind(navigation.ObjectKind)
	// Should have Object1 selected
	h.assert.ViewContains("Object1")
}

func TestSelectInputTypeByName(t *testing.T) {
	h := New(t, testSchema, WithSelection("Mutation1Input", ""))

	// Should switch to Input kind
	h.assert.CurrentKind(navigation.InputKind)
	// Should have Mutation1Input selected
	h.assert.ViewContains("Mutation1Input")
}

func TestSelectEnumTypeByName(t *testing.T) {
	h := New(t, testSchema, WithSelection("Enum1", ""))

	// Should switch to Enum kind
	h.assert.CurrentKind(navigation.EnumKind)
	// Should have Enum1 selected
	h.assert.ViewContains("Enum1")
}

func TestSelectFieldWithinType(t *testing.T) {
	h := New(t, testSchema, WithSelection("Query", "query1"))

	// Should be in Query kind
	h.assert.CurrentKind(navigation.QueryKind)
	// Query type shows fields directly, so query1 should be selected
	// After selection, we navigate forward which adds the breadcrumb
	h.assert.ViewContains("query1")
	// The breadcrumb will show once we navigate to the detail panel
	// In this case, the selection pre-navigates to show query1's details
	h.assert.BreadcrumbsEquals("query1")
}

func TestSelectFieldWithinMutation(t *testing.T) {
	h := New(t, testSchema, WithSelection("Mutation", "mutation1"))

	// Should be in Mutation kind
	h.assert.CurrentKind(navigation.MutationKind)
	// Mutation type shows fields directly, so mutation1 should be selected
	h.assert.ViewContains("mutation1")
	// The breadcrumb will show once we navigate to the detail panel
	h.assert.BreadcrumbsEquals("mutation1")
}

func TestSelectNonExistentType(t *testing.T) {
	h := New(t, testSchema, WithSelection("NonExistent", ""))

	// Should gracefully fallback to default (Query type)
	h.assert.CurrentKind(navigation.QueryKind)
	// Should have no special selection
	h.assert.BreadcrumbsEquals("")
}

func TestSelectNonExistentField(t *testing.T) {
	h := New(t, testSchema, WithSelection("Query", "nonExistentField"))

	// Should be in Query kind
	h.assert.CurrentKind(navigation.QueryKind)
	// Should have Query selected but field not found - no breadcrumbs for missing field
	// The exact behavior is graceful fallback - no error, just Query selected
	h.assert.ViewContains("query1") // Query type should still be visible
}

func TestSelectWithEmptyTarget(t *testing.T) {
	h := New(t, testSchema, WithSelection("", ""))

	// Should gracefully handle empty selection and default to Query
	h.assert.CurrentKind(navigation.QueryKind)
	h.assert.BreadcrumbsEquals("")
}

func TestSelectObjectTypeAndNavigate(t *testing.T) {
	h := New(t, testSchema, WithSelection("Object1", ""))

	// Verify Object1 is selected
	h.assert.CurrentKind(navigation.ObjectKind)

	// Navigate to the detail panel
	h.nav.NextPanel()
	// Should show fields of Object1
	h.assert.ViewContains("field1", "field2")
}
