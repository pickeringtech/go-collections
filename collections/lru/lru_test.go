package lru_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/collections/lru"
)

// factory builds a fresh MutableCache for one of the three implementations, so
// every behavioural test runs identically against LRU, ConcurrentLRU and
// ConcurrentLRURW — "learn one, learn all" verified, not just claimed.
type factory struct {
	name string
	make func(capacity int, opts ...lru.Option[string, int]) lru.MutableCache[string, int]
}

func factories() []factory {
	return []factory{
		{"LRU", func(c int, o ...lru.Option[string, int]) lru.MutableCache[string, int] {
			return lru.NewLRU(c, o...)
		}},
		{"ConcurrentLRU", func(c int, o ...lru.Option[string, int]) lru.MutableCache[string, int] {
			return lru.NewConcurrentLRU(c, o...)
		}},
		{"ConcurrentLRURW", func(c int, o ...lru.Option[string, int]) lru.MutableCache[string, int] {
			return lru.NewConcurrentLRURW(c, o...)
		}},
	}
}

func pair(key string, value int) lru.Pair[string, int] {
	return lru.Pair[string, int]{Key: key, Value: value}
}

func TestCache_PutInPlaceAndPeek(t *testing.T) {
	type op struct {
		key   string
		value int
	}
	tests := []struct {
		name     string
		capacity int
		puts     []op
		wantKeys []int // values in most- to least-recently-used order
	}{
		{
			name:     "single put",
			capacity: 3,
			puts:     []op{{"a", 1}},
			wantKeys: []int{1},
		},
		{
			name:     "puts order newest first",
			capacity: 3,
			puts:     []op{{"a", 1}, {"b", 2}, {"c", 3}},
			wantKeys: []int{3, 2, 1},
		},
		{
			name:     "overwrite updates value and promotes",
			capacity: 3,
			puts:     []op{{"a", 1}, {"b", 2}, {"a", 11}},
			wantKeys: []int{11, 2},
		},
		{
			name:     "eviction drops least-recently-used",
			capacity: 2,
			puts:     []op{{"a", 1}, {"b", 2}, {"c", 3}},
			wantKeys: []int{3, 2},
		},
	}
	for _, f := range factories() {
		for _, tt := range tests {
			t.Run(f.name+"/"+tt.name, func(t *testing.T) {
				cache := f.make(tt.capacity)
				for _, p := range tt.puts {
					cache.PutInPlace(p.key, p.value)
				}
				got := cache.Values()
				if !reflect.DeepEqual(got, tt.wantKeys) {
					t.Errorf("Values() = %v, want %v", got, tt.wantKeys)
				}
			})
		}
	}
}

func TestCache_Get(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			cache := f.make(2, lru.WithEntries(pair("a", 1), pair("b", 2)))

			value, found := cache.Get("a")
			if value != 1 || !found {
				t.Errorf("Get(a) = (%d, %t), want (1, true)", value, found)
			}

			// Get promoted "a", so the next insert evicts "b".
			cache.PutInPlace("c", 3)
			if cache.Contains("b") {
				t.Errorf("expected b to be evicted after promoting a")
			}
			if !cache.Contains("a") {
				t.Errorf("expected a to survive after promotion")
			}

			missing, found := cache.Get("zzz")
			if missing != 0 || found {
				t.Errorf("Get(missing) = (%d, %t), want (0, false)", missing, found)
			}
		})
	}
}

func TestCache_Peek_DoesNotPromote(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			cache := f.make(2, lru.WithEntries(pair("a", 1), pair("b", 2)))

			value, found := cache.Peek("a")
			if value != 1 || !found {
				t.Errorf("Peek(a) = (%d, %t), want (1, true)", value, found)
			}

			// Peek did not promote "a", so it remains the eviction candidate.
			cache.PutInPlace("c", 3)
			if cache.Contains("a") {
				t.Errorf("expected a to be evicted; Peek must not promote")
			}

			missing, found := cache.Peek("zzz")
			if missing != 0 || found {
				t.Errorf("Peek(missing) = (%d, %t), want (0, false)", missing, found)
			}
		})
	}
}

func TestCache_ContainsLengthIsEmptyCapacity(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			cache := f.make(5)
			if !cache.IsEmpty() {
				t.Errorf("new cache should be empty")
			}
			if cache.Length() != 0 {
				t.Errorf("Length() = %d, want 0", cache.Length())
			}
			if cache.Capacity() != 5 {
				t.Errorf("Capacity() = %d, want 5", cache.Capacity())
			}
			if cache.Contains("a") {
				t.Errorf("empty cache should not contain a")
			}

			cache.PutInPlace("a", 1)
			if cache.IsEmpty() {
				t.Errorf("cache with an entry should not be empty")
			}
			if cache.Length() != 1 {
				t.Errorf("Length() = %d, want 1", cache.Length())
			}
			if !cache.Contains("a") {
				t.Errorf("cache should contain a")
			}
		})
	}
}

