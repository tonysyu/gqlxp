package navigation

import (
	"testing"

	"github.com/tonysyu/gqlxp/tui/components"
)

func TestPanelStack_NewPanelStack(t *testing.T) {
	stack := NewPanelStack(2)
	if stack.Len() != 0 {
		t.Errorf("expected empty stack, got length %d", stack.Len())
	}
	if stack.Position() != 0 {
		t.Errorf("expected position 0, got %d", stack.Position())
	}
}

func TestPanelStack_Push(t *testing.T) {
	stack := NewPanelStack(2)
	p1 := components.NewEmptyListPanel("1")
	p2 := components.NewEmptyListPanel("2")

	stack.Push(p1)
	if stack.Len() != 1 {
		t.Errorf("expected length 1, got %d", stack.Len())
	}
	if stack.Current() != p1 {
		t.Error("expected current panel to be p1")
	}

	stack.Push(p2)
	if stack.Len() != 2 {
		t.Errorf("expected length 2, got %d", stack.Len())
	}
	if stack.Current() != p1 {
		t.Error("expected current panel to still be p1 (position unchanged)")
	}
}

func TestPanelStack_Push_TruncatesAfterCurrent(t *testing.T) {
	stack := NewPanelStack(3)
	p1 := components.NewEmptyListPanel("1")
	p2 := components.NewEmptyListPanel("2")
	p3 := components.NewEmptyListPanel("3")
	p4 := components.NewEmptyListPanel("4")

	stack.Push(p1)
	stack.Push(p2)
	stack.Push(p3)
	// Stack: [p1, p2, p3], position: 0

	stack.MoveForward()
	// Stack: [p1, p2, p3], position: 1 (on p2)

	stack.Push(p4)
	// Should truncate p3 and add p4: [p1, p2, p4]
	if stack.Len() != 3 {
		t.Errorf("expected length 3 after truncate, got %d", stack.Len())
	}
	if stack.All()[2] != p4 {
		t.Error("expected last panel to be p4")
	}
}

func TestPanelStack_MoveForward(t *testing.T) {
	stack := NewPanelStack(2)
	p1 := components.NewEmptyListPanel("1")
	p2 := components.NewEmptyListPanel("2")

	stack.Push(p1)
	stack.Push(p2)

	moved := stack.MoveForward()
	if !moved {
		t.Error("expected move forward to succeed")
	}
	if stack.Position() != 1 {
		t.Errorf("expected position 1, got %d", stack.Position())
	}
	if stack.Current() != p2 {
		t.Error("expected current panel to be p2")
	}

	// Can't move past end
	moved = stack.MoveForward()
	if moved {
		t.Error("expected move forward to fail at end of stack")
	}
	if stack.Position() != 1 {
		t.Errorf("expected position to remain 1, got %d", stack.Position())
	}
}

func TestPanelStack_MoveBackward(t *testing.T) {
	stack := NewPanelStack(2)
	p1 := components.NewEmptyListPanel("1")
	p2 := components.NewEmptyListPanel("2")

	stack.Push(p1)
	stack.Push(p2)
	stack.MoveForward()
	// position: 1 (on p2)

	moved := stack.MoveBackward()
	if !moved {
		t.Error("expected move backward to succeed")
	}
	if stack.Position() != 0 {
		t.Errorf("expected position 0, got %d", stack.Position())
	}
	if stack.Current() != p1 {
		t.Error("expected current panel to be p1")
	}

	// Can't move before beginning
	moved = stack.MoveBackward()
	if moved {
		t.Error("expected move backward to fail at beginning of stack")
	}
	if stack.Position() != 0 {
		t.Errorf("expected position to remain 0, got %d", stack.Position())
	}
}

func TestPanelStack_Current(t *testing.T) {
	stack := NewPanelStack(2)
	if stack.Current() != nil {
		t.Error("expected nil for empty stack")
	}

	p1 := components.NewEmptyListPanel("1")
	stack.Push(p1)
	if stack.Current() != p1 {
		t.Error("expected current to be p1")
	}
}

func TestPanelStack_Next(t *testing.T) {
	stack := NewPanelStack(2)
	p1 := components.NewEmptyListPanel("1")
	p2 := components.NewEmptyListPanel("2")

	stack.Push(p1)
	if stack.Next() != nil {
		t.Error("expected nil when no next panel")
	}

	stack.Push(p2)
	if stack.Next() != p2 {
		t.Error("expected next to be p2")
	}

	stack.MoveForward()
	if stack.Next() != nil {
		t.Error("expected nil when at end of stack")
	}
}

func TestPanelStack_Replace(t *testing.T) {
	stack := NewPanelStack(2)
	p1 := components.NewEmptyListPanel("1")
	p2 := components.NewEmptyListPanel("2")
	p3 := components.NewEmptyListPanel("3")

	stack.Push(p1)
	stack.Push(p2)
	stack.MoveForward()
	// position: 1

	newPanels := []components.Panel{p3}
	stack.Replace(newPanels)

	if stack.Len() != 1 {
		t.Errorf("expected length 1 after replace, got %d", stack.Len())
	}
	if stack.Position() != 0 {
		t.Errorf("expected position reset to 0, got %d", stack.Position())
	}
	if stack.Current() != p3 {
		t.Error("expected current panel to be p3")
	}
}

func TestPanelStack_All(t *testing.T) {
	stack := NewPanelStack(3)
	p1 := components.NewEmptyListPanel("1")
	p2 := components.NewEmptyListPanel("2")

	stack.Push(p1)
	stack.Push(p2)

	all := stack.All()
	if len(all) != 2 {
		t.Errorf("expected 2 panels, got %d", len(all))
	}
	if all[0] != p1 || all[1] != p2 {
		t.Error("expected panels in order [p1, p2]")
	}
}
