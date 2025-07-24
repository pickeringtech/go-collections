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
