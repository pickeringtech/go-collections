package sets_test

import (
	"fmt"
	"github.com/pickeringtech/go-collections/collections/sets"
	"sort"
)

func ExampleHash_IsSubsetOf() {
	s1 := sets.NewHash(1, 2)
	s2 := sets.NewHash(1, 2, 3, 4)

	fmt.Printf("s1 is subset of s2: %v\n", s1.IsSubsetOf(s2))
	fmt.Printf("s2 is subset of s1: %v\n", s2.IsSubsetOf(s1))

	// Output:
	// s1 is subset of s2: true
	// s2 is subset of s1: false
}

func ExampleHash_IsSupersetOf() {
	s1 := sets.NewHash(1, 2, 3, 4)
	s2 := sets.NewHash(1, 2)

	fmt.Printf("s1 is superset of s2: %v\n", s1.IsSupersetOf(s2))
	fmt.Printf("s2 is superset of s1: %v\n", s2.IsSupersetOf(s1))

	// Output:
	// s1 is superset of s2: true
	// s2 is superset of s1: false
}

func ExampleHash_IsDisjoint() {
	s1 := sets.NewHash(1, 2, 3)
	s2 := sets.NewHash(4, 5, 6)
	s3 := sets.NewHash(3, 4, 5)

	fmt.Printf("s1 and s2 are disjoint: %v\n", s1.IsDisjoint(s2))
	fmt.Printf("s1 and s3 are disjoint: %v\n", s1.IsDisjoint(s3))

	// Output:
	// s1 and s2 are disjoint: true
	// s1 and s3 are disjoint: false
}

func ExampleHash_Equals() {
	s1 := sets.NewHash(1, 2, 3)
	s2 := sets.NewHash(3, 2, 1) // Order doesn't matter in sets
	s3 := sets.NewHash(1, 2, 4)

	fmt.Printf("s1 equals s2: %v\n", s1.Equals(s2))
	fmt.Printf("s1 equals s3: %v\n", s1.Equals(s3))

	// Output:
	// s1 equals s2: true
	// s1 equals s3: false
}

func ExampleHash_Find() {
	s := sets.NewHash(1, 2, 3, 4, 5)

	// Find first even number
	element, found := s.Find(func(n int) bool {
		return n%2 == 0
	})

	if found {
		fmt.Printf("Found even number: %d\n", element)
	} else {
		fmt.Println("No even number found")
	}

	// Output:
	// Found even number: 2
}

func ExampleHash_AllMatch() {
	s1 := sets.NewHash(2, 4, 6, 8)
	s2 := sets.NewHash(1, 2, 3, 4)

	allEven1 := s1.AllMatch(func(n int) bool {
		return n%2 == 0
	})

	allEven2 := s2.AllMatch(func(n int) bool {
		return n%2 == 0
	})

	fmt.Printf("All elements in s1 are even: %v\n", allEven1)
	fmt.Printf("All elements in s2 are even: %v\n", allEven2)

	// Output:
	// All elements in s1 are even: true
	// All elements in s2 are even: false
}

func ExampleHash_AnyMatch() {
	s := sets.NewHash(1, 3, 5, 7, 8)

	hasEven := s.AnyMatch(func(n int) bool {
		return n%2 == 0
	})

	fmt.Printf("Set has any even number: %v\n", hasEven)

	// Output:
	// Set has any even number: true
}

func ExampleHash_ForEach() {
	s := sets.NewHash("apple", "banana", "cherry")

	fmt.Println("Fruits in the set:")
	s.ForEach(func(fruit string) {
		fmt.Printf("- %s\n", fruit)
	})

	// Unordered output:
	// Fruits in the set:
	// - apple
	// - banana
	// - cherry
}

func ExampleHash_AsSlice() {
	s := sets.NewHash(3, 1, 4, 1, 5) // Note: duplicate 1 will be ignored

	slice := s.AsSlice()
	sort.Ints(slice) // Sort for consistent output

	fmt.Printf("Set as sorted slice: %v\n", slice)
	fmt.Printf("Length: %d\n", len(slice))

	// Output:
	// Set as sorted slice: [1 3 4 5]
	// Length: 4
}

