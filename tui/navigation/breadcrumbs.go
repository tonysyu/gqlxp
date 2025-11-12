package navigation

// breadcrumbsModel maintains breadcrumb trail state for navigation
type breadcrumbsModel struct {
	// Stack of breadcrumb titles representing the navigation path
	crumbs []string
}

// newBreadcrumbsModel creates a new breadcrumbs model
func newBreadcrumbsModel() breadcrumbsModel {
	return breadcrumbsModel{
		crumbs: []string{},
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

// Get returns a copy of the breadcrumb trail
func (b *breadcrumbsModel) Get() []string {
	if len(b.crumbs) == 0 {
		return nil
	}
	// Return a copy to prevent external modification
	result := make([]string, len(b.crumbs))
	copy(result, b.crumbs)
	return result
}
