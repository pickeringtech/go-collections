package lru

import "iter"

// node is one entry in the intrusive doubly-linked recency list. The list is
// what makes both reads and writes O(1): the map finds a node by key in
// constant time, and the node's prev/next pointers let it be re-linked to the
// front (most-recently-used) without scanning anything.
type node[K comparable, V any] struct {
	key   K
	value V
	prev  *node[K, V]
	next  *node[K, V]
}

// LRU is a bounded, single-threaded least-recently-used cache. It keeps at most
// Capacity entries; inserting beyond that evicts the least-recently-used one.
//
// The implementation pairs a map (key -> node) with an intrusive doubly-linked
// list ordered most- to least-recently-used, giving O(1) Get, Peek, Put and
// Remove. head and tail are sentinel nodes that never hold data, so the list is
// never empty and link/unlink code needs no nil checks at the ends:
//
//	head <-> (MRU) <-> ... <-> (LRU) <-> tail
//
// LRU is not safe for concurrent use. For that, choose ConcurrentLRU (mutex) or
// ConcurrentLRURW (read-write mutex), which wrap this type — see the package doc.
type LRU[K comparable, V any] struct {
	items    map[K]*node[K, V]
	head     *node[K, V]
	tail     *node[K, V]
	capacity int
	onEvict  EvictFunc[K, V]
}

// Interface guards: a pointer to LRU satisfies both the immutable and mutable
// cache contracts. A missing method breaks the build here, not at a call site.
var _ Cache[string, int] = &LRU[string, int]{}
var _ MutableCache[string, int] = &LRU[string, int]{}

// NewLRU creates an empty LRU cache bounded to capacity entries. A capacity
// below 1 is treated as 1, since a cache that can hold nothing is never useful.
//
// Configure optional behaviour with Option values:
//
//	cache := lru.NewLRU[string, int](2,
//		lru.WithOnEvict(func(k string, v int) { fmt.Printf("evicted %s\n", k) }),
//		lru.WithEntries(lru.Pair[string, int]{Key: "a", Value: 1}),
//	)
func NewLRU[K comparable, V any](capacity int, opts ...Option[K, V]) *LRU[K, V] {
	cfg := resolveConfig(opts)
	c := newEmptyLRU[K, V](capacity, cfg.onEvict)
	for _, entry := range cfg.entries {
		c.fireEvict(c.putInPlace(entry.Key, entry.Value))
	}
	return c
}

// newEmptyLRU builds an empty cache with its sentinels linked together.
func newEmptyLRU[K comparable, V any](capacity int, onEvict EvictFunc[K, V]) *LRU[K, V] {
	if capacity < 1 {
		capacity = 1
	}
	head := &node[K, V]{}
	tail := &node[K, V]{}
	head.next = tail
	tail.prev = head
	return &LRU[K, V]{
		items:    make(map[K]*node[K, V]),
		head:     head,
		tail:     tail,
		capacity: capacity,
		onEvict:  onEvict,
	}
}

// pushFront links n in just after the head sentinel, making it most-recently-used.
func (c *LRU[K, V]) pushFront(n *node[K, V]) {
	n.prev = c.head
	n.next = c.head.next
	c.head.next.prev = n
	c.head.next = n
}

// unlink removes n from the recency list, leaving its pointers dangling for the
// caller to discard or re-link.
func (c *LRU[K, V]) unlink(n *node[K, V]) {
	n.prev.next = n.next
	n.next.prev = n.prev
}

// moveToFront promotes an already-linked node to most-recently-used.
func (c *LRU[K, V]) moveToFront(n *node[K, V]) {
	c.unlink(n)
	c.pushFront(n)
}

// evictLRU drops the least-recently-used entry and returns it so the caller can
// fire the eviction callback once any lock is released. It is only called when
// the cache is over capacity, so tail.prev is always a real entry rather than
// the head sentinel.
func (c *LRU[K, V]) evictLRU() Pair[K, V] {
	victim := c.tail.prev
	c.unlink(victim)
	delete(c.items, victim.key)
	return Pair[K, V]{Key: victim.key, Value: victim.value}
}

// putInPlace is the unsynchronised insert shared by the constructor and the
// exported mutating methods. It returns the evicted entry, if any, without
// invoking the eviction callback — callers fire onEvict via fireEvict after any
// lock is released, so the callback may safely re-enter the cache.
func (c *LRU[K, V]) putInPlace(key K, value V) (Pair[K, V], bool) {
	existing, ok := c.items[key]
	if ok {
		existing.value = value
		c.moveToFront(existing)
		return Pair[K, V]{}, false
	}
	fresh := &node[K, V]{key: key, value: value}
	c.items[key] = fresh
	c.pushFront(fresh)
	if len(c.items) > c.capacity {
		return c.evictLRU(), true
	}
	return Pair[K, V]{}, false
}

