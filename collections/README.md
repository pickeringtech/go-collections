# Collections - Core Data Structures

The `collections` package provides the core of the Go Collections library with four data structures: **Dicts** (key-value mappings), **Sets** (unique collections), **Lists** (ordered sequences), and **Heaps** (priority queues).

## Quick Start

```go
import "github.com/pickeringtech/go-collections/collections"

// Create collections with simple factory functions
dict := collections.NewDict(
    collections.Pair[string, int]{Key: "apples", Value: 5},
    collections.Pair[string, int]{Key: "oranges", Value: 3},
)

set := collections.NewSet("red", "green", "blue")
list := collections.NewList(1, 2, 3, 4, 5)
```

## What's Inside

### Dicts - Key-Value Mappings
Use for lookups, caching, and associative data.

```go
// Create a user database
users := collections.NewDict(
    collections.Pair[int, string]{Key: 1, Value: "Alice"},
    collections.Pair[int, string]{Key: 2, Value: "Bob"},
)

// Fast lookups
name, found := users.Get(1, "Unknown")
fmt.Printf("User 1: %s\n", name) // User 1: Alice

// Rich operations
activeUsers := users.Filter(func(id int, name string) bool {
    return len(name) > 3 // Names longer than 3 characters
})
```

**Available Implementations:**
- `NewDict()` - Fast hash-based dictionary
- `NewConcurrentDict()` - Thread-safe with mutex
- `NewConcurrentRWDict()` - Thread-safe with read-write mutex (best for read-heavy)

### Sets - Unique Collections
Use for membership testing and mathematical operations.

```go
// Create permission sets
adminPerms := collections.NewSet("read", "write", "delete", "admin")
userPerms := collections.NewSet("read", "write")

// Mathematical operations
commonPerms := adminPerms.Intersection(userPerms)
extraPerms := adminPerms.Difference(userPerms)

fmt.Printf("Common: %v\n", commonPerms.AsSlice()) // [read write]
fmt.Printf("Admin-only: %v\n", extraPerms.AsSlice()) // [delete admin]

// Membership testing
canDelete := adminPerms.Contains("delete") // true
```

**Available Implementations:**
- `NewSet()` - Fast hash-based set
- `NewConcurrentSet()` - Thread-safe with mutex
- `NewConcurrentRWSet()` - Thread-safe with read-write mutex

### Lists - Ordered Sequences
Use for stacks, queues, and ordered data.

```go
// Create a task queue
tasks := collections.NewList("design", "implement", "test")

// Stack operations (LIFO)
tasks.PushInPlace("deploy")
lastTask, found := tasks.PopInPlace() // "deploy"

// Queue operations (FIFO)
tasks.EnqueueInPlace("monitor")
firstTask, found := tasks.DequeueInPlace() // "design"

// Rich operations
longTasks := tasks.Filter(func(task string) bool {
    return len(task) > 4
})
```

**Available Implementations:**
- `NewList()` - Linked list
- `NewDoublyLinkedList()` - Bidirectional linked list
- `NewConcurrentList()` - Thread-safe linked list
- `NewConcurrentDoublyLinkedList()` - Thread-safe bidirectional list
- `NewConcurrentRWList()` - Read-optimized thread-safe list

### Heaps - Priority Queues
Use when you always need the most- (or least-) extreme item next. See the
[`heaps` package](./heaps/README.md).

```go
import "github.com/pickeringtech/go-collections/collections/heaps"

// Smallest-first by default
pq := heaps.NewMin(5, 1, 3)
pq.PushInPlace(0)
next, _ := pq.PopInPlace() // 0

// Or order by any comparator
tasks := heaps.New(func(a, b Task) bool { return a.Priority > b.Priority })
```

**Available Implementations:**
- `heaps.NewMin()` / `heaps.NewMax()` - Min / max binary heap
- `heaps.New(less, ...)` - Comparator-driven binary heap
- `heaps.NewConcurrentMin()` / `heaps.NewConcurrentMax()` - Thread-safe (mutex)
- `heaps.NewConcurrentRWMin()` / `heaps.NewConcurrentRWMax()` - Read-optimized thread-safe

## Common Patterns

### Immutable vs Mutable Operations

All collections support both paradigms:

