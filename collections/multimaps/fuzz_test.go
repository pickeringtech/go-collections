package multimaps_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/collections/multimaps"
)

// op codes derived from fuzz bytes.
const (
	opPut = iota
	opRemove
	opRemoveAll
	opFilter
	opCount
)

func fuzzKey(b byte) string {
	return string(rune('a' + int(b%4)))
}

func fuzzValue(b byte) int {
	return int(b % 6)
}

// listMultimapFuzzFactories returns the list-backed multimap constructors to
// fuzz, each producing a fresh, empty MutableMultimap[string, int]. The
// concurrent variants share the plain ListMultimap's semantics, so they must
// satisfy the same oracle.
func listMultimapFuzzFactories() []struct {
	name string
	make func() multimaps.MutableMultimap[string, int]
} {
	return []struct {
		name string
		make func() multimaps.MutableMultimap[string, int]
	}{
		{"ListMultimap", func() multimaps.MutableMultimap[string, int] { return multimaps.NewListMultimap[string, int]() }},
		{"ConcurrentListMultimap", func() multimaps.MutableMultimap[string, int] {
			return multimaps.NewConcurrentListMultimap[string, int]()
		}},
		{"ConcurrentRWListMultimap", func() multimaps.MutableMultimap[string, int] {
			return multimaps.NewConcurrentRWListMultimap[string, int]()
		}},
	}
}

// setMultimapFuzzFactories returns the set-backed multimap constructors to fuzz,
// each producing a fresh, empty MutableMultimap[string, int]. The concurrent
// variants share the plain SetMultimap's semantics, so they must satisfy the
// same oracle.
func setMultimapFuzzFactories() []struct {
	name string
	make func() multimaps.MutableMultimap[string, int]
} {
	return []struct {
		name string
		make func() multimaps.MutableMultimap[string, int]
	}{
		{"SetMultimap", func() multimaps.MutableMultimap[string, int] { return multimaps.NewSetMultimap[string, int]() }},
		{"ConcurrentSetMultimap", func() multimaps.MutableMultimap[string, int] {
			return multimaps.NewConcurrentSetMultimap[string, int]()
		}},
		{"ConcurrentRWSetMultimap", func() multimaps.MutableMultimap[string, int] {
			return multimaps.NewConcurrentRWSetMultimap[string, int]()
		}},
	}
}

// FuzzListMultimap drives every list-backed multimap backend with arbitrary
// operation sequences and checks it against a native map[string][]int oracle
// that models list-backed (ordered, duplicate-keeping) semantics.
func FuzzListMultimap(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{opPut, 0, 0})
	f.Add([]byte{opPut, 0, 0, opPut, 0, 0, opRemove, 0, 0})
	f.Add([]byte{opPut, 1, 2, opRemoveAll, 1, 0, opFilter, 0, 0})

	f.Fuzz(func(t *testing.T, data []byte) {
		for _, backend := range listMultimapFuzzFactories() {
			m := backend.make()
			model := map[string][]int{}

			for i := 0; i+2 < len(data); i += 3 {
				op := int(data[i]) % opCount
				key := fuzzKey(data[i+1])
				value := fuzzValue(data[i+2])

				switch op {
				case opPut:
					m.PutInPlace(key, value)
					model[key] = append(model[key], value)
				case opRemove:
					removed := m.RemoveInPlace(key, value)
					modelRemoved := removeFirst(model, key, value)
					if removed != modelRemoved {
						t.Fatalf("%s: RemoveInPlace(%q,%d) = %v, model = %v", backend.name, key, value, removed, modelRemoved)
					}
				case opRemoveAll:
					_, ok := m.RemoveAllInPlace(key)
					_, modelOK := model[key]
					delete(model, key)
					if ok != modelOK {
						t.Fatalf("%s: RemoveAllInPlace(%q) ok = %v, model = %v", backend.name, key, ok, modelOK)
					}
				case opFilter:
					m.FilterInPlace(func(_ string, v int) bool { return v%2 == 0 })
					filterModel(model, func(v int) bool { return v%2 == 0 })
				}

				assertListMatchesModel(t, backend.name, m, model)
			}
		}
	})
}

// FuzzSetMultimap drives every set-backed multimap backend with arbitrary
// operation sequences and checks it against a native map[string]map[int]struct{}
// oracle that models set-backed (deduping) semantics.
func FuzzSetMultimap(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{opPut, 0, 0})
	f.Add([]byte{opPut, 0, 0, opPut, 0, 0, opRemove, 0, 0})
	f.Add([]byte{opPut, 1, 2, opRemoveAll, 1, 0, opFilter, 0, 0})

	f.Fuzz(func(t *testing.T, data []byte) {
		for _, backend := range setMultimapFuzzFactories() {
			m := backend.make()
			model := map[string]map[int]struct{}{}

			for i := 0; i+2 < len(data); i += 3 {
				op := int(data[i]) % opCount
				key := fuzzKey(data[i+1])
				value := fuzzValue(data[i+2])

				switch op {
				case opPut:
					m.PutInPlace(key, value)
					if model[key] == nil {
						model[key] = map[int]struct{}{}
					}
					model[key][value] = struct{}{}
				case opRemove:
					removed := m.RemoveInPlace(key, value)
					modelRemoved := removeFromSet(model, key, value)
					if removed != modelRemoved {
						t.Fatalf("%s: RemoveInPlace(%q,%d) = %v, model = %v", backend.name, key, value, removed, modelRemoved)
					}
				case opRemoveAll:
					_, ok := m.RemoveAllInPlace(key)
					_, modelOK := model[key]
					delete(model, key)
					if ok != modelOK {
						t.Fatalf("%s: RemoveAllInPlace(%q) ok = %v, model = %v", backend.name, key, ok, modelOK)
					}
				case opFilter:
					m.FilterInPlace(func(_ string, v int) bool { return v%2 == 0 })
					filterSetModel(model, func(v int) bool { return v%2 == 0 })
				}

				assertSetMatchesModel(t, backend.name, m, model)
			}
		}
	})
}

