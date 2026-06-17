# Dicts - Key-Value Mappings Made Simple

The `dicts` package provides type-safe dictionaries (maps) that extend Go's built-in maps. It adds thread-safe access, operations like filtering and searching, and sorted iteration.

## Quick Start

```go
import "github.com/pickeringtech/go-collections/collections/dicts"

// Create a dictionary with initial data
inventory := dicts.NewHash(
    dicts.Pair[string, int]{Key: "apples", Value: 50},
    dicts.Pair[string, int]{Key: "oranges", Value: 30},
    dicts.Pair[string, int]{Key: "bananas", Value: 20},
)

// Rich operations that native maps can't do
lowStock := inventory.Filter(func(item string, count int) bool {
    return count < 25
})

fmt.Printf("Low stock items: %v\n", lowStock.Keys()) // [oranges bananas]
```

## Why Use Dicts?

**Native Go maps are limited:**
```go
// Native maps - basic and limited
m := map[string]int{"a": 1, "b": 2}
// No filtering, no thread safety, no rich operations
```

**Dicts are powerful:**
```go
// Dicts - rich and flexible
d := dicts.NewHash(
    dicts.Pair[string, int]{Key: "a", Value: 1},
    dicts.Pair[string, int]{Key: "b", Value: 2},
)

// Rich operations
filtered := d.Filter(func(k string, v int) bool { return v > 1 })
key, value, found := d.Find(func(k string, v int) bool { return v == 2 })
concurrent := dicts.NewConcurrentHash(d.Items()...) // Thread-safe
```

## Available Implementations

### Hash Dictionary - Fast and Simple
**Perfect for**: Most use cases, fast lookups, general-purpose key-value storage

```go
// Hash-based dictionary
users := dicts.NewHash(
    dicts.Pair[int, string]{Key: 1, Value: "Alice"},
    dicts.Pair[int, string]{Key: 2, Value: "Bob"},
)

// O(1) operations
name, found := users.Get(1, "Unknown")  // Fast lookup
users.PutInPlace(3, "Charlie")          // Fast insertion
exists := users.Contains(2)             // Fast membership test
```

**Performance**: O(1) average case for all operations

### Concurrent Hash Dictionary - Thread-Safe
**Perfect for**: Multi-threaded applications, shared state, balanced read/write workloads

```go
// Thread-safe dictionary for concurrent access
counter := dicts.NewConcurrentHash(
    dicts.Pair[string, int]{Key: "requests", Value: 0},
)

// Safe from multiple goroutines. UpdateInPlace performs the read-modify-write
// atomically (under a single lock acquisition), so concurrent increments don't
// lose updates. A separate Get + PutInPlace would race: each call locks
// independently, so goroutines read the same value and overwrite each other.
var wg sync.WaitGroup
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        counter.UpdateInPlace("requests", func(old int, _ bool) int {
            return old + 1
        })
    }()
}
wg.Wait()
// counter.Get("requests", 0) == (100, true)
```

**Performance**: O(1) with mutex overhead (~2x slower than Hash)

### Concurrent RW Hash Dictionary - Read-Optimized
**Perfect for**: Read-heavy workloads, caching, configuration data

```go
// Optimized for concurrent reads
cache := dicts.NewConcurrentHashRW(
    dicts.Pair[string, []byte]{Key: "config", Value: configData},
)

// Multiple readers can access simultaneously
var wg sync.WaitGroup
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        data, _ := cache.Get("config", nil) // Concurrent reads!
        processConfig(data)
    }()
}
wg.Wait()
```

**Performance**: O(1) with read-write mutex (concurrent reads, exclusive writes)

### Tree Dictionary - Sorted Keys
**Perfect for**: Sorted iteration, range queries, ordered data

```go
// Maintains keys in sorted order
scores := dicts.NewTree(
    dicts.Pair[string, int]{Key: "charlie", Value: 85},
    dicts.Pair[string, int]{Key: "alice", Value: 92},
    dicts.Pair[string, int]{Key: "bob", Value: 78},
)

// Always iterate in sorted key order
fmt.Println("Leaderboard:")
scores.ForEach(func(name string, score int) {
    fmt.Printf("%s: %d\n", name, score)
})
// Output (always sorted):
// alice: 92
// bob: 78
// charlie: 85

// Get sorted keys
names := scores.Keys() // ["alice", "bob", "charlie"]
```

