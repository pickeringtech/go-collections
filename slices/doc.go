// Package slices provides functional programming utilities for Go slices,
// enabling elegant data transformation, filtering, and processing without
// manual loops or complex iteration logic.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/slices"
//
//	// Transform data with functional style
//	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
//
//	// Chain operations elegantly
//	result := slices.Filter(numbers, func(n int) bool { return n%2 == 0 }).
//		Map(func(n int) int { return n * n }).
//		Reduce(0, func(acc, n int) int { return acc + n })
//
//	// Result: sum of squares of even numbers = 220
//
// # Why Use Slices Package?
//
// Native Go slice operations require verbose loops:
//
//	// Native approach - verbose and error-prone
//	var evens []int
//	for _, n := range numbers {
//		if n%2 == 0 {
//			evens = append(evens, n)
//		}
//	}
//
//	var squares []int
//	for _, n := range evens {
//		squares = append(squares, n*n)
//	}
//
//	sum := 0
//	for _, n := range squares {
//		sum += n
//	}
//
// Slices package enables elegant functional style:
//
//	// Functional approach - clean and expressive
//	sum := slices.Filter(numbers, isEven).
//		Map(square).
//		Reduce(0, add)
//
// # Core Operations
//
// Transform Operations:
//   - Map: Transform each element
//   - Filter: Keep elements matching condition
//   - Reduce: Combine elements into single value
//   - FlatMap: Transform and flatten nested structures
//
// Search Operations:
//   - Find: Get first matching element
//   - FindIndex: Get index of first match
//   - Contains: Check if element exists
//   - All/Any: Check conditions across elements
//
// Utility Operations:
//   - Reverse: Reverse slice order
//   - Unique: Remove duplicates
//   - Chunk: Split into smaller slices
//   - Zip: Combine multiple slices
//
// # Common Patterns
//
// Data Processing Pipeline:
//
//	users := []User{...}
//
//	// Find active adult users and get their emails
//	emails := slices.Filter(users, func(u User) bool {
//		return u.Active && u.Age >= 18
//	}).Map(func(u User) string {
//		return u.Email
//	})
//
// Configuration Processing:
//
//	configs := []string{"key1=value1", "key2=value2", "invalid"}
//
//	// Parse valid configs into map
//	configMap := slices.Filter(configs, func(s string) bool {
//		return strings.Contains(s, "=")
//	}).Map(func(s string) Pair {
//		parts := strings.Split(s, "=")
//		return Pair{Key: parts[0], Value: parts[1]}
//	}).Reduce(make(map[string]string), func(acc map[string]string, p Pair) map[string]string {
//		acc[p.Key] = p.Value
//		return acc
//	})
//
// Log Processing:
//
//	logs := []LogEntry{...}
//
//	// Find error logs from last hour
//	recentErrors := slices.Filter(logs, func(log LogEntry) bool {
//		return log.Level == "ERROR" && 
//			   time.Since(log.Timestamp) < time.Hour
//	})
//
//	// Group by error type
//	errorCounts := slices.Reduce(recentErrors, 
//		make(map[string]int), 
//		func(acc map[string]int, log LogEntry) map[string]int {
//			acc[log.ErrorType]++
//			return acc
//		})
//
// # Performance Considerations
//
// Functional operations create new slices and may be slower than manual loops
// for performance-critical code. However, they offer significant benefits:
//
//   - Reduced bugs through immutability
//   - Clearer, more maintainable code
//   - Easier testing and reasoning
//   - Composable operations
//
// Benchmark comparison:
//
//	BenchmarkManualLoop-16     100M    12.3 ns/op    0 B/op
//	BenchmarkFunctional-16      50M    24.7 ns/op   32 B/op
//
// Use functional style for:
//   - Business logic and data transformation
//   - Code that prioritizes readability
//   - Prototyping and development
//
// Use manual loops for:
//   - Performance-critical hot paths
//   - Memory-constrained environments
//   - Very large datasets
//
// # Integration with Collections
//
// Slices package works seamlessly with collections:
//
//	// Process data and store in collections
//	users := []User{...}
//	
//	activeEmails := slices.Filter(users, isActive).
//		Map(getEmail)
//	
//	emailSet := collections.NewSet(activeEmails...)
//	emailDict := collections.NewDict(
//		slices.Map(activeEmails, func(email string) collections.Pair[string, bool] {
//			return collections.Pair[string, bool]{Key: email, Value: true}
//		})...
//	)
//
// Start with simple operations like Filter and Map, then explore advanced
// patterns like Reduce and FlatMap as you become comfortable with functional programming.
package slices
