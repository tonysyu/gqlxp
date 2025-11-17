package navigation

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/tui/xplr/components"
)

func TestNavigationManager_NewNavigationManager(t *testing.T) {
	is := is.New(t)
	nm := NewNavigationManager(2)

	is.True(nm.Stack() != nil)
	is.Equal(nm.CurrentType(), QueryType)
	is.Equal(nm.Breadcrumbs(), nil)
}

func TestNavigationManager_NavigateForward(t *testing.T) {
	is := is.New(t)
	nm := NewNavigationManager(2)

	p1 := components.NewEmptyPanel("1")
	p2 := components.NewEmptyPanel("2")
	nm.OpenPanel(p1)
	nm.OpenPanel(p2)

	moved := nm.NavigateForward()
	is.True(moved)

	is.Equal(nm.Stack().Position(), 1)

	// Can't move past end
	moved = nm.NavigateForward()
	is.True(!moved)
}

func TestNavigationManager_NavigateBackward(t *testing.T) {
	is := is.New(t)
	nm := NewNavigationManager(2)

	p1 := components.NewEmptyPanel("1")
	p2 := components.NewEmptyPanel("2")
	nm.OpenPanel(p1)
	nm.OpenPanel(p2)
	nm.NavigateForward()

	moved := nm.NavigateBackward()
	is.True(moved)
	is.Equal(nm.Stack().Position(), 0)
}

func TestNavigationManager_OpenPanel(t *testing.T) {
	is := is.New(t)
	nm := NewNavigationManager(2)

	p1 := components.NewEmptyPanel("1")
	nm.OpenPanel(p1)

	is.Equal(nm.Stack().Len(), 1)
	is.Equal(nm.CurrentPanel(), p1)
}

func TestNavigationManager_SwitchType(t *testing.T) {
	is := is.New(t)
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

	is.Equal(nm.CurrentType(), MutationType)
	breadcrumbs := nm.Breadcrumbs()
	is.Equal(breadcrumbs, nil)
}

func TestNavigationManager_CycleTypeForward(t *testing.T) {
	is := is.New(t)
	nm := NewNavigationManager(2)

	is.Equal(nm.CurrentType(), QueryType)

	newType := nm.CycleTypeForward()
	is.Equal(newType, MutationType)
	is.Equal(nm.CurrentType(), MutationType)
}

func TestNavigationManager_CycleTypeBackward(t *testing.T) {
	is := is.New(t)
	nm := NewNavigationManager(2)

	is.Equal(nm.CurrentType(), QueryType)

	newType := nm.CycleTypeBackward()
	is.Equal(newType, DirectiveType)
	is.Equal(nm.CurrentType(), DirectiveType)
}

func TestNavigationManager_AllTypes(t *testing.T) {
	is := is.New(t)
	nm := NewNavigationManager(2)

	allTypes := nm.AllTypes()
	is.Equal(len(allTypes), 9)
	is.Equal(allTypes[0], QueryType)
}
