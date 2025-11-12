package navigation

import (
	"testing"
)

func TestTypeSelector_NewTypeSelector(t *testing.T) {
	ts := NewTypeSelector()
	if ts.Current() != QueryType {
		t.Errorf("expected default type to be QueryType, got %v", ts.Current())
	}
	if len(ts.All()) != 9 {
		t.Errorf("expected 9 types, got %d", len(ts.All()))
	}
}

func TestTypeSelector_Set(t *testing.T) {
	ts := NewTypeSelector()
	ts.Set(MutationType)
	if ts.Current() != MutationType {
		t.Errorf("expected current type to be MutationType, got %v", ts.Current())
	}
}

func TestTypeSelector_Next(t *testing.T) {
	ts := NewTypeSelector()
	if ts.Current() != QueryType {
		t.Errorf("expected initial type to be QueryType, got %v", ts.Current())
	}

	next := ts.Next()
	if next != MutationType {
		t.Errorf("expected next type to be MutationType, got %v", next)
	}
	if ts.Current() != MutationType {
		t.Errorf("expected current type to be MutationType, got %v", ts.Current())
	}

	// Cycle through all types
	for i := 0; i < 7; i++ {
		ts.Next()
	}
	if ts.Current() != DirectiveType {
		t.Errorf("expected current type to be DirectiveType, got %v", ts.Current())
	}

	// Test wraparound
	next = ts.Next()
	if next != QueryType {
		t.Errorf("expected wraparound to QueryType, got %v", next)
	}
}

func TestTypeSelector_Previous(t *testing.T) {
	ts := NewTypeSelector()
	if ts.Current() != QueryType {
		t.Errorf("expected initial type to be QueryType, got %v", ts.Current())
	}

	// Test wraparound at beginning
	prev := ts.Previous()
	if prev != DirectiveType {
		t.Errorf("expected wraparound to DirectiveType, got %v", prev)
	}
	if ts.Current() != DirectiveType {
		t.Errorf("expected current type to be DirectiveType, got %v", ts.Current())
	}

	// Move back one more
	prev = ts.Previous()
	if prev != UnionType {
		t.Errorf("expected previous type to be UnionType, got %v", prev)
	}
}

func TestTypeSelector_All(t *testing.T) {
	ts := NewTypeSelector()
	all := ts.All()

	expected := []GQLType{
		QueryType, MutationType, ObjectType, InputType,
		EnumType, ScalarType, InterfaceType, UnionType, DirectiveType,
	}

	if len(all) != len(expected) {
		t.Errorf("expected %d types, got %d", len(expected), len(all))
	}

	for i, typ := range expected {
		if all[i] != typ {
			t.Errorf("expected type at index %d to be %v, got %v", i, typ, all[i])
		}
	}
}

func TestTypeSelector_CurrentIndex(t *testing.T) {
	ts := NewTypeSelector()

	// Test first type
	idx := ts.currentIndex()
	if idx != 0 {
		t.Errorf("expected index 0 for QueryType, got %d", idx)
	}

	// Test middle type
	ts.Set(EnumType)
	idx = ts.currentIndex()
	if idx != 4 {
		t.Errorf("expected index 4 for EnumType, got %d", idx)
	}

	// Test last type
	ts.Set(DirectiveType)
	idx = ts.currentIndex()
	if idx != 8 {
		t.Errorf("expected index 8 for DirectiveType, got %d", idx)
	}
}
