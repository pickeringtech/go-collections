// Package maps provides functional programming utilities for Go's native maps,
// enabling elegant data transformation and processing without manual iteration.
// These utilities work directly with Go's built-in map[K]V type.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/maps"
//
//	// Transform native Go maps functionally.
//	inventory := map[string]int{
//		"apples":  50,
//		"oranges": 30,
//		"bananas": 20,
//	}
//
//	// Filter for low-stock items.
//	lowStock := maps.Filter(inventory, func(item string, count int) bool {
//		return count < 40
//	})
//	// lowStock: {"oranges": 30, "bananas": 20}
//
//	// Transform entries. Map rebuilds the map and can change the key, the
//	// value, or both; here it doubles each value and keeps the key.
//	doubled := maps.Map(inventory, func(item string, count int) (string, int) {
//		return item, count * 2
//	})
//	// doubled: {"apples": 100, "oranges": 60, "bananas": 40}
//
// This Quick Start is compiled and run as Example_quickStart in the package's
// test suite, so it is guaranteed to track the real API.
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
// Transform operations (each returns a new map without mutating the input):
//   - Filter: keep key-value pairs matching a predicate
//   - Map: rebuild the map, transforming each key and/or value
//   - Update: merge another map's entries over a copy of the input
//
// Lookup and extraction:
//   - Keys / Values: extract all keys or values as a slice
//   - Items: extract all entries as a slice of Entry[K, V]
//   - GetOrDefault / GetMany / GetManyOrDefault: read values with fallbacks
//   - ContainsValue: check whether a value is present
//
// Construction and utilities:
//   - FromKeys: build a map from a slice of keys and a default value
//   - Copy: shallow-copy a map
//   - Clear: remove every entry from a map in place
//
// # Common Patterns
//
// Transform values with Map (swap, recompute, or relabel entries):
//
//	userScores := map[string]int{
//		"alice":   95,
//		"bob":     87,
//		"charlie": 92,
//	}
//
//	// Convert scores to letter grades, keeping the user as the key.
//	grades := maps.Map(userScores, func(user string, score int) (string, string) {
//		switch {
//		case score >= 90:
//			return user, "A"
//		case score >= 80:
//			return user, "B"
//		default:
//			return user, "C"
//		}
//	})
//
// Invert a map by swapping keys and values in the mapping function:
//
//	idToName := map[int]string{1: "alice", 2: "bob"}
//	nameToID := maps.Map(idToName, func(id int, name string) (string, int) {
//		return name, id
//	})
//
// Merge maps with Update (entries from the second argument win):
//
//	defaults := map[string]string{"host": "localhost", "port": "8080"}
//	overrides := map[string]string{"port": "9090"}
//	settings := maps.Update(defaults, overrides)
//	// settings: {"host": "localhost", "port": "9090"}
//
// Filter then extract:
//
//	config := map[string]string{
//		"database.host": "localhost",
//		"database.port": "5432",
//		"server.port":   "8080",
//	}
//
//	dbConfig := maps.Filter(config, func(key, value string) bool {
//		return strings.HasPrefix(key, "database.")
//	})
//	dbKeys := maps.Keys(dbConfig)
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
// The maps package works well with slices and collections:
//
//	// Extract and process map data.
//	inventory := map[string]int{"apples": 50, "oranges": 30}
//
//	// Get keys and process them with the slices package.
//	items := maps.Keys(inventory)
//	slices.SortOrderedAscInPlace(items)
//
//	// Store the data in collections.
//	itemSet := collections.NewSet(items...)
//	itemDict := collections.NewDict(
//		slices.Map(items, func(item string) collections.Pair[string, int] {
//			return collections.Pair[string, int]{Key: item, Value: inventory[item]}
//		})...,
//	)
//
// Start with simple Filter and Map operations, then reach for Update to merge
// maps and the Keys/Values/Items extractors as needed.
package maps
