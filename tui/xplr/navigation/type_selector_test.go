package navigation

import (
	"testing"

	"github.com/matryer/is"
)

func TestTypeSelector_newTypeSelector(t *testing.T) {
	is := is.New(t)
	ts := newTypeSelector()
	is.Equal(ts.Current(), QueryType)
	is.Equal(len(ts.All()), 9)
}

func TestTypeSelector_Set(t *testing.T) {
	is := is.New(t)
	ts := newTypeSelector()
	ts.Set(MutationType)
	is.Equal(ts.Current(), MutationType)
}

func TestTypeSelector_Next(t *testing.T) {
	is := is.New(t)
	ts := newTypeSelector()
	is.Equal(ts.Current(), QueryType)

	next := ts.Next()
	is.Equal(next, MutationType)
	is.Equal(ts.Current(), MutationType)

	// Cycle through all types
	for i := 0; i < 7; i++ {
		ts.Next()
	}
	is.Equal(ts.Current(), DirectiveType)

	// Test wraparound
	next = ts.Next()
	is.Equal(next, QueryType)
}

func TestTypeSelector_Previous(t *testing.T) {
	is := is.New(t)
	ts := newTypeSelector()
	is.Equal(ts.Current(), QueryType)

	// Test wraparound at beginning
	prev := ts.Previous()
	is.Equal(prev, DirectiveType)
	is.Equal(ts.Current(), DirectiveType)

	// Move back one more
	prev = ts.Previous()
	is.Equal(prev, UnionType)
}

func TestTypeSelector_All(t *testing.T) {
	is := is.New(t)
	ts := newTypeSelector()
	all := ts.All()

	expected := []GQLType{
		QueryType, MutationType, ObjectType, InputType,
		EnumType, ScalarType, InterfaceType, UnionType, DirectiveType,
	}

	is.Equal(len(all), len(expected))

	for i, typ := range expected {
		is.Equal(all[i], typ)
	}
}

func TestTypeSelector_CurrentIndex(t *testing.T) {
	is := is.New(t)
	ts := newTypeSelector()

	// Test first type
	idx := ts.currentIndex()
	is.Equal(idx, 0)

	// Test middle type
	ts.Set(EnumType)
	idx = ts.currentIndex()
	is.Equal(idx, 4)

	// Test last type
	ts.Set(DirectiveType)
	idx = ts.currentIndex()
	is.Equal(idx, 8)
}
