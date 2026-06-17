package multimaps

import (
	"iter"
	"sync"
)

// ConcurrentSetMultimap is a thread-safe, set-backed multimap guarded by a
// single mutex. Every operation acquires the lock, making it a good fit for
// workloads with a balanced mix of reads and writes. For read-heavy workloads
// prefer ConcurrentRWSetMultimap.
//
// Immutable operations return a new ConcurrentSetMultimap, preserving thread
// safety of results (see the concurrency standards).
type ConcurrentSetMultimap[K comparable, V comparable] struct {
	data SetMultimap[K, V]
	lock *sync.Mutex
}

// NewConcurrentSetMultimap creates a new thread-safe, set-backed multimap seeded
// with the given entries.
func NewConcurrentSetMultimap[K comparable, V comparable](entries ...Entry[K, V]) *ConcurrentSetMultimap[K, V] {
	return &ConcurrentSetMultimap[K, V]{
		data: NewSetMultimap(entries...),
		lock: &sync.Mutex{},
	}
}

// Interface guards to ensure ConcurrentSetMultimap implements the required interfaces.
var _ Multimap[string, int] = &ConcurrentSetMultimap[string, int]{}
var _ MutableMultimap[string, int] = &ConcurrentSetMultimap[string, int]{}

func (c *ConcurrentSetMultimap[K, V]) wrap(data SetMultimap[K, V]) *ConcurrentSetMultimap[K, V] {
	return &ConcurrentSetMultimap[K, V]{data: data, lock: &sync.Mutex{}}
}

// Get returns a copy of all values bound to the given key.
func (c *ConcurrentSetMultimap[K, V]) Get(key K) []V {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.data.Get(key)
}

// ContainsKey reports whether the given key has at least one value bound to it.
func (c *ConcurrentSetMultimap[K, V]) ContainsKey(key K) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.data.ContainsKey(key)
}

// ContainsEntry reports whether the given key is bound to the given value.
func (c *ConcurrentSetMultimap[K, V]) ContainsEntry(key K, value V) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.data.ContainsEntry(key, value)
}

// Length returns the total number of entries (distinct key-value associations).
func (c *ConcurrentSetMultimap[K, V]) Length() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.data.Length()
}

// KeyCount returns the number of distinct keys.
func (c *ConcurrentSetMultimap[K, V]) KeyCount() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.data.KeyCount()
}

// IsEmpty returns true if the multimap contains no entries.
func (c *ConcurrentSetMultimap[K, V]) IsEmpty() bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.data.IsEmpty()
}

// ForEach executes the given function once for every entry. fn is invoked after
// the lock is released, against a point-in-time snapshot taken under the lock, so
// fn may safely call back into the multimap.
func (c *ConcurrentSetMultimap[K, V]) ForEach(fn func(key K, value V)) {
	c.lock.Lock()
	entries := c.data.Entries()
	c.lock.Unlock()

	for _, entry := range entries {
		fn(entry.Key, entry.Value)
	}
}

// ForEachKey executes the given function once per distinct key. fn is invoked
// after the lock is released, against a point-in-time snapshot taken under the
// lock, so fn may safely call back into the multimap.
func (c *ConcurrentSetMultimap[K, V]) ForEachKey(fn func(key K, values []V)) {
	c.lock.Lock()
	type keyValues struct {
		key    K
		values []V
	}
	snapshot := make([]keyValues, 0, len(c.data))
	for key, values := range c.data {
		snapshot = append(snapshot, keyValues{key: key, values: valuesSlice(values)})
	}
	c.lock.Unlock()

	for _, item := range snapshot {
		fn(item.key, item.values)
	}
}

// All returns an iterator over every entry. The entries are snapshotted under
// the lock, so iteration is safe against concurrent mutation and never holds the
// lock while calling yield (yield may safely call back into the multimap).
func (c *ConcurrentSetMultimap[K, V]) All() iter.Seq2[K, V] {
	c.lock.Lock()
	entries := c.data.Entries()
	c.lock.Unlock()

	return func(yield func(K, V) bool) {
		for _, entry := range entries {
			if !yield(entry.Key, entry.Value) {
				return
			}
		}
	}
}

