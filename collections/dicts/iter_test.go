package dicts_test

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/collections/dicts"
)

type dictConstructor struct {
	name string
	make func(...dicts.Pair[string, int]) dicts.Dict[string, int]
}

func allDictConstructors() []dictConstructor {
	return []dictConstructor{
		{"Hash", func(p ...dicts.Pair[string, int]) dicts.Dict[string, int] { return dicts.NewHash(p...) }},
		{"Tree", func(p ...dicts.Pair[string, int]) dicts.Dict[string, int] { return dicts.NewTree(p...) }},
		{"ConcurrentHash", func(p ...dicts.Pair[string, int]) dicts.Dict[string, int] { return dicts.NewConcurrentHash(p...) }},
		{"ConcurrentHashRW", func(p ...dicts.Pair[string, int]) dicts.Dict[string, int] { return dicts.NewConcurrentHashRW(p...) }},
	}
}

func seedPairs() []dicts.Pair[string, int] {
	return []dicts.Pair[string, int]{
		{Key: "a", Value: 1},
		{Key: "b", Value: 2},
		{Key: "c", Value: 3},
	}
}

func TestDict_All(t *testing.T) {
	for _, ctor := range allDictConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			d := ctor.make(seedPairs()...)

			got := map[string]int{}
			for k, v := range d.All() {
				got[k] = v
			}
			want := map[string]int{"a": 1, "b": 2, "c": 3}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("All() = %v, want %v", got, want)
			}
		})
	}
}

func TestDict_KeysSeq(t *testing.T) {
	for _, ctor := range allDictConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			var got []string
			for k := range ctor.make(seedPairs()...).KeysSeq() {
				got = append(got, k)
			}
			sort.Strings(got)
			if want := []string{"a", "b", "c"}; !reflect.DeepEqual(got, want) {
				t.Errorf("KeysSeq() = %v, want %v", got, want)
			}
		})
	}
}

func TestDict_ValuesSeq(t *testing.T) {
	for _, ctor := range allDictConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			var got []int
			for v := range ctor.make(seedPairs()...).ValuesSeq() {
				got = append(got, v)
			}
			sort.Ints(got)
			if want := []int{1, 2, 3}; !reflect.DeepEqual(got, want) {
				t.Errorf("ValuesSeq() = %v, want %v", got, want)
			}
		})
	}
}

func TestDict_All_Empty(t *testing.T) {
	for _, ctor := range allDictConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			count := 0
			for range ctor.make().All() {
				count++
			}
			if count != 0 {
				t.Errorf("All() over empty dict yielded %d times, want 0", count)
			}
		})
	}
}

func TestDict_All_EarlyBreak(t *testing.T) {
	for _, ctor := range allDictConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			count := 0
			for range ctor.make(seedPairs()...).All() {
				count++
				break
			}
			if count != 1 {
				t.Errorf("All() with early break yielded %d times, want 1", count)
			}
		})
	}
}

func TestDict_KeysSeq_EarlyBreak(t *testing.T) {
	for _, ctor := range allDictConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			count := 0
			for range ctor.make(seedPairs()...).KeysSeq() {
				count++
				break
			}
			if count != 1 {
				t.Errorf("KeysSeq() with early break yielded %d times, want 1", count)
			}
		})
	}
}

func TestDict_ValuesSeq_EarlyBreak(t *testing.T) {
	for _, ctor := range allDictConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			count := 0
			for range ctor.make(seedPairs()...).ValuesSeq() {
				count++
				break
			}
			if count != 1 {
				t.Errorf("ValuesSeq() with early break yielded %d times, want 1", count)
			}
		})
	}
}

func TestTree_All_SortedOrder(t *testing.T) {
	tree := dicts.NewTree(
		dicts.Pair[string, int]{Key: "charlie", Value: 3},
		dicts.Pair[string, int]{Key: "alice", Value: 1},
		dicts.Pair[string, int]{Key: "bob", Value: 2},
	)
	var keys []string
	for k := range tree.All() {
		keys = append(keys, k)
	}
	if want := []string{"alice", "bob", "charlie"}; !reflect.DeepEqual(keys, want) {
		t.Errorf("Tree.All() keys = %v, want sorted %v", keys, want)
	}
}

func TestFromSeq2(t *testing.T) {
	for _, ctor := range allDictConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			source := ctor.make(seedPairs()...)
			got := dicts.FromSeq2(source.All())
			if !reflect.DeepEqual(got.AsMap(), source.AsMap()) {
				t.Errorf("FromSeq2 round-trip = %v, want %v", got.AsMap(), source.AsMap())
			}
		})
	}
}

