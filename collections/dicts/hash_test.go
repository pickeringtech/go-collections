package dicts_test

import (
	"fmt"
	"github.com/pickeringtech/go-collections/collections/dicts"
	"reflect"
	"testing"
)

func ExampleHash_Get() {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "one", Value: 1},
		dicts.Pair[string, int]{Key: "two", Value: 2},
		dicts.Pair[string, int]{Key: "three", Value: 3},
	)

	value, found := h.Get("two", -1)
	fmt.Printf("Value: %v, Found: %v\n", value, found)

	value, found = h.Get("four", -1)
	fmt.Printf("Value: %v, Found: %v\n", value, found)

	// Output:
	// Value: 2, Found: true
	// Value: -1, Found: false
}

func TestHash_Get(t *testing.T) {
	type testCase[K comparable, V any] struct {
		name         string
		hash         dicts.Hash[K, V]
		key          K
		defaultValue V
		wantValue    V
		wantFound    bool
	}

	tests := []testCase[string, int]{
		{
			name: "existing key returns value and true",
			hash: dicts.NewHash(
				dicts.Pair[string, int]{Key: "one", Value: 1},
				dicts.Pair[string, int]{Key: "two", Value: 2},
			),
			key:          "one",
			defaultValue: -1,
			wantValue:    1,
			wantFound:    true,
		},
		{
			name: "non-existing key returns default and false",
			hash: dicts.NewHash(
				dicts.Pair[string, int]{Key: "one", Value: 1},
			),
			key:          "two",
			defaultValue: -1,
			wantValue:    -1,
			wantFound:    false,
		},
		{
			name:         "empty hash returns default and false",
			hash:         dicts.NewHash[string, int](),
			key:          "one",
			defaultValue: -1,
			wantValue:    -1,
			wantFound:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotFound := tt.hash.Get(tt.key, tt.defaultValue)
			if gotValue != tt.wantValue {
				t.Errorf("Get() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotFound != tt.wantFound {
				t.Errorf("Get() gotFound = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}

func ExampleHash_Contains() {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "one", Value: 1},
		dicts.Pair[string, int]{Key: "two", Value: 2},
	)

	fmt.Printf("Contains 'one': %v\n", h.Contains("one"))
	fmt.Printf("Contains 'three': %v\n", h.Contains("three"))

	// Output:
	// Contains 'one': true
	// Contains 'three': false
}

func TestHash_Contains(t *testing.T) {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "one", Value: 1},
		dicts.Pair[string, int]{Key: "two", Value: 2},
	)

	if !h.Contains("one") {
		t.Error("Contains() should return true for existing key")
	}
	if h.Contains("three") {
		t.Error("Contains() should return false for non-existing key")
	}
}

func ExampleHash_Length() {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "one", Value: 1},
		dicts.Pair[string, int]{Key: "two", Value: 2},
		dicts.Pair[string, int]{Key: "three", Value: 3},
	)

	fmt.Printf("Length: %v\n", h.Length())

	// Output:
	// Length: 3
}

func TestHash_Length(t *testing.T) {
	tests := []struct {
		name string
		hash dicts.Hash[string, int]
		want int
	}{
		{
			name: "empty hash has length 0",
			hash: dicts.NewHash[string, int](),
			want: 0,
		},
		{
			name: "hash with elements has correct length",
			hash: dicts.NewHash(
				dicts.Pair[string, int]{Key: "one", Value: 1},
				dicts.Pair[string, int]{Key: "two", Value: 2},
				dicts.Pair[string, int]{Key: "three", Value: 3},
			),
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hash.Length(); got != tt.want {
				t.Errorf("Length() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleHash_IsEmpty() {
	empty := dicts.NewHash[string, int]()
	notEmpty := dicts.NewHash(dicts.Pair[string, int]{Key: "one", Value: 1})

	fmt.Printf("Empty hash is empty: %v\n", empty.IsEmpty())
	fmt.Printf("Non-empty hash is empty: %v\n", notEmpty.IsEmpty())

	// Output:
	// Empty hash is empty: true
	// Non-empty hash is empty: false
}

func TestHash_IsEmpty(t *testing.T) {
	empty := dicts.NewHash[string, int]()
	notEmpty := dicts.NewHash(dicts.Pair[string, int]{Key: "one", Value: 1})

	if !empty.IsEmpty() {
		t.Error("IsEmpty() should return true for empty hash")
	}
	if notEmpty.IsEmpty() {
		t.Error("IsEmpty() should return false for non-empty hash")
	}
}

func ExampleHash_ForEach() {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "one", Value: 1},
		dicts.Pair[string, int]{Key: "two", Value: 2},
	)

	h.ForEach(func(key string, value int) {
		fmt.Printf("%s: %d\n", key, value)
	})

	// Unordered output:
	// one: 1
	// two: 2
}

func TestHash_ForEach(t *testing.T) {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "one", Value: 1},
		dicts.Pair[string, int]{Key: "two", Value: 2},
	)

	visited := make(map[string]int)
	h.ForEach(func(key string, value int) {
		visited[key] = value
	})

	expected := map[string]int{"one": 1, "two": 2}
	if !reflect.DeepEqual(visited, expected) {
		t.Errorf("ForEach() visited = %v, want %v", visited, expected)
	}
}

