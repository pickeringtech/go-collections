package dicts

import "sync"

// ConcurrentHashRW is a thread-safe dictionary implementation using Go's built-in map
// with a read-write mutex for synchronization. Read operations use read locks for better
// performance when there are many concurrent readers.
type ConcurrentHashRW[K comparable, V any] struct {
	data map[K]V
	lock *sync.RWMutex
}

// NewConcurrentHashRW creates a new ConcurrentHashRW dictionary with the given key-value pairs.
func NewConcurrentHashRW[K comparable, V any](entries ...Pair[K, V]) *ConcurrentHashRW[K, V] {
	m := &ConcurrentHashRW[K, V]{
		data: make(map[K]V),
		lock: &sync.RWMutex{},
	}
	for _, entry := range entries {
		m.data[entry.Key] = entry.Value
	}
	return m
}

// Interface guards to ensure ConcurrentHashRW implements the required interfaces
var _ Dict[string, int] = &ConcurrentHashRW[string, int]{}
var _ MutableDict[string, int] = &ConcurrentHashRW[string, int]{}

// Get retrieves the value associated with the given key.
// If the key is not found, returns the default value and false.
func (ch *ConcurrentHashRW[K, V]) Get(key K, defaultValue V) (V, bool) {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	if value, exists := ch.data[key]; exists {
		return value, true
	}
	return defaultValue, false
}

// Contains checks if the given key exists in the dictionary.
func (ch *ConcurrentHashRW[K, V]) Contains(key K) bool {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	_, exists := ch.data[key]
	return exists
}

// Length returns the number of key-value pairs in the dictionary.
func (ch *ConcurrentHashRW[K, V]) Length() int {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	return len(ch.data)
}

// IsEmpty returns true if the dictionary contains no key-value pairs.
func (ch *ConcurrentHashRW[K, V]) IsEmpty() bool {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	return len(ch.data) == 0
}

// ForEach executes the given function for each key-value pair.
func (ch *ConcurrentHashRW[K, V]) ForEach(fn func(key K, value V)) {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	for key, value := range ch.data {
		fn(key, value)
	}
}

// ForEachKey executes the given function for each key.
func (ch *ConcurrentHashRW[K, V]) ForEachKey(fn func(key K)) {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	for key := range ch.data {
		fn(key)
	}
}

// ForEachValue executes the given function for each value.
func (ch *ConcurrentHashRW[K, V]) ForEachValue(fn func(value V)) {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	for _, value := range ch.data {
		fn(value)
	}
}

// Filter returns a new dictionary containing only the key-value pairs
// that satisfy the given predicate function.
func (ch *ConcurrentHashRW[K, V]) Filter(fn func(key K, value V) bool) Dict[K, V] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	result := make(Hash[K, V])
	for key, value := range ch.data {
		if fn(key, value) {
			result[key] = value
		}
	}
	return result
}

// FilterInPlace removes all key-value pairs that do not satisfy
// the given predicate function, modifying the dictionary in place.
func (ch *ConcurrentHashRW[K, V]) FilterInPlace(fn func(key K, value V) bool) {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	for key, value := range ch.data {
		if !fn(key, value) {
			delete(ch.data, key)
		}
	}
}

