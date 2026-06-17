package dicts

import (
	"reflect"
	"sync"
)

// ConcurrentHash is a thread-safe dictionary implementation using Go's built-in map
// with a mutex for synchronization. All operations are protected by a single mutex.
//
// Zero value: always construct with NewConcurrentHash. The embedded mutex is a
// value, so a bare &ConcurrentHash{} is at least lock-safe, but its backing map
// is nil until the constructor runs, so writes (PutInPlace) panic. Reads on the
// zero value return empty results.
type ConcurrentHash[K comparable, V any] struct {
	data map[K]V
	lock sync.Mutex
}

// NewConcurrentHash creates a new ConcurrentHash dictionary with the given key-value pairs.
func NewConcurrentHash[K comparable, V any](entries ...Pair[K, V]) *ConcurrentHash[K, V] {
	m := &ConcurrentHash[K, V]{
		data: make(map[K]V),
	}
	for _, entry := range entries {
		m.data[entry.Key] = entry.Value
	}
	return m
}

// Interface guards to ensure ConcurrentHash implements the required interfaces
var _ Dict[string, int] = &ConcurrentHash[string, int]{}
var _ MutableDict[string, int] = &ConcurrentHash[string, int]{}

// Get retrieves the value associated with the given key.
// If the key is not found, returns the default value and false.
func (ch *ConcurrentHash[K, V]) Get(key K, defaultValue V) (V, bool) {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	if value, exists := ch.data[key]; exists {
		return value, true
	}
	return defaultValue, false
}

// Contains checks if the given key exists in the dictionary.
func (ch *ConcurrentHash[K, V]) Contains(key K) bool {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	_, exists := ch.data[key]
	return exists
}

// Length returns the number of key-value pairs in the dictionary.
func (ch *ConcurrentHash[K, V]) Length() int {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	return len(ch.data)
}

// IsEmpty returns true if the dictionary contains no key-value pairs.
func (ch *ConcurrentHash[K, V]) IsEmpty() bool {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	return len(ch.data) == 0
}

// ForEach executes the given function for each key-value pair. fn is invoked
// after the lock is released, against a point-in-time snapshot taken under the
// lock, so fn may safely call back into the collection.
func (ch *ConcurrentHash[K, V]) ForEach(fn func(key K, value V)) {
	ch.lock.Lock()
	items := make([]Pair[K, V], 0, len(ch.data))
	for key, value := range ch.data {
		items = append(items, Pair[K, V]{Key: key, Value: value})
	}
	ch.lock.Unlock()

	for _, item := range items {
		fn(item.Key, item.Value)
	}
}

// ForEachKey executes the given function for each key. fn is invoked after the
// lock is released, against a point-in-time snapshot taken under the lock, so
// fn may safely call back into the collection.
func (ch *ConcurrentHash[K, V]) ForEachKey(fn func(key K)) {
	ch.lock.Lock()
	keys := make([]K, 0, len(ch.data))
	for key := range ch.data {
		keys = append(keys, key)
	}
	ch.lock.Unlock()

	for _, key := range keys {
		fn(key)
	}
}

// ForEachValue executes the given function for each value. fn is invoked after
// the lock is released, against a point-in-time snapshot taken under the lock,
// so fn may safely call back into the collection.
func (ch *ConcurrentHash[K, V]) ForEachValue(fn func(value V)) {
	ch.lock.Lock()
	values := make([]V, 0, len(ch.data))
	for _, value := range ch.data {
		values = append(values, value)
	}
	ch.lock.Unlock()

	for _, value := range values {
		fn(value)
	}
}

