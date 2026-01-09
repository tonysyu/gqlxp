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
	keyNextType  = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'}'}}
	keyPrevType  = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'{'}}
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
	nav      *Navigator
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
	nav      *Navigator
}

// Navigator provides methods for navigating the TUI
type Navigator struct {
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
		h.explorer.model.SetOverlayStyle(lipgloss.NewStyle())
	}
}

// WithSelection applies a selection target to the model
func WithSelection(typeName, fieldName string) Option {
	return func(h *Harness) {
		target := xplr.SelectionTarget{
			TypeName:  typeName,
			FieldName: fieldName,
		}
		h.explorer.model.ApplySelection(target)
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
	navigator := &Navigator{explorer: explorer}
	h := &Harness{
		explorer: explorer,
		assert:   &Assert{explorer: explorer, t: t, is: is.New(t)},
		overlay:  &Overlay{explorer: explorer, nav: navigator},
		nav:      navigator,
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

// CurrentType returns the currently selected GraphQL type
func (e *Explorer) CurrentType() navigation.GQLType {
	return navigation.GQLType(e.model.CurrentType())
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
	o.nav.GoToGqlType(gqlType)
	o.Open()
}

// ============================================================================
// Component Helpers - Navigator
// ============================================================================

// NextPanel navigates forward to the next panel (Tab key)
func (n *Navigator) NextPanel() {
	n.explorer.Update(keyNextPanel)
}

// PrevPanel navigates backward to the previous panel (Shift+Tab key)
func (n *Navigator) PrevPanel() {
	n.explorer.Update(keyPrevPanel)
}

// NextGqlType cycles to the next GraphQL type (} key)
func (n *Navigator) NextGqlType() {
	n.explorer.Update(keyNextType)
}

// PrevGqlType cycles to the previous GraphQL type ({ key)
func (n *Navigator) PrevGqlType() {
	n.explorer.Update(keyPrevType)
}

// GoToGqlType switches directly to a specific GraphQL type
func (n *Navigator) GoToGqlType(gqlType navigation.GQLType) {
	n.explorer.model.SwitchToType(string(gqlType))
}

// SelectItemAtIndex moves the cursor to the item at the given index (0-based)
func (n *Navigator) SelectItemAtIndex(idx int) {
	// Move cursor to top first (safety measure)
	for i := 0; i < 100; i++ {
		n.explorer.Update(keyUp)
	}
	// Move down to the desired index
	for i := 0; i < idx; i++ {
		n.explorer.Update(keyDown)
	}
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
	if !a.explorer.model.Overlay().IsActive() {
		a.t.Error("overlay is not visible")
	}
}

// OverlayContains checks if the normalized overlay view contains the expected text
// Both the view and expected text are normalized (whitespace trimmed, excess spaces removed)
// before comparison. This is useful for comparing content structure without exact formatting.
func (a *Assert) OverlayContains(expected string) {
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
	if expected != a.explorer.CurrentType() {
		a.t.Errorf("current type is not %q", expected)
	}
}
