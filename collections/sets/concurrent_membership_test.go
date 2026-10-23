package sets_test

import (
	"fmt"
	"math"
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sets"
)

// describeFloatSet renders a set's elements in a stable, NaN-safe form.
func describeFloatSet(s sets.Set[float64]) string {
	elements := s.AsSlice()
	sort.Float64s(elements) // NaN sorts first
	return fmt.Sprint(elements)
}

// TestConcurrentSet_BinaryOpsUseOtherSetMembership verifies that every binary
// set operation on a concurrent set decides membership in the other set the
// same way its non-concurrent sibling does — through the other set's own
// Contains — so the two agree even when the other set's equality is not ==.
// ConcurrentHash used to snapshot other.AsMap() and look elements up natively,
// so a NaN that a TreeSet operand contains was treated as absent; the concurrent
// tree sets converted other into a TreeSet, so a NaN that a Hash operand cannot
// find was treated as present.
func TestConcurrentSet_BinaryOpsUseOtherSetMembership(t *testing.T) {
	nan := math.NaN()
	ops := []struct {
		name string
		run  func(receiver sets.MutableSet[float64], other sets.Set[float64]) string
	}{
		{"Union", func(r sets.MutableSet[float64], o sets.Set[float64]) string { return describeFloatSet(r.Union(o)) }},
		{"UnionInPlace", func(r sets.MutableSet[float64], o sets.Set[float64]) string {
			r.UnionInPlace(o)
			return describeFloatSet(r)
		}},
		{"Difference", func(r sets.MutableSet[float64], o sets.Set[float64]) string {
			return describeFloatSet(r.Difference(o))
		}},
		{"DifferenceInPlace", func(r sets.MutableSet[float64], o sets.Set[float64]) string {
			r.DifferenceInPlace(o)
			return describeFloatSet(r)
		}},
		{"Intersection", func(r sets.MutableSet[float64], o sets.Set[float64]) string {
			return describeFloatSet(r.Intersection(o))
		}},
		{"IntersectionInPlace", func(r sets.MutableSet[float64], o sets.Set[float64]) string {
			r.IntersectionInPlace(o)
			return describeFloatSet(r)
		}},
		{"IsSubsetOf", func(r sets.MutableSet[float64], o sets.Set[float64]) string { return fmt.Sprint(r.IsSubsetOf(o)) }},
		{"IsSupersetOf", func(r sets.MutableSet[float64], o sets.Set[float64]) string {
			return fmt.Sprint(r.IsSupersetOf(o))
		}},
		{"IsDisjoint", func(r sets.MutableSet[float64], o sets.Set[float64]) string { return fmt.Sprint(r.IsDisjoint(o)) }},
		{"Equals", func(r sets.MutableSet[float64], o sets.Set[float64]) string { return fmt.Sprint(r.Equals(o)) }},
	}
	type setMaker func(elements ...float64) sets.MutableSet[float64]
	hashFamily := map[string]setMaker{
		"ConcurrentHash":   func(e ...float64) sets.MutableSet[float64] { return sets.NewConcurrentHash(e...) },
		"ConcurrentHashRW": func(e ...float64) sets.MutableSet[float64] { return sets.NewConcurrentHashRW(e...) },
	}
	treeFamily := map[string]setMaker{
		"ConcurrentTreeSet":   func(e ...float64) sets.MutableSet[float64] { return sets.NewConcurrentTreeSet(e...) },
		"ConcurrentTreeSetRW": func(e ...float64) sets.MutableSet[float64] { return sets.NewConcurrentTreeSetRW(e...) },
	}
	families := []struct {
		name       string
		sibling    setMaker
		concurrent map[string]setMaker
		other      func() sets.Set[float64]
	}{
		{
			name:       "hash receiver, TreeSet operand holding NaN",
			sibling:    func(e ...float64) sets.MutableSet[float64] { return sets.NewHash(e...) },
			concurrent: hashFamily,
			other:      func() sets.Set[float64] { return sets.NewTreeSet(nan, 1.0) },
		},
		{
			name:       "tree receiver, Hash operand holding NaN",
			sibling:    func(e ...float64) sets.MutableSet[float64] { return sets.NewTreeSet(e...) },
			concurrent: treeFamily,
			other:      func() sets.Set[float64] { return sets.NewHash(nan, 1.0) },
		},
	}
	receivers := [][]float64{{nan, 1, 2}, {nan, 1}, {nan, 2}}
	for _, family := range families {
		for _, elements := range receivers {
			for _, op := range ops {
				want := op.run(family.sibling(elements...), family.other())
				for name, newSet := range family.concurrent {
					t.Run(fmt.Sprintf("%s/%v/%s/%s", family.name, elements, op.name, name), func(t *testing.T) {
						got := op.run(newSet(elements...), family.other())
						if got != want {
							t.Errorf("%s() = %s, want %s (the non-concurrent sibling's result)", op.name, got, want)
						}
					})
				}
			}
		}
	}
}

// TestConcurrentHash_IntersectionLengthWithNaNTreeSet pins the reported case:
// intersecting a set holding NaN and 1 with a TreeSet holding the same elements
// keeps both, as Hash does.
func TestConcurrentHash_IntersectionLengthWithNaNTreeSet(t *testing.T) {
	tests := []struct {
		name string
		set  sets.Set[float64]
		want int
	}{
		{name: "Hash", set: sets.NewHash(math.NaN(), 1.0), want: 2},
		{name: "ConcurrentHash", set: sets.NewConcurrentHash(math.NaN(), 1.0), want: 2},
		{name: "ConcurrentHashRW", set: sets.NewConcurrentHashRW(math.NaN(), 1.0), want: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.set.Intersection(sets.NewTreeSet(math.NaN(), 1.0)).Length()
			if got != tt.want {
				t.Errorf("Intersection(TreeSet(NaN, 1)).Length() = %d, want %d", got, tt.want)
			}
		})
	}
}
