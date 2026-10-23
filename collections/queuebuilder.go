package collections

import (
	"github.com/pickeringtech/go-collections/collections/lists"
	"github.com/pickeringtech/go-collections/slices"
)

// QueueBuilder fluently configures and constructs a Queue, allowing the
// concurrency characteristics and initial values to be chosen before building.
type QueueBuilder[T any] struct {
	isConcurrent bool
	isRW         bool
	values       []T
}

// NewQueueBuilder creates an empty QueueBuilder that builds a plain (non-concurrent) Queue by default.
func NewQueueBuilder[T any]() QueueBuilder[T] {
	return QueueBuilder[T]{}
}

// Concurrent marks the builder to produce a thread-safe Queue, returning the updated builder.
func (b QueueBuilder[T]) Concurrent() QueueBuilder[T] {
	b.isConcurrent = true
	return b
}

// RW selects the read-optimised (RWMutex-guarded) concurrent variant, returning the updated builder.
// It takes effect only together with Concurrent: without Concurrent, Build
// returns the plain, unsynchronised Queue whether or not RW was called.
func (b QueueBuilder[T]) RW() QueueBuilder[T] {
	b.isRW = true
	return b
}

// Add appends the given values to those the built Queue will contain, returning the updated builder.
// The result is independent of the receiver: builders branched from a common
// base by separate Add calls never share or overwrite each other's values.
func (b QueueBuilder[T]) Add(values ...T) QueueBuilder[T] {
	// PushCopy never appends into the receiver's spare capacity, which a sibling
	// builder branched from the same base may also be using.
	b.values = slices.PushCopy(b.values, values...)
	return b
}

// Build constructs the Queue described by the builder.
func (b QueueBuilder[T]) Build() lists.Queue[T] {
	if b.isConcurrent && b.isRW {
		return NewConcurrentRWQueue[T](b.values...)
	}
	if b.isConcurrent {
		return NewConcurrentQueue[T](b.values...)
	}
	return NewQueue[T](b.values...)
}
