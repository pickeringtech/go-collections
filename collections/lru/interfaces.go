package lru

import "iter"

// Indexable provides recency-neutral reads of a cache. None of these methods
// promote an entry to most-recently-used, so they are safe to call when you
// only want to inspect the cache without disturbing its eviction order.
type Indexable[K comparable, V any] interface {
	// Peek returns the value stored for key without marking it as recently
	// used. It returns (value, true) when the key is present, or (zero, false)
	// otherwise. Use Get (on MutableCache) when the lookup should count as a use.
	Peek(key K) (V, bool)

	// Contains reports whether key is present, without affecting recency.
	Contains(key K) bool

	// Length returns the number of entries currently held.
	Length() int

	// IsEmpty reports whether the cache holds no entries.
	IsEmpty() bool

	// Capacity returns the maximum number of entries the cache retains before
	// evicting the least-recently-used entry. It is always at least 1.
	Capacity() int
}

// Iterable provides iteration over a cache, always from most-recently-used to
// least-recently-used.
type Iterable[K comparable, V any] interface {
	// ForEach calls fn for each entry, ordered most- to least-recently-used.
	// Iteration is recency-neutral: it does not re-order the entries it visits.
	ForEach(fn EachFunc[K, V])

	// All returns a range-over-func iterator over the entries, ordered most- to
	// least-recently-used. Iteration is recency-neutral.
	All() iter.Seq2[K, V]
}

// Convertible exports a cache's contents into native Go containers, ordered
// most- to least-recently-used where order is observable.
type Convertible[K comparable, V any] interface {
	// Keys returns the keys, ordered most- to least-recently-used.
	Keys() []K

	// Values returns the values, ordered most- to least-recently-used.
	Values() []V

	// Items returns the entries as Pairs, ordered most- to least-recently-used.
	Items() []Pair[K, V]

	// AsMap returns the entries as a native Go map. The map is unordered, like
	// every Go map.
	AsMap() map[K]V
}

// Insertable provides immutable insertion: each method returns a new cache and
// leaves the receiver untouched.
type Insertable[K comparable, V any] interface {
	// Put returns a new cache with key set to value and promoted to
	// most-recently-used. If adding the entry exceeds the capacity, the new
	// cache's least-recently-used entry is evicted (and any eviction callback
	// fires for it). The receiver is not modified.
	Put(key K, value V) Cache[K, V]
}

// MutableInsertable provides in-place insertion.
type MutableInsertable[K comparable, V any] interface {
	// PutInPlace sets key to value and promotes it to most-recently-used,
	// modifying the receiver. If this exceeds the capacity, the
	// least-recently-used entry is evicted and any eviction callback fires.
	PutInPlace(key K, value V)
}

// Accessible provides the recency-marking read. It lives on MutableCache rather
// than the immutable base because a successful Get promotes the entry to
// most-recently-used, which mutates the cache's eviction order — something an
// immutable value cannot do. For a read that does not affect recency, use Peek.
type Accessible[K comparable, V any] interface {
	// Get returns the value stored for key and, on a hit, promotes that entry
	// to most-recently-used. It returns (value, true) when present, or
	// (zero, false) otherwise.
	Get(key K) (V, bool)
}

// Removable provides immutable removal: it returns a new cache and leaves the
// receiver untouched.
type Removable[K comparable, V any] interface {
	// Remove returns a new cache with key absent. The remaining entries keep
	// their relative recency order. The receiver is not modified. Explicit
	// removal never triggers the eviction callback.
	Remove(key K) Cache[K, V]
}

// MutableRemovable provides in-place removal.
type MutableRemovable[K comparable, V any] interface {
	// RemoveInPlace removes key from the receiver, returning the removed value
	// and true if it was present, or (zero, false) otherwise. Explicit removal
	// never triggers the eviction callback.
	RemoveInPlace(key K) (V, bool)

	// Clear removes every entry from the receiver. It does not trigger the
	// eviction callback.
	Clear()
}

// Cache is the immutable LRU cache interface: a bounded key-value store whose
// transforming operations (Put, Remove) return a new cache and leave the
// receiver unchanged. Reads exposed here are recency-neutral — promoting an
// entry is a mutation and lives on MutableCache.Get.
type Cache[K comparable, V any] interface {
	Indexable[K, V]
	Iterable[K, V]
	Convertible[K, V]
	Insertable[K, V]
	Removable[K, V]
}

// MutableCache is the mutable LRU cache interface. It embeds the immutable
// Cache and adds the in-place capabilities plus Get, the recency-marking read.
// This is the interface most callers want: an LRU cache is an inherently
// stateful, mutating structure.
type MutableCache[K comparable, V any] interface {
	Cache[K, V]
	Accessible[K, V]
	MutableInsertable[K, V]
	MutableRemovable[K, V]
}
