package dicts_test

import (
	"fmt"
	"github.com/pickeringtech/go-collections/collections/dicts"
	"reflect"
	"testing"
)

func ExampleTree_ForEach() {
	tree := dicts.NewTree(
		dicts.Pair[string, int]{Key: "charlie", Value: 3},
		dicts.Pair[string, int]{Key: "alice", Value: 1},
		dicts.Pair[string, int]{Key: "bob", Value: 2},
	)

	tree.ForEach(func(key string, value int) {
		fmt.Printf("%s: %d\n", key, value)
	})

	// Output:
	// alice: 1
	// bob: 2
	// charlie: 3
}

func TestTree_Get(t *testing.T) {
	tree := dicts.NewTree(
		dicts.Pair[string, int]{Key: "one", Value: 1},
		dicts.Pair[string, int]{Key: "two", Value: 2},
	)

	// Test existing key
	value, found := tree.Get("one", -1)
	if !found || value != 1 {
		t.Errorf("Get() = %v, %v, want 1, true", value, found)
	}

	// Test non-existing key
	value, found = tree.Get("three", -1)
	if found || value != -1 {
		t.Errorf("Get() = %v, %v, want -1, false", value, found)
	}
}

func TestTree_PutInPlace(t *testing.T) {
	tree := dicts.NewTree[string, int]()

	tree.PutInPlace("two", 2)
	tree.PutInPlace("one", 1)
	tree.PutInPlace("three", 3)

	if tree.Length() != 3 {
		t.Errorf("Length() = %v, want 3", tree.Length())
	}

	// Test that keys are accessible
	for _, key := range []string{"one", "two", "three"} {
		if !tree.Contains(key) {
			t.Errorf("Contains(%s) = false, want true", key)
		}
	}
}

func TestTree_RemoveInPlace(t *testing.T) {
	tree := dicts.NewTree(
		dicts.Pair[string, int]{Key: "one", Value: 1},
		dicts.Pair[string, int]{Key: "two", Value: 2},
		dicts.Pair[string, int]{Key: "three", Value: 3},
	)

	// Remove existing key
	value, found := tree.RemoveInPlace("two")
	if !found || value != 2 {
		t.Errorf("RemoveInPlace() = %v, %v, want 2, true", value, found)
	}

	if tree.Length() != 2 {
		t.Errorf("Length() after remove = %v, want 2", tree.Length())
	}

	if tree.Contains("two") {
		t.Error("Contains('two') = true, want false after removal")
	}

	// Remove non-existing key
	value, found = tree.RemoveInPlace("four")
	if found || value != 0 {
		t.Errorf("RemoveInPlace() non-existing = %v, %v, want 0, false", value, found)
	}
}

func TestTree_Keys_Sorted(t *testing.T) {
	tree := dicts.NewTree(
		dicts.Pair[string, int]{Key: "charlie", Value: 3},
		dicts.Pair[string, int]{Key: "alice", Value: 1},
		dicts.Pair[string, int]{Key: "bob", Value: 2},
		dicts.Pair[string, int]{Key: "david", Value: 4},
	)

	keys := tree.Keys()
	expected := []string{"alice", "bob", "charlie", "david"}

	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("Keys() = %v, want %v", keys, expected)
	}
}

func TestTree_Filter(t *testing.T) {
	tree := dicts.NewTree(
		dicts.Pair[string, int]{Key: "one", Value: 1},
		dicts.Pair[string, int]{Key: "two", Value: 2},
		dicts.Pair[string, int]{Key: "three", Value: 3},
		dicts.Pair[string, int]{Key: "four", Value: 4},
	)

	filtered := tree.Filter(func(key string, value int) bool {
		return value%2 == 0
	})

	if filtered.Length() != 2 {
		t.Errorf("Filter() length = %v, want 2", filtered.Length())
	}

	if !filtered.Contains("two") || !filtered.Contains("four") {
		t.Error("Filter() should contain keys with even values")
	}

	if filtered.Contains("one") || filtered.Contains("three") {
		t.Error("Filter() should not contain keys with odd values")
	}
}

func TestTree_IntegerKeys(t *testing.T) {
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 3, Value: "three"},
		dicts.Pair[int, string]{Key: 1, Value: "one"},
		dicts.Pair[int, string]{Key: 4, Value: "four"},
		dicts.Pair[int, string]{Key: 2, Value: "two"},
	)

	keys := tree.Keys()
	expected := []int{1, 2, 3, 4}

	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("Keys() = %v, want %v", keys, expected)
	}
}

func TestTree_Clear(t *testing.T) {
	tree := dicts.NewTree(
		dicts.Pair[string, int]{Key: "one", Value: 1},
		dicts.Pair[string, int]{Key: "two", Value: 2},
	)

	tree.Clear()

	if !tree.IsEmpty() {
		t.Error("IsEmpty() = false, want true after Clear()")
	}

	if tree.Length() != 0 {
		t.Errorf("Length() = %v, want 0 after Clear()", tree.Length())
	}
}

func ExampleTree_Keys() {
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 3, Value: "three"},
		dicts.Pair[int, string]{Key: 1, Value: "one"},
		dicts.Pair[int, string]{Key: 2, Value: "two"},
	)

	keys := tree.Keys()
	fmt.Printf("Sorted keys: %v\n", keys)

	// Output:
	// Sorted keys: [1 2 3]
}

