# Sets - Mathematical Collections

The `sets` package provides mathematical sets for Go. Use it for membership testing, removing duplicates, and set operations like union and intersection.

## Quick Start

```go
import "github.com/pickeringtech/go-collections/collections/sets"

// Create sets with unique elements
permissions := sets.NewHash("read", "write", "execute")
userPerms := sets.NewHash("read", "write")

// Mathematical operations
canExecute := permissions.Contains("execute")           // true
common := permissions.Intersection(userPerms)          // {read, write}
missing := permissions.Difference(userPerms)           // {execute}
isSubset := userPerms.IsSubsetOf(permissions)          // true
```

## Why Use Sets?

**Native Go maps for sets are clunky:**
```go
// Native approach - verbose and error-prone
permissions := map[string]struct{}{
    "read": {}, "write": {}, "execute": {},
}
if _, exists := permissions["read"]; exists {
    // Can read
}
// No built-in operations for union, intersection, etc.
```

**Sets are concise:**
```go
// Clean and expressive
permissions := sets.NewHash("read", "write", "execute")
if permissions.Contains("read") {
    // Can read
}

// Rich mathematical operations built-in
adminPerms := permissions.Union(sets.NewHash("admin", "delete"))
```

## Implementations

### Hash Set (`Hash[T]`)

A hash set implementation using Go's built-in map. Provides O(1) average case performance for basic operations.

```go
// Create a new hash set
s := sets.NewHash("apple", "banana", "cherry")

// Check membership
exists := s.Contains("apple")
fmt.Printf("Contains 'apple': %t\n", exists) // Contains 'apple': true

// Add elements (immutable)
newS := s.Add("date")
fmt.Printf("Original length: %d, New length: %d\n", s.Length(), newS.Length())

// Add elements (mutable)
s.AddInPlace("elderberry")
```

### Concurrent Hash Set (`ConcurrentHash[T]`)

A thread-safe hash set using a mutex for synchronization. All operations are atomic.

```go
cs := sets.NewConcurrentHash("apple", "banana")

// Safe to use from multiple goroutines
go func() {
    cs.AddInPlace("cherry")
}()

go func() {
    exists := cs.Contains("apple")
    fmt.Printf("Contains 'apple': %t\n", exists)
}()
```

### Concurrent RW Hash Set (`ConcurrentHashRW[T]`)

A thread-safe hash set using a read-write mutex. Read operations can proceed concurrently, while write operations are exclusive.

```go
crws := sets.NewConcurrentHashRW(1, 2, 3, 4, 5)

// Multiple readers can access concurrently
// Writers get exclusive access
```

## Interface Overview

### Core Operations

```go
// Basic access
exists := set.Contains(element)
length := set.Length()
isEmpty := set.IsEmpty()

// Iteration
set.ForEach(func(element T) {
    // Process each element
})
```

### Immutable Operations

```go
// Adding elements (returns new set)
newSet := set.Add(element)
newSet = set.AddMany(elem1, elem2, elem3)

// Removing elements (returns new set)
newSet = set.Remove(element)
newSet = set.RemoveMany(elem1, elem2)

// Filtering (returns new set)
filtered := set.Filter(func(element T) bool {
    return someCondition(element)
})
```

### Mutable Operations

```go
// Adding elements (modifies original)
set.AddInPlace(element)
set.AddManyInPlace(elem1, elem2, elem3)

// Removing elements (modifies original)
removed := set.RemoveInPlace(element)
set.RemoveManyInPlace(elem1, elem2)
set.Clear()

// Filtering (modifies original)
set.FilterInPlace(func(element T) bool {
    return someCondition(element)
})
```

### Mathematical Set Operations

```go
// Set operations (return new sets)
union := set1.Union(set2)
intersection := set1.Intersection(set2)
difference := set1.Difference(set2)

// Set relationships
isSubset := set1.IsSubsetOf(set2)
isSuperset := set1.IsSupersetOf(set2)
areDisjoint := set1.IsDisjoint(set2)
areEqual := set1.Equals(set2)

// In-place operations
set1.UnionInPlace(set2)
set1.IntersectionInPlace(set2)
set1.DifferenceInPlace(set2)
```

### Search Operations

```go
// Find first matching element
element, found := set.Find(func(e T) bool {
    return someCondition(e)
})

// Check if all/any elements match
allMatch := set.AllMatch(func(e T) bool {
    return someCondition(e)
})

anyMatch := set.AnyMatch(func(e T) bool {
    return someCondition(e)
})

noneMatch := set.NoneMatch(func(e T) bool {
    return someCondition(e)
})
```

`Find`, `AllMatch`, `AnyMatch` and `NoneMatch` form the search core shared
across the `lists`, `dicts` and `sets` families.

### Conversion Operations

```go
// Convert to slice
slice := set.AsSlice()

// Convert to native Go map
nativeMap := set.AsMap()
```

### Transforming to a New Type: Map / Reduce

`Filter` is a method because it keeps the same element type. A general `Map` is
`T -> U` with a **different** element type, which Go methods cannot express
([golang/go#49085](https://github.com/golang/go/issues/49085)), so `Map` and
`Reduce` are **free functions** over the `Set` interface. `Map` returns the
`Set` interface (backed by `NewHash`) so results chain on.

```go
words := sets.NewHash("a", "bb", "cc")

// Map: T -> U. Because a Set deduplicates, distinct inputs that map to the
// same value collapse — here "bb" and "cc" both -> 2, so the result is {1, 2}.
lengths := sets.Map(words, func(s string) int { return len(s) })

// Reduce: fold every element into a single value
total := sets.Reduce(words, 0, func(acc int, s string) int { return acc + len(s) })
```

**`Map` always returns a `Hash`-backed `Set`, even for a sorted (`TreeSet`)
input.** This is deliberate: `Map` may change the element type, and the output
element type is only constrained to `comparable`, not `Ordered` — so a sorted
output cannot be guaranteed in general. When your output element type *is*
`Ordered` and you want to keep sorted iteration, use `MapSorted`, which returns a
`TreeSet`-backed `SortedSet`:

```go
// MapSorted: like Map, but U must be Ordered and the result stays sorted.
lengths := sets.MapSorted(words, func(s string) int { return len(s) })  // SortedSet[int]
// lengths.AsSlice() == []int{1, 2} (ascending, deduplicated)
```

Iteration order over a `Set` is unspecified, so a reduction should be
order-independent.

## Performance Characteristics

| Implementation | Contains | Add | Remove | Memory | Thread-Safe |
|---------------|----------|-----|--------|---------|-------------|
| Hash | O(1) | O(1) | O(1) | Low | No |
| ConcurrentHash | O(1) | O(1) | O(1) | Low | Yes |
| ConcurrentHashRW | O(1) | O(1) | O(1) | Low | Yes |

## Usage Examples

### Basic Set Operations

```go
package main

import (
    "fmt"
    "github.com/pickeringtech/go-collections/collections/sets"
)

func main() {
    // Create sets
    evens := sets.NewHash(2, 4, 6, 8)
    primes := sets.NewHash(2, 3, 5, 7)

    // Mathematical operations
    union := evens.Union(primes)
    intersection := evens.Intersection(primes)
    difference := evens.Difference(primes)

    fmt.Printf("Union: %d elements\n", union.Length())
    fmt.Printf("Intersection: %d elements\n", intersection.Length())
    fmt.Printf("Difference: %d elements\n", difference.Length())
}
```


## How Do They Work?

## When Should I Use Them?

## Implementations

### Hash Set

## Interfaces
