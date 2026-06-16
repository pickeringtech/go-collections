// Package deques provides a generic double-ended queue (deque) backed by a ring
// buffer, with O(1) push and pop at both ends and an optional bounded mode.
//
// A deque generalises both the stack and the queue: you can add and remove
// elements at either end. The bounded variant doubles as a circular buffer,
// useful for fixed-size sliding windows, rate limiters, and recent-item caches.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/collections/deques"
//
//	d := deques.NewRingBuffer[int]()
//	d.PushBackInPlace(1)              // back:  [1]
//	d.PushBackInPlace(2)              // back:  [1 2]
//	d.PushFrontInPlace(0)            // front: [0 1 2]
//	front, _ := d.PopFrontInPlace()  // front: 0, deque [1 2]
//	back, _ := d.PopBackInPlace()    // back:  2, deque [1]
//	// Result: front == 0, back == 2
//
// # Why Use the Deques Package?
//
// Native approach — a slice gives O(n) removal from the front, and a fixed
// circular buffer means hand-rolling head/tail/wrap arithmetic every time:
//
//	queue := []int{1, 2, 3}
//	front := queue[0]
//	queue = queue[1:]               // O(n) over time, or leaks the backing array
//
// With this package:
//
//	d := deques.NewRingBuffer(1, 2, 3)
//	front, _, _ := d.PopFront()      // O(1), no manual index bookkeeping
//
// # Bounded vs Unbounded
//
// NewRingBuffer creates an unbounded deque that grows on demand; its Capacity is
// Unbounded (-1) and IsFull is always false.
//
// NewBoundedRingBuffer fixes the capacity and chooses an OverflowPolicy for what
// happens when a push arrives while full:
//
//	// Classic ring buffer: a push when full evicts the opposite end.
//	win := deques.NewBoundedRingBuffer[int](3, deques.OverwriteOldest)
//	win.PushBackInPlace(1)           // [1]
//	win.PushBackInPlace(2)           // [1 2]
//	win.PushBackInPlace(3)           // [1 2 3] (full)
//	win.PushBackInPlace(4)           // [2 3 4] — front 1 evicted, returns true
//
//	// Reject-when-full: a push when full is a no-op that reports false.
//	buf := deques.NewBoundedRingBuffer[int](3, deques.RejectWhenFull)
//	buf.PushBackInPlace(1)           // [1]
//	buf.PushBackInPlace(2)           // [1 2]
//	buf.PushBackInPlace(3)           // [1 2 3] (full)
//	ok := buf.PushBackInPlace(4)     // unchanged [1 2 3], ok == false
//
// # Immutable vs Mutable Operations
//
// Immutable operations return a new deque and never modify the receiver:
//
//	d2 := d.PushBack(element)         // returns Deque[T]
//	v, ok, d3 := d.PopFront()        // returns element, present?, Deque[T]
//
// In-place operations modify the receiver and report only a status:
//
//	accepted := d.PushBackInPlace(element)  // bool — false only on RejectWhenFull full
//	v, ok := d.PopFrontInPlace()            // element, present?
//	d.Clear()                               // empties the deque
//
// For a full bounded RejectWhenFull deque, the immutable PushFront/PushBack
// return an unchanged copy; use the in-place form when you need to know whether
// the element was accepted.
//
// # Iteration
//
// Deques are iterator-native. Range over values or index/value pairs, front to
// back or back to front:
//
//	for i, v := range d.All() { /* front to back */ }
//	for v := range d.Values() { /* front to back */ }
//	for i, v := range d.Backward() { /* back to front, i still counts from front */ }
//
// ForEach and ForEachWithIndex offer the same traversal via callbacks.
//
// # Available Implementations
//
//	RingBuffer              — plain, lock-free; not safe for concurrent use.
//	ConcurrentRingBuffer    — safe for concurrent use, guarded by a sync.Mutex.
//	ConcurrentRWRingBuffer  — safe for concurrent use, guarded by a sync.RWMutex
//	                          (reads take a read lock); favour for read-heavy use.
//
// Each has an unbounded constructor (New…RingBuffer) and a bounded one
// (NewBounded…RingBuffer) taking a capacity and OverflowPolicy.
//
// # Thread Safety
//
// The plain RingBuffer is not safe for concurrent use. Choose a concurrent
// variant when sharing a deque across goroutines:
//
//	// Balanced read/write workloads
//	work := deques.NewConcurrentRingBuffer[Task]()
//
//	// Read-heavy workloads
//	recent := deques.NewBoundedConcurrentRWRingBuffer[Event](100, deques.OverwriteOldest)
//
// Operating on a concurrent deque yields a concurrent deque of the same type:
// immutable operations return a new instance behind Deque[T], so thread-safe in
// means thread-safe out.
package deques
