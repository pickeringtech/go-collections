package sets_test

import (
	"reflect"
	"sync"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sets"
)

// sortedSetFactory names a MutableSortedSet implementation and a constructor for
// it, so the shared behavioural suite can run against every backend.
type sortedSetFactory struct {
	name string
	make func(elements ...int) sets.MutableSortedSet[int]
}

func sortedSetFactories() []sortedSetFactory {
	return []sortedSetFactory{
		{"TreeSet", func(e ...int) sets.MutableSortedSet[int] { return sets.NewTreeSet(e...) }},
		{"ConcurrentTreeSet", func(e ...int) sets.MutableSortedSet[int] { return sets.NewConcurrentTreeSet(e...) }},
		{"ConcurrentTreeSetRW", func(e ...int) sets.MutableSortedSet[int] { return sets.NewConcurrentTreeSetRW(e...) }},
	}
}

// sortedSetSample has sorted order 1,3,4,6,7,9 — with interior gaps to probe.
func sortedSetSample() []int {
	return []int{6, 3, 9, 1, 7, 4}
}

func collectSeq(seq func(yield func(int) bool)) []int {
	out := []int{}
	seq(func(e int) bool {
		out = append(out, e)
		return true
	})
	return out
}

func TestTreeSet_MinMax(t *testing.T) {
	for _, f := range sortedSetFactories() {
		t.Run(f.name, func(t *testing.T) {
			empty := f.make()
			if e, ok := empty.Min(); ok || e != 0 {
				t.Errorf("empty Min() = (%d, %v), want (0, false)", e, ok)
			}
			if e, ok := empty.Max(); ok || e != 0 {
				t.Errorf("empty Max() = (%d, %v), want (0, false)", e, ok)
			}

			s := f.make(sortedSetSample()...)
			if e, ok := s.Min(); !ok || e != 1 {
				t.Errorf("Min() = (%d, %v), want (1, true)", e, ok)
			}
			if e, ok := s.Max(); !ok || e != 9 {
				t.Errorf("Max() = (%d, %v), want (9, true)", e, ok)
			}
		})
	}
}

func TestTreeSet_FloorCeiling(t *testing.T) {
	type want struct {
		element int
		ok      bool
	}
	tests := []struct {
		name        string
		arg         int
		floor, ceil want
	}{
		{"exact match", 4, want{4, true}, want{4, true}},
		{"between keys", 5, want{4, true}, want{6, true}},
		{"below min", 0, want{0, false}, want{1, true}},
		{"above max", 10, want{9, true}, want{0, false}},
	}
	for _, f := range sortedSetFactories() {
		for _, tt := range tests {
			t.Run(f.name+"/"+tt.name, func(t *testing.T) {
				s := f.make(sortedSetSample()...)
				if e, ok := s.Floor(tt.arg); e != tt.floor.element || ok != tt.floor.ok {
					t.Errorf("Floor(%d) = (%d, %v), want (%d, %v)", tt.arg, e, ok, tt.floor.element, tt.floor.ok)
				}
				if e, ok := s.Ceiling(tt.arg); e != tt.ceil.element || ok != tt.ceil.ok {
					t.Errorf("Ceiling(%d) = (%d, %v), want (%d, %v)", tt.arg, e, ok, tt.ceil.element, tt.ceil.ok)
				}
			})
		}
	}
}

func TestTreeSet_Range(t *testing.T) {
	tests := []struct {
		name   string
		lo, hi int
		want   []int
	}{
		{"interior inclusive", 3, 7, []int{3, 4, 6, 7}},
		{"bounds outside data", -1, 100, []int{1, 3, 4, 6, 7, 9}},
		{"single element", 4, 4, []int{4}},
		{"empty when lo greater than hi", 7, 3, []int{}},
		{"gap yields empty", 100, 200, []int{}},
	}
	for _, f := range sortedSetFactories() {
		for _, tt := range tests {
			t.Run(f.name+"/"+tt.name, func(t *testing.T) {
				s := f.make(sortedSetSample()...)
				if got := s.Range(tt.lo, tt.hi); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Range(%d, %d) = %v, want %v", tt.lo, tt.hi, got, tt.want)
				}
			})
		}
	}
}

