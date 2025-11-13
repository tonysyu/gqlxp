package navigation

import (
	"testing"

	"github.com/tonysyu/gqlxp/tui/components"
)

func TestNavigationManager_NewNavigationManager(t *testing.T) {
	nm := NewNavigationManager(2)

	if nm.Stack() == nil {
		t.Error("expected stack to be initialized")
	}
	if nm.CurrentType() != QueryType {
		t.Errorf("expected default type to be QueryType, got %v", nm.CurrentType())
	}
	if nm.Breadcrumbs() != nil {
		t.Error("expected initial breadcrumbs to be nil/empty")
	}
}

func TestNavigationManager_NavigateForward(t *testing.T) {
	nm := NewNavigationManager(2)

	p1 := components.NewEmptyPanel("1")
	p2 := components.NewEmptyPanel("2")
	nm.OpenPanel(p1)
	nm.OpenPanel(p2)

	moved := nm.NavigateForward()
	if !moved {
		t.Error("expected navigate forward to succeed")
	}

	if nm.Stack().Position() != 1 {
		t.Errorf("expected position 1, got %d", nm.Stack().Position())
	}

	// Can't move past end
	moved = nm.NavigateForward()
	if moved {
		t.Error("expected navigate forward to fail at end")
	}
}

func TestNavigationManager_NavigateBackward(t *testing.T) {
	nm := NewNavigationManager(2)

	p1 := components.NewEmptyPanel("1")
	p2 := components.NewEmptyPanel("2")
	nm.OpenPanel(p1)
	nm.OpenPanel(p2)
	nm.NavigateForward()

	moved := nm.NavigateBackward()
	if !moved {
		t.Error("expected navigate backward to succeed")
	}
	if nm.Stack().Position() != 0 {
		t.Errorf("expected position 0, got %d", nm.Stack().Position())
	}
}

func TestNavigationManager_OpenPanel(t *testing.T) {
	nm := NewNavigationManager(2)

	p1 := components.NewEmptyPanel("1")
	nm.OpenPanel(p1)

	if nm.Stack().Len() != 1 {
		t.Errorf("expected 1 panel in stack, got %d", nm.Stack().Len())
	}
	if nm.CurrentPanel() != p1 {
		t.Error("expected current panel to be p1")
	}
}

func TestNavigationManager_SwitchType(t *testing.T) {
	nm := NewNavigationManager(2)

	// Add panels and navigate
	// Need to navigate between OpenPanel calls to avoid truncation
	p1 := components.NewEmptyPanel("1")
	p2 := components.NewEmptyPanel("2")
	p3 := components.NewEmptyPanel("3")
	nm.OpenPanel(p1)
	nm.OpenPanel(p2)
	nm.NavigateForward() // Move to position 1
	nm.OpenPanel(p3)     // Add p3 at position 2
	nm.NavigateForward() // Move to position 2

	nm.SwitchType(MutationType)

	if nm.CurrentType() != MutationType {
		t.Errorf("expected current type to be MutationType, got %v", nm.CurrentType())
	}
	breadcrumbs := nm.Breadcrumbs()
	if breadcrumbs != nil {
		t.Errorf("expected breadcrumbs to be reset (nil/empty), got %v", breadcrumbs)
	}
}

func TestNavigationManager_CycleTypeForward(t *testing.T) {
	nm := NewNavigationManager(2)

	if nm.CurrentType() != QueryType {
		t.Error("setup error: should start with QueryType")
	}

	newType := nm.CycleTypeForward()
	if newType != MutationType {
		t.Errorf("expected MutationType, got %v", newType)
	}
	if nm.CurrentType() != MutationType {
		t.Errorf("expected current type to be MutationType, got %v", nm.CurrentType())
	}
}

func TestNavigationManager_CycleTypeBackward(t *testing.T) {
	nm := NewNavigationManager(2)

	if nm.CurrentType() != QueryType {
		t.Error("setup error: should start with QueryType")
	}

	newType := nm.CycleTypeBackward()
	if newType != DirectiveType {
		t.Errorf("expected DirectiveType (wraparound), got %v", newType)
	}
	if nm.CurrentType() != DirectiveType {
		t.Errorf("expected current type to be DirectiveType, got %v", nm.CurrentType())
	}
}

func TestNavigationManager_AllTypes(t *testing.T) {
	nm := NewNavigationManager(2)

	allTypes := nm.AllTypes()
	if len(allTypes) != 9 {
		t.Errorf("expected 9 types, got %d", len(allTypes))
	}
	if allTypes[0] != QueryType {
		t.Errorf("expected first type to be QueryType, got %v", allTypes[0])
	}
}
