package navigation

import (
	"testing"

	"github.com/matryer/is"
)

func TestTypeSelector_newTypeSelector(t *testing.T) {
	is := is.New(t)
	ts := newTypeSelector()
	is.Equal(ts.Current(), QueryType)
	is.Equal(len(ts.All()), 10) // Updated to include SearchType
}

func TestTypeSelector_Set(t *testing.T) {
	is := is.New(t)
	ts := newTypeSelector()
	ts = ts.Set(MutationType)
	is.Equal(ts.Current(), MutationType)
}

func TestTypeSelector_Next(t *testing.T) {
	is := is.New(t)
	ts := newTypeSelector()
	is.Equal(ts.Current(), QueryType)

	var next GQLType
	ts, next = ts.Next()
	is.Equal(next, MutationType)
	is.Equal(ts.Current(), MutationType)

	// Cycle through all types
	for i := 0; i < 7; i++ {
		ts, _ = ts.Next()
	}
	is.Equal(ts.Current(), DirectiveType)

	// Next should be SearchType
	ts, next = ts.Next()
	is.Equal(next, SearchType)

	// Test wraparound
	_, next = ts.Next()
	is.Equal(next, QueryType)
}

func TestTypeSelector_Previous(t *testing.T) {
	is := is.New(t)
	ts := newTypeSelector()
	is.Equal(ts.Current(), QueryType)

	// Test wraparound at beginning - should go to SearchType (last in list)
	var prev GQLType
	ts, prev = ts.Previous()
	is.Equal(prev, SearchType)
	is.Equal(ts.Current(), SearchType)

	// Previous from SearchType should be DirectiveType
	ts, prev = ts.Previous()
	is.Equal(prev, DirectiveType)

	// Move back one more
	_, prev = ts.Previous()
	is.Equal(prev, UnionType)
}

func TestTypeSelector_All(t *testing.T) {
	is := is.New(t)
	ts := newTypeSelector()
	all := ts.All()

	expected := []GQLType{
		QueryType, MutationType, ObjectType, InputType,
		EnumType, ScalarType, InterfaceType, UnionType, DirectiveType, SearchType,
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
	ts = ts.Set(EnumType)
	idx = ts.currentIndex()
	is.Equal(idx, 4)

	// Test DirectiveType
	ts = ts.Set(DirectiveType)
	idx = ts.currentIndex()
	is.Equal(idx, 8)

	// Test last type (SearchType)
	ts = ts.Set(SearchType)
	idx = ts.currentIndex()
	is.Equal(idx, 9)
}
