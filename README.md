# Go Collections

<a href="https://github.com/pickeringtech/go-collections/actions/workflows/ci.yml"><img src="https://github.com/pickeringtech/go-collections/actions/workflows/ci.yml/badge.svg?branch=main" alt="CI"></a>
<a href="https://github.com/pickeringtech/go-collections/actions/workflows/mutation.yml"><img src="https://github.com/pickeringtech/go-collections/actions/workflows/mutation.yml/badge.svg?branch=main" alt="Mutation testing"></a>
<a href="https://codecov.io/gh/pickeringtech/go-collections"><img src="https://codecov.io/gh/pickeringtech/go-collections/graph/badge.svg?token=J2EZ0A9GUI" alt="Coverage"></a>
<a href="https://pkg.go.dev/github.com/pickeringtech/go-collections"><img src="https://img.shields.io/badge/api-reference-blue.svg?style=flat-square" alt="GoDoc"></a>

**A comprehensive, type-safe, and high-performance collections library for Go**

Go Collections provides data structures and utilities for working with collections in Go. It includes thread-safe maps, mathematical sets, and flexible lists.

## Features

- **Thread-Safe**: Concurrent implementations for multi-threaded applications
- **High Performance**: Optimized implementations with detailed benchmarks
- **Type-Safe**: Full generic support with compile-time type checking
- **Rich APIs**: Operations for filtering, mapping, and transforming data
- **Well Documented**: Examples and clear documentation
- **Zero Dependencies**: Pure Go implementation with no external dependencies
- **Familiar**: APIs inspired by popular languages (Java, Python, JavaScript)

## Quick Start

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

    fmt.Printf("Tasks: %v\n", tasks.AsSlice())
}
```

## What's Included

### Core Collections

| Package | Description | Thread-Safe Options |
|---------|-------------|---------------------|
| **[Dicts](./collections/dicts/)** | Key-value mappings with rich operations | Mutex & RWMutex |
| **[Sets](./collections/sets/)** | Mathematical sets with union, intersection | Mutex & RWMutex |
| **[Lists](./collections/lists/)** | Flexible sequences with stack/queue operations | Mutex & RWMutex |
| **[Multimaps](./collections/multimaps/)** | One key to many values (list- or set-backed) | Mutex & RWMutex |
| **[Deques](./collections/deques/)** | Double-ended queue / bounded ring buffer | Mutex & RWMutex |
| **[Heaps](./collections/heaps/)** | Binary heap / priority queue (min, max, comparator) | Mutex & RWMutex |
| **[LRU](./collections/lru/)** | Bounded cache with least-recently-used eviction | Mutex & RWMutex |

### Utilities

| Package | Description | Use Cases |
|---------|-------------|-----------|
| **[Slices](./slices/)** | Enhanced slice operations | Filtering, mapping, reducing |
| **[Maps](./maps/)** | Native map utilities | Key extraction, value transformation |
| **[Channels](./channels/)** | Channel-based pipelines | Stream processing, fan-out/fan-in |
| **[Constraints](./constraints/)** | Type constraints for generics | Custom generic functions |

## Choose Your Data Structure

### When to Use Dicts (Maps)
```go
// Perfect for key-value relationships
userRoles := collections.NewDict(
    collections.Pair[string, string]{Key: "admin", Value: "Administrator"},
    collections.Pair[string, string]{Key: "user", Value: "Regular User"},
)

