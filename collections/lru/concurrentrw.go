package lru

import (
	"iter"
	"sync"

	"github.com/pickeringtech/go-collections/collections/internal/nocopy"
)

// ConcurrentLRURW is a thread-safe LRU cache guarded by a read-write mutex.
// Recency-neutral reads (Peek, Contains, Length, the exporters) take a shared
// read lock and can proceed concurrently; anything that mutates state — every
// write, plus Get, which promotes the entry it returns — takes the exclusive
// write lock. Prefer it over ConcurrentLRU for read-heavy workloads.
//
// Note that Get is a write here: marking an entry most-recently-used re-orders
// the recency list. A workload dominated by Get therefore sees little benefit
// from the RW variant; reach for Peek when a lookup need not count as a use.
//
// ConcurrentLRURW must not be copied after first use; copying after construction
// produces an independent lock over shared backing data, which breaks the
// thread-safety contract. go vet reports any such copy.
type ConcurrentLRURW[K comparable, V any] struct {
	_     nocopy.NoCopy
	inner *LRU[K, V]
	lock  sync.RWMutex
}

// Interface guards: a pointer to ConcurrentLRURW satisfies both cache contracts.
var _ Cache[string, int] = &ConcurrentLRURW[string, int]{}
var _ MutableCache[string, int] = &ConcurrentLRURW[string, int]{}

// NewConcurrentLRURW creates an empty thread-safe LRU cache bounded to capacity
// entries, backed by a read-write mutex. A capacity below 1 is treated as 1. It
// accepts the same Options as NewLRU.
func NewConcurrentLRURW[K comparable, V any](capacity int, opts ...Option[K, V]) *ConcurrentLRURW[K, V] {
	return &ConcurrentLRURW[K, V]{
		inner: NewLRU[K, V](capacity, opts...),
	}
}

// Peek returns the value stored for key without marking it as recently used.
func (c *ConcurrentLRURW[K, V]) Peek(key K) (V, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.inner.Peek(key)
}

// Contains reports whether key is present, without affecting recency.
func (c *ConcurrentLRURW[K, V]) Contains(key K) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.inner.Contains(key)
}

// Length returns the number of entries currently held.
func (c *ConcurrentLRURW[K, V]) Length() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.inner.Length()
}

// IsEmpty reports whether the cache holds no entries.
func (c *ConcurrentLRURW[K, V]) IsEmpty() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.inner.IsEmpty()
}

// Capacity returns the maximum number of entries retained before eviction.
func (c *ConcurrentLRURW[K, V]) Capacity() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.inner.Capacity()
}

// Get returns the value for key and, on a hit, promotes it to most-recently-used.
// Because that promotion mutates the recency order, Get takes the write lock.
func (c *ConcurrentLRURW[K, V]) Get(key K) (V, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.Get(key)
}

// ForEach calls fn for each entry, most- to least-recently-used. fn is invoked
// after the lock is released, against a point-in-time snapshot taken under the
// read lock, so fn may safely call back into the cache. Snapshotting does not
// affect recency, exactly as Items does.
func (c *ConcurrentLRURW[K, V]) ForEach(fn EachFunc[K, V]) {
	items := c.Items()
	for _, entry := range items {
		fn(entry.Key, entry.Value)
	}
}

// All returns a range-over-func iterator over a snapshot of the entries, most-
// to least-recently-used. The snapshot is taken under the read lock, so
// iteration does not hold the lock.
func (c *ConcurrentLRURW[K, V]) All() iter.Seq2[K, V] {
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
func (c *ConcurrentLRURW[K, V]) Keys() []K {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.inner.Keys()
}

// Values returns the values, most- to least-recently-used.
func (c *ConcurrentLRURW[K, V]) Values() []V {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.inner.Values()
}

// Items returns the entries as Pairs, most- to least-recently-used.
func (c *ConcurrentLRURW[K, V]) Items() []Pair[K, V] {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.inner.Items()
}

// AsMap returns the entries as a native, unordered Go map.
func (c *ConcurrentLRURW[K, V]) AsMap() map[K]V {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.inner.AsMap()
}

// Put returns a new thread-safe cache with key set to value and promoted to
// most-recently-used, evicting the least-recently-used entry if the capacity is
// exceeded. The receiver is read-locked and not modified.
func (c *ConcurrentLRURW[K, V]) Put(key K, value V) Cache[K, V] {
	c.lock.RLock()
	defer c.lock.RUnlock()
	dup := c.inner.clone()
	dup.putInPlace(key, value)
	return &ConcurrentLRURW[K, V]{inner: dup}
}

// PutInPlace sets key to value and promotes it to most-recently-used, evicting
// the least-recently-used entry if the capacity is exceeded.
func (c *ConcurrentLRURW[K, V]) PutInPlace(key K, value V) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.inner.PutInPlace(key, value)
}

// Remove returns a new thread-safe cache with key absent; the receiver is
// read-locked and not modified.
func (c *ConcurrentLRURW[K, V]) Remove(key K) Cache[K, V] {
	c.lock.RLock()
	defer c.lock.RUnlock()
	dup := c.inner.clone()
	dup.removeInPlace(key)
	return &ConcurrentLRURW[K, V]{inner: dup}
}

// RemoveInPlace removes key, returning the removed value and true if it was
// present, or (zero, false) otherwise.
func (c *ConcurrentLRURW[K, V]) RemoveInPlace(key K) (V, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.RemoveInPlace(key)
}

// Clear removes every entry. The capacity and eviction callback are unchanged.
func (c *ConcurrentLRURW[K, V]) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.inner.Clear()
}
