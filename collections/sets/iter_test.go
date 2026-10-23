package sets_test

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sets"
)

type setConstructor struct {
	name string
	make func(...int) sets.Set[int]
}

func allSetConstructors() []setConstructor {
	return []setConstructor{
		{"Hash", func(e ...int) sets.Set[int] { return sets.NewHash(e...) }},
		{"ConcurrentHash", func(e ...int) sets.Set[int] { return sets.NewConcurrentHash(e...) }},
		{"ConcurrentHashRW", func(e ...int) sets.Set[int] { return sets.NewConcurrentHashRW(e...) }},
	}
}

func TestSet_All(t *testing.T) {
	for _, ctor := range allSetConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			var got []int
			for v := range ctor.make(3, 1, 2).All() {
				got = append(got, v)
			}
			sort.Ints(got)
			if want := []int{1, 2, 3}; !reflect.DeepEqual(got, want) {
				t.Errorf("All() = %v, want %v", got, want)
			}
		})
	}
}

func TestSet_All_Empty(t *testing.T) {
	for _, ctor := range allSetConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			count := 0
			for range ctor.make().All() {
				count++
			}
			if count != 0 {
				t.Errorf("All() over empty set yielded %d times, want 0", count)
			}
		})
	}
}

func TestSet_All_EarlyBreak(t *testing.T) {
	for _, ctor := range allSetConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			count := 0
			for range ctor.make(1, 2, 3, 4, 5).All() {
				count++
				break
			}
			if count != 1 {
				t.Errorf("All() with early break yielded %d times, want 1", count)
			}
		})
	}
}

func TestFromSeq(t *testing.T) {
	for _, ctor := range allSetConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			source := ctor.make(1, 2, 3)
			got := sets.FromSeq(source.All())

			gotSlice := got.AsSlice()
			sort.Ints(gotSlice)
			if want := []int{1, 2, 3}; !reflect.DeepEqual(gotSlice, want) {
				t.Errorf("FromSeq round-trip = %v, want %v", gotSlice, want)
			}
		})
	}
}

func TestFromSeq_CollapsesDuplicates(t *testing.T) {
	seq := func(yield func(int) bool) {
		for _, v := range []int{1, 1, 2, 2, 3} {
			if !yield(v) {
				return
			}
		}
	}
	got := sets.FromSeq(seq)
	if got.Length() != 3 {
		t.Errorf("FromSeq with duplicates Length() = %d, want 3", got.Length())
	}
}

func ExampleHash_All() {
	set := sets.NewHash(42)
	for v := range set.All() {
		fmt.Println(v)
	}
	// Output: 42
}

func ExampleFromSeq() {
	source := sets.NewHash(1, 2, 3)
	set := sets.FromSeq(source.All())
	fmt.Println(set.Length())
	// Output: 3
}

// mutableSetConstructors builds every MutableSet implementation, plain and
// concurrent; sorted reports whether it also implements SortedSet.
func mutableSetConstructors() []struct {
	name   string
	sorted bool
	make   func(...int) sets.MutableSet[int]
} {
	return []struct {
		name   string
		sorted bool
		make   func(...int) sets.MutableSet[int]
	}{
		{"Hash", false, func(e ...int) sets.MutableSet[int] { return sets.NewHash(e...) }},
		{"ConcurrentHash", false, func(e ...int) sets.MutableSet[int] { return sets.NewConcurrentHash(e...) }},
		{"ConcurrentHashRW", false, func(e ...int) sets.MutableSet[int] { return sets.NewConcurrentHashRW(e...) }},
		{"TreeSet", true, func(e ...int) sets.MutableSet[int] { return sets.NewTreeSet(e...) }},
		{"ConcurrentTreeSet", true, func(e ...int) sets.MutableSet[int] { return sets.NewConcurrentTreeSet(e...) }},
		{"ConcurrentTreeSetRW", true, func(e ...int) sets.MutableSet[int] { return sets.NewConcurrentTreeSetRW(e...) }},
	}
}