func removeFirst(model map[string][]int, key string, value int) bool {
	values, ok := model[key]
	if !ok {
		return false
	}
	for index, existing := range values {
		if existing != value {
			continue
		}
		remaining := append(values[:index], values[index+1:]...)
		if len(remaining) == 0 {
			delete(model, key)
		} else {
			model[key] = remaining
		}
		return true
	}
	return false
}

func removeFromSet(model map[string]map[int]struct{}, key string, value int) bool {
	values, ok := model[key]
	if !ok {
		return false
	}
	_, found := values[value]
	if !found {
		return false
	}
	delete(values, value)
	if len(values) == 0 {
		delete(model, key)
	}
	return true
}

func filterModel(model map[string][]int, keep func(int) bool) {
	for key, values := range model {
		kept := make([]int, 0, len(values))
		for _, v := range values {
			if keep(v) {
				kept = append(kept, v)
			}
		}
		if len(kept) == 0 {
			delete(model, key)
		} else {
			model[key] = kept
		}
	}
}

func filterSetModel(model map[string]map[int]struct{}, keep func(int) bool) {
	for key, values := range model {
		for v := range values {
			if !keep(v) {
				delete(values, v)
			}
		}
		if len(values) == 0 {
			delete(model, key)
		}
	}
}

func assertListMatchesModel(t *testing.T, name string, m multimaps.Multimap[string, int], model map[string][]int) {
	t.Helper()
	modelLen := 0
	for _, values := range model {
		modelLen += len(values)
	}
	if m.Length() != modelLen {
		t.Fatalf("%s: Length() = %d, model = %d", name, m.Length(), modelLen)
	}
	if m.KeyCount() != len(model) {
		t.Fatalf("%s: KeyCount() = %d, model = %d", name, m.KeyCount(), len(model))
	}
	// List values must match the model in exact insertion order: the
	// list-multimap contract guarantees per-key values come back in the order
	// they were inserted, so the oracle compares without sorting.
	for key, values := range model {
		got := m.Get(key)
		if !reflect.DeepEqual(got, values) {
			t.Fatalf("%s: Get(%q) = %v, model = %v", name, key, got, values)
		}
	}

	// ListMultimapFromSeq2(All) round-trips back to the same entries, again in
	// insertion order.
	rt := multimaps.ListMultimapFromSeq2(m.All())
	if rt.Length() != m.Length() {
		t.Fatalf("%s: FromSeq2 round-trip Length() = %d, want %d", name, rt.Length(), m.Length())
	}
	for key, values := range model {
		if got := rt.Get(key); !reflect.DeepEqual(got, values) {
			t.Fatalf("%s: FromSeq2 round-trip Get(%q) = %v, model = %v", name, key, got, values)
		}
	}
}

func assertSetMatchesModel(t *testing.T, name string, m multimaps.Multimap[string, int], model map[string]map[int]struct{}) {
	t.Helper()
	modelLen := 0
	for _, values := range model {
		modelLen += len(values)
	}
	if m.Length() != modelLen {
		t.Fatalf("%s: Length() = %d, model = %d", name, m.Length(), modelLen)
	}
	if m.KeyCount() != len(model) {
		t.Fatalf("%s: KeyCount() = %d, model = %d", name, m.KeyCount(), len(model))
	}
	for key, values := range model {
		want := make([]int, 0, len(values))
		for v := range values {
			want = append(want, v)
		}
		sort.Ints(want)
		got := sortedInts(m.Get(key))
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("%s: Get(%q) = %v, model = %v", name, key, got, want)
		}
	}

	// SetMultimapFromSeq2(All) round-trips back to the same entries.
	rt := multimaps.SetMultimapFromSeq2(m.All())
	if rt.Length() != m.Length() {
		t.Fatalf("%s: FromSeq2 round-trip Length() = %d, want %d", name, rt.Length(), m.Length())
	}
	for key := range model {
		if got, want := sortedInts(rt.Get(key)), sortedInts(m.Get(key)); !reflect.DeepEqual(got, want) {
			t.Fatalf("%s: FromSeq2 round-trip Get(%q) = %v, want %v", name, key, got, want)
		}
	}
}
