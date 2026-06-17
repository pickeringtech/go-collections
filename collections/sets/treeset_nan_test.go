package sets_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sets"
)

// TreeSet is backed by dicts.Tree, so it inherits the tree's NaN contract: all
// NaN values are a single element that sorts as the minimum, ahead of every
// non-NaN value including -Inf. These tests pin that contract for the set API.
// See dicts/tree_nan_test.go for the full rationale.

// TestTreeSet_NaNElement_RoundTrips verifies basic add/contains/remove
// consistency for a NaN element.
func TestTreeSet_NaNElement_RoundTrips(t *testing.T) {
	nan := math.NaN()
	s := sets.NewTreeSet[float64]()

	s.AddInPlace(nan)
	if !s.Contains(nan) {
		t.Fatal("Contains(NaN) = false, want true")
	}
	if s.Length() != 1 {
		t.Fatalf("Length() = %d, want 1", s.Length())
	}

	if !s.RemoveInPlace(nan) {
		t.Fatal("RemoveInPlace(NaN) = false, want true")
	}
	if s.Contains(nan) {
		t.Fatal("Contains(NaN) = true after removal, want false")
	}
}

// TestTreeSet_NaNElements_Deduplicate verifies that distinct NaN bit patterns
// collapse to a single set element — unlike a native map-backed set, which
// would treat every NaN as unique.
func TestTreeSet_NaNElements_Deduplicate(t *testing.T) {
	nan1 := math.NaN()
	nan2 := math.Inf(1) - math.Inf(1)

	s := sets.NewTreeSet(nan1, nan2, 1.0, 2.0)
	if s.Length() != 3 {
		t.Fatalf("Length() = %d, want 3 (the two NaNs are one element)", s.Length())
	}
}

// TestTreeSet_NaNElement_SortsAsMinimum verifies NaN leads the sorted order and
// is reported as Min.
func TestTreeSet_NaNElement_SortsAsMinimum(t *testing.T) {
	nan := math.NaN()
	negInf := math.Inf(-1)

	s := sets.NewTreeSet(2.0, nan, negInf, 1.0)

	minElem, ok := s.Min()
	if !ok || !math.IsNaN(minElem) {
		t.Fatalf("Min() = %v (ok=%v), want NaN", minElem, ok)
	}

	ordered := s.AsSlice()
	want := []float64{nan, negInf, 1.0, 2.0}
	if len(ordered) != len(want) {
		t.Fatalf("AsSlice() = %v, want %v", ordered, want)
	}
	if !math.IsNaN(ordered[0]) {
		t.Fatalf("AsSlice()[0] = %v, want NaN", ordered[0])
	}
	for i := 1; i < len(want); i++ {
		if ordered[i] != want[i] {
			t.Fatalf("AsSlice()[%d] = %v, want %v", i, ordered[i], want[i])
		}
	}
}
