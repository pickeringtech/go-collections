package dicts

import (
	"reflect"
	"sync"
)

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

// ForEach executes the given function for each key-value pair. fn is invoked
// after the lock is released, against a point-in-time snapshot taken under the
// lock, so fn may safely call back into the collection.
func (ch *ConcurrentHashRW[K, V]) ForEach(fn func(key K, value V)) {
	ch.lock.RLock()
	items := make([]Pair[K, V], 0, len(ch.data))
	for key, value := range ch.data {
		items = append(items, Pair[K, V]{Key: key, Value: value})
	}
	ch.lock.RUnlock()

	for _, item := range items {
		fn(item.Key, item.Value)
	}
}

// ForEachKey executes the given function for each key. fn is invoked after the
// lock is released, against a point-in-time snapshot taken under the lock, so
// fn may safely call back into the collection.
func (ch *ConcurrentHashRW[K, V]) ForEachKey(fn func(key K)) {
	ch.lock.RLock()
	keys := make([]K, 0, len(ch.data))
	for key := range ch.data {
		keys = append(keys, key)
	}
	ch.lock.RUnlock()

	for _, key := range keys {
		fn(key)
	}
}

// ForEachValue executes the given function for each value. fn is invoked after
// the lock is released, against a point-in-time snapshot taken under the lock,
// so fn may safely call back into the collection.
func (ch *ConcurrentHashRW[K, V]) ForEachValue(fn func(value V)) {
	ch.lock.RLock()
	values := make([]V, 0, len(ch.data))
	for _, value := range ch.data {
		values = append(values, value)
	}
	ch.lock.RUnlock()

	for _, value := range values {
		fn(value)
	}
}

// Filter returns a new dictionary containing only the key-value pairs
// that satisfy the given predicate function. The returned dictionary is a new
// thread-safe ConcurrentHashRW, independent of the receiver. The predicate is
// evaluated after the lock is released, against a point-in-time snapshot taken
// under the lock, so it may safely call back into the collection.
func (ch *ConcurrentHashRW[K, V]) Filter(fn func(key K, value V) bool) Dict[K, V] {
	ch.lock.RLock()
	items := make([]Pair[K, V], 0, len(ch.data))
	for key, value := range ch.data {
		items = append(items, Pair[K, V]{Key: key, Value: value})
	}
	ch.lock.RUnlock()

	result := NewConcurrentHashRW[K, V]()
	for _, item := range items {
		if fn(item.Key, item.Value) {
			result.data[item.Key] = item.Value
		}
	}
	return result
}

// FilterInPlace removes all key-value pairs that do not satisfy
// the given predicate function, modifying the dictionary in place. The
// predicate is evaluated after the lock is released, against a point-in-time
// snapshot taken under the lock, so it may safely call back into the
// collection. Modifications made concurrently with evaluation are not
// reflected in the retained set.
func (ch *ConcurrentHashRW[K, V]) FilterInPlace(fn func(key K, value V) bool) {
	ch.lock.Lock()
	items := make([]Pair[K, V], 0, len(ch.data))
	for key, value := range ch.data {
		items = append(items, Pair[K, V]{Key: key, Value: value})
	}
	ch.lock.Unlock()

	var toRemove []K
	for _, item := range items {
		if !fn(item.Key, item.Value) {
			toRemove = append(toRemove, item.Key)
		}
	}

	ch.lock.Lock()
	for _, key := range toRemove {
		delete(ch.data, key)
	}
	ch.lock.Unlock()
}

// AllMatch returns true if every key-value pair satisfies the given predicate.
// It is vacuously true for an empty dictionary. The predicate is evaluated
// after the lock is released, against a point-in-time snapshot taken under the
// lock, so it may safely call back into the collection.
func (ch *ConcurrentHashRW[K, V]) AllMatch(fn func(key K, value V) bool) bool {
	ch.lock.RLock()
	items := make([]Pair[K, V], 0, len(ch.data))
	for key, value := range ch.data {
		items = append(items, Pair[K, V]{Key: key, Value: value})
	}
	ch.lock.RUnlock()

	for _, item := range items {
		if !fn(item.Key, item.Value) {
			return false
		}
	}
	return true
}

// AnyMatch returns true if at least one key-value pair satisfies the given
// predicate. It is false for an empty dictionary. The predicate is evaluated
// after the lock is released, against a point-in-time snapshot taken under the
// lock, so it may safely call back into the collection.
func (ch *ConcurrentHashRW[K, V]) AnyMatch(fn func(key K, value V) bool) bool {
	ch.lock.RLock()
	items := make([]Pair[K, V], 0, len(ch.data))
	for key, value := range ch.data {
		items = append(items, Pair[K, V]{Key: key, Value: value})
	}
	ch.lock.RUnlock()

	for _, item := range items {
		if fn(item.Key, item.Value) {
			return true
		}
	}
	return false
}

