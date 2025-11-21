// Package acceptance provides functional, domain-specific helpers for acceptance testing of the TUI.
//
// The test harness simplifies writing high-level tests that verify complete user workflows
// by providing explorer helpers (navigation, type cycling, item selection) and screen
// verification helpers (panel assertions, breadcrumb checks, overlay verification).
//
// Example usage:
//
//	h := New(t, `
//		type Query { users: [User!]! }
//		type User { id: ID!, name: String! }
//	`)
//
//	h.SelectItem("users")
//	h.NavigateToNextPanel()
//	h.assert.PanelContains(1, "User")
package acceptance

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/tui/adapters"
	"github.com/tonysyu/gqlxp/tui/xplr"
	"github.com/tonysyu/gqlxp/tui/xplr/navigation"
	"github.com/tonysyu/gqlxp/utils/testx"
	"github.com/tonysyu/gqlxp/utils/text"
)

// Key messages for simulating user input
var (
	keyNextPanel = tea.KeyMsg{Type: tea.KeyTab}
	keyPrevPanel = tea.KeyMsg{Type: tea.KeyShiftTab}
	keyNextType  = tea.KeyMsg{Type: tea.KeyCtrlT}
	keyPrevType  = tea.KeyMsg{Type: tea.KeyCtrlR}
	keySpace     = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}}
	keyEscape    = tea.KeyMsg{Type: tea.KeyEsc}
	keyDown      = tea.KeyMsg{Type: tea.KeyDown}
	keyUp        = tea.KeyMsg{Type: tea.KeyUp}
)

// Harness provides high-level test utilities for acceptance testing
type Harness struct {
	explorer *Explorer
	assert   *Assert
	overlay  *Overlay
}

// Explorer provides methods for navigating and interacting with the TUI model
type Explorer struct {
	model xplr.Model
}

// Assert provides assertion helpers for test verification
type Assert struct {
	explorer *Explorer
	t        *testing.T
	is       *is.I
}

// Overlay provides helpers for overlay interaction
type Overlay struct {
	explorer *Explorer
}

// Option is a functional option for configuring the test harness
type Option func(*Harness)

// WithWindowSize sets the window dimensions for the test harness
func WithWindowSize(width, height int) Option {
	return func(h *Harness) {
		h.explorer.model, _ = h.explorer.model.Update(tea.WindowSizeMsg{Width: width, Height: height})
	}
}

// WithoutOverlayBorders removes overlay border styling to simplify test comparisons
// This allows tests to focus on content verification without dealing with box-drawing characters
func WithoutOverlayBorders() Option {
	return func(h *Harness) {
		h.explorer.model.Overlay.Styles.Overlay = lipgloss.NewStyle()
	}
}

// New creates a new test harness with the given options
// If no schema is provided, a default test schema is used
func New(t *testing.T, schema string, opts ...Option) *Harness {
	schemaView, err := adapters.ParseSchemaString(schema)
	if err != nil {
		t.Fatalf("failed to parse schema: %v", err)
	}

	explorer := &Explorer{
		model: xplr.New(schemaView),
	}
	h := &Harness{
		explorer: explorer,
		assert:   &Assert{explorer: explorer, t: t, is: is.New(t)},
		overlay:  &Overlay{explorer: explorer},
	}

	// Apply options
	for _, opt := range opts {
		opt(h)
	}

	// Set default window size if not already set
	if h.explorer.model.Width() == 0 {
		WithWindowSize(120, 40)(h)
	}

	return h
}

// ============================================================================
// Explorer Helpers - Navigation and Interaction
// ============================================================================

// Update sends a message to the model and handles any resulting commands
func (e *Explorer) Update(msg tea.Msg) {
	var cmd tea.Cmd
	e.model, cmd = e.model.Update(msg)
	if cmd != nil {
		// If Update returns a cmd, execute it and recursively Update model
		e.Update(cmd())
	}
}