// fireEvict invokes the eviction callback for the pair returned by putInPlace,
// but only when an eviction actually happened and a callback is registered.
// Callers invoke it after releasing any lock, so onEvict never runs while the
// cache is locked.
func (c *LRU[K, V]) fireEvict(evicted Pair[K, V], evictedOK bool) {
	if evictedOK && c.onEvict != nil {
		c.onEvict(evicted.Key, evicted.Value)
	}
}

// clone makes an independent copy with the same capacity, callback and recency
// order — the basis of the immutable Put/Remove operations.
func (c *LRU[K, V]) clone() *LRU[K, V] {
	dup := newEmptyLRU[K, V](c.capacity, c.onEvict)
	// Walk the receiver least- to most-recently-used and push each node to the
	// front, so the copy ends up in the same most- to least-recently-used order.
	for n := c.tail.prev; n != c.head; n = n.prev {
		fresh := &node[K, V]{key: n.key, value: n.value}
		dup.items[n.key] = fresh
		dup.pushFront(fresh)
	}
	return dup
}

// Peek returns the value stored for key without marking it as recently used.
func (c *LRU[K, V]) Peek(key K) (V, bool) {
	n, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}
	return n.value, true
}

// Contains reports whether key is present, without affecting recency.
func (c *LRU[K, V]) Contains(key K) bool {
	_, ok := c.items[key]
	return ok
}

// Length returns the number of entries currently held.
func (c *LRU[K, V]) Length() int {
	return len(c.items)
}

// IsEmpty reports whether the cache holds no entries.
func (c *LRU[K, V]) IsEmpty() bool {
	return len(c.items) == 0
}

// Capacity returns the maximum number of entries retained before eviction.
func (c *LRU[K, V]) Capacity() int {
	return c.capacity
}

// Get returns the value for key and, on a hit, promotes it to most-recently-used.
func (c *LRU[K, V]) Get(key K) (V, bool) {
	n, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}
	c.moveToFront(n)
	return n.value, true
}

// ForEach calls fn for each entry, most- to least-recently-used.
func (c *LRU[K, V]) ForEach(fn EachFunc[K, V]) {
	for n := c.head.next; n != c.tail; n = n.next {
		fn(n.key, n.value)
	}
}

// All returns a range-over-func iterator over the entries, most- to
// least-recently-used.
func (c *LRU[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for n := c.head.next; n != c.tail; n = n.next {
			if !yield(n.key, n.value) {
				return
			}
		}
	}
}

// Keys returns the keys, most- to least-recently-used.
func (c *LRU[K, V]) Keys() []K {
	keys := make([]K, 0, len(c.items))
	for n := c.head.next; n != c.tail; n = n.next {
		keys = append(keys, n.key)
	}
	return keys
}

// Values returns the values, most- to least-recently-used.
func (c *LRU[K, V]) Values() []V {
	values := make([]V, 0, len(c.items))
	for n := c.head.next; n != c.tail; n = n.next {
		values = append(values, n.value)
	}
	return values
}

// Items returns the entries as Pairs, most- to least-recently-used.
func (c *LRU[K, V]) Items() []Pair[K, V] {
	items := make([]Pair[K, V], 0, len(c.items))
	for n := c.head.next; n != c.tail; n = n.next {
		items = append(items, Pair[K, V]{Key: n.key, Value: n.value})
	}
	return items
}

// AsMap returns the entries as a native, unordered Go map.
func (c *LRU[K, V]) AsMap() map[K]V {
	result := make(map[K]V, len(c.items))
	for key, n := range c.items {
		result[key] = n.value
	}
	return result
}

// Put returns a new cache with key set to value and promoted to
// most-recently-used, evicting the least-recently-used entry if the capacity is
// exceeded. The receiver is not modified.
func (c *LRU[K, V]) Put(key K, value V) Cache[K, V] {
	dup := c.clone()
	dup.fireEvict(dup.putInPlace(key, value))
	return dup
}

// PutInPlace sets key to value and promotes it to most-recently-used, evicting
// the least-recently-used entry if the capacity is exceeded.
func (c *LRU[K, V]) PutInPlace(key K, value V) {
	c.fireEvict(c.putInPlace(key, value))
}

// Remove returns a new cache with key absent; the receiver is not modified.
func (c *LRU[K, V]) Remove(key K) Cache[K, V] {
	dup := c.clone()
	dup.removeInPlace(key)
	return dup
}

// RemoveInPlace removes key, returning the removed value and true if it was
// present, or (zero, false) otherwise.
func (c *LRU[K, V]) RemoveInPlace(key K) (V, bool) {
	return c.removeInPlace(key)
}

// removeInPlace is the unsynchronised removal shared by Remove and RemoveInPlace.
func (c *LRU[K, V]) removeInPlace(key K) (V, bool) {
	n, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}
	c.unlink(n)
	delete(c.items, key)
	return n.value, true
}

// Clear removes every entry. The capacity and eviction callback are unchanged.
func (c *LRU[K, V]) Clear() {
	c.items = make(map[K]*node[K, V])
	c.head.next = c.tail
	c.tail.prev = c.head
}