// Find returns the first key-value pair that satisfies the given predicate.
// Returns the key, value, and true if found; zero values and false otherwise.
func (ch *ConcurrentHashRW[K, V]) Find(fn func(key K, value V) bool) (K, V, bool) {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	for key, value := range ch.data {
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
func (ch *ConcurrentHashRW[K, V]) FindKey(fn func(key K) bool) (K, bool) {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	for key := range ch.data {
		if fn(key) {
			return key, true
		}
	}
	var zeroK K
	return zeroK, false
}

// FindValue returns the first value that satisfies the given predicate.
// Returns the value and true if found; zero value and false otherwise.
func (ch *ConcurrentHashRW[K, V]) FindValue(fn func(value V) bool) (V, bool) {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	for _, value := range ch.data {
		if fn(value) {
			return value, true
		}
	}
	var zeroV V
	return zeroV, false
}

// ContainsValue checks if the given value exists in the dictionary.
func (ch *ConcurrentHashRW[K, V]) ContainsValue(value V) bool {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	for _, v := range ch.data {
		if any(v) == any(value) {
			return true
		}
	}
	return false
}

// Keys returns a slice containing all keys in the dictionary.
func (ch *ConcurrentHashRW[K, V]) Keys() []K {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	keys := make([]K, 0, len(ch.data))
	for key := range ch.data {
		keys = append(keys, key)
	}
	return keys
}

// Values returns a slice containing all values in the dictionary.
func (ch *ConcurrentHashRW[K, V]) Values() []V {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	values := make([]V, 0, len(ch.data))
	for _, value := range ch.data {
		values = append(values, value)
	}
	return values
}

// Items returns a slice containing all key-value pairs as Pair structs.
func (ch *ConcurrentHashRW[K, V]) Items() []Pair[K, V] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	items := make([]Pair[K, V], 0, len(ch.data))
	for key, value := range ch.data {
		items = append(items, Pair[K, V]{Key: key, Value: value})
	}
	return items
}

// AsMap returns the dictionary as a native Go map.
func (ch *ConcurrentHashRW[K, V]) AsMap() map[K]V {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	result := make(map[K]V, len(ch.data))
	for key, value := range ch.data {
		result[key] = value
	}
	return result
}

// Put creates a new dictionary with the given key-value pair added or updated.
// Returns the new dictionary without modifying the original.
func (ch *ConcurrentHashRW[K, V]) Put(key K, value V) Dict[K, V] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	result := make(Hash[K, V], len(ch.data)+1)
	for k, v := range ch.data {
		result[k] = v
	}
	result[key] = value
	return result
}

// PutMany creates a new dictionary with all given key-value pairs added or updated.
// Returns the new dictionary without modifying the original.
func (ch *ConcurrentHashRW[K, V]) PutMany(pairs ...Pair[K, V]) Dict[K, V] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	result := make(Hash[K, V], len(ch.data)+len(pairs))
	for k, v := range ch.data {
		result[k] = v
	}
	for _, pair := range pairs {
		result[pair.Key] = pair.Value
	}
	return result
}

// PutInPlace adds or updates the given key-value pair in the dictionary.
func (ch *ConcurrentHashRW[K, V]) PutInPlace(key K, value V) {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	ch.data[key] = value
}

// PutManyInPlace adds or updates all given key-value pairs in the dictionary.
func (ch *ConcurrentHashRW[K, V]) PutManyInPlace(pairs ...Pair[K, V]) {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	for _, pair := range pairs {
		ch.data[pair.Key] = pair.Value
	}
}

// Remove creates a new dictionary with the given key removed.
// Returns the new dictionary without modifying the original.
func (ch *ConcurrentHashRW[K, V]) Remove(key K) Dict[K, V] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	result := make(Hash[K, V], len(ch.data))
	for k, v := range ch.data {
		if k != key {
			result[k] = v
		}
	}
	return result
}

// RemoveMany creates a new dictionary with all given keys removed.
// Returns the new dictionary without modifying the original.
func (ch *ConcurrentHashRW[K, V]) RemoveMany(keys ...K) Dict[K, V] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	// Create a set of keys to remove for efficient lookup
	toRemove := make(map[K]struct{}, len(keys))
	for _, key := range keys {
		toRemove[key] = struct{}{}
	}

	result := make(Hash[K, V], len(ch.data))
	for k, v := range ch.data {
		if _, shouldRemove := toRemove[k]; !shouldRemove {
			result[k] = v
		}
	}
	return result
}

// RemoveInPlace removes the given key from the dictionary.
// Returns the removed value and true if the key existed; zero value and false otherwise.
func (ch *ConcurrentHashRW[K, V]) RemoveInPlace(key K) (V, bool) {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	if value, exists := ch.data[key]; exists {
		delete(ch.data, key)
		return value, true
	}
	var zeroV V
	return zeroV, false
}

// RemoveManyInPlace removes all given keys from the dictionary.
func (ch *ConcurrentHashRW[K, V]) RemoveManyInPlace(keys ...K) {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	for _, key := range keys {
		delete(ch.data, key)
	}
}

// Clear removes all key-value pairs from the dictionary.
func (ch *ConcurrentHashRW[K, V]) Clear() {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	for key := range ch.data {
		delete(ch.data, key)
	}
}