func ExampleHash_Keys() {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "one", Value: 1},
		dicts.Pair[string, int]{Key: "two", Value: 2},
	)

	keys := h.Keys()
	fmt.Printf("Number of keys: %d\n", len(keys))

	// Output:
	// Number of keys: 2
}

func TestHash_Keys(t *testing.T) {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "one", Value: 1},
		dicts.Pair[string, int]{Key: "two", Value: 2},
	)

	keys := h.Keys()
	if len(keys) != 2 {
		t.Errorf("Keys() length = %v, want 2", len(keys))
	}

	keySet := make(map[string]bool)
	for _, key := range keys {
		keySet[key] = true
	}

	if !keySet["one"] || !keySet["two"] {
		t.Errorf("Keys() = %v, should contain 'one' and 'two'", keys)
	}
}

func ExampleHash_Values() {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "one", Value: 1},
		dicts.Pair[string, int]{Key: "two", Value: 2},
	)

	values := h.Values()
	fmt.Printf("Number of values: %d\n", len(values))

	// Output:
	// Number of values: 2
}

func TestHash_Values(t *testing.T) {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "one", Value: 1},
		dicts.Pair[string, int]{Key: "two", Value: 2},
	)

	values := h.Values()
	if len(values) != 2 {
		t.Errorf("Values() length = %v, want 2", len(values))
	}

	valueSet := make(map[int]bool)
	for _, value := range values {
		valueSet[value] = true
	}

	if !valueSet[1] || !valueSet[2] {
		t.Errorf("Values() = %v, should contain 1 and 2", values)
	}
}

func ExampleHash_Put() {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "one", Value: 1},
	)

	newH := h.Put("two", 2)
	
	fmt.Printf("Original length: %d\n", h.Length())
	fmt.Printf("New length: %d\n", newH.Length())

	// Output:
	// Original length: 1
	// New length: 2
}

func TestHash_Put(t *testing.T) {
	original := dicts.NewHash(
		dicts.Pair[string, int]{Key: "one", Value: 1},
	)

	newDict := original.Put("two", 2)

	// Original should be unchanged
	if original.Length() != 1 {
		t.Errorf("Put() modified original dict, length = %v, want 1", original.Length())
	}

	// New dict should have the new entry
	if newDict.Length() != 2 {
		t.Errorf("Put() new dict length = %v, want 2", newDict.Length())
	}

	if !newDict.Contains("two") {
		t.Error("Put() new dict should contain the new key")
	}

	value, found := newDict.Get("two", -1)
	if !found || value != 2 {
		t.Errorf("Put() new dict Get('two') = %v, %v, want 2, true", value, found)
	}
}

func ExampleHash_PutInPlace() {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "one", Value: 1},
	)

	fmt.Printf("Before: length = %d\n", h.Length())
	h.PutInPlace("two", 2)
	fmt.Printf("After: length = %d\n", h.Length())

	// Output:
	// Before: length = 1
	// After: length = 2
}

func ExampleHash_Remove() {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "one", Value: 1},
		dicts.Pair[string, int]{Key: "two", Value: 2},
	)

	newH := h.Remove("one")

	fmt.Printf("Original length: %d\n", h.Length())
	fmt.Printf("New length: %d\n", newH.Length())
	fmt.Printf("New contains 'one': %v\n", newH.Contains("one"))

	// Output:
	// Original length: 2
	// New length: 1
	// New contains 'one': false
}

func ExampleHash_Filter() {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "one", Value: 1},
		dicts.Pair[string, int]{Key: "two", Value: 2},
		dicts.Pair[string, int]{Key: "three", Value: 3},
		dicts.Pair[string, int]{Key: "four", Value: 4},
	)

	filtered := h.Filter(func(key string, value int) bool {
		return value%2 == 0 // Keep only even values
	})

	fmt.Printf("Original length: %d\n", h.Length())
	fmt.Printf("Filtered length: %d\n", filtered.Length())

	// Output:
	// Original length: 4
	// Filtered length: 2
}

func ExampleHash_Find() {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "one", Value: 1},
		dicts.Pair[string, int]{Key: "two", Value: 2},
		dicts.Pair[string, int]{Key: "three", Value: 3},
	)

	key, value, found := h.Find(func(k string, v int) bool {
		return v > 2
	})

	if found {
		fmt.Printf("Found: %s = %d\n", key, value)
	} else {
		fmt.Println("Not found")
	}

	// Output:
	// Found: three = 3
}

func ExampleHash_Items() {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "apple", Value: 5},
		dicts.Pair[string, int]{Key: "banana", Value: 3},
	)

	items := h.Items()
	fmt.Printf("Number of items: %d\n", len(items))

	for _, item := range items {
		fmt.Printf("%s: %d\n", item.Key, item.Value)
	}

	// Unordered output:
	// Number of items: 2
	// apple: 5
	// banana: 3
}

func ExampleNewHash() {
	// Create an empty hash
	empty := dicts.NewHash[string, int]()
	fmt.Printf("Empty hash length: %d\n", empty.Length())

	// Create a hash with initial data
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "red", Value: 1},
		dicts.Pair[string, int]{Key: "green", Value: 2},
		dicts.Pair[string, int]{Key: "blue", Value: 3},
	)
	fmt.Printf("Initialized hash length: %d\n", h.Length())

	// Output:
	// Empty hash length: 0
	// Initialized hash length: 3
}
