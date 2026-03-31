package navigation

type GQLKind string

const (
	QueryKind     GQLKind = "Query"
	MutationKind  GQLKind = "Mutation"
	ObjectKind    GQLKind = "Object"
	InputKind     GQLKind = "Input"
	EnumKind      GQLKind = "Enum"
	ScalarKind    GQLKind = "Scalar"
	InterfaceKind GQLKind = "Interface"
	UnionKind     GQLKind = "Union"
	DirectiveKind GQLKind = "Directive"
	SearchKind    GQLKind = "Search"
)

// kindSelector manages selection among available GQL kinds
type kindSelector struct {
	kinds    []GQLKind
	selected GQLKind
}

func newKindSelector() kindSelector {
	kinds := []GQLKind{
		QueryKind, MutationKind, ObjectKind, InputKind,
		EnumKind, ScalarKind, InterfaceKind, UnionKind, DirectiveKind, SearchKind,
	}
	return kindSelector{
		kinds:    kinds,
		selected: QueryKind,
	}
}

// Current returns currently selected kind
func (ts kindSelector) Current() GQLKind {
	return ts.selected
}

// Set changes selected kind
func (ts kindSelector) Set(gqlKind GQLKind) kindSelector {
	ts.selected = gqlKind
	return ts
}

// Next cycles to next kind (with wraparound)
func (ts kindSelector) Next() (kindSelector, GQLKind) {
	idx := ts.currentIndex()
	nextIdx := (idx + 1) % len(ts.kinds)
	ts.selected = ts.kinds[nextIdx]
	return ts, ts.selected
}

// Previous cycles to previous kind (with wraparound)
func (ts kindSelector) Previous() (kindSelector, GQLKind) {
	idx := ts.currentIndex()
	prevIdx := (idx - 1 + len(ts.kinds)) % len(ts.kinds)
	ts.selected = ts.kinds[prevIdx]
	return ts, ts.selected
}

// All returns all available kinds
func (ts kindSelector) All() []GQLKind {
	return ts.kinds
}

func (ts kindSelector) currentIndex() int {
	for i, t := range ts.kinds {
		if t == ts.selected {
			return i
		}
	}
	return 0
}