func TestFromSeq2_LastWins(t *testing.T) {
	seq := func(yield func(string, int) bool) {
		_ = yield("k", 1) && yield("k", 2)
	}
	got := dicts.FromSeq2(seq)
	if v, _ := got.Get("k", 0); v != 2 {
		t.Errorf("FromSeq2 duplicate key: got %d, want last value 2", v)
	}
}

func ExampleFromSeq2() {
	source := dicts.NewTree(
		dicts.Pair[string, int]{Key: "a", Value: 1},
		dicts.Pair[string, int]{Key: "b", Value: 2},
	)
	d := dicts.FromSeq2(source.All())
	fmt.Println(d.Length())
	// Output: 2
}

// mutatingTraversal drives one callback or iterator method of d, calling visit
// with each key it is handed (every dictionary under test maps key to key, so a
// value doubles as its key).
type mutatingTraversal struct {
	name   string
	sorted bool // only SortedDict implementations have this method
	run    func(d dicts.MutableDict[int, int], visit func(key int))
}

func mutatingTraversals() []mutatingTraversal {
	return []mutatingTraversal{
		{name: "ForEach", run: func(d dicts.MutableDict[int, int], visit func(int)) { d.ForEach(func(k, _ int) { visit(k) }) }},
		{name: "ForEachKey", run: func(d dicts.MutableDict[int, int], visit func(int)) { d.ForEachKey(visit) }},
		{name: "ForEachValue", run: func(d dicts.MutableDict[int, int], visit func(int)) { d.ForEachValue(visit) }},
		{name: "Filter", run: func(d dicts.MutableDict[int, int], visit func(int)) {
			d.Filter(func(k, _ int) bool { visit(k); return true })
		}},
		{name: "FilterInPlace", run: func(d dicts.MutableDict[int, int], visit func(int)) {
			d.FilterInPlace(func(k, _ int) bool { visit(k); return true })
		}},
		{name: "AllMatch", run: func(d dicts.MutableDict[int, int], visit func(int)) {
			d.AllMatch(func(k, _ int) bool { visit(k); return true })
		}},
		{name: "AnyMatch", run: func(d dicts.MutableDict[int, int], visit func(int)) {
			d.AnyMatch(func(k, _ int) bool { visit(k); return false })
		}},
		{name: "NoneMatch", run: func(d dicts.MutableDict[int, int], visit func(int)) {
			d.NoneMatch(func(k, _ int) bool { visit(k); return false })
		}},
		{name: "Find", run: func(d dicts.MutableDict[int, int], visit func(int)) {
			d.Find(func(k, _ int) bool { visit(k); return false })
		}},
		{name: "FindKey", run: func(d dicts.MutableDict[int, int], visit func(int)) {
			d.FindKey(func(k int) bool { visit(k); return false })
		}},
		{name: "FindValue", run: func(d dicts.MutableDict[int, int], visit func(int)) {
			d.FindValue(func(v int) bool { visit(v); return false })
		}},
		{name: "All", run: func(d dicts.MutableDict[int, int], visit func(int)) {
			for k := range d.All() {
				visit(k)
			}
		}},
		{name: "KeysSeq", run: func(d dicts.MutableDict[int, int], visit func(int)) {
			for k := range d.KeysSeq() {
				visit(k)
			}
		}},
		{name: "ValuesSeq", run: func(d dicts.MutableDict[int, int], visit func(int)) {
			for v := range d.ValuesSeq() {
				visit(v)
			}
		}},
		{name: "Backward", sorted: true, run: func(d dicts.MutableDict[int, int], visit func(int)) {
			for k := range d.(dicts.SortedDict[int, int]).Backward() {
				visit(k)
			}
		}},
		{name: "RangeAll", sorted: true, run: func(d dicts.MutableDict[int, int], visit func(int)) {
			for k := range d.(dicts.SortedDict[int, int]).RangeAll(1, 7) {
				visit(k)
			}
		}},
	}
}

