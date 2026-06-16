package sets_test

import (
	"sort"
	"sync"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sets"
)

// sortedInts returns a sorted copy of the given slice for deterministic comparison.
func sortedInts(in []int) []int {
	out := make([]int, len(in))
	copy(out, in)
	sort.Ints(out)
	return out
}

func intSlicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	sa := sortedInts(a)
	sb := sortedInts(b)
	for i := range sa {
		if sa[i] != sb[i] {
			return false
		}
	}
	return true
}

// concurrentFactory describes how to build a fresh MutableSet of a particular
// concurrent implementation, so the same table-driven tests run against both.
type concurrentFactory struct {
	name string
	make func(values ...int) sets.MutableSet[int]
}

func concurrentFactories() []concurrentFactory {
	return []concurrentFactory{
		{
			name: "ConcurrentHash",
			make: func(values ...int) sets.MutableSet[int] {
				return sets.NewConcurrentHash(values...)
			},
		},
		{
			name: "ConcurrentHashRW",
			make: func(values ...int) sets.MutableSet[int] {
				return sets.NewConcurrentHashRW(values...)
			},
		},
	}
}

func TestConcurrent_IndexableAndIterable(t *testing.T) {
	for _, f := range concurrentFactories() {
		t.Run(f.name, func(t *testing.T) {
			s := f.make(1, 2, 3)

			if !s.Contains(2) {
				t.Error("Contains(2) = false, want true")
			}
			if s.Contains(99) {
				t.Error("Contains(99) = true, want false")
			}
			if s.Length() != 3 {
				t.Errorf("Length() = %d, want 3", s.Length())
			}
			if s.IsEmpty() {
				t.Error("IsEmpty() = true, want false")
			}

			empty := f.make()
			if !empty.IsEmpty() {
				t.Error("IsEmpty() = false on empty set, want true")
			}

			seen := make([]int, 0, 3)
			s.ForEach(func(e int) {
				seen = append(seen, e)
			})
			if !intSlicesEqual(seen, []int{1, 2, 3}) {
				t.Errorf("ForEach visited %v, want {1,2,3}", seen)
			}
		})
	}
}

func TestConcurrent_FilterAndFilterInPlace(t *testing.T) {
	for _, f := range concurrentFactories() {
		t.Run(f.name, func(t *testing.T) {
			s := f.make(1, 2, 3, 4)
			filtered := s.Filter(func(e int) bool { return e%2 == 0 })
			if !intSlicesEqual(filtered.AsSlice(), []int{2, 4}) {
				t.Errorf("Filter = %v, want {2,4}", filtered.AsSlice())
			}
			// Original unchanged.
			if s.Length() != 4 {
				t.Errorf("Filter mutated receiver, Length() = %d, want 4", s.Length())
			}

			s.FilterInPlace(func(e int) bool { return e > 2 })
			if !intSlicesEqual(s.AsSlice(), []int{3, 4}) {
				t.Errorf("FilterInPlace = %v, want {3,4}", s.AsSlice())
			}
		})
	}
}

func TestConcurrent_Searchable(t *testing.T) {
	for _, f := range concurrentFactories() {
		t.Run(f.name, func(t *testing.T) {
			s := f.make(2, 4, 6)

			found, ok := s.Find(func(e int) bool { return e > 3 })
			if !ok {
				t.Error("Find found nothing, want a match")
			}
			if found != 4 && found != 6 {
				t.Errorf("Find = %d, want 4 or 6", found)
			}

			_, ok = s.Find(func(e int) bool { return e > 100 })
			if ok {
				t.Error("Find found a match, want none")
			}

			if !s.AllMatch(func(e int) bool { return e%2 == 0 }) {
				t.Error("AllMatch(even) = false, want true")
			}
			if s.AllMatch(func(e int) bool { return e > 4 }) {
				t.Error("AllMatch(>4) = true, want false")
			}

			if !s.AnyMatch(func(e int) bool { return e == 4 }) {
				t.Error("AnyMatch(==4) = false, want true")
			}
			if s.AnyMatch(func(e int) bool { return e == 5 }) {
				t.Error("AnyMatch(==5) = true, want false")
			}
		})
	}
}

