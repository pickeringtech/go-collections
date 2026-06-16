package multimaps_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/collections/multimaps"
)

// factory builds a fresh MutableMultimap seeded with the given entries. dedupe
// reports whether the implementation collapses duplicate (key, value) pairs
// (set-backed) or keeps them (list-backed).
type factory struct {
	name   string
	dedupe bool
	make   func(entries ...multimaps.Entry[string, int]) multimaps.MutableMultimap[string, int]
}

func factories() []factory {
	return []factory{
		{
			name:   "ListMultimap",
			dedupe: false,
			make: func(entries ...multimaps.Entry[string, int]) multimaps.MutableMultimap[string, int] {
				return multimaps.NewListMultimap(entries...)
			},
		},
		{
			name:   "ConcurrentListMultimap",
			dedupe: false,
			make: func(entries ...multimaps.Entry[string, int]) multimaps.MutableMultimap[string, int] {
				return multimaps.NewConcurrentListMultimap(entries...)
			},
		},
		{
			name:   "ConcurrentRWListMultimap",
			dedupe: false,
			make: func(entries ...multimaps.Entry[string, int]) multimaps.MutableMultimap[string, int] {
				return multimaps.NewConcurrentRWListMultimap(entries...)
			},
		},
		{
			name:   "SetMultimap",
			dedupe: true,
			make: func(entries ...multimaps.Entry[string, int]) multimaps.MutableMultimap[string, int] {
				return multimaps.NewSetMultimap(entries...)
			},
		},
		{
			name:   "ConcurrentSetMultimap",
			dedupe: true,
			make: func(entries ...multimaps.Entry[string, int]) multimaps.MutableMultimap[string, int] {
				return multimaps.NewConcurrentSetMultimap(entries...)
			},
		},
		{
			name:   "ConcurrentRWSetMultimap",
			dedupe: true,
			make: func(entries ...multimaps.Entry[string, int]) multimaps.MutableMultimap[string, int] {
				return multimaps.NewConcurrentRWSetMultimap(entries...)
			},
		},
	}
}

func sortedInts(in []int) []int {
	out := append([]int(nil), in...)
	sort.Ints(out)
	return out
}

func sortedStrings(in []string) []string {
	out := append([]string(nil), in...)
	sort.Strings(out)
	return out
}

func TestMultimap_EmptyState(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			m := f.make()
			if !m.IsEmpty() {
				t.Errorf("IsEmpty() = false, want true")
			}
			if m.Length() != 0 {
				t.Errorf("Length() = %d, want 0", m.Length())
			}
			if m.KeyCount() != 0 {
				t.Errorf("KeyCount() = %d, want 0", m.KeyCount())
			}
			if m.ContainsKey("missing") {
				t.Errorf("ContainsKey(missing) = true, want false")
			}
			if m.ContainsEntry("missing", 1) {
				t.Errorf("ContainsEntry(missing, 1) = true, want false")
			}
			got := m.Get("missing")
			if got == nil || len(got) != 0 {
				t.Errorf("Get(missing) = %v, want empty non-nil slice", got)
			}
		})
	}
}

func TestMultimap_PutAndGet(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			m := f.make()
			m.PutInPlace("a", 1)
			m.PutInPlace("a", 2)
			m.PutInPlace("b", 3)

			if m.IsEmpty() {
				t.Errorf("IsEmpty() = true, want false")
			}
			if m.KeyCount() != 2 {
				t.Errorf("KeyCount() = %d, want 2", m.KeyCount())
			}
			if m.Length() != 3 {
				t.Errorf("Length() = %d, want 3", m.Length())
			}
			if !m.ContainsKey("a") {
				t.Errorf("ContainsKey(a) = false, want true")
			}
			if !m.ContainsEntry("a", 1) {
				t.Errorf("ContainsEntry(a, 1) = false, want true")
			}
			if m.ContainsEntry("a", 99) {
				t.Errorf("ContainsEntry(a, 99) = true, want false")
			}
			got := sortedInts(m.Get("a"))
			want := []int{1, 2}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Get(a) = %v, want %v", got, want)
			}
		})
	}
}

