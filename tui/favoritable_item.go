package tui

import (
	"slices"

	"github.com/tonysyu/gqlxp/tui/components"
)

// favoritableItem wraps a ListItem and adds a favorite indicator to the title
type favoritableItem struct {
	wrapped    components.ListItem
	isFavorite bool
}

// wrapItemsWithFavorites wraps items to add favorite indicators
func wrapItemsWithFavorites(items []components.ListItem, favorites []string) []components.ListItem {
	wrapped := make([]components.ListItem, len(items))
	for i, item := range items {
		isFavorite := slices.Contains(favorites, item.TypeName())
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
	return f.wrapped.OpenPanel()
}