func ExampleHash_AsMap() {
	s := sets.NewHash("red", "green", "blue")

	m := s.AsMap()
	fmt.Printf("Set as map: %v\n", m)
	fmt.Printf("Map length: %d\n", len(m))

	// Check if key exists
	_, exists := m["red"]
	fmt.Printf("'red' exists in map: %v\n", exists)

	// Output:
	// Set as map: map[blue:{} green:{} red:{}]
	// Map length: 3
	// 'red' exists in map: true
}

func ExampleHash_AddMany() {
	s := sets.NewHash("apple")

	newS := s.AddMany("banana", "cherry", "apple") // duplicate "apple" ignored

	fmt.Printf("Original length: %d\n", s.Length())
	fmt.Printf("New length: %d\n", newS.Length())

	// Output:
	// Original length: 1
	// New length: 3
}

func ExampleHash_RemoveMany() {
	s := sets.NewHash(1, 2, 3, 4, 5)

	newS := s.RemoveMany(2, 4, 6) // 6 doesn't exist, ignored

	fmt.Printf("Original length: %d\n", s.Length())
	fmt.Printf("New length: %d\n", newS.Length())
	fmt.Printf("New contains 2: %v\n", newS.Contains(2))
	fmt.Printf("New contains 3: %v\n", newS.Contains(3))

	// Output:
	// Original length: 5
	// New length: 3
	// New contains 2: false
	// New contains 3: true
}

func ExampleHash_FilterInPlace() {
	s := sets.NewHash(1, 2, 3, 4, 5, 6)

	fmt.Printf("Before filter: length = %d\n", s.Length())

	// Keep only even numbers
	s.FilterInPlace(func(n int) bool {
		return n%2 == 0
	})

	fmt.Printf("After filter: length = %d\n", s.Length())
	fmt.Printf("Contains 1: %v\n", s.Contains(1))
	fmt.Printf("Contains 2: %v\n", s.Contains(2))

	// Output:
	// Before filter: length = 6
	// After filter: length = 3
	// Contains 1: false
	// Contains 2: true
}

func ExampleHash_Clear() {
	s := sets.NewHash(1, 2, 3, 4, 5)

	fmt.Printf("Before clear: length = %d\n", s.Length())
	fmt.Printf("Before clear: empty = %v\n", s.IsEmpty())

	s.Clear()

	fmt.Printf("After clear: length = %d\n", s.Length())
	fmt.Printf("After clear: empty = %v\n", s.IsEmpty())

	// Output:
	// Before clear: length = 5
	// Before clear: empty = false
	// After clear: length = 0
	// After clear: empty = true
}

func Example_setOperations() {
	// Mathematical set operations example
	evens := sets.NewHash(2, 4, 6, 8)
	primes := sets.NewHash(2, 3, 5, 7)
	odds := sets.NewHash(1, 3, 5, 7, 9)

	// Union: all elements from both sets
	evenOrPrime := evens.Union(primes)
	fmt.Printf("Even or prime numbers: %d elements\n", evenOrPrime.Length())

	// Intersection: elements in both sets
	evenPrimes := evens.Intersection(primes)
	fmt.Printf("Even prime numbers: %d elements\n", evenPrimes.Length())

	// Difference: elements in first set but not in second
	evenNotPrime := evens.Difference(primes)
	fmt.Printf("Even non-prime numbers: %d elements\n", evenNotPrime.Length())

	// Set relationships
	fmt.Printf("Evens and odds are disjoint: %v\n", evens.IsDisjoint(odds))

	// Output:
	// Even or prime numbers: 7 elements
	// Even prime numbers: 1 elements
	// Even non-prime numbers: 3 elements
	// Evens and odds are disjoint: true
}

func Example_setComparison() {
	// Demonstrate different ways to compare sets
	s1 := sets.NewHash(1, 2, 3)
	s2 := sets.NewHash(1, 2)
	s3 := sets.NewHash(1, 2, 3)

	fmt.Printf("s1 equals s3: %v\n", s1.Equals(s3))
	fmt.Printf("s2 is subset of s1: %v\n", s2.IsSubsetOf(s1))
	fmt.Printf("s1 is superset of s2: %v\n", s1.IsSupersetOf(s2))

	// Empty set is subset of any set
	empty := sets.NewHash[int]()
	fmt.Printf("Empty set is subset of s1: %v\n", empty.IsSubsetOf(s1))

	// Output:
	// s1 equals s3: true
	// s2 is subset of s1: true
	// s1 is superset of s2: true
	// Empty set is subset of s1: true
}
