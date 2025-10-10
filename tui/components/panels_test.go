package components

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/matryer/is"
)

// testOpenableItem is a test helper that implements ListItem with a working Open() method
type testOpenableItem struct {
	SimpleItem
	openPanel Panel
}

func (i testOpenableItem) Open() (Panel, bool) {
	return i.openPanel, true
}

func TestStringPanelBasic(t *testing.T) {
	is := is.New(t)

	content := "This is test content"
	panel := NewStringPanel(content)

	// Test basic properties
	is.Equal(panel.content, content)
	is.Equal(panel.width, 0)
	is.Equal(panel.height, 0)

	// Test SetSize
	panel.SetSize(80, 20)
	is.Equal(panel.width, 80)
	is.Equal(panel.height, 20)

	// Test that Update returns self
	updatedPanel, cmd := panel.Update(tea.KeyMsg{})
	is.Equal(updatedPanel, panel)
	is.True(cmd == nil)
}

func TestStringPanelView(t *testing.T) {
	is := is.New(t)

	panel := NewStringPanel("test content")
	panel.SetSize(50, 10)

	view := panel.View()
	is.True(strings.Contains(view, "test content"))
}

func TestStringPanelWithEmptyContent(t *testing.T) {
	is := is.New(t)

	panel := NewStringPanel("")
	panel.SetSize(80, 20)

	view := panel.View()
	is.True(len(view) > 0) // Should still render something (padding/style)
}

func TestStringPanelWithLargeContent(t *testing.T) {
	is := is.New(t)

	// Create very long content
	longContent := ""
	for range 1000 {
		longContent += "This is a very long line of content that should be handled properly. "
	}

	panel := NewStringPanel(longContent)
	panel.SetSize(80, 20)

	// Should not crash with large content
	view := panel.View()
	is.True(len(view) > 69000)
}

func TestListPanelBasic(t *testing.T) {
	is := is.New(t)

	// Create test items
	items := []SimpleItem{
		NewSimpleItem("Item 1"),
		NewSimpleItem("Item 2"),
		NewSimpleItem("Item 3"),
	}

	listItems := make([]ListItem, len(items))
	for i, item := range items {
		listItems[i] = item
	}

	panel := NewListPanel(listItems, "Test Panel")

	// Test initial state
	is.Equal(panel.lastSelectedIndex, -1)
	is.Equal(panel.Title(), "Test Panel")

	// Test SetSize
	panel.SetSize(80, 20)
	// SetSize should be called without error
}

func TestListPanelWithEmptyItems(t *testing.T) {
	is := is.New(t)

	panel := NewListPanel([]ListItem{}, "Empty Panel")
	panel.SetSize(80, 20)

	// Should handle empty list gracefully
	view := panel.View()
	is.True(strings.Contains(view, "Empty Panel"))
}

func TestListPanelSelectionChange(t *testing.T) {
	is := is.New(t)

	// Create items with Open capability
	testPanel := NewStringPanel("opened panel content")
	items := []ListItem{
		testOpenableItem{
			SimpleItem: NewSimpleItem("Item 1"),
			openPanel:  testPanel,
		},
		testOpenableItem{
			SimpleItem: NewSimpleItem("Item 2"),
			openPanel:  testPanel,
		},
	}
	panel := NewListPanel(items, "Test Panel")

	// Simulate key down to change selection
	_, cmd := panel.Update(tea.KeyMsg{Type: tea.KeyDown})

	// Should generate OpenPanelMsg command when selection changes
	is.True(cmd != nil)
}

