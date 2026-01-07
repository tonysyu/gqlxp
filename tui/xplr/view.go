package xplr

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	"github.com/tonysyu/gqlxp/tui/xplr/navigation"
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
	var tabs []string

	for _, fieldType := range m.nav.AllTypes() {
		var style lipgloss.Style
		if m.nav.CurrentType() == fieldType {
			style = m.Styles.ActiveTab
		} else {
			style = m.Styles.InactiveTab
		}
		tabs = append(tabs, style.Render(string(fieldType)))
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
