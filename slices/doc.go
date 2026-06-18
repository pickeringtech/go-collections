// Package slices provides functional programming utilities for Go slices,
// enabling elegant data transformation, filtering, and processing without
// manual loops or complex iteration logic.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/slices"
//
//	// Transform data with functional style.
//	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
//
//	// Each operation is a standalone function that takes a slice and returns a
//	// new one, so compose them by nesting the calls.
//	evens := slices.Filter(numbers, func(n int) bool { return n%2 == 0 })
//	squares := slices.Map(evens, func(n int) int { return n * n })
//	sum := slices.Reduce(squares, func(acc, n int) int { return acc + n })
//
//	// sum of the squares of the even numbers = 220
//
// This Quick Start is compiled and run as Example_quickStart in the package's
// test suite, so it is guaranteed to track the real API.
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
// The slices package collapses each loop into a named operation:
//
//	// Functional approach - clean and expressive
//	evens := slices.Filter(numbers, isEven)
//	squares := slices.Map(evens, square)
//	sum := slices.Reduce(squares, add)
//
// # Core Operations
//
// Transform operations (each returns a new slice or value):
//   - Map: transform each element into a new slice
//   - Filter: keep only the elements matching a predicate
//   - Reduce: fold the elements into a single value
//   - Reverse: reverse the slice order
//
// Search operations:
//   - Find / FindLast: get the first (or last) matching element
//   - FindIndex / FindLastIndex / IndexOf: locate a match by index
//   - Includes: check whether an element is present
//   - AllMatch / AnyMatch: check a predicate across the slice
//
// Access and structure:
//   - First / PeekFront / PeekEnd / Get: read elements safely
//   - SubSlice / Paginate: take a window of elements
//   - Push / Pop / PushFront / PopFront / Insert / Delete: grow or shrink
//   - Concatenate / Copy / Fill / JoinToString: combine and format
//
// Numeric helpers:
//   - Min / Max and the SortOrdered* family (ordering reductions)
//   - Numeric summaries such as Sum and Mean live in the stats package; the
//     NumericSlice accessor type delegates to them.
//
// # Common Patterns
//
// Data Processing Pipeline:
//
//	users := []User{...}
//
//	// Find active adult users and get their emails.
//	adults := slices.Filter(users, func(u User) bool {
//		return u.Active && u.Age >= 18
//	})
//	emails := slices.Map(adults, func(u User) string {
//		return u.Email
//	})
//
// Building a Map by Reduction:
//
//	configs := []string{"key1=value1", "key2=value2", "invalid"}
//
//	// Parse valid "key=value" configs into a map.
//	valid := slices.Filter(configs, func(s string) bool {
//		return strings.Contains(s, "=")
//	})
//	configMap := slices.Reduce(valid, func(acc map[string]string, s string) map[string]string {
//		if acc == nil {
//			acc = map[string]string{}
//		}
//		parts := strings.SplitN(s, "=", 2)
//		acc[parts[0]] = parts[1]
//		return acc
//	})
//
// Grouping by Reduction:
//
//	logs := []LogEntry{...}
//
//	// Find error logs from the last hour.
//	recentErrors := slices.Filter(logs, func(log LogEntry) bool {
//		return log.Level == "ERROR" &&
//			time.Since(log.Timestamp) < time.Hour
//	})
//
//	// Count them by error type.
//	errorCounts := slices.Reduce(recentErrors, func(acc map[string]int, log LogEntry) map[string]int {
//		if acc == nil {
//			acc = map[string]int{}
//		}
//		acc[log.ErrorType]++
//		return acc
//	})
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
// The slices package works seamlessly with collections:
//
//	// Process data and store it in collections.
//	users := []User{...}
//
//	active := slices.Filter(users, isActive)
//	activeEmails := slices.Map(active, getEmail)
//
//	emailSet := collections.NewSet(activeEmails...)
//	emailDict := collections.NewDict(
//		slices.Map(activeEmails, func(email string) collections.Pair[string, bool] {
//			return collections.Pair[string, bool]{Key: email, Value: true}
//		})...,
//	)
//
// Start with simple operations like Filter and Map, then reach for Reduce to
// fold a slice into a single value as you become comfortable with the style.
package slices