// Rich operations
admins := userRoles.Filter(func(role, title string) bool {
    return role == "admin"
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

### When to Use Heaps (Priority Queues)
```go
import (
    "github.com/pickeringtech/go-collections/collections"
    "github.com/pickeringtech/go-collections/collections/heaps" // for the in-place API
)

// Smallest-first scheduling — reachable straight from the facade.
// Pop is immutable: it returns the element, an ok flag, and the remaining heap.
pq := collections.NewMinHeap(5, 1, 3)
next, ok, rest := pq.Pop() // next == 1, ok == true; rest is pq without it

// Or order by any comparator (e.g. a struct field)
tasks := collections.NewHeap(func(a, b Task) bool { return a.Priority > b.Priority })

// The heaps subpackage adds the in-place, mutating API (PushInPlace / PopInPlace)
mpq := heaps.NewMin(5, 1, 3)
mpq.PushInPlace(0)
top, _ := mpq.PopInPlace() // 0
```

**Use when**: You always need the most- (or least-) extreme item next —
scheduling, Dijkstra / A* frontiers, streaming top-k, or merging sorted streams.

### When to Use an LRU Cache
```go
import (
    "github.com/pickeringtech/go-collections/collections"
    "github.com/pickeringtech/go-collections/collections/lru" // for eviction options
)

// Bounded cache that evicts the least-recently-used entry — from the facade
cache := collections.NewLRU[string, int](2)
cache.PutInPlace("a", 1)
cache.PutInPlace("b", 2)
v, ok := cache.Get("a")  // promotes "a"; adding a third entry now evicts "b"

// Eviction callbacks and seed entries via lru.Option
cache = collections.NewLRU[string, int](100,
    lru.WithOnEvict(func(k string, v int) { /* ... */ }),
)
```

**Use when**: You need a fixed-memory cache with automatic eviction —
hot-key caches, memoisation with a budget, or any bounded most-recently-seen store.

## Iterator Bridge (Go 1.23+)

Build collections straight from any range-over-func iterator — the inbound
counterpart to each collection's `Values`/`All` iterators:

```go
import "slices" // std-lib, for slices.Values

// From a value iterator (iter.Seq)
list := collections.ListFromSeq(slices.Values([]int{3, 1, 2}))
set := collections.SetFromSeq(slices.Values([]string{"a", "b", "a"}))

// From a key/value iterator (iter.Seq2) — last value wins per key
prices := collections.DictFromSeq2(other.All())
```

`DequeFromSeq`, `HeapFromSeq`, `ListMultimapFromSeq2` and `SetMultimapFromSeq2`
cover the remaining families.

## Thread Safety

All collections offer thread-safe variants:

```go
// Each collection has two thread-safe variants. Pick one per use:

// Balanced read/write:
dict := collections.NewConcurrentDict(pairs...)
set := collections.NewConcurrentSet(items...)
list := collections.NewConcurrentList(items...)

// Or, for read-heavy workloads:
dict = collections.NewConcurrentRWDict(pairs...)
set = collections.NewConcurrentRWSet(items...)
list = collections.NewConcurrentRWList(items...)
```

## Performance

All collection implementations are continuously benchmarked. The headline
preview below — a curated table and chart — is **regenerated and committed on
every push to `main`** by the `benchreport` tool, so it never goes stale. The
full per-operation matrix lives in [BENCHMARKS.md](BENCHMARKS.md).

<!-- BENCH:START -->

<!-- Generated by tools/benchreport — do not edit by hand. Regenerate with `make bench-report`. -->

_Headline numbers are from the controlled **Reference — Framework Desktop** baseline; the shared-runner **CI** numbers are indicative only. See [BENCHMARKS.md](BENCHMARKS.md) for both full matrices and caveats._

| Operation | ns/op | B/op | allocs/op |
|---|--:|--:|--:|
| Dict — Hash.Get | 4.44 | 0 | 0 |
| Dict — ConcurrentHash.Get | 12.8 | 0 | 0 |
| Dict — ConcurrentHashRW.Get | 12.6 | 0 | 0 |
| Dict — Tree.Get | 411 | 0 | 0 |
| List — Array.Get | 1.41 | 0 | 0 |
| Set — Hash.Contains | 5.23 | 0 | 0 |

![Benchmark chart](docs/bench.svg)

Reference — Framework Desktop: `31d657d` · 2026-06-16 · linux/amd64 · Go go1.25.5

CI — GitHub-hosted runner (ubuntu-latest): `e5d4871` · 2026-06-18 · linux/amd64 · Go go1.24.13

Full report → [BENCHMARKS.md](BENCHMARKS.md)

Performance trend across recent commits → [BENCHMARKS.md](BENCHMARKS.md#trend-recent-main-commits)
<!-- BENCH:END -->

## Documentation & Examples

Each package includes:
- **Comprehensive README** with usage examples
- **GoDoc examples** for every major operation
- **Performance characteristics** and best practices
- **Real-world use cases** and patterns

### Runnable Examples

The [`examples/`](./examples) directory holds small, focused, **runnable apps**
that exercise the public API in realistic, cross-package flows — and are
build-and-run E2E-tested against golden output in CI on every PR:

- **[word-frequency](./examples/cmd/word-frequency)** — tokenise text, count words, print top-N (`slices` + `maps`)
- **[set-algebra](./examples/cmd/set-algebra)** — union / intersection / difference / subset over `sets`
- **[worker-pipeline](./examples/cmd/worker-pipeline)** — fan-out/fan-in a stream through a bounded worker pool (`channels` + `concurrency`)
- **[ordered-processing](./examples/cmd/ordered-processing)** — reverse, replay and sort with `lists`

They consume the library as a separate module, so they double as a
downstream-consumer smoke test. See the [examples README](./examples/README.md).

### Package Documentation
- **[Collections Overview](./collections/README.md)** - Start here for core data structures
- **[Dicts Documentation](./collections/dicts/README.md)** - Key-value mappings
- **[Sets Documentation](./collections/sets/README.md)** - Mathematical sets
- **[Lists Documentation](./collections/lists/README.md)** - Ordered sequences
- **[Multimaps Documentation](./collections/multimaps/README.md)** - One key, many values
- **[LRU Documentation](./collections/lru/README.md)** - Bounded cache with eviction
- **[Slices Utilities](./slices/README.md)** - Enhanced slice operations
- **[Maps Utilities](./maps/README.md)** - Native map helpers
- **[Channels Utilities](./channels/README.md)** - Pipeline processing
- **[Mutation Testing](./docs/mutation-testing.md)** - How we verify the tests catch regressions, not just run lines

## Contributing

Contributions are welcome. See the [Contributing Guide](CONTRIBUTING.md) for details.

## License

Go Collections is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

Made by [Pickering Technologies](https://www.picktech.co.uk).