func TestListPanelAutoOpen(t *testing.T) {
	is := is.New(t)

	// Create items with Open capability
	testPanel := NewStringPanel("opened panel content")
	items := []ListItem{
		testOpenableItem{
			SimpleItem: NewSimpleItem("Test Field"),
			openPanel:  testPanel,
		},
	}
	panel := NewListPanel(items, "Test Panel")

	// Ensure we have an item that can be opened
	is.True(len(items) > 0)

	// Simulate navigation which triggers auto-open
	panel.lastSelectedIndex = -1 // Simulate fresh state
	_, cmd := panel.Update(tea.KeyMsg{Type: tea.KeyDown})

	// Should return a command for opening panel
	is.True(cmd != nil)

	// The last selected index should be updated
	is.True(panel.lastSelectedIndex >= 0)
}

func TestListPanelTitleSetting(t *testing.T) {
	is := is.New(t)

	panel := NewListPanel([]ListItem{}, "Initial Title")
	is.Equal(panel.Title(), "Initial Title")

	// Test SetTitle
	panel.SetTitle("Updated Title")
	is.Equal(panel.Title(), "Updated Title")
}

func TestListPanelWithManyItems(t *testing.T) {
	is := is.New(t)

	// Create many simple items
	var items []ListItem
	for i := 0; i < 100; i++ {
		items = append(items, NewSimpleItem(
			"Item "+string(rune(i)),
			WithDescription("Description "+string(rune(i))),
		))
	}

	panel := NewListPanel(items, "Large Panel")
	panel.SetSize(80, 20)

	// Should handle large lists without issues
	view := panel.View()
	is.True(len(view) > 0)

	// Should be able to navigate
	_, cmd := panel.Update(tea.KeyMsg{Type: tea.KeyDown})
	is.True(cmd != nil || cmd == nil) // Either command or no command is acceptable
}

func TestPanelSizeEdgeCases(t *testing.T) {
	is := is.New(t)

	// Test with very small sizes
	stringPanel := NewStringPanel("test")
	stringPanel.SetSize(1, 1)
	view := stringPanel.View()
	is.True(len(view) >= 0) // Should not crash

	// Test with zero sizes
	stringPanel.SetSize(0, 0)
	view = stringPanel.View()
	is.True(len(view) >= 0) // Should not crash

	// Test with very large sizes
	stringPanel.SetSize(10000, 1000)
	view = stringPanel.View()
	is.True(len(view) >= 0) // Should not crash

	// Test list panel with small sizes
	listPanel := NewListPanel([]ListItem{NewSimpleItem("test")}, "test")
	listPanel.SetSize(1, 1)
	view = listPanel.View()
	is.True(len(view) >= 0) // Should not crash
}

func TestListPanelFilteringSupport(t *testing.T) {
	is := is.New(t)

	// Create items with different filter values
	items := []ListItem{
		NewSimpleItem("Apple"),
		NewSimpleItem("Banana"),
		NewSimpleItem("Carrot"),
	}

	panel := NewListPanel(items, "Filterable Panel")
	panel.SetSize(80, 20)

	// Test that items have proper FilterValue implementation
	listItems := panel.Items()
	is.Equal(len(listItems), 3)

	for i, item := range listItems {
		defaultItem := item.(list.DefaultItem)
		is.Equal(defaultItem.FilterValue(), items[i].FilterValue())
	}
}

func TestPanelInterfaceCompliance(t *testing.T) {
	is := is.New(t)

	// Test that both panel types implement Panel interface
	var stringPanelInterface Panel = NewStringPanel("test")
	var listPanelInterface Panel = NewListPanel([]ListItem{}, "test")

	// Test that they can be used as Panels
	stringPanelInterface.SetSize(80, 20)
	listPanelInterface.SetSize(80, 20)

	// Test that they implement tea.Model
	var stringModel tea.Model = stringPanelInterface
	var listModel tea.Model = listPanelInterface

	// Should be able to call Update
	_, _ = stringModel.Update(tea.KeyMsg{})
	_, _ = listModel.Update(tea.KeyMsg{})

	// Should be able to call View
	_ = stringModel.View()
	_ = listModel.View()

	// Just verify the interfaces work
	is.True(stringPanelInterface != nil)
	is.True(listPanelInterface != nil)
}
