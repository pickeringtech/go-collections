package lru

import (
	"iter"
	"sync"
)

// ConcurrentLRU is a thread-safe LRU cache guarded by a single mutex. Every
// operation takes the lock, so it is the right choice when reads and writes are
// roughly balanced. For read-heavy workloads prefer ConcurrentLRURW.
//
// It wraps a plain LRU and adds synchronisation only; the eviction semantics,
// O(1) costs and recency ordering are exactly those of LRU.
type ConcurrentLRU[K comparable, V any] struct {
	inner *LRU[K, V]
	lock  *sync.Mutex
}

// Interface guards: a pointer to ConcurrentLRU satisfies both cache contracts.
var _ Cache[string, int] = &ConcurrentLRU[string, int]{}
var _ MutableCache[string, int] = &ConcurrentLRU[string, int]{}

// NewConcurrentLRU creates an empty thread-safe LRU cache bounded to capacity
// entries. A capacity below 1 is treated as 1. It accepts the same Options as
// NewLRU.
func NewConcurrentLRU[K comparable, V any](capacity int, opts ...Option[K, V]) *ConcurrentLRU[K, V] {
	return &ConcurrentLRU[K, V]{
		inner: NewLRU[K, V](capacity, opts...),
		lock:  &sync.Mutex{},
	}
}

// Peek returns the value stored for key without marking it as recently used.
func (c *ConcurrentLRU[K, V]) Peek(key K) (V, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.Peek(key)
}

// Contains reports whether key is present, without affecting recency.
func (c *ConcurrentLRU[K, V]) Contains(key K) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.Contains(key)
}

// Length returns the number of entries currently held.
func (c *ConcurrentLRU[K, V]) Length() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.Length()
}

// IsEmpty reports whether the cache holds no entries.
func (c *ConcurrentLRU[K, V]) IsEmpty() bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.IsEmpty()
}

// Capacity returns the maximum number of entries retained before eviction.
func (c *ConcurrentLRU[K, V]) Capacity() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.Capacity()
}

// Get returns the value for key and, on a hit, promotes it to most-recently-used.
func (c *ConcurrentLRU[K, V]) Get(key K) (V, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.Get(key)
}

// ForEach calls fn for each entry, most- to least-recently-used. The lock is
// held for the whole walk, so fn must not call back into this cache.
func (c *ConcurrentLRU[K, V]) ForEach(fn EachFunc[K, V]) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.inner.ForEach(fn)
}

// All returns a range-over-func iterator over a snapshot of the entries, most-
// to least-recently-used. The snapshot is taken under the lock, so iteration
// does not hold the lock and is safe to interleave with other operations.
func (c *ConcurrentLRU[K, V]) All() iter.Seq2[K, V] {
	items := c.Items()
	return func(yield func(K, V) bool) {
		for _, entry := range items {
			if !yield(entry.Key, entry.Value) {
				return
			}
		}
	}
}

// Keys returns the keys, most- to least-recently-used.
func (c *ConcurrentLRU[K, V]) Keys() []K {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.Keys()
}

// Values returns the values, most- to least-recently-used.
func (c *ConcurrentLRU[K, V]) Values() []V {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.Values()
}

// Items returns the entries as Pairs, most- to least-recently-used.
func (c *ConcurrentLRU[K, V]) Items() []Pair[K, V] {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.Items()
}

// AsMap returns the entries as a native, unordered Go map.
func (c *ConcurrentLRU[K, V]) AsMap() map[K]V {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.AsMap()
}

// Put returns a new thread-safe cache with key set to value and promoted to
// most-recently-used, evicting the least-recently-used entry if the capacity is
// exceeded. The receiver is not modified.
func (c *ConcurrentLRU[K, V]) Put(key K, value V) Cache[K, V] {
	c.lock.Lock()
	defer c.lock.Unlock()
	dup := c.inner.clone()
	dup.putInPlace(key, value)
	return &ConcurrentLRU[K, V]{inner: dup, lock: &sync.Mutex{}}
}

// PutInPlace sets key to value and promotes it to most-recently-used, evicting
// the least-recently-used entry if the capacity is exceeded.
func (c *ConcurrentLRU[K, V]) PutInPlace(key K, value V) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.inner.PutInPlace(key, value)
}

// Remove returns a new thread-safe cache with key absent; the receiver is not
// modified.
func (c *ConcurrentLRU[K, V]) Remove(key K) Cache[K, V] {
	c.lock.Lock()
	defer c.lock.Unlock()
	dup := c.inner.clone()
	dup.removeInPlace(key)
	return &ConcurrentLRU[K, V]{inner: dup, lock: &sync.Mutex{}}
}

// RemoveInPlace removes key, returning the removed value and true if it was
// present, or (zero, false) otherwise.
func (c *ConcurrentLRU[K, V]) RemoveInPlace(key K) (V, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.RemoveInPlace(key)
}

// Clear removes every entry. The capacity and eviction callback are unchanged.
func (c *ConcurrentLRU[K, V]) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.inner.Clear()
}
