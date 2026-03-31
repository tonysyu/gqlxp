package navigation

import (
	tea "charm.land/bubbletea/v2"
	"github.com/tonysyu/gqlxp/tui/xplr/components"
)

// NavigationManager coordinates panel stack, breadcrumbs, and type selection
type NavigationManager struct {
	stack         panelStack
	breadcrumbs   breadcrumbsModel
	kindSelector  kindSelector
	visiblePanels int
}

func NewNavigationManager(visiblePanels int) NavigationManager {
	return NavigationManager{
		stack:         newPanelStack(visiblePanels),
		breadcrumbs:   newBreadcrumbsModel(),
		kindSelector:  newKindSelector(),
		visiblePanels: visiblePanels,
	}
}

// syncFocus sets the current panel as focused and all others as blurred.
// Called internally after any mutation that changes which panel is current.
func (nm NavigationManager) syncFocus() NavigationManager {
	for _, p := range nm.stack.All() {
		if p != nil {
			p.SetBlurred()
		}
	}
	if current := nm.stack.Current(); current != nil {
		current.SetFocused()
	}
	return nm
}

// NavigateForward moves forward in panel stack
// Extracts breadcrumb title from current panel's selected item and adds it BEFORE moving forward
func (nm NavigationManager) NavigateForward() (NavigationManager, bool) {
	if nm.stack.position+1 >= len(nm.stack.panels) {
		return nm, false
	}
	// Extract breadcrumb title from current panel before moving forward
	breadcrumbTitle := ""
	if currentPanel := nm.stack.Current(); currentPanel != nil {
		if selectedItem := currentPanel.SelectedItem(); selectedItem != nil {
			if listItem, ok := selectedItem.(components.ListItem); ok {
				breadcrumbTitle = listItem.RefName()
			}
		}
	}
	if breadcrumbTitle != "" {
		nm.breadcrumbs = nm.breadcrumbs.Push(breadcrumbTitle)
	}
	nm.stack.position++
	return nm.syncFocus(), true
}

// NavigateBackward moves backward in panel stack
func (nm NavigationManager) NavigateBackward() (NavigationManager, bool) {
	stack, ok := nm.stack.MoveBackward()
	if !ok {
		return nm, false
	}
	nm.stack = stack
	nm.breadcrumbs = nm.breadcrumbs.Pop()
	return nm.syncFocus(), true
}

// GoForward navigates forward in the panel stack, syncs focus, and returns any
// open command from the newly focused panel's selected item.
func (nm NavigationManager) GoForward() (NavigationManager, bool, tea.Cmd) {
	nm, moved := nm.NavigateForward()
	if !moved {
		return nm, false, nil
	}
	var openCmd tea.Cmd
	if current := nm.stack.Current(); current != nil {
		openCmd = current.OpenSelectedItem()
	}
	return nm, true, openCmd
}

// Load replaces the panel at the current position and syncs focus state.
// Use this when loading a new panel as part of navigation (instead of SetCurrentPanel).
func (nm NavigationManager) Load(panel *components.Panel) NavigationManager {
	nm.stack.panels[nm.stack.position] = panel
	return nm.syncFocus()
}

// ResetAndLoad resets the panel stack, loads the main panel at position 0,
// and optionally pushes a child panel. Syncs focus state.
func (nm NavigationManager) ResetAndLoad(panel *components.Panel, childPanel *components.Panel) NavigationManager {
	initialPanels := make([]*components.Panel, nm.visiblePanels)
	for i := range nm.visiblePanels {
		initialPanels[i] = components.NewPanel([]components.ListItem{}, "")
	}
	nm.stack = nm.stack.Replace(initialPanels)
	nm.breadcrumbs = nm.breadcrumbs.Reset()
	nm.stack.panels[nm.stack.position] = panel
	if childPanel != nil {
		nm.stack = nm.stack.Push(childPanel)
	}
	return nm.syncFocus()
}

// OpenPanel pushes new panel onto stack
func (nm NavigationManager) OpenPanel(panel *components.Panel) NavigationManager {
	nm.stack = nm.stack.Push(panel)
	return nm
}

// SwitchKind changes selected GQL type and resets breadcrumbs
func (nm NavigationManager) SwitchKind(gqlType GQLKind) NavigationManager {
	nm.kindSelector = nm.kindSelector.Set(gqlType)
	nm.breadcrumbs = nm.breadcrumbs.Reset()
	return nm
}

// CycleKindForward moves to next GQL type and resets breadcrumbs
func (nm NavigationManager) CycleKindForward() (NavigationManager, GQLKind) {
	nm.breadcrumbs = nm.breadcrumbs.Reset()
	var t GQLKind
	nm.kindSelector, t = nm.kindSelector.Next()
	return nm, t
}

// CycleKindBackward moves to previous GQL type and resets breadcrumbs
func (nm NavigationManager) CycleKindBackward() (NavigationManager, GQLKind) {
	nm.breadcrumbs = nm.breadcrumbs.Reset()
	var t GQLKind
	nm.kindSelector, t = nm.kindSelector.Previous()
	return nm, t
}

// Stack returns the panel stack (for rendering)
func (nm NavigationManager) Stack() panelStack {
	return nm.stack
}

// CurrentPanel returns the current panel from the stack
func (nm NavigationManager) CurrentPanel() *components.Panel {
	return nm.stack.Current()
}

// SetCurrentPanel sets the panel at the current stack position without syncing focus.
// Use this for panel content updates (e.g. when a panel receives a message).
// Use Load when replacing a panel as part of navigation.
func (nm NavigationManager) SetCurrentPanel(panel *components.Panel) NavigationManager {
	nm.stack.panels[nm.stack.position] = panel
	return nm
}

// NextPanel returns the panel after the current position (right panel)
func (nm NavigationManager) NextPanel() *components.Panel {
	return nm.stack.Next()
}

// CurrentKind returns currently selected GQL type
func (nm NavigationManager) CurrentKind() GQLKind {
	return nm.kindSelector.Current()
}

// AllKinds returns all available GQL types
func (nm NavigationManager) AllKinds() []GQLKind {
	return nm.kindSelector.All()
}

// Breadcrumbs returns the breadcrumb trail
func (nm NavigationManager) Breadcrumbs() []string {
	return nm.breadcrumbs.Get()
}

// Reset clears the panel stack to initial state with empty panels and resets breadcrumbs
func (nm NavigationManager) Reset() NavigationManager {
	initialPanels := make([]*components.Panel, nm.visiblePanels)
	for i := range nm.visiblePanels {
		initialPanels[i] = components.NewPanel([]components.ListItem{}, "")
	}
	nm.stack = nm.stack.Replace(initialPanels)
	nm.breadcrumbs = nm.breadcrumbs.Reset()
	return nm.syncFocus()
}

// IsAtTopLevelPanel returns true if the current panel is the first panel (position 0)
// which corresponds to top-level GQL type panels (Query, Mutation, Object, etc.)
func (nm NavigationManager) IsAtTopLevelPanel() bool {
	return nm.stack.Position() == 0
}
