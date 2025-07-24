# Dicts - Key-Value Mappings Made Simple

The `dicts` package provides powerful, type-safe dictionaries (maps) that go far beyond Go's built-in maps. Whether you need thread-safe access, rich operations like filtering and searching, or sorted iteration, dicts has you covered.

## 🚀 Quick Start

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

## ✨ Why Use Dicts?

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
concurrent := dicts.NewConcurrentHash(d.Items()...) // Thread-safe!
```

## 📦 Available Implementations

### 🏃‍♂️ Hash Dictionary - Fast & Simple
**Perfect for**: Most use cases, fast lookups, general-purpose key-value storage

```go
// Lightning-fast hash-based dictionary
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

### 🔒 Concurrent Hash Dictionary - Thread-Safe
**Perfect for**: Multi-threaded applications, shared state, balanced read/write workloads

```go
// Thread-safe dictionary for concurrent access
counter := dicts.NewConcurrentHash(
    dicts.Pair[string, int]{Key: "requests", Value: 0},
)

// Safe from multiple goroutines
var wg sync.WaitGroup
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        current, _ := counter.Get("requests", 0)
        counter.PutInPlace("requests", current+1)
    }()
}
wg.Wait()
```

**Performance**: O(1) with mutex overhead (~2x slower than Hash)

### 📖 Concurrent RW Hash Dictionary - Read-Optimized
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

### 🌳 Tree Dictionary - Sorted Keys
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

## 🎯 Choose Your Implementation

| Implementation | Use When | Performance | Thread-Safe |
|---------------|----------|-------------|-------------|
| `NewHash()` | Single-threaded, general use | Fastest | ❌ |
| `NewConcurrentHash()` | Multi-threaded, balanced R/W | Fast | ✅ |
| `NewConcurrentHashRW()` | Multi-threaded, read-heavy | Fast reads | ✅ |
| `NewTree()` | Need sorted keys/iteration | Slower | ❌ |

## 🔄 Two Ways to Work: Immutable vs Mutable

### 🧊 Immutable Style (Functional Programming)
Returns new dictionaries, original unchanged - perfect for functional programming:

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

### ⚡ Mutable Style (Performance-Focused)
Modifies the dictionary in place - perfect for performance-critical code:

```go
dict := dicts.NewHash[string, int]()

// Direct modifications for speed
dict.PutInPlace("count", 1)                        // Fast insertion
dict.PutInPlace("total", 100)                      // Fast insertion
removed, found := dict.RemoveInPlace("total")      // Fast removal

fmt.Printf("Count: %d, Removed: %d\n", dict.Get("count", 0), removed)
```

## 🛠️ Essential Operations

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

### 🔍 Smart Filtering - Find What You Need
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

### 🎯 Advanced Search - Find Exactly What You Want
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

// Find specific key pattern
fruitKey, found := inventory.FindKey(func(name string) bool {
    return strings.HasPrefix(name, "a")  // Fruits starting with 'a'
})

// Check if we have any item with exactly 30 units
has30 := inventory.ContainsValue(30)
fmt.Printf("Has item with 30 units: %t\n", has30)
```

### 📊 Data Extraction - Get What You Need
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

## 🌟 Real-World Examples

### Web Application Cache
```go
// Thread-safe cache for web application
cache := dicts.NewConcurrentRWHash[string, []byte]()

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

## 📊 Performance Guide

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
| Hash | O(1) | O(1) | O(1) | Low | ❌ | Single-threaded, high performance |
| ConcurrentHash | O(1) | O(1) | O(1) | Low | ✅ | Balanced read/write workloads |
| ConcurrentHashRW | O(1) | O(1) | O(1) | Low | ✅ | Read-heavy workloads |
| Tree | O(log n) | O(log n) | O(log n) | Medium | ❌ | Sorted iteration needed |

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

## 🎯 Best Practices & Quick Reference

### 🏆 Choose the Right Implementation

| Scenario | Use | Why |
|----------|-----|-----|
| Single-threaded app | `NewHash()` | Fastest performance |
| Multi-threaded, balanced R/W | `NewConcurrentHash()` | Simple thread safety |
| Multi-threaded, read-heavy | `NewConcurrentHashRW()` | Concurrent reads |
| Need sorted iteration | `NewTree()` | Maintains key order |

### ⚡ Performance Tips

```go
// ✅ Use defaults with Get() for safety
value, found := dict.Get(key, defaultValue)

// ✅ Batch operations for better performance
for key, value := range updates {
    dict.PutInPlace(key, value)
}

// ✅ Use Contains() for existence checks
if dict.Contains(key) {
    // Process existing key
}

// ✅ Choose mutable operations for performance-critical code
dict.FilterInPlace(predicate)  // Faster than Filter()
```

### 🔄 Immutable vs Mutable Strategy

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

### 🚀 Quick Reference

| Operation | Immutable | Mutable | Use Case |
|-----------|-----------|---------|----------|
| Add/Update | `Put(k, v)` | `PutInPlace(k, v)` | Insert or update |
| Remove | `Remove(k)` | `RemoveInPlace(k)` | Delete key |
| Filter | `Filter(fn)` | `FilterInPlace(fn)` | Conditional removal |
| Access | `Get(k, default)` | `Get(k, default)` | Safe retrieval |
| Check | `Contains(k)` | `Contains(k)` | Membership test |

**Start with `NewHash()` and upgrade to concurrent versions only when you need thread safety!**
