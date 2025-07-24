// Package dicts provides comprehensive dictionary/map implementations with a focus on
// performance, thread-safety, and rich functionality.
//
// The package offers both immutable and mutable interfaces, allowing you to choose
// the right approach for your use case. All implementations are fully generic and
// type-safe.
//
// # Implementations
//
// Hash Dictionary (Hash[K, V]): A hash table implementation using Go's built-in map.
// Provides O(1) average case performance for basic operations.
//
// Concurrent Hash Dictionary (ConcurrentHash[K, V]): A thread-safe hash dictionary
// using a mutex for synchronization. All operations are atomic.
//
// Concurrent RW Hash Dictionary (ConcurrentHashRW[K, V]): A thread-safe hash dictionary
// using a read-write mutex. Read operations can proceed concurrently, while write
// operations are exclusive.
//
// Tree Dictionary (Tree[K, V]): A binary search tree implementation that maintains
// keys in sorted order. Provides O(log n) average case performance and ordered iteration.
// Keys must implement constraints.Ordered (integers, floats, strings).
//
// # Interfaces
//
// Dict[K, V]: Immutable dictionary interface that provides comprehensive key-value
// operations without modifying the original dictionary.
//
// MutableDict[K, V]: Mutable dictionary interface that provides comprehensive key-value
// operations with the ability to modify the dictionary in place.
//
// # Performance Characteristics
//
//	Implementation     | Get   | Put   | Remove | Memory | Thread-Safe
//	-------------------|-------|-------|--------|--------|-------------
//	Hash               | O(1)  | O(1)  | O(1)   | Low    | No
//	ConcurrentHash     | O(1)  | O(1)  | O(1)   | Low    | Yes
//	ConcurrentHashRW   | O(1)  | O(1)  | O(1)   | Low    | Yes
//	Tree               | O(log n) | O(log n) | O(log n) | Medium | No
//
// # Benchmark Results
//
//	BenchmarkComparison_Get/Hash-16                228195298    5.248 ns/op    0 B/op    0 allocs/op
//	BenchmarkComparison_Get/ConcurrentHash-16      100000000   10.41 ns/op     0 B/op    0 allocs/op
//	BenchmarkComparison_Get/ConcurrentHashRW-16    100000000   10.30 ns/op     0 B/op    0 allocs/op
//	BenchmarkComparison_Get/NativeMap-16           226149249    5.437 ns/op    0 B/op    0 allocs/op
//
// # Usage Examples
//
// Basic Hash Dictionary:
//
//	h := dicts.NewHash(
//	    dicts.Pair[string, int]{Key: "one", Value: 1},
//	    dicts.Pair[string, int]{Key: "two", Value: 2},
//	)
//
//	value, found := h.Get("one", -1)
//	fmt.Printf("Value: %d, Found: %t\n", value, found) // Value: 1, Found: true
//
//	// Immutable operations
//	newH := h.Put("three", 3)
//	filtered := h.Filter(func(key string, value int) bool {
//	    return value > 1
//	})
//
// Concurrent Dictionary:
//
//	ch := dicts.NewConcurrentHash[string, int]()
//	
//	// Safe to use from multiple goroutines
//	go func() {
//	    if mutableDict, ok := ch.(dicts.MutableDict[string, int]); ok {
//	        mutableDict.PutInPlace("key", 42)
//	    }
//	}()
//
// Tree Dictionary (Sorted):
//
//	tree := dicts.NewTree(
//	    dicts.Pair[string, int]{Key: "charlie", Value: 3},
//	    dicts.Pair[string, int]{Key: "alice", Value: 1},
//	    dicts.Pair[string, int]{Key: "bob", Value: 2},
//	)
//
//	// Iterate in sorted order
//	tree.ForEach(func(key string, value int) {
//	    fmt.Printf("%s: %d\n", key, value)
//	})
//	// Output:
//	// alice: 1
//	// bob: 2
//	// charlie: 3
//
// # Best Practices
//
// Choose the Right Implementation:
//   - Use Hash for single-threaded, high-performance scenarios
//   - Use ConcurrentHash for multi-threaded scenarios with balanced read/write
//   - Use ConcurrentHashRW for read-heavy multi-threaded scenarios
//   - Use Tree when you need sorted iteration or range queries
//
// Immutable vs Mutable:
//   - Use immutable operations (Put, Remove, Filter) for functional programming style
//   - Use mutable operations (PutInPlace, RemoveInPlace, FilterInPlace) for performance-critical scenarios
//
// Memory Management:
//   - Immutable operations create new dictionaries; be mindful of memory usage
//   - Use Clear() to reset large dictionaries instead of creating new ones
//
// Error Handling:
//   - Always check the boolean return value from Get() operations
//   - Use Contains() for existence checks when you don't need the value
// Package dicts provides powerful key-value mappings that extend Go's built-in maps
// with rich operations, thread safety, and both immutable and mutable interfaces.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/collections/dicts"
//
//	// Create a dictionary with initial data
//	inventory := dicts.NewHash(
//		dicts.Pair[string, int]{Key: "apples", Value: 50},
//		dicts.Pair[string, int]{Key: "oranges", Value: 30},
//	)
//
//	// Rich operations that native maps can't do
//	lowStock := inventory.Filter(func(item string, count int) bool {
//		return count < 40
//	})
//
// # Available Implementations
//
// Hash Dictionary (dicts.Hash):
//   - Fast O(1) operations using Go's built-in map
//   - Perfect for general-purpose key-value storage
//   - Single-threaded use
//
// Concurrent Hash Dictionary (dicts.ConcurrentHash):
//   - Thread-safe with mutex protection
//   - O(1) operations with locking overhead
//   - Perfect for balanced read/write workloads
//
// Concurrent RW Hash Dictionary (dicts.ConcurrentHashRW):
//   - Thread-safe with read-write mutex
//   - Concurrent reads, exclusive writes
//   - Perfect for read-heavy workloads (10:1 read/write ratio)
//
// Tree Dictionary (dicts.Tree):
//   - Maintains keys in sorted order
//   - O(log n) operations
//   - Perfect when you need sorted iteration
//
// # Immutable vs Mutable Operations
//
// Immutable operations return new dictionaries:
//
//	newDict := dict.Put("key", value)        // Returns new dict
//	filtered := dict.Filter(predicate)       // Returns new dict
//	removed := dict.Remove("key")            // Returns new dict
//
// Mutable operations modify the original:
//
//	dict.PutInPlace("key", value)            // Modifies original
//	dict.FilterInPlace(predicate)            // Modifies original
//	removed, found := dict.RemoveInPlace("key") // Modifies original
//
// # Thread Safety
//
// Choose the right concurrent implementation:
//
//	// Balanced read/write workloads
//	cache := dicts.NewConcurrentHash[string, []byte]()
//
//	// Read-heavy workloads (concurrent reads)
//	config := dicts.NewConcurrentHashRW[string, string]()
//
// # Common Patterns
//
// Web application cache:
//
//	cache := dicts.NewConcurrentHashRW[string, []byte]()
//
//	func getPage(url string) []byte {
//		if content, found := cache.Get(url, nil); found {
//			return content
//		}
//		content := fetchFromNetwork(url)
//		cache.PutInPlace(url, content)
//		return content
//	}
//
// Configuration management:
//
//	config := dicts.NewTree(
//		dicts.Pair[string, string]{Key: "database.host", Value: "localhost"},
//		dicts.Pair[string, string]{Key: "server.port", Value: "8080"},
//	)
//
//	// Iterate in sorted order
//	config.ForEach(func(key, value string) {
//		fmt.Printf("%s = %s\n", key, value)
//	})
//
// # Performance
//
//	BenchmarkDict_Get/Hash-16                231M    5.188 ns/op
//	BenchmarkDict_Get/ConcurrentHash-16      121M   10.02 ns/op
//	BenchmarkDict_Get/ConcurrentHashRW-16    123M    9.814 ns/op
//	BenchmarkDict_Get/Tree-16                 50M   25.67 ns/op
//
// Start with NewHash() and upgrade to concurrent versions only when needed.
package dicts
