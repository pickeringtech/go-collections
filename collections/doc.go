// Package collections provides comprehensive, type-safe data structures for Go.
//
// This package offers three core data structures with rich APIs and multiple implementations:
//
// # Dicts - Key-Value Mappings
//
// Dicts provide powerful key-value mappings that go beyond Go's built-in maps:
//
//	// Create a dictionary with rich operations
//	users := collections.NewDict(
//		collections.Pair[int, string]{Key: 1, Value: "Alice"},
//		collections.Pair[int, string]{Key: 2, Value: "Bob"},
//	)
//
//	// Rich operations not available in native maps
//	activeUsers := users.Filter(func(id int, name string) bool {
//		return isActive(id)
//	})
//
//	// Thread-safe variants available
//	cache := collections.NewConcurrentDict(pairs...)
//
// # Sets - Unique Collections
//
// Sets provide mathematical set operations with automatic deduplication:
//
//	// Create sets with unique elements
//	permissions := collections.NewSet("read", "write", "execute")
//	userPerms := collections.NewSet("read", "write")
//
//	// Mathematical operations
//	common := permissions.Intersection(userPerms)    // {read, write}
//	missing := permissions.Difference(userPerms)     // {execute}
//	isSubset := userPerms.IsSubsetOf(permissions)    // true
//
// # Lists - Ordered Sequences
//
// Lists provide flexible ordered collections with stack and queue operations:
//
//	// Create a task queue
//	tasks := collections.NewList("design", "implement", "test")
//
//	// Stack operations (LIFO)
//	tasks.PushInPlace("deploy")
//	lastTask, found := tasks.PopInPlace()
//
//	// Queue operations (FIFO)
//	tasks.EnqueueInPlace("monitor")
//	firstTask, found := tasks.DequeueInPlace()
//
//	// Rich operations
//	longTasks := tasks.Filter(func(task string) bool {
//		return len(task) > 4
//	})
//
// # Multimaps - One Key, Many Values
//
// Multimaps map a single key to many values, replacing hand-rolled
// map[K][]V plumbing:
//
//	// Group orders by customer (ordered, duplicates kept)
//	orders := collections.NewListMultimap[string, string]()
//	orders.PutInPlace("alice", "book")
//	orders.PutInPlace("alice", "pen")
//	alice := orders.Get("alice") // [book pen]
//
//	// Tag documents (duplicates collapsed)
//	tags := collections.NewSetMultimap[string, string]()
//
// # Heaps - Priority Queues
//
// Heaps always hand back the most- (or least-) extreme element next:
//
//	// Smallest-first (min-heap) over an ordered type
//	pq := collections.NewMinHeap(5, 1, 3)
//	next, ok, rest := pq.Pop() // next == 1, ok == true; rest is pq without it
//
//	// Or order by any comparator, e.g. a struct field
//	tasks := collections.NewHeap(func(a, b Task) bool { return a.Priority > b.Priority })
//
// The facade constructors return the immutable heaps.Heap interface (Push/Pop
// return a new heap). For the in-place mutating API (PushInPlace/PopInPlace),
// reach for the heaps subpackage directly.
//
// # LRU - Bounded Caches
//
// LRU caches keep at most a fixed number of entries, evicting the
// least-recently-used one to make room:
//
//	cache := collections.NewLRU[string, int](2)
//	cache.PutInPlace("a", 1)
//	cache.PutInPlace("b", 2)
//	v, ok := cache.Get("a") // promotes "a"; a third insert now evicts "b"
//
// The facade constructors return lru.MutableCache because an LRU is inherently
// stateful: its defining recency-marking read, Get, is a mutation. Pass
// lru.Option values (e.g. lru.WithOnEvict) to configure eviction callbacks and
// seed entries.
//
// # Thread Safety
//
// All data structures offer thread-safe variants:
//
//	// Choose your concurrency model
//	dict := collections.NewConcurrentDict(...)     // Balanced read/write
//	dict := collections.NewConcurrentRWDict(...)   // Read-heavy workloads
//
//	set := collections.NewConcurrentSet(...)       // Balanced read/write
//	set := collections.NewConcurrentRWSet(...)     // Read-heavy workloads
//
//	list := collections.NewConcurrentList(...)     // Balanced read/write
//	list := collections.NewConcurrentRWList(...)   // Read-heavy workloads
//
// # Immutable vs Mutable Operations
//
// All collections support both paradigms:
//
//	// Immutable style - returns new collections
//	newDict := dict.Put("key", value)
//	filtered := set.Filter(predicate)
//	newList := list.Push(element)
//
//	// Mutable style - modifies in place
//	dict.PutInPlace("key", value)
//	set.FilterInPlace(predicate)
//	list.PushInPlace(element)
//
// # Construction and Zero Values
//
// The library-wide convention is to build every collection with its New*
// constructor (NewDict, NewConcurrentDict, NewList, and so on). Treat the zero
// value as not ready for use unless a type's own documentation says otherwise:
// most implementations leave their backing map, tree, or wrapped collection nil
// until the constructor runs, so a write to a bare CollectionType{} panics.
//
// Concurrent types embed their sync.Mutex / sync.RWMutex by value. The lock is
// always safe to take on the zero value and reads return empty results, but the
// backing data is not initialized until the constructor runs, so the constructor
// is still required before writing.
//
// Important: concurrent types must NOT be copied after first use. A copy
// produces an independent lock while both values share the same backing data
// (map, slice, inner pointer), silently breaking the thread-safety guarantee.
// Every concurrent type embeds a nocopy sentinel so that go vet's copylocks
// analyser reports any value-copy after construction. Always pass concurrent
// collections by pointer.
//
// A few types document a usable zero value as part of their contract — for
// example deques.RingBuffer (a valid empty, unbounded deque) and dicts.Tree (a
// valid empty tree). Even for these, the constructor remains the recommended
// entry point so call sites read consistently.
//
// # Performance
//
// All implementations are benchmarked and optimized:
//
//	BenchmarkDict_Get/Hash-16                228M    5.248 ns/op
//	BenchmarkDict_Get/ConcurrentHash-16      100M   10.41 ns/op
//	BenchmarkSet_Contains/Hash-16            200M    6.123 ns/op
//	BenchmarkList_Push/Linked-16             150M    8.456 ns/op
//
// Start with the simple variants (NewDict, NewSet, NewList) and upgrade to
// concurrent versions only when you need thread safety.
package collections
