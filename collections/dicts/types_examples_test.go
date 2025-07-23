package dicts_test

import (
	"fmt"
	"github.com/pickeringtech/go-collections/collections/dicts"
)

func ExamplePair() {
	// Create individual pairs
	pair1 := dicts.Pair[string, int]{Key: "apple", Value: 5}
	pair2 := dicts.Pair[string, int]{Key: "banana", Value: 3}

	fmt.Printf("Pair 1: %s = %d\n", pair1.Key, pair1.Value)
	fmt.Printf("Pair 2: %s = %d\n", pair2.Key, pair2.Value)

	// Use pairs to create a dictionary
	dict := dicts.NewHash(pair1, pair2)
	fmt.Printf("Dictionary length: %d\n", dict.Length())

	// Output:
	// Pair 1: apple = 5
	// Pair 2: banana = 3
	// Dictionary length: 2
}

// Example showing how to work with the Dict interface
func ExampleDict() {
	// All implementations satisfy the Dict interface
	var dict dicts.Dict[string, int]

	// Can assign any implementation
	dict = dicts.NewHash(
		dicts.Pair[string, int]{Key: "red", Value: 1},
		dicts.Pair[string, int]{Key: "green", Value: 2},
	)

	// Use interface methods
	value, found := dict.Get("red", -1)
	fmt.Printf("Red: %d, found: %t\n", value, found)

	// Immutable operations return new Dict
	newDict := dict.Put("blue", 3)
	fmt.Printf("Original length: %d\n", dict.Length())
	fmt.Printf("New length: %d\n", newDict.Length())

	// Output:
	// Red: 1, found: true
	// Original length: 2
	// New length: 3
}

// Example showing how to work with the MutableDict interface
func ExampleMutableDict() {
	// Create a mutable dictionary
	var mutableDict dicts.MutableDict[string, int]
	mutableDict = dicts.NewHash(
		dicts.Pair[string, int]{Key: "count", Value: 0},
	)

	// Use mutable operations
	mutableDict.PutInPlace("count", 42)
	mutableDict.PutInPlace("total", 100)

	fmt.Printf("Count: %d\n", func() int {
		v, _ := mutableDict.Get("count", -1)
		return v
	}())
	fmt.Printf("Length: %d\n", mutableDict.Length())

	// Remove in place
	removed, found := mutableDict.RemoveInPlace("total")
	fmt.Printf("Removed: %d, found: %t\n", removed, found)
	fmt.Printf("Final length: %d\n", mutableDict.Length())

	// Output:
	// Count: 42
	// Length: 2
	// Removed: 100, found: true
	// Final length: 1
}

// Example showing polymorphic usage
func Example_polymorphicUsage() {
	// Function that works with any Dict implementation
	processDict := func(d dicts.Dict[string, int], name string) {
		fmt.Printf("%s - Length: %d\n", name, d.Length())
		
		// Find all values greater than 5
		filtered := d.Filter(func(key string, value int) bool {
			return value > 5
		})
		fmt.Printf("%s - High values: %d\n", name, filtered.Length())
	}

	// Works with Hash
	hash := dicts.NewHash(
		dicts.Pair[string, int]{Key: "small", Value: 3},
		dicts.Pair[string, int]{Key: "large", Value: 10},
	)
	processDict(hash, "Hash")

	// Works with Tree
	tree := dicts.NewTree(
		dicts.Pair[string, int]{Key: "small", Value: 3},
		dicts.Pair[string, int]{Key: "large", Value: 10},
	)
	processDict(tree, "Tree")

	// Works with Concurrent implementations
	concurrent := dicts.NewConcurrentHash(
		dicts.Pair[string, int]{Key: "small", Value: 3},
		dicts.Pair[string, int]{Key: "large", Value: 10},
	)
	processDict(concurrent, "Concurrent")

	// Output:
	// Hash - Length: 2
	// Hash - High values: 1
	// Tree - Length: 2
	// Tree - High values: 1
	// Concurrent - Length: 2
	// Concurrent - High values: 1
}

