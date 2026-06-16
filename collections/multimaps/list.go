package multimaps

import "iter"

// ListMultimap is a list-backed multimap: each key maps to an ordered slice of
// values that preserves insertion order and allows the same value to be bound to
// a key more than once.
//
// Use ListMultimap when value order matters or duplicate bindings are
// meaningful (for example, an ordered log of events per user). For
// duplicate-free value collections, use SetMultimap instead.
//
// Key iteration order is unspecified (it follows Go's map iteration); the order
// of values within a single key is the order in which they were inserted.
//
// Example usage:
//
//	// Group orders by customer, preserving order and duplicates.
//	orders := multimaps.NewListMultimap[string, string]()
//	orders.PutInPlace("alice", "book")
//	orders.PutInPlace("alice", "pen")
//	orders.PutInPlace("alice", "book") // duplicate kept
//	alice := orders.Get("alice")       // ["book", "pen", "book"]
type ListMultimap[K comparable, V comparable] map[K][]V

// NewListMultimap creates a new list-backed multimap seeded with the given
// entries, preserving their order and any duplicate bindings.
//
// Example:
//
//	// Empty multimap
//	empty := multimaps.NewListMultimap[string, int]()
//
//	// Seeded with entries
//	scores := multimaps.NewListMultimap(
//		multimaps.Entry[string, int]{Key: "alice", Value: 10},
//		multimaps.Entry[string, int]{Key: "alice", Value: 20},
//	)
func NewListMultimap[K comparable, V comparable](entries ...Entry[K, V]) ListMultimap[K, V] {
	m := make(ListMultimap[K, V])
	for _, entry := range entries {
		m[entry.Key] = append(m[entry.Key], entry.Value)
	}
	return m
}

// Interface guards to ensure ListMultimap implements the required interfaces.
var _ Multimap[string, int] = ListMultimap[string, int]{}
var _ MutableMultimap[string, int] = ListMultimap[string, int]{}

// Get returns a copy of all values bound to the given key, in insertion order.
// Returns an empty (non-nil) slice if the key has no values.
func (m ListMultimap[K, V]) Get(key K) []V {
	return cloneValues(m[key])
}

// ContainsKey reports whether the given key has at least one value bound to it.
func (m ListMultimap[K, V]) ContainsKey(key K) bool {
	_, exists := m[key]
	return exists
}

// ContainsEntry reports whether the given key is bound to the given value.
func (m ListMultimap[K, V]) ContainsEntry(key K, value V) bool {
	for _, existing := range m[key] {
		if existing == value {
			return true
		}
	}
	return false
}

// Length returns the total number of entries (key-value associations).
func (m ListMultimap[K, V]) Length() int {
	total := 0
	for _, values := range m {
		total += len(values)
	}
	return total
}

// KeyCount returns the number of distinct keys.
func (m ListMultimap[K, V]) KeyCount() int {
	return len(m)
}

// IsEmpty returns true if the multimap contains no entries.
func (m ListMultimap[K, V]) IsEmpty() bool {
	return len(m) == 0
}

// ForEach executes the given function once for every entry.
func (m ListMultimap[K, V]) ForEach(fn func(key K, value V)) {
	for key, values := range m {
		for _, value := range values {
			fn(key, value)
		}
	}
}

// ForEachKey executes the given function once per distinct key, passing a copy
// of all values bound to that key.
func (m ListMultimap[K, V]) ForEachKey(fn func(key K, values []V)) {
	for key, values := range m {
		fn(key, cloneValues(values))
	}
}

// All returns an iterator over every entry, suitable for range-over-func.
func (m ListMultimap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for key, values := range m {
			for _, value := range values {
				if !yield(key, value) {
					return
				}
			}
		}
	}
}