// Filter returns a new dictionary containing only the key-value pairs
// that satisfy the given predicate function. The returned dictionary is a new
// thread-safe ConcurrentHash, independent of the receiver. The predicate is
// evaluated after the lock is released, against a point-in-time snapshot taken
// under the lock, so it may safely call back into the collection.
func (ch *ConcurrentHash[K, V]) Filter(fn func(key K, value V) bool) Dict[K, V] {
	ch.lock.Lock()
	items := make([]Pair[K, V], 0, len(ch.data))
	for key, value := range ch.data {
		items = append(items, Pair[K, V]{Key: key, Value: value})
	}
	ch.lock.Unlock()

	result := NewConcurrentHash[K, V]()
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
func (ch *ConcurrentHash[K, V]) FilterInPlace(fn func(key K, value V) bool) {
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
func (ch *ConcurrentHash[K, V]) AllMatch(fn func(key K, value V) bool) bool {
	ch.lock.Lock()
	items := make([]Pair[K, V], 0, len(ch.data))
	for key, value := range ch.data {
		items = append(items, Pair[K, V]{Key: key, Value: value})
	}
	ch.lock.Unlock()

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
func (ch *ConcurrentHash[K, V]) AnyMatch(fn func(key K, value V) bool) bool {
	ch.lock.Lock()
	items := make([]Pair[K, V], 0, len(ch.data))
	for key, value := range ch.data {
		items = append(items, Pair[K, V]{Key: key, Value: value})
	}
	ch.lock.Unlock()

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
func (ch *ConcurrentHash[K, V]) NoneMatch(fn func(key K, value V) bool) bool {
	ch.lock.Lock()
	items := make([]Pair[K, V], 0, len(ch.data))
	for key, value := range ch.data {
		items = append(items, Pair[K, V]{Key: key, Value: value})
	}
	ch.lock.Unlock()

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
func (ch *ConcurrentHash[K, V]) Find(fn func(key K, value V) bool) (K, V, bool) {
	ch.lock.Lock()
	items := make([]Pair[K, V], 0, len(ch.data))
	for key, value := range ch.data {
		items = append(items, Pair[K, V]{Key: key, Value: value})
	}
	ch.lock.Unlock()

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
func (ch *ConcurrentHash[K, V]) FindKey(fn func(key K) bool) (K, bool) {
	ch.lock.Lock()
	keys := make([]K, 0, len(ch.data))
	for key := range ch.data {
		keys = append(keys, key)
	}
	ch.lock.Unlock()

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
func (ch *ConcurrentHash[K, V]) FindValue(fn func(value V) bool) (V, bool) {
	ch.lock.Lock()
	values := make([]V, 0, len(ch.data))
	for _, value := range ch.data {
		values = append(values, value)
	}
	ch.lock.Unlock()

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
//
// This is the deliberate counterpart to maps.ContainsValue, which uses == and
// requires a comparable V: dicts trades that speed for the ability to compare
// nested and non-comparable values structurally.
func (ch *ConcurrentHash[K, V]) ContainsValue(value V) bool {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	for _, v := range ch.data {
		if reflect.DeepEqual(v, value) {
			return true
		}
	}
	return false
}

// Keys returns a slice containing all keys in the dictionary.
func (ch *ConcurrentHash[K, V]) Keys() []K {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	keys := make([]K, 0, len(ch.data))
	for key := range ch.data {
		keys = append(keys, key)
	}
	return keys
}

// Values returns a slice containing all values in the dictionary.
func (ch *ConcurrentHash[K, V]) Values() []V {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	values := make([]V, 0, len(ch.data))
	for _, value := range ch.data {
		values = append(values, value)
	}
	return values
}

// Items returns a slice containing all key-value pairs as Pair structs.
func (ch *ConcurrentHash[K, V]) Items() []Pair[K, V] {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	items := make([]Pair[K, V], 0, len(ch.data))
	for key, value := range ch.data {
		items = append(items, Pair[K, V]{Key: key, Value: value})
	}
	return items
}

// AsMap returns the dictionary as a native Go map.
func (ch *ConcurrentHash[K, V]) AsMap() map[K]V {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	result := make(map[K]V, len(ch.data))
	for key, value := range ch.data {
		result[key] = value
	}
	return result
}

// Put creates a new dictionary with the given key-value pair added or updated.
// Returns a new thread-safe ConcurrentHash without modifying the original.
func (ch *ConcurrentHash[K, V]) Put(key K, value V) Dict[K, V] {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	result := NewConcurrentHash[K, V]()
	for k, v := range ch.data {
		result.data[k] = v
	}
	result.data[key] = value
	return result
}

// PutMany creates a new dictionary with all given key-value pairs added or updated.
// Returns a new thread-safe ConcurrentHash without modifying the original.
func (ch *ConcurrentHash[K, V]) PutMany(pairs ...Pair[K, V]) Dict[K, V] {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	result := NewConcurrentHash[K, V]()
	for k, v := range ch.data {
		result.data[k] = v
	}
	for _, pair := range pairs {
		result.data[pair.Key] = pair.Value
	}
	return result
}

// PutInPlace adds or updates the given key-value pair in the dictionary.
func (ch *ConcurrentHash[K, V]) PutInPlace(key K, value V) {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	ch.data[key] = value
}

// PutManyInPlace adds or updates all given key-value pairs in the dictionary.
func (ch *ConcurrentHash[K, V]) PutManyInPlace(pairs ...Pair[K, V]) {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	for _, pair := range pairs {
		ch.data[pair.Key] = pair.Value
	}
}

// UpdateInPlace atomically reads the value at key, applies fn to it, and stores
// the result back under key, returning the new value. fn receives the current
// value (the zero value if the key is absent) and whether the key existed. The
// whole read-modify-write runs under a single lock acquisition, so concurrent
// updates compose without losing writes. fn must not call back into the
// dictionary, which would deadlock on the held lock.
func (ch *ConcurrentHash[K, V]) UpdateInPlace(key K, fn func(old V, existed bool) V) V {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	old, existed := ch.data[key]
	newValue := fn(old, existed)
	ch.data[key] = newValue
	return newValue
}

// Remove creates a new dictionary with the given key removed.
// Returns a new thread-safe ConcurrentHash without modifying the original.
func (ch *ConcurrentHash[K, V]) Remove(key K) Dict[K, V] {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	result := NewConcurrentHash[K, V]()
	for k, v := range ch.data {
		if k != key {
			result.data[k] = v
		}
	}
	return result
}

// RemoveMany creates a new dictionary with all given keys removed.
// Returns a new thread-safe ConcurrentHash without modifying the original.
func (ch *ConcurrentHash[K, V]) RemoveMany(keys ...K) Dict[K, V] {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	// Create a set of keys to remove for efficient lookup
	toRemove := make(map[K]struct{}, len(keys))
	for _, key := range keys {
		toRemove[key] = struct{}{}
	}

	result := NewConcurrentHash[K, V]()
	for k, v := range ch.data {
		if _, shouldRemove := toRemove[k]; !shouldRemove {
			result.data[k] = v
		}
	}
	return result
}

// RemoveInPlace removes the given key from the dictionary.
// Returns the removed value and true if the key existed; zero value and false otherwise.
func (ch *ConcurrentHash[K, V]) RemoveInPlace(key K) (V, bool) {
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
func (ch *ConcurrentHash[K, V]) RemoveManyInPlace(keys ...K) {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	for _, key := range keys {
		delete(ch.data, key)
	}
}

// Clear removes all key-value pairs from the dictionary.
func (ch *ConcurrentHash[K, V]) Clear() {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	for key := range ch.data {
		delete(ch.data, key)
	}
}
