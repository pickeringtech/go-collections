package dicts

import (
	"iter"
	"sync"

	"github.com/pickeringtech/go-collections/constraints"
)

// ConcurrentTreeRW is a thread-safe ordered dictionary backed by a binary search
// tree (Tree) with a read-write mutex for synchronization. Read operations use
// read locks so many readers can proceed concurrently, while writes are
// exclusive. Keys are maintained in sorted order. Prefer it over ConcurrentTree
// for read-heavy workloads.
type ConcurrentTreeRW[K constraints.Ordered, V any] struct {
	tree *Tree[K, V]
	lock *sync.RWMutex
}

// NewConcurrentTreeRW creates a new ConcurrentTreeRW dictionary with the given key-value pairs.
func NewConcurrentTreeRW[K constraints.Ordered, V any](entries ...Pair[K, V]) *ConcurrentTreeRW[K, V] {
	return &ConcurrentTreeRW[K, V]{
		tree: NewTree[K, V](entries...),
		lock: &sync.RWMutex{},
	}
}

// Interface guards to ensure ConcurrentTreeRW implements the required interfaces.
var _ SortedDict[string, int] = &ConcurrentTreeRW[string, int]{}
var _ MutableSortedDict[string, int] = &ConcurrentTreeRW[string, int]{}

// wrapConcurrentTreeRW builds a new ConcurrentTreeRW, with its own lock, around the given tree.
func wrapConcurrentTreeRW[K constraints.Ordered, V any](tree *Tree[K, V]) *ConcurrentTreeRW[K, V] {
	return &ConcurrentTreeRW[K, V]{tree: tree, lock: &sync.RWMutex{}}
}

// Get retrieves the value associated with the given key.
// If the key is not found, returns the default value and false.
func (ch *ConcurrentTreeRW[K, V]) Get(key K, defaultValue V) (V, bool) {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return ch.tree.Get(key, defaultValue)
}

// Contains checks if the given key exists in the dictionary.
func (ch *ConcurrentTreeRW[K, V]) Contains(key K) bool {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return ch.tree.Contains(key)
}

// Length returns the number of key-value pairs in the dictionary.
func (ch *ConcurrentTreeRW[K, V]) Length() int {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return ch.tree.Length()
}

// IsEmpty returns true if the dictionary contains no key-value pairs.
func (ch *ConcurrentTreeRW[K, V]) IsEmpty() bool {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return ch.tree.IsEmpty()
}

// ForEach executes the given function for each key-value pair in sorted order.
// fn is invoked after the lock is released, against a point-in-time snapshot
// taken under the lock, so fn may safely call back into the collection.
func (ch *ConcurrentTreeRW[K, V]) ForEach(fn func(key K, value V)) {
	ch.lock.RLock()
	items := ch.tree.Items()
	ch.lock.RUnlock()

	for _, item := range items {
		fn(item.Key, item.Value)
	}
}

// ForEachKey executes the given function for each key in sorted order. fn is
// invoked after the lock is released, against a point-in-time snapshot taken
// under the lock, so fn may safely call back into the collection.
func (ch *ConcurrentTreeRW[K, V]) ForEachKey(fn func(key K)) {
	ch.lock.RLock()
	items := ch.tree.Items()
	ch.lock.RUnlock()

	for _, item := range items {
		fn(item.Key)
	}
}

// ForEachValue executes the given function for each value in key-sorted order.
// fn is invoked after the lock is released, against a point-in-time snapshot
// taken under the lock, so fn may safely call back into the collection.
func (ch *ConcurrentTreeRW[K, V]) ForEachValue(fn func(value V)) {
	ch.lock.RLock()
	items := ch.tree.Items()
	ch.lock.RUnlock()

	for _, item := range items {
		fn(item.Value)
	}
}

// Filter returns a new dictionary containing only the key-value pairs that
// satisfy the given predicate. The result is a new thread-safe ConcurrentTreeRW.
// The predicate is evaluated after the lock is released, against a
// point-in-time snapshot taken under the lock, so it may safely call back into
// the collection.
func (ch *ConcurrentTreeRW[K, V]) Filter(fn func(key K, value V) bool) Dict[K, V] {
	ch.lock.RLock()
	items := ch.tree.Items()
	ch.lock.RUnlock()

	var retained []Pair[K, V]
	for _, item := range items {
		if fn(item.Key, item.Value) {
			retained = append(retained, item)
		}
	}
	return wrapConcurrentTreeRW(NewTree[K, V](retained...))
}