func TestConcurrent_Convertible(t *testing.T) {
	for _, f := range concurrentFactories() {
		t.Run(f.name, func(t *testing.T) {
			s := f.make(1, 2, 3)

			if !intSlicesEqual(s.AsSlice(), []int{1, 2, 3}) {
				t.Errorf("AsSlice = %v, want {1,2,3}", s.AsSlice())
			}

			m := s.AsMap()
			if len(m) != 3 {
				t.Errorf("AsMap len = %d, want 3", len(m))
			}
			for _, want := range []int{1, 2, 3} {
				_, ok := m[want]
				if !ok {
					t.Errorf("AsMap missing %d", want)
				}
			}
		})
	}
}

func TestConcurrent_InsertableImmutable(t *testing.T) {
	for _, f := range concurrentFactories() {
		t.Run(f.name, func(t *testing.T) {
			s := f.make(1, 2)
			other := f.make(2, 3)

			added := s.Add(3)
			if !intSlicesEqual(added.AsSlice(), []int{1, 2, 3}) {
				t.Errorf("Add = %v, want {1,2,3}", added.AsSlice())
			}

			addedMany := s.AddMany(3, 4)
			if !intSlicesEqual(addedMany.AsSlice(), []int{1, 2, 3, 4}) {
				t.Errorf("AddMany = %v, want {1,2,3,4}", addedMany.AsSlice())
			}

			union := s.Union(other)
			if !intSlicesEqual(union.AsSlice(), []int{1, 2, 3}) {
				t.Errorf("Union = %v, want {1,2,3}", union.AsSlice())
			}

			// Receiver unchanged by immutable ops.
			if s.Length() != 2 {
				t.Errorf("immutable op mutated receiver, Length() = %d, want 2", s.Length())
			}
		})
	}
}

func TestConcurrent_MutableInsertable(t *testing.T) {
	for _, f := range concurrentFactories() {
		t.Run(f.name, func(t *testing.T) {
			s := f.make(1)

			s.AddInPlace(2)
			if !s.Contains(2) {
				t.Error("AddInPlace(2) did not add element")
			}

			s.AddManyInPlace(3, 4)
			if !intSlicesEqual(s.AsSlice(), []int{1, 2, 3, 4}) {
				t.Errorf("AddManyInPlace = %v, want {1,2,3,4}", s.AsSlice())
			}

			other := f.make(4, 5, 6)
			s.UnionInPlace(other)
			if !intSlicesEqual(s.AsSlice(), []int{1, 2, 3, 4, 5, 6}) {
				t.Errorf("UnionInPlace = %v, want {1,2,3,4,5,6}", s.AsSlice())
			}
		})
	}
}

func TestConcurrent_RemovableImmutable(t *testing.T) {
	for _, f := range concurrentFactories() {
		t.Run(f.name, func(t *testing.T) {
			s := f.make(1, 2, 3, 4)
			other := f.make(2, 4)

			removed := s.Remove(2)
			if !intSlicesEqual(removed.AsSlice(), []int{1, 3, 4}) {
				t.Errorf("Remove = %v, want {1,3,4}", removed.AsSlice())
			}

			removedMany := s.RemoveMany(1, 2)
			if !intSlicesEqual(removedMany.AsSlice(), []int{3, 4}) {
				t.Errorf("RemoveMany = %v, want {3,4}", removedMany.AsSlice())
			}

			diff := s.Difference(other)
			if !intSlicesEqual(diff.AsSlice(), []int{1, 3}) {
				t.Errorf("Difference = %v, want {1,3}", diff.AsSlice())
			}

			if s.Length() != 4 {
				t.Errorf("immutable op mutated receiver, Length() = %d, want 4", s.Length())
			}
		})
	}
}

func TestConcurrent_MutableRemovable(t *testing.T) {
	for _, f := range concurrentFactories() {
		t.Run(f.name, func(t *testing.T) {
			s := f.make(1, 2, 3, 4, 5)

			if !s.RemoveInPlace(3) {
				t.Error("RemoveInPlace(3) = false, want true")
			}
			if s.RemoveInPlace(99) {
				t.Error("RemoveInPlace(99) = true, want false")
			}
			if s.Contains(3) {
				t.Error("RemoveInPlace did not remove 3")
			}

			s.RemoveManyInPlace(1, 2)
			if !intSlicesEqual(s.AsSlice(), []int{4, 5}) {
				t.Errorf("after RemoveManyInPlace = %v, want {4,5}", s.AsSlice())
			}

			other := f.make(5)
			s.DifferenceInPlace(other)
			if !intSlicesEqual(s.AsSlice(), []int{4}) {
				t.Errorf("after DifferenceInPlace = %v, want {4}", s.AsSlice())
			}

			s.Clear()
			if !s.IsEmpty() {
				t.Errorf("after Clear IsEmpty() = false, want true (len %d)", s.Length())
			}
		})
	}
}

