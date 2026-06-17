package lru

// Pair represents a key-value entry held by a cache. It mirrors dicts.Pair so
// the families read alike — "learn one, learn all".
type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

// EachFunc operates on a single key-value entry during iteration.
type EachFunc[K comparable, V any] func(key K, value V)

// EvictFunc is invoked with the key and value of an entry the moment it leaves
// the cache because the capacity bound was exceeded. It is the optional
// eviction callback described in the package overview: use it to release the
// resources an evicted value owned (close a handle, decrement a refcount, emit
// a metric). It runs synchronously within the operation that caused the
// eviction, so keep it cheap and non-blocking. On the concurrent caches it is
// invoked after the lock is released, against the entry snapshotted under it, so
// the callback may safely call back into the same cache without deadlocking.
type EvictFunc[K comparable, V any] func(key K, value V)

// Option configures a cache at construction time. Options are applied in the
// order given to the NewXxx constructors.
type Option[K comparable, V any] func(*config[K, V])

// config holds the resolved construction-time settings shared by every cache
// variant. It is unexported; callers shape it exclusively through Option values.
type config[K comparable, V any] struct {
	onEvict EvictFunc[K, V]
	entries []Pair[K, V]
}

// WithOnEvict registers an eviction callback fired whenever an entry is dropped
// to honour the capacity bound. Passing nil leaves the cache without a callback,
// which is also the default.
//
// The callback fires for capacity-driven evictions only — not for entries you
// remove explicitly via Remove/RemoveInPlace/Clear, which are deliberate and
// already visible to the caller.
func WithOnEvict[K comparable, V any](fn EvictFunc[K, V]) Option[K, V] {
	return func(c *config[K, V]) {
		c.onEvict = fn
	}
}

// WithEntries seeds the cache with the given pairs, inserted left-to-right as
// though each were passed to PutInPlace. Later pairs are therefore more
// recently used than earlier ones, and seeding more pairs than the capacity
// evicts the earliest just as live inserts would.
func WithEntries[K comparable, V any](entries ...Pair[K, V]) Option[K, V] {
	return func(c *config[K, V]) {
		c.entries = append(c.entries, entries...)
	}
}

// resolveConfig folds the given options into a config.
func resolveConfig[K comparable, V any](opts []Option[K, V]) config[K, V] {
	var cfg config[K, V]
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}
