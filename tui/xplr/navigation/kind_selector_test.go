package navigation

import (
	"testing"

	"github.com/matryer/is"
)

func TestKindSelector_newKindSelector(t *testing.T) {
	is := is.New(t)
	ts := newKindSelector()
	is.Equal(ts.Current(), QueryKind)
	is.Equal(len(ts.All()), 10) // Updated to include SearchKind
}

func TestKindSelector_Set(t *testing.T) {
	is := is.New(t)
	ts := newKindSelector()
	ts = ts.Set(MutationKind)
	is.Equal(ts.Current(), MutationKind)
}

func TestKindSelector_Next(t *testing.T) {
	is := is.New(t)
	ts := newKindSelector()
	is.Equal(ts.Current(), QueryKind)

	var next GQLKind
	ts, next = ts.Next()
	is.Equal(next, MutationKind)
	is.Equal(ts.Current(), MutationKind)

	// Cycle through all kinds
	for i := 0; i < 7; i++ {
		ts, _ = ts.Next()
	}
	is.Equal(ts.Current(), DirectiveKind)

	// Next should be SearchKind
	ts, next = ts.Next()
	is.Equal(next, SearchKind)

	// Test wraparound
	_, next = ts.Next()
	is.Equal(next, QueryKind)
}

func TestKindSelector_Previous(t *testing.T) {
	is := is.New(t)
	ts := newKindSelector()
	is.Equal(ts.Current(), QueryKind)

	// Test wraparound at beginning - should go to SearchKind (last in list)
	var prev GQLKind
	ts, prev = ts.Previous()
	is.Equal(prev, SearchKind)
	is.Equal(ts.Current(), SearchKind)

	// Previous from SearchKind should be DirectiveKind
	ts, prev = ts.Previous()
	is.Equal(prev, DirectiveKind)

	// Move back one more
	_, prev = ts.Previous()
	is.Equal(prev, UnionKind)
}

func TestKindSelector_All(t *testing.T) {
	is := is.New(t)
	ts := newKindSelector()
	all := ts.All()

	expected := []GQLKind{
		QueryKind, MutationKind, ObjectKind, InputKind,
		EnumKind, ScalarKind, InterfaceKind, UnionKind, DirectiveKind, SearchKind,
	}

	is.Equal(len(all), len(expected))

	for i, typ := range expected {
		is.Equal(all[i], typ)
	}
}

func TestKindSelector_CurrentIndex(t *testing.T) {
	is := is.New(t)
	ts := newKindSelector()

	// Test first kind
	idx := ts.currentIndex()
	is.Equal(idx, 0)

	// Test middle kind
	ts = ts.Set(EnumKind)
	idx = ts.currentIndex()
	is.Equal(idx, 4)

	// Test DirectiveKind
	ts = ts.Set(DirectiveKind)
	idx = ts.currentIndex()
	is.Equal(idx, 8)

	// Test last kind (SearchKind)
	ts = ts.Set(SearchKind)
	idx = ts.currentIndex()
	is.Equal(idx, 9)
}
