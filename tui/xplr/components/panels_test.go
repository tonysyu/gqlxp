package components

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/utils/testx/assert"
)

// testOpenableItem is a test helper that implements ListItem with a working Open() method
type testOpenableItem struct {
	SimpleItem
	openPanel *Panel
}

func (i testOpenableItem) OpenPanel() (*Panel, bool) {
	return i.openPanel, true
}

func TestPanelBasic(t *testing.T) {
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

	panel := NewPanel(listItems, "Test Panel")

	// Test initial state
	is.Equal(panel.lastSelectedIndex, -1)
	is.Equal(panel.Title(), "Test Panel")

	// Test SetSize
	panel.SetSize(80, 20)
	// SetSize should be called without error
}

func TestPanel_SelectItemByName(t *testing.T) {
	is := is.New(t)

	// Create test items with different RefNames
	items := []ListItem{
		NewSimpleItem("User", WithTypeName("User")),
		NewSimpleItem("Post", WithTypeName("Post")),
		NewSimpleItem("Comment", WithTypeName("Comment")),
	}

	panel := NewPanel(items, "Test Panel")

	// Test selecting existing item
	found := panel.SelectItemByName("Post")
	is.True(found)
	is.Equal(panel.ListModel.Index(), 1)
	is.Equal(panel.lastSelectedIndex, 1)

	// Test selecting another item
	found = panel.SelectItemByName("User")
	is.True(found)
	is.Equal(panel.ListModel.Index(), 0)
	is.Equal(panel.lastSelectedIndex, 0)

	// Test selecting non-existent item
	found = panel.SelectItemByName("NonExistent")
	is.True(!found)                      // Should return false
	is.Equal(panel.ListModel.Index(), 0) // Should remain at previous selection
}

func TestPanelWithEmptyItems(t *testing.T) {
	is := is.New(t)

	panel := NewPanel([]ListItem{}, "Empty Panel")
	panel.SetSize(80, 20)

	// Should handle empty list gracefully
	view := panel.View()
	is.True(strings.Contains(view, "Empty Panel"))
}

func TestPanelSelectionChange(t *testing.T) {
	is := is.New(t)

	// Create items with Open capability
	testPanel := NewEmptyPanel("opened panel content")
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
	panel := NewPanel(items, "Test Panel")

	// Simulate key down to change selection
	_, cmd := panel.Update(tea.KeyMsg{Type: tea.KeyDown})

	// Should generate OpenPanelMsg command when selection changes
	is.True(cmd != nil)
}