// KeysSeq returns an iterator over the distinct keys. The keys are snapshotted
// under the lock, so iteration is safe against concurrent mutation and never
// holds the lock while calling yield (yield may safely call back into the
// multimap).
func (c *ConcurrentSetMultimap[K, V]) KeysSeq() iter.Seq[K] {
	c.lock.Lock()
	keys := c.data.Keys()
	c.lock.Unlock()

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
func (c *ConcurrentSetMultimap[K, V]) Filter(fn func(key K, value V) bool) Multimap[K, V] {
	c.lock.Lock()
	entries := c.data.Entries()
	c.lock.Unlock()

	result := make(SetMultimap[K, V])
	for _, entry := range entries {
		if fn(entry.Key, entry.Value) {
			result.PutInPlace(entry.Key, entry.Value)
		}
	}
	return c.wrap(result)
}

// FilterInPlace removes every entry that does not satisfy the predicate. The
// predicate is evaluated after the lock is released, against a point-in-time
// snapshot taken under the lock, so it may safely call back into the multimap.
// Modifications made concurrently with evaluation are not reflected in the
// retained set.
func (c *ConcurrentSetMultimap[K, V]) FilterInPlace(fn func(key K, value V) bool) {
	c.lock.Lock()
	entries := c.data.Entries()
	c.lock.Unlock()

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
func (c *ConcurrentSetMultimap[K, V]) AllMatch(fn func(key K, value V) bool) bool {
	c.lock.Lock()
	entries := c.data.Entries()
	c.lock.Unlock()

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
func (c *ConcurrentSetMultimap[K, V]) AnyMatch(fn func(key K, value V) bool) bool {
	c.lock.Lock()
	entries := c.data.Entries()
	c.lock.Unlock()

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
func (c *ConcurrentSetMultimap[K, V]) NoneMatch(fn func(key K, value V) bool) bool {
	return !c.AnyMatch(fn)
}

// Find returns the first entry that satisfies the predicate. The predicate is
// evaluated after the lock is released, against a point-in-time snapshot taken
// under the lock, so it may safely call back into the multimap.
func (c *ConcurrentSetMultimap[K, V]) Find(fn func(key K, value V) bool) (K, V, bool) {
	c.lock.Lock()
	entries := c.data.Entries()
	c.lock.Unlock()

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
func (c *ConcurrentSetMultimap[K, V]) Keys() []K {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.data.Keys()
}

// Values returns a slice containing every value across all keys.
func (c *ConcurrentSetMultimap[K, V]) Values() []V {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.data.Values()
}

// Entries returns a slice containing every entry.
func (c *ConcurrentSetMultimap[K, V]) Entries() []Entry[K, V] {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.data.Entries()
}

// AsMap returns the multimap as a native Go map from each key to a copy of its values.
func (c *ConcurrentSetMultimap[K, V]) AsMap() map[K][]V {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.data.AsMap()
}

// Put returns a new thread-safe multimap with the given value bound to the key.
func (c *ConcurrentSetMultimap[K, V]) Put(key K, value V) Multimap[K, V] {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.wrap(c.data.Put(key, value).(SetMultimap[K, V]))
}

// PutAll returns a new thread-safe multimap with all the given values bound to the key.
func (c *ConcurrentSetMultimap[K, V]) PutAll(key K, values ...V) Multimap[K, V] {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.wrap(c.data.PutAll(key, values...).(SetMultimap[K, V]))
}

// PutInPlace binds the given value to the given key, modifying the multimap in place.
func (c *ConcurrentSetMultimap[K, V]) PutInPlace(key K, value V) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data.PutInPlace(key, value)
}

// PutAllInPlace binds all the given values to the given key, modifying the multimap in place.
func (c *ConcurrentSetMultimap[K, V]) PutAllInPlace(key K, values ...V) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data.PutAllInPlace(key, values...)
}

// Remove returns a new thread-safe multimap with the given value unbound from the key.
func (c *ConcurrentSetMultimap[K, V]) Remove(key K, value V) Multimap[K, V] {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.wrap(c.data.Remove(key, value).(SetMultimap[K, V]))
}

// RemoveAll returns a new thread-safe multimap with the given key removed.
func (c *ConcurrentSetMultimap[K, V]) RemoveAll(key K) Multimap[K, V] {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.wrap(c.data.RemoveAll(key).(SetMultimap[K, V]))
}

// RemoveInPlace unbinds the given value from the given key, modifying the multimap in place.
func (c *ConcurrentSetMultimap[K, V]) RemoveInPlace(key K, value V) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.data.RemoveInPlace(key, value)
}

// RemoveAllInPlace removes the given key and all its values, modifying the multimap in place.
func (c *ConcurrentSetMultimap[K, V]) RemoveAllInPlace(key K) ([]V, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.data.RemoveAllInPlace(key)
}

// Clear removes all entries from the multimap.
func (c *ConcurrentSetMultimap[K, V]) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data.Clear()
}
