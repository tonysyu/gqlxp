package xplr

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	"github.com/tonysyu/gqlxp/tui/xplr/navigation"
	"github.com/tonysyu/gqlxp/utils/text"
)

// View renders the main TUI view
func (m Model) View() string {
	// Build help key bindings
	helpBindings := []key.Binding{
		m.keymap.NextPanel,
		m.keymap.PrevPanel,
		m.keymap.ToggleGQLType,
		m.keymap.ToggleOverlay,
		m.keymap.Quit,
	}

	help := m.help.ShortHelpView(helpBindings)

	// Show overlay if active, and return immediately
	if m.Overlay().IsActive() {
		return m.Overlay().View()
	}

	var views []string
	if m.nav.CurrentPanel() != nil {
		views = append(views, m.nav.CurrentPanel().View())
	}
	if m.nav.NextPanel() != nil {
		views = append(views, m.nav.NextPanel().View())
	}

	navbar := m.renderGQLTypeNavbar()
	breadcrumbs := m.renderBreadcrumbs()
	panels := lipgloss.JoinHorizontal(lipgloss.Top, views...)

	// Add search input if on Search tab
	var mainView string
	if m.nav.CurrentType() == navigation.SearchType {
		searchInput := m.searchInput.View()
		mainView = lipgloss.JoinVertical(0, navbar, breadcrumbs, panels, searchInput, help)
	} else {
		mainView = lipgloss.JoinVertical(0, navbar, breadcrumbs, panels, help)
	}
	return mainView
}

// renderGQLTypeNavbar creates the navbar showing GQL types
func (m *Model) renderGQLTypeNavbar() string {
	allTypes := m.nav.AllTypes()
	if len(allTypes) == 0 {
		return ""
	}

	// Calculate max width per tab to fit all tabs in available width
	// Account for padding in tab styles by using a small buffer
	const styleOverhead = 4 // Approximate horizontal padding per tab
	availableWidth := m.width - (len(allTypes) * styleOverhead)
	if availableWidth < len(allTypes) {
		availableWidth = len(allTypes) // Ensure at least 1 char per tab
	}
	maxTabWidth := availableWidth / len(allTypes)

	var tabs []string
	for _, fieldType := range allTypes {
		var style lipgloss.Style
		if m.nav.CurrentType() == fieldType {
			style = m.Styles.ActiveTab
		} else {
			style = m.Styles.InactiveTab
		}
		// Truncate label to fit within calculated max width
		label := text.Truncate(string(fieldType), maxTabWidth)
		tabs = append(tabs, style.Render(label))
	}

	navbar := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	return m.Styles.Navbar.Render(navbar)
}

// renderBreadcrumbs renders the breadcrumb trail
func (m *Model) renderBreadcrumbs() string {
	crumbs := m.nav.Breadcrumbs()
	if len(crumbs) == 0 {
		return ""
	}

	// Build breadcrumb parts with separators
	var parts []string
	separator := " > "

	for i, crumb := range crumbs {
		if i > 0 {
			parts = append(parts, separator)
		}
		// Apply special color to the last breadcrumb
		if i == len(crumbs)-1 {
			parts = append(parts, m.Styles.CurrentBreadcrumb.Render(crumb))
		} else {
			parts = append(parts, crumb)
		}
	}

	breadcrumbText := lipgloss.JoinHorizontal(lipgloss.Left, parts...)
	return m.Styles.Breadcrumbs.Render(breadcrumbText)
}