func TestMultimap_GetReturnsIndependentCopy(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			m := f.make()
			m.PutInPlace("a", 1)
			got := m.Get("a")
			got[0] = 999 // mutating the copy must not affect the multimap
			if !m.ContainsEntry("a", 1) {
				t.Errorf("mutating Get() result changed the multimap")
			}
		})
	}
}

func TestMultimap_DuplicateSemantics(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			m := f.make()
			m.PutInPlace("a", 1)
			m.PutInPlace("a", 1)

			wantLen := 2
			if f.dedupe {
				wantLen = 1
			}
			if m.Length() != wantLen {
				t.Errorf("Length() after duplicate Put = %d, want %d", m.Length(), wantLen)
			}
		})
	}
}

func TestMultimap_PutAllInPlace(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			m := f.make()
			m.PutAllInPlace("a", 1, 2, 3)
			m.PutAllInPlace("a") // no values: no-op
			m.PutAllInPlace("b", 4)

			if m.Length() != 4 {
				t.Errorf("Length() = %d, want 4", m.Length())
			}
			got := sortedInts(m.Get("a"))
			want := []int{1, 2, 3}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Get(a) = %v, want %v", got, want)
			}
			if m.ContainsKey("c") {
				t.Errorf("empty PutAllInPlace must not create a key")
			}
		})
	}
}

func TestMultimap_RemoveInPlace(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			m := f.make()
			m.PutAllInPlace("a", 1, 2)

			if m.RemoveInPlace("missing", 1) {
				t.Errorf("RemoveInPlace(missing) = true, want false")
			}
			if m.RemoveInPlace("a", 99) {
				t.Errorf("RemoveInPlace(a, 99) = true, want false")
			}
			if !m.RemoveInPlace("a", 1) {
				t.Errorf("RemoveInPlace(a, 1) = false, want true")
			}
			got := sortedInts(m.Get("a"))
			want := []int{2}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Get(a) = %v, want %v", got, want)
			}
			// Removing the final value drops the key entirely.
			if !m.RemoveInPlace("a", 2) {
				t.Errorf("RemoveInPlace(a, 2) = false, want true")
			}
			if m.ContainsKey("a") {
				t.Errorf("ContainsKey(a) = true after removing last value, want false")
			}
		})
	}
}

func TestMultimap_RemoveAllInPlace(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			m := f.make()
			m.PutAllInPlace("a", 1, 2)

			values, ok := m.RemoveAllInPlace("missing")
			if ok {
				t.Errorf("RemoveAllInPlace(missing) ok = true, want false")
			}
			if values == nil || len(values) != 0 {
				t.Errorf("RemoveAllInPlace(missing) values = %v, want empty non-nil", values)
			}

			values, ok = m.RemoveAllInPlace("a")
			if !ok {
				t.Errorf("RemoveAllInPlace(a) ok = false, want true")
			}
			got := sortedInts(values)
			want := []int{1, 2}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("RemoveAllInPlace(a) values = %v, want %v", got, want)
			}
			if !m.IsEmpty() {
				t.Errorf("multimap not empty after removing only key")
			}
		})
	}
}

func TestMultimap_Clear(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			m := f.make()
			m.PutAllInPlace("a", 1, 2)
			m.PutInPlace("b", 3)
			m.Clear()
			if !m.IsEmpty() {
				t.Errorf("IsEmpty() = false after Clear, want true")
			}
		})
	}
}

