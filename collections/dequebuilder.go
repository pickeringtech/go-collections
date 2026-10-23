package collections

import (
	"github.com/pickeringtech/go-collections/collections/deques"
	"github.com/pickeringtech/go-collections/slices"
)

// DequeBuilder fluently configures and constructs a Deque, allowing the
// concurrency characteristics, optional bounded ring-buffer mode, and initial
// values to be chosen before building.
type DequeBuilder[T any] struct {
	isConcurrent bool
	isRW         bool
	bounded      bool
	capacity     int
	policy       deques.OverflowPolicy
	values       []T
}

// NewDequeBuilder creates an empty DequeBuilder that builds a plain (non-concurrent), unbounded Deque by default.
func NewDequeBuilder[T any]() DequeBuilder[T] {
	return DequeBuilder[T]{}
}

// Concurrent marks the builder to produce a thread-safe Deque, returning the updated builder.
func (b DequeBuilder[T]) Concurrent() DequeBuilder[T] {
	b.isConcurrent = true
	return b
}

// RW selects the read-optimised (RWMutex-guarded) concurrent variant, returning the updated builder.
// It takes effect only together with Concurrent: without Concurrent, Build
// returns the plain, unsynchronised Deque whether or not RW was called.
func (b DequeBuilder[T]) RW() DequeBuilder[T] {
	b.isRW = true
	return b
}

// Bounded marks the builder to produce a bounded ring buffer with the given capacity and overflow policy, returning the updated builder.
func (b DequeBuilder[T]) Bounded(capacity int, policy deques.OverflowPolicy) DequeBuilder[T] {
	b.bounded = true
	b.capacity = capacity
	b.policy = policy
	return b
}

// Add appends the given values to those the built Deque will contain (front to back), returning the updated builder.
// The result is independent of the receiver: builders branched from a common
// base by separate Add calls never share or overwrite each other's values.
func (b DequeBuilder[T]) Add(values ...T) DequeBuilder[T] {
	// PushCopy never appends into the receiver's spare capacity, which a sibling
	// builder branched from the same base may also be using.
	b.values = slices.PushCopy(b.values, values...)
	return b
}

// Build constructs the Deque described by the builder.
func (b DequeBuilder[T]) Build() deques.Deque[T] {
	if b.bounded {
		if b.isConcurrent && b.isRW {
			return deques.NewBoundedConcurrentRWRingBuffer[T](b.capacity, b.policy, b.values...)
		}
		if b.isConcurrent {
			return deques.NewBoundedConcurrentRingBuffer[T](b.capacity, b.policy, b.values...)
		}
		return deques.NewBoundedRingBuffer[T](b.capacity, b.policy, b.values...)
	}
	if b.isConcurrent && b.isRW {
		return deques.NewConcurrentRWRingBuffer[T](b.values...)
	}
	if b.isConcurrent {
		return deques.NewConcurrentRingBuffer[T](b.values...)
	}
	return deques.NewRingBuffer[T](b.values...)
}
