package xplr

import (
	"slices"

	"github.com/tonysyu/gqlxp/tui/xplr/components"
)

// favoritableItem wraps a ListItem and adds a favorite indicator to the title
type favoritableItem struct {
	wrapped    components.ListItem
	isFavorite bool
}

// wrapItemsWithFavorites wraps items to add favorite indicators
// Only used for top-level panels - checks RefName() against favorites (field names)
func wrapItemsWithFavorites(items []components.ListItem, favorites []string, isTopLevel bool) []components.ListItem {
	wrapped := make([]components.ListItem, len(items))
	for i, item := range items {
		var isFavorite bool
		if isTopLevel {
			isFavorite = slices.Contains(favorites, item.RefName())
		}
		wrapped[i] = newFavoritableItem(item, isFavorite)
	}
	return wrapped
}

func newFavoritableItem(item components.ListItem, isFavorite bool) components.ListItem {
	return &favoritableItem{
		wrapped:    item,
		isFavorite: isFavorite,
	}
}

func (f *favoritableItem) Title() string {
	title := f.wrapped.Title()
	if f.isFavorite {
		return "★ " + title
	}
	return title
}

func (f *favoritableItem) Description() string {
	return f.wrapped.Description()
}

func (f *favoritableItem) FilterValue() string {
	return f.wrapped.FilterValue()
}

func (f *favoritableItem) TypeName() string {
	return f.wrapped.TypeName()
}

func (f *favoritableItem) RefName() string {
	return f.wrapped.RefName()
}

func (f *favoritableItem) Details() string {
	return f.wrapped.Details()
}

func (f *favoritableItem) OpenPanel() (*components.Panel, bool) {
	panel, ok := f.wrapped.OpenPanel()
	if ok && panel != nil && f.isFavorite {
		// Add favorite indicator to panel title
		panel.SetTitle("★ " + panel.Title())
	}
	return panel, ok
}

// unwrapFavoritableItem extracts the original item from a favoritableItem wrapper
func unwrapFavoritableItem(item components.ListItem) components.ListItem {
	if f, ok := item.(*favoritableItem); ok {
		return f.wrapped
	}
	return item
}
