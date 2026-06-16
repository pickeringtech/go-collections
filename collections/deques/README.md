# Deques

Generic double-ended queue (deque) backed by a ring buffer, with O(1) push and
pop at **both** ends and an optional **bounded** (circular buffer) mode.

```go
import "github.com/pickeringtech/go-collections/collections/deques"
```

## Quick Start

```go
d := deques.NewRingBuffer[int]()
d.PushBackInPlace(1)              // [1]
d.PushBackInPlace(2)             // [1 2]
d.PushFrontInPlace(0)           // [0 1 2]
front, _ := d.PopFrontInPlace() // 0, deque [1 2]
back, _ := d.PopBackInPlace()   // 2, deque [1]
```

## Bounded ring buffer

A bounded deque fixes its capacity and chooses an `OverflowPolicy` for what
happens when a push arrives while full:

```go
// Classic ring buffer — a push when full evicts the opposite end.
win := deques.NewBoundedRingBuffer[int](3, deques.OverwriteOldest)
win.PushBackInPlace(1)           // [1]
win.PushBackInPlace(2)          // [1 2]
win.PushBackInPlace(3)          // [1 2 3] (full)
win.PushBackInPlace(4)          // [2 3 4] — front evicted, returns true

// Reject-when-full — a push when full is a no-op that reports false.
buf := deques.NewBoundedRingBuffer[int](3, deques.RejectWhenFull)
buf.PushBackInPlace(1)       // [1]
buf.PushBackInPlace(2)      // [1 2]
buf.PushBackInPlace(3)      // [1 2 3] (full)
ok := buf.PushBackInPlace(4) // unchanged [1 2 3], ok == false
```

`Capacity()` returns the bound (or `deques.Unbounded` / `-1` for unbounded
deques) and `IsFull()` reports whether a bounded deque is at capacity.

## Immutable vs mutable

| Form | Method | Returns |
|------|--------|---------|
| Immutable | `PushBack(e)` | new `Deque[T]` |
| Immutable | `PopFront()` | `(element, present?, Deque[T])` |
| In-place | `PushBackInPlace(e)` | `bool` (acceptance) |
| In-place | `PopFrontInPlace()` | `(element, present?)` |

Immutable operations never modify the receiver; in-place operations mutate it
and return only a status.

## Iteration

Deques are iterator-native:

```go
for i, v := range d.All() { /* front to back */ }
for v := range d.Values() { /* front to back */ }
for i, v := range d.Backward() { /* back to front */ }
```

`ForEach` and `ForEachWithIndex` offer the same traversal via callbacks.

## Implementations

| Type | Thread safety |
|------|---------------|
| `RingBuffer` | none (lock-free, single goroutine) |
| `ConcurrentRingBuffer` | `sync.Mutex` |
| `ConcurrentRWRingBuffer` | `sync.RWMutex` (favour for read-heavy use) |

Each has an unbounded constructor (`New…RingBuffer`) and a bounded one
(`NewBounded…RingBuffer`). Operating on a concurrent deque returns a concurrent
deque of the same type — thread-safe in, thread-safe out.

The top-level `collections` package also offers `NewDeque`,
`NewBoundedDeque`, their concurrent variants, and a fluent `DequeBuilder`.