// TestDict_RemovingDuringIterationVisitsEveryEntry verifies that removing each
// visited entry from inside a callback or range loop still visits every entry
// exactly once and leaves the dictionary empty. Tree used to traverse its live
// nodes, so the AVL rebalancing and the successor copy of a two-child deletion
// made it skip entries: removing each key of 1..7 from ForEach left [6 7].
func TestDict_RemovingDuringIterationVisitsEveryEntry(t *testing.T) {
	seed := make([]dicts.Pair[int, int], 0, 7)
	for key := 1; key <= 7; key++ {
		seed = append(seed, dicts.Pair[int, int]{Key: key, Value: key})
	}
	implementations := []struct {
		name   string
		sorted bool
		make   func() dicts.MutableDict[int, int]
	}{
		{"Hash", false, func() dicts.MutableDict[int, int] { return dicts.NewHash(seed...) }},
		{"ConcurrentHash", false, func() dicts.MutableDict[int, int] { return dicts.NewConcurrentHash(seed...) }},
		{"ConcurrentHashRW", false, func() dicts.MutableDict[int, int] { return dicts.NewConcurrentHashRW(seed...) }},
		{"Tree", true, func() dicts.MutableDict[int, int] { return dicts.NewTree(seed...) }},
		{"ConcurrentTree", true, func() dicts.MutableDict[int, int] { return dicts.NewConcurrentTree(seed...) }},
		{"ConcurrentTreeRW", true, func() dicts.MutableDict[int, int] { return dicts.NewConcurrentTreeRW(seed...) }},
	}
	for _, traversal := range mutatingTraversals() {
		for _, impl := range implementations {
			if traversal.sorted && !impl.sorted {
				continue
			}
			t.Run(traversal.name+"/"+impl.name, func(t *testing.T) {
				d := impl.make()
				var visited []int
				traversal.run(d, func(key int) {
					visited = append(visited, key)
					d.RemoveInPlace(key)
				})

				sort.Ints(visited)
				if want := []int{1, 2, 3, 4, 5, 6, 7}; !reflect.DeepEqual(visited, want) {
					t.Errorf("%s() removing as it goes visited %v, want %v", traversal.name, visited, want)
				}
				if d.Length() != 0 {
					t.Errorf("%s() removing as it goes left %v", traversal.name, d.Keys())
				}
			})
		}
	}
}

// TestDict_IteratorSnapshotsWhenRangingBegins verifies that every iterator
// reads the dictionary when a range over it begins, not when the iterator is
// created: a write made between creating the iterator and ranging it is seen.
// Tree.KeysSeq/ValuesSeq and the concurrent dictionaries' iterators used to copy
// the entries when the method was called, so a stored iterator was stale.
func TestDict_IteratorSnapshotsWhenRangingBegins(t *testing.T) {
	iterators := []struct {
		name   string
		sorted bool
		create func(d dicts.MutableDict[int, int]) func(visit func(key int))
	}{
		{name: "All", create: func(d dicts.MutableDict[int, int]) func(func(int)) {
			seq := d.All()
			return func(visit func(int)) {
				for k := range seq {
					visit(k)
				}
			}
		}},
		{name: "KeysSeq", create: func(d dicts.MutableDict[int, int]) func(func(int)) {
			seq := d.KeysSeq()
			return func(visit func(int)) {
				for k := range seq {
					visit(k)
				}
			}
		}},
		{name: "ValuesSeq", create: func(d dicts.MutableDict[int, int]) func(func(int)) {
			seq := d.ValuesSeq()
			return func(visit func(int)) {
				for v := range seq {
					visit(v)
				}
			}
		}},
		{name: "Backward", sorted: true, create: func(d dicts.MutableDict[int, int]) func(func(int)) {
			seq := d.(dicts.SortedDict[int, int]).Backward()
			return func(visit func(int)) {
				for k := range seq {
					visit(k)
				}
			}
		}},
		{name: "RangeAll", sorted: true, create: func(d dicts.MutableDict[int, int]) func(func(int)) {
			seq := d.(dicts.SortedDict[int, int]).RangeAll(1, 7)
			return func(visit func(int)) {
				for k := range seq {
					visit(k)
				}
			}
		}},
	}
	seed := []dicts.Pair[int, int]{{Key: 1, Value: 1}, {Key: 2, Value: 2}}
	implementations := []struct {
		name   string
		sorted bool
		make   func() dicts.MutableDict[int, int]
	}{
		{"Hash", false, func() dicts.MutableDict[int, int] { return dicts.NewHash(seed...) }},
		{"ConcurrentHash", false, func() dicts.MutableDict[int, int] { return dicts.NewConcurrentHash(seed...) }},
		{"ConcurrentHashRW", false, func() dicts.MutableDict[int, int] { return dicts.NewConcurrentHashRW(seed...) }},
		{"Tree", true, func() dicts.MutableDict[int, int] { return dicts.NewTree(seed...) }},
		{"ConcurrentTree", true, func() dicts.MutableDict[int, int] { return dicts.NewConcurrentTree(seed...) }},
		{"ConcurrentTreeRW", true, func() dicts.MutableDict[int, int] { return dicts.NewConcurrentTreeRW(seed...) }},
	}
	for _, it := range iterators {
		for _, impl := range implementations {
			if it.sorted && !impl.sorted {
				continue
			}
			t.Run(it.name+"/"+impl.name, func(t *testing.T) {
				d := impl.make()
				rangeOver := it.create(d)
				d.PutInPlace(3, 3) // written after the iterator was created

				var visited []int
				rangeOver(func(key int) { visited = append(visited, key) })

				sort.Ints(visited)
				if want := []int{1, 2, 3}; !reflect.DeepEqual(visited, want) {
					t.Errorf("%s() created before PutInPlace then ranged = %v, want %v", it.name, visited, want)
				}
			})
		}
	}
}
