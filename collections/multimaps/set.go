package multimaps

import "iter"

// SetMultimap is a set-backed multimap: each key maps to a set of distinct
// values. Binding the same value to a key more than once is a no-op, so each
// (key, value) pair appears at most once.
//
// Use SetMultimap when duplicate bindings are meaningless and fast membership
// checks matter (for example, the set of tags applied to a document). For
// ordered values or meaningful duplicates, use ListMultimap instead.
//
// Both key and value iteration order is unspecified (it follows Go's map
// iteration).
//
// Example usage:
//
//	// Tag documents, ignoring duplicate tags.
//	tags := multimaps.NewSetMultimap[string, string]()
//	tags.PutInPlace("doc1", "go")
//	tags.PutInPlace("doc1", "go") // duplicate ignored
//	tags.PutInPlace("doc1", "testing")
//	count := tags.Length() // 2
type SetMultimap[K comparable, V comparable] map[K]map[V]struct{}

// NewSetMultimap creates a new set-backed multimap seeded with the given
// entries. Duplicate (key, value) pairs are collapsed.
//
// Example:
//
//	// Empty multimap
//	empty := multimaps.NewSetMultimap[string, int]()
//
//	// Seeded with entries (duplicates collapsed)
//	groups := multimaps.NewSetMultimap(
//		multimaps.Entry[string, int]{Key: "even", Value: 2},
//		multimaps.Entry[string, int]{Key: "even", Value: 4},
//	)
func NewSetMultimap[K comparable, V comparable](entries ...Entry[K, V]) SetMultimap[K, V] {
	m := make(SetMultimap[K, V])
	for _, entry := range entries {
		m.PutInPlace(entry.Key, entry.Value)
	}
	return m
}

// Interface guards to ensure SetMultimap implements the required interfaces.
var _ Multimap[string, int] = SetMultimap[string, int]{}
var _ MutableMultimap[string, int] = SetMultimap[string, int]{}

// Get returns a copy of all values bound to the given key. Returns an empty
// (non-nil) slice if the key has no values. Order is unspecified.
func (m SetMultimap[K, V]) Get(key K) []V {
	return valuesSlice(m[key])
}

// ContainsKey reports whether the given key has at least one value bound to it.
func (m SetMultimap[K, V]) ContainsKey(key K) bool {
	_, exists := m[key]
	return exists
}

// ContainsEntry reports whether the given key is bound to the given value.
func (m SetMultimap[K, V]) ContainsEntry(key K, value V) bool {
	values, exists := m[key]
	if !exists {
		return false
	}
	_, found := values[value]
	return found
}

// Length returns the total number of entries (distinct key-value associations).
func (m SetMultimap[K, V]) Length() int {
	total := 0
	for _, values := range m {
		total += len(values)
	}
	return total
}

// KeyCount returns the number of distinct keys.
func (m SetMultimap[K, V]) KeyCount() int {
	return len(m)
}

// IsEmpty returns true if the multimap contains no entries.
func (m SetMultimap[K, V]) IsEmpty() bool {
	return len(m) == 0
}

// ForEach executes the given function once for every entry.
func (m SetMultimap[K, V]) ForEach(fn func(key K, value V)) {
	for key, values := range m {
		for value := range values {
			fn(key, value)
		}
	}
}

// ForEachKey executes the given function once per distinct key, passing a copy
// of all values bound to that key.
func (m SetMultimap[K, V]) ForEachKey(fn func(key K, values []V)) {
	for key, values := range m {
		fn(key, valuesSlice(values))
	}
}

// All returns an iterator over every entry, suitable for range-over-func.
func (m SetMultimap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for key, values := range m {
			for value := range values {
				if !yield(key, value) {
					return
				}
			}
		}
	}
}

// KeysSeq returns an iterator over the distinct keys, suitable for
// range-over-func.
func (m SetMultimap[K, V]) KeysSeq() iter.Seq[K] {
	return func(yield func(K) bool) {
		for key := range m {
			if !yield(key) {
				return
			}
		}
	}
}

// Filter returns a new multimap containing only the entries that satisfy the
// predicate. Keys left with no values are dropped.
func (m SetMultimap[K, V]) Filter(fn func(key K, value V) bool) Multimap[K, V] {
	result := make(SetMultimap[K, V])
	for key, values := range m {
		for value := range values {
			if fn(key, value) {
				result.PutInPlace(key, value)
			}
		}
	}
	return result
}

// FilterInPlace removes every entry that does not satisfy the predicate. Keys
// left with no values are dropped.
func (m SetMultimap[K, V]) FilterInPlace(fn func(key K, value V) bool) {
	for key, values := range m {
		for value := range values {
			if !fn(key, value) {
				delete(values, value)
			}
		}
		if len(values) == 0 {
			delete(m, key)
		}
	}
}

// AllMatch returns true if every entry satisfies the predicate.
func (m SetMultimap[K, V]) AllMatch(fn func(key K, value V) bool) bool {
	for key, values := range m {
		for value := range values {
			if !fn(key, value) {
				return false
			}
		}
	}
	return true
}