func TestConcurrent_SetOperations(t *testing.T) {
	for _, f := range concurrentFactories() {
		t.Run(f.name, func(t *testing.T) {
			s := f.make(1, 2, 3)
			overlap := f.make(2, 3, 4)
			subset := f.make(1, 2)
			disjoint := f.make(7, 8)
			equal := f.make(3, 2, 1)

			inter := s.Intersection(overlap)
			if !intSlicesEqual(inter.AsSlice(), []int{2, 3}) {
				t.Errorf("Intersection = %v, want {2,3}", inter.AsSlice())
			}

			if !subset.IsSubsetOf(s) {
				t.Error("IsSubsetOf = false, want true")
			}
			if s.IsSubsetOf(subset) {
				t.Error("IsSubsetOf = true, want false")
			}

			if !s.IsSupersetOf(subset) {
				t.Error("IsSupersetOf = false, want true")
			}
			if subset.IsSupersetOf(s) {
				t.Error("IsSupersetOf = true, want false")
			}

			if !s.IsDisjoint(disjoint) {
				t.Error("IsDisjoint = false, want true")
			}
			if s.IsDisjoint(overlap) {
				t.Error("IsDisjoint = true, want false")
			}

			if !s.Equals(equal) {
				t.Error("Equals = false, want true")
			}
			if s.Equals(overlap) {
				t.Error("Equals (same length, different elems) = true, want false")
			}
			if s.Equals(subset) {
				t.Error("Equals (different length) = true, want false")
			}
		})
	}
}

func TestConcurrent_IntersectionInPlace(t *testing.T) {
	for _, f := range concurrentFactories() {
		t.Run(f.name, func(t *testing.T) {
			s := f.make(1, 2, 3, 4)
			other := f.make(2, 4, 6)

			s.IntersectionInPlace(other)
			if !intSlicesEqual(s.AsSlice(), []int{2, 4}) {
				t.Errorf("IntersectionInPlace = %v, want {2,4}", s.AsSlice())
			}
		})
	}
}

// TestConcurrent_RaceAccess exercises concurrent reads and writes from many
// goroutines so that the -race detector can flag any unsynchronized access.
func TestConcurrent_RaceAccess(t *testing.T) {
	for _, f := range concurrentFactories() {
		t.Run(f.name, func(t *testing.T) {
			s := f.make()

			const goroutines = 16
			const perGoroutine = 200

			var wg sync.WaitGroup
			wg.Add(goroutines)
			for g := 0; g < goroutines; g++ {
				start := g * perGoroutine
				go func(base int) {
					defer wg.Done()
					for i := 0; i < perGoroutine; i++ {
						v := base + i
						s.AddInPlace(v)
						s.Contains(v)
						s.Length()
						s.IsEmpty()
						_ = s.AsSlice()
						s.ForEach(func(int) {})
						s.RemoveInPlace(v)
					}
				}(start)
			}
			wg.Wait()
		})
	}
}

// TestConcurrent_RaceMixedOperations runs a wider mix of read-only and mutating
// operations concurrently to maximise the surface checked by the race detector.
func TestConcurrent_RaceMixedOperations(t *testing.T) {
	for _, f := range concurrentFactories() {
		t.Run(f.name, func(t *testing.T) {
			s := f.make(1, 2, 3, 4, 5)
			other := f.make(3, 4, 5, 6, 7)

			const goroutines = 12
			var wg sync.WaitGroup
			wg.Add(goroutines)
			for g := 0; g < goroutines; g++ {
				idx := g
				go func(id int) {
					defer wg.Done()
					for i := 0; i < 100; i++ {
						switch id % 6 {
						case 0:
							s.AddManyInPlace(id, id+100)
						case 1:
							s.RemoveManyInPlace(id, id+100)
						case 2:
							_ = s.Union(other)
						case 3:
							_ = s.Intersection(other)
						case 4:
							_ = s.Equals(other)
						case 5:
							s.UnionInPlace(other)
						}
					}
				}(idx)
			}
			wg.Wait()
		})
	}
}