// TestSet_RemovingDuringIterationVisitsEveryElement verifies that removing each
// visited element from inside a callback or range loop still visits every
// element exactly once and leaves the set empty. TreeSet used to traverse the
// live nodes of its backing tree, so AVL rebalancing and the successor copy of
// a two-child deletion made it skip elements.
func TestSet_RemovingDuringIterationVisitsEveryElement(t *testing.T) {
	traversals := []struct {
		name   string
		sorted bool
		run    func(s sets.MutableSet[int], visit func(element int))
	}{
		{name: "ForEach", run: func(s sets.MutableSet[int], visit func(int)) { s.ForEach(visit) }},
		{name: "Filter", run: func(s sets.MutableSet[int], visit func(int)) {
			s.Filter(func(e int) bool { visit(e); return true })
		}},
		{name: "FilterInPlace", run: func(s sets.MutableSet[int], visit func(int)) {
			s.FilterInPlace(func(e int) bool { visit(e); return true })
		}},
		{name: "AllMatch", run: func(s sets.MutableSet[int], visit func(int)) {
			s.AllMatch(func(e int) bool { visit(e); return true })
		}},
		{name: "AnyMatch", run: func(s sets.MutableSet[int], visit func(int)) {
			s.AnyMatch(func(e int) bool { visit(e); return false })
		}},
		{name: "NoneMatch", run: func(s sets.MutableSet[int], visit func(int)) {
			s.NoneMatch(func(e int) bool { visit(e); return false })
		}},
		{name: "Find", run: func(s sets.MutableSet[int], visit func(int)) {
			s.Find(func(e int) bool { visit(e); return false })
		}},
		{name: "All", run: func(s sets.MutableSet[int], visit func(int)) {
			for e := range s.All() {
				visit(e)
			}
		}},
		{name: "Backward", sorted: true, run: func(s sets.MutableSet[int], visit func(int)) {
			for e := range s.(sets.SortedSet[int]).Backward() {
				visit(e)
			}
		}},
		{name: "RangeAll", sorted: true, run: func(s sets.MutableSet[int], visit func(int)) {
			for e := range s.(sets.SortedSet[int]).RangeAll(1, 7) {
				visit(e)
			}
		}},
	}
	for _, traversal := range traversals {
		for _, ctor := range mutableSetConstructors() {
			if traversal.sorted && !ctor.sorted {
				continue
			}
			t.Run(traversal.name+"/"+ctor.name, func(t *testing.T) {
				s := ctor.make(1, 2, 3, 4, 5, 6, 7)
				var visited []int
				traversal.run(s, func(element int) {
					visited = append(visited, element)
					s.RemoveInPlace(element)
				})

				sort.Ints(visited)
				if want := []int{1, 2, 3, 4, 5, 6, 7}; !reflect.DeepEqual(visited, want) {
					t.Errorf("%s() removing as it goes visited %v, want %v", traversal.name, visited, want)
				}
				if s.Length() != 0 {
					t.Errorf("%s() removing as it goes left %v", traversal.name, s.AsSlice())
				}
			})
		}
	}
}

// TestSet_IteratorSnapshotsWhenRangingBegins verifies that every iterator reads
// the set when a range over it begins, not when the iterator is created: an
// element added between creating the iterator and ranging it is seen. The
// concurrent sets used to copy their elements when All/Backward/RangeAll was
// called, so a stored iterator was stale.
func TestSet_IteratorSnapshotsWhenRangingBegins(t *testing.T) {
	iterators := []struct {
		name   string
		sorted bool
		create func(s sets.MutableSet[int]) func(yield func(int) bool)
	}{
		{name: "All", create: func(s sets.MutableSet[int]) func(func(int) bool) { return s.All() }},
		{name: "Backward", sorted: true, create: func(s sets.MutableSet[int]) func(func(int) bool) {
			return s.(sets.SortedSet[int]).Backward()
		}},
		{name: "RangeAll", sorted: true, create: func(s sets.MutableSet[int]) func(func(int) bool) {
			return s.(sets.SortedSet[int]).RangeAll(1, 7)
		}},
	}
	for _, it := range iterators {
		for _, ctor := range mutableSetConstructors() {
			if it.sorted && !ctor.sorted {
				continue
			}
			t.Run(it.name+"/"+ctor.name, func(t *testing.T) {
				s := ctor.make(1, 2)
				seq := it.create(s)
				s.AddInPlace(3) // added after the iterator was created

				var visited []int
				for e := range seq {
					visited = append(visited, e)
				}

				sort.Ints(visited)
				if want := []int{1, 2, 3}; !reflect.DeepEqual(visited, want) {
					t.Errorf("%s() created before AddInPlace then ranged = %v, want %v", it.name, visited, want)
				}
			})
		}
	}
}
