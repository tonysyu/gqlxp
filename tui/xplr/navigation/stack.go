package navigation

import "github.com/tonysyu/gqlxp/tui/xplr/components"

// panelStack manages a stack of panels with navigation
type panelStack struct {
	panels   []*components.Panel
	position int
}

func newPanelStack(initialCapacity int) panelStack {
	return panelStack{
		panels:   make([]*components.Panel, 0, initialCapacity),
		position: 0,
	}
}

// Current returns the currently focused panel
func (s panelStack) Current() *components.Panel {
	if s.position >= 0 && s.position < len(s.panels) {
		return s.panels[s.position]
	}
	return nil
}

// Next returns the panel after the current position (right panel)
func (s panelStack) Next() *components.Panel {
	nextPos := s.position + 1
	if nextPos < len(s.panels) {
		return s.panels[nextPos]
	}
	return nil
}

// MoveForward advances position if possible, returns updated stack and success
func (s panelStack) MoveForward() (panelStack, bool) {
	if s.position+1 < len(s.panels) {
		s.position++
		return s, true
	}
	return s, false
}

// MoveBackward moves position back if possible, returns updated stack and success
func (s panelStack) MoveBackward() (panelStack, bool) {
	if s.position > 0 {
		s.position--
		return s, true
	}
	return s, false
}

// Push adds a panel after current position, truncating rest
func (s panelStack) Push(panel *components.Panel) panelStack {
	// Only truncate if we're not at the end
	if s.position+1 < len(s.panels) {
		s.panels = s.panels[:s.position+1]
	}
	s.panels = append(s.panels, panel)
	return s
}

// Replace replaces all panels with new set
func (s panelStack) Replace(panels []*components.Panel) panelStack {
	s.panels = panels
	s.position = 0
	return s
}

// Position returns current position
func (s panelStack) Position() int {
	return s.position
}

// Len returns number of panels
func (s panelStack) Len() int {
	return len(s.panels)
}

// All returns all panels (for iteration)
func (s panelStack) All() []*components.Panel {
	return s.panels
}
