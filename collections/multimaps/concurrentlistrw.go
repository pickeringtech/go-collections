package multimaps

import (
	"iter"
	"sync"

	"github.com/pickeringtech/go-collections/collections/internal/nocopy"
)

// ConcurrentRWListMultimap is a thread-safe, list-backed multimap guarded by a
// read-write mutex. Reads use a read lock and writes use a write lock, making it
// a good fit for read-heavy workloads. For balanced read/write workloads
// ConcurrentListMultimap is simpler.
//
// Immutable operations return a new ConcurrentRWListMultimap, preserving thread
// safety of results (see the concurrency standards).
//
// ConcurrentRWListMultimap must not be copied after first use; copying after construction
// produces an independent lock over shared backing data, which breaks the
// thread-safety contract. go vet reports any such copy.
type ConcurrentRWListMultimap[K comparable, V any] struct {
	_    nocopy.NoCopy
	data ListMultimap[K, V]
	lock sync.RWMutex
}

// NewConcurrentRWListMultimap creates a new thread-safe, list-backed multimap
// optimised for concurrent reads, seeded with the given entries.
func NewConcurrentRWListMultimap[K comparable, V any](entries ...Entry[K, V]) *ConcurrentRWListMultimap[K, V] {
	return &ConcurrentRWListMultimap[K, V]{
		data: NewListMultimap(entries...),
	}
}

// Interface guards to ensure ConcurrentRWListMultimap implements the required interfaces.
var _ Multimap[string, int] = &ConcurrentRWListMultimap[string, int]{}
var _ MutableMultimap[string, int] = &ConcurrentRWListMultimap[string, int]{}

func (c *ConcurrentRWListMultimap[K, V]) wrap(data ListMultimap[K, V]) *ConcurrentRWListMultimap[K, V] {
	return &ConcurrentRWListMultimap[K, V]{data: data}
}

// Get returns a copy of all values bound to the given key, in insertion order.
func (c *ConcurrentRWListMultimap[K, V]) Get(key K) []V {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.data.Get(key)
}

// ContainsKey reports whether the given key has at least one value bound to it.
func (c *ConcurrentRWListMultimap[K, V]) ContainsKey(key K) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.data.ContainsKey(key)
}

// ContainsEntry reports whether the given key is bound to the given value.
func (c *ConcurrentRWListMultimap[K, V]) ContainsEntry(key K, value V) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.data.ContainsEntry(key, value)
}

// Length returns the total number of entries (key-value associations).
func (c *ConcurrentRWListMultimap[K, V]) Length() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.data.Length()
}

// KeyCount returns the number of distinct keys.
func (c *ConcurrentRWListMultimap[K, V]) KeyCount() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.data.KeyCount()
}

// IsEmpty returns true if the multimap contains no entries.
func (c *ConcurrentRWListMultimap[K, V]) IsEmpty() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.data.IsEmpty()
}

// ForEach executes the given function once for every entry. fn is invoked after
// the lock is released, against a point-in-time snapshot taken under the lock, so
// fn may safely call back into the multimap.
func (c *ConcurrentRWListMultimap[K, V]) ForEach(fn func(key K, value V)) {
	c.lock.RLock()
	entries := c.data.Entries()
	c.lock.RUnlock()

	for _, entry := range entries {
		fn(entry.Key, entry.Value)
	}
}

// ForEachKey executes the given function once per distinct key. fn is invoked
// after the lock is released, against a point-in-time snapshot taken under the
// lock, so fn may safely call back into the multimap.
func (c *ConcurrentRWListMultimap[K, V]) ForEachKey(fn func(key K, values []V)) {
	c.lock.RLock()
	type keyValues struct {
		key    K
		values []V
	}
	snapshot := make([]keyValues, 0, len(c.data))
	for key, values := range c.data {
		snapshot = append(snapshot, keyValues{key: key, values: cloneValues(values)})
	}
	c.lock.RUnlock()

	for _, item := range snapshot {
		fn(item.key, item.values)
	}
}

// All returns an iterator over every entry. The entries are snapshotted under
// the read lock, so iteration is safe against concurrent mutation and never
// holds the lock while calling yield (yield may safely call back into the
// multimap).
func (c *ConcurrentRWListMultimap[K, V]) All() iter.Seq2[K, V] {
	c.lock.RLock()
	entries := c.data.Entries()
	c.lock.RUnlock()

	return func(yield func(K, V) bool) {
		for _, entry := range entries {
			if !yield(entry.Key, entry.Value) {
				return
			}
		}
	}
}

// KeysSeq returns an iterator over the distinct keys. The keys are snapshotted
// under the read lock, so iteration is safe against concurrent mutation and
// never holds the lock while calling yield (yield may safely call back into the
// multimap).
func (c *ConcurrentRWListMultimap[K, V]) KeysSeq() iter.Seq[K] {
	c.lock.RLock()
	keys := c.data.Keys()
	c.lock.RUnlock()

	return func(yield func(K) bool) {
		for _, key := range keys {
			if !yield(key) {
				return
			}
		}
	}
}

// Filter returns a new thread-safe multimap containing only entries that satisfy
// the predicate. The predicate is evaluated after the lock is released, against a
// point-in-time snapshot taken under the lock, so it may safely call back into
// the multimap.
func (c *ConcurrentRWListMultimap[K, V]) Filter(fn func(key K, value V) bool) Multimap[K, V] {
	c.lock.RLock()
	entries := c.data.Entries()
	c.lock.RUnlock()

	result := make(ListMultimap[K, V])
	for _, entry := range entries {
		if fn(entry.Key, entry.Value) {
			result[entry.Key] = append(result[entry.Key], entry.Value)
		}
	}
	return c.wrap(result)
}

