// Package heaps provides a generic binary-heap priority queue — the structure
// the standard library only exposes through the clunky, non-generic
// container/heap interface.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/collections/heaps"
//
//	// A min-heap: the smallest element leaves first.
//	pq := heaps.NewMin(5, 1, 3, 2, 4)
//	next, _ := pq.Peek()          // 1
//	smallest, _ := pq.PopInPlace()  // 1
//
//	// A max-heap: the largest element leaves first.
//	mx := heaps.NewMax(5, 1, 3, 2, 4)
//	largest, _ := mx.Peek()       // 5
//
//	// A comparator-driven heap over any type.
//	type Task struct{ Name string; Priority int }
//	tasks := heaps.New(func(a, b Task) bool { return a.Priority > b.Priority })
//	tasks.PushInPlace(Task{"deploy", 10})
//
// # Why a heap (vs the standard library)
//
// container/heap works, but it is awkward: you implement five methods
// (Len/Less/Swap/Push/Pop) on your own slice type, push and pop through
// package-level functions, and you get no generics — every heap is a bespoke
// type with interface{} plumbing. heaps.Binary is a single generic type with
// ordinary methods and a func(T, T) bool comparator that matches the lists Sort
// convention.
//
// When to reach for it: scheduling by priority, Dijkstra / A* frontiers,
// streaming top-k, merging sorted streams, or any "always take the most/least
// extreme item next" loop.
//
// # Available Implementations
//
// Binary (heaps.Binary):
//   - Single-threaded binary heap backed by a slice.
//   - O(log n) Push/Pop, O(1) Peek, O(n) construction (heapify).
//
// Concurrent Binary (heaps.ConcurrentBinary):
//   - Thread-safe with a single mutex.
//   - Best for balanced push/pop workloads.
//
// Concurrent RW Binary (heaps.ConcurrentRWBinary):
//   - Thread-safe with a read-write mutex; concurrent reads, exclusive writes.
//   - Best for read-heavy (Peek-heavy) workloads.
//
// # Comparators
//
// Ordering is supplied as a LessFunc — func(a, b T) bool reporting whether a
// has higher priority than b (leaves the heap first):
//
//	heaps.Min[int]   // smallest first (a < b)
//	heaps.Max[int]   // largest first  (a > b)
//
// The NewMin / NewMax constructors wire these up for any constraints.Ordered
// type; New takes an arbitrary comparator for custom or struct ordering.
//
// # Immutable vs Mutable Operations
//
// Immutable operations return a new heap and leave the receiver untouched:
//
//	bigger := pq.Push(7)         // returns a new heap
//	v, ok, rest := pq.Pop()      // returns the element and a new heap
//
// In-place operations mutate the receiver and carry the InPlace suffix:
//
//	pq.PushInPlace(7)            // modifies pq
//	v, ok := pq.PopInPlace()     // modifies pq
//
// # Draining in Priority Order
//
// The heap-array order is unspecified beyond the heap invariant. For priority
// order, drain the heap:
//
//	for v := range pq.Drain() {  // iterator, receiver untouched
//		fmt.Println(v)
//	}
//	sorted := pq.AsSortedSlice() // []T in priority order, receiver untouched
//
// # Thread Safety
//
//	balanced := heaps.NewConcurrentMin[int]()    // mutex
//	readHeavy := heaps.NewConcurrentRWMin[int]()  // read-write mutex
//
// Start with NewMin / NewMax / New and upgrade to a concurrent variant only
// when shared across goroutines.
//
// Callbacks passed to ForEach and the All iterator run after the lock is
// released, against a point-in-time snapshot taken under the lock. They may
// therefore safely re-enter the same heap (read it, or mutate it) without
// deadlocking.
package heaps
