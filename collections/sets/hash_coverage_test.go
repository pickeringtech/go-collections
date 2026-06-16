package sets_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/sets"
)

func TestHash_MutableInPlaceOps(t *testing.T) {
	s := sets.NewHash(1)

	s.AddManyInPlace(2, 3)
	if !intSlicesEqual(s.AsSlice(), []int{1, 2, 3}) {
		t.Errorf("AddManyInPlace = %v, want {1,2,3}", s.AsSlice())
	}

	other := sets.NewHash(3, 4, 5)
	s.UnionInPlace(other)
	if !intSlicesEqual(s.AsSlice(), []int{1, 2, 3, 4, 5}) {
		t.Errorf("UnionInPlace = %v, want {1,2,3,4,5}", s.AsSlice())
	}

	s.RemoveManyInPlace(1, 2)
	if !intSlicesEqual(s.AsSlice(), []int{3, 4, 5}) {
		t.Errorf("RemoveManyInPlace = %v, want {3,4,5}", s.AsSlice())
	}

	s.DifferenceInPlace(sets.NewHash(4))
	if !intSlicesEqual(s.AsSlice(), []int{3, 5}) {
		t.Errorf("DifferenceInPlace = %v, want {3,5}", s.AsSlice())
	}

	s.IntersectionInPlace(sets.NewHash(5, 9))
	if !intSlicesEqual(s.AsSlice(), []int{5}) {
		t.Errorf("IntersectionInPlace = %v, want {5}", s.AsSlice())
	}
}

func TestHash_SearchEdgeCases(t *testing.T) {
	s := sets.NewHash(2, 4, 6)

	_, ok := s.Find(func(e int) bool { return e%2 == 1 })
	if ok {
		t.Error("Find(odd) found a match, want none")
	}

	if s.AnyMatch(func(e int) bool { return e%2 == 1 }) {
		t.Error("AnyMatch(odd) = true, want false")
	}
}

func TestHash_EqualsDifferentLength(t *testing.T) {
	s := sets.NewHash(1, 2, 3)
	if s.Equals(sets.NewHash(1, 2)) {
		t.Error("Equals (different length) = true, want false")
	}
}
