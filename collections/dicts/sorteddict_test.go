package dicts_test

import (
	"reflect"
	"sync"
	"testing"

	"github.com/pickeringtech/go-collections/collections/dicts"
)

// sortedDictFactory names a MutableSortedDict implementation and a constructor
// for it, so the shared behavioural suite can run against every backend.
type sortedDictFactory struct {
	name string
	make func(entries ...dicts.Pair[int, string]) dicts.MutableSortedDict[int, string]
}

func sortedDictFactories() []sortedDictFactory {
	return []sortedDictFactory{
		{"Tree", func(e ...dicts.Pair[int, string]) dicts.MutableSortedDict[int, string] {
			return dicts.NewTree(e...)
		}},
		{"ConcurrentTree", func(e ...dicts.Pair[int, string]) dicts.MutableSortedDict[int, string] {
			return dicts.NewConcurrentTree(e...)
		}},
		{"ConcurrentTreeRW", func(e ...dicts.Pair[int, string]) dicts.MutableSortedDict[int, string] {
			return dicts.NewConcurrentTreeRW(e...)
		}},
	}
}

// sampleEntries is an unsorted set of entries whose sorted key order is
// 1,3,4,6,7,9 — chosen so floor/ceiling/range have interior gaps to probe.
func orderedSampleEntries() []dicts.Pair[int, string] {
	return []dicts.Pair[int, string]{
		{Key: 6, Value: "six"},
		{Key: 3, Value: "three"},
		{Key: 9, Value: "nine"},
		{Key: 1, Value: "one"},
		{Key: 7, Value: "seven"},
		{Key: 4, Value: "four"},
	}
}

func collectSeq2(seq func(yield func(int, string) bool)) ([]int, []string) {
	keys := []int{}
	values := []string{}
	seq(func(k int, v string) bool {
		keys = append(keys, k)
		values = append(values, v)
		return true
	})
	return keys, values
}

func TestSortedDict_MinMax(t *testing.T) {
	for _, f := range sortedDictFactories() {
		t.Run(f.name, func(t *testing.T) {
			empty := f.make()
			if k, v, ok := empty.Min(); ok || k != 0 || v != "" {
				t.Errorf("empty Min() = (%d, %q, %v), want (0, \"\", false)", k, v, ok)
			}
			if k, v, ok := empty.Max(); ok || k != 0 || v != "" {
				t.Errorf("empty Max() = (%d, %q, %v), want (0, \"\", false)", k, v, ok)
			}

			d := f.make(orderedSampleEntries()...)
			if k, v, ok := d.Min(); !ok || k != 1 || v != "one" {
				t.Errorf("Min() = (%d, %q, %v), want (1, \"one\", true)", k, v, ok)
			}
			if k, v, ok := d.Max(); !ok || k != 9 || v != "nine" {
				t.Errorf("Max() = (%d, %q, %v), want (9, \"nine\", true)", k, v, ok)
			}
		})
	}
}

func TestSortedDict_Floor(t *testing.T) {
	type want struct {
		key   int
		value string
		ok    bool
	}
	tests := []struct {
		name string
		arg  int
		want want
	}{
		{"exact match", 4, want{4, "four", true}},
		{"between keys rounds down", 5, want{4, "four", true}},
		{"above max returns max", 100, want{9, "nine", true}},
		{"equals min", 1, want{1, "one", true}},
		{"below min has no floor", 0, want{0, "", false}},
	}
	for _, f := range sortedDictFactories() {
		for _, tt := range tests {
			t.Run(f.name+"/"+tt.name, func(t *testing.T) {
				d := f.make(orderedSampleEntries()...)
				k, v, ok := d.Floor(tt.arg)
				if k != tt.want.key || v != tt.want.value || ok != tt.want.ok {
					t.Errorf("Floor(%d) = (%d, %q, %v), want (%d, %q, %v)",
						tt.arg, k, v, ok, tt.want.key, tt.want.value, tt.want.ok)
				}
			})
		}
	}
}

