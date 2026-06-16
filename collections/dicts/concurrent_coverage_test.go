package dicts_test

import (
	"sort"
	"sync"
	"testing"

	"github.com/pickeringtech/go-collections/collections/dicts"
)

// concurrentDict captures the subset of the MutableDict interface that both
// ConcurrentHash and ConcurrentHashRW satisfy, so we can table-drive the same
// behavioural tests across both implementations.
type concurrentFactory func(entries ...dicts.Pair[string, int]) dicts.MutableDict[string, int]

func concurrentFactories() map[string]concurrentFactory {
	return map[string]concurrentFactory{
		"ConcurrentHash": func(entries ...dicts.Pair[string, int]) dicts.MutableDict[string, int] {
			return dicts.NewConcurrentHash(entries...)
		},
		"ConcurrentHashRW": func(entries ...dicts.Pair[string, int]) dicts.MutableDict[string, int] {
			return dicts.NewConcurrentHashRW(entries...)
		},
	}
}

func sampleEntries() []dicts.Pair[string, int] {
	return []dicts.Pair[string, int]{
		{Key: "a", Value: 1},
		{Key: "b", Value: 2},
		{Key: "c", Value: 3},
	}
}

func TestConcurrent_IsEmpty(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			empty := factory()
			if !empty.IsEmpty() {
				t.Error("IsEmpty() = false, want true for empty dict")
			}

			nonEmpty := factory(sampleEntries()...)
			if nonEmpty.IsEmpty() {
				t.Error("IsEmpty() = true, want false for populated dict")
			}
		})
	}
}

func TestConcurrent_Length(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory(sampleEntries()...)
			if d.Length() != 3 {
				t.Errorf("Length() = %d, want 3", d.Length())
			}
		})
	}
}

func TestConcurrent_Get(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory(sampleEntries()...)

			v, ok := d.Get("a", -1)
			if !ok || v != 1 {
				t.Errorf("Get(a) = %d, %v, want 1, true", v, ok)
			}

			v, ok = d.Get("missing", -1)
			if ok || v != -1 {
				t.Errorf("Get(missing) = %d, %v, want -1, false", v, ok)
			}
		})
	}
}

func TestConcurrent_Contains(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory(sampleEntries()...)
			if !d.Contains("b") {
				t.Error("Contains(b) = false, want true")
			}
			if d.Contains("missing") {
				t.Error("Contains(missing) = true, want false")
			}
		})
	}
}

func TestConcurrent_ForEach(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory(sampleEntries()...)

			seen := make(map[string]int)
			d.ForEach(func(key string, value int) {
				seen[key] = value
			})

			if len(seen) != 3 || seen["a"] != 1 || seen["b"] != 2 || seen["c"] != 3 {
				t.Errorf("ForEach visited %v, want a=1 b=2 c=3", seen)
			}
		})
	}
}

func TestConcurrent_ForEachKey(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory(sampleEntries()...)

			var keys []string
			d.ForEachKey(func(key string) {
				keys = append(keys, key)
			})
			sort.Strings(keys)

			want := []string{"a", "b", "c"}
			if len(keys) != 3 || keys[0] != want[0] || keys[1] != want[1] || keys[2] != want[2] {
				t.Errorf("ForEachKey keys = %v, want %v", keys, want)
			}
		})
	}
}

func TestConcurrent_ForEachValue(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory(sampleEntries()...)

			sum := 0
			d.ForEachValue(func(value int) {
				sum += value
			})

			if sum != 6 {
				t.Errorf("ForEachValue sum = %d, want 6", sum)
			}
		})
	}
}

func TestConcurrent_Filter(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory(sampleEntries()...)

			filtered := d.Filter(func(key string, value int) bool {
				return value%2 == 1
			})

			if filtered.Length() != 2 {
				t.Errorf("Filter length = %d, want 2", filtered.Length())
			}
			if !filtered.Contains("a") || !filtered.Contains("c") {
				t.Error("Filter should retain odd-valued keys a and c")
			}
			// Original unchanged.
			if d.Length() != 3 {
				t.Errorf("original length = %d, want 3 (Filter must not mutate)", d.Length())
			}
		})
	}
}