func TestCache_CapacityBelowOneClampsToOne(t *testing.T) {
	for _, f := range factories() {
		for _, capacity := range []int{0, -5} {
			t.Run(f.name, func(t *testing.T) {
				cache := f.make(capacity)
				if cache.Capacity() != 1 {
					t.Errorf("Capacity() = %d, want 1 for input %d", cache.Capacity(), capacity)
				}
				cache.PutInPlace("a", 1)
				cache.PutInPlace("b", 2)
				if cache.Length() != 1 {
					t.Errorf("Length() = %d, want 1 (capacity clamped)", cache.Length())
				}
				if !cache.Contains("b") || cache.Contains("a") {
					t.Errorf("expected only b to remain")
				}
			})
		}
	}
}

func TestCache_RemoveInPlace(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			cache := f.make(3, lru.WithEntries(pair("a", 1), pair("b", 2), pair("c", 3)))

			value, ok := cache.RemoveInPlace("b")
			if value != 2 || !ok {
				t.Errorf("RemoveInPlace(b) = (%d, %t), want (2, true)", value, ok)
			}
			if cache.Contains("b") {
				t.Errorf("b should be gone after removal")
			}
			if got := cache.Values(); !reflect.DeepEqual(got, []int{3, 1}) {
				t.Errorf("Values() = %v, want [3 1]", got)
			}

			missing, ok := cache.RemoveInPlace("zzz")
			if missing != 0 || ok {
				t.Errorf("RemoveInPlace(missing) = (%d, %t), want (0, false)", missing, ok)
			}
		})
	}
}

func TestCache_Clear(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			cache := f.make(3, lru.WithEntries(pair("a", 1), pair("b", 2)))
			cache.Clear()
			if !cache.IsEmpty() {
				t.Errorf("cache should be empty after Clear")
			}
			if got := cache.Keys(); !reflect.DeepEqual(got, []string{}) {
				t.Errorf("Keys() = %v, want []", got)
			}
			// Cache is reusable after Clear.
			cache.PutInPlace("x", 9)
			if !cache.Contains("x") {
				t.Errorf("cache should be usable after Clear")
			}
		})
	}
}

func TestCache_KeysValuesItems(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			cache := f.make(3, lru.WithEntries(pair("a", 1), pair("b", 2), pair("c", 3)))

			if got := cache.Keys(); !reflect.DeepEqual(got, []string{"c", "b", "a"}) {
				t.Errorf("Keys() = %v, want [c b a]", got)
			}
			if got := cache.Values(); !reflect.DeepEqual(got, []int{3, 2, 1}) {
				t.Errorf("Values() = %v, want [3 2 1]", got)
			}
			wantItems := []lru.Pair[string, int]{pair("c", 3), pair("b", 2), pair("a", 1)}
			if got := cache.Items(); !reflect.DeepEqual(got, wantItems) {
				t.Errorf("Items() = %v, want %v", got, wantItems)
			}
		})
	}
}

func TestCache_KeysValuesItemsEmpty(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			cache := f.make(3)
			if got := cache.Keys(); !reflect.DeepEqual(got, []string{}) {
				t.Errorf("Keys() = %v, want []", got)
			}
			if got := cache.Values(); !reflect.DeepEqual(got, []int{}) {
				t.Errorf("Values() = %v, want []", got)
			}
			if got := cache.Items(); !reflect.DeepEqual(got, []lru.Pair[string, int]{}) {
				t.Errorf("Items() = %v, want []", got)
			}
			if got := cache.AsMap(); !reflect.DeepEqual(got, map[string]int{}) {
				t.Errorf("AsMap() = %v, want empty map", got)
			}
		})
	}
}

func TestCache_AsMap(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			cache := f.make(3, lru.WithEntries(pair("a", 1), pair("b", 2)))
			got := cache.AsMap()
			want := map[string]int{"a": 1, "b": 2}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("AsMap() = %v, want %v", got, want)
			}
		})
	}
}

func TestCache_ForEach(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			cache := f.make(3, lru.WithEntries(pair("a", 1), pair("b", 2), pair("c", 3)))
			var keys []string
			var values []int
			cache.ForEach(func(key string, value int) {
				keys = append(keys, key)
				values = append(values, value)
			})
			if !reflect.DeepEqual(keys, []string{"c", "b", "a"}) {
				t.Errorf("ForEach keys = %v, want [c b a]", keys)
			}
			if !reflect.DeepEqual(values, []int{3, 2, 1}) {
				t.Errorf("ForEach values = %v, want [3 2 1]", values)
			}
		})
	}
}

