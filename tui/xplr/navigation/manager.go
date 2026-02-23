package navigation

import "github.com/tonysyu/gqlxp/tui/xplr/components"

// NavigationManager coordinates panel stack, breadcrumbs, and type selection
type NavigationManager struct {
	stack         panelStack
	breadcrumbs   breadcrumbsModel
	typeSelector  typeSelector
	visiblePanels int
}

func NewNavigationManager(visiblePanels int) NavigationManager {
	return NavigationManager{
		stack:         newPanelStack(visiblePanels),
		breadcrumbs:   newBreadcrumbsModel(),
		typeSelector:  newTypeSelector(),
		visiblePanels: visiblePanels,
	}
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
	return nm, true
}

// NavigateBackward moves backward in panel stack
func (nm NavigationManager) NavigateBackward() (NavigationManager, bool) {
	stack, ok := nm.stack.MoveBackward()
	if !ok {
		return nm, false
	}
	nm.stack = stack
	nm.breadcrumbs = nm.breadcrumbs.Pop()
	return nm, true
}

// OpenPanel pushes new panel onto stack
func (nm NavigationManager) OpenPanel(panel *components.Panel) NavigationManager {
	nm.stack = nm.stack.Push(panel)
	return nm
}

// SwitchType changes selected GQL type and resets breadcrumbs
func (nm NavigationManager) SwitchType(gqlType GQLType) NavigationManager {
	nm.typeSelector = nm.typeSelector.Set(gqlType)
	nm.breadcrumbs = nm.breadcrumbs.Reset()
	return nm
}

// CycleTypeForward moves to next GQL type and resets breadcrumbs
func (nm NavigationManager) CycleTypeForward() (NavigationManager, GQLType) {
	nm.breadcrumbs = nm.breadcrumbs.Reset()
	var t GQLType
	nm.typeSelector, t = nm.typeSelector.Next()
	return nm, t
}

// CycleTypeBackward moves to previous GQL type and resets breadcrumbs
func (nm NavigationManager) CycleTypeBackward() (NavigationManager, GQLType) {
	nm.breadcrumbs = nm.breadcrumbs.Reset()
	var t GQLType
	nm.typeSelector, t = nm.typeSelector.Previous()
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

// SetCurrentPanel sets the panel at the current stack position
func (nm NavigationManager) SetCurrentPanel(panel *components.Panel) NavigationManager {
	nm.stack.panels[nm.stack.position] = panel
	return nm
}

// NextPanel returns the panel after the current position (right panel)
func (nm NavigationManager) NextPanel() *components.Panel {
	return nm.stack.Next()
}

// CurrentType returns currently selected GQL type
func (nm NavigationManager) CurrentType() GQLType {
	return nm.typeSelector.Current()
}

// AllTypes returns all available GQL types
func (nm NavigationManager) AllTypes() []GQLType {
	return nm.typeSelector.All()
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
	return nm
}

// IsAtTopLevelPanel returns true if the current panel is the first panel (position 0)
// which corresponds to top-level GQL type panels (Query, Mutation, Object, etc.)
func (nm NavigationManager) IsAtTopLevelPanel() bool {
	return nm.stack.Position() == 0
}
