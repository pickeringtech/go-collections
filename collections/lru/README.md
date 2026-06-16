# LRU Cache - Bounded Maps with Eviction

The `lru` package provides a generic, bounded **least-recently-used cache**: a
fixed-capacity key-value store that evicts the least-recently-used entry when it
would otherwise overflow. It is the cache Go's standard library never ships —
type-safe, O(1), and offered in the same immutable/mutable and thread-safe
shapes as the rest of go-collections.

## Quick Start

```go
import "github.com/pickeringtech/go-collections/collections/lru"

cache := lru.NewLRU[string, int](2) // holds at most 2 entries

cache.PutInPlace("a", 1)
cache.PutInPlace("b", 2)
cache.PutInPlace("c", 3) // capacity exceeded → evicts the least-recently-used "a"

value, found := cache.Get("b") // value == 2, found == true
_, stillThere := cache.Peek("a") // stillThere == false (evicted)
```

## Why Use an LRU Cache?

**A native map never forgets, so it leaks memory as a cache:**
```go
cache := map[string][]byte{}
cache[key] = fetch(key) // nothing is ever evicted — grows without bound
```

**An LRU cache is bounded, with O(1) eviction handled for you:**
```go
cache := lru.NewLRU[string, []byte](1000)
cache.PutInPlace(key, fetch(key)) // the coldest entry is dropped at capacity
```

Under the hood it pairs a map (for O(1) lookup) with an intrusive doubly-linked
list (for O(1) recency updates), so every `Get`, `Peek`, `Put` and `Remove` is
constant-time regardless of size.

## Recency: `Get` vs `Peek`

`Get` counts as a use — a hit promotes the entry to most-recently-used, so it is
evicted last. `Peek` reads the same value **without** touching recency.

```go
cache := lru.NewLRU[string, int](2,
    lru.WithEntries(
        lru.Pair[string, int]{Key: "a", Value: 1},
        lru.Pair[string, int]{Key: "b", Value: 2},
    ),
)

cache.Get("a")          // promotes "a"; "b" is now the eviction candidate
cache.PutInPlace("c", 3) // evicts "b", not "a"
```

`Peek`, `Contains`, `Length`, `IsEmpty`, `Capacity` and the exporters (`Keys`,
`Values`, `Items`, `AsMap`, `ForEach`, `All`) are all recency-neutral.

## Eviction Callback

Register a callback to learn when capacity forces an entry out — to close a
handle, decrement a refcount, or emit a metric:

```go
cache := lru.NewLRU[string, *os.File](100,
    lru.WithOnEvict(func(key string, f *os.File) { f.Close() }),
)
```

The callback fires **only** for capacity-driven evictions, never for entries you
remove explicitly via `Remove`, `RemoveInPlace` or `Clear`.

## Implementations

| Type                  | Thread-safe | Backing lock     | Best for                       |
|-----------------------|-------------|------------------|--------------------------------|
| `LRU[K, V]`           | No          | —                | Single-threaded use            |
| `ConcurrentLRU[K, V]` | Yes         | `sync.Mutex`     | Balanced read/write workloads  |
| `ConcurrentLRURW[K, V]`| Yes        | `sync.RWMutex`   | Read-heavy workloads           |

All three implement the same `Cache` and `MutableCache` interfaces, so you can
swap between them without touching call sites.

> On `ConcurrentLRURW`, `Get` takes the **write** lock because it re-orders
> recency. Use `Peek` when a lookup need not count as a use, to keep reads
> concurrent.

## Immutable vs Mutable Operations

Mutable operations modify the receiver — the common case for a cache:

```go
cache.PutInPlace("key", value)            // insert/update in place
value, found := cache.Get("key")          // read + promote in place
removed, ok := cache.RemoveInPlace("key")
cache.Clear()
```

Immutable operations leave the receiver untouched and return a new cache:

```go
updated := cache.Put("key", value)        // returns a new Cache
smaller := cache.Remove("key")            // returns a new Cache
```

On the concurrent types the immutable operations return a new value of the
**same** concurrent type, so thread-safety is never silently dropped.

## Iteration

`ForEach` and the range-over-func iterator `All` both visit entries from
most- to least-recently-used, without affecting recency:

```go
for key, value := range cache.All() {
    fmt.Printf("%s=%d\n", key, value)
}
```

## Performance

| Operation         | Cost                                          |
|-------------------|-----------------------------------------------|
| `Get` / `Peek`    | O(1)                                          |
| `Put` / `PutInPlace` | O(1)                                       |
| `Remove`          | O(1) in place, O(n) for the immutable copy    |
| `Keys`/`Values`/… | O(n)                                          |

Start with `NewLRU` and upgrade to a concurrent variant only when you actually
share the cache across goroutines.
