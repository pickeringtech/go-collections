# Multimaps - One Key, Many Values

The `multimaps` package provides generic multimaps for Go: collections that map a
single key to many values. Use it whenever you would otherwise hand-roll a
`map[K][]V` or `map[K]map[V]struct{}` and reimplement append, remove, counting,
and iteration by hand.

## Quick Start

```go
import "github.com/pickeringtech/go-collections/collections/multimaps"

// Group orders by customer, preserving order and duplicates
orders := multimaps.NewListMultimap[string, string]()
orders.PutInPlace("alice", "book")
orders.PutInPlace("alice", "pen")

alice := orders.Get("alice")   // [book pen]
entries := orders.Length()     // 2 (total key-value associations)
keys := orders.KeyCount()      // 1 (distinct keys)
```

## Why Use Multimaps?

**Native Go maps of slices are clunky and error-prone:**
```go
groups := map[string][]string{}
groups["alice"] = append(groups["alice"], "book") // manual append every time
vals := groups["alice"]                           // leaks the internal slice
// removing a single value, counting entries vs keys, deduping — all by hand
```

**Multimaps are concise:**
```go
groups := multimaps.NewListMultimap[string, string]()
groups.PutInPlace("alice", "book")
vals := groups.Get("alice")    // independent copy, safe to mutate
groups.RemoveInPlace("alice", "book")
```

## Value Semantics: ListMultimap vs SetMultimap

| | `ListMultimap` | `SetMultimap` |
|---|---|---|
| Backing | `map[K][]V` | `map[K]map[V]struct{}` |
| Duplicate values per key | Kept | Collapsed |
| Value order within a key | Insertion order | Unspecified |
| `ContainsEntry` cost | O(values per key) | O(1) average |
| Use when | Order/duplicates matter (event logs) | Duplicates are noise (tags, members) |

Key iteration order is unspecified for both (it follows Go's map iteration).

```go
list := multimaps.NewListMultimap[string, int]()
list.PutInPlace("a", 1)
list.PutInPlace("a", 1)
list.Length() // 2 — duplicate kept

set := multimaps.NewSetMultimap[string, int]()
set.PutInPlace("a", 1)
set.PutInPlace("a", 1)
set.Length() // 1 — duplicate collapsed
```

## `Length` vs `KeyCount`

- `Length()` — total number of entries (key-value associations). A key bound to
  three values contributes three.
- `KeyCount()` — number of distinct keys.

## Immutable vs Mutable Operations

Every type supports both styles. Methods without a suffix return a new multimap
and leave the receiver untouched; `InPlace` methods modify the receiver.

```go
// Immutable — returns a new multimap
updated := m.Put("key", value)
evens := m.Filter(func(k string, v int) bool { return v%2 == 0 })

// In-place — modifies the receiver
m.PutInPlace("key", value)
m.FilterInPlace(func(k string, v int) bool { return v%2 == 0 })
```

## Iteration

```go
m := multimaps.NewListMultimap[string, int]()
m.PutAllInPlace("a", 1, 2, 3)

// Flat iteration over every entry (callback)
m.ForEach(func(key string, value int) { /* ... */ })

// Per-key iteration with all values
m.ForEachKey(func(key string, values []int) { /* ... */ })

// Range-over-func (Go 1.23+ iterators)
for key, value := range m.All() { /* ... */ }
for key := range m.KeysSeq() { /* ... */ }
```

## Thread Safety

Each backing offers a lock-free type plus mutex and RWMutex variants:

```go
// List-backed
multimaps.NewListMultimap[K, V]()              // not thread-safe
multimaps.NewConcurrentListMultimap[K, V]()    // *sync.Mutex (balanced)
multimaps.NewConcurrentRWListMultimap[K, V]()  // *sync.RWMutex (read-heavy)

// Set-backed
multimaps.NewSetMultimap[K, V]()
multimaps.NewConcurrentSetMultimap[K, V]()
multimaps.NewConcurrentRWSetMultimap[K, V]()
```

Operating on a thread-safe multimap yields a thread-safe result: immutable
operations return a new multimap of the same concurrent type.

## Implementations

All six implementations satisfy the `Multimap` (immutable) and
`MutableMultimap` (immutable + in-place) interfaces, so you can program against
the interface and swap implementations freely:

- `ListMultimap`, `ConcurrentListMultimap`, `ConcurrentRWListMultimap`
- `SetMultimap`, `ConcurrentSetMultimap`, `ConcurrentRWSetMultimap`