// View returns the rendered view
func (e *Explorer) View() string {
	return e.model.View()
}

// NavigateToNextPanel navigates forward to the next panel (Tab key)
func (e *Explorer) NavigateToNextPanel() {
	e.Update(keyNextPanel)
}

// NavigateToPreviousPanel navigates backward to the previous panel (Shift+Tab key)
func (e *Explorer) NavigateToPreviousPanel() {
	e.Update(keyPrevPanel)
}

// CycleTypeForward cycles to the next GraphQL type (Ctrl+T key)
func (e *Explorer) CycleTypeForward() {
	e.Update(keyNextType)
}

// CycleTypeBackward cycles to the previous GraphQL type (Ctrl+R key)
func (e *Explorer) CycleTypeBackward() {
	e.Update(keyPrevType)
}

// SwitchToType switches directly to a specific GraphQL type
func (e *Explorer) SwitchToType(gqlType navigation.GQLType) {
	e.model.SwitchToType(string(gqlType))
}

// GetCurrentType returns the currently selected GraphQL type
func (e *Explorer) GetCurrentType() navigation.GQLType {
	return navigation.GQLType(e.model.CurrentType())
}

// SelectItemAtIndex moves the cursor to the item at the given index (0-based)
func (e *Explorer) SelectItemAtIndex(idx int) {
	// Move cursor to top first (safety measure)
	for i := 0; i < 100; i++ {
		e.Update(keyUp)
	}
	// Move down to the desired index
	for i := 0; i < idx; i++ {
		e.Update(keyDown)
	}
}

// SelectItem moves the cursor to the item with the given name by checking the view
// This is useful for selecting items by their display name
func (e *Explorer) SelectItem(name string) {
	// Since we can't directly access panel items, we'll try selecting by index
	// based on what appears in the view. For now, we use a simpler approach:
	// just move down until we find the item in the view.
	// A more robust implementation would parse the view to find the exact index.

	// Reset to top
	for i := 0; i < 100; i++ {
		e.Update(keyUp)
	}

	// Try moving down and checking if the item is selected
	// This is a simplified heuristic - in practice, we'd need better view parsing
	for i := 0; i < 20; i++ {
		// For now, we'll just move to the first item by default
		// A real implementation would parse the view to find the exact position
		if i > 0 {
			e.Update(keyDown)
		}
		// Check if view contains the selected item (this is approximate)
		view := e.View()
		if strings.Contains(view, name) {
			// Assume we found it - this is a simplified implementation
			break
		}
	}
}

// ============================================================================
// Component Helpers - Overlay
// ============================================================================
// Open opens the overlay for the currently selected item (Space key)

func (o *Overlay) Open() {
	o.explorer.Update(keySpace)
}

// Close closes the currently open overlay (Escape key)
func (o *Overlay) Close() {
	o.explorer.Update(keyEscape)
}

// OpenForType switches to the specified type and opens the overlay for the first item
// This is a convenience method for testing overlay content for different GraphQL types
func (o *Overlay) OpenForType(gqlType navigation.GQLType) {
	o.explorer.SwitchToType(gqlType)
	o.explorer.SelectItemAtIndex(0)
	o.Open()
}

// ============================================================================
// Screen Verification Helpers - Assertions
// ============================================================================

// getPanelContent returns the rendered content of the panel at the given index
func (a *Assert) getPanelContent(panelIdx int) string {
	// For now, we extract panel content from the full view
	// This is a simplified approach - in a real implementation, we might need
	// more sophisticated view parsing
	view := a.explorer.View()
	lines := text.SplitLines(view)

	// Skip navbar (line 0), breadcrumbs (line 1-2), get panel content
	if len(lines) > 3 {
		panelContent := strings.Join(lines[3:len(lines)-1], "\n") // Exclude help line
		return panelContent
	}
	return ""
}

