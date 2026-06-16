package collections

import "github.com/pickeringtech/go-collections/collections/lists"

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

// RW marks the builder to favour concurrent reads (read-write locking), returning the updated builder.
func (b StackBuilder[T]) RW() StackBuilder[T] {
	b.isRW = true
	return b
}

// Add appends the given values to those the built Stack will contain, returning the updated builder.
func (b StackBuilder[T]) Add(values ...T) StackBuilder[T] {
	b.values = append(b.values, values...)
	return b
}

// Build constructs the Stack described by the builder.
func (b StackBuilder[T]) Build() lists.Stack[T] {
	//if b.isRW {
	//	return NewConcurrentRWList[T](b.values...)
	//}
	if b.isConcurrent {
		return NewConcurrentList[T](b.values...)
	}
	return NewList[T](b.values...)
}
