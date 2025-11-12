package navigation

import "github.com/tonysyu/gqlxp/tui/components"

// PanelStack manages a stack of panels with navigation
type PanelStack struct {
	panels   []components.Panel
	position int
}

func NewPanelStack(initialCapacity int) *PanelStack {
	return &PanelStack{
		panels:   make([]components.Panel, 0, initialCapacity),
		position: 0,
	}
}

// Current returns the currently focused panel
func (s *PanelStack) Current() components.Panel {
	if s.position >= 0 && s.position < len(s.panels) {
		return s.panels[s.position]
	}
	return nil
}

// Next returns the panel after the current position (right panel)
func (s *PanelStack) Next() components.Panel {
	nextPos := s.position + 1
	if nextPos < len(s.panels) {
		return s.panels[nextPos]
	}
	return nil
}

// MoveForward advances position if possible, returns success
func (s *PanelStack) MoveForward() bool {
	if s.position+1 < len(s.panels) {
		s.position++
		return true
	}
	return false
}

// MoveBackward moves position back if possible, returns success
func (s *PanelStack) MoveBackward() bool {
	if s.position > 0 {
		s.position--
		return true
	}
	return false
}

// Push adds a panel after current position, truncating rest
func (s *PanelStack) Push(panel components.Panel) {
	// Only truncate if we're not at the end
	if s.position+1 < len(s.panels) {
		s.panels = s.panels[:s.position+1]
	}
	s.panels = append(s.panels, panel)
}

// Replace replaces all panels with new set
func (s *PanelStack) Replace(panels []components.Panel) {
	s.panels = panels
	s.position = 0
}

// Position returns current position
func (s *PanelStack) Position() int {
	return s.position
}

// Len returns number of panels
func (s *PanelStack) Len() int {
	return len(s.panels)
}

// All returns all panels (for iteration)
func (s *PanelStack) All() []components.Panel {
	return s.panels
}