func TestSortedDict_Ceiling(t *testing.T) {
	type want struct {
		key   int
		value string
		ok    bool
	}
	tests := []struct {
		name string
		arg  int
		want want
	}{
		{"exact match", 4, want{4, "four", true}},
		{"between keys rounds up", 5, want{6, "six", true}},
		{"below min returns min", -5, want{1, "one", true}},
		{"equals max", 9, want{9, "nine", true}},
		{"above max has no ceiling", 10, want{0, "", false}},
	}
	for _, f := range sortedDictFactories() {
		for _, tt := range tests {
			t.Run(f.name+"/"+tt.name, func(t *testing.T) {
				d := f.make(orderedSampleEntries()...)
				k, v, ok := d.Ceiling(tt.arg)
				if k != tt.want.key || v != tt.want.value || ok != tt.want.ok {
					t.Errorf("Ceiling(%d) = (%d, %q, %v), want (%d, %q, %v)",
						tt.arg, k, v, ok, tt.want.key, tt.want.value, tt.want.ok)
				}
			})
		}
	}
}

func TestSortedDict_Range(t *testing.T) {
	tests := []struct {
		name   string
		lo, hi int
		want   []dicts.Pair[int, string]
	}{
		{"interior inclusive", 3, 7, []dicts.Pair[int, string]{{3, "three"}, {4, "four"}, {6, "six"}, {7, "seven"}}},
		{"bounds outside data", -1, 100, []dicts.Pair[int, string]{{1, "one"}, {3, "three"}, {4, "four"}, {6, "six"}, {7, "seven"}, {9, "nine"}}},
		{"single key", 4, 4, []dicts.Pair[int, string]{{4, "four"}}},
		{"missing endpoints", 2, 5, []dicts.Pair[int, string]{{3, "three"}, {4, "four"}}},
		{"empty when lo greater than hi", 7, 3, []dicts.Pair[int, string]{}},
		{"gap yields empty", 100, 200, []dicts.Pair[int, string]{}},
	}
	for _, f := range sortedDictFactories() {
		for _, tt := range tests {
			t.Run(f.name+"/"+tt.name, func(t *testing.T) {
				d := f.make(orderedSampleEntries()...)
				got := d.Range(tt.lo, tt.hi)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Range(%d, %d) = %v, want %v", tt.lo, tt.hi, got, tt.want)
				}
			})
		}
	}
}

func TestSortedDict_Iterators(t *testing.T) {
	for _, f := range sortedDictFactories() {
		t.Run(f.name, func(t *testing.T) {
			d := f.make(orderedSampleEntries()...)

			ascKeys, ascValues := collectSeq2(d.All())
			wantAscKeys := []int{1, 3, 4, 6, 7, 9}
			wantAscValues := []string{"one", "three", "four", "six", "seven", "nine"}
			if !reflect.DeepEqual(ascKeys, wantAscKeys) || !reflect.DeepEqual(ascValues, wantAscValues) {
				t.Errorf("All() = %v/%v, want %v/%v", ascKeys, ascValues, wantAscKeys, wantAscValues)
			}

			descKeys, _ := collectSeq2(d.Backward())
			wantDescKeys := []int{9, 7, 6, 4, 3, 1}
			if !reflect.DeepEqual(descKeys, wantDescKeys) {
				t.Errorf("Backward() = %v, want %v", descKeys, wantDescKeys)
			}

			rangeKeys, _ := collectSeq2(d.RangeAll(3, 7))
			wantRangeKeys := []int{3, 4, 6, 7}
			if !reflect.DeepEqual(rangeKeys, wantRangeKeys) {
				t.Errorf("RangeAll(3, 7) = %v, want %v", rangeKeys, wantRangeKeys)
			}
		})
	}
}

