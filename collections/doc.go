// Package collections provides comprehensive, type-safe data structures for Go.
//
// This package is the facade hub for seven families of data structure — dicts,
// sets, lists, multimaps, deques, heaps and LRU caches — each with rich APIs and
// multiple implementations:
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
//	// Stack operations (LIFO). The facade returns a MutableList, so the
//	// InPlace forms modify it directly.
//	tasks.PushInPlace("deploy")
//	lastTask, found := tasks.PopInPlace() // "deploy"
//
//	// Queue operations (FIFO).
//	tasks.EnqueueInPlace("monitor")
//	firstTask, ok := tasks.DequeueInPlace() // "design"
//
//	// The immutable forms return a new list and leave tasks untouched.
//	withDeploy := tasks.Push("deploy")
//	last, popped, remaining := withDeploy.Pop()
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
//	// Group orders by customer (ordered, duplicates kept). The facade returns a
//	// MutableMultimap, so PutInPlace adds to it directly (Put would instead
//	// return a new multimap).
//	orders := collections.NewListMultimap[string, string]()
//	orders.PutInPlace("alice", "book")
//	orders.PutInPlace("alice", "pen")
//	alice := orders.Get("alice") // [book pen]
//
//	// Tag documents (duplicates collapsed)
//	tags := collections.NewSetMultimap[string, string]()
//
// # Deques - Double-Ended Queues
//
// Deques push and pop at both ends in O(1), backed by a ring buffer, and can be
// bounded to act as a fixed-capacity circular buffer:
//
//	// Unbounded, growable at both ends. The facade returns a MutableDeque, so
//	// the InPlace forms modify it directly; PushFront/PushBack/PopFront would
//	// instead return a new deque.
//	dq := collections.NewDeque(1, 2, 3)
//	dq.PushFrontInPlace(0)            // front: 0 1 2 3
//	first, ok := dq.PopFrontInPlace() // first == 0, ok == true
//
//	// Bounded: keep at most the last N elements. With OverwriteOldest a push
//	// onto a full deque drops the element at the opposite end (a sliding window);
//	// RejectWhenFull makes the push a no-op instead.
//	window := collections.NewBoundedDeque[int](2, deques.OverwriteOldest)
//	window.PushBackInPlace(1)
//	window.PushBackInPlace(2)
//	window.PushBackInPlace(3) // drops 1; window holds 2 3
//
// # Heaps - Priority Queues
//
// Heaps always hand back the most- (or least-) extreme element next:
//
//	// Smallest-first (min-heap) over an ordered type
//	pq := collections.NewMinHeap(5, 1, 3)
//	pq.PushInPlace(0)
//	next, ok := pq.PopInPlace() // next == 0, ok == true
//
//	// Or order by any comparator, e.g. a struct field
//	tasks := collections.NewHeap(func(a, b Task) bool { return a.Priority > b.Priority })
//
// The facade constructors return heaps.MutableHeap, so the O(log n) in-place
// PushInPlace/PopInPlace are reachable directly. The immutable Push/Pop are
// still available; they return a new heap, copying the receiver (O(n)).
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
// # Sorted and Alternative Implementations
//
// The facade constructors return the everyday default for each family (a
// hash-backed dict or set, an array-backed list). When you need a sorted or
// otherwise alternative implementation, reach into the spoke package directly:
//
//	tree := dicts.NewTree(...)      // sorted, ordered-iteration dictionary
//	sorted := sets.NewTreeSet(...)  // sorted set with ordered iteration
//	linked := lists.NewLinked(...)  // singly linked list (also collections.NewLinkedList)
//
// See each spoke package's documentation for the full set of implementations
// and their trade-offs.
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
// Concurrent types embed their sync.Mutex / sync.RWMutex by value rather than by
// pointer, so the lock alone is always safe to take. They do not initialize the
// backing data, though, so the constructor is still required before writing.
//
// Zero-value reads split into two groups by how a type stores its data. The map-
// and slice-backed concurrent types hold their backing collection by value, so a
// nil map or slice reads as empty and reads on their zero value return empty
// results: ConcurrentHash (dicts and sets), ConcurrentArray (lists), the
// multimaps (ConcurrentListMultimap, ConcurrentSetMultimap) and all their RW
// variants. The remaining concurrent types wrap an inner collection by pointer:
// the linked lists (ConcurrentLinked, ConcurrentDoublyLinked), the tree-backed
// dicts and sets (ConcurrentTree, ConcurrentTreeSet), the heaps
// (ConcurrentBinary), the ring buffers (ConcurrentRingBuffer), the LRU caches
// (ConcurrentLRU) and all their RW variants. Their inner pointer is nil until the
// constructor runs, so any operation — reads included — dereferences a nil
// pointer and panics. Construct those with their New* constructor before use;
// each type's own documentation repeats the warning.
//
// Important: concurrent types must NOT be copied after first use. A copy
// produces an independent lock while both values share the same backing data
// (map, slice, inner pointer), silently breaking the thread-safety guarantee.
// Every concurrent type carries a nocopy sentinel field so that go vet's
// copylocks analyser reports any value-copy of the type. Always pass concurrent
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
