# Heaps - Priority Queues

The `heaps` package provides a generic binary-heap priority queue: the structure
the standard library only exposes through the clunky, non-generic
`container/heap`. Use it whenever you always need the most- (or least-) extreme
item next — scheduling, Dijkstra / A* frontiers, streaming top-k, or merging
sorted streams.

## Quick Start

```go
import "github.com/pickeringtech/go-collections/collections/heaps"

// A min-heap: the smallest element leaves first.
pq := heaps.NewMin(5, 1, 3, 2, 4)
next, _ := pq.Peek()            // 1
smallest, _ := pq.PopInPlace()  // 1

// A max-heap: the largest element leaves first.
mx := heaps.NewMax(5, 1, 3, 2, 4)
largest, _ := mx.Peek()         // 5

// A comparator-driven heap over any type.
type Task struct {
    Name     string
    Priority int
}
tasks := heaps.New(func(a, b Task) bool { return a.Priority > b.Priority })
tasks.PushInPlace(Task{"deploy", 10})
```

## Why Use a Heap?

**`container/heap` is awkward:** you implement five methods
(`Len`/`Less`/`Swap`/`Push`/`Pop`) on your own slice type, push and pop through
package-level functions, and you get no generics — every heap is a bespoke type
with `interface{}` plumbing.

```go
// Native approach - a bespoke type plus boilerplate per element type
type IntHeap []int
func (h IntHeap) Len() int            { return len(h) }
func (h IntHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x any)         { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() any           { /* ... */ }
```

**`heaps.Binary` is one generic type with ordinary methods:**

```go
pq := heaps.NewMin(5, 1, 3)
pq.PushInPlace(0)
v, _ := pq.PopInPlace()
```

The comparator is a `func(a, b T) bool` reporting whether `a` has higher
priority than `b` — the same less-than convention as the lists `Sort` methods.

## Implementations

### Binary Heap (`Binary[T]`)

Single-threaded binary heap backed by a slice. `Push`/`Pop` are O(log n), `Peek`
is O(1), and construction via `New`/`NewMin`/`NewMax` is O(n) (bottom-up
heapify).

```go
pq := heaps.NewMin(5, 1, 3, 2, 4)
sorted := pq.AsSortedSlice() // [1 2 3 4 5], receiver untouched
```

### Concurrent Binary Heap (`ConcurrentBinary[T]`)

Thread-safe with a single mutex. Best for balanced push/pop workloads.

```go
pq := heaps.NewConcurrentMin[int]()
pq.PushInPlace(7)
```

### Concurrent RW Binary Heap (`ConcurrentRWBinary[T]`)

Thread-safe with a read-write mutex — concurrent reads, exclusive writes. Best
for read-heavy (Peek-heavy) workloads.

```go
pq := heaps.NewConcurrentRWMin[int]()
```

Immutable operations on a concurrent heap return another concurrent heap of the
same kind, so thread-safe in yields thread-safe out.

## Comparators

| Constructor                | Ordering                              |
|----------------------------|---------------------------------------|
| `NewMin[T]` / `Min[T]`     | smallest first (`a < b`)              |
| `NewMax[T]` / `Max[T]`     | largest first (`a > b`)               |
| `New(less, ...)`           | any `func(a, b T) bool` comparator    |

`Min`/`Max` work for any `constraints.Ordered` type; `New` takes an arbitrary
comparator for custom or struct ordering.

## Immutable vs Mutable

Immutable operations return a new heap and leave the receiver untouched:

```go
bigger := pq.Push(7)        // returns a new heap
v, ok, rest := pq.Pop()     // returns the element and a new heap
```

In-place operations mutate the receiver and carry the `InPlace` suffix:

```go
pq.PushInPlace(7)           // modifies pq
v, ok := pq.PopInPlace()    // modifies pq
```

## Draining in Priority Order

The heap-array order is unspecified beyond the heap invariant. For priority
order, drain the heap (both forms leave the receiver intact):

```go
for v := range pq.Drain() { // iterator (iter.Seq)
    fmt.Println(v)
}
sorted := pq.AsSortedSlice() // []T in priority order
```
