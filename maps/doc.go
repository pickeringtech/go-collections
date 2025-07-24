// Package maps provides functional programming utilities for Go's native maps,
// enabling elegant data transformation and processing without manual iteration.
// These utilities work directly with Go's built-in map[K]V type.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/maps"
//
//	// Transform native Go maps functionally
//	inventory := map[string]int{
//		"apples":  50,
//		"oranges": 30,
//		"bananas": 20,
//	}
//
//	// Filter for low stock items
//	lowStock := maps.Filter(inventory, func(item string, count int) bool {
//		return count < 40
//	})
//	// Result: {"oranges": 30, "bananas": 20}
//
//	// Transform values
//	doubled := maps.MapValues(inventory, func(count int) int {
//		return count * 2
//	})
//	// Result: {"apples": 100, "oranges": 60, "bananas": 40}
//
// # When to Use Maps vs Collections/Dicts
//
// Use maps package when:
//   - Working with existing native Go maps
//   - Need simple transformations on map data
//   - Integrating with APIs that return map[K]V
//   - Want functional operations without changing data structures
//
// Use collections/dicts when:
//   - Need rich operations like Find, Contains, etc.
//   - Want thread-safe concurrent access
//   - Need both immutable and mutable operations
//   - Building new applications from scratch
//
// # Core Operations
//
// Transform Operations:
//   - Filter: Keep key-value pairs matching condition
//   - MapKeys: Transform all keys
//   - MapValues: Transform all values
//   - Map: Transform both keys and values
//
// Utility Operations:
//   - Keys: Extract all keys as slice
//   - Values: Extract all values as slice
//   - Invert: Swap keys and values
//   - Merge: Combine multiple maps
//
// # Common Patterns
//
// Configuration Processing:
//
//	config := map[string]string{
//		"database.host": "localhost",
//		"database.port": "5432",
//		"server.port":   "8080",
//		"debug.enabled": "true",
//	}
//
//	// Extract database config
//	dbConfig := maps.Filter(config, func(key, value string) bool {
//		return strings.HasPrefix(key, "database.")
//	})
//
//	// Convert to integers where needed
//	ports := maps.Filter(config, func(key, value string) bool {
//		return strings.HasSuffix(key, ".port")
//	})
//	portInts := maps.MapValues(ports, func(port string) int {
//		p, _ := strconv.Atoi(port)
//		return p
//	})
//
// Data Transformation:
//
//	userScores := map[string]int{
//		"alice":   95,
//		"bob":     87,
//		"charlie": 92,
//	}
//
//	// Convert scores to grades
//	grades := maps.MapValues(userScores, func(score int) string {
//		if score >= 90 { return "A" }
//		if score >= 80 { return "B" }
//		return "C"
//	})
//
//	// Create reverse lookup (grade -> users)
//	gradeUsers := make(map[string][]string)
//	for user, grade := range grades {
//		gradeUsers[grade] = append(gradeUsers[grade], user)
//	}
//
// API Response Processing:
//
//	// Process API response
//	response := map[string]interface{}{
//		"user_id":    123,
//		"user_name":  "alice",
//		"user_email": "alice@example.com",
//		"created_at": "2023-01-01T00:00:00Z",
//	}
//
//	// Extract user fields only
//	userFields := maps.Filter(response, func(key string, value interface{}) bool {
//		return strings.HasPrefix(key, "user_")
//	})
//
//	// Clean up field names
//	cleanFields := maps.MapKeys(userFields, func(key string) string {
//		return strings.TrimPrefix(key, "user_")
//	})
//
// # Performance Considerations
//
// Maps package operations create new maps and may be slower than manual iteration
// for performance-critical code. However, they offer significant benefits:
//
//   - Reduced bugs through functional style
//   - Clearer, more maintainable code
//   - Easier testing and reasoning
//   - Composable operations
//
// Use functional style for:
//   - Business logic and data transformation
//   - Configuration processing
//   - API response handling
//   - Code that prioritizes readability
//
// Use manual iteration for:
//   - Performance-critical hot paths
//   - Very large maps (>10k entries)
//   - Memory-constrained environments
//
// # Integration with Other Packages
//
// Maps package works well with slices and collections:
//
//	// Extract and process map data
//	inventory := map[string]int{"apples": 50, "oranges": 30}
//	
//	// Get keys and process with slices
//	items := maps.Keys(inventory)
//	sortedItems := slices.Sort(items, func(a, b string) bool { return a < b })
//	
//	// Store in collections
//	itemSet := collections.NewSet(items...)
//	itemDict := collections.NewDict(
//		slices.Map(items, func(item string) collections.Pair[string, int] {
//			return collections.Pair[string, int]{Key: item, Value: inventory[item]}
//		})...
//	)
//
// Start with simple Filter and MapValues operations, then explore advanced
// patterns like key transformation and map merging as needed.
package maps