func TestTreeSet_Iterators(t *testing.T) {
	for _, f := range sortedSetFactories() {
		t.Run(f.name, func(t *testing.T) {
			s := f.make(sortedSetSample()...)

			if got := collectSeq(s.All()); !reflect.DeepEqual(got, []int{1, 3, 4, 6, 7, 9}) {
				t.Errorf("All() = %v", got)
			}
			if got := collectSeq(s.Backward()); !reflect.DeepEqual(got, []int{9, 7, 6, 4, 3, 1}) {
				t.Errorf("Backward() = %v", got)
			}
			if got := collectSeq(s.RangeAll(3, 7)); !reflect.DeepEqual(got, []int{3, 4, 6, 7}) {
				t.Errorf("RangeAll(3, 7) = %v", got)
			}
		})
	}
}

func TestTreeSet_IteratorEarlyStop(t *testing.T) {
	for _, f := range sortedSetFactories() {
		t.Run(f.name, func(t *testing.T) {
			s := f.make(sortedSetSample()...)

			countAsc, firstAsc := 0, 0
			s.All()(func(e int) bool { firstAsc = e; countAsc++; return false })
			if countAsc != 1 || firstAsc != 1 {
				t.Errorf("All() early stop = %d items first %d, want 1/1", countAsc, firstAsc)
			}

			countDesc, firstDesc := 0, 0
			s.Backward()(func(e int) bool { firstDesc = e; countDesc++; return false })
			if countDesc != 1 || firstDesc != 9 {
				t.Errorf("Backward() early stop = %d items first %d, want 1/9", countDesc, firstDesc)
			}

			countRange := 0
			s.RangeAll(3, 7)(func(int) bool { countRange++; return false })
			if countRange != 1 {
				t.Errorf("RangeAll() early stop = %d items, want 1", countRange)
			}
		})
	}
}

// TestTreeSet_Surface exercises the full Set/MutableSet surface so the concurrent
// backends' delegating methods are covered with real assertions.
func TestTreeSet_Surface(t *testing.T) {
	for _, f := range sortedSetFactories() {
		t.Run(f.name, func(t *testing.T) {
			s := f.make(sortedSetSample()...)

			if s.IsEmpty() || s.Length() != 6 {
				t.Errorf("IsEmpty=%v Length=%d, want false/6", s.IsEmpty(), s.Length())
			}
			if !s.Contains(7) || s.Contains(2) {
				t.Error("Contains mismatch for 7/2")
			}

			sum := 0
			s.ForEach(func(e int) { sum += e })
			if sum != 30 {
				t.Errorf("ForEach sum = %d, want 30", sum)
			}

			if got := s.AsSlice(); !reflect.DeepEqual(got, []int{1, 3, 4, 6, 7, 9}) {
				t.Errorf("AsSlice() = %v", got)
			}
			wantMap := map[int]struct{}{1: {}, 3: {}, 4: {}, 6: {}, 7: {}, 9: {}}
			if got := s.AsMap(); !reflect.DeepEqual(got, wantMap) {
				t.Errorf("AsMap() = %v", got)
			}

			// Search predicates.
			if !s.AllMatch(func(e int) bool { return e > 0 }) || s.AllMatch(func(e int) bool { return e > 5 }) {
				t.Error("AllMatch mismatch")
			}
			if !s.AnyMatch(func(e int) bool { return e == 6 }) || s.AnyMatch(func(e int) bool { return e == 2 }) {
				t.Error("AnyMatch mismatch")
			}
			if !s.NoneMatch(func(e int) bool { return e == 2 }) || s.NoneMatch(func(e int) bool { return e == 6 }) {
				t.Error("NoneMatch mismatch")
			}
			if e, ok := s.Find(func(e int) bool { return e >= 4 }); !ok || e != 4 {
				t.Errorf("Find(e>=4) = (%d, %v), want (4, true)", e, ok)
			}

			// Filtering.
			filtered := s.Filter(func(e int) bool { return e%2 == 0 })
			if got := filtered.AsSlice(); !reflect.DeepEqual(got, []int{4, 6}) {
				t.Errorf("Filter even = %v, want [4 6]", got)
			}
			if s.Length() != 6 {
				t.Error("Filter mutated the receiver")
			}
		})
	}
}

