// Package lru provides a generic, bounded least-recently-used (LRU) cache: a
// fixed-capacity key-value store that evicts the least-recently-used entry when
// it would otherwise overflow. It is the cache Go's standard library never ships
// — type-safe, O(1), and offered in the same immutable/mutable and thread-safe
// shapes as the rest of go-collections.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/collections/lru"
//
//	cache := lru.NewLRU[string, int](2)
//	cache.PutInPlace("a", 1)
//	cache.PutInPlace("b", 2)
//	cache.PutInPlace("c", 3) // capacity exceeded → evicts "a"
//
//	value, found := cache.Get("b")
//	// Result: value == 2, found == true
//	_, stillThere := cache.Peek("a")
//	// Result: stillThere == false (evicted)
//
// # Why Use an LRU Cache?
//
// A plain map grows without bound and forgets nothing, so using one as a cache
// leaks memory. Hand-rolling the eviction logic is fiddly — you need a map for
// O(1) lookup and a linked list for O(1) recency updates, kept in lockstep:
//
//	// Native approach — a map that never forgets, so it grows forever:
//	cache := map[string][]byte{}
//	cache[key] = fetch(key) // nothing is ever evicted
//
//	// With this package — bounded, with O(1) eviction handled for you:
//	cache := lru.NewLRU[string, []byte](1000)
//	cache.PutInPlace(key, fetch(key)) // the coldest entry is dropped at capacity
//
// # Recency: Get vs Peek
//
// Get counts as a use: a hit promotes the entry to most-recently-used, so it
// will be evicted last. Peek reads the same value without touching recency, for
// when an inspection should not keep an entry alive. Peek, Contains, Length,
// IsEmpty, Capacity and the exporters (Keys, Values, Items, AsMap, ForEach, All)
// are all recency-neutral.
//
// Because a successful Get mutates the recency order, it lives on MutableCache,
// not the immutable Cache interface — an immutable value cannot re-order itself.
//
// # Eviction Callback
//
// Register a callback to learn when capacity forces an entry out — to close a
// handle, decrement a refcount, or emit a metric:
//
//	cache := lru.NewLRU[string, *os.File](100,
//		lru.WithOnEvict(func(key string, f *os.File) { f.Close() }),
//	)
//
// The callback fires only for capacity-driven evictions, never for entries you
// remove explicitly via Remove, RemoveInPlace or Clear.
//
// # Available Implementations
//
// LRU[K, V]:
//   - Single-threaded, O(1) Get/Peek/Put/Remove
//   - A map plus an intrusive doubly-linked recency list
//
// ConcurrentLRU[K, V]:
//   - Thread-safe via a single mutex; every operation locks
//   - Best when reads and writes are balanced
//
// ConcurrentLRURW[K, V]:
//   - Thread-safe via a read-write mutex; recency-neutral reads share a read lock
//   - Best for read-heavy workloads — but note Get takes the write lock, since
//     it re-orders recency; reach for Peek when a lookup need not count as a use
//
// All three implement the same Cache and MutableCache interfaces, so you can
// swap between them without touching call sites.
//
// # Immutable vs Mutable Operations
//
// Mutable operations modify the receiver and are the common case for a cache:
//
//	cache.PutInPlace("key", value)         // insert/update in place
//	value, found := cache.Get("key")       // read + promote in place
//	removed, ok := cache.RemoveInPlace("key")
//	cache.Clear()
//
// Immutable operations leave the receiver untouched and return a new cache:
//
//	updated := cache.Put("key", value)     // returns a new Cache
//	smaller := cache.Remove("key")         // returns a new Cache
//
// On the concurrent types the immutable operations return a new value of the
// same concurrent type, so thread-safety is never silently dropped.
//
// # Thread Safety
//
// The plain LRU is not safe for concurrent use. For shared access choose a
// concurrent variant:
//
//	// Balanced read/write workloads
//	sessions := lru.NewConcurrentLRU[string, Session](10_000)
//
//	// Read-heavy workloads
//	pages := lru.NewConcurrentLRURW[string, []byte](1_000)
//
// Callbacks passed to ForEach and the iterator methods (All, Keys, Values) run
// after the lock is released, against a point-in-time snapshot taken under the
// lock. They may therefore safely re-enter the same cache (read it, or mutate
// it) without deadlocking. The eviction callback registered with WithOnEvict
// behaves the same way: it fires after the lock is released, so it may safely
// call back into the same cache without deadlocking — but keep it cheap and
// non-blocking to avoid stalling the operation that triggered the eviction.
//
// # Performance
//
//	Operation       | Cost
//	----------------|------
//	Get / Peek      | O(1)
//	Put / PutInPlace| O(1)
//	Remove          | O(1) in place, O(n) for the immutable copy
//	Keys/Values/... | O(n)
//
// Start with NewLRU and upgrade to a concurrent variant only when you actually
// share the cache across goroutines.
package lru