**Performance**: O(log n) operations, sorted iteration
**Constraint**: Keys must be comparable (strings, numbers)

## Choose Your Implementation

| Implementation | Use When | Performance | Thread-Safe |
|---------------|----------|-------------|-------------|
| `NewHash()` | Single-threaded, general use | Fastest | No |
| `NewConcurrentHash()` | Multi-threaded, balanced R/W | Fast | Yes |
| `NewConcurrentHashRW()` | Multi-threaded, read-heavy | Fast reads | Yes |
| `NewTree()` | Need sorted keys/iteration | Slower | No |

## Two Ways to Work: Immutable vs Mutable

### Immutable Style (Functional Programming)
Returns new dictionaries and leaves the original unchanged:

```go
dict := dicts.NewHash(
    dicts.Pair[string, int]{Key: "count", Value: 1},
)

// Chain operations functionally
result := dict.
    Put("total", 100).                              // Add new key
    Put("count", 2).                                // Update existing
    Filter(func(k string, v int) bool { return v > 1 }) // Keep values > 1

fmt.Printf("Original count: %d\n", dict.Get("count", 0))    // Original: 1
fmt.Printf("Result count: %d\n", result.Get("count", 0))    // Result: 2
```

### Mutable Style (Performance-Focused)
Modifies the dictionary in place, which is faster for performance-critical code:

```go
dict := dicts.NewHash[string, int]()

// Direct modifications for speed
dict.PutInPlace("count", 1)                        // Fast insertion
dict.PutInPlace("total", 100)                      // Fast insertion
removed, found := dict.RemoveInPlace("total")      // Fast removal

fmt.Printf("Count: %d, Removed: %d\n", dict.Get("count", 0), removed)
```

## Essential Operations

### Basic Access - The Building Blocks
```go
inventory := dicts.NewHash(
    dicts.Pair[string, int]{Key: "apples", Value: 50},
    dicts.Pair[string, int]{Key: "oranges", Value: 30},
)

// Safe access with defaults
apples, found := inventory.Get("apples", 0)
if found {
    fmt.Printf("We have %d apples\n", apples)
}

// Quick existence check
if inventory.Contains("bananas") {
    fmt.Println("We have bananas!")
} else {
    fmt.Println("No bananas in stock")
}

// Size information
fmt.Printf("Items in inventory: %d\n", inventory.Length())
if inventory.IsEmpty() {
    fmt.Println("Inventory is empty!")
}
```

### Filtering - Find What You Need
```go
products := dicts.NewHash(
    dicts.Pair[string, float64]{Key: "laptop", Value: 999.99},
    dicts.Pair[string, float64]{Key: "mouse", Value: 29.99},
    dicts.Pair[string, float64]{Key: "keyboard", Value: 79.99},
)

// Find expensive items (immutable - returns new dict)
expensive := products.Filter(func(name string, price float64) bool {
    return price > 50.0
})

// Remove cheap items (mutable - modifies original)
products.FilterInPlace(func(name string, price float64) bool {
    return price > 50.0  // Keep only expensive items
})

fmt.Printf("Expensive products: %v\n", expensive.Keys())
```

