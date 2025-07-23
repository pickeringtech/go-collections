# Dicts - Associative Arrays

The `dicts` package provides comprehensive dictionary/map implementations with a focus on performance, thread-safety, and rich functionality. It offers both immutable and mutable interfaces, allowing you to choose the right approach for your use case.

## Features

- **Multiple Implementations**: Hash-based and Tree-based dictionaries
- **Thread-Safe Options**: Concurrent implementations with mutex and read-write mutex
- **Rich Interface**: Comprehensive set of operations including filtering, searching, and iteration
- **Immutable & Mutable**: Choose between functional-style immutable operations or in-place mutations
- **Type-Safe**: Full generic type support for keys and values
- **Performance Optimized**: Benchmarked against native Go maps

## Implementations

### Hash Dictionary (`Hash[K, V]`)

A hash table implementation using Go's built-in map. Provides O(1) average case performance for basic operations.

```go
// Create a new hash dictionary
h := dicts.NewHash(
    dicts.Pair[string, int]{Key: "one", Value: 1},
    dicts.Pair[string, int]{Key: "two", Value: 2},
)

// Get a value
value, found := h.Get("one", -1)
fmt.Printf("Value: %d, Found: %t\n", value, found) // Value: 1, Found: true

// Add a new key-value pair (immutable)
newH := h.Put("three", 3)
fmt.Printf("Original length: %d, New length: %d\n", h.Length(), newH.Length())

// Add a new key-value pair (mutable)
h.PutInPlace("four", 4)
```

### Concurrent Hash Dictionary (`ConcurrentHash[K, V]`)

A thread-safe hash dictionary using a mutex for synchronization. All operations are atomic.

```go
ch := dicts.NewConcurrentHash(
    dicts.Pair[string, int]{Key: "one", Value: 1},
)

// Safe to use from multiple goroutines
go func() {
    ch.PutInPlace("two", 2)
}()

go func() {
    value, found := ch.Get("one", -1)
    fmt.Printf("Value: %d\n", value)
}()
```

### Concurrent RW Hash Dictionary (`ConcurrentHashRW[K, V]`)

A thread-safe hash dictionary using a read-write mutex. Read operations can proceed concurrently, while write operations are exclusive.

```go
chrw := dicts.NewConcurrentHashRW(
    dicts.Pair[string, int]{Key: "one", Value: 1},
)

// Multiple readers can access concurrently
// Writers get exclusive access
```

### Tree Dictionary (`Tree[K, V]`)

A binary search tree implementation that maintains keys in sorted order. Provides O(log n) average case performance and ordered iteration.

```go
tree := dicts.NewTree(
    dicts.Pair[string, int]{Key: "charlie", Value: 3},
    dicts.Pair[string, int]{Key: "alice", Value: 1},
    dicts.Pair[string, int]{Key: "bob", Value: 2},
)

// Iterate in sorted order
tree.ForEach(func(key string, value int) {
    fmt.Printf("%s: %d\n", key, value)
})
// Output:
// alice: 1
// bob: 2
// charlie: 3
```

## Interface Overview

### Core Operations

```go
// Basic access
value, found := dict.Get("key", defaultValue)
exists := dict.Contains("key")
length := dict.Length()
isEmpty := dict.IsEmpty()

// Iteration
dict.ForEach(func(key K, value V) {
    // Process each key-value pair
})
dict.ForEachKey(func(key K) {
    // Process each key
})
dict.ForEachValue(func(value V) {
    // Process each value
})
```

### Immutable Operations

```go
// Adding/updating (returns new dictionary)
newDict := dict.Put("key", value)
newDict = dict.PutMany(
    dicts.Pair[K, V]{Key: "key1", Value: value1},
    dicts.Pair[K, V]{Key: "key2", Value: value2},
)

// Removing (returns new dictionary)
newDict = dict.Remove("key")
newDict = dict.RemoveMany("key1", "key2", "key3")

// Filtering (returns new dictionary)
filtered := dict.Filter(func(key K, value V) bool {
    return value > 10
})
```

### Mutable Operations