func BenchmarkTree_Get(b *testing.B) {
	// Setup
	pairs := make([]dicts.Pair[int, string], 1000)
	for i := 0; i < 1000; i++ {
		pairs[i] = dicts.Pair[int, string]{Key: i, Value: fmt.Sprintf("value_%d", i)}
	}
	tree := dicts.NewTree(pairs...)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := i % 1000
		_, _ = tree.Get(key, "default")
	}
}

func BenchmarkTree_PutInPlace(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		tree := dicts.NewTree[int, string]()
		b.StartTimer()

		tree.PutInPlace(i, fmt.Sprintf("value_%d", i))
	}
}

func ExampleNewTree() {
	// Create an empty tree
	empty := dicts.NewTree[string, int]()
	fmt.Printf("Empty tree length: %d\n", empty.Length())

	// Create a tree with initial data
	tree := dicts.NewTree(
		dicts.Pair[string, int]{Key: "charlie", Value: 3},
		dicts.Pair[string, int]{Key: "alice", Value: 1},
		dicts.Pair[string, int]{Key: "bob", Value: 2},
	)
	fmt.Printf("Initialized tree length: %d\n", tree.Length())

	// Output:
	// Empty tree length: 0
	// Initialized tree length: 3
}

func ExampleTree_Get() {
	tree := dicts.NewTree(
		dicts.Pair[string, int]{Key: "apple", Value: 5},
		dicts.Pair[string, int]{Key: "banana", Value: 3},
	)

	value, found := tree.Get("apple", -1)
	fmt.Printf("Apple: %d, found: %t\n", value, found)

	value, found = tree.Get("cherry", -1)
	fmt.Printf("Cherry: %d, found: %t\n", value, found)

	// Output:
	// Apple: 5, found: true
	// Cherry: -1, found: false
}

func ExampleTree_PutInPlace() {
	tree := dicts.NewTree[int, string]()

	tree.PutInPlace(3, "three")
	tree.PutInPlace(1, "one")
	tree.PutInPlace(2, "two")

	fmt.Printf("Length: %d\n", tree.Length())

	// Keys are maintained in sorted order
	keys := tree.Keys()
	fmt.Printf("Sorted keys: %v\n", keys)

	// Output:
	// Length: 3
	// Sorted keys: [1 2 3]
}

func ExampleTree_Remove() {
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 1, Value: "one"},
		dicts.Pair[int, string]{Key: 2, Value: "two"},
		dicts.Pair[int, string]{Key: 3, Value: "three"},
	)

	newTree := tree.Remove(2)

	fmt.Printf("Original length: %d\n", tree.Length())
	fmt.Printf("New tree length: %d\n", newTree.Length())
	fmt.Printf("New tree keys: %v\n", newTree.Keys())

	// Output:
	// Original length: 3
	// New tree length: 2
	// New tree keys: [1 3]
}

func ExampleTree_Filter() {
	tree := dicts.NewTree(
		dicts.Pair[string, int]{Key: "a", Value: 1},
		dicts.Pair[string, int]{Key: "b", Value: 2},
		dicts.Pair[string, int]{Key: "c", Value: 3},
		dicts.Pair[string, int]{Key: "d", Value: 4},
	)

	// Filter for even values
	filtered := tree.Filter(func(key string, value int) bool {
		return value%2 == 0
	})

	fmt.Printf("Original keys: %v\n", tree.Keys())
	fmt.Printf("Filtered keys: %v\n", filtered.Keys())

	// Output:
	// Original keys: [a b c d]
	// Filtered keys: [b d]
}

func ExampleTree_Find() {
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 10, Value: "ten"},
		dicts.Pair[int, string]{Key: 5, Value: "five"},
		dicts.Pair[int, string]{Key: 15, Value: "fifteen"},
	)

	// Find first key greater than 7
	key, value, found := tree.Find(func(k int, v string) bool {
		return k > 7
	})

	if found {
		fmt.Printf("Found: %d = %s\n", key, value)
	}

	// Output:
	// Found: 10 = ten
}

func ExampleTree_Values() {
	tree := dicts.NewTree(
		dicts.Pair[string, int]{Key: "z", Value: 26},
		dicts.Pair[string, int]{Key: "a", Value: 1},
		dicts.Pair[string, int]{Key: "m", Value: 13},
	)

	// Values are returned in key-sorted order
	values := tree.Values()
	fmt.Printf("Values in key order: %v\n", values)

	// Output:
	// Values in key order: [1 13 26]
}

func ExampleTree_Items() {
	tree := dicts.NewTree(
		dicts.Pair[float64, string]{Key: 3.14, Value: "pi"},
		dicts.Pair[float64, string]{Key: 2.71, Value: "e"},
		dicts.Pair[float64, string]{Key: 1.41, Value: "sqrt2"},
	)

	// Items are returned in key-sorted order
	items := tree.Items()
	for _, item := range items {
		fmt.Printf("%.2f: %s\n", item.Key, item.Value)
	}

	// Output:
	// 1.41: sqrt2
	// 2.71: e
	// 3.14: pi
}