// FilterInPlace removes every entry that does not satisfy the predicate. The
// predicate is evaluated after the lock is released, against a point-in-time
// snapshot taken under the lock, so it may safely call back into the multimap.
//
// Each rejected (key, value) entry removes one matching entry from the multimap
// as it stands at apply time, so entries added concurrently in the evaluation
// window are preserved.
func (c *ConcurrentRWListMultimap[K, V]) FilterInPlace(fn func(key K, value V) bool) {
	c.lock.RLock()
	entries := c.data.Entries()
	c.lock.RUnlock()

	var toRemove []Entry[K, V]
	for _, entry := range entries {
		if !fn(entry.Key, entry.Value) {
			toRemove = append(toRemove, entry)
		}
	}

	c.lock.Lock()
	for _, entry := range toRemove {
		c.data.RemoveInPlace(entry.Key, entry.Value)
	}
	c.lock.Unlock()
}

// AllMatch returns true if every entry satisfies the predicate. The predicate is
// evaluated after the lock is released, against a point-in-time snapshot taken
// under the lock, so it may safely call back into the multimap.
func (c *ConcurrentRWListMultimap[K, V]) AllMatch(fn func(key K, value V) bool) bool {
	c.lock.RLock()
	entries := c.data.Entries()
	c.lock.RUnlock()

	for _, entry := range entries {
		if !fn(entry.Key, entry.Value) {
			return false
		}
	}
	return true
}

// AnyMatch returns true if at least one entry satisfies the predicate. The
// predicate is evaluated after the lock is released, against a point-in-time
// snapshot taken under the lock, so it may safely call back into the multimap.
func (c *ConcurrentRWListMultimap[K, V]) AnyMatch(fn func(key K, value V) bool) bool {
	c.lock.RLock()
	entries := c.data.Entries()
	c.lock.RUnlock()

	for _, entry := range entries {
		if fn(entry.Key, entry.Value) {
			return true
		}
	}
	return false
}

// NoneMatch returns true if no entry satisfies the predicate. The predicate is
// evaluated after the lock is released, against a point-in-time snapshot taken
// under the lock, so it may safely call back into the multimap.
func (c *ConcurrentRWListMultimap[K, V]) NoneMatch(fn func(key K, value V) bool) bool {
	return !c.AnyMatch(fn)
}

// Find returns the first entry that satisfies the predicate. The predicate is
// evaluated after the lock is released, against a point-in-time snapshot taken
// under the lock, so it may safely call back into the multimap.
func (c *ConcurrentRWListMultimap[K, V]) Find(fn func(key K, value V) bool) (K, V, bool) {
	c.lock.RLock()
	entries := c.data.Entries()
	c.lock.RUnlock()

	for _, entry := range entries {
		if fn(entry.Key, entry.Value) {
			return entry.Key, entry.Value, true
		}
	}
	var zeroKey K
	var zeroValue V
	return zeroKey, zeroValue, false
}

// Keys returns a slice containing each distinct key once.
func (c *ConcurrentRWListMultimap[K, V]) Keys() []K {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.data.Keys()
}

// Values returns a slice containing every value across all keys.
func (c *ConcurrentRWListMultimap[K, V]) Values() []V {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.data.Values()
}

// Entries returns a slice containing every entry.
func (c *ConcurrentRWListMultimap[K, V]) Entries() []Entry[K, V] {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.data.Entries()
}

// AsMap returns the multimap as a native Go map from each key to a copy of its values.
func (c *ConcurrentRWListMultimap[K, V]) AsMap() map[K][]V {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.data.AsMap()
}

// Put returns a new thread-safe multimap with the given value bound to the key.
func (c *ConcurrentRWListMultimap[K, V]) Put(key K, value V) Multimap[K, V] {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.wrap(c.data.Put(key, value).(ListMultimap[K, V]))
}

// PutAll returns a new thread-safe multimap with all the given values bound to the key.
func (c *ConcurrentRWListMultimap[K, V]) PutAll(key K, values ...V) Multimap[K, V] {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.wrap(c.data.PutAll(key, values...).(ListMultimap[K, V]))
}

// PutInPlace binds the given value to the given key, modifying the multimap in place.
func (c *ConcurrentRWListMultimap[K, V]) PutInPlace(key K, value V) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data.PutInPlace(key, value)
}

// PutAllInPlace binds all the given values to the given key, modifying the multimap in place.
func (c *ConcurrentRWListMultimap[K, V]) PutAllInPlace(key K, values ...V) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data.PutAllInPlace(key, values...)
}

// Remove returns a new thread-safe multimap with a single binding removed.
func (c *ConcurrentRWListMultimap[K, V]) Remove(key K, value V) Multimap[K, V] {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.wrap(c.data.Remove(key, value).(ListMultimap[K, V]))
}

// RemoveAll returns a new thread-safe multimap with the given key removed.
func (c *ConcurrentRWListMultimap[K, V]) RemoveAll(key K) Multimap[K, V] {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.wrap(c.data.RemoveAll(key).(ListMultimap[K, V]))
}

// RemoveInPlace removes a single binding, modifying the multimap in place.
func (c *ConcurrentRWListMultimap[K, V]) RemoveInPlace(key K, value V) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.data.RemoveInPlace(key, value)
}

// RemoveAllInPlace removes the given key and all its values, modifying the multimap in place.
func (c *ConcurrentRWListMultimap[K, V]) RemoveAllInPlace(key K) ([]V, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.data.RemoveAllInPlace(key)
}

// Clear removes all entries from the multimap.
func (c *ConcurrentRWListMultimap[K, V]) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data.Clear()
}
