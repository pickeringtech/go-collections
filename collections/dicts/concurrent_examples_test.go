package dicts_test

import (
	"fmt"
	"github.com/pickeringtech/go-collections/collections/dicts"
	"sync"
)

func ExampleNewConcurrentHash() {
	// Create a concurrent hash dictionary
	ch := dicts.NewConcurrentHash(
		dicts.Pair[string, int]{Key: "counter", Value: 0},
		dicts.Pair[string, int]{Key: "total", Value: 100},
	)

	fmt.Printf("Initial length: %d\n", ch.Length())
	
	value, found := ch.Get("counter", -1)
	fmt.Printf("Counter: %d, found: %t\n", value, found)

	// Output:
	// Initial length: 2
	// Counter: 0, found: true
}

func ExampleConcurrentHash_PutInPlace() {
	ch := dicts.NewConcurrentHash[string, int]()
	var wg sync.WaitGroup

	// Simulate concurrent writes
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ch.PutInPlace(fmt.Sprintf("key%d", id), id*10)
		}(i)
	}

	wg.Wait()
	fmt.Printf("Final length: %d\n", ch.Length())

	// Output:
	// Final length: 3
}

func ExampleConcurrentHash_ForEach() {
	ch := dicts.NewConcurrentHash(
		dicts.Pair[string, int]{Key: "apple", Value: 5},
		dicts.Pair[string, int]{Key: "banana", Value: 3},
		dicts.Pair[string, int]{Key: "cherry", Value: 8},
	)

	// Safe concurrent iteration
	ch.ForEach(func(key string, value int) {
		fmt.Printf("%s: %d\n", key, value)
	})

	// Unordered output:
	// apple: 5
	// banana: 3
	// cherry: 8
}

func ExampleNewConcurrentHashRW() {
	// Create a concurrent hash dictionary with read-write mutex
	chrw := dicts.NewConcurrentHashRW(
		dicts.Pair[string, int]{Key: "readers", Value: 10},
		dicts.Pair[string, int]{Key: "writers", Value: 2},
	)

	fmt.Printf("Initial length: %d\n", chrw.Length())
	
	// Multiple concurrent reads are efficient
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			value, _ := chrw.Get("readers", 0)
			fmt.Printf("Read value: %d\n", value)
		}()
	}

	wg.Wait()

	// Output:
	// Initial length: 2
	// Read value: 10
	// Read value: 10
	// Read value: 10
	// Read value: 10
	// Read value: 10
}

func ExampleConcurrentHashRW_Filter() {
	chrw := dicts.NewConcurrentHashRW(
		dicts.Pair[string, int]{Key: "small", Value: 5},
		dicts.Pair[string, int]{Key: "medium", Value: 15},
		dicts.Pair[string, int]{Key: "large", Value: 25},
	)

	// Filter returns a new regular Hash (not concurrent)
	filtered := chrw.Filter(func(key string, value int) bool {
		return value > 10
	})

	fmt.Printf("Original length: %d\n", chrw.Length())
	fmt.Printf("Filtered length: %d\n", filtered.Length())

	// Output:
	// Original length: 3
	// Filtered length: 2
}

func ExampleConcurrentHash_RemoveInPlace() {
	ch := dicts.NewConcurrentHash(
		dicts.Pair[string, int]{Key: "temp1", Value: 1},
		dicts.Pair[string, int]{Key: "temp2", Value: 2},
		dicts.Pair[string, int]{Key: "keep", Value: 3},
	)

	// Safe concurrent removal
	value, found := ch.RemoveInPlace("temp1")
	fmt.Printf("Removed value: %d, found: %t\n", value, found)
	fmt.Printf("Length after removal: %d\n", ch.Length())

	// Try to remove non-existing key
	value, found = ch.RemoveInPlace("nonexistent")
	fmt.Printf("Removed value: %d, found: %t\n", value, found)

	// Output:
	// Removed value: 1, found: true
	// Length after removal: 2
	// Removed value: 0, found: false
}

func ExampleConcurrentHash_Keys() {
	ch := dicts.NewConcurrentHash(
		dicts.Pair[int, string]{Key: 3, Value: "three"},
		dicts.Pair[int, string]{Key: 1, Value: "one"},
		dicts.Pair[int, string]{Key: 2, Value: "two"},
	)

	keys := ch.Keys()
	fmt.Printf("Number of keys: %d\n", len(keys))

	// Note: Hash iteration order is not guaranteed
	// Output:
	// Number of keys: 3
}

func ExampleConcurrentHashRW_AsMap() {
	chrw := dicts.NewConcurrentHashRW(
		dicts.Pair[string, int]{Key: "x", Value: 10},
		dicts.Pair[string, int]{Key: "y", Value: 20},
	)

	// Convert to native Go map
	nativeMap := chrw.AsMap()
	fmt.Printf("Native map length: %d\n", len(nativeMap))
	fmt.Printf("x value: %d\n", nativeMap["x"])

	// Output:
	// Native map length: 2
	// x value: 10
}

// Example showing the difference between concurrent implementations
func Example_concurrentComparison() {
	// ConcurrentHash - uses sync.Mutex (exclusive access for all operations)
	ch := dicts.NewConcurrentHash(
		dicts.Pair[string, int]{Key: "data", Value: 42},
	)

	// ConcurrentHashRW - uses sync.RWMutex (concurrent reads, exclusive writes)
	chrw := dicts.NewConcurrentHashRW(
		dicts.Pair[string, int]{Key: "data", Value: 42},
	)

	// Both are safe for concurrent access
	fmt.Printf("ConcurrentHash value: %d\n", func() int {
		v, _ := ch.Get("data", 0)
		return v
	}())

	fmt.Printf("ConcurrentHashRW value: %d\n", func() int {
		v, _ := chrw.Get("data", 0)
		return v
	}())

	// Output:
	// ConcurrentHash value: 42
	// ConcurrentHashRW value: 42
}
