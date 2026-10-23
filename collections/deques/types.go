package deques

// EachFunc is the callback invoked for each element during iteration.
type EachFunc[T any] func(element T)

// IndexedEachFunc is the callback invoked for each element during iteration,
// receiving both the element's index (counting from the front, 0-based) and its
// value.
type IndexedEachFunc[T any] func(idx int, element T)

// OverflowPolicy selects what a bounded deque does when an element is pushed
// while it is already at capacity. It has no effect on unbounded deques, which
// simply grow.
type OverflowPolicy int

const (
	// OverwriteOldest makes a push onto a full bounded deque succeed by evicting
	// the element at the opposite end: PushBack drops the current front, PushFront
	// drops the current back. This is the classic ring-buffer behaviour — a push
	// never blocks, and never fails unless the capacity is zero: a zero-capacity
	// deque has no slot to overwrite, so it rejects every push whatever its
	// policy.
	OverwriteOldest OverflowPolicy = iota

	// RejectWhenFull makes a push onto a full bounded deque a no-op: the deque is
	// left unchanged and the in-place push reports false. The caller decides what
	// to do with the rejected element.
	RejectWhenFull
)

// Unbounded is the Capacity reported by deques that have no fixed bound; such
// deques grow on demand and never report IsFull.
const Unbounded = -1