func TestConcurrent_FilterInPlace(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory(sampleEntries()...)

			d.FilterInPlace(func(key string, value int) bool {
				return value > 1
			})

			if d.Length() != 2 {
				t.Errorf("FilterInPlace length = %d, want 2", d.Length())
			}
			if d.Contains("a") {
				t.Error("FilterInPlace should have removed key a")
			}
		})
	}
}

func TestConcurrent_Find(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory(sampleEntries()...)

			k, v, ok := d.Find(func(key string, value int) bool {
				return value == 2
			})
			if !ok || k != "b" || v != 2 {
				t.Errorf("Find = %q, %d, %v, want b, 2, true", k, v, ok)
			}

			k, v, ok = d.Find(func(key string, value int) bool {
				return value == 99
			})
			if ok || k != "" || v != 0 {
				t.Errorf("Find(missing) = %q, %d, %v, want \"\", 0, false", k, v, ok)
			}
		})
	}
}

func TestConcurrent_FindKey(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory(sampleEntries()...)

			k, ok := d.FindKey(func(key string) bool {
				return key == "c"
			})
			if !ok || k != "c" {
				t.Errorf("FindKey = %q, %v, want c, true", k, ok)
			}

			k, ok = d.FindKey(func(key string) bool {
				return key == "zzz"
			})
			if ok || k != "" {
				t.Errorf("FindKey(missing) = %q, %v, want \"\", false", k, ok)
			}
		})
	}
}

func TestConcurrent_FindValue(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory(sampleEntries()...)

			v, ok := d.FindValue(func(value int) bool {
				return value == 3
			})
			if !ok || v != 3 {
				t.Errorf("FindValue = %d, %v, want 3, true", v, ok)
			}

			v, ok = d.FindValue(func(value int) bool {
				return value == 99
			})
			if ok || v != 0 {
				t.Errorf("FindValue(missing) = %d, %v, want 0, false", v, ok)
			}
		})
	}
}

func TestConcurrent_ContainsValue(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory(sampleEntries()...)

			if !d.ContainsValue(2) {
				t.Error("ContainsValue(2) = false, want true")
			}
			if d.ContainsValue(99) {
				t.Error("ContainsValue(99) = true, want false")
			}
		})
	}
}

func TestConcurrent_Keys(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory(sampleEntries()...)

			keys := d.Keys()
			sort.Strings(keys)
			want := []string{"a", "b", "c"}
			if len(keys) != 3 || keys[0] != want[0] || keys[1] != want[1] || keys[2] != want[2] {
				t.Errorf("Keys() = %v, want %v", keys, want)
			}
		})
	}
}

func TestConcurrent_Values(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory(sampleEntries()...)

			values := d.Values()
			sort.Ints(values)
			want := []int{1, 2, 3}
			if len(values) != 3 || values[0] != want[0] || values[1] != want[1] || values[2] != want[2] {
				t.Errorf("Values() = %v, want %v", values, want)
			}
		})
	}
}

func TestConcurrent_Items(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory(sampleEntries()...)

			items := d.Items()
			if len(items) != 3 {
				t.Fatalf("Items() len = %d, want 3", len(items))
			}
			got := make(map[string]int)
			for _, item := range items {
				got[item.Key] = item.Value
			}
			if got["a"] != 1 || got["b"] != 2 || got["c"] != 3 {
				t.Errorf("Items() = %v, want a=1 b=2 c=3", got)
			}
		})
	}
}

func TestConcurrent_AsMap(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory(sampleEntries()...)

			m := d.AsMap()
			if len(m) != 3 || m["a"] != 1 || m["b"] != 2 || m["c"] != 3 {
				t.Errorf("AsMap() = %v, want a=1 b=2 c=3", m)
			}

			// Returned map is an independent copy.
			m["a"] = 100
			v, _ := d.Get("a", -1)
			if v != 1 {
				t.Error("AsMap returned map is not independent of the dict")
			}
		})
	}
}

func TestConcurrent_Put(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory(sampleEntries()...)

			result := d.Put("d", 4)
			if result.Length() != 4 {
				t.Errorf("Put result length = %d, want 4", result.Length())
			}
			if d.Length() != 3 {
				t.Errorf("original length = %d, want 3 (Put must not mutate)", d.Length())
			}
		})
	}
}

