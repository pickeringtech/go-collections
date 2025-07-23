package dicts

// Hash is a dictionary implementation using Go's built-in map.
// It provides an immutable interface where operations return new instances.
type Hash[K comparable, V any] map[K]V

// NewHash creates a new Hash dictionary with the given key-value pairs.
func NewHash[K comparable, V any](entries ...Pair[K, V]) Hash[K, V] {
	m := make(Hash[K, V])
	for _, entry := range entries {
		m[entry.Key] = entry.Value
	}
	return m
}

// Interface guards to ensure Hash implements the required interfaces
var _ Dict[string, int] = Hash[string, int]{}
var _ MutableDict[string, int] = Hash[string, int]{}

// Get retrieves the value associated with the given key.
// If the key is not found, returns the default value and false.
func (h Hash[K, V]) Get(key K, defaultValue V) (V, bool) {
	if value, exists := h[key]; exists {
		return value, true
	}
	return defaultValue, false
}

// Contains checks if the given key exists in the dictionary.
func (h Hash[K, V]) Contains(key K) bool {
	_, exists := h[key]
	return exists
}

// Length returns the number of key-value pairs in the dictionary.
func (h Hash[K, V]) Length() int {
	return len(h)
}

// IsEmpty returns true if the dictionary contains no key-value pairs.
func (h Hash[K, V]) IsEmpty() bool {
	return len(h) == 0
}

// ForEach executes the given function for each key-value pair.
func (h Hash[K, V]) ForEach(fn func(key K, value V)) {
	for key, value := range h {
		fn(key, value)
	}
}

// ForEachKey executes the given function for each key.
func (h Hash[K, V]) ForEachKey(fn func(key K)) {
	for key := range h {
		fn(key)
	}
}

// ForEachValue executes the given function for each value.
func (h Hash[K, V]) ForEachValue(fn func(value V)) {
	for _, value := range h {
		fn(value)
	}
}

// Filter returns a new dictionary containing only the key-value pairs
// that satisfy the given predicate function.
func (h Hash[K, V]) Filter(fn func(key K, value V) bool) Dict[K, V] {
	result := make(Hash[K, V])
	for key, value := range h {
		if fn(key, value) {
			result[key] = value
		}
	}
	return result
}

// FilterInPlace removes all key-value pairs that do not satisfy
// the given predicate function, modifying the dictionary in place.
func (h Hash[K, V]) FilterInPlace(fn func(key K, value V) bool) {
	for key, value := range h {
		if !fn(key, value) {
			delete(h, key)
		}
	}
}

// Find returns the first key-value pair that satisfies the given predicate.
// Returns the key, value, and true if found; zero values and false otherwise.
func (h Hash[K, V]) Find(fn func(key K, value V) bool) (K, V, bool) {
	for key, value := range h {
		if fn(key, value) {
			return key, value, true
		}
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

// FindKey returns the first key that satisfies the given predicate.
// Returns the key and true if found; zero value and false otherwise.
func (h Hash[K, V]) FindKey(fn func(key K) bool) (K, bool) {
	for key := range h {
		if fn(key) {
			return key, true
		}
	}
	var zeroK K
	return zeroK, false
}

// FindValue returns the first value that satisfies the given predicate.
// Returns the value and true if found; zero value and false otherwise.
func (h Hash[K, V]) FindValue(fn func(value V) bool) (V, bool) {
	for _, value := range h {
		if fn(value) {
			return value, true
		}
	}
	var zeroV V
	return zeroV, false
}

// ContainsValue checks if the given value exists in the dictionary.
func (h Hash[K, V]) ContainsValue(value V) bool {
	for _, v := range h {
		// Note: This requires V to be comparable for equality check
		// For non-comparable types, this would need a different approach
		if any(v) == any(value) {
			return true
		}
	}
	return false
}

// Keys returns a slice containing all keys in the dictionary.
func (h Hash[K, V]) Keys() []K {
	keys := make([]K, 0, len(h))
	for key := range h {
		keys = append(keys, key)
	}
	return keys
}

// Values returns a slice containing all values in the dictionary.
func (h Hash[K, V]) Values() []V {
	values := make([]V, 0, len(h))
	for _, value := range h {
		values = append(values, value)
	}
	return values
}

// Items returns a slice containing all key-value pairs as Pair structs.
func (h Hash[K, V]) Items() []Pair[K, V] {
	items := make([]Pair[K, V], 0, len(h))
	for key, value := range h {
		items = append(items, Pair[K, V]{Key: key, Value: value})
	}
	return items
}

// AsMap returns the dictionary as a native Go map.
func (h Hash[K, V]) AsMap() map[K]V {
	result := make(map[K]V, len(h))
	for key, value := range h {
		result[key] = value
	}
	return result
}

// Put creates a new dictionary with the given key-value pair added or updated.
// Returns the new dictionary without modifying the original.
func (h Hash[K, V]) Put(key K, value V) Dict[K, V] {
	result := make(Hash[K, V], len(h)+1)
	for k, v := range h {
		result[k] = v
	}
	result[key] = value
	return result
}

// PutMany creates a new dictionary with all given key-value pairs added or updated.
// Returns the new dictionary without modifying the original.
func (h Hash[K, V]) PutMany(pairs ...Pair[K, V]) Dict[K, V] {
	result := make(Hash[K, V], len(h)+len(pairs))
	for k, v := range h {
		result[k] = v
	}
	for _, pair := range pairs {
		result[pair.Key] = pair.Value
	}
	return result
}

// PutInPlace adds or updates the given key-value pair in the dictionary.
func (h Hash[K, V]) PutInPlace(key K, value V) {
	h[key] = value
}

// PutManyInPlace adds or updates all given key-value pairs in the dictionary.
func (h Hash[K, V]) PutManyInPlace(pairs ...Pair[K, V]) {
	for _, pair := range pairs {
		h[pair.Key] = pair.Value
	}
}

// Remove creates a new dictionary with the given key removed.
// Returns the new dictionary without modifying the original.
func (h Hash[K, V]) Remove(key K) Dict[K, V] {
	result := make(Hash[K, V], len(h))
	for k, v := range h {
		if k != key {
			result[k] = v
		}
	}
	return result
}

// RemoveMany creates a new dictionary with all given keys removed.
// Returns the new dictionary without modifying the original.
func (h Hash[K, V]) RemoveMany(keys ...K) Dict[K, V] {
	// Create a set of keys to remove for efficient lookup
	toRemove := make(map[K]struct{}, len(keys))
	for _, key := range keys {
		toRemove[key] = struct{}{}
	}

	result := make(Hash[K, V], len(h))
	for k, v := range h {
		if _, shouldRemove := toRemove[k]; !shouldRemove {
			result[k] = v
		}
	}
	return result
}

// RemoveInPlace removes the given key from the dictionary.
// Returns the removed value and true if the key existed; zero value and false otherwise.
func (h Hash[K, V]) RemoveInPlace(key K) (V, bool) {
	if value, exists := h[key]; exists {
		delete(h, key)
		return value, true
	}
	var zeroV V
	return zeroV, false
}

// RemoveManyInPlace removes all given keys from the dictionary.
func (h Hash[K, V]) RemoveManyInPlace(keys ...K) {
	for _, key := range keys {
		delete(h, key)
	}
}

// Clear removes all key-value pairs from the dictionary.
func (h Hash[K, V]) Clear() {
	for key := range h {
		delete(h, key)
	}
}