```go
// Adding/updating (modifies original)
dict.PutInPlace("key", value)
dict.PutManyInPlace(pairs...)

// Removing (modifies original)
removedValue, found := dict.RemoveInPlace("key")
dict.RemoveManyInPlace("key1", "key2")
dict.Clear()

// Filtering (modifies original)
dict.FilterInPlace(func(key K, value V) bool {
    return value > 10
})
```

### Search Operations

```go
// Find first matching pair
key, value, found := dict.Find(func(k K, v V) bool {
    return v > 100
})

// Find first matching key
key, found := dict.FindKey(func(k K) bool {
    return strings.HasPrefix(k, "prefix")
})

// Find first matching value
value, found := dict.FindValue(func(v V) bool {
    return v%2 == 0
})

// Check if value exists
exists := dict.ContainsValue(42)
```

### Conversion Operations

```go
// Get all keys, values, or pairs
keys := dict.Keys()
values := dict.Values()
pairs := dict.Items()

// Convert to native Go map
nativeMap := dict.AsMap()
```

## Performance Characteristics

| Implementation | Get | Put | Remove | Memory | Thread-Safe |
|---------------|-----|-----|--------|---------|-------------|
| Hash | O(1) | O(1) | O(1) | Low | No |
| ConcurrentHash | O(1) | O(1) | O(1) | Low | Yes |
| ConcurrentHashRW | O(1) | O(1) | O(1) | Low | Yes |
| Tree | O(log n) | O(log n) | O(log n) | Medium | No |

## Benchmark Results

```
BenchmarkComparison_Get/Hash-16                231361258    5.188 ns/op    0 B/op    0 allocs/op
BenchmarkComparison_Get/ConcurrentHash-16      121083358   10.02 ns/op     0 B/op    0 allocs/op
BenchmarkComparison_Get/ConcurrentHashRW-16    123224967    9.814 ns/op    0 B/op    0 allocs/op
BenchmarkComparison_Get/NativeMap-16           233014706    5.203 ns/op    0 B/op    0 allocs/op
```

## Usage Examples

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/pickeringtech/go-collections/collections/dicts"
)

func main() {
    // Create a dictionary
    dict := dicts.NewHash(
        dicts.Pair[string, int]{Key: "apple", Value: 5},
        dicts.Pair[string, int]{Key: "banana", Value: 3},
        dicts.Pair[string, int]{Key: "cherry", Value: 8},
    )

    // Filter fruits with count > 4
    expensive := dict.Filter(func(fruit string, count int) bool {
        return count > 4
    })

    fmt.Printf("Expensive fruits: %v\n", expensive.Keys())
    // Output: Expensive fruits: [apple cherry]
}
```

### Concurrent Usage

```go
package main

import (
    "fmt"
    "sync"
    "github.com/pickeringtech/go-collections/collections/dicts"
)

func main() {
    dict := dicts.NewConcurrentHash[string, int]()
    var wg sync.WaitGroup

    // Multiple goroutines adding data
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            dict.PutInPlace(fmt.Sprintf("key%d", id), id)
        }(i)
    }

    wg.Wait()
    fmt.Printf("Final size: %d\n", dict.Length())
}
```

## Best Practices

1. **Choose the Right Implementation**:
   - Use `Hash` for single-threaded, high-performance scenarios
   - Use `ConcurrentHash` for multi-threaded scenarios with balanced read/write
   - Use `ConcurrentHashRW` for read-heavy multi-threaded scenarios
   - Use `Tree` when you need sorted iteration or range queries

2. **Immutable vs Mutable**:
   - Use immutable operations (`Put`, `Remove`, `Filter`) for functional programming style
   - Use mutable operations (`PutInPlace`, `RemoveInPlace`, `FilterInPlace`) for performance-critical scenarios

3. **Memory Management**:
   - Immutable operations create new dictionaries; be mindful of memory usage
   - Use `Clear()` to reset large dictionaries instead of creating new ones

4. **Error Handling**:
   - Always check the boolean return value from `Get()` operations
   - Use `Contains()` for existence checks when you don't need the value
