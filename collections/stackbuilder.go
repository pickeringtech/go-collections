package collections

import (
	"github.com/pickeringtech/go-collections/collections/lists"
	"github.com/pickeringtech/go-collections/slices"
)

// StackBuilder fluently configures and constructs a Stack, allowing the
// concurrency characteristics and initial values to be chosen before building.
type StackBuilder[T any] struct {
	isConcurrent bool
	isRW         bool
	values       []T
}

// NewStackBuilder creates an empty StackBuilder that builds a plain (non-concurrent) Stack by default.
func NewStackBuilder[T any]() StackBuilder[T] {
	return StackBuilder[T]{}
}

// Concurrent marks the builder to produce a thread-safe Stack, returning the updated builder.
func (b StackBuilder[T]) Concurrent() StackBuilder[T] {
	b.isConcurrent = true
	return b
}

// RW selects the read-optimised (RWMutex-guarded) concurrent variant, returning the updated builder.
// It takes effect only together with Concurrent: without Concurrent, Build
// returns the plain, unsynchronised Stack whether or not RW was called.
func (b StackBuilder[T]) RW() StackBuilder[T] {
	b.isRW = true
	return b
}

// Add appends the given values to those the built Stack will contain, returning the updated builder.
// The result is independent of the receiver: builders branched from a common
// base by separate Add calls never share or overwrite each other's values.
func (b StackBuilder[T]) Add(values ...T) StackBuilder[T] {
	// PushCopy never appends into the receiver's spare capacity, which a sibling
	// builder branched from the same base may also be using.
	b.values = slices.PushCopy(b.values, values...)
	return b
}

// Build constructs the Stack described by the builder.
func (b StackBuilder[T]) Build() lists.Stack[T] {
	if b.isConcurrent && b.isRW {
		return NewConcurrentRWStack[T](b.values...)
	}
	if b.isConcurrent {
		return NewConcurrentStack[T](b.values...)
	}
	return NewStack[T](b.values...)
}