// AnyMatch returns true if at least one entry satisfies the predicate.
func (m SetMultimap[K, V]) AnyMatch(fn func(key K, value V) bool) bool {
	for key, values := range m {
		for value := range values {
			if fn(key, value) {
				return true
			}
		}
	}
	return false
}

// NoneMatch returns true if no entry satisfies the predicate.
func (m SetMultimap[K, V]) NoneMatch(fn func(key K, value V) bool) bool {
	return !m.AnyMatch(fn)
}

// Find returns the first entry that satisfies the predicate.
func (m SetMultimap[K, V]) Find(fn func(key K, value V) bool) (K, V, bool) {
	for key, values := range m {
		for value := range values {
			if fn(key, value) {
				return key, value, true
			}
		}
	}
	var zeroKey K
	var zeroValue V
	return zeroKey, zeroValue, false
}

// Keys returns a slice containing each distinct key once.
func (m SetMultimap[K, V]) Keys() []K {
	result := make([]K, 0, len(m))
	for key := range m {
		result = append(result, key)
	}
	return result
}

// Values returns a slice containing every value across all keys.
func (m SetMultimap[K, V]) Values() []V {
	result := make([]V, 0, m.Length())
	for _, values := range m {
		for value := range values {
			result = append(result, value)
		}
	}
	return result
}

// Entries returns a slice containing every entry.
func (m SetMultimap[K, V]) Entries() []Entry[K, V] {
	result := make([]Entry[K, V], 0, m.Length())
	for key, values := range m {
		for value := range values {
			result = append(result, Entry[K, V]{Key: key, Value: value})
		}
	}
	return result
}

// AsMap returns the multimap as a native Go map from each key to a copy of its
// values.
func (m SetMultimap[K, V]) AsMap() map[K][]V {
	result := make(map[K][]V, len(m))
	for key, values := range m {
		result[key] = valuesSlice(values)
	}
	return result
}

// Put returns a new multimap with the given value bound to the given key. If the
// binding already exists the result is an equivalent copy.
func (m SetMultimap[K, V]) Put(key K, value V) Multimap[K, V] {
	result := m.clone()
	result.PutInPlace(key, value)
	return result
}

// PutAll returns a new multimap with all the given values bound to the given
// key. With no values it returns an equivalent copy.
func (m SetMultimap[K, V]) PutAll(key K, values ...V) Multimap[K, V] {
	result := m.clone()
	result.PutAllInPlace(key, values...)
	return result
}

// PutInPlace binds the given value to the given key, modifying the multimap in
// place. Binding an existing value is a no-op.
func (m SetMultimap[K, V]) PutInPlace(key K, value V) {
	values, exists := m[key]
	if !exists {
		values = make(map[V]struct{})
		m[key] = values
	}
	values[value] = struct{}{}
}

// PutAllInPlace binds all the given values to the given key, modifying the
// multimap in place.
func (m SetMultimap[K, V]) PutAllInPlace(key K, values ...V) {
	if len(values) == 0 {
		return
	}
	existing, exists := m[key]
	if !exists {
		existing = make(map[V]struct{})
		m[key] = existing
	}
	for _, value := range values {
		existing[value] = struct{}{}
	}
}

// Remove returns a new multimap with the given value unbound from the given key.
func (m SetMultimap[K, V]) Remove(key K, value V) Multimap[K, V] {
	result := m.clone()
	result.RemoveInPlace(key, value)
	return result
}

// RemoveAll returns a new multimap with the given key and all of its values
// removed.
func (m SetMultimap[K, V]) RemoveAll(key K) Multimap[K, V] {
	result := m.clone()
	delete(result, key)
	return result
}

// RemoveInPlace unbinds the given value from the given key. Returns true if the
// binding was present and removed.
func (m SetMultimap[K, V]) RemoveInPlace(key K, value V) bool {
	values, exists := m[key]
	if !exists {
		return false
	}
	_, found := values[value]
	if !found {
		return false
	}
	delete(values, value)
	if len(values) == 0 {
		delete(m, key)
	}
	return true
}

// RemoveAllInPlace removes the given key and all of its values. Returns the
// removed values and true if the key was present.
func (m SetMultimap[K, V]) RemoveAllInPlace(key K) ([]V, bool) {
	values, exists := m[key]
	if !exists {
		return []V{}, false
	}
	delete(m, key)
	return valuesSlice(values), true
}

// Clear removes all entries from the multimap.
func (m SetMultimap[K, V]) Clear() {
	for key := range m {
		delete(m, key)
	}
}

// clone returns a deep copy of the multimap, with independent value sets.
func (m SetMultimap[K, V]) clone() SetMultimap[K, V] {
	result := make(SetMultimap[K, V], len(m))
	for key, values := range m {
		cloned := make(map[V]struct{}, len(values))
		for value := range values {
			cloned[value] = struct{}{}
		}
		result[key] = cloned
	}
	return result
}

// valuesSlice returns an independent, non-nil slice of the given value set.
func valuesSlice[V comparable](values map[V]struct{}) []V {
	result := make([]V, 0, len(values))
	for value := range values {
		result = append(result, value)
	}
	return result
}
