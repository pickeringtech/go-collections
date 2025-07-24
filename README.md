# Go Collections

<a href="https://pkg.go.dev/github.com/pickeringtech/go-collections"><img src="https://img.shields.io/badge/api-reference-blue.svg?style=flat-square" alt="GoDoc"></a>

**A comprehensive, type-safe, and high-performance collections library for Go**

Go Collections provides powerful data structures and utilities that make working with collections in Go simple, safe, and efficient. Whether you need thread-safe maps, mathematical sets, or flexible lists, this library has you covered.

## ✨ Features

- **🔒 Thread-Safe**: Concurrent implementations for multi-threaded applications
- **⚡ High Performance**: Optimized implementations with detailed benchmarks
- **🎯 Type-Safe**: Full generic support with compile-time type checking
- **🧩 Rich APIs**: Comprehensive operations for filtering, mapping, and transforming data
- **📚 Well Documented**: Extensive examples and clear documentation
- **🔧 Zero Dependencies**: Pure Go implementation with no external dependencies
- **🎨 Familiar**: APIs inspired by popular languages (Java, Python, JavaScript)

## 🚀 Quick Start

```bash
go get github.com/pickeringtech/go-collections
```

```go
package main

import (
    "fmt"
    "github.com/pickeringtech/go-collections/collections"
)

func main() {
    // Create a thread-safe dictionary
    users := collections.NewConcurrentDict(
        collections.Pair[string, string]{Key: "john", Value: "John Doe"},
        collections.Pair[string, string]{Key: "jane", Value: "Jane Smith"},
    )

    // Safe concurrent access
    go func() {
        users.PutInPlace("bob", "Bob Wilson")
    }()

    name, found := users.Get("john", "Unknown")
    fmt.Printf("User: %s (found: %t)\n", name, found)

    // Create a mathematical set
    numbers := collections.NewSet(1, 2, 3, 4, 5)
    evens := numbers.Filter(func(n int) bool { return n%2 == 0 })

    fmt.Printf("Even numbers: %v\n", evens.AsSlice())

    // Create a flexible list
    tasks := collections.NewList("design", "code", "test")
    tasks.PushInPlace("deploy")

    fmt.Printf("Tasks: %v\n", tasks.GetAsSlice())
}
```

## 📦 What's Included

### Core Collections

| Package | Description | Thread-Safe Options |
|---------|-------------|---------------------|
| **[Dicts](./collections/dicts/)** | Key-value mappings with rich operations | ✅ Mutex & RWMutex |
| **[Sets](./collections/sets/)** | Mathematical sets with union, intersection | ✅ Mutex & RWMutex |
| **[Lists](./collections/lists/)** | Flexible sequences with stack/queue operations | ✅ Mutex & RWMutex |

### Utilities

| Package | Description | Use Cases |
|---------|-------------|-----------|
| **[Slices](./slices/)** | Enhanced slice operations | Filtering, mapping, reducing |
| **[Maps](./maps/)** | Native map utilities | Key extraction, value transformation |
| **[Channels](./channels/)** | Channel-based pipelines | Stream processing, fan-out/fan-in |
| **[Constraints](./constraints/)** | Type constraints for generics | Custom generic functions |

## 🎯 Choose Your Data Structure

### When to Use Dicts (Maps)
```go
// Perfect for key-value relationships
userRoles := collections.NewDict(
    collections.Pair[string, string]{Key: "admin", Value: "Administrator"},
    collections.Pair[string, string]{Key: "user", Value: "Regular User"},
)

// Rich operations
admins := userRoles.Filter(func(role, title string) bool {
    return strings.Contains(title, "Admin")
})
```

**Use when**: You need fast lookups, key-value relationships, or caching.

### When to Use Sets
```go
// Perfect for unique collections and mathematical operations
allowed := collections.NewSet("read", "write", "execute")
requested := collections.NewSet("read", "write", "delete")

// Mathematical operations
granted := allowed.Intersection(requested)
denied := requested.Difference(allowed)
```

**Use when**: You need unique elements, set operations, or membership testing.

### When to Use Lists
```go
// Perfect for ordered sequences
queue := collections.NewConcurrentList[string]()

// Thread-safe queue operations
go func() {
    queue.EnqueueInPlace("task1")
    queue.EnqueueInPlace("task2")
}()

task, found := queue.DequeueInPlace()
```

**Use when**: You need ordered data, stacks, queues, or sequential processing.

## 🔒 Thread Safety Made Simple

All collections offer thread-safe variants:

```go
// Choose your concurrency model
dict := collections.NewConcurrentDict(...)     // Balanced read/write
dict := collections.NewConcurrentRWDict(...)   // Read-heavy workloads

set := collections.NewConcurrentSet(...)       // Balanced read/write
set := collections.NewConcurrentRWSet(...)     // Read-heavy workloads

list := collections.NewConcurrentList(...)     // Balanced read/write
list := collections.NewConcurrentRWList(...)   // Read-heavy workloads
```

## 📊 Performance

All implementations are benchmarked and optimized:

```
BenchmarkDict_Get/Hash-16                228M    5.248 ns/op    0 B/op
BenchmarkDict_Get/ConcurrentHash-16      100M   10.41 ns/op     0 B/op
BenchmarkDict_Get/ConcurrentHashRW-16    100M   10.30 ns/op     0 B/op
BenchmarkDict_Get/Tree-16                 50M   25.67 ns/op     0 B/op
```

## 📚 Documentation & Examples

Each package includes:
- **Comprehensive README** with usage examples
- **GoDoc examples** for every major operation
- **Performance characteristics** and best practices
- **Real-world use cases** and patterns

### Package Documentation
- **[Collections Overview](./collections/README.md)** - Start here for core data structures
- **[Dicts Documentation](./collections/dicts/README.md)** - Key-value mappings
- **[Sets Documentation](./collections/sets/README.md)** - Mathematical sets
- **[Lists Documentation](./collections/lists/README.md)** - Ordered sequences
- **[Slices Utilities](./slices/README.md)** - Enhanced slice operations
- **[Maps Utilities](./maps/README.md)** - Native map helpers
- **[Channels Utilities](./channels/README.md)** - Pipeline processing

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## 📄 License

Go Collections is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

Made with ♥ by [Pickering Technologies](https://www.picktech.co.uk) - Your Strategic Technology Partner