func TestCache_All(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			cache := f.make(3, lru.WithEntries(pair("a", 1), pair("b", 2), pair("c", 3)))
			collected := map[string]int{}
			var order []string
			for key, value := range cache.All() {
				collected[key] = value
				order = append(order, key)
			}
			if !reflect.DeepEqual(order, []string{"c", "b", "a"}) {
				t.Errorf("All order = %v, want [c b a]", order)
			}
			if !reflect.DeepEqual(collected, map[string]int{"a": 1, "b": 2, "c": 3}) {
				t.Errorf("All values = %v", collected)
			}
		})
	}
}

func TestCache_All_EarlyBreak(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			cache := f.make(3, lru.WithEntries(pair("a", 1), pair("b", 2), pair("c", 3)))
			var visited []string
			for key := range cache.All() {
				visited = append(visited, key)
				break // exercises the yield-returns-false path
			}
			if !reflect.DeepEqual(visited, []string{"c"}) {
				t.Errorf("early-break visited = %v, want [c]", visited)
			}
		})
	}
}

func TestCache_OnEvictCallback(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			var evicted []lru.Pair[string, int]
			cache := f.make(2, lru.WithOnEvict(func(key string, value int) {
				evicted = append(evicted, pair(key, value))
			}))
			cache.PutInPlace("a", 1)
			cache.PutInPlace("b", 2)
			cache.PutInPlace("c", 3) // evicts a
			cache.PutInPlace("d", 4) // evicts b

			want := []lru.Pair[string, int]{pair("a", 1), pair("b", 2)}
			if !reflect.DeepEqual(evicted, want) {
				t.Errorf("evicted = %v, want %v", evicted, want)
			}
		})
	}
}

func TestCache_RemoveAndClearDoNotFireEvict(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			fired := false
			cache := f.make(3, lru.WithOnEvict(func(string, int) { fired = true }))
			cache.PutInPlace("a", 1)
			cache.PutInPlace("b", 2)
			cache.RemoveInPlace("a")
			cache.Clear()
			if fired {
				t.Errorf("eviction callback must not fire for explicit Remove/Clear")
			}
		})
	}
}

func TestCache_ImmutablePut(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			base := f.make(2, lru.WithEntries(pair("a", 1), pair("b", 2)))
			updated := base.Put("c", 3)

			// Receiver is unchanged.
			if !base.Contains("a") {
				t.Errorf("Put must not modify the receiver")
			}
			if base.Length() != 2 {
				t.Errorf("receiver Length = %d, want 2", base.Length())
			}
			// Result reflects the insert and the eviction.
			if updated.Contains("a") {
				t.Errorf("updated cache should have evicted a")
			}
			if got := updated.Keys(); !reflect.DeepEqual(got, []string{"c", "b"}) {
				t.Errorf("updated Keys = %v, want [c b]", got)
			}
		})
	}
}

func TestCache_ImmutableRemove(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			base := f.make(3, lru.WithEntries(pair("a", 1), pair("b", 2), pair("c", 3)))
			updated := base.Remove("b")

			if !base.Contains("b") {
				t.Errorf("Remove must not modify the receiver")
			}
			if updated.Contains("b") {
				t.Errorf("updated cache should not contain b")
			}
			if got := updated.Keys(); !reflect.DeepEqual(got, []string{"c", "a"}) {
				t.Errorf("updated Keys = %v, want [c a]", got)
			}
			// Removing an absent key yields an equivalent copy.
			same := base.Remove("zzz")
			if !reflect.DeepEqual(same.Keys(), base.Keys()) {
				t.Errorf("Remove(absent) Keys = %v, want %v", same.Keys(), base.Keys())
			}
		})
	}
}

func TestCache_ImmutablePutPreservesEvictCallback(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			var evicted []string
			base := f.make(1, lru.WithOnEvict(func(key string, _ int) {
				evicted = append(evicted, key)
			}))
			base.PutInPlace("a", 1)

			// The immutable Put exceeds capacity on the copy and must fire the
			// inherited callback for the evicted entry.
			base.Put("b", 2)
			if !reflect.DeepEqual(evicted, []string{"a"}) {
				t.Errorf("evicted = %v, want [a] from immutable Put", evicted)
			}
		})
	}
}

// TestCache_AsMapMatchesKeys checks AsMap and the ordered exporters agree on
// contents (ignoring order, since maps are unordered).
func TestCache_AsMapMatchesKeys(t *testing.T) {
	for _, f := range factories() {
		t.Run(f.name, func(t *testing.T) {
			cache := f.make(4, lru.WithEntries(pair("a", 1), pair("b", 2), pair("c", 3)))
			keys := cache.Keys()
			sort.Strings(keys)
			mapKeys := make([]string, 0, len(cache.AsMap()))
			for k := range cache.AsMap() {
				mapKeys = append(mapKeys, k)
			}
			sort.Strings(mapKeys)
			if !reflect.DeepEqual(keys, mapKeys) {
				t.Errorf("AsMap keys = %v, want %v", mapKeys, keys)
			}
		})
	}
}
