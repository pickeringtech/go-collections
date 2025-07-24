// Package sets provides mathematical set operations with automatic deduplication
// and rich set operations like union, intersection, and difference.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/collections/sets"
//
//	// Create sets with unique elements
//	permissions := sets.NewHash("read", "write", "execute")
//	userPerms := sets.NewHash("read", "write")
//
//	// Mathematical operations
//	canExecute := permissions.Contains("execute")           // true
//	common := permissions.Intersection(userPerms)          // {read, write}
//	missing := permissions.Difference(userPerms)           // {execute}
//	isSubset := userPerms.IsSubsetOf(permissions)          // true
//
// # Available Implementations
//
// Hash Set (sets.Hash):
//   - Fast O(1) operations using Go's built-in map
//   - Perfect for general-purpose unique collections
//   - Single-threaded use
//
// Concurrent Hash Set (sets.ConcurrentHash):
//   - Thread-safe with mutex protection
//   - O(1) operations with locking overhead
//   - Perfect for balanced read/write workloads
//
// Concurrent RW Hash Set (sets.ConcurrentHashRW):
//   - Thread-safe with read-write mutex
//   - Concurrent reads, exclusive writes
//   - Perfect for read-heavy workloads
//
// # Mathematical Set Operations
//
// Sets provide all standard mathematical operations:
//
//	s1 := sets.NewHash(1, 2, 3, 4)
//	s2 := sets.NewHash(3, 4, 5, 6)
//
//	union := s1.Union(s2)                    // {1, 2, 3, 4, 5, 6}
//	intersection := s1.Intersection(s2)      // {3, 4}
//	difference := s1.Difference(s2)          // {1, 2}
//
//	isSubset := s1.IsSubsetOf(s2)            // false
//	isSuperset := s1.IsSupersetOf(s2)        // false
//	areDisjoint := s1.IsDisjoint(s2)         // false (they share 3, 4)
//	areEqual := s1.Equals(s2)                // false
//
// # Immutable vs Mutable Operations
//
// Immutable operations return new sets:
//
//	newSet := set.Add(element)               // Returns new set
//	filtered := set.Filter(predicate)       // Returns new set
//	union := set.Union(otherSet)             // Returns new set
//
// Mutable operations modify the original:
//
//	set.AddInPlace(element)                  // Modifies original
//	set.FilterInPlace(predicate)             // Modifies original
//	set.UnionInPlace(otherSet)               // Modifies original
//
// # Thread Safety
//
// Choose the right concurrent implementation:
//
//	// Balanced read/write workloads
//	activeUsers := sets.NewConcurrentHash[string]()
//
//	// Read-heavy workloads (concurrent reads)
//	permissions := sets.NewConcurrentHashRW[string]()
//
// # Common Patterns
//
// Permission system:
//
//	adminPerms := sets.NewHash("read", "write", "delete", "admin")
//	userPerms := sets.NewHash("read", "write")
//
//	func canPerform(userRole sets.Set[string], action string) bool {
//		return userRole.Contains(action)
//	}
//
//	// Find common permissions
//	common := adminPerms.Intersection(userPerms)
//
// Deduplication:
//
//	// Remove duplicates from slice
//	items := []string{"apple", "banana", "apple", "cherry", "banana"}
//	unique := sets.NewHash(items...)
//	deduplicated := unique.AsSlice()  // ["apple", "banana", "cherry"]
//
// Tag management:
//
//	postTags := sets.NewHash("go", "programming", "tutorial")
//	userInterests := sets.NewHash("go", "rust", "programming")
//
//	// Find matching interests
//	matches := postTags.Intersection(userInterests)  // {"go", "programming"}
//	relevanceScore := float64(matches.Length()) / float64(postTags.Length())
//
// # Performance
//
//	BenchmarkSet_Contains/Hash-16            200M    6.123 ns/op
//	BenchmarkSet_Add/Hash-16                 180M    7.456 ns/op
//	BenchmarkSet_Union/Hash-16                50M   28.34 ns/op
//
// Start with NewHash() and upgrade to concurrent versions only when needed.
package sets