// Example showing functional programming style with immutable operations
func Example_functionalStyle() {
	// Start with initial data
	original := dicts.NewHash(
		dicts.Pair[string, int]{Key: "a", Value: 1},
		dicts.Pair[string, int]{Key: "b", Value: 2},
		dicts.Pair[string, int]{Key: "c", Value: 3},
	)

	// Chain immutable operations
	result := original.
		Put("d", 4).                    // Add new item
		Remove("a").                    // Remove item
		Filter(func(k string, v int) bool { // Keep only even values
			return v%2 == 0
		})

	fmt.Printf("Original length: %d\n", original.Length())
	fmt.Printf("Result length: %d\n", result.Length())
	fmt.Printf("Result contains 'b': %t\n", result.Contains("b"))
	fmt.Printf("Result contains 'd': %t\n", result.Contains("d"))

	// Output:
	// Original length: 3
	// Result length: 2
	// Result contains 'b': true
	// Result contains 'd': true
}

// Example showing different key types with Tree
func Example_treeKeyTypes() {
	// String keys
	stringTree := dicts.NewTree(
		dicts.Pair[string, int]{Key: "zebra", Value: 1},
		dicts.Pair[string, int]{Key: "apple", Value: 2},
	)
	fmt.Printf("String keys: %v\n", stringTree.Keys())

	// Integer keys
	intTree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 100, Value: "hundred"},
		dicts.Pair[int, string]{Key: 50, Value: "fifty"},
	)
	fmt.Printf("Int keys: %v\n", intTree.Keys())

	// Float keys
	floatTree := dicts.NewTree(
		dicts.Pair[float64, string]{Key: 3.14, Value: "pi"},
		dicts.Pair[float64, string]{Key: 2.71, Value: "e"},
	)
	fmt.Printf("Float keys: %v\n", floatTree.Keys())

	// Output:
	// String keys: [apple zebra]
	// Int keys: [50 100]
	// Float keys: [2.71 3.14]
}

// Example showing conversion between implementations
func Example_conversion() {
	// Start with a Hash
	hash := dicts.NewHash(
		dicts.Pair[string, int]{Key: "c", Value: 3},
		dicts.Pair[string, int]{Key: "a", Value: 1},
		dicts.Pair[string, int]{Key: "b", Value: 2},
	)

	// Convert to Tree for sorted iteration
	tree := dicts.NewTree(hash.Items()...)

	fmt.Printf("Hash keys (unordered): %v\n", hash.Keys())
	fmt.Printf("Tree keys (sorted): %v\n", tree.Keys())

	// Convert back to native Go map
	nativeMap := tree.AsMap()
	fmt.Printf("Native map: %v\n", nativeMap)

	// Output:
	// Hash keys (unordered): [c a b]
	// Tree keys (sorted): [a b c]
	// Native map: map[a:1 b:2 c:3]
}

// Example showing search operations
func Example_searchOperations() {
	dict := dicts.NewHash(
		dicts.Pair[string, int]{Key: "apple", Value: 5},
		dicts.Pair[string, int]{Key: "banana", Value: 3},
		dicts.Pair[string, int]{Key: "cherry", Value: 8},
	)

	// Find by key-value predicate
	key, value, found := dict.Find(func(k string, v int) bool {
		return v > 6
	})
	if found {
		fmt.Printf("Found high value: %s = %d\n", key, value)
	}

	// Find by key predicate
	_, found = dict.FindKey(func(k string) bool {
		return len(k) > 5
	})
	fmt.Printf("Found long key: %t\n", found)

	// Find by value predicate
	foundValue, found := dict.FindValue(func(v int) bool {
		return v%2 == 1 // odd number
	})
	if found {
		fmt.Printf("Found odd value: %d\n", foundValue)
	}

	// Check if value exists
	hasValue := dict.ContainsValue(3)
	fmt.Printf("Contains value 3: %t\n", hasValue)

	// Output:
	// Found high value: cherry = 8
	// Found long key: true
	// Found odd value: 5
	// Contains value 3: true
}
