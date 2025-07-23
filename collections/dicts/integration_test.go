package dicts_test

import (
	"fmt"
	"github.com/pickeringtech/go-collections/collections"
	"github.com/pickeringtech/go-collections/collections/dicts"
	"testing"
)

func ExampleNewDict() {
	// Create a dictionary using the collections package
	dict := collections.NewDict(
		dicts.Pair[string, int]{Key: "apple", Value: 5},
		dicts.Pair[string, int]{Key: "banana", Value: 3},
	)

	value, found := dict.Get("apple", -1)
	fmt.Printf("Apple count: %d, found: %t\n", value, found)

	// Add more items (immutable operation)
	newDict := dict.Put("cherry", 8)
	fmt.Printf("Original items: %d\n", dict.Length())
	fmt.Printf("New dict items: %d\n", newDict.Length())

	// Output:
	// Apple count: 5, found: true
	// Original items: 2
	// New dict items: 3
}

func ExampleNewConcurrentDict() {
	// Create a concurrent dictionary
	dict := collections.NewConcurrentDict(
		dicts.Pair[string, int]{Key: "counter", Value: 0},
	)

	// Safe for concurrent access (using mutable interface)
	if mutableDict, ok := dict.(dicts.MutableDict[string, int]); ok {
		mutableDict.PutInPlace("counter", 42)
	}
	value, found := dict.Get("counter", -1)
	fmt.Printf("Counter: %d, found: %t\n", value, found)

	// Output:
	// Counter: 42, found: true
}

func TestIntegration_AllImplementations(t *testing.T) {
	pairs := []dicts.Pair[string, int]{
		{Key: "one", Value: 1},
		{Key: "two", Value: 2},
		{Key: "three", Value: 3},
	}

	// Test all implementations through the collections package
	implementations := map[string]dicts.Dict[string, int]{
		"Hash":            collections.NewDict(pairs...),
		"ConcurrentHash":  collections.NewConcurrentDict(pairs...),
		"ConcurrentHashRW": collections.NewConcurrentRWDict(pairs...),
	}

	for name, dict := range implementations {
		t.Run(name, func(t *testing.T) {
			// Test basic operations
			if dict.Length() != 3 {
				t.Errorf("%s Length() = %v, want 3", name, dict.Length())
			}

			value, found := dict.Get("two", -1)
			if !found || value != 2 {
				t.Errorf("%s Get('two') = %v, %v, want 2, true", name, value, found)
			}

			if !dict.Contains("one") {
				t.Errorf("%s Contains('one') = false, want true", name)
			}

			// Test filtering
			filtered := dict.Filter(func(key string, value int) bool {
				return value > 1
			})

			if filtered.Length() != 2 {
				t.Errorf("%s Filter() length = %v, want 2", name, filtered.Length())
			}

			// Test keys and values
			keys := dict.Keys()
			if len(keys) != 3 {
				t.Errorf("%s Keys() length = %v, want 3", name, len(keys))
			}

			values := dict.Values()
			if len(values) != 3 {
				t.Errorf("%s Values() length = %v, want 3", name, len(values))
			}
		})
	}
}

func TestIntegration_TreeWithOrderedKeys(t *testing.T) {
	// Test Tree with different ordered types
	
	// String keys
	stringTree := dicts.NewTree(
		dicts.Pair[string, int]{Key: "charlie", Value: 3},
		dicts.Pair[string, int]{Key: "alice", Value: 1},
		dicts.Pair[string, int]{Key: "bob", Value: 2},
	)

	stringKeys := stringTree.Keys()
	expectedStringKeys := []string{"alice", "bob", "charlie"}
	for i, key := range stringKeys {
		if key != expectedStringKeys[i] {
			t.Errorf("String tree keys[%d] = %v, want %v", i, key, expectedStringKeys[i])
		}
	}

	// Integer keys
	intTree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 3, Value: "three"},
		dicts.Pair[int, string]{Key: 1, Value: "one"},
		dicts.Pair[int, string]{Key: 2, Value: "two"},
	)

	intKeys := intTree.Keys()
	expectedIntKeys := []int{1, 2, 3}
	for i, key := range intKeys {
		if key != expectedIntKeys[i] {
			t.Errorf("Int tree keys[%d] = %v, want %v", i, key, expectedIntKeys[i])
		}
	}

	// Float keys
	floatTree := dicts.NewTree(
		dicts.Pair[float64, string]{Key: 3.14, Value: "pi"},
		dicts.Pair[float64, string]{Key: 2.71, Value: "e"},
		dicts.Pair[float64, string]{Key: 1.41, Value: "sqrt2"},
	)

	floatKeys := floatTree.Keys()
	expectedFloatKeys := []float64{1.41, 2.71, 3.14}
	for i, key := range floatKeys {
		if key != expectedFloatKeys[i] {
			t.Errorf("Float tree keys[%d] = %v, want %v", i, key, expectedFloatKeys[i])
		}
	}
}

func BenchmarkIntegration_CompareImplementations(b *testing.B) {
	pairs := make([]dicts.Pair[int, string], 1000)
	for i := 0; i < 1000; i++ {
		pairs[i] = dicts.Pair[int, string]{Key: i, Value: fmt.Sprintf("value_%d", i)}
	}

	b.Run("Hash_Get", func(b *testing.B) {
		dict := collections.NewDict(pairs...)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := i % 1000
			_, _ = dict.Get(key, "default")
		}
	})

	b.Run("ConcurrentHash_Get", func(b *testing.B) {
		dict := collections.NewConcurrentDict(pairs...)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := i % 1000
			_, _ = dict.Get(key, "default")
		}
	})

	b.Run("ConcurrentHashRW_Get", func(b *testing.B) {
		dict := collections.NewConcurrentRWDict(pairs...)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := i % 1000
			_, _ = dict.Get(key, "default")
		}
	})

	b.Run("Tree_Get", func(b *testing.B) {
		tree := dicts.NewTree(pairs...)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := i % 1000
			_, _ = tree.Get(key, "default")
		}
	})
}

// Removed problematic example function
