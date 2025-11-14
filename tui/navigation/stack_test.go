package navigation

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/tui/components"
)

func TestPanelStack_NewPanelStack(t *testing.T) {
	is := is.New(t)
	stack := NewPanelStack(2)
	is.Equal(stack.Len(), 0)
	is.Equal(stack.Position(), 0)
}

func TestPanelStack_Push(t *testing.T) {
	is := is.New(t)
	stack := NewPanelStack(2)
	p1 := components.NewEmptyPanel("1")
	p2 := components.NewEmptyPanel("2")

	stack.Push(p1)
	is.Equal(stack.Len(), 1)
	is.Equal(stack.Current(), p1)

	stack.Push(p2)
	is.Equal(stack.Len(), 2)
	is.Equal(stack.Current(), p1)
}

func TestPanelStack_Push_TruncatesAfterCurrent(t *testing.T) {
	is := is.New(t)
	stack := NewPanelStack(3)
	p1 := components.NewEmptyPanel("1")
	p2 := components.NewEmptyPanel("2")
	p3 := components.NewEmptyPanel("3")
	p4 := components.NewEmptyPanel("4")

	stack.Push(p1)
	stack.Push(p2)
	stack.Push(p3)
	// Stack: [p1, p2, p3], position: 0

	stack.MoveForward()
	// Stack: [p1, p2, p3], position: 1 (on p2)

	stack.Push(p4)
	// Should truncate p3 and add p4: [p1, p2, p4]
	is.Equal(stack.Len(), 3)
	is.Equal(stack.All()[2], p4)
}

func TestPanelStack_MoveForward(t *testing.T) {
	is := is.New(t)
	stack := NewPanelStack(2)
	p1 := components.NewEmptyPanel("1")
	p2 := components.NewEmptyPanel("2")

	stack.Push(p1)
	stack.Push(p2)

	moved := stack.MoveForward()
	is.True(moved)
	is.Equal(stack.Position(), 1)
	is.Equal(stack.Current(), p2)

	// Can't move past end
	moved = stack.MoveForward()
	is.True(!moved)
	is.Equal(stack.Position(), 1)
}

func TestPanelStack_MoveBackward(t *testing.T) {
	is := is.New(t)
	stack := NewPanelStack(2)
	p1 := components.NewEmptyPanel("1")
	p2 := components.NewEmptyPanel("2")

	stack.Push(p1)
	stack.Push(p2)
	stack.MoveForward()
	// position: 1 (on p2)

	moved := stack.MoveBackward()
	is.True(moved)
	is.Equal(stack.Position(), 0)
	is.Equal(stack.Current(), p1)

	// Can't move before beginning
	moved = stack.MoveBackward()
	is.True(!moved)
	is.Equal(stack.Position(), 0)
}

func TestPanelStack_Current(t *testing.T) {
	is := is.New(t)
	stack := NewPanelStack(2)
	is.Equal(stack.Current(), nil)

	p1 := components.NewEmptyPanel("1")
	stack.Push(p1)
	is.Equal(stack.Current(), p1)
}

func TestPanelStack_Next(t *testing.T) {
	is := is.New(t)
	stack := NewPanelStack(2)
	p1 := components.NewEmptyPanel("1")
	p2 := components.NewEmptyPanel("2")

	stack.Push(p1)
	is.Equal(stack.Next(), nil)

	stack.Push(p2)
	is.Equal(stack.Next(), p2)

	stack.MoveForward()
	is.Equal(stack.Next(), nil)
}

func TestPanelStack_Replace(t *testing.T) {
	is := is.New(t)
	stack := NewPanelStack(2)
	p1 := components.NewEmptyPanel("1")
	p2 := components.NewEmptyPanel("2")
	p3 := components.NewEmptyPanel("3")

	stack.Push(p1)
	stack.Push(p2)
	stack.MoveForward()
	// position: 1

	newPanels := []*components.Panel{p3}
	stack.Replace(newPanels)

	is.Equal(stack.Len(), 1)
	is.Equal(stack.Position(), 0)
	is.Equal(stack.Current(), p3)
}

func TestPanelStack_All(t *testing.T) {
	is := is.New(t)
	stack := NewPanelStack(3)
	p1 := components.NewEmptyPanel("1")
	p2 := components.NewEmptyPanel("2")

	stack.Push(p1)
	stack.Push(p2)

	all := stack.All()
	is.Equal(len(all), 2)
	is.Equal(all[0], p1)
	is.Equal(all[1], p2)
}