func TestPanelAutoOpen(t *testing.T) {
	is := is.New(t)

	// Create items with Open capability
	testPanel := NewEmptyPanel("opened panel content")
	items := []ListItem{
		testOpenableItem{
			SimpleItem: NewSimpleItem("Test Field"),
			openPanel:  testPanel,
		},
	}
	panel := NewPanel(items, "Test Panel")

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

func TestPanelTitleSetting(t *testing.T) {
	is := is.New(t)

	panel := NewPanel([]ListItem{}, "Initial Title")
	is.Equal(panel.Title(), "Initial Title")

	// Test SetTitle
	panel.SetTitle("Updated Title")
	is.Equal(panel.Title(), "Updated Title")
}

func TestPanelWithManyItems(t *testing.T) {
	is := is.New(t)

	// Create many simple items
	var items []ListItem
	for i := 0; i < 100; i++ {
		items = append(items, NewSimpleItem(
			"Item "+string(rune(i)),
			WithDescription("Description "+string(rune(i))),
		))
	}

	panel := NewPanel(items, "Large Panel")
	panel.SetSize(80, 20)

	// Should handle large lists without issues
	view := panel.View()
	is.True(len(view) > 0)

	// Should be able to navigate
	_, _ = panel.Update(tea.KeyMsg{Type: tea.KeyDown})
}

func TestPanelSizeEdgeCases(t *testing.T) {
	is := is.New(t)

	// Test with very small sizes
	panel := NewEmptyPanel("test")
	panel.SetSize(1, 1)
	view := panel.View()
	is.True(len(view) >= 0) // Should not crash

	// Test with zero sizes
	panel.SetSize(0, 0)
	view = panel.View()
	is.True(len(view) >= 0) // Should not crash

	// Test with very large sizes
	panel.SetSize(10000, 1000)
	view = panel.View()
	is.True(len(view) >= 0) // Should not crash
}

func TestPanelFilteringSupport(t *testing.T) {
	is := is.New(t)

	// Create items with different filter values
	items := []ListItem{
		NewSimpleItem("Apple"),
		NewSimpleItem("Banana"),
		NewSimpleItem("Carrot"),
	}

	panel := NewPanel(items, "Filterable Panel")
	panel.SetSize(80, 20)

	// Test that items have proper FilterValue implementation
	listItems := panel.Items()
	is.Equal(len(listItems), 3)

	for i, item := range listItems {
		defaultItem := item.(list.DefaultItem)
		is.Equal(defaultItem.FilterValue(), items[i].FilterValue())
	}
}

func TestPanelFilterExitRefresh(t *testing.T) {
	is := is.New(t)

	// Create items with Open capability
	items := []ListItem{
		testOpenableItem{
			SimpleItem: NewSimpleItem("Apple"),
			openPanel:  NewEmptyPanel("Apple subpanel"),
		},
		testOpenableItem{
			SimpleItem: NewSimpleItem("Banana"),
			openPanel:  NewEmptyPanel("Banana subpanel"),
		},
		testOpenableItem{
			SimpleItem: NewSimpleItem("Carrot"),
			openPanel:  NewEmptyPanel("Carrot subpanel"),
		},
	}
	panel := NewPanel(items, "Test Panel")
	panel.SetSize(80, 20)

	// The list.Model starts at index 0 by default, but lastSelectedIndex is -1
	// So we need to trigger an initial update to sync them
	_, cmd := panel.Update(tea.KeyMsg{Type: tea.KeyDown})
	is.True(cmd != nil) // Should open panel for selected item
	is.Equal(panel.lastSelectedIndex, 1)

	// Start filtering by pressing "/"
	_, _ = panel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	// wasFiltering should now be true (we're in filter mode)
	is.True(panel.wasFiltering)

	// Type a filter - the cursor might move during filtering
	_, _ = panel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	is.True(panel.wasFiltering) // Still filtering

	// Accept the filter by pressing Enter
	// This should exit filter mode and trigger a refresh even if cursor is at the same index
	_, cmd = panel.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// After accepting filter:
	// 1. wasFiltering should be false (no longer in filter mode)
	// 2. A command should be generated to refresh the inactive panel
	is.True(!panel.wasFiltering) // Should have exited filter mode
	is.True(cmd != nil)          // Should have generated OpenPanel command

	// Verify the command is an OpenPanelMsg
	if cmd != nil {
		msg := cmd()
		_, ok := msg.(OpenPanelMsg)
		is.True(ok) // Should be OpenPanelMsg
	}
}

func TestPanelTitleTruncation(t *testing.T) {
	assert := assert.New(t)

	panel := NewPanel([]ListItem{}, "This is a very long title that should be truncated")
	// Set a narrow width that should trigger truncation
	panel.SetSize(20, 10)

	view := panel.View()

	// PanelTitleHPadding is 1, so effective width for text is width - 2*1 = 18
	assert.StringContains(view, "This is a very lo…")
}

func TestPanelResultTypeTruncation(t *testing.T) {
	assert := assert.New(t)

	panel := NewPanel([]ListItem{}, "Test Panel")
	panel.SetTabs([]Tab{
		{Label: "Type", Content: []ListItem{NewSimpleItem("VeryLongResultTypeNameThatShouldBeTruncated")}},
	})
	// Set a narrow width that should trigger truncation
	panel.SetSize(20, 10)

	view := panel.View()

	// ItemLeftPadding is 2, so effective width for text is width - 2 = 18
	assert.StringContains(view, "VeryLongResultTyp…")
}