// TestSortedDict_IteratorEarlyStop verifies that breaking out of the range-over-func
// stops iteration promptly (exercises the early-stop branches in every backend).
func TestSortedDict_IteratorEarlyStop(t *testing.T) {
	for _, f := range sortedDictFactories() {
		t.Run(f.name, func(t *testing.T) {
			d := f.make(orderedSampleEntries()...)

			var firstAsc int
			countAsc := 0
			d.All()(func(k int, _ string) bool {
				firstAsc = k
				countAsc++
				return false
			})
			if countAsc != 1 || firstAsc != 1 {
				t.Errorf("All() early stop visited %d items, first %d; want 1 item, first 1", countAsc, firstAsc)
			}

			var firstDesc int
			countDesc := 0
			d.Backward()(func(k int, _ string) bool {
				firstDesc = k
				countDesc++
				return false
			})
			if countDesc != 1 || firstDesc != 9 {
				t.Errorf("Backward() early stop visited %d items, first %d; want 1 item, first 9", countDesc, firstDesc)
			}

			countRange := 0
			d.RangeAll(3, 7)(func(_ int, _ string) bool {
				countRange++
				return false
			})
			if countRange != 1 {
				t.Errorf("RangeAll() early stop visited %d items, want 1", countRange)
			}
		})
	}
}

// TestSortedDict_DictSurface exercises the full Dict/MutableDict surface so the
// concurrent backends' delegating methods are covered with real assertions.
func TestSortedDict_DictSurface(t *testing.T) {
	for _, f := range sortedDictFactories() {
		t.Run(f.name, func(t *testing.T) {
			d := f.make(orderedSampleEntries()...)

			if d.IsEmpty() {
				t.Error("IsEmpty() = true, want false")
			}
			if d.Length() != 6 {
				t.Errorf("Length() = %d, want 6", d.Length())
			}
			if v, ok := d.Get(4, "?"); !ok || v != "four" {
				t.Errorf("Get(4) = (%q, %v), want (\"four\", true)", v, ok)
			}
			if v, ok := d.Get(2, "?"); ok || v != "?" {
				t.Errorf("Get(2) = (%q, %v), want (\"?\", false)", v, ok)
			}
			if !d.Contains(7) || d.Contains(2) {
				t.Error("Contains mismatch for keys 7/2")
			}

			// Iteration callbacks.
			var keySum int
			d.ForEachKey(func(k int) { keySum += k })
			if keySum != 30 {
				t.Errorf("ForEachKey sum = %d, want 30", keySum)
			}
			valueCount := 0
			d.ForEachValue(func(string) { valueCount++ })
			if valueCount != 6 {
				t.Errorf("ForEachValue count = %d, want 6", valueCount)
			}
			pairCount := 0
			d.ForEach(func(int, string) { pairCount++ })
			if pairCount != 6 {
				t.Errorf("ForEach count = %d, want 6", pairCount)
			}

			// Conversions (all sorted by key).
			if got := d.Keys(); !reflect.DeepEqual(got, []int{1, 3, 4, 6, 7, 9}) {
				t.Errorf("Keys() = %v", got)
			}
			if got := d.Values(); !reflect.DeepEqual(got, []string{"one", "three", "four", "six", "seven", "nine"}) {
				t.Errorf("Values() = %v", got)
			}
			if got := d.Items(); !reflect.DeepEqual(got, sortedItems()) {
				t.Errorf("Items() = %v", got)
			}
			if got := d.AsMap(); !reflect.DeepEqual(got, sampleMap()) {
				t.Errorf("AsMap() = %v", got)
			}

			// Search predicates.
			if !d.AllMatch(func(k int, _ string) bool { return k > 0 }) {
				t.Error("AllMatch(k>0) = false, want true")
			}
			if d.AllMatch(func(k int, _ string) bool { return k > 5 }) {
				t.Error("AllMatch(k>5) = true, want false")
			}
			if !d.AnyMatch(func(k int, _ string) bool { return k == 6 }) {
				t.Error("AnyMatch(k==6) = false, want true")
			}
			if !d.NoneMatch(func(k int, _ string) bool { return k == 2 }) {
				t.Error("NoneMatch(k==2) = false, want true")
			}
			if k, v, ok := d.Find(func(k int, _ string) bool { return k >= 4 }); !ok || k != 4 || v != "four" {
				t.Errorf("Find(k>=4) = (%d, %q, %v), want (4, \"four\", true)", k, v, ok)
			}
			if k, ok := d.FindKey(func(k int) bool { return k > 6 }); !ok || k != 7 {
				t.Errorf("FindKey(k>6) = (%d, %v), want (7, true)", k, ok)
			}
			if v, ok := d.FindValue(func(v string) bool { return v == "nine" }); !ok || v != "nine" {
				t.Errorf("FindValue(nine) = (%q, %v), want (\"nine\", true)", v, ok)
			}
			if !d.ContainsValue("seven") || d.ContainsValue("zero") {
				t.Error("ContainsValue mismatch for seven/zero")
			}
		})
	}
}