// NoneMatch returns true if no key-value pair satisfies the given predicate.
// It is vacuously true for an empty dictionary. The predicate is evaluated
// after the lock is released, against a point-in-time snapshot taken under the
// lock, so it may safely call back into the collection.
func (ch *ConcurrentHashRW[K, V]) NoneMatch(fn func(key K, value V) bool) bool {
	ch.lock.RLock()
	items := make([]Pair[K, V], 0, len(ch.data))
	for key, value := range ch.data {
		items = append(items, Pair[K, V]{Key: key, Value: value})
	}
	ch.lock.RUnlock()

	for _, item := range items {
		if fn(item.Key, item.Value) {
			return false
		}
	}
	return true
}

// Find returns the first key-value pair that satisfies the given predicate.
// Returns the key, value, and true if found; zero values and false otherwise.
// The predicate is evaluated after the lock is released, against a
// point-in-time snapshot taken under the lock, so it may safely call back into
// the collection.
func (ch *ConcurrentHashRW[K, V]) Find(fn func(key K, value V) bool) (K, V, bool) {
	ch.lock.RLock()
	items := make([]Pair[K, V], 0, len(ch.data))
	for key, value := range ch.data {
		items = append(items, Pair[K, V]{Key: key, Value: value})
	}
	ch.lock.RUnlock()

	for _, item := range items {
		if fn(item.Key, item.Value) {
			return item.Key, item.Value, true
		}
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

// FindKey returns the first key that satisfies the given predicate.
// Returns the key and true if found; zero value and false otherwise. The
// predicate is evaluated after the lock is released, against a point-in-time
// snapshot taken under the lock, so it may safely call back into the
// collection.
func (ch *ConcurrentHashRW[K, V]) FindKey(fn func(key K) bool) (K, bool) {
	ch.lock.RLock()
	keys := make([]K, 0, len(ch.data))
	for key := range ch.data {
		keys = append(keys, key)
	}
	ch.lock.RUnlock()

	for _, key := range keys {
		if fn(key) {
			return key, true
		}
	}
	var zeroK K
	return zeroK, false
}

// FindValue returns the first value that satisfies the given predicate.
// Returns the value and true if found; zero value and false otherwise. The
// predicate is evaluated after the lock is released, against a point-in-time
// snapshot taken under the lock, so it may safely call back into the
// collection.
func (ch *ConcurrentHashRW[K, V]) FindValue(fn func(value V) bool) (V, bool) {
	ch.lock.RLock()
	values := make([]V, 0, len(ch.data))
	for _, value := range ch.data {
		values = append(values, value)
	}
	ch.lock.RUnlock()

	for _, value := range values {
		if fn(value) {
			return value, true
		}
	}
	var zeroV V
	return zeroV, false
}

// ContainsValue checks if the given value exists in the dictionary.
//
// Values are compared with reflect.DeepEqual, matching the equality semantics
// used by list removal. This supports non-comparable value types (slices, maps,
// funcs) without panicking.
func (ch *ConcurrentHashRW[K, V]) ContainsValue(value V) bool {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	for _, v := range ch.data {
		if reflect.DeepEqual(v, value) {
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
// Returns a new thread-safe ConcurrentHashRW without modifying the original.
func (ch *ConcurrentHashRW[K, V]) Put(key K, value V) Dict[K, V] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	result := NewConcurrentHashRW[K, V]()
	for k, v := range ch.data {
		result.data[k] = v
	}
	result.data[key] = value
	return result
}

// PutMany creates a new dictionary with all given key-value pairs added or updated.
// Returns a new thread-safe ConcurrentHashRW without modifying the original.
func (ch *ConcurrentHashRW[K, V]) PutMany(pairs ...Pair[K, V]) Dict[K, V] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	result := NewConcurrentHashRW[K, V]()
	for k, v := range ch.data {
		result.data[k] = v
	}
	for _, pair := range pairs {
		result.data[pair.Key] = pair.Value
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
// Returns a new thread-safe ConcurrentHashRW without modifying the original.
func (ch *ConcurrentHashRW[K, V]) Remove(key K) Dict[K, V] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	result := NewConcurrentHashRW[K, V]()
	for k, v := range ch.data {
		if k != key {
			result.data[k] = v
		}
	}
	return result
}

// RemoveMany creates a new dictionary with all given keys removed.
// Returns a new thread-safe ConcurrentHashRW without modifying the original.
func (ch *ConcurrentHashRW[K, V]) RemoveMany(keys ...K) Dict[K, V] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	// Create a set of keys to remove for efficient lookup
	toRemove := make(map[K]struct{}, len(keys))
	for _, key := range keys {
		toRemove[key] = struct{}{}
	}

	result := NewConcurrentHashRW[K, V]()
	for k, v := range ch.data {
		if _, shouldRemove := toRemove[k]; !shouldRemove {
			result.data[k] = v
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