// FilterInPlace removes all key-value pairs that do not satisfy the given
// predicate, modifying the dictionary in place. The predicate is evaluated
// after the lock is released, against a point-in-time snapshot taken under the
// lock, so it may safely call back into the collection. Modifications made
// concurrently with evaluation are not reflected in the retained set.
func (ch *ConcurrentTreeRW[K, V]) FilterInPlace(fn func(key K, value V) bool) {
	ch.lock.RLock()
	items := ch.tree.Items()
	ch.lock.RUnlock()

	var toRemove []K
	for _, item := range items {
		if !fn(item.Key, item.Value) {
			toRemove = append(toRemove, item.Key)
		}
	}

	ch.lock.Lock()
	ch.tree.RemoveManyInPlace(toRemove...)
	ch.lock.Unlock()
}

// AllMatch returns true if every key-value pair satisfies the given predicate.
// It is vacuously true for an empty dictionary. The predicate is evaluated
// after the lock is released, against a point-in-time snapshot taken under the
// lock, so it may safely call back into the collection.
func (ch *ConcurrentTreeRW[K, V]) AllMatch(fn func(key K, value V) bool) bool {
	ch.lock.RLock()
	items := ch.tree.Items()
	ch.lock.RUnlock()

	for _, item := range items {
		if !fn(item.Key, item.Value) {
			return false
		}
	}
	return true
}

// AnyMatch returns true if at least one key-value pair satisfies the given
// predicate. The predicate is evaluated after the lock is released, against a
// point-in-time snapshot taken under the lock, so it may safely call back into
// the collection.
func (ch *ConcurrentTreeRW[K, V]) AnyMatch(fn func(key K, value V) bool) bool {
	ch.lock.RLock()
	items := ch.tree.Items()
	ch.lock.RUnlock()

	for _, item := range items {
		if fn(item.Key, item.Value) {
			return true
		}
	}
	return false
}

// NoneMatch returns true if no key-value pair satisfies the given predicate.
// The predicate is evaluated after the lock is released, against a
// point-in-time snapshot taken under the lock, so it may safely call back into
// the collection.
func (ch *ConcurrentTreeRW[K, V]) NoneMatch(fn func(key K, value V) bool) bool {
	ch.lock.RLock()
	items := ch.tree.Items()
	ch.lock.RUnlock()

	for _, item := range items {
		if fn(item.Key, item.Value) {
			return false
		}
	}
	return true
}

// Find returns the first key-value pair (in sorted order) that satisfies the
// predicate. The predicate is evaluated after the lock is released, against a
// point-in-time snapshot taken under the lock, so it may safely call back into
// the collection.
func (ch *ConcurrentTreeRW[K, V]) Find(fn func(key K, value V) bool) (K, V, bool) {
	ch.lock.RLock()
	items := ch.tree.Items()
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

// FindKey returns the first key (in sorted order) that satisfies the predicate.
// The predicate is evaluated after the lock is released, against a
// point-in-time snapshot taken under the lock, so it may safely call back into
// the collection.
func (ch *ConcurrentTreeRW[K, V]) FindKey(fn func(key K) bool) (K, bool) {
	ch.lock.RLock()
	items := ch.tree.Items()
	ch.lock.RUnlock()

	for _, item := range items {
		if fn(item.Key) {
			return item.Key, true
		}
	}
	var zeroK K
	return zeroK, false
}

// FindValue returns the first value (in key-sorted order) that satisfies the
// predicate. The predicate is evaluated after the lock is released, against a
// point-in-time snapshot taken under the lock, so it may safely call back into
// the collection.
func (ch *ConcurrentTreeRW[K, V]) FindValue(fn func(value V) bool) (V, bool) {
	ch.lock.RLock()
	items := ch.tree.Items()
	ch.lock.RUnlock()

	for _, item := range items {
		if fn(item.Value) {
			return item.Value, true
		}
	}
	var zeroV V
	return zeroV, false
}

// ContainsValue checks if the given value exists in the dictionary.
func (ch *ConcurrentTreeRW[K, V]) ContainsValue(value V) bool {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return ch.tree.ContainsValue(value)
}

// Keys returns a slice containing all keys in sorted order.
func (ch *ConcurrentTreeRW[K, V]) Keys() []K {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return ch.tree.Keys()
}

// Values returns a slice containing all values in key-sorted order.
func (ch *ConcurrentTreeRW[K, V]) Values() []V {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return ch.tree.Values()
}

// Items returns a slice containing all key-value pairs as Pair structs in sorted order.
func (ch *ConcurrentTreeRW[K, V]) Items() []Pair[K, V] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return ch.tree.Items()
}

// AsMap returns the dictionary as a native Go map.
func (ch *ConcurrentTreeRW[K, V]) AsMap() map[K]V {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return ch.tree.AsMap()
}

