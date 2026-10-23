package collections

import (
	"github.com/pickeringtech/go-collections/collections/lists"
	"github.com/pickeringtech/go-collections/slices"
)

// ListBuilder fluently configures and constructs a List, allowing the
// concurrency characteristics and initial values to be chosen before building.
type ListBuilder[T any] struct {
	isConcurrent bool
	isRW         bool
	values       []T
}

// NewListBuilder creates an empty ListBuilder that builds a plain (non-concurrent) List by default.
func NewListBuilder[T any]() ListBuilder[T] {
	return ListBuilder[T]{}
}

// Concurrent marks the builder to produce a thread-safe List, returning the updated builder.
func (b ListBuilder[T]) Concurrent() ListBuilder[T] {
	b.isConcurrent = true
	return b
}

// RW selects the read-optimised (RWMutex-guarded) concurrent variant, returning the updated builder.
// It takes effect only together with Concurrent: without Concurrent, Build
// returns the plain, unsynchronised List whether or not RW was called.
func (b ListBuilder[T]) RW() ListBuilder[T] {
	b.isRW = true
	return b
}

// Add appends the given values to those the built List will contain, returning the updated builder.
// The result is independent of the receiver: builders branched from a common
// base by separate Add calls never share or overwrite each other's values.
func (b ListBuilder[T]) Add(values ...T) ListBuilder[T] {
	// PushCopy never appends into the receiver's spare capacity, which a sibling
	// builder branched from the same base may also be using.
	b.values = slices.PushCopy(b.values, values...)
	return b
}

// Build constructs the List described by the builder.
func (b ListBuilder[T]) Build() lists.List[T] {
	if b.isConcurrent && b.isRW {
		return NewConcurrentRWList[T](b.values...)
	}
	if b.isConcurrent {
		return NewConcurrentList[T](b.values...)
	}
	return NewList[T](b.values...)
}