func TestMultimap_ImmutableOpsDoNotMutateOriginal(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			m := f.make()
			m.PutInPlace("a", 1)

			put := m.Put("a", 2)
			putAll := m.PutAll("b", 3, 4)
			removed := m.Remove("a", 1)
			removedAll := put.RemoveAll("a")

			// Original is untouched by every immutable operation.
			if m.Length() != 1 {
				t.Errorf("original Length() = %d, want 1 (immutable op mutated receiver)", m.Length())
			}
			if !put.ContainsEntry("a", 2) {
				t.Errorf("Put result missing new entry")
			}
			if !putAll.ContainsEntry("b", 3) || !putAll.ContainsEntry("b", 4) {
				t.Errorf("PutAll result missing new entries")
			}
			if removed.ContainsEntry("a", 1) {
				t.Errorf("Remove result still contains removed entry")
			}
			if removedAll.ContainsKey("a") {
				t.Errorf("RemoveAll result still contains removed key")
			}
		})
	}
}

func TestMultimap_PutAllImmutableNoValues(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			m := f.make()
			m.PutInPlace("a", 1)
			result := m.PutAll("a") // no values
			if result.Length() != 1 {
				t.Errorf("PutAll() with no values Length() = %d, want 1", result.Length())
			}
		})
	}
}

func TestMultimap_Filter(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			m := f.make()
			m.PutAllInPlace("a", 1, 2, 3)
			m.PutAllInPlace("b", 10, 20)

			even := m.Filter(func(_ string, value int) bool { return value%2 == 0 })
			gotA := sortedInts(even.Get("a"))
			wantA := []int{2}
			if !reflect.DeepEqual(gotA, wantA) {
				t.Errorf("Filter Get(a) = %v, want %v", gotA, wantA)
			}
			gotB := sortedInts(even.Get("b"))
			wantB := []int{10, 20}
			if !reflect.DeepEqual(gotB, wantB) {
				t.Errorf("Filter Get(b) = %v, want %v", gotB, wantB)
			}
			// Original unchanged.
			if m.Length() != 5 {
				t.Errorf("original Length() = %d after Filter, want 5", m.Length())
			}

			// A predicate that excludes every value for a key drops the key.
			none := m.Filter(func(_ string, _ int) bool { return false })
			if !none.IsEmpty() {
				t.Errorf("Filter(false) not empty")
			}
		})
	}
}

func TestMultimap_FilterInPlace(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			m := f.make()
			m.PutAllInPlace("a", 1, 2, 3)
			m.PutInPlace("b", 1)

			m.FilterInPlace(func(_ string, value int) bool { return value >= 3 })
			if m.ContainsKey("b") {
				t.Errorf("FilterInPlace did not drop key with no surviving values")
			}
			got := sortedInts(m.Get("a"))
			want := []int{3}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("FilterInPlace Get(a) = %v, want %v", got, want)
			}
		})
	}
}

func TestMultimap_SearchMatchers(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			empty := f.make()
			if !empty.AllMatch(func(string, int) bool { return false }) {
				t.Errorf("AllMatch on empty = false, want true (vacuous)")
			}
			if empty.AnyMatch(func(string, int) bool { return true }) {
				t.Errorf("AnyMatch on empty = true, want false")
			}
			if !empty.NoneMatch(func(string, int) bool { return true }) {
				t.Errorf("NoneMatch on empty = false, want true (vacuous)")
			}

			m := f.make()
			m.PutAllInPlace("a", 2, 4, 6)
			if !m.AllMatch(func(_ string, v int) bool { return v%2 == 0 }) {
				t.Errorf("AllMatch even = false, want true")
			}
			if m.AllMatch(func(_ string, v int) bool { return v > 4 }) {
				t.Errorf("AllMatch >4 = true, want false")
			}
			if !m.AnyMatch(func(_ string, v int) bool { return v == 4 }) {
				t.Errorf("AnyMatch ==4 = false, want true")
			}
			if m.AnyMatch(func(_ string, v int) bool { return v == 5 }) {
				t.Errorf("AnyMatch ==5 = true, want false")
			}
			if !m.NoneMatch(func(_ string, v int) bool { return v == 5 }) {
				t.Errorf("NoneMatch ==5 = false, want true")
			}
		})
	}
}