```go
dict := collections.NewDict(
    collections.Pair[string, int]{Key: "count", Value: 1},
)

// Immutable - returns new collection
newDict := dict.Put("count", 2)
fmt.Printf("Original: %d, New: %d\n",
    dict.Get("count", 0), newDict.Get("count", 0)) // Original: 1, New: 2

// Mutable - modifies in place
dict.PutInPlace("count", 3)
fmt.Printf("Modified: %d\n", dict.Get("count", 0)) // Modified: 3
```

### Functional Programming Style

```go
numbers := collections.NewSet(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

result := numbers.
    Filter(func(n int) bool { return n%2 == 0 }).    // Keep evens
    Filter(func(n int) bool { return n > 4 })        // Keep > 4

fmt.Printf("Even numbers > 4: %v\n", result.AsSlice()) // [6 8 10]
```

### Thread-Safe Processing

```go
// Create a thread-safe counter
counter := collections.NewConcurrentDict(
    collections.Pair[string, int]{Key: "requests", Value: 0},
)

// Safe concurrent updates
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

total, _ := counter.Get("requests", 0)
fmt.Printf("Total requests: %d\n", total)
```

## Performance Guide

### Choose the Right Implementation

| Use Case | Recommendation | Why |
|----------|---------------|-----|
| Single-threaded, high performance | `NewDict()`, `NewSet()`, `NewList()` | No locking overhead |
| Balanced read/write, multi-threaded | `NewConcurrentDict()`, etc. | Simple mutex protection |
| Read-heavy, multi-threaded | `NewConcurrentRWDict()`, etc. | Concurrent reads |
| Need sorted iteration | `NewTreeDict()` | Maintains key order |
| Need bidirectional access | `NewDoublyLinkedList()` | O(n/2) average access |

### Performance Characteristics

| Operation | Dict | Set | List | Concurrent Overhead |
|-----------|------|-----|------|-------------------|
| Get/Contains | O(1) | O(1) | O(n) | ~2x slower |
| Put/Add | O(1) | O(1) | O(1) at ends | ~2x slower |
| Remove | O(1) | O(1) | O(n) | ~2x slower |
| Iteration | O(n) | O(n) | O(n) | Minimal |

## Detailed Documentation

Each data structure has comprehensive documentation:

- **[Dicts Documentation](./dicts/README.md)** - Complete guide to key-value mappings
- **[Sets Documentation](./sets/README.md)** - Complete guide to mathematical sets
- **[Lists Documentation](./lists/README.md)** - Complete guide to ordered sequences
- **[Heaps Documentation](./heaps/README.md)** - Complete guide to priority queues
- **[LRU Documentation](./lru/README.md)** - Complete guide to bounded caches with eviction

## Real-World Examples

### Web Application Cache
```go
// Thread-safe cache for web application
cache := collections.NewConcurrentRWDict[string, []byte]()

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
```

### Permission System
```go
// Role-based permissions
type Role struct {
    Name        string
    Permissions collections.Set[string]
}

admin := Role{
    Name: "admin",
    Permissions: collections.NewSet("read", "write", "delete", "manage"),
}

user := Role{
    Name: "user",
    Permissions: collections.NewSet("read", "write"),
}

// Check if user can perform action
func canPerform(role Role, action string) bool {
    return role.Permissions.Contains(action)
}

// Find common permissions
common := admin.Permissions.Intersection(user.Permissions)
```

### Task Processing Queue
```go
// Producer-consumer pattern with thread-safe queue
queue := collections.NewConcurrentList[Task]()

// Producer
go func() {
    for task := range taskChannel {
        queue.EnqueueInPlace(task)
    }
}()

// Consumer
go func() {
    for {
        if task, found := queue.DequeueInPlace(); found {
            processTask(task)
        } else {
            time.Sleep(10 * time.Millisecond)
        }
    }
}()
```

## Best Practices

1. **Choose Immutable for Functional Style**: Use `Put()`, `Add()`, `Filter()` for functional programming
2. **Choose Mutable for Performance**: Use `PutInPlace()`, `AddInPlace()` for high-performance scenarios
3. **Use RW Variants for Read-Heavy**: `NewConcurrentRWDict()` when reads outnumber writes 10:1
4. **Prefer Sets for Membership**: Use sets instead of maps when you only need to check existence
5. **Use Lists for Ordered Data**: When insertion order or sequential access matters

Start with the simple variants (`NewDict`, `NewSet`, `NewList`) and switch to concurrent versions only when needed.

