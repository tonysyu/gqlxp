package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/tonysyu/gqlxp/tui/config"
)

// breadcrumbsModel maintains breadcrumb trail state for navigation
type breadcrumbsModel struct {
	// Stack of breadcrumb titles representing the navigation path
	crumbs []string
	styles config.Styles
}

// newBreadcrumbsModel creates a new breadcrumbs model
func newBreadcrumbsModel(styles config.Styles) breadcrumbsModel {
	return breadcrumbsModel{
		crumbs: []string{},
		styles: styles,
	}
}

// Push adds a new breadcrumb to the trail
func (b *breadcrumbsModel) Push(title string) {
	b.crumbs = append(b.crumbs, title)
}

// Pop removes the last breadcrumb from the trail
func (b *breadcrumbsModel) Pop() {
	if len(b.crumbs) > 0 {
		b.crumbs = b.crumbs[:b.Len()-1]
	}
}

// Reset clears all breadcrumbs
func (b *breadcrumbsModel) Reset() {
	b.crumbs = []string{}
}

// Len returns the number of breadcrumbs
func (b *breadcrumbsModel) Len() int {
	return len(b.crumbs)
}

// Render creates the breadcrumb trail view
func (b *breadcrumbsModel) Render() string {
	if len(b.crumbs) == 0 {
		return ""
	}

	// Build breadcrumb parts with separators
	var parts []string
	separator := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(" > ")

	for i, crumb := range b.crumbs {
		if i > 0 {
			parts = append(parts, separator)
		}
		parts = append(parts, crumb)
	}

	breadcrumbText := lipgloss.JoinHorizontal(lipgloss.Left, parts...)
	return b.styles.Breadcrumbs.Render(breadcrumbText)
}