func TestMultimap_Find(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			m := f.make()
			m.PutAllInPlace("a", 1, 2, 3)

			key, value, ok := m.Find(func(_ string, v int) bool { return v == 2 })
			if !ok {
				t.Errorf("Find(==2) ok = false, want true")
			}
			if key != "a" || value != 2 {
				t.Errorf("Find(==2) = (%q, %d), want (a, 2)", key, value)
			}

			_, _, ok = m.Find(func(_ string, v int) bool { return v == 99 })
			if ok {
				t.Errorf("Find(==99) ok = true, want false")
			}
		})
	}
}

func TestMultimap_Conversions(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			m := f.make()
			m.PutAllInPlace("a", 1, 2)
			m.PutInPlace("b", 3)

			gotKeys := sortedStrings(m.Keys())
			wantKeys := []string{"a", "b"}
			if !reflect.DeepEqual(gotKeys, wantKeys) {
				t.Errorf("Keys() = %v, want %v", gotKeys, wantKeys)
			}
			gotValues := sortedInts(m.Values())
			wantValues := []int{1, 2, 3}
			if !reflect.DeepEqual(gotValues, wantValues) {
				t.Errorf("Values() = %v, want %v", gotValues, wantValues)
			}
			if len(m.Entries()) != 3 {
				t.Errorf("len(Entries()) = %d, want 3", len(m.Entries()))
			}

			asMap := m.AsMap()
			gotMapA := sortedInts(asMap["a"])
			wantMapA := []int{1, 2}
			if !reflect.DeepEqual(gotMapA, wantMapA) {
				t.Errorf("AsMap()[a] = %v, want %v", asMap["a"], wantMapA)
			}
			// AsMap is independent of the multimap.
			asMap["a"][0] = 999
			if !m.ContainsEntry("a", 1) {
				t.Errorf("mutating AsMap() result changed the multimap")
			}
		})
	}
}

func TestMultimap_Entries(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			m := f.make()
			m.PutAllInPlace("a", 1, 2)

			entries := m.Entries()
			got := make([]int, 0, len(entries))
			for _, e := range entries {
				if e.Key != "a" {
					t.Errorf("Entry key = %q, want a", e.Key)
				}
				got = append(got, e.Value)
			}
			want := []int{1, 2}
			if !reflect.DeepEqual(sortedInts(got), want) {
				t.Errorf("entry values = %v, want %v", sortedInts(got), want)
			}
		})
	}
}

func TestMultimap_ForEach(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			m := f.make()
			m.PutAllInPlace("a", 1, 2)
			m.PutInPlace("b", 3)

			var seen []int
			m.ForEach(func(_ string, value int) {
				seen = append(seen, value)
			})
			got := sortedInts(seen)
			want := []int{1, 2, 3}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("ForEach values = %v, want %v", got, want)
			}
		})
	}
}

func TestMultimap_ForEachKey(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			m := f.make()
			m.PutAllInPlace("a", 1, 2)
			m.PutInPlace("b", 3)

			counts := map[string]int{}
			m.ForEachKey(func(key string, values []int) {
				counts[key] = len(values)
				// Mutating the supplied copy must not affect the multimap.
				if len(values) > 0 {
					values[0] = -1
				}
			})
			if counts["a"] != 2 || counts["b"] != 1 {
				t.Errorf("ForEachKey counts = %v, want a:2 b:1", counts)
			}
			if !m.ContainsEntry("a", 1) || !m.ContainsEntry("a", 2) {
				t.Errorf("ForEachKey leaked a mutable view of the values")
			}
		})
	}
}

func TestMultimap_AllIterator(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			m := f.make()
			m.PutAllInPlace("a", 1, 2)
			m.PutInPlace("b", 3)

			var seen []int
			for _, value := range m.All() {
				seen = append(seen, value)
			}
			got := sortedInts(seen)
			want := []int{1, 2, 3}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("All() values = %v, want %v", got, want)
			}

			// Early termination must stop iteration.
			count := 0
			for range m.All() {
				count++
				break
			}
			if count != 1 {
				t.Errorf("All() early-return visited %d entries, want 1", count)
			}
		})
	}
}

