package navigation

type GQLType string

const (
	QueryType     GQLType = "Query"
	MutationType  GQLType = "Mutation"
	ObjectType    GQLType = "Object"
	InputType     GQLType = "Input"
	EnumType      GQLType = "Enum"
	ScalarType    GQLType = "Scalar"
	InterfaceType GQLType = "Interface"
	UnionType     GQLType = "Union"
	DirectiveType GQLType = "Directive"
)

// TypeSelector manages selection among available GQL types
type TypeSelector struct {
	types    []GQLType
	selected GQLType
}

func NewTypeSelector() *TypeSelector {
	types := []GQLType{
		QueryType, MutationType, ObjectType, InputType,
		EnumType, ScalarType, InterfaceType, UnionType, DirectiveType,
	}
	return &TypeSelector{
		types:    types,
		selected: QueryType,
	}
}

// Current returns currently selected type
func (ts *TypeSelector) Current() GQLType {
	return ts.selected
}

// Set changes selected type
func (ts *TypeSelector) Set(gqlType GQLType) {
	ts.selected = gqlType
}

// Next cycles to next type (with wraparound)
func (ts *TypeSelector) Next() GQLType {
	idx := ts.currentIndex()
	nextIdx := (idx + 1) % len(ts.types)
	ts.selected = ts.types[nextIdx]
	return ts.selected
}

// Previous cycles to previous type (with wraparound)
func (ts *TypeSelector) Previous() GQLType {
	idx := ts.currentIndex()
	prevIdx := (idx - 1 + len(ts.types)) % len(ts.types)
	ts.selected = ts.types[prevIdx]
	return ts.selected
}

// All returns all available types
func (ts *TypeSelector) All() []GQLType {
	return ts.types
}

func (ts *TypeSelector) currentIndex() int {
	for i, t := range ts.types {
		if t == ts.selected {
			return i
		}
	}
	return 0
}