// Put creates a new dictionary with the given key-value pair added or updated.
// Returns a new thread-safe ConcurrentTreeRW without modifying the original.
func (ch *ConcurrentTreeRW[K, V]) Put(key K, value V) Dict[K, V] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return wrapConcurrentTreeRW(ch.tree.Put(key, value).(*Tree[K, V]))
}

// PutMany creates a new dictionary with all given key-value pairs added or updated.
// Returns a new thread-safe ConcurrentTreeRW without modifying the original.
func (ch *ConcurrentTreeRW[K, V]) PutMany(pairs ...Pair[K, V]) Dict[K, V] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return wrapConcurrentTreeRW(ch.tree.PutMany(pairs...).(*Tree[K, V]))
}

// PutInPlace adds or updates the given key-value pair in the dictionary.
func (ch *ConcurrentTreeRW[K, V]) PutInPlace(key K, value V) {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	ch.tree.PutInPlace(key, value)
}

// PutManyInPlace adds or updates all given key-value pairs in the dictionary.
func (ch *ConcurrentTreeRW[K, V]) PutManyInPlace(pairs ...Pair[K, V]) {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	ch.tree.PutManyInPlace(pairs...)
}

// Remove creates a new dictionary with the given key removed.
// Returns a new thread-safe ConcurrentTreeRW without modifying the original.
func (ch *ConcurrentTreeRW[K, V]) Remove(key K) Dict[K, V] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return wrapConcurrentTreeRW(ch.tree.Remove(key).(*Tree[K, V]))
}

// RemoveMany creates a new dictionary with all given keys removed.
// Returns a new thread-safe ConcurrentTreeRW without modifying the original.
func (ch *ConcurrentTreeRW[K, V]) RemoveMany(keys ...K) Dict[K, V] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return wrapConcurrentTreeRW(ch.tree.RemoveMany(keys...).(*Tree[K, V]))
}

// RemoveInPlace removes the given key from the dictionary.
// Returns the removed value and true if the key existed; zero value and false otherwise.
func (ch *ConcurrentTreeRW[K, V]) RemoveInPlace(key K) (V, bool) {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	return ch.tree.RemoveInPlace(key)
}

// RemoveManyInPlace removes all given keys from the dictionary.
func (ch *ConcurrentTreeRW[K, V]) RemoveManyInPlace(keys ...K) {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	ch.tree.RemoveManyInPlace(keys...)
}

// Clear removes all key-value pairs from the dictionary.
func (ch *ConcurrentTreeRW[K, V]) Clear() {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	ch.tree.Clear()
}

// Min returns the entry with the smallest key.
func (ch *ConcurrentTreeRW[K, V]) Min() (K, V, bool) {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return ch.tree.Min()
}

// Max returns the entry with the largest key.
func (ch *ConcurrentTreeRW[K, V]) Max() (K, V, bool) {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return ch.tree.Max()
}

// Floor returns the entry with the largest key less than or equal to the given key.
func (ch *ConcurrentTreeRW[K, V]) Floor(key K) (K, V, bool) {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return ch.tree.Floor(key)
}

// Ceiling returns the entry with the smallest key greater than or equal to the given key.
func (ch *ConcurrentTreeRW[K, V]) Ceiling(key K) (K, V, bool) {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return ch.tree.Ceiling(key)
}

// Range returns all entries whose key is within the inclusive range [lo, hi],
// in ascending key order. Returns a non-nil (possibly empty) slice.
func (ch *ConcurrentTreeRW[K, V]) Range(lo, hi K) []Pair[K, V] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return ch.tree.Range(lo, hi)
}

// All returns an iterator over all entries in ascending key order. The entries
// are snapshotted under the read lock, so iteration is safe against concurrent
// mutation and never holds the lock while calling the consumer.
func (ch *ConcurrentTreeRW[K, V]) All() iter.Seq2[K, V] {
	ch.lock.RLock()
	items := ch.tree.Items()
	ch.lock.RUnlock()
	return seq2FromPairs(items)
}

// Backward returns an iterator over all entries in descending key order,
// snapshotted under the read lock.
func (ch *ConcurrentTreeRW[K, V]) Backward() iter.Seq2[K, V] {
	ch.lock.RLock()
	items := ch.tree.Items()
	ch.lock.RUnlock()
	return seq2FromPairsReverse(items)
}

// RangeAll returns an iterator over the entries whose key is within the inclusive
// range [lo, hi], in ascending key order, snapshotted under the read lock.
func (ch *ConcurrentTreeRW[K, V]) RangeAll(lo, hi K) iter.Seq2[K, V] {
	ch.lock.RLock()
	items := ch.tree.Range(lo, hi)
	ch.lock.RUnlock()
	return seq2FromPairs(items)
}