// KeysSeq returns an iterator over the distinct keys, suitable for
// range-over-func.
func (m ListMultimap[K, V]) KeysSeq() iter.Seq[K] {
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
func (m ListMultimap[K, V]) Filter(fn func(key K, value V) bool) Multimap[K, V] {
	result := make(ListMultimap[K, V])
	for key, values := range m {
		for _, value := range values {
			if fn(key, value) {
				result[key] = append(result[key], value)
			}
		}
	}
	return result
}

// FilterInPlace removes every entry that does not satisfy the predicate. Keys
// left with no values are dropped.
func (m ListMultimap[K, V]) FilterInPlace(fn func(key K, value V) bool) {
	for key, values := range m {
		kept := make([]V, 0, len(values))
		for _, value := range values {
			if fn(key, value) {
				kept = append(kept, value)
			}
		}
		if len(kept) == 0 {
			delete(m, key)
		} else {
			m[key] = kept
		}
	}
}

// AllMatch returns true if every entry satisfies the predicate.
func (m ListMultimap[K, V]) AllMatch(fn func(key K, value V) bool) bool {
	for key, values := range m {
		for _, value := range values {
			if !fn(key, value) {
				return false
			}
		}
	}
	return true
}

// AnyMatch returns true if at least one entry satisfies the predicate.
func (m ListMultimap[K, V]) AnyMatch(fn func(key K, value V) bool) bool {
	for key, values := range m {
		for _, value := range values {
			if fn(key, value) {
				return true
			}
		}
	}
	return false
}

// NoneMatch returns true if no entry satisfies the predicate.
func (m ListMultimap[K, V]) NoneMatch(fn func(key K, value V) bool) bool {
	return !m.AnyMatch(fn)
}

// Find returns the first entry that satisfies the predicate.
func (m ListMultimap[K, V]) Find(fn func(key K, value V) bool) (K, V, bool) {
	for key, values := range m {
		for _, value := range values {
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
func (m ListMultimap[K, V]) Keys() []K {
	result := make([]K, 0, len(m))
	for key := range m {
		result = append(result, key)
	}
	return result
}

// Values returns a slice containing every value across all keys.
func (m ListMultimap[K, V]) Values() []V {
	result := make([]V, 0, m.Length())
	for _, values := range m {
		result = append(result, values...)
	}
	return result
}

// Entries returns a slice containing every entry.
func (m ListMultimap[K, V]) Entries() []Entry[K, V] {
	result := make([]Entry[K, V], 0, m.Length())
	for key, values := range m {
		for _, value := range values {
			result = append(result, Entry[K, V]{Key: key, Value: value})
		}
	}
	return result
}

// AsMap returns the multimap as a native Go map from each key to a copy of its
// values.
func (m ListMultimap[K, V]) AsMap() map[K][]V {
	result := make(map[K][]V, len(m))
	for key, values := range m {
		result[key] = cloneValues(values)
	}
	return result
}

// Put returns a new multimap with the given value bound to the given key.
func (m ListMultimap[K, V]) Put(key K, value V) Multimap[K, V] {
	result := m.clone()
	result[key] = append(result[key], value)
	return result
}

// PutAll returns a new multimap with all the given values bound to the given
// key. With no values it returns an equivalent copy.
func (m ListMultimap[K, V]) PutAll(key K, values ...V) Multimap[K, V] {
	result := m.clone()
	if len(values) > 0 {
		result[key] = append(result[key], values...)
	}
	return result
}

// PutInPlace binds the given value to the given key, modifying the multimap in
// place.
func (m ListMultimap[K, V]) PutInPlace(key K, value V) {
	m[key] = append(m[key], value)
}

// PutAllInPlace binds all the given values to the given key, modifying the
// multimap in place.
func (m ListMultimap[K, V]) PutAllInPlace(key K, values ...V) {
	if len(values) > 0 {
		m[key] = append(m[key], values...)
	}
}

// Remove returns a new multimap with a single (first) binding of the given
// value to the given key removed.
func (m ListMultimap[K, V]) Remove(key K, value V) Multimap[K, V] {
	result := m.clone()
	result.RemoveInPlace(key, value)
	return result
}

// RemoveAll returns a new multimap with the given key and all of its values
// removed.
func (m ListMultimap[K, V]) RemoveAll(key K) Multimap[K, V] {
	result := m.clone()
	delete(result, key)
	return result
}

// RemoveInPlace removes a single (first) binding of the given value to the given
// key. Returns true if a binding was removed.
func (m ListMultimap[K, V]) RemoveInPlace(key K, value V) bool {
	values, exists := m[key]
	if !exists {
		return false
	}
	for index, existing := range values {
		if existing != value {
			continue
		}
		remaining := append(values[:index], values[index+1:]...)
		if len(remaining) == 0 {
			delete(m, key)
		} else {
			m[key] = remaining
		}
		return true
	}
	return false
}

// RemoveAllInPlace removes the given key and all of its values. Returns the
// removed values and true if the key was present.
func (m ListMultimap[K, V]) RemoveAllInPlace(key K) ([]V, bool) {
	values, exists := m[key]
	if !exists {
		return []V{}, false
	}
	delete(m, key)
	return cloneValues(values), true
}

// Clear removes all entries from the multimap.
func (m ListMultimap[K, V]) Clear() {
	for key := range m {
		delete(m, key)
	}
}

// clone returns a deep copy of the multimap, with independent value slices.
func (m ListMultimap[K, V]) clone() ListMultimap[K, V] {
	result := make(ListMultimap[K, V], len(m))
	for key, values := range m {
		result[key] = cloneValues(values)
	}
	return result
}

// cloneValues returns an independent, non-nil copy of the given values.
func cloneValues[V comparable](values []V) []V {
	result := make([]V, len(values))
	copy(result, values)
	return result
}