func TestConcurrent_PutMany(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory(sampleEntries()...)

			result := d.PutMany(
				dicts.Pair[string, int]{Key: "d", Value: 4},
				dicts.Pair[string, int]{Key: "e", Value: 5},
			)
			if result.Length() != 5 {
				t.Errorf("PutMany result length = %d, want 5", result.Length())
			}
			if d.Length() != 3 {
				t.Errorf("original length = %d, want 3 (PutMany must not mutate)", d.Length())
			}
		})
	}
}

func TestConcurrent_PutInPlace(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory(sampleEntries()...)

			d.PutInPlace("d", 4)
			v, ok := d.Get("d", -1)
			if !ok || v != 4 {
				t.Errorf("after PutInPlace Get(d) = %d, %v, want 4, true", v, ok)
			}
		})
	}
}

func TestConcurrent_PutManyInPlace(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory(sampleEntries()...)

			d.PutManyInPlace(
				dicts.Pair[string, int]{Key: "d", Value: 4},
				dicts.Pair[string, int]{Key: "e", Value: 5},
			)
			if d.Length() != 5 {
				t.Errorf("Length after PutManyInPlace = %d, want 5", d.Length())
			}
		})
	}
}

func TestConcurrent_Remove(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory(sampleEntries()...)

			result := d.Remove("a")
			if result.Length() != 2 {
				t.Errorf("Remove result length = %d, want 2", result.Length())
			}
			if result.Contains("a") {
				t.Error("Remove result should not contain a")
			}
			if d.Length() != 3 {
				t.Errorf("original length = %d, want 3 (Remove must not mutate)", d.Length())
			}
		})
	}
}

func TestConcurrent_RemoveMany(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory(sampleEntries()...)

			result := d.RemoveMany("a", "b")
			if result.Length() != 1 {
				t.Errorf("RemoveMany result length = %d, want 1", result.Length())
			}
			if !result.Contains("c") {
				t.Error("RemoveMany result should retain c")
			}
			if d.Length() != 3 {
				t.Errorf("original length = %d, want 3 (RemoveMany must not mutate)", d.Length())
			}
		})
	}
}

func TestConcurrent_RemoveInPlace(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory(sampleEntries()...)

			v, ok := d.RemoveInPlace("b")
			if !ok || v != 2 {
				t.Errorf("RemoveInPlace(b) = %d, %v, want 2, true", v, ok)
			}
			if d.Contains("b") {
				t.Error("RemoveInPlace should have removed b")
			}

			v, ok = d.RemoveInPlace("missing")
			if ok || v != 0 {
				t.Errorf("RemoveInPlace(missing) = %d, %v, want 0, false", v, ok)
			}
		})
	}
}

func TestConcurrent_RemoveManyInPlace(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory(sampleEntries()...)

			d.RemoveManyInPlace("a", "c")
			if d.Length() != 1 {
				t.Errorf("Length after RemoveManyInPlace = %d, want 1", d.Length())
			}
			if !d.Contains("b") {
				t.Error("RemoveManyInPlace should retain b")
			}
		})
	}
}

func TestConcurrent_Clear(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory(sampleEntries()...)

			d.Clear()
			if !d.IsEmpty() {
				t.Error("Clear should leave the dict empty")
			}
			if d.Length() != 0 {
				t.Errorf("Length after Clear = %d, want 0", d.Length())
			}
		})
	}
}

// TestConcurrent_RaceAccess exercises each implementation from many goroutines
// at once so the -race detector can flag any unsynchronised access.
func TestConcurrent_RaceAccess(t *testing.T) {
	for name, factory := range concurrentFactories() {
		t.Run(name, func(t *testing.T) {
			d := factory()

			const goroutines = 16
			const iterations = 200

			var wg sync.WaitGroup
			wg.Add(goroutines)

			for g := 0; g < goroutines; g++ {
				go func(id int) {
					defer wg.Done()
					for i := 0; i < iterations; i++ {
						key := "k"
						d.PutInPlace(key, i)
						_, _ = d.Get(key, -1)
						_ = d.Length()
						_ = d.IsEmpty()
						d.ForEach(func(string, int) {})
						_ = d.Keys()
						_ = d.Values()
						_ = d.AsMap()
						_ = d.Filter(func(string, int) bool { return true })
						_, _ = d.RemoveInPlace(key)
					}
				}(g)
			}

			wg.Wait()
		})
	}
}