// TestSortedDict_Immutability verifies the returns-new methods leave the receiver
// untouched, while the InPlace methods mutate it.
func TestSortedDict_Immutability(t *testing.T) {
	for _, f := range sortedDictFactories() {
		t.Run(f.name, func(t *testing.T) {
			d := f.make(orderedSampleEntries()...)

			added := d.Put(5, "five")
			if d.Contains(5) {
				t.Error("Put mutated the receiver")
			}
			if !added.Contains(5) {
				t.Error("Put result missing the new key")
			}

			addedMany := d.PutMany(dicts.Pair[int, string]{Key: 2, Value: "two"}, dicts.Pair[int, string]{Key: 8, Value: "eight"})
			if d.Contains(2) || d.Contains(8) {
				t.Error("PutMany mutated the receiver")
			}
			if !addedMany.Contains(2) || !addedMany.Contains(8) {
				t.Error("PutMany result missing new keys")
			}

			removed := d.Remove(4)
			if !d.Contains(4) {
				t.Error("Remove mutated the receiver")
			}
			if removed.Contains(4) {
				t.Error("Remove result still has the key")
			}

			removedMany := d.RemoveMany(1, 9)
			if !d.Contains(1) || !d.Contains(9) {
				t.Error("RemoveMany mutated the receiver")
			}
			if removedMany.Contains(1) || removedMany.Contains(9) {
				t.Error("RemoveMany result still has keys")
			}

			filtered := d.Filter(func(k int, _ string) bool { return k < 5 })
			if filtered.Length() != 3 {
				t.Errorf("Filter length = %d, want 3", filtered.Length())
			}
			if d.Length() != 6 {
				t.Error("Filter mutated the receiver")
			}

			// In-place mutations.
			d.PutInPlace(5, "five")
			if !d.Contains(5) {
				t.Error("PutInPlace did not add key")
			}
			d.PutManyInPlace(dicts.Pair[int, string]{Key: 2, Value: "two"})
			if !d.Contains(2) {
				t.Error("PutManyInPlace did not add key")
			}
			if v, ok := d.RemoveInPlace(5); !ok || v != "five" {
				t.Errorf("RemoveInPlace(5) = (%q, %v), want (\"five\", true)", v, ok)
			}
			d.RemoveManyInPlace(2, 3)
			if d.Contains(2) || d.Contains(3) {
				t.Error("RemoveManyInPlace did not remove keys")
			}
			d.FilterInPlace(func(k int, _ string) bool { return k%2 == 0 })
			if d.AnyMatch(func(k int, _ string) bool { return k%2 != 0 }) {
				t.Error("FilterInPlace left odd keys")
			}
			d.Clear()
			if !d.IsEmpty() {
				t.Error("Clear did not empty the dict")
			}
		})
	}
}

func sortedItems() []dicts.Pair[int, string] {
	return []dicts.Pair[int, string]{
		{Key: 1, Value: "one"}, {Key: 3, Value: "three"}, {Key: 4, Value: "four"},
		{Key: 6, Value: "six"}, {Key: 7, Value: "seven"}, {Key: 9, Value: "nine"},
	}
}