// TestTreeSet_SetOperations checks the mathematical set operations and their
// in-place twins against a second set.
func TestTreeSet_SetOperations(t *testing.T) {
	for _, f := range sortedSetFactories() {
		t.Run(f.name, func(t *testing.T) {
			a := f.make(1, 2, 3, 4)
			b := f.make(3, 4, 5, 6)

			if got := a.Union(b).AsSlice(); !reflect.DeepEqual(got, []int{1, 2, 3, 4, 5, 6}) {
				t.Errorf("Union = %v", got)
			}
			if got := a.Intersection(b).AsSlice(); !reflect.DeepEqual(got, []int{3, 4}) {
				t.Errorf("Intersection = %v", got)
			}
			if got := a.Difference(b).AsSlice(); !reflect.DeepEqual(got, []int{1, 2}) {
				t.Errorf("Difference = %v", got)
			}

			if a.IsSubsetOf(b) || !f.make(3, 4).IsSubsetOf(a) {
				t.Error("IsSubsetOf mismatch")
			}
			if a.IsSupersetOf(b) || !a.IsSupersetOf(f.make(1, 2)) {
				t.Error("IsSupersetOf mismatch")
			}
			if a.IsDisjoint(b) || !a.IsDisjoint(f.make(7, 8)) {
				t.Error("IsDisjoint mismatch")
			}
			if a.Equals(b) || !a.Equals(f.make(4, 3, 2, 1)) {
				t.Error("Equals mismatch")
			}
			// Different length short-circuit in Equals.
			if a.Equals(f.make(1, 2, 3)) {
				t.Error("Equals true for different-length sets")
			}

			// Immutable ops left the receiver untouched.
			if a.Length() != 4 {
				t.Errorf("set operations mutated the receiver, length = %d", a.Length())
			}

			// In-place mutations.
			c := f.make(1, 2, 3, 4)
			c.IntersectionInPlace(b)
			if got := c.AsSlice(); !reflect.DeepEqual(got, []int{3, 4}) {
				t.Errorf("IntersectionInPlace = %v, want [3 4]", got)
			}
			c.UnionInPlace(f.make(8, 9))
			if got := c.AsSlice(); !reflect.DeepEqual(got, []int{3, 4, 8, 9}) {
				t.Errorf("UnionInPlace = %v", got)
			}
			c.DifferenceInPlace(f.make(8, 9))
			if got := c.AsSlice(); !reflect.DeepEqual(got, []int{3, 4}) {
				t.Errorf("DifferenceInPlace = %v", got)
			}
		})
	}
}

// TestTreeSet_Mutations covers Add/Remove and their in-place twins.
func TestTreeSet_Mutations(t *testing.T) {
	for _, f := range sortedSetFactories() {
		t.Run(f.name, func(t *testing.T) {
			s := f.make(1, 2, 3)

			if got := s.Add(4).AsSlice(); !reflect.DeepEqual(got, []int{1, 2, 3, 4}) {
				t.Errorf("Add(4) = %v", got)
			}
			if got := s.AddMany(4, 5).AsSlice(); !reflect.DeepEqual(got, []int{1, 2, 3, 4, 5}) {
				t.Errorf("AddMany(4,5) = %v", got)
			}
			if got := s.Remove(2).AsSlice(); !reflect.DeepEqual(got, []int{1, 3}) {
				t.Errorf("Remove(2) = %v", got)
			}
			if got := s.RemoveMany(1, 3).AsSlice(); !reflect.DeepEqual(got, []int{2}) {
				t.Errorf("RemoveMany(1,3) = %v", got)
			}
			if s.Length() != 3 {
				t.Error("immutable Add/Remove mutated the receiver")
			}

			s.AddInPlace(4)
			s.AddManyInPlace(5, 6)
			if got := s.AsSlice(); !reflect.DeepEqual(got, []int{1, 2, 3, 4, 5, 6}) {
				t.Errorf("after AddInPlace/AddManyInPlace = %v", got)
			}
			if !s.RemoveInPlace(6) || s.RemoveInPlace(99) {
				t.Error("RemoveInPlace return mismatch")
			}
			s.RemoveManyInPlace(4, 5)
			if got := s.AsSlice(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
				t.Errorf("after RemoveManyInPlace = %v", got)
			}
			s.FilterInPlace(func(e int) bool { return e != 2 })
			if got := s.AsSlice(); !reflect.DeepEqual(got, []int{1, 3}) {
				t.Errorf("after FilterInPlace = %v", got)
			}
			s.Clear()
			if !s.IsEmpty() {
				t.Error("Clear did not empty the set")
			}
		})
	}
}