// getBreadcrumbs returns the breadcrumb text from the rendered view
func (a *Assert) getBreadcrumbs() string {
	view := a.explorer.View()
	lines := text.SplitLines(view)
	// Breadcrumbs are on line index 2 (after navbar and empty line)
	// But we need to look for the line that contains breadcrumb separators " > "
	// or check a few lines after the navbar
	if len(lines) > 3 {
		// Check lines 1-3 for breadcrumbs (accounting for empty lines)
		for i := 1; i < 4 && i < len(lines); i++ {
			line := strings.TrimSpace(lines[i])
			// Breadcrumbs contain " > " separator or are non-empty text between navbar and panels
			// Empty lines and panel borders start with special characters
			if line != "" && !strings.HasPrefix(line, "╭") && !strings.Contains(line, "Query") {
				return line
			}
		}
	}
	return ""
}

// PanelContains checks if the panel at the given index contains the expected text
// panelIdx is 0-based (0 = left panel, 1 = right panel)
func (a *Assert) PanelContains(panelIdx int, expected string) {
	a.t.Helper()
	content := a.getPanelContent(panelIdx)
	if !strings.Contains(content, expected) {
		a.t.Errorf("panel %d does not contain %q\nPanel content:\n%s", panelIdx, expected, content)
	}
}

// PanelEquals checks if the normalized panel content equals the expected text
func (a *Assert) PanelEquals(panelIdx int, expected string) {
	a.t.Helper()
	content := testx.NormalizeView(a.getPanelContent(panelIdx))
	expectedNormalized := testx.NormalizeView(expected)
	a.is.Equal(content, expectedNormalized)
}

// BreadcrumbsEquals checks if breadcrumbs exactly match the expected text
func (a *Assert) BreadcrumbsEquals(expected string) {
	a.t.Helper()
	breadcrumbs := a.getBreadcrumbs()
	a.is.Equal(breadcrumbs, expected)
}

// OverlayVisible checks if the overlay is currently visible
func (a *Assert) OverlayVisible() {
	a.t.Helper()
	if !a.explorer.model.Overlay.IsActive() {
		a.t.Error("overlay is not visible")
	}
}

// OverlayContains checks if the overlay contains the expected text
func (a *Assert) OverlayContains(expected string) {
	a.t.Helper()
	a.OverlayVisible()
	view := a.explorer.View()
	if !strings.Contains(view, expected) {
		a.t.Errorf("overlay does not contain %q\nOverlay content:\n%s", expected, view)
	}
}

// OverlayContainsNormalized checks if the normalized overlay view contains the expected text
// Both the view and expected text are normalized (whitespace trimmed, excess spaces removed)
// before comparison. This is useful for comparing content structure without exact formatting.
func (a *Assert) OverlayContainsNormalized(expected string) {
	a.t.Helper()
	a.OverlayVisible()
	normalizedView := testx.NormalizeView(a.explorer.View())
	normalizedExpected := testx.NormalizeView(expected)
	if !strings.Contains(normalizedView, normalizedExpected) {
		a.t.Errorf("normalized overlay does not contain expected content\nExpected:\n%s\n\nActual overlay:\n%s",
			normalizedExpected, normalizedView)
	}
}

// ViewContains checks if the entire view contains the expected text
func (a *Assert) ViewContains(expectedStrings ...string) {
	a.t.Helper()
	view := a.explorer.View()
	for _, expected := range expectedStrings {
		if !strings.Contains(view, expected) {
			a.t.Errorf("view does not contain %q\nView:\n%s", expected, view)
		}
	}
}

// CurrentType checks if the current GraphQL type matches the expected type
func (a *Assert) CurrentType(expected navigation.GQLType) {
	a.t.Helper()
	// Check if the type tab is visible in the view
	view := a.explorer.View()
	if !strings.Contains(view, string(expected)) {
		a.t.Errorf("current type is not %q\nView:\n%s", expected, view)
	}
}