func sampleMap() map[int]string {
	return map[int]string{1: "one", 3: "three", 4: "four", 6: "six", 7: "seven", 9: "nine"}
}

// TestConcurrentTree_ReturnsConcurrentType asserts the thread-safe-in →
// thread-safe-out contract: immutable ops return the same concurrent type.
func TestConcurrentTree_ReturnsConcurrentType(t *testing.T) {
	ct := dicts.NewConcurrentTree(orderedSampleEntries()...)
	if _, ok := ct.Put(2, "two").(*dicts.ConcurrentTree[int, string]); !ok {
		t.Error("ConcurrentTree.Put did not return *ConcurrentTree")
	}
	if _, ok := ct.PutMany(dicts.Pair[int, string]{Key: 2, Value: "two"}).(*dicts.ConcurrentTree[int, string]); !ok {
		t.Error("ConcurrentTree.PutMany did not return *ConcurrentTree")
	}
	if _, ok := ct.Remove(1).(*dicts.ConcurrentTree[int, string]); !ok {
		t.Error("ConcurrentTree.Remove did not return *ConcurrentTree")
	}
	if _, ok := ct.RemoveMany(1).(*dicts.ConcurrentTree[int, string]); !ok {
		t.Error("ConcurrentTree.RemoveMany did not return *ConcurrentTree")
	}
	if _, ok := ct.Filter(func(int, string) bool { return true }).(*dicts.ConcurrentTree[int, string]); !ok {
		t.Error("ConcurrentTree.Filter did not return *ConcurrentTree")
	}

	rw := dicts.NewConcurrentTreeRW(orderedSampleEntries()...)
	if _, ok := rw.Put(2, "two").(*dicts.ConcurrentTreeRW[int, string]); !ok {
		t.Error("ConcurrentTreeRW.Put did not return *ConcurrentTreeRW")
	}
	if _, ok := rw.PutMany(dicts.Pair[int, string]{Key: 2, Value: "two"}).(*dicts.ConcurrentTreeRW[int, string]); !ok {
		t.Error("ConcurrentTreeRW.PutMany did not return *ConcurrentTreeRW")
	}
	if _, ok := rw.Remove(1).(*dicts.ConcurrentTreeRW[int, string]); !ok {
		t.Error("ConcurrentTreeRW.Remove did not return *ConcurrentTreeRW")
	}
	if _, ok := rw.RemoveMany(1).(*dicts.ConcurrentTreeRW[int, string]); !ok {
		t.Error("ConcurrentTreeRW.RemoveMany did not return *ConcurrentTreeRW")
	}
	if _, ok := rw.Filter(func(int, string) bool { return true }).(*dicts.ConcurrentTreeRW[int, string]); !ok {
		t.Error("ConcurrentTreeRW.Filter did not return *ConcurrentTreeRW")
	}
}

// TestConcurrentTree_RaceSafety drives both concurrent trees from many
// goroutines; run with -race it asserts there are no data races or panics.
func TestConcurrentTree_RaceSafety(t *testing.T) {
	mutexTree := dicts.NewConcurrentTree[int, string]()
	rwTree := dicts.NewConcurrentTreeRW[int, string]()
	dictsUnderTest := []dicts.MutableSortedDict[int, string]{mutexTree, rwTree}

	for _, d := range dictsUnderTest {
		var wg sync.WaitGroup
		for w := 0; w < 8; w++ {
			wg.Add(1)
			go func(base int) {
				defer wg.Done()
				for i := 0; i < 50; i++ {
					key := base*50 + i
					d.PutInPlace(key, "v")
					d.Contains(key)
					d.Min()
					d.Floor(key)
					for range d.All() {
						break
					}
					d.RemoveInPlace(key)
				}
			}(w)
		}
		wg.Wait()
		if !d.IsEmpty() {
			t.Errorf("expected empty dict after balanced put/remove, got length %d", d.Length())
		}
	}
}
