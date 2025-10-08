package components

// SimpleItem is a ListItem implementation with arbitrary title and description and no-op Open() function.
type SimpleItem struct {
	Title_       string
	Description_ string
}

// NewSimpleItem creates a new SimpleItem with the given title and description
func NewSimpleItem(title, description string) SimpleItem {
	return SimpleItem{
		Title_:       title,
		Description_: description,
	}
}

func (si SimpleItem) Title() string       { return si.Title_ }
func (si SimpleItem) Description() string { return si.Description_ }
func (si SimpleItem) FilterValue() string { return si.Title_ }
func (si SimpleItem) Details() string {
	if si.Description_ != "" {
		return "# " + si.Title() + "\n\n" + si.Description()
	}
	return "# " + si.Title()
}
func (si SimpleItem) Open() (Panel, bool) { return nil, false }