// TestTreeSet_DifferenceInPlaceSelf guards against mutating the backing tree
// mid-traversal: removing the receiver from itself must empty it cleanly rather
// than corrupt the in-order walk.
func TestTreeSet_DifferenceInPlaceSelf(t *testing.T) {
	s := sets.NewTreeSet(5, 3, 9, 1, 7, 4)
	s.DifferenceInPlace(s)
	if !s.IsEmpty() {
		t.Errorf("DifferenceInPlace(self) left %v, want empty", s.AsSlice())
	}
}

// TestConcurrentTreeSet_ReturnsConcurrentType asserts the thread-safe-in →
// thread-safe-out contract for immutable operations.
func TestConcurrentTreeSet_ReturnsConcurrentType(t *testing.T) {
	cs := sets.NewConcurrentTreeSet(1, 2, 3)
	other := sets.NewTreeSet(2, 3, 4)
	checks := []struct {
		name string
		got  sets.Set[int]
	}{
		{"Add", cs.Add(4)},
		{"AddMany", cs.AddMany(4, 5)},
		{"Remove", cs.Remove(1)},
		{"RemoveMany", cs.RemoveMany(1)},
		{"Filter", cs.Filter(func(int) bool { return true })},
		{"Union", cs.Union(other)},
		{"Intersection", cs.Intersection(other)},
		{"Difference", cs.Difference(other)},
	}
	for _, c := range checks {
		if _, ok := c.got.(*sets.ConcurrentTreeSet[int]); !ok {
			t.Errorf("ConcurrentTreeSet.%s did not return *ConcurrentTreeSet", c.name)
		}
	}

	rw := sets.NewConcurrentTreeSetRW(1, 2, 3)
	otherRW := sets.NewTreeSet(2, 3, 4)
	checksRW := []struct {
		name string
		got  sets.Set[int]
	}{
		{"Add", rw.Add(4)},
		{"AddMany", rw.AddMany(4, 5)},
		{"Remove", rw.Remove(1)},
		{"RemoveMany", rw.RemoveMany(1)},
		{"Filter", rw.Filter(func(int) bool { return true })},
		{"Union", rw.Union(otherRW)},
		{"Intersection", rw.Intersection(otherRW)},
		{"Difference", rw.Difference(otherRW)},
	}
	for _, c := range checksRW {
		if _, ok := c.got.(*sets.ConcurrentTreeSetRW[int]); !ok {
			t.Errorf("ConcurrentTreeSetRW.%s did not return *ConcurrentTreeSetRW", c.name)
		}
	}
}

// TestConcurrentTreeSet_RaceSafety drives both concurrent sets from many
// goroutines; run with -race it asserts there are no data races or panics.
func TestConcurrentTreeSet_RaceSafety(t *testing.T) {
	setsUnderTest := []sets.MutableSortedSet[int]{
		sets.NewConcurrentTreeSet[int](),
		sets.NewConcurrentTreeSetRW[int](),
	}
	for _, s := range setsUnderTest {
		var wg sync.WaitGroup
		for w := 0; w < 8; w++ {
			wg.Add(1)
			go func(base int) {
				defer wg.Done()
				for i := 0; i < 50; i++ {
					element := base*50 + i
					s.AddInPlace(element)
					s.Contains(element)
					s.Min()
					s.Floor(element)
					for range s.All() {
						break
					}
					s.RemoveInPlace(element)
				}
			}(w)
		}
		wg.Wait()
		if !s.IsEmpty() {
			t.Errorf("expected empty set after balanced add/remove, got length %d", s.Length())
		}
	}
}