### Transforming to New Types: Map / Reduce
`Filter` is a method because it keeps the same key/value types. A general `Map`
to a **different** key/value type needs type parameters Go methods cannot have
([golang/go#49085](https://github.com/golang/go/issues/49085)), so `Map` and
`Reduce` are **free functions** over the `Dict` interface. `Map` returns the
`Dict` interface (backed by `NewHash`) so results chain on.

```go
prices := dicts.NewHash(
    dicts.Pair[string, float64]{Key: "laptop", Value: 999.99},
    dicts.Pair[string, float64]{Key: "mouse", Value: 29.99},
)

// Map: (K, V) -> (OK, OV) — here to (name, rounded-price-as-int)
rounded := dicts.Map(prices, func(name string, price float64) (string, int) {
    return name, int(price)
})                                                                   // Dict[string, int]

// Reduce: fold every pair into a single value
total := dicts.Reduce(prices, 0.0, func(acc float64, _ string, price float64) float64 {
    return acc + price
})
```

**`Map` always returns a `Hash`-backed `Dict`, even for a sorted (`Tree`) input.**
This is deliberate: `Map` may change the key type, and the output key type is only
constrained to `comparable`, not `Ordered` — so a sorted output cannot be
guaranteed in general. When your output key type *is* `Ordered` and you want to
keep sorted iteration, use `MapSorted`, which returns a `Tree`-backed
`SortedDict`:

```go
// MapSorted: like Map, but OK must be Ordered and the result stays sorted.
byPrice := dicts.MapSorted(prices, func(name string, price float64) (int, string) {
    return int(price), name
})                                                                   // SortedDict[int, string]
// byPrice iterates in ascending key order: 29, 999
```

Iteration order over a `Dict` is unspecified, so a reduction should be
order-independent. (`FlatMap` is intentionally omitted — flattening a dict of
dicts has no unambiguous key-merging rule.)

### Advanced Search - Find Exactly What You Want
```go
inventory := dicts.NewHash(
    dicts.Pair[string, int]{Key: "apples", Value: 50},
    dicts.Pair[string, int]{Key: "oranges", Value: 5},   // Low stock!
    dicts.Pair[string, int]{Key: "bananas", Value: 30},
)

// Find first low-stock item
item, count, found := inventory.Find(func(name string, qty int) bool {
    return qty < 10
})
if found {
    fmt.Printf("Low stock alert: %s (%d remaining)\n", item, count)
}

// Predicate checks over (key, value) pairs
allStocked := inventory.AllMatch(func(name string, qty int) bool { return qty > 0 })
anyLow := inventory.AnyMatch(func(name string, qty int) bool { return qty < 10 })
noneEmpty := inventory.NoneMatch(func(name string, qty int) bool { return qty == 0 })

// Find specific key pattern
fruitKey, found := inventory.FindKey(func(name string) bool {
    return strings.HasPrefix(name, "a")  // Fruits starting with 'a'
})

// Check if we have any item with exactly 30 units
has30 := inventory.ContainsValue(30)
fmt.Printf("Has item with 30 units: %t\n", has30)
```

`Find`, `AllMatch`, `AnyMatch` and `NoneMatch` form the search core shared
across the `lists`, `dicts` and `sets` families. `FindKey`, `FindValue` and
`ContainsValue` are deliberate dict-specific extensions that reflect a
dictionary's key-value shape.

### Data Extraction - Get What You Need
```go
scores := dicts.NewHash(
    dicts.Pair[string, int]{Key: "alice", Value: 95},
    dicts.Pair[string, int]{Key: "bob", Value: 87},
    dicts.Pair[string, int]{Key: "charlie", Value: 92},
)

// Extract data in different formats
students := scores.Keys()                    // ["alice", "bob", "charlie"]
allScores := scores.Values()                 // [95, 87, 92]
pairs := scores.Items()                      // []Pair[string, int]

// Convert to native Go map for interop
nativeMap := scores.AsMap()                  // map[string]int
fmt.Printf("Native map: %v\n", nativeMap)
```

## Real-World Examples

### Web Application Cache
```go
// Thread-safe cache for web application
cache := dicts.NewConcurrentHashRW[string, []byte]()

func getPage(url string) []byte {
    // Check cache first (concurrent reads are fast)
    if content, found := cache.Get(url, nil); found {
        return content
    }

    // Fetch and cache (exclusive write)
    content := fetchFromNetwork(url)
    cache.PutInPlace(url, content)
    return content
}

// Cache cleanup
cache.FilterInPlace(func(url string, content []byte) bool {
    return time.Since(getLastAccess(url)) < time.Hour
})
```

### Configuration Management
```go
// Application configuration with defaults
config := dicts.NewTree(
    dicts.Pair[string, string]{Key: "database.host", Value: "localhost"},
    dicts.Pair[string, string]{Key: "database.port", Value: "5432"},
    dicts.Pair[string, string]{Key: "server.port", Value: "8080"},
)

// Get config with fallback
dbHost := config.Get("database.host", "localhost")
dbPort := config.Get("database.port", "5432")

// Iterate in sorted order (thanks to Tree)
fmt.Println("Configuration:")
config.ForEach(func(key, value string) {
    fmt.Printf("  %s = %s\n", key, value)
})
```

### User Session Management
```go
// Thread-safe session store
sessions := dicts.NewConcurrentHash[string, Session]()

func createSession(userID string) string {
    sessionID := generateSessionID()
    session := Session{
        UserID:    userID,
        CreatedAt: time.Now(),
        LastSeen:  time.Now(),
    }
    sessions.PutInPlace(sessionID, session)
    return sessionID
}

func validateSession(sessionID string) (Session, bool) {
    session, found := sessions.Get(sessionID, Session{})
    if !found {
        return Session{}, false
    }

    // Update last seen
    session.LastSeen = time.Now()
    sessions.PutInPlace(sessionID, session)
    return session, true
}

// Cleanup expired sessions
func cleanupSessions() {
    sessions.FilterInPlace(func(id string, session Session) bool {
        return time.Since(session.LastSeen) < 24*time.Hour
    })
}
```

## Performance Guide

### Benchmark Results
```
BenchmarkComparison_Get/Hash-16                231M    5.188 ns/op    0 B/op    0 allocs/op
BenchmarkComparison_Get/ConcurrentHash-16      121M   10.02 ns/op     0 B/op    0 allocs/op
BenchmarkComparison_Get/ConcurrentHashRW-16    123M    9.814 ns/op    0 B/op    0 allocs/op
BenchmarkComparison_Get/Tree-16                 50M   25.67 ns/op     0 B/op    0 allocs/op
BenchmarkComparison_Get/NativeMap-16           233M    5.203 ns/op    0 B/op    0 allocs/op
```

### Performance Characteristics

| Implementation | Get | Put | Remove | Memory | Thread-Safe | Best For |
|---------------|-----|-----|--------|---------|-------------|----------|
| Hash | O(1) | O(1) | O(1) | Low | No | Single-threaded, high performance |
| ConcurrentHash | O(1) | O(1) | O(1) | Low | Yes | Balanced read/write workloads |
| ConcurrentHashRW | O(1) | O(1) | O(1) | Low | Yes | Read-heavy workloads |
| Tree | O(log n) | O(log n) | O(log n) | Medium | No | Sorted iteration needed |

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

## Best Practices and Quick Reference

### Choose the Right Implementation

| Scenario | Use | Why |
|----------|-----|-----|
| Single-threaded app | `NewHash()` | Fastest performance |
| Multi-threaded, balanced R/W | `NewConcurrentHash()` | Simple thread safety |
| Multi-threaded, read-heavy | `NewConcurrentHashRW()` | Concurrent reads |
| Need sorted iteration | `NewTree()` | Maintains key order |

### Performance Tips

```go
// Use defaults with Get() for safety
value, found := dict.Get(key, defaultValue)

// Batch operations for better performance
for key, value := range updates {
    dict.PutInPlace(key, value)
}

// Use Contains() for existence checks
if dict.Contains(key) {
    // Process existing key
}

// Choose mutable operations for performance-critical code
dict.FilterInPlace(predicate)  // Faster than Filter()
```

### Immutable vs Mutable Strategy

```go
// Functional style - use immutable operations
result := dict.
    Put("new", value).
    Filter(predicate).
    Remove("old")

// Performance style - use mutable operations
dict.PutInPlace("new", value)
dict.FilterInPlace(predicate)
dict.RemoveInPlace("old")
```

### Quick Reference

| Operation | Immutable | Mutable | Use Case |
|-----------|-----------|---------|----------|
| Add/Update | `Put(k, v)` | `PutInPlace(k, v)` | Insert or update |
| Atomic update | — | `UpdateInPlace(k, fn)` | Race-free read-modify-write (counters) |
| Remove | `Remove(k)` | `RemoveInPlace(k)` | Delete key |
| Filter | `Filter(fn)` | `FilterInPlace(fn)` | Conditional removal |
| Access | `Get(k, default)` | `Get(k, default)` | Safe retrieval |
| Check | `Contains(k)` | `Contains(k)` | Membership test |

**Start with `NewHash()` and upgrade to concurrent versions only when you need thread safety!**