func TestMultimap_KeysSeqIterator(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			m := f.make()
			m.PutInPlace("a", 1)
			m.PutInPlace("b", 2)

			var seen []string
			for key := range m.KeysSeq() {
				seen = append(seen, key)
			}
			got := sortedStrings(seen)
			want := []string{"a", "b"}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("KeysSeq() = %v, want %v", got, want)
			}

			count := 0
			for range m.KeysSeq() {
				count++
				break
			}
			if count != 1 {
				t.Errorf("KeysSeq() early-return visited %d keys, want 1", count)
			}
		})
	}
}

func TestMultimap_ConstructorWithEntries(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			m := f.make(
				multimaps.Entry[string, int]{Key: "a", Value: 1},
				multimaps.Entry[string, int]{Key: "a", Value: 1},
				multimaps.Entry[string, int]{Key: "a", Value: 2},
			)
			wantLen := 3
			if f.dedupe {
				wantLen = 2
			}
			if m.Length() != wantLen {
				t.Errorf("constructor Length() = %d, want %d", m.Length(), wantLen)
			}
		})
	}
}

// TestListMultimap_OrderPreserved verifies the list-backed multimap keeps
// insertion order and duplicates, which the set-backed variant does not.
func TestListMultimap_OrderPreserved(t *testing.T) {
	m := multimaps.NewListMultimap[string, int]()
	m.PutInPlace("a", 3)
	m.PutInPlace("a", 1)
	m.PutInPlace("a", 3)
	m.PutInPlace("a", 2)

	got := m.Get("a")
	want := []int{3, 1, 3, 2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Get(a) = %v, want %v (order + duplicates preserved)", got, want)
	}
}

// TestListMultimap_RemoveFirstOccurrence verifies that Remove drops only a
// single (first) occurrence of a duplicated value.
func TestListMultimap_RemoveFirstOccurrence(t *testing.T) {
	m := multimaps.NewListMultimap[string, int]()
	m.PutAllInPlace("a", 1, 2, 1, 3)

	m.RemoveInPlace("a", 1)
	got := m.Get("a")
	want := []int{2, 1, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Get(a) = %v, want %v (only first occurrence removed)", got, want)
	}
}

// TestListMultimap_NonComparableValues verifies that the list-backed multimap
// accepts non-comparable value types (here, []int) and compares them with
// reflect.DeepEqual for ContainsEntry / Remove.
func TestListMultimap_NonComparableValues(t *testing.T) {
	m := multimaps.NewListMultimap[string, []int]()
	m.PutInPlace("a", []int{1, 2})
	m.PutInPlace("a", []int{3, 4})

	if !m.ContainsEntry("a", []int{1, 2}) {
		t.Errorf("ContainsEntry(a, [1 2]) = false, want true")
	}
	if m.ContainsEntry("a", []int{9, 9}) {
		t.Errorf("ContainsEntry(a, [9 9]) = true, want false")
	}
	if !m.RemoveInPlace("a", []int{1, 2}) {
		t.Errorf("RemoveInPlace(a, [1 2]) = false, want true")
	}
	got := m.Get("a")
	want := [][]int{{3, 4}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Get(a) = %v, want %v", got, want)
	}
}

// TestSetMultimap_Dedupes verifies the set-backed multimap collapses duplicates.
func TestSetMultimap_Dedupes(t *testing.T) {
	m := multimaps.NewSetMultimap[string, int]()
	m.PutInPlace("a", 1)
	m.PutInPlace("a", 1)
	m.PutAllInPlace("a", 1, 2, 2)

	got := sortedInts(m.Get("a"))
	want := []int{1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Get(a) = %v, want %v (deduped)", got, want)
	}
}
